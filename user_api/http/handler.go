package http

import (
	"github.com/obarbier/awesome-crypto/user_api/domain"
	"net/http"
)

type HandlerProperties struct {
	userService  domain.IUserService
	recoveryMode bool
}

func Handler(hp HandlerProperties) http.Handler {
	// TODO: add recovering resources / admin resources
	mux := http.NewServeMux()
	switch {
	case hp.recoveryMode:
	default:
		// TODO: adding auth wrapper
		mux.Handle("/v1/user", handleUserOperation(hp.userService))
	}
	return mux
}
