package db

import (
	"context"
	"log"
	"time"

	"github.com/uptrace/bun"
)

// QueryLogger implements bun.QueryHook to log executed queries.
type QueryLogger struct{}

// BeforeQuery is called before a query is executed.
func (h *QueryLogger) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

// AfterQuery is called after a query is executed.
func (h *QueryLogger) AfterQuery(_ context.Context, event *bun.QueryEvent) {
	// You can replace the standard `log` with your project's logger.
	log.Printf(
		"[Repository] SQL Query: %s | Duration: %s | Error: %v",
		event.Query,
		time.Since(event.StartTime),
		event.Err,
	)
}
