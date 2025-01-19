package cli

import (
	"context"
	"database/sql"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"goph-keeper/internal/api/client/handlers/auth"
	"goph-keeper/internal/api/client/handlers/save"
	"log/slog"
)

type getService interface {
	GetAllData(ctx context.Context, token, tableName string) (*sql.Rows, error)
}

type CLI struct {
	log    *slog.Logger
	auth   *auth.Handlers
	save   *save.Handler
	getAll getService
	conn   *grpc.ClientConn
	token  string
}

func NewCLI(log *slog.Logger, auth *auth.Handlers, save *save.Handler, get getService, conn *grpc.ClientConn) *CLI {
	return &CLI{
		log:    log,
		auth:   auth,
		save:   save,
		getAll: get,
		conn:   conn,
	}
}

func (c *CLI) RunCLI(ctx context.Context) {
	app := tview.NewApplication()
	pages := tview.NewPages()

	pages.AddPage("Buttons", c.buttonsStart(ctx, app, pages), true, true)

	app.SetRoot(pages, true).EnableMouse(true).EnablePaste(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
