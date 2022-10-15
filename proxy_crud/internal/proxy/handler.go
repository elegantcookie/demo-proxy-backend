package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy_crud/internal/proxy/apperror"
	"proxy_crud/pkg/logging"
)

type Handler struct {
	Logger       logging.Logger
	ProxyService Service
}

const (
	getAllUrl = "/get/all"
)

func (h *Handler) Register(group *gin.RouterGroup) {
	group.Handle(http.MethodGet, getAllUrl, gin.WrapF(apperror.Middleware(h.GetProxies)))
}

// GetUsers swaggo
// @Summary Returns data of all user_service
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/v1/user_service/get/all [get]
func (h *Handler) GetProxies(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	proxies, err := h.ProxyService.GetAll(r.Context())
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(proxies)
	if err != nil {
		return fmt.Errorf("failed to marshall proxies. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return nil
}
