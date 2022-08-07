package db

import (
	"analytic-service/internal/db/mem"
	"analytic-service/internal/db/pg"
	"analytic-service/internal/logger"
	"analytic-service/internal/report"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewReport(ctx context.Context, source string, connstr string, logger *logger.Logger) report.Report {
	switch source {
	case "PG":
		obj := pg.ReportPG{Logger: logger}
		if connstr == "" {
			return &obj
		}

		dbpool, err := pgxpool.Connect(ctx, connstr)
		if err != nil {
			logger.Fatal(err.Error())
		}
		logger.Debug(fmt.Sprintf("Success connect to postgres: %s ", dbpool.Config().ConnString()))
		obj.Pool = dbpool
		return &obj
	case "MEM":
		obj := mem.ReportMEM{}
		return &obj
	default:
		return nil
	}
}

func Dispose(obj report.Report) {

}
