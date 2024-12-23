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
		AddButton("Find", func() {
			// открывает перечень всего сохраненного
			c.getResource(ctx, app, pages, form)
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
					c.log.Error("failed save credentials", "error", err)
					c.errorsSave(ctx, app, pages)
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

func (c *CLI) getResource(ctx context.Context,
	app *tview.Application,
	pages *tview.Pages,
	form *tview.Form) {

	// создаем макет базы.
	databases := tview.NewList().ShowSecondaryText(false)
	databases.SetBorder(true).SetTitle("Database")

	// создаем колонки таблицы.
	columns := tview.NewList().ShowSecondaryText(true)
	columns.SetBorder(true).SetTitle("Columns")

	// создаем список таблиц.
	tables := tview.NewList()
	tables.ShowSecondaryText(false).
		SetDoneFunc(func() {
			tables.Clear()
			columns.Clear()
			app.SetFocus(databases)
		})
	tables.SetBorder(true).SetTitle("Tables")

	//// Создаем слои.
	//flex := tview.NewFlex().
	//	AddItem(databases, 0, 1, true).
	//	AddItem(tables, 0, 1, false).
	//	AddItem(columns, 0, 3, false)
}
