package metrics

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const HeartbeatURL = "api/proxy_crud/heartbeat"

type Handler struct {
}

func (h *Handler) Register(router gin.IRouter) {
	router.Handle(http.MethodGet, HeartbeatURL, gin.WrapF(h.Heartbeat))
}

// Heartbeat checks if the service is up
// @Summary Heartbeat metric
// @Tags Metrics
// @Success 204
// @Failure 400
// @Router /api/proxy_crud/heartbeat [get]
func (h *Handler) Heartbeat(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}
