package cli

import (
	"context"
	"github.com/rivo/tview"
	"log/slog"
)

type Resource struct {
	Resource string
	Login    string
	Password string
}

// Основная форма для ввода данных
func (c *CLI) credentials(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {
	const op = "cli.credentials"
	c.log.With(slog.String("op", op))

	c.log.Info("start credentials")
	var resource Resource

	form := tview.NewForm()
	form.
		AddInputField("Resource", "", 20, nil, func(text string) {
			resource.Resource = text
		}).
		AddInputField("Login", "", 20, nil, func(text string) {
			resource.Login = text
		}).
		AddInputField("Password", "", 20, nil, func(text string) {
			resource.Password = text
		}).
		AddButton("Save", func() {
			// Открываем модальное окно подтверждения
			c.saveResource(ctx, app, pages, form, &resource)
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("Buttons_data")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)

	return form
}

// Модальное окно подтверждения сохранения
func (c *CLI) saveResource(
	ctx context.Context,
	app *tview.Application,
	pages *tview.Pages,
	form *tview.Form,
	resource *Resource,
) {
	c.log.Info("cli.saveResource Start")

	modal := tview.NewModal().
		SetText("Вы хотите сохранить данные?\n" +
			"Resource: " + resource.Resource + "\n" +
			"Login: " + resource.Login + "\n" +
			"Password: " + resource.Password).
		AddButtons([]string{"Save", "Correct", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("SaveConfirmation") // Удаляем страницу с модальным окном
			switch buttonLabel {
			case "Save":
				err := c.save.PostLoginAndPassword(ctx, c.token, resource.Resource, resource.Login, resource.Password)
				if err != nil {
					c.log.Error("failed handlers credentials", "error", err)
					c.errorsSaveCredentials(ctx, app, pages)
				} else {
					pages.SwitchToPage("Buttons_data")
				}
			case "Correct":
				// Возвращаем пользователя к заполнению данных
				pages.SwitchToPage("Credentials")
			default:
				// Очищаем форму и возвращаемся к заполнению
				clearFormResource(form, resource)
				pages.SwitchToPage("Credentials")
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("SaveConfirmation", modal, true, true)
}

// Сбрасывает данные в форме и структуре
func clearFormResource(form *tview.Form, resource *Resource) {
	resource.Resource = ""
	resource.Login = ""
	resource.Password = ""

	form.GetFormItem(0).(*tview.InputField).SetText("") // Очищаем поле Resource
	form.GetFormItem(1).(*tview.InputField).SetText("") // Очищаем поле Login
	form.GetFormItem(2).(*tview.InputField).SetText("") // Очищаем поле Password
}
