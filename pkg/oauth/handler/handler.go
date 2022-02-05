package handler

import (
	"fmt"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/obarbier/awesome-crypto/internal"
	"net/http"
)

type PropertiesHandler struct {
	Validate    *validator.Validate
	OauthServer *server.Server
}

func NewPropertiesHandler() PropertiesHandler {
	return PropertiesHandler{
		Validate:    common.NewValidator(),
		OauthServer: NewDefaultManageServer(),
	}
}
func Handler(hp PropertiesHandler) http.Handler {
	// TODO: add recovering resources / admin resources
	basePath := "/oauth"
	mux := http.NewServeMux()
	// TODO: adding auth wrapper
	mux.Handle(fmt.Sprintf("%s/%s", basePath, "v1/authorize"), handleOauthAuthorize(hp))
	mux.Handle(fmt.Sprintf("%s/%s", basePath, "v1/token"), handleOauthToken(hp))
	mux.Handle(fmt.Sprintf("%s/%s", basePath, "v1/login"), handleLogin(hp))
	mux.Handle(fmt.Sprintf("%s/%s", basePath, "v1/auth"), handleLoginAuth(hp))
	mux.Handle(fmt.Sprintf("%s/%s", basePath, "v1/test"), handleTestToken(hp))
	return mux
}
