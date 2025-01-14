package cli

import (
	"context"
	"database/sql"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"goph-keeper/internal/api/client/handlers"
	"goph-keeper/internal/api/client/repositories/auth"
	"goph-keeper/internal/api/client/repositories/health"
	"log/slog"
)

type getService interface {
	GetAllData(ctx context.Context, token, tableName string) (*sql.Rows, error)
}

type deletedService interface {
	DeletedData(ctx context.Context, token, tableName string, id int) error
}

type CLI struct {
	log         *slog.Logger
	auth        *auth.Handlers
	save        *handlers.Handler
	deleted     deletedService
	getAll      getService
	healthCheck *health.Handler
	conn        *grpc.ClientConn
	token       string
}

func NewCLI(log *slog.Logger,
	auth *auth.Handlers,
	save *handlers.Handler,
	deleted deletedService,
	get getService,
	healthCheck *health.Handler,
	conn *grpc.ClientConn) *CLI {
	return &CLI{
		log:         log,
		auth:        auth,
		save:        save,
		deleted:     deleted,
		getAll:      get,
		healthCheck: healthCheck,
		conn:        conn,
	}
}

func (c *CLI) RunCLI(ctx context.Context) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	if err := c.healthCheck.Health(ctx, c.conn); err != nil {
		pages.AddPage("ButtonsOffLine", c.buttonsOffline(ctx, app, pages), true, true)
	} else {
		pages.AddPage("Buttons", c.buttonsStart(ctx, app, pages), true, true)
	}

	app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
