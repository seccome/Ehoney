<p align="center">

<img  width="600" src="https://github.com/kataras/go-sessions/raw/master/logo_900_273_bg_white.png">
<br/><br/>

 <a href="https://travis-ci.org/kataras/go-sessions"><img src="https://img.shields.io/travis/kataras/go-sessions.svg?style=flat-square" alt="Build Status"></a>
 <a href="https://github.com/kataras/go-sessions/blob/master/LICENSE"><img src="https://img.shields.io/badge/%20license-MIT%20%20License%20-E91E63.svg?style=flat-square" alt="License"></a>
 <a href="https://github.com/kataras/go-sessions/releases"><img src="https://img.shields.io/badge/%20release%20-%20v3.3.0-blue.svg?style=flat-square" alt="Releases"></a>
 <a href="#documentation"><img src="https://img.shields.io/badge/%20docs-reference-5272B4.svg?style=flat-square" alt="Read me docs"></a>
 <a href="https://kataras.rocket.chat/channel/go-sessions"><img src="https://img.shields.io/badge/%20community-chat-00BCD4.svg?style=flat-square" alt="Build Status"></a>
 <a href="https://golang.org"><img src="https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square" alt="Built with GoLang"></a>
 <a href="#"><img src="https://img.shields.io/badge/platform-All-yellow.svg?style=flat-square" alt="Platforms"></a>
<br/>

<a href="#features" >Fast</a> http sessions manager for Go.<br/>
Simple <a href ="#outline">API</a>, while providing robust set of features such as immutability, expiration time (can be shifted), [databases](sessiondb) like badger and redis as back-end storage.<br/>

</p>

Quick view
-----------

```go
import "github.com/kataras/go-sessions/v3"

sess := sessions.Start(http.ResponseWriter, *http.Request)
sess.
  ID() string
  Get(string) interface{}
  HasFlash() bool
  GetFlash(string) interface{}
  GetFlashString(string) string
  GetString(key string) string
  GetInt(key string) (int, error)
  GetInt64(key string) (int64, error)
  GetFloat32(key string) (float32, error)
  GetFloat64(key string) (float64, error)
  GetBoolean(key string) (bool, error)
  GetAll() map[string]interface{}
  GetFlashes() map[string]interface{}
  VisitAll(cb func(k string, v interface{}))
  Set(string, interface{})
  SetImmutable(key string, value interface{})
  SetFlash(string, interface{})
  Delete(string)
  Clear()
  ClearFlashes()
```

Installation
------------

