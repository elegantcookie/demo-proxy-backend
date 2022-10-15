package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"proxy_crud/internal/proxy"
	"proxy_crud/pkg/client/postgresql"
	"proxy_crud/pkg/logging"
	"proxy_crud/pkg/utils/convertor"
	"strings"
)

type db struct {
	client postgresql.Client
	logger *logging.Logger
}

func queryToDebug(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func proxiesToArray(proxies []proxy.Proxy) []any {
	s := make([]any, len(proxies))
	for i, v := range proxies {
		s[i] = v
	}
	return s
}

func (d db) Insert(ctx context.Context, proxies []proxy.Proxy) error {
	columns := []string{
		"id", "ip", "port", "external_ip", "country", "open_ports", "active",
		"ping", "created_at", "checked_at", "valid_at", "bl_check", "processing_status"}
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

func (d db) InsertOne(ctx context.Context, proxy proxy.Proxy) error {
	q := `
		INSERT INTO proxy (ip, port, external_ip, created_at, country)
		VALUES ($1, $2, $3, $4, $5)
		`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))
	d.client.QueryRow(ctx, q, proxy.Ip, proxy.Port, proxy.ExternalIP, proxy.CreatedAt, proxy.Country)
	return nil
}

func (d db) FindById(ctx context.Context, id string) (proxy.Proxy, error) {
	q := `
			SELECT id, ip, port, external_ip, country, open_ports, active, ping, created_at, checked_at, valid_at, bl_check
			FROM public.proxy WHERE id=$1
	`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))

	var pr proxy.Proxy

	err := d.client.QueryRow(ctx, q, id).Scan(
		&pr.ID, &pr.Ip, &pr.Port, &pr.ExternalIP, &pr.Country, &pr.OpenPorts,
		&pr.Active, &pr.Ping, &pr.CreatedAt, &pr.CheckedAt, &pr.ValidAt, &pr.BLCheck)
	if err != nil {
		return pr, err
	}
	return pr, nil
}

func (d db) FindAll(ctx context.Context) ([]proxy.Proxy, error) {
	q := `
			SELECT id, ip, port, external_ip, country, open_ports, active, ping, created_at, checked_at, valid_at, bl_check
			FROM public.proxy
	`
	d.logger.Tracef("SQL Query: %s", queryToDebug(q))

	rows, err := d.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	proxies := make([]proxy.Proxy, 0)

	for rows.Next() {
		var pr proxy.Proxy

		err := rows.Scan(
			&pr.ID, &pr.Ip, &pr.Port, &pr.ExternalIP, &pr.Country, &pr.OpenPorts,
			&pr.Active, &pr.Ping, &pr.CreatedAt, &pr.CheckedAt, &pr.ValidAt, &pr.BLCheck)
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

func (d db) Update(ctx context.Context, proxy proxy.Proxy) error {
	return nil
}

func (d db) Delete(ctx context.Context, id string) error {
	return nil
}

func NewStorage(client postgresql.Client, logger *logging.Logger) proxy.Storage {
	return &db{
		client: client,
		logger: logger,
	}
}
