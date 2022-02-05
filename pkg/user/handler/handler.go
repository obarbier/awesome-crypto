package handler

import (
	"fmt"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-playground/validator/v10"
	"github.com/obarbier/awesome-crypto/internal"
	"github.com/obarbier/awesome-crypto/pkg/user/domain"
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
	}
}
func Handler(hp PropertiesHandler) http.Handler {
	// TODO: add recovering resources / admin resources
	basePath := "/user"
	mux := http.NewServeMux()
	switch {
	case hp.RecoveryMode:
	default:
		// TODO: adding auth wrapper
		mux.Handle(fmt.Sprintf("%s/%s", basePath, "v1/user"), handleUserOperation(hp))
	}
	jwtChecker := common.NewJwtMiddleWare()
	return jwtChecker.Handle(mux)
}
