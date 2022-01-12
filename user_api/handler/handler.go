package handler

import (
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/obarbier/awesome-crypto/common"
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"net/http"
)

type PropertiesHandler struct {
	UserService  domain.IUserService
	RecoveryMode bool
	Validate     *validator.Validate
	OauthServer  *server.Server
}

func NewPropertiesHandler(service domain.IUserService, recoveryMode bool) PropertiesHandler {
	return PropertiesHandler{
		UserService:  service,
		RecoveryMode: recoveryMode,
		Validate:     common.NewValidator(),
		OauthServer:  NewDefaultManageServer(),
	}
}
func Handler(hp PropertiesHandler) http.Handler {
	// TODO: add recovering resources / admin resources
	mux := http.NewServeMux()
	switch {
	case hp.RecoveryMode:
	default:
		// TODO: adding auth wrapper
		mux.Handle("/v1/user", handleUserOperation(hp))
		mux.Handle("/v1/oauth/authorize", handleOauthAuthorize(hp))
		mux.Handle("/v1/oauth/token", handleOauthToken(hp))
		mux.Handle("/v1/login", handleLogin(hp))
		mux.Handle("/v1/auth", handleLoginAuth(hp))
		mux.Handle("/v1/test", handleTestToken(hp))
	}
	return mux
}
