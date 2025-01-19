package cli

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (c *CLI) getResource(ctx context.Context,
	app *tview.Application,
	pages *tview.Pages) {

	columns := tview.NewTable().SetBorders(true)
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
			}
		}(tableName))
	}

	tables.AddItem("Back", "", 0, func() {
		pages.SwitchToPage("Buttons_data")
	})

	tables.AddItem("Quit", "", 0, func() {
		app.Stop()
	})

	flex := tview.NewFlex().
		AddItem(tables, 0, 1, true).
		AddItem(columns, 0, 3, false)

	pages.AddPage("GetAll", flex, true, true)
	app.SetRoot(pages, true)
}
