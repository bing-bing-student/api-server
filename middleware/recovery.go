package middleware

import (
	"blog/global"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"runtime/debug"
	"strings"
)

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if e, ok := err.(error); ok {
					var netErr net.Error
					if errors.As(e, &netErr) {
						if strings.Contains(strings.ToLower(e.Error()), "broken pipe") || strings.Contains(strings.ToLower(e.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					global.LOGGER.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					_ = c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					global.LOGGER.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					global.LOGGER.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
