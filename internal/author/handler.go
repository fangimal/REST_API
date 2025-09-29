package author

import (
	"REST_API/internal/apperror"
	"REST_API/internal/handlers"
	"REST_API/pkg/logging"
	"context"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

//var _ handlers.Handler = &handler{}

const (
	authorsURL = "/authors"
	authorURL  = "/authors/:uuid"
)

type handler struct {
	logger     *logging.Logger
	repository Repository
}

func NewHandler(repository Repository, logger *logging.Logger) handlers.Handler {
	return &handler{
		repository: repository,
		logger:     logger,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, authorsURL, apperror.Middleware(h.GetList))
	//router.HandlerFunc(http.MethodPost, authorsURL, apperror.Middleware(h.Createauthor))
	//router.HandlerFunc(http.MethodGet, authorURL, apperror.Middleware(h.GetauthorByUUID))
	//router.HandlerFunc(http.MethodPut, authorURL, apperror.Middleware(h.Updateauthor))
	//router.HandlerFunc(http.MethodPatch, authorURL, apperror.Middleware(h.PartiallyUpdateauthor))
	//router.HandlerFunc(http.MethodDelete, authorURL, apperror.Middleware(h.Deleteauthor))
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	all, err := h.repository.FindAll(context.TODO())
	if err != nil {
		w.WriteHeader(400)
		return err
	}

	allBytes, err := json.Marshal(all)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(allBytes)

	return nil
}

//func (h *handler) Createauthor(w http.ResponseWriter, r *http.Request) error {
//	return fmt.Errorf("this is API error")
//}
//
//func (h *handler) GetauthorByUUID(w http.ResponseWriter, r *http.Request) error {
//	return apperror.NewAppError(nil, "test", "test", "t13")
//}
//
//func (h *handler) Updateauthor(w http.ResponseWriter, r *http.Request) error {
//	w.WriteHeader(204)
//	w.Write([]byte("this is update authors"))
//
//	return nil
//}
//
//func (h *handler) PartiallyUpdateauthor(w http.ResponseWriter, r *http.Request) error {
//	w.WriteHeader(204)
//	w.Write([]byte("this is PartiallyUpdate authors"))
//
//	return nil
//}
//
//func (h *handler) Deleteauthor(w http.ResponseWriter, r *http.Request) error {
//	w.WriteHeader(204)
//	w.Write([]byte("this is delete authors"))
//
//	return nil
//}
