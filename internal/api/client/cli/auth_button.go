package cli

import (
	"context"
	"fmt"
	"github.com/rivo/tview"
)

type Register struct {
	Login    string
	Password string
}

func (c *CLI) register(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {

	var reg Register

	form := tview.NewForm().
		AddInputField("Login", "", 20, nil, func(text string) {
			reg.Login = text
		}).
		AddInputField("Password", "", 20, nil, func(text string) {
			reg.Password = text
		}).
		AddButton("Save", func() {
			ok := c.registerAPI(ctx, reg)

			fmt.Println(ok)

			if ok {
				pages.AddPage("AuthUser", c.authUser(ctx, app, pages), true, false)
				pages.SwitchToPage("AuthUser")
			} else {
				c.errorsRegister(ctx, app, pages)
			}

		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Зарегистрироваться").SetTitleAlign(tview.AlignLeft)

	return form
}

func (c *CLI) registerAPI(ctx context.Context, reg Register) bool {
	err := c.auth.RegisterUser(ctx, c.conn, reg.Login, reg.Password)

	if err != nil {
		return false
	}

	return true
}

func (c *CLI) authUser(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {

	var reg Register

	form := tview.NewForm().
		AddInputField("Login", "", 20, nil, func(text string) {
			reg.Login = text
		}).
		AddInputField("Password", "", 20, nil, func(text string) {
			reg.Password = text
		}).
		AddButton("Save", func() {
			token, err := c.auth.AuthUser(ctx, c.conn, reg.Login, reg.Password)
			if err != nil {
				c.errorsAuth(ctx, app, pages)
			} else {
				c.token = token
				pages.AddPage("Buttons_data", c.buttonsData(ctx, app, pages), true, false)
				pages.SwitchToPage("Buttons_data")
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Авторизоваться на сервере").SetTitleAlign(tview.AlignCenter)

	return form
}

func (c *CLI) authUserOffline(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {

	var reg Register

	form := tview.NewForm().
		AddInputField("Login", "", 20, nil, func(text string) {
			reg.Login = text
		}).
		AddInputField("Password", "", 20, nil, func(text string) {
			reg.Password = text
		}).
		AddButton("Save", func() {
			token, err := c.auth.AuthUserOffLine(ctx, reg.Login, reg.Password)
			if err != nil {
				c.errorsAuth(ctx, app, pages)
			} else {
				c.token = token
				pages.AddPage("Buttons_data", c.buttonsData(ctx, app, pages), true, false)
				pages.SwitchToPage("Buttons_data")
			}
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Авторизоваться на сервере").SetTitleAlign(tview.AlignCenter)

	return form
}
