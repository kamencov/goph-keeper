package cli

import (
	"context"
	"github.com/rivo/tview"
)

// buttonsStart - основная форма online.
func (c *CLI) buttonsStart(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.
		AddButton("Register", func() {
			c.log.Info("Switching to Register page")
			pages.AddPage("Register", c.register(ctx, app, pages), true, false)
			pages.SwitchToPage("Register")
		}).
		AddButton("Auth", func() {
			c.log.Info("Switching to AuthUser page")
			pages.AddPage("AuthUser", c.authUser(ctx, app, pages), true, false)
			pages.SwitchToPage("AuthUser")
		}).
		AddButton("Quit", func() {
			c.log.Info("Stopping application")
			app.Stop()
		})
	form.SetBorder(true).SetTitle("goph-keeper").
		SetTitleAlign(tview.AlignCenter)

	form.AddFormItem(tview.NewTextView().SetText("Добро пожаловать! \nВыберите действие:\n" +
		"1. Регистрация: Если вы впервые, то требуется пройти регистрацию\n" +
		"2. Авторизация: Если вы уже зарегистрированы, то требуется пройти авторизацию\n"))

	return form
}

// buttonsOffline - основная форма offline.
func (c *CLI) buttonsOffline(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.
		AddButton("Auth", func() {
			c.log.Info("Switching to AuthUser page")
			pages.AddPage("AuthUser", c.authUserOffline(ctx, app, pages), true, false)
			pages.SwitchToPage("AuthUser")
		}).
		AddButton("Quit", func() {
			c.log.Info("Stopping application")
			app.Stop()
		})

	form.SetBorder(true).SetTitle("goph-keeper").
		SetTitleAlign(tview.AlignCenter)

	form.AddFormItem(tview.NewTextView().SetText("Добро пожаловать в оффлаин режим!\n" +
		"Данный режим только для уже зарегистрированных пользователей\n"))

	return form
}

// buttonsData - основная форма для работы с данными.
func (c *CLI) buttonsData(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {
	form := tview.NewForm()
	form.
		AddButton("Find and delete", func() {
			// открывает перечень всего сохраненного
			c.getResource(ctx, app, pages)
		}).
		AddButton("Credentials", func() {
			pages.AddPage("Credentials", c.credentials(ctx, app, pages), true, false)
			pages.SwitchToPage("Credentials")
		}).
		AddButton("Text", func() {
			pages.AddPage("Text", c.textButton(ctx, app, pages), true, false)
			pages.SwitchToPage("Text")
		}).
		AddButton("Binary", func() {
			pages.AddPage("Binary", c.binaryButton(ctx, app, pages), true, false)
			pages.SwitchToPage("Binary")
		}).
		AddButton("Card", func() {
			pages.AddPage("Card", c.cardButton(ctx, app, pages), true, false)
			pages.SwitchToPage("Card")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("goph-keeper").
		SetTitleAlign(tview.AlignCenter)

	form.AddFormItem(tview.NewTextView().SetText("Выберите действие:\n" +
		"1. Find and delete: Найти или удалить данные\n" +
		"2. Credentials: Если вы хотите сохранить данные\n" +
		"3. Text: Если вы хотите сохранить текст\n" +
		"4. Binary: Если вы хотите сохранить бинарные данные\n" +
		"5. Card: Если вы хотите сохранить данные карты\n"))

	return form
}
