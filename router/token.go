package router

import (
	"bytes"
	"io"
	"strings"

	"github.com/buger/jsonparser"
	"github.com/gin-gonic/gin"
	"github.com/pingcap/errors"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
	"github.com/romberli/go-util/middleware"
	"github.com/romberli/log"
	"github.com/spf13/viper"

	"github.com/romberli/db-operator/config"
	"github.com/romberli/db-operator/global"
	"github.com/romberli/db-operator/pkg/message"
	"github.com/romberli/db-operator/pkg/message/router"
	"github.com/romberli/db-operator/pkg/resp"
)

const (
	tokenTokenJSON   = "token"
	tokenPProfPath   = "/debug/pprof/"
	tokenStatusPath  = "/status"
	tokenHealthPath  = "/api/v1/health/"
	tokenSwaggerPath = "/swagger"
)

type TokenAuth struct {
	Database middleware.Pool
}

func NewTokenAuth(database middleware.Pool) *TokenAuth {
	return newTokenAuth(database)
}

func NewTokenAuthWithGlobal() *TokenAuth {
	return newTokenAuth(global.DBOMySQLPool)
}

func newTokenAuth(database middleware.Pool) *TokenAuth {
	return &TokenAuth{database}
}

func (ta *TokenAuth) Execute(command string, args ...interface{}) (middleware.Result, error) {
	conn, err := ta.Database.Get()
	if err != nil {
		return nil, err
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Errorf("router TokenAuth.Execute(): close database connection failed.\n%+v", err)
		}
	}()

	return conn.Execute(command, args...)
}

func (ta *TokenAuth) GetTokens() ([]string, error) {
	var tokens []string

	sql := `SELECT token FROM t_sys_token_info WHERE del_flag = 0;`
	log.Debugf("router TokenAuth.GetTokens() sql: \n%s", sql)

	result, err := ta.Execute(sql)
	if err != nil {
		return nil, err
	}

	for i := constant.ZeroInt; i < result.RowNumber(); i++ {
		token, err := result.GetString(i, constant.ZeroInt)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token)
	}

	return tokens, nil
}

func (ta *TokenAuth) GetHandlerFunc(tokens []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path

		if ta.IsSafePath(path) {
			return
		}

		// get data
		data, err := c.GetRawData()
		if err != nil {
			resp.ResponseNOK(c, message.ErrGetRawData, errors.Trace(err))
			c.Abort()
			return
		}

		// get alternative body
		bodyPath := viper.GetString(config.ServerRouterAlternativeBodyPathKey)
		if bodyPath != constant.EmptyString {
			bodyString, err := jsonparser.GetString(data, strings.Split(bodyPath, constant.DotString)...)
			if err == nil {
				// alternative body path exists
				data = []byte(bodyString)
			}
		}

		// set body back so that body can be read later in the router
		c.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		token, err := jsonparser.GetString(data, tokenTokenJSON)
		if err != nil {
			log.Errorf(string(data))
			resp.ResponseNOK(c, message.ErrFieldNotExistsOrWrongType, tokenTokenJSON)
			c.Abort()
			return
		}

		if !common.StringInSlice(tokens, token) {
			// not a valid token
			resp.ResponseNOK(c, router.ErrRouterValidateToken, token, c.ClientIP())
			c.Abort()
			return
		}
	}
}

func (ta *TokenAuth) IsSafePath(path string) bool {
	if strings.HasPrefix(path, tokenStatusPath) ||
		strings.HasPrefix(path, tokenHealthPath) ||
		strings.HasPrefix(path, tokenSwaggerPath) ||
		strings.HasPrefix(path, tokenPProfPath) {
		// do not check token for swagger
		return true
	}

	alternativeBasePath := viper.GetString(config.ServerRouterAlternativeBasePathKey)
	if alternativeBasePath != constant.EmptyString {
		if strings.HasPrefix(path, alternativeBasePath+tokenStatusPath) ||
			strings.HasPrefix(path, alternativeBasePath+tokenHealthPath) ||
			strings.HasPrefix(path, alternativeBasePath+tokenSwaggerPath) {
			// do not check token for swagger
			return true
		}
	}

	return false
}
