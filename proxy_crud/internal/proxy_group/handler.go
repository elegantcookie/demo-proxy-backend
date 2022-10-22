package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"proxy_crud/internal/apperror"
	"proxy_crud/internal/proxy_group/model"
	"proxy_crud/internal/proxy_group/service"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/logging"
)

type Handler struct {
	Logger            logging.Logger
	ProxyGroupService service.Service
}

const (
	addProxyGroupURL = "/add/"
	getAllUrl        = "/get/all"
	getByIDUrl       = "/get/id/:id"
	deleteAllURL     = "/delete/all"
)

func (h *Handler) Register(group *gin.RouterGroup) {
	group.Handle(http.MethodPost, addProxyGroupURL, gin.WrapF(apperror.Middleware(h.AddProxyGroup)))
	group.Handle(http.MethodGet, getAllUrl, gin.WrapF(
		filter.Middleware(
			apperror.Middleware(h.GetProxyGroups), "name", "ASC", 10)))
	group.Handle(http.MethodGet, getByIDUrl, apperror.GinMiddleware(h.GetProxyGroupByID))
	group.Handle(http.MethodDelete, deleteAllURL, gin.WrapF(apperror.Middleware(h.DeleteAll)))
}

func (h *Handler) AddProxyGroup(w http.ResponseWriter, r *http.Request) error {
	var dto model.CreateProxyGroupDTO
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("invalid JSON scheme. check swagger API")
	}
	err := h.ProxyGroupService.AddProxyGroup(r.Context(), dto)
	if err != nil {
		return err
	}
	w.WriteHeader(http.StatusNoContent)

	return nil
}

// GetProxyGroups swaggo
// @Summary Returns data of all user_service
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/v1/user_service/get/all [get]
func (h *Handler) GetProxyGroups(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var options filter.Options
	if opts, ok := r.Context().Value(filter.OptionsKey).(filter.Options); ok {
		options = opts
	}
	val := r.URL.Query().Get("name")
	options.ValidateStringAndAdd("name", val, filter.OperatorLike)

	proxyGroups, err := h.ProxyGroupService.GetAll(r.Context(), options)
	if err != nil {
		return err
	}

	bytes, err := json.Marshal(proxyGroups)
	if err != nil {
		return fmt.Errorf("failed to marshall proxy groups. error: %w", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
	return nil
}

func (h *Handler) GetProxyGroupByID(c *gin.Context) error {
	id := c.Param("id")
	proxy, err := h.ProxyGroupService.GetById(c, id)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, proxy)
	return nil
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	err := h.ProxyGroupService.DeleteAll(r.Context())
	if err != nil {
		return fmt.Errorf("failed to delete all: %v", err)
	}
	return nil
}
