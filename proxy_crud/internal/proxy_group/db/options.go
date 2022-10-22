package db

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"proxy_crud/internal/proxy_group/pgstorage"
	"proxy_crud/pkg/api/filter"
)

type Options struct {
	FilterOptions filter.IFOptions
	SortOptions   filter.SOptions
}

func (o *Options) GetOrderBy() string {
	return fmt.Sprintf("%s %s", o.SortOptions.Field, o.SortOptions.Order)
}

func (o *Options) GetLimit() uint64 {
	return uint64(o.FilterOptions.Limit())
}
func (o *Options) GetOffset() uint64 {
	return uint64(o.FilterOptions.Limit() * (o.FilterOptions.Page() - 1))
}

func (o *Options) MapOptions() pgstorage.Sqlizer {
	fields := o.FilterOptions.Fields()
	foLen := len(fields)
	some := sq.And{}
	for i := 0; i < foLen; i++ {
		some = sqlize(some, fields[i].Key, fields[i].Value, fields[i].Operator)
	}
	return some
}
