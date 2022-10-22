package db

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"proxy_crud/internal/apperror"
	"proxy_crud/internal/proxy_group/model"
	"proxy_crud/internal/proxy_group/pgstorage"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/client/postgresql"
	"proxy_crud/pkg/logging"
	"strings"
)

type db struct {
	client postgresql.Client
	logger *logging.Logger
}

// TODO: refactor
func sqlize(s sq.And, key string, value any, operator string) sq.And {
	switch operator {
	case filter.OperatorLike:
		s = append(s, sq.Like{key: value})
	case filter.OperatorGreaterThan:
		s = append(s, sq.Gt{key: value})
	case filter.OperatorGreaterThanEq:
		s = append(s, sq.GtOrEq{key: value})
	case filter.OperatorLowerThan:
		s = append(s, sq.Lt{key: value})
	case filter.OperatorLowerThanEq:
		s = append(s, sq.LtOrEq{key: value})
	case filter.OperatorEqual:
		s = append(s, sq.Eq{key: value})
	case filter.OperatorNotEqual:
		s = append(s, sq.NotEq{key: value})
	}

	return s
}

func NewPSQLFilterOptions(options filter.Options) pgstorage.IOptions {
	return &Options{
		FilterOptions: options.FilterOptions,
		SortOptions:   options.SortOptions,
	}
}

func queryToDebug(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func proxiesToArray(proxies []model.ProxyGroup) []any {
	s := make([]any, len(proxies))
	for i, v := range proxies {
		s[i] = v
	}
	return s
}

func (d db) InsertOne(ctx context.Context, pg model.ProxyGroup) error {
	q := `
		INSERT INTO public.proxy_group (id, name)
		VALUES ($1, $2)
		`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))
	d.client.QueryRow(ctx, q, pg.ID, pg.Name)
	return nil
}

func (d db) FindById(ctx context.Context, id string) (model.ProxyGroup, error) {
	q := `
			SELECT id, name
			FROM public.proxy_group WHERE id=$1
	`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))

	var pg model.ProxyGroup

	err := d.client.QueryRow(ctx, q, id).Scan(
		&pg.ID, &pg.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pg, apperror.ErrNotFound
		}
		return pg, err
	}
	return pg, nil
}

func (d db) FindAll(ctx context.Context, options pgstorage.IOptions) ([]model.ProxyGroup, error) {
	qb := sq.Select("id, name").From("public.proxy_group")

	if options != nil {
		sql := options.MapOptions()
		qb = qb.Where(sql).
			OrderBy(options.GetOrderBy()).
			PlaceholderFormat(sq.Dollar).
			Limit(options.GetLimit()).
			Offset(options.GetOffset())
	}
	sql, i, err := qb.ToSql()
	fmt.Println(i)

	if err != nil {
		return nil, err
	}

	d.logger.Tracef("SQL Query: %s", queryToDebug(sql))
	rows, err := d.client.Query(ctx, sql, i...)
	if err != nil {
		return nil, err
	}

	proxyGroups := make([]model.ProxyGroup, 0)

	for rows.Next() {
		var pg model.ProxyGroup

		err := rows.Scan(
			&pg.ID, &pg.Name)
		if err != nil {
			return nil, err
		}

		proxyGroups = append(proxyGroups, pg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return proxyGroups, nil
}

func (d db) Update(ctx context.Context, proxy model.ProxyGroup) error {
	return nil
}

func (d db) Delete(ctx context.Context, id string) error {
	return nil
}

func (d db) DeleteAll(ctx context.Context) error {
	q := `
			TRUNCATE public.proxy_group
	`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))

	result, err := d.client.Exec(ctx, q)
	if err != nil {
		return err
	}
	d.logger.Info(result)
	return nil
}

func NewStorage(client postgresql.Client, logger *logging.Logger) pgstorage.Storage {
	return &db{
		client: client,
		logger: logger,
	}
}
