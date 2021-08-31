package https

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"strconv"
)

func TLSHandler(port int) gin.HandlerFunc {
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect : true,
			SSLHost     : ":" + strconv.Itoa(port),
		})
		err := secureMiddleware.Process(c.Writer, c.Request)
		if err != nil {
			return
		}
		c.Next()
	}
}
