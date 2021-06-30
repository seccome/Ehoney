// Package sessions provides sessions support for net/http and valyala/fasthttp
// unique with auto-GC, register unlimited number of databases to Load and Update/Save the sessions in external server or to an external (no/or/and sql) database
// Usage net/http:
// // init a new sessions manager( if you use only one web framework inside your app then you can use the package-level functions like: sessions.Start/sessions.Destroy)
// manager := sessions.New(sessions.Config{})
// // start a session for a particular client
// manager.Start(http.ResponseWriter, *http.Request)
//
// // destroy a session from the server and client,
//  // don't call it on each handler, only on the handler you want the client to 'logout' or something like this:
// manager.Destroy(http.ResponseWriter, *http.Request)
//
//
// Usage valyala/fasthttp:
// // init a new sessions manager( if you use only one web framework inside your app then you can use the package-level functions like: sessions.Start/sessions.Destroy)
// manager := sessions.New(sessions.Config{})
// // start a session for a particular client
// manager.StartFasthttp(*fasthttp.RequestCtx)
//
// // destroy a session from the server and client,
//  // don't call it on each handler, only on the handler you want the client to 'logout' or something like this:
// manager.DestroyFasthttp(*fasthttp.Request)
//
// Note that, now, you can use both fasthttp and net/http within the same sessions manager(.New) instance!
// So now, you can share sessions between a net/http app and valyala/fasthttp app
package sessions

import (
	"net/http"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	// Version current semantic version string of the go-sessions package.
	Version = "3.3.0"
)

// A Sessions manager should be responsible to Start a sesion, based
// on a Context, which should return
// a compatible Session interface, type. If the external session manager
// doesn't qualifies, then the user should code the rest of the functions with empty implementation.
//
// Sessions should be responsible to Destroy a session based
// on the Context.
type Sessions struct {
	config   Config
	provider *provider
}

// Default instance of the sessions, used for package-level functions.
var Default = New(Config{}.Validate())

// New returns the fast, feature-rich sessions manager.
func New(cfg Config) *Sessions {
	return &Sessions{
		config:   cfg.Validate(),
		provider: newProvider(),
	}
}

// UseDatabase adds a session database to the manager's provider.
func UseDatabase(db Database) {
	Default.UseDatabase(db)
}

// UseDatabase adds a session database to the manager's provider.
func (s *Sessions) UseDatabase(db Database) {
	s.provider.RegisterDatabase(db)
}

// updateCookie gains the ability of updating the session browser cookie to any method which wants to update it
func (s *Sessions) updateCookie(w http.ResponseWriter, r *http.Request, sid string, expires time.Duration) {
	cookie := &http.Cookie{}

	// The RFC makes no mention of encoding url value, so here I think to encode both sessionid key and the value using the safe(to put and to use as cookie) url-encoding
	cookie.Name = s.config.Cookie

	cookie.Value = sid
	cookie.Path = "/"
	if !s.config.DisableSubdomainPersistence {

		requestDomain := r.URL.Host
		if portIdx := strings.IndexByte(requestDomain, ':'); portIdx > 0 {
			requestDomain = requestDomain[0:portIdx]
		}
		if IsValidCookieDomain(requestDomain) {

			// RFC2109, we allow level 1 subdomains, but no further
			// if we have localhost.com , we want the localhost.cos.
			// so if we have something like: mysubdomain.localhost.com we want the localhost here
			// if we have mysubsubdomain.mysubdomain.localhost.com we want the .mysubdomain.localhost.com here
			// slow things here, especially the 'replace' but this is a good and understable( I hope) way to get the be able to set cookies from subdomains & domain with 1-level limit
			if dotIdx := strings.LastIndexByte(requestDomain, '.'); dotIdx > 0 {
				// is mysubdomain.localhost.com || mysubsubdomain.mysubdomain.localhost.com
				s := requestDomain[0:dotIdx] // set mysubdomain.localhost || mysubsubdomain.mysubdomain.localhost
				if secondDotIdx := strings.LastIndexByte(s, '.'); secondDotIdx > 0 {
					//is mysubdomain.localhost ||  mysubsubdomain.mysubdomain.localhost
					s = s[secondDotIdx+1:] // set to localhost || mysubdomain.localhost
				}
				// replace the s with the requestDomain before the domain's siffux
				subdomainSuff := strings.LastIndexByte(requestDomain, '.')
				if subdomainSuff > len(s) { // if it is actual exists as subdomain suffix
					requestDomain = strings.Replace(requestDomain, requestDomain[0:subdomainSuff], s, 1) // set to localhost.com || mysubdomain.localhost.com
				}
			}
			// finally set the .localhost.com (for(1-level) || .mysubdomain.localhost.com (for 2-level subdomain allow)
			cookie.Domain = "." + requestDomain // . to allow persistence
		}
	}

	cookie.Domain = formatCookieDomain(r.URL.Host, s.config.DisableSubdomainPersistence)
	cookie.HttpOnly = true
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	if expires >= 0 {
		if expires == 0 { // unlimited life
			cookie.Expires = CookieExpireUnlimited
		} else { // > 0
			cookie.Expires = time.Now().Add(expires)
		}
		cookie.MaxAge = int(cookie.Expires.Sub(time.Now()).Seconds())
	}

	// set the cookie to secure if this is a tls wrapped request
	// and the configuration allows it.
	if r.TLS != nil && s.config.CookieSecureTLS {
		cookie.Secure = true
	}

	// encode the session id cookie client value right before send it.
	cookie.Value = s.encodeCookieValue(cookie.Value)
	AddCookie(w, r, cookie, s.config.AllowReclaim)
}

