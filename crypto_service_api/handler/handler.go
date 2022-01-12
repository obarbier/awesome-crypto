package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/obarbier/awesome-crypto/common"
	"github.com/obarbier/awesome-crypto/crypto_service_api/domain"
	"net/http"
)

type PropertiesHandler struct {
	UserService  domain.IKeyService
	RecoveryMode bool
	Validate     *validator.Validate
}

func NewPropertiesHandler(service domain.IKeyService, recoveryMode bool) PropertiesHandler {
	return PropertiesHandler{
		UserService:  service,
		RecoveryMode: recoveryMode,
		Validate:     common.NewValidator(),
	}
}
func Handler(hp PropertiesHandler) http.Handler {
	mux := http.NewServeMux()
	switch {
	case hp.RecoveryMode:
	default:
		mux.Handle("/v1/key", keyOperationHandler(hp))
	}
	jwtChecker := common.NewJwtMiddleWare()
	return jwtChecker.Handle(mux)
}
