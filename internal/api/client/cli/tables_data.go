package cli

import (
	"context"
	"fmt"
	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"strconv"
)

func (c *CLI) getResource(ctx context.Context,
	app *tview.Application,
	pages *tview.Pages) {

	columns := tview.NewTable().
		SetBorders(true).
		SetSelectable(true, false) // Включаем выделение строк
	columns.SetBorder(true).SetTitle("Columns")
	tables := tview.NewList()
	tables.ShowSecondaryText(false).
		SetDoneFunc(func() {
			columns.Clear()
			app.SetFocus(tables)
		})
	tables.SetBorder(true).SetTitle("Tables")

	// Список таблиц
	tableNames := []string{"credentials", "text_data", "binary_data", "cards"}
	for _, tableName := range tableNames {
		tables.AddItem(tableName, "", 0, func(tableName string) func() {
			return func() {
				columns.Clear()

				rows, err := c.getAll.GetAllData(ctx, c.token, tableName)
				if err != nil {
					c.log.Error("failed to getAll data from database", "error", err)
					return
				}

				defer rows.Close()

				// Получение имен столбцов
				columnNames, err := rows.Columns()
				if err != nil {
					panic(err)
				}

				// Добавляем заголовки для таблицы
				for colIndex, colName := range columnNames {
					columns.SetCell(0, colIndex, &tview.TableCell{
						Text:  colName,
						Align: tview.AlignCenter,
						Color: tcell.ColorBlue,
					})
				}

				// Чтение строк из таблицы
				rowIndex := 1
				for rows.Next() {
					// Создаем массив для хранения значений
					values := make([]any, len(columnNames))
					valuePtrs := make([]any, len(columnNames))
					for i := range values {
						valuePtrs[i] = &values[i]
					}

					// Сканируем строку
					if err := rows.Scan(valuePtrs...); err != nil {
						panic(err)
					}

					// Добавляем строку в таблицу для отображения
					for colIndex, value := range values {
						var text string
						if value != nil {
							text = fmt.Sprintf("%v", value)
						}
						columns.SetCell(rowIndex, colIndex, &tview.TableCell{
							Text:  text,
							Align: tview.AlignLeft,
							Color: tcell.ColorWhite,
						})
					}
					rowIndex++
				}

				// Проверка на ошибки после итерации
				if err := rows.Err(); err != nil {
					panic(err)
				}

				// Обработка выбора строки
				columns.SetSelectedFunc(func(row, column int) {
					if row == 0 {
						return // Пропускаем заголовок
					}

					// Получаем данные из выбранной строки
					data := []string{}
					for col := 0; col < columns.GetColumnCount(); col++ {
						cell := columns.GetCell(row, col)
						data = append(data, cell.Text)
					}

					// Показываем контекстное меню
					menu := tview.NewModal().
						SetText(fmt.Sprintf("Действия с данными: %v", data)).
						AddButtons([]string{"Copy", "Delete", "Edit", "Cancel"}).
						SetDoneFunc(func(buttonIndex int, buttonLabel string) {
							switch buttonLabel {
							case "Copy":
								c.copyData(ctx, data, pages)
							case "Delete":
								c.delData(ctx, tableName, data, pages)
								columns.Clear()
							case "Edit":

							}
							pages.RemovePage("Menu")
						})

					pages.AddPage("Menu", menu, true, true)
				})
			}
		}(tableName))
	}

	buttonBack := tview.NewButton("BACK").SetSelectedFunc(func() {
		pages.SwitchToPage("Buttons_data")
	})

	buttonQuit := tview.NewButton("QUIT").SetSelectedFunc(func() {
		app.Stop()
	})

	flex := tview.NewFlex().
		SetDirection(tview.FlexRow). // Устанавливаем вертикальную ориентацию
		AddItem(
			tview.NewFlex(). // Вложенный Flex для таблиц
						AddItem(tables, 0, 1, true).
						AddItem(columns, 0, 3, false),
			0, 1, true,
		).
		AddItem(tview.NewBox().SetBorder(false), 0, 0, false).
		AddItem(
			tview.NewFlex(). // Вложенный Flex для кнопок
						AddItem(buttonBack, 0, 1, false).
						AddItem(buttonQuit, 0, 1, false),
			1, 0, false, // Высота 1 строки
		)

	pages.AddPage("GetAll", flex, true, true)
	app.SetRoot(pages, true)
}

func (c *CLI) copyData(ctx context.Context, data []string, pages *tview.Pages) {
	// Объединяем данные строки в одну строку
	textToCopy := fmt.Sprintf("%v", data)

	// Сохраняем в буфер обмена
	if err := clipboard.WriteAll(textToCopy); err != nil {
		c.log.Error("failed to copy to clipboard", "error", err)
		return
	}

	// Показываем уведомление пользователю
	modal := tview.NewModal().
		SetText("Данные скопированы в буфер обмена!").
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			pages.RemovePage("CopiedNotification")
		})
	pages.AddPage("CopiedNotification", modal, true, true)
}

func (c *CLI) delData(ctx context.Context, tableName string, data []string, pages *tview.Pages) {
	id, err := strconv.Atoi(data[0])
	if err != nil {
		c.log.Error("failed to convert id to int", "error", err)
		return
	}

	// Подтверждение удаления
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Вы уверены, что хотите удалить данные с ID %d?", id)).
		AddButtons([]string{"OK", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "OK" {
				// Удаление данных
				err := c.deleted.DeletedData(ctx, tableName, id)
				if err != nil {
					c.log.Error("failed to delete data", "error", err)
					// Отображаем сообщение об ошибке
					errorModal := tview.NewModal().
						SetText("Не удалось удалить данные!").
						AddButtons([]string{"OK"}).
						SetDoneFunc(func(int, string) {
							pages.RemovePage("ErrorNotification")
						})
					pages.AddPage("ErrorNotification", errorModal, true, true)
					return
				}

				// Уведомление об успешном удалении
				successModal := tview.NewModal().
					SetText("Данные успешно удалены!").
					AddButtons([]string{"OK"}).
					SetDoneFunc(func(int, string) {
						pages.RemovePage("SuccessNotification")
					})
				pages.AddPage("SuccessNotification", successModal, true, true)
			}
			// Закрываем подтверждение удаления
			pages.RemovePage("DeletedConfirmation")
		})

	// Добавляем окно подтверждения
	pages.AddPage("DeletedConfirmation", modal, true, true)
}
