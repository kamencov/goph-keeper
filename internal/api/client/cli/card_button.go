package cli

import (
	"context"
	"github.com/rivo/tview"
)

type CardData struct {
	Data string
}

func (c *CLI) cardButton(ctx context.Context, app *tview.Application, pages *tview.Pages) *tview.Form {
	var cardData CardData
	form := tview.NewForm()
	form.
		AddInputField("Card", "", 20, nil, func(text string) {
			cardData.Data = text
		}).
		AddButton("Save", func() {
			// Показываем подтверждение сохранения
			c.saveCardData(pages, form, &cardData)
		}).
		AddButton("Find", func() {
			// Показываем подтверждение сохранения
			//c.findCardData(pages, form, &cardData)
		}).
		AddButton("Delete", func() {
			// Показываем подтверждение сохранения
			//c.deleteCardData(pages, form, &cardData)
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	return form
}

func (c *CLI) saveCardData(pages *tview.Pages, form *tview.Form, cardData *CardData) {
	model := tview.NewModal()
	model.SetText("Вы хотите сохранить данные?\n" +
		"Card: " + cardData.Data)
	model.AddButtons([]string{"Save", "Correct", "Cancel"})
	model.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Save" {
			// Возврат в главное меню после сохранения
			clearFormCard(form, cardData)
			pages.SwitchToPage("Buttons_data")
		} else if buttonLabel == "Correct" {
			pages.SwitchToPage("Card")
		} else {
			// Возврат к форме ввода данных
			clearFormCard(form, cardData)
			pages.SwitchToPage("Card")
		}
	})

	// Добавляем модальное окно как новую страницу
	pages.AddPage("SaveConfirmation", model, true, true)
}

// Сбрасывает данные в форме и структуре
func clearFormCard(form *tview.Form, cardData *CardData) {
	cardData.Data = ""

	form.GetFormItem(0).(*tview.InputField).SetText("") // Очищаем поле Text
}
