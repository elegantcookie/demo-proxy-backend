package filter

import (
	"context"
	"net/http"
	"proxy_crud/internal/proxy/apperror"
	"strconv"
	"strings"
)

const (
	ascendingSort  = "ASC"
	descendingSort = "DESC"

	OptionsKey = "options"
)

func getSortOptions(r *http.Request, defaultSortField, defaultSortOrder string) (opts SOptions, err *apperror.AppError) {
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
			return opts, apperror.WrongSortingOptions
		}
	}
	opts = SOptions{
		Field: sortBy,
		Order: sortOrder,
	}
	return opts, nil
}

func getFilterOptions(r *http.Request, defaultLimit int) (opts *FOptions, err *apperror.AppError) {
	limitFromQuery := r.URL.Query().Get("limit")

	var limit int
	if limitFromQuery == "" {
		limit = defaultLimit
	} else {
		var parseErr error
		if limit, parseErr = strconv.Atoi(limitFromQuery); parseErr != nil {
			return opts, apperror.WrongFilterOptions
		}
	}
	opts = &FOptions{
		apply:  true,
		limit:  limit,
		fields: make([]Field, 0),
	}
	return opts, nil
}

func Middleware(h http.HandlerFunc, defaultSortField, defaultSortOrder string, defaultLimit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortOptions, err := getSortOptions(r, defaultSortField, defaultSortOrder)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(err.Marshal())
			return
		}
		filterOptions, err := getFilterOptions(r, defaultLimit)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(err.Marshal())
			return
		}

		//fmt.Printf("%+v", sortOptions)
		opts := Options{
			SortOptions:   sortOptions,
			FilterOptions: filterOptions,
		}
		ctx := context.WithValue(r.Context(), OptionsKey, opts)
		//ctx = context.WithValue(ctx, OptionsFilterKey, filterOptions)
		r = r.WithContext(ctx)
		h(w, r)
	}
}
