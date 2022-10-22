package db

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/jackc/pgx/v5"
	"proxy_crud/internal/apperror"
	"proxy_crud/internal/proxy/model"
	"proxy_crud/internal/proxy/pstorage"
	"proxy_crud/pkg/api/filter"
	"proxy_crud/pkg/client/postgresql"
	"proxy_crud/pkg/logging"
	"proxy_crud/pkg/utils/convertor"
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

func NewPSQLFilterOptions(options filter.Options) pstorage.IOptions {
	return &Options{
		FilterOptions: options.FilterOptions,
		SortOptions:   options.SortOptions,
	}
}

func queryToDebug(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func proxiesToArray(proxies []model.Proxy) []any {
	s := make([]any, len(proxies))
	for i, v := range proxies {
		s[i] = v
	}
	return s
}

func (d db) Insert(ctx context.Context, proxies []model.Proxy) error {
	columns := []string{
		"id", "ip", "port", "external_ip", "country", "open_ports", "active",
		"ping", "created_at", "checked_at", "valid_at", "bl_check", "processing_status", "proxy_group_id"}
	arr := proxiesToArray(proxies)
	ddproxies := convertor.StructsToArrays(arr)
	_, err := d.client.CopyFrom(
		ctx,
		pgx.Identifier{"proxy"},
		columns,
		pgx.CopyFromRows(ddproxies),
	)
	if err != nil {
		return err
	}
	return nil
}

func (d db) InsertOne(ctx context.Context, proxy model.Proxy) error {
	q := `
		INSERT INTO proxy (ip, port, external_ip, created_at, country, proxy_group_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))
	d.client.QueryRow(ctx, q, proxy.Ip, proxy.Port, proxy.ExternalIP, proxy.CreatedAt, proxy.Country)
	return nil
}

func (d db) FindById(ctx context.Context, id string) (model.Proxy, error) {
	q := `
			SELECT id, ip, port, external_ip, country, open_ports, active, ping, created_at, checked_at, valid_at, bl_check, processing_status, proxy_group_id
			FROM public.proxy WHERE id=$1
	`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))

	var pr model.Proxy

	err := d.client.QueryRow(ctx, q, id).Scan(
		&pr.ID, &pr.Ip, &pr.Port, &pr.ExternalIP, &pr.Country, &pr.OpenPorts,
		&pr.Active, &pr.Ping, &pr.CreatedAt, &pr.CheckedAt, &pr.ValidAt, &pr.BLCheck, &pr.ProcessingStatus, &pr.ProxyGroupID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pr, apperror.ErrNotFound
		}
		return pr, err
	}
	return pr, nil
}

func (d db) FindAll(ctx context.Context, options pstorage.IOptions) ([]model.Proxy, error) {
	qb := sq.Select("id, ip, port, external_ip, country, open_ports, active," +
		" ping, created_at, checked_at, valid_at, bl_check, processing_status, proxy_group_id").From("public.proxy")

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

	proxies := make([]model.Proxy, 0)

	for rows.Next() {
		var pr model.Proxy

		err := rows.Scan(
			&pr.ID, &pr.Ip, &pr.Port, &pr.ExternalIP, &pr.Country, &pr.OpenPorts,
			&pr.Active, &pr.Ping, &pr.CreatedAt, &pr.CheckedAt, &pr.ValidAt, &pr.BLCheck, &pr.ProcessingStatus, &pr.ProxyGroupID)
		if err != nil {
			return nil, err
		}

		proxies = append(proxies, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return proxies, nil
}

func (d db) Update(ctx context.Context, proxy model.Proxy) error {
	m := structs.Map(proxy)
	delete(m, "id")
	fmt.Printf("%+v", proxy)
	qb := sq.Update("public.proxy").SetMap(m).
		Where(sq.Eq{"id": proxy.ID}).
		PlaceholderFormat(sq.Dollar)

	sql, i, err := qb.ToSql()
	if err != nil {
		return err
	}
	d.logger.Tracef("SQL Query: %s", queryToDebug(sql))

	result, err := d.client.Exec(ctx, sql, i...)
	if err != nil {
		return err
	}
	d.logger.Info(result)
	return nil
}

func (d db) Delete(ctx context.Context, id string) error {
	return nil
}

func (d db) DeleteAll(ctx context.Context) error {
	q := `
			TRUNCATE public.proxy
	`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))

	result, err := d.client.Exec(ctx, q)
	if err != nil {
		return err
	}
	d.logger.Info(result)
	return nil
}

func NewStorage(client postgresql.Client, logger *logging.Logger) pstorage.Storage {
	return &db{
		client: client,
		logger: logger,
	}
}
