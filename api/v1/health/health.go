package health

import (
	"github.com/gin-gonic/gin"
	msghealth "github.com/romberli/go-template-web/pkg/message/health"
	"github.com/romberli/go-template-web/pkg/resp"
)

const (
	ZeroString = "0"
	pongString = `{"ping": "pong"}`
)

// @Tags health
// @Summary check health status
// @Accept	application/json
// @Produce application/json
// @Success 200 {string} string "0"
// @Router	/api/v1/health/status [get]
func Status(c *gin.Context) {
	resp.ResponseOKWithDebug(c, ZeroString, msghealth.InfoHealthStatus)
}

// @Tags health
// @Summary ping
// @Accept	application/json
// @Produce application/json
// @Success 200 {string} string "{"ping": "pong"}"
// @Router	/api/v1/health/ping [post]
func Ping(c *gin.Context) {
	resp.ResponseOK(c, pongString, msghealth.InfoHealthPing)
}
