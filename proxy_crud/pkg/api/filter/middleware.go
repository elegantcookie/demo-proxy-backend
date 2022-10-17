package filter

import (
	"context"
	"fmt"
	"net/http"
	"proxy_crud/internal/proxy/apperror"
	"strings"
)

const (
	ascendingSort  = "ASC"
	descendingSort = "DESC"

	OptionsSortKey = "options_sort"
)

func Middleware(h http.HandlerFunc, defaultSortField, defaultSortOrder string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortBy := r.URL.Query().Get("sort_by")
		sortOrder := r.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = defaultSortField
		}
		if sortOrder == "" {
			sortOrder = defaultSortOrder
		} else {
			upperSortOrder := strings.ToUpper(sortOrder)
			if upperSortOrder != ascendingSort && upperSortOrder != descendingSort {
				w.WriteHeader(http.StatusBadRequest)
				w.Write(apperror.WrongSortingOptions.Marshal())
				return
			}
		}
		options := Options{
			Field: sortBy,
			Order: sortOrder,
		}
		fmt.Printf("%+v", options)

		value := context.WithValue(r.Context(), OptionsSortKey, options)
		r = r.WithContext(value)
		h(w, r)
	}
}

type Options struct {
	Field string
	Order string
}
