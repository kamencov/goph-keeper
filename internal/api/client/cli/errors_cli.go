package cli

import (
	"context"
	"github.com/rivo/tview"
)

func (c *CLI) errorsRegister(ctx context.Context, app *tview.Application, pages *tview.Pages) {

	model := tview.NewModal()
	model.SetText("Данный пользователь уже есть или введены неверные данные\n" +
		"Выберите действие:\n" +
		"1. Register: Повторно пройти регистрацию\n" +
		"2. Cancel: Закрыть приложение\n").
		AddButtons([]string{"Register", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Register" {
				pages.AddPage("Register", c.register(ctx, app, pages), true, false)
				pages.SwitchToPage("Register")
			} else {
				app.Stop()
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("ErrorsRegister", model, true, true)
}

func (c *CLI) errorsAuth(ctx context.Context, app *tview.Application, pages *tview.Pages) {
	model := tview.NewModal()
	model.SetText("Что-то пошло не так\n" +
		"Выберите действие:\n" +
		"1. Auth: Повторно пройти авторизацию\n" +
		"2. Cancel: Закрыть приложение\n").
		AddButtons([]string{"Auth", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Auth" {
				pages.AddPage("Auth", c.authUser(ctx, app, pages), true, false)
				pages.SwitchToPage("Auth")
			} else {
				app.Stop()
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("ErrorsAuth", model, true, true)
}

func (c *CLI) errorsAuthOffLine(ctx context.Context, app *tview.Application, pages *tview.Pages) {
	model := tview.NewModal()
	model.SetText("Что-то пошло не так\n" +
		"Выберите действие:\n" +
		"1. Auth: Повторно пройти авторизацию\n" +
		"2. Cancel: Закрыть приложение\n").
		AddButtons([]string{"Auth", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Auth" {
				pages.AddPage("AuthOffLine", c.authUserOffline(ctx, app, pages), true, false)
				pages.SwitchToPage("AuthOffLine")
			} else {
				app.Stop()
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("ErrorsAuthOffLine", model, true, true)
}

func (c *CLI) errorsSaveCredentials(ctx context.Context, app *tview.Application, pages *tview.Pages) {
	model := tview.NewModal()
	model.SetText("Что-то пошло не так\n" +
		"Выберите действие:\n" +
		"1. Save: Повторно сохранить данные\n" +
		"2. Cancel: Закрыть приложение\n").
		AddButtons([]string{"Save", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				pages.AddPage("Credentials", c.credentials(ctx, app, pages), true, false)
				pages.SwitchToPage("Credentials")
			} else {
				app.Stop()
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("ErrorsSaveCredentials", model, true, true)
}

func (c *CLI) errorsSaveCards(ctx context.Context, app *tview.Application, pages *tview.Pages) {
	model := tview.NewModal()
	model.SetText("Что-то пошло не так\n" +
		"Выберите действие:\n" +
		"1. Save: Повторно сохранить данные\n" +
		"2. Cancel: Закрыть приложение\n").
		AddButtons([]string{"Save", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				pages.AddPage("Card", c.cardButton(ctx, app, pages), true, false)
				pages.SwitchToPage("Card")
			} else {
				app.Stop()
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("ErrorsSaveCards", model, true, true)
}

func (c *CLI) errorsSaveBinary(ctx context.Context, app *tview.Application, pages *tview.Pages) {
	model := tview.NewModal()
	model.SetText("Что-то пошло не так\n" +
		"Выберите действие:\n" +
		"1. Save: Повторно сохранить данные\n" +
		"2. Cancel: Закрыть приложение\n").
		AddButtons([]string{"Save", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				pages.AddPage("Binary", c.binaryButton(ctx, app, pages), true, false)
				pages.SwitchToPage("Binary")
			} else {
				app.Stop()
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("ErrorsSaveBinary", model, true, true)
}

func (c *CLI) errorsSaveText(ctx context.Context, app *tview.Application, pages *tview.Pages) {
	model := tview.NewModal()
	model.SetText("Что-то пошло не так\n" +
		"Выберите действие:\n" +
		"1. Save: Повторно сохранить данные\n" +
		"2. Cancel: Закрыть приложение\n").
		AddButtons([]string{"Save", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				pages.AddPage("Text", c.textButton(ctx, app, pages), true, false)
				pages.SwitchToPage("Text")
			} else {
				app.Stop()
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("ErrorsSaveText", model, true, true)
}