The only requirement is the [Go Programming Language](https://golang.org/dl), at least 1.14.

```bash
$ go get github.com/kataras/go-sessions/v3
```

**go.mod**

```sh
module your_app

go 1.14

require (
	github.com/kataras/go-sessions/v3 v3.3.0
)
```

Features
------------

- Focus on simplicity and performance.
- Flash messages.
- Supports any type of [external database](_examples/database).
- Works with both [net/http](https://golang.org/pkg/net/http/) and [valyala/fasthttp](https://github.com/valyala/fasthttp).

Documentation
------------

Take a look at the [./examples](https://github.com/kataras/go-sessions/tree/master/_examples) folder.

- [Overview](_examples/overview/main.go)
- [Standalone](_examples/standalone/main.go)
- [Fasthttp](_examples/fasthttp/main.go)
- [Secure Cookie](_examples/securecookie/main.go)
- [Flash Messages](_examples/flash-messages/main.go)
- [Databases](_examples/database)
	* [Redis](_examples/database/redis/main.go)

Outline
------------

```go
// Start starts the session for the particular net/http request
Start(w http.ResponseWriter,r *http.Request) Session
// ShiftExpiration move the expire date of a session to a new date
// by using session default timeout configuration.
ShiftExpiration(w http.ResponseWriter, r *http.Request)
// UpdateExpiration change expire date of a session to a new date
// by using timeout value passed by `expires` receiver.
UpdateExpiration(w http.ResponseWriter, r *http.Request, expires time.Duration)
// Destroy kills the net/http session and remove the associated cookie
Destroy(w http.ResponseWriter,r  *http.Request)

// Start starts the session for the particular valyala/fasthttp request
StartFasthttp(ctx *fasthttp.RequestCtx) Session
// ShiftExpirationFasthttp move the expire date of a session to a new date
// by using session default timeout configuration.
ShiftExpirationFasthttp(ctx *fasthttp.RequestCtx)
// UpdateExpirationFasthttp change expire date of a session to a new date
// by using timeout value passed by `expires` receiver.
UpdateExpirationFasthttp(ctx *fasthttp.RequestCtx, expires time.Duration)
// Destroy kills the valyala/fasthttp session and remove the associated cookie
DestroyFasthttp(ctx *fasthttp.RequestCtx)

// DestroyByID removes the session entry
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
//
// It's safe to use it even if you are not sure if a session with that id exists.
// Works for both net/http & fasthttp
DestroyByID(string)
// DestroyAll removes all sessions
// from the server-side memory (and database if registered).
// Client's session cookie will still exist but it will be reseted on the next request.
// Works for both net/http & fasthttp
DestroyAll()

// UseDatabase ,optionally, adds a session database to the manager's provider,
// a session db doesn't have write access
// see https://github.com/kataras/go-sessions/tree/master/sessiondb
UseDatabase(Database)
```

### Configuration

```go
// Config is the configuration for sessions. Please read it before using sessions.
Config struct {
	// Cookie string, the session's client cookie name, for example: "mysessionid"
	//
	// Defaults to "irissessionid".
	Cookie string

	// CookieSecureTLS set to true if server is running over TLS
	// and you need the session's cookie "Secure" field to be setted true.
	//
	// Note: The user should fill the Decode configuation field in order for this to work.
	// Recommendation: You don't need this to be setted to true, just fill the Encode and Decode fields
	// with a third-party library like secure cookie, example is provided at the _examples folder.
	//
	// Defaults to false.
	CookieSecureTLS bool

	// AllowReclaim will allow to
	// Destroy and Start a session in the same request handler.
	// All it does is that it removes the cookie for both `Request` and `ResponseWriter` while `Destroy`
	// or add a new cookie to `Request` while `Start`.
	//
	// Defaults to false.
	AllowReclaim bool

	// Encode the cookie value if not nil.
	// Should accept as first argument the cookie name (config.Cookie)
	//         as second argument the server's generated session id.
	// Should return the new session id, if error the session id setted to empty which is invalid.
	//
	// Note: Errors are not printed, so you have to know what you're doing,
	// and remember: if you use AES it only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	//
	// Defaults to nil.
	Encode func(cookieName string, value interface{}) (string, error)
	// Decode the cookie value if not nil.
	// Should accept as first argument the cookie name (config.Cookie)
	//               as second second accepts the client's cookie value (the encoded session id).
	// Should return an error if decode operation failed.
	//
	// Note: Errors are not printed, so you have to know what you're doing,
	// and remember: if you use AES it only supports key sizes of 16, 24 or 32 bytes.
	// You either need to provide exactly that amount or you derive the key from what you type in.
	//
	// Defaults to nil.
	Decode func(cookieName string, cookieValue string, v interface{}) error

	// Encoding same as Encode and Decode but receives a single instance which
	// completes the "CookieEncoder" interface, `Encode` and `Decode` functions.
	//
	// Defaults to nil.
	Encoding Encoding

	// Expires the duration of which the cookie must expires (created_time.Add(Expires)).
	// If you want to delete the cookie when the browser closes, set it to -1.
	//
	// 0 means no expire, (24 years)
	// -1 means when browser closes
	// > 0 is the time.Duration which the session cookies should expire.
	//
	// Defaults to infinitive/unlimited life duration(0).
	Expires time.Duration

	// SessionIDGenerator should returns a random session id.
	// By default we will use a uuid impl package to generate
	// that, but developers can change that with simple assignment.
	SessionIDGenerator func() string

	// DisableSubdomainPersistence set it to true in order dissallow your subdomains to have access to the session cookie
	//
	// Defaults to false.
	DisableSubdomainPersistence bool
}
```


Usage NET/HTTP
------------


`Start` returns a `Session`, **Session outline**

```go
ID() string
Get(string) interface{}
HasFlash() bool
GetFlash(string) interface{}
GetString(key string) string
GetFlashString(string) string
GetInt(key string) (int, error)
GetInt64(key string) (int64, error)
GetFloat32(key string) (float32, error)
GetFloat64(key string) (float64, error)
GetBoolean(key string) (bool, error)
GetAll() map[string]interface{}
GetFlashes() map[string]interface{}
VisitAll(cb func(k string, v interface{}))
Set(string, interface{})
SetImmutable(key string, value interface{})
SetFlash(string, interface{})
Delete(string)
Clear()
ClearFlashes()
```

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kataras/go-sessions/v3"
)

type businessModel struct {
	Name string
}

func main() {
	app := http.NewServeMux()
	sess := sessions.New(sessions.Config{
		// Cookie string, the session's client cookie name, for example: "mysessionid"
		//
		// Defaults to "gosessionid"
		Cookie: "mysessionid",
		// it's time.Duration, from the time cookie is created, how long it can be alive?
		// 0 means no expire.
		// -1 means expire when browser closes
		// or set a value, like 2 hours:
		Expires: time.Hour * 2,
		// if you want to invalid cookies on different subdomains
		// of the same host, then enable it
		DisableSubdomainPersistence: false,
		// want to be crazy safe? Take a look at the "securecookie" example folder.
	})

	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("You should navigate to the /set, /get, /delete, /clear,/destroy instead")))
	})
	app.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {

		//set session values.
		s := sess.Start(w, r)
		s.Set("name", "iris")

		//test if setted here
		w.Write([]byte(fmt.Sprintf("All ok session setted to: %s", s.GetString("name"))))

		// Set will set the value as-it-is,
		// if it's a slice or map
		// you will be able to change it on .Get directly!
		// Keep note that I don't recommend saving big data neither slices or maps on a session
		// but if you really need it then use the `SetImmutable` instead of `Set`.
		// Use `SetImmutable` consistently, it's slower.
		// Read more about muttable and immutable go types: https://stackoverflow.com/a/8021081
	})

	app.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		// get a specific value, as string, if no found returns just an empty string
		name := sess.Start(w, r).GetString("name")

		w.Write([]byte(fmt.Sprintf("The name on the /set was: %s", name)))
	})

	app.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		// delete a specific key
		sess.Start(w, r).Delete("name")
	})

	app.HandleFunc("/clear", func(w http.ResponseWriter, r *http.Request) {
		// removes all entries
		sess.Start(w, r).Clear()
	})

	app.HandleFunc("/update", func(w http.ResponseWriter, r *http.Request) {
		// updates expire date
		sess.ShiftExpiration(w, r)
	})

	app.HandleFunc("/destroy", func(w http.ResponseWriter, r *http.Request) {

		//destroy, removes the entire session data and cookie
		sess.Destroy(w, r)
	})
	// Note about Destroy:
	//
	// You can destroy a session outside of a handler too, using the:
	// mySessions.DestroyByID
	// mySessions.DestroyAll

	// remember: slices and maps are muttable by-design
	// The `SetImmutable` makes sure that they will be stored and received
	// as immutable, so you can't change them directly by mistake.
	//
	// Use `SetImmutable` consistently, it's slower than `Set`.
	// Read more about muttable and immutable go types: https://stackoverflow.com/a/8021081
	app.HandleFunc("/set_immutable", func(w http.ResponseWriter, r *http.Request) {
		business := []businessModel{{Name: "Edward"}, {Name: "value 2"}}
		s := sess.Start(w, r)
		s.SetImmutable("businessEdit", business)
		businessGet := s.Get("businessEdit").([]businessModel)

		// try to change it, if we used `Set` instead of `SetImmutable` this
		// change will affect the underline array of the session's value "businessEdit", but now it will not.
		businessGet[0].Name = "Gabriel"

	})

	app.HandleFunc("/get_immutable", func(w http.ResponseWriter, r *http.Request) {
		valSlice := sess.Start(w, r).Get("businessEdit")
		if valSlice == nil {
			w.Header().Set("Content-Type", "text/html; charset=UTF-8")
			w.Write([]byte("please navigate to the <a href='/set_immutable'>/set_immutable</a> first"))
			return
		}

		firstModel := valSlice.([]businessModel)[0]
		// businessGet[0].Name is equal to Edward initially
		if firstModel.Name != "Edward" {
			panic("Report this as a bug, immutable data cannot be changed from the caller without re-SetImmutable")
		}

		w.Write([]byte(fmt.Sprintf("[]businessModel[0].Name remains: %s", firstModel.Name)))

		// the name should remains "Edward"
	})

	http.ListenAndServe(":8080", app)
}
```


Usage FASTHTTP
------------

`StartFasthttp` returns the same object as `Start`, `Session`.

```go
ID() string
Get(string) interface{}
HasFlash() bool
GetFlash(string) interface{}
GetString(key string) string
GetFlashString(string) string
GetInt(key string) (int, error)
GetInt64(key string) (int64, error)
GetFloat32(key string) (float32, error)
GetFloat64(key string) (float64, error)
GetBoolean(key string) (bool, error)
GetAll() map[string]interface{}
GetFlashes() map[string]interface{}
VisitAll(cb func(k string, v interface{}))
Set(string, interface{})
SetImmutable(key string, value interface{})
SetFlash(string, interface{})
Delete(string)
Clear()
ClearFlashes()
```

We have only one simple example because the API is the same, the returned session is the same for both net/http and valyala/fasthttp.

Just append the word "Fasthttp", the rest of the API remains as it's with net/http.

`Start` for net/http, `StartFasthttp` for valyala/fasthttp.
`ShiftExpiration` for net/http, `ShiftExpirationFasthttp` for valyala/fasthttp.
`UpdateExpiration` for net/http, `UpdateExpirationFasthttp` for valyala/fasthttp.
`Destroy` for net/http, `DestroyFasthttp` for valyala/fasthttp.

```go
package main

