package order

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	orderURL = "/order/:uuid"
)

type IHandler interface {
	Register(router *httprouter.Router)
}

type SHandler struct {
	logger     interface{} //*logging.Logger
	repository IHandler
}

func NewHandler(repository IHandler, logger interface{}) IHandler {
	return repository
}

func (h *SHandler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, orderURL, Middleware(h.GetById))
}

func (h *SHandler) GetById(w http.ResponseWriter, r *http.Request) error {
	// all, err := h.repository.FindAll(context.TODO())
	// if err != nil {
	// 	w.WriteHeader(400)
	// 	return err
	// }

	// allBytes, err := json.Marshal(all)
	// if err != nil {
	// 	return err
	// }

	// w.WriteHeader(http.StatusOK)
	// w.Write(allBytes)

	return nil
}
