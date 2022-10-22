package filter

import (
	"context"
	"net/http"
	"strconv"
	"strings"
)

const (
	ascendingSort  = "ASC"
	descendingSort = "DESC"

	OptionsKey = "options"
)

// TODO: refactor

func getSortOptions(r *http.Request, defaultSortField, defaultSortOrder string) (opts SOptions) {
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
			sortOrder = defaultSortOrder
		}
	}
	opts = SOptions{
		Field: sortBy,
		Order: sortOrder,
	}
	return opts
}

func getFilterOptions(r *http.Request, defaultLimit, defaultPage int) (opts *FOptions) {
	limitFromQuery := r.URL.Query().Get("limit")
	pageFromQuery := r.URL.Query().Get("page")

	var limit int
	if limitFromQuery == "" {
		limit = defaultLimit
	} else {
		var parseErr error
		if limit, parseErr = strconv.Atoi(limitFromQuery); parseErr != nil || limit < 1 {
			limit = defaultLimit
		}
	}
	var page int
	if pageFromQuery == "" {
		page = defaultPage
	} else {
		var parseErr error
		if page, parseErr = strconv.Atoi(pageFromQuery); parseErr != nil || page < 1 {
			page = defaultPage
		}
	}

	opts = &FOptions{
		apply:  true,
		limit:  limit,
		page:   page,
		fields: make([]Field, 0),
	}
	return opts
}

func Middleware(h http.HandlerFunc, defaultSortField, defaultSortOrder string, defaultLimit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		sortOptions := getSortOptions(r, defaultSortField, defaultSortOrder)
		filterOptions := getFilterOptions(r, defaultLimit, 1)

		opts := Options{
			SortOptions:   sortOptions,
			FilterOptions: filterOptions,
		}
		ctx := context.WithValue(r.Context(), OptionsKey, opts)
		r = r.WithContext(ctx)
		h(w, r)
	}
}