import (
	"fmt"

	"github.com/kataras/go-sessions/v3"
	"github.com/valyala/fasthttp"
)

func main() {
// set some values to the session
setHandler := func(reqCtx *fasthttp.RequestCtx) {
	values := map[string]interface{}{
		"Name":   "go-sessions",
		"Days":   "1",
		"Secret": "dsads£2132215£%%Ssdsa",
	}

	sess := sessions.StartFasthttp(reqCtx) // init the session
	// sessions.StartFasthttp returns the, same, Session interface we saw before too

	for k, v := range values {
		sess.Set(k, v) // fill session, set each of the key-value pair
	}
	reqCtx.WriteString("Session saved, go to /get to view the results")
}

// get the values from the session
getHandler := func(reqCtx *fasthttp.RequestCtx) {
	sess := sessions.StartFasthttp(reqCtx) // init the session
	sessValues := sess.GetAll()            // get all values from this session

	reqCtx.WriteString(fmt.Sprintf("%#v", sessValues))
}

// clear all values from the session
clearHandler := func(reqCtx *fasthttp.RequestCtx) {
	sess := sessions.StartFasthttp(reqCtx)
	sess.Clear()
}

// destroys the session, clears the values and removes the server-side entry and client-side sessionid cookie
destroyHandler := func(reqCtx *fasthttp.RequestCtx) {
	sessions.DestroyFasthttp(reqCtx)
}

fmt.Println("Open a browser tab and navigate to the localhost:8080/set")
fasthttp.ListenAndServe(":8080", func(reqCtx *fasthttp.RequestCtx) {
	path := string(reqCtx.Path())

	if path == "/set" {
		setHandler(reqCtx)
	} else if path == "/get" {
		getHandler(reqCtx)
	} else if path == "/clear" {
		clearHandler(reqCtx)
	} else if path == "/destroy" {
		destroyHandler(reqCtx)
	} else {
		reqCtx.WriteString("Please navigate to /set or /get or /clear or /destroy")
	}
})
}
```

FAQ
------------

If you'd like to discuss this package, or ask questions about it, feel free to

 * Explore [these questions](https://github.com/kataras/go-sessions/issues?go-sessions=label%3Aquestion).
 * Post an issue or  idea [here](https://github.com/kataras/go-sessions/issues).
 * Navigate to the [Chat][Chat].



Versioning
------------

Current: **v3.3.0**

Read more about Semantic Versioning 2.0.0

 - http://semver.org/
 - https://en.wikipedia.org/wiki/Software_versioning
 - https://wiki.debian.org/UpstreamGuide#Releases_and_Versions

People
------------

The author of go-sessions is [@kataras](https://github.com/kataras).

Contributing
------------

If you are interested in contributing to the go-sessions project, please make a PR.

License
------------

This project is licensed under the MIT License.

License can be found [here](LICENSE).

[Travis Widget]: https://img.shields.io/travis/kataras/go-sessions.svg?style=flat-square
[Travis]: http://travis-ci.org/kataras/go-sessions
[License Widget]: https://img.shields.io/badge/license-MIT%20%20License%20-E91E63.svg?style=flat-square
[License]: https://github.com/kataras/go-sessions/blob/master/LICENSE
[Release Widget]: https://img.shields.io/badge/release-v3.3.0-blue.svg?style=flat-square
[Release]: https://github.com/kataras/go-sessions/releases
[Chat Widget]: https://img.shields.io/badge/community-chat-00BCD4.svg?style=flat-square
[Chat]: https://kataras.rocket.chat/channel/go-sessions
[ChatMain]: https://kataras.rocket.chat/channel/go-sessions
[ChatAlternative]: https://gitter.im/kataras/go-sessions
[Report Widget]: https://img.shields.io/badge/report%20card-A%2B-F44336.svg?style=flat-square
[Report]: http://goreportcard.com/report/kataras/go-sessions
[Documentation Widget]: https://img.shields.io/badge/docs-reference-5272B4.svg?style=flat-square
[Documentation]: https://godoc.org/github.com/kataras/go-sessions
[Language Widget]: https://img.shields.io/badge/powered_by-Go-3362c2.svg?style=flat-square
[Language]: http://golang.org
[Platform Widget]: https://img.shields.io/badge/platform-Any--OS-yellow.svg?style=flat-square
