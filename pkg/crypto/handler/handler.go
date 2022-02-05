package handler

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	common "github.com/obarbier/awesome-crypto/internal"
	"github.com/obarbier/awesome-crypto/pkg/crypto/domain"
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
	basePath := "/crypto"
	mux := http.NewServeMux()
	switch {
	case hp.RecoveryMode:
	default:
		mux.Handle(fmt.Sprintf("%s/%s", basePath, "v1/key"), keyOperationHandler(hp))
	}
	jwtChecker := common.NewJwtMiddleWare()
	return jwtChecker.Handle(mux)
}
