package cli

import (
	"context"
	"github.com/rivo/tview"
)

type TextData struct {
	Text string
}

func (c *CLI) textButton(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {
	var textData TextData
	form := tview.NewForm()

	form.AddInputField("Text", "", 20, nil, func(text string) {
		textData.Text = text
	}).
		AddButton("Save", func() {
			// Показываем подтверждение сохранения
			c.saveTextData(pages, form, &textData)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	return form
}

// Модальное окно подтверждения сохранения
func (c *CLI) saveTextData(pages *tview.Pages, form *tview.Form, textData *TextData) {
	modal := tview.NewModal().
		SetText("Вы хотите сохранить данные?\n" +
			"Text: " + textData.Text).
		AddButtons([]string{"Save", "Correct", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Save" {
				// Возврат в главное меню после сохранения
				clearFormTextData(form, textData)
				pages.SwitchToPage("Buttons_data")
			} else if buttonLabel == "Correct" {
				pages.SwitchToPage("Text")
			} else {
				// Возврат к форме ввода данных
				clearFormTextData(form, textData)
				pages.SwitchToPage("Text")
			}
		})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("SaveConfirmation", modal, true, true)
}

// Сбрасывает данные в форме и структуре
func clearFormTextData(form *tview.Form, textData *TextData) {
	textData.Text = ""

	form.GetFormItem(0).(*tview.InputField).SetText("") // Очищаем поле Text
}
