package router

import (
	errs "errors"
	"net"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pingcap/errors"
	"github.com/romberli/go-template-web/pkg/message"
	"github.com/romberli/go-template-web/pkg/resp"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/log"
	"go.uber.org/zap/zapcore"
)

const (
	defaultBrokenPipeMessage = "broken pipe"
	defaultResetMessage      = "connection reset by peer"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		if log.GetLevel() != zapcore.DebugLevel {
			return
		}

		// start time
		start := time.Now()
		c.Next()
		// end time
		end := time.Now()
		// execution tme
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		log.Debugf("path: %s, clientIP: %s, method: %s, statusCode: %d, latency: %s", path, clientIP, method, statusCode, latency)
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			err := recover()
			if err != nil {
				var brokenPipe bool

				ne, ok := err.(*net.OpError)
				if ok {
					var se *os.SyscallError
					if errs.As(ne, &se) {
						if strings.Contains(strings.ToLower(se.Error()), defaultBrokenPipeMessage) || strings.Contains(strings.ToLower(se.Error()), defaultResetMessage) {
							brokenPipe = true
						}
					}
				}

				request, dumpErr := httputil.DumpRequest(c.Request, true)
				if dumpErr != nil {
					log.Errorf("gin Recovery: dump http request failed. error:\n%+v", errors.Trace(dumpErr))
				} else {
					log.Errorf("gin Recovery: http request:\n%s", string(request))
				}

				errWithStack := errors.New(err.(error).Error())
				if brokenPipe {
					log.Errorf(constant.LogWithStackString, message.NewMessage(message.ErrGinRecovery, errWithStack))
					c.Abort()
					return
				}

				resp.ResponseNOK(c, message.ErrGinRecovery, errWithStack)
				c.Abort()
				return
			}
		}()

		c.Next()
	}
}
