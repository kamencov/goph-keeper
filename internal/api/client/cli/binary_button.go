package cli

import (
	"context"
	"github.com/rivo/tview"
)

type BinaryDataCLI struct {
	Data string
}

func (c *CLI) binaryButton(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {
	var binaryData BinaryDataCLI
	form := tview.NewForm()
	form.
		AddInputField("Binary", "", 20, nil, func(text string) {
			binaryData.Data = text
		}).
		AddButton("Save", func() {
			// Показываем подтверждение сохранения
			c.saveBinaryData(ctx, app, pages, form, &binaryData)
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("Buttons_data")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	return form

}

func (c *CLI) saveBinaryData(ctx context.Context, app *tview.Application, pages *tview.Pages, form *tview.Form, binaryData *BinaryDataCLI) {
	model := tview.NewModal()
	model.SetText("Вы хотите сохранить данные?\n" +
		"Binary: " + binaryData.Data)
	model.AddButtons([]string{"Save", "Correct", "Cancel"})
	model.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Save" {
			err := c.save.PostBinaryData(ctx, c.token, binaryData.Data)
			if err != nil {
				c.log.Error("failed save credentials", "error", err)
				c.errorsSaveBinary(ctx, app, pages)
			} else {
				pages.SwitchToPage("Buttons_data")
			}
		} else if buttonLabel == "Correct" {
			pages.SwitchToPage("Binary")
		} else {
			// Возврат к форме ввода данных
			clearFormBinary(form, binaryData)
			pages.SwitchToPage("Binary")
		}
	})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("SaveConfirmation", model, true, true)
}

// Сбрасывает данные в форме и структуре
func clearFormBinary(form *tview.Form, binaryData *BinaryDataCLI) {
	binaryData.Data = ""

	form.GetFormItem(0).(*tview.InputField).SetText("") // Очищаем поле Text
}