// Start starts the session for the particular request.
func Start(w http.ResponseWriter, r *http.Request) *Session {
	return Default.Start(w, r)
}

// Start starts the session for the particular request.
func (s *Sessions) Start(w http.ResponseWriter, r *http.Request) *Session {
	cookieValue := s.decodeCookieValue(GetCookie(r, s.config.Cookie))

	if cookieValue == "" { // cookie doesn't exists, let's generate a session and add set a cookie
		sid := s.config.SessionIDGenerator()

		sess := s.provider.Init(sid, s.config.Expires)
		sess.isNew = s.provider.db.Len(sid) == 0

		s.updateCookie(w, r, sid, s.config.Expires)

		return sess
	}

	sess := s.provider.Read(cookieValue, s.config.Expires)

	return sess
}

func (s *Sessions) updateCookieFasthttp(ctx *fasthttp.RequestCtx, sid string, expires time.Duration) {
	cookie := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(cookie)

	// The RFC makes no mention of encoding url value, so here I think to encode both sessionid key and the value using the safe(to put and to use as cookie) url-encoding
	cookie.SetKey(s.config.Cookie)

	cookie.SetValue(sid)
	cookie.SetPath("/")
	cookie.SetDomain(formatCookieDomain(string(ctx.Host()), s.config.DisableSubdomainPersistence))

	cookie.SetHTTPOnly(true)
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
	// MaxAge>0 means Max-Age attribute present and given in seconds
	if expires >= 0 {
		if expires == 0 { // unlimited life
			cookie.SetExpire(CookieExpireUnlimited)
		} else { // > 0
			cookie.SetExpire(time.Now().Add(expires))
		}
	}

	// set the cookie to secure if this is a tls wrapped request
	// and the configuration allows it.

	if ctx.IsTLS() && s.config.CookieSecureTLS {
		cookie.SetSecure(true)
	}

	// encode the session id cookie client value right before send it.
	cookie.SetValue(s.encodeCookieValue(string(cookie.Value())))
	AddCookieFasthttp(ctx, cookie)
}

// StartFasthttp starts the session for the particular request.
func StartFasthttp(ctx *fasthttp.RequestCtx) *Session {
	return Default.StartFasthttp(ctx)
}

// StartFasthttp starts the session for the particular request.
func (s *Sessions) StartFasthttp(ctx *fasthttp.RequestCtx) *Session {
	cookieValue := s.decodeCookieValue(GetCookieFasthttp(ctx, s.config.Cookie))

	if cookieValue == "" { // cookie doesn't exists, let's generate a session and add set a cookie
		sid := s.config.SessionIDGenerator()

		sess := s.provider.Init(sid, s.config.Expires)
		sess.isNew = s.provider.db.Len(sid) == 0

		s.updateCookieFasthttp(ctx, sid, s.config.Expires)

		return sess
	}

	sess := s.provider.Read(cookieValue, s.config.Expires)

	return sess
}

// ShiftExpiration move the expire date of a session to a new date
// by using session default timeout configuration.
func ShiftExpiration(w http.ResponseWriter, r *http.Request) {
	Default.ShiftExpiration(w, r)
}

// ShiftExpiration move the expire date of a session to a new date
// by using session default timeout configuration.
func (s *Sessions) ShiftExpiration(w http.ResponseWriter, r *http.Request) {
	s.UpdateExpiration(w, r, s.config.Expires)
}

// ShiftExpirationFasthttp move the expire date of a session to a new date
// by using session default timeout configuration.
func ShiftExpirationFasthttp(ctx *fasthttp.RequestCtx) {
	Default.ShiftExpirationFasthttp(ctx)
}

// ShiftExpirationFasthttp move the expire date of a session to a new date
// by using session default timeout configuration.
func (s *Sessions) ShiftExpirationFasthttp(ctx *fasthttp.RequestCtx) {
	s.UpdateExpirationFasthttp(ctx, s.config.Expires)
}

// UpdateExpiration change expire date of a session to a new date
// by using timeout value passed by `expires` receiver.
func UpdateExpiration(w http.ResponseWriter, r *http.Request, expires time.Duration) {
	Default.UpdateExpiration(w, r, expires)
}

// UpdateExpiration change expire date of a session to a new date
// by using timeout value passed by `expires` receiver.
// It will return `ErrNotFound` when trying to update expiration on a non-existence or not valid session entry.
// It will return `ErrNotImplemented` if a database is used and it does not support this feature, yet.
func (s *Sessions) UpdateExpiration(w http.ResponseWriter, r *http.Request, expires time.Duration) error {
	cookieValue := s.decodeCookieValue(GetCookie(r, s.config.Cookie))
	if cookieValue == "" {
		return ErrNotFound
	}

	// we should also allow it to expire when the browser closed
	err := s.provider.UpdateExpiration(cookieValue, expires)
	if err == nil || expires == -1 {
		s.updateCookie(w, r, cookieValue, expires)
	}

	return err
}

