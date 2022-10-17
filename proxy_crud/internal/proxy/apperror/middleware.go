package apperror

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error
type ginHandler func(c *gin.Context) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				//if errors.Is(err, ErrNotFound) {
				//	w.WriteHeader(http.StatusNotFound)
				//	w.Write(ErrNotFound.Marshal())
				//	return
				//}
				err := err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(err.Marshal())
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			w.Write(systemError(err.Error()).Marshal())
		}
	}
}

func GinMiddleware(h ginHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var appErr *AppError
		err := h(c)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					c.JSON(http.StatusNotFound, ErrNotFound)
					//w.Write(ErrNotFound.Marshal())
					return
				}
				err := err.(*AppError)
				c.JSON(http.StatusBadRequest, err)
				//c.Writer.Write(err.Marshal())
				return
			}
			c.JSON(http.StatusBadRequest, systemError(err.Error()))
			//c.Writer.Write(systemError(err.Error()).Marshal())
			return
		}
	}
}
