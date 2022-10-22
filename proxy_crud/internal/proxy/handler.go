package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"proxy_crud/internal/apperror"
	"proxy_crud/internal/proxy/model"
	"proxy_crud/internal/proxy/service"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/logging"
	"proxy_crud/pkg/utils/parser"
)

type Handler struct {
	Logger       logging.Logger
	ProxyService service.Service
}

const (
	getAllUrl       = "/get/all"
	addProxiesByURL = "/add/many/url"
	updateProxyURL  = "/upd/id/:id"
	getByIDUrl      = "/get/id/:id"
	deleteAllURL    = "/delete/all"
)

func (h *Handler) Register(group *gin.RouterGroup) {
	group.Handle(http.MethodGet, getAllUrl, gin.WrapF(
		filter.Middleware(
			apperror.Middleware(h.GetProxies), "port", "ASC", 10)))
	group.Handle(http.MethodGet, getByIDUrl, apperror.GinMiddleware(h.GetProxyByID))
	group.Handle(http.MethodPut, updateProxyURL, apperror.GinMiddleware(h.UpdateProxy))
	group.Handle(http.MethodPost, addProxiesByURL, gin.WrapF(apperror.Middleware(h.AddProxies)))
	group.Handle(http.MethodDelete, deleteAllURL, gin.WrapF(apperror.Middleware(h.DeleteAll)))
}

func (h *Handler) AddProxies(w http.ResponseWriter, r *http.Request) error {
	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read body: %v", err)
	}
	str := string(bytes)
	proxyLines := parser.SplitText(str)
	var proxyDTOs []model.CreateProxyDTO
	//h.Logger.Println(proxyLines)
	for i := 0; i < len(proxyLines); i++ {
		pf, err := parser.ParseLine(proxyLines[i])
		if err != nil {
			h.Logger.Errorf("failed to parse line: %v", err)
			continue
		}

		proxyDTO, err := model.NewCreateProxyDTO(pf.Ip, pf.Port, pf.ExternalIp, pf.Country, pf.ProxyGroupID)
		if err != nil {
			h.Logger.Errorf("%v", err)
			continue
		}
		proxyDTOs = append(proxyDTOs, proxyDTO)
	}
	//h.Logger.Info(proxyDTOs)
	err = h.ProxyService.AddProxies(r.Context(), proxyDTOs)
	if err != nil {
		return fmt.Errorf("failed to add proxies: %v", err)
	}
	return nil
}

// GetProxies swaggo
// @Summary Returns data of all user_service
// @Accept json
// @Produce json
// @Tags Users
// @Success 200
// @Failure 400
// @Router /api/v1/user_service/get/all [get]
func (h *Handler) GetProxies(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var options filter.Options
	if opts, ok := r.Context().Value(filter.OptionsKey).(filter.Options); ok {
		options = opts
	}
	intFields := []string{"ping", "processing_status", "bl_check"}
	for _, fieldName := range intFields {
		val := r.URL.Query().Get(fieldName)
		options.ValidateIntAndAdd(fieldName, val, filter.OperatorEqual)
	}

	strFields := []string{"ip", "external_ip", "country"}
	for _, fieldName := range strFields {
		val := r.URL.Query().Get(fieldName)
		options.ValidateStringAndAdd(fieldName, val, filter.OperatorLike)
	}

	boolFields := []string{"active"}
	for _, fieldName := range boolFields {
		val := r.URL.Query().Get(fieldName)
		options.ValidateBoolAndAdd(fieldName, val, filter.OperatorEqual)
	}

	proxies, err := h.ProxyService.GetAll(r.Context(), options)
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

func (h *Handler) UpdateProxy(c *gin.Context) error {
	id := c.Param("id")
	_, err := h.ProxyService.GetById(c, id)
	if err != nil {
		return err
	}
	var dto model.UpdateProxyDTO
	err = json.NewDecoder(c.Request.Body).Decode(&dto)
	if err != nil {
		return err
	}
	err = h.ProxyService.Update(c, id, dto)
	if err != nil {
		return err
	}

	return nil

}

func (h *Handler) GetProxyByID(c *gin.Context) error {
	id := c.Param("id")
	proxy, err := h.ProxyService.GetById(c, id)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, proxy)
	return nil
}

func (h *Handler) DeleteAll(w http.ResponseWriter, r *http.Request) error {
	err := h.ProxyService.DeleteAll(r.Context())
	if err != nil {
		return fmt.Errorf("failed to delete all: %v", err)
	}
	return nil
}
