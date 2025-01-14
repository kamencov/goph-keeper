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
		AddInputField("Card [only numbers]", "", 20, func(textToCheck string, lastChar rune) bool {
			// Проверяем, что каждый вводимый символ является цифрой
			return lastChar >= '0' && lastChar <= '9'
		}, func(text string) {
			// Проверяем текст полностью после завершения ввода
			if len(text) == 16 {
				cardData.Data = text
				c.log.Info("Card number saved", "card", cardData.Data)
			}
		}).
		AddButton("Save", func() {
			// Показываем подтверждение сохранения
			c.saveCardData(ctx, app, pages, form, &cardData)
		}).
		AddButton("Back", func() {
			pages.SwitchToPage("Buttons_data")
		}).
		AddButton("Quit", func() {
			app.Stop()
		})
	return form
}

func (c *CLI) saveCardData(ctx context.Context, app *tview.Application, pages *tview.Pages, form *tview.Form, cardData *CardData) {
	model := tview.NewModal()
	model.SetText("Вы хотите сохранить данные?\n" +
		"Card: " + cardData.Data)
	model.AddButtons([]string{"Save", "Correct", "Cancel"})
	model.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		if buttonLabel == "Save" {
			err := c.save.PostCards(ctx, c.token, cardData.Data)
			if err != nil {
				c.log.Error("failed handlers credentials", "error", err)
				c.errorsSaveCards(ctx, app, pages)
			} else {
				pages.SwitchToPage("Buttons_data")
			}
		} else if buttonLabel == "Correct" {
			// Возвращаем пользователя к заполнению данных
			pages.SwitchToPage("Card")
		} else {
			// Очищаем форму и возвращаемся к заполнению
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