// UpdateExpirationFasthttp change expire date of a session to a new date
// by using timeout value passed by `expires` receiver.
func UpdateExpirationFasthttp(ctx *fasthttp.RequestCtx, expires time.Duration) {
	Default.UpdateExpirationFasthttp(ctx, expires)
}

// UpdateExpirationFasthttp change expire date of a session to a new date
// by using timeout value passed by `expires` receiver.
func (s *Sessions) UpdateExpirationFasthttp(ctx *fasthttp.RequestCtx, expires time.Duration) error {
	cookieValue := s.decodeCookieValue(GetCookieFasthttp(ctx, s.config.Cookie))
	if cookieValue == "" {
		return ErrNotFound
	}

	// we should also allow it to expire when the browser closed
	err := s.provider.UpdateExpiration(cookieValue, expires)
	if err == nil || expires == -1 {
		s.updateCookieFasthttp(ctx, cookieValue, expires)
	}

	return err
}

func (s *Sessions) destroy(cookieValue string) {
	// decode the client's cookie value in order to find the server's session id
	// to destroy the session data.
	cookieValue = s.decodeCookieValue(cookieValue)
	if cookieValue == "" { // nothing to destroy
		return
	}

	s.provider.Destroy(cookieValue)
}

// DestroyListener is the form of a destroy listener.
// Look `OnDestroy` for more.
type DestroyListener func(sid string)

// OnDestroy registers one or more destroy listeners.
// A destroy listener is fired when a session has been removed entirely from the server (the entry) and client-side (the cookie).
// Note that if a destroy listener is blocking, then the session manager will delay respectfully,
// use a goroutine inside the listener to avoid that behavior.
func (s *Sessions) OnDestroy(listeners ...DestroyListener) {
	for _, ln := range listeners {
		s.provider.registerDestroyListener(ln)
	}
}

// OnDestroy registers one or more destroy listeners.
// A destroy listener is fired when a session has been removed entirely from the server (the entry) and client-side (the cookie).
// Note that if a destroy listener is blocking, then the session manager will delay respectfully,
// use a goroutine inside the listener to avoid that behavior.
func OnDestroy(listeners ...DestroyListener) {
	Default.OnDestroy(listeners...)
}

// Destroy remove the session data and remove the associated cookie.
func Destroy(w http.ResponseWriter, r *http.Request) {
	Default.Destroy(w, r)
}

// Destroy remove the session data and remove the associated cookie.
func (s *Sessions) Destroy(w http.ResponseWriter, r *http.Request) {
	cookieValue := GetCookie(r, s.config.Cookie)
	s.destroy(cookieValue)
	RemoveCookie(w, r, s.config)
}

// DestroyFasthttp remove the session data and remove the associated cookie.
func DestroyFasthttp(ctx *fasthttp.RequestCtx) {
	Default.DestroyFasthttp(ctx)
}

// DestroyFasthttp remove the session data and remove the associated cookie.
func (s *Sessions) DestroyFasthttp(ctx *fasthttp.RequestCtx) {
	cookieValue := GetCookieFasthttp(ctx, s.config.Cookie)
	s.destroy(cookieValue)
	RemoveCookieFasthttp(ctx, s.config)
}

// DestroyByID removes the session entry
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
//
// It's safe to use it even if you are not sure if a session with that id exists.
//
// Note: the sid should be the original one (i.e: fetched by a store )
// it's not decoded.
func DestroyByID(sid string) {
	Default.DestroyByID(sid)
}

// DestroyByID removes the session entry
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
//
// It's safe to use it even if you are not sure if a session with that id exists.
//
// Note: the sid should be the original one (i.e: fetched by a store )
// it's not decoded.
func (s *Sessions) DestroyByID(sid string) {
	s.provider.Destroy(sid)
}

// DestroyAll removes all sessions
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
func DestroyAll() {
	Default.DestroyAll()
}

// DestroyAll removes all sessions
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
func (s *Sessions) DestroyAll() {
	s.provider.DestroyAll()
}

// let's keep these funcs simple, we can do it with two lines but we may add more things in the future.
func (s *Sessions) decodeCookieValue(cookieValue string) string {
	if cookieValue == "" {
		return ""
	}

	var cookieValueDecoded *string

	if decode := s.config.Decode; decode != nil {
		err := decode(s.config.Cookie, cookieValue, &cookieValueDecoded)
		if err == nil {
			cookieValue = *cookieValueDecoded
		} else {
			cookieValue = ""
		}
	}

	return cookieValue
}

func (s *Sessions) encodeCookieValue(cookieValue string) string {
	if encode := s.config.Encode; encode != nil {
		newVal, err := encode(s.config.Cookie, cookieValue)
		if err == nil {
			cookieValue = newVal
		} else {
			cookieValue = ""
		}
	}

	return cookieValue
}
