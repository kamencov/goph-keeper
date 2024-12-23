package cli

import (
	"context"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"goph-keeper/internal/api/client/handlers/auth"
	"goph-keeper/internal/api/client/handlers/save"
	"log/slog"
)

type CLI struct {
	log   *slog.Logger
	auth  *auth.Handlers
	save  *save.Handler
	conn  *grpc.ClientConn
	token string
}

func NewCLI(log *slog.Logger, auth *auth.Handlers, save *save.Handler, conn *grpc.ClientConn) *CLI {
	return &CLI{
		log:  log,
		auth: auth,
		save: save,
		conn: conn,
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
