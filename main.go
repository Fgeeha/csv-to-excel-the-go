package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Создаем приложение Fyne
	a := app.New()
	w := a.NewWindow("CSV to Excel Converter")
	w.Resize(fyne.NewSize(400, 200))

	// Поле для отображения статуса (доступно для копирования)
	statusEntry := widget.NewEntry()
	statusEntry.SetText("Перетащите CSV-файл сюда или нажмите кнопку ниже")
	statusEntry.Disable() // Делаем только для чтения
	statusEntry.MultiLine = true
	statusEntry.Wrapping = fyne.TextWrapWord

	// Кнопка для выбора файла
	selectButton := widget.NewButton("Выбрать CSV-файл", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				statusEntry.SetText("Ошибка при выборе файла")
				return
			}
			csvPath := reader.URI().Path()
			reader.Close()
			processFile(csvPath, statusEntry)
		}, w)
	})

	// Устанавливаем содержимое окна
	w.SetContent(container.NewVBox(
		statusEntry,
		selectButton,
	))

	// Поддержка drag-and-drop
	w.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			csvPath := uris[0].Path()
			processFile(csvPath, statusEntry)
		}
	})

	// Запускаем приложение
	w.ShowAndRun()
}

// Обработка файла
func processFile(csvPath string, statusEntry *widget.Entry) {
	// Открываем CSV файл
	file, err := os.Open(csvPath)
	if err != nil {
		statusEntry.SetText(fmt.Sprintf("Ошибка открытия CSV файла: %v", err))
		return
	}
	defer file.Close()

	// Читаем CSV с разделителем ";"
	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true // Удаляем ведущие пробелы
	records, err := reader.ReadAll()
	if err != nil {
		// Если ошибка парсинга, показываем проблемную строку
		if parseErr, ok := err.(*csv.ParseError); ok {
			statusEntry.SetText(fmt.Sprintf("Ошибка парсинга CSV на строке %d, колонке %d: %v\nПроверьте формат файла (разделитель ';', кодировка UTF-8)", parseErr.Line, parseErr.Column, err))
		} else {
			statusEntry.SetText(fmt.Sprintf("Ошибка чтения CSV файла: %v", err))
		}
		return
	}

	if len(records) == 0 {
		statusEntry.SetText("CSV файл пуст")
		return
	}

	// Создаем новый Excel файл
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Устанавливаем названия столбцов (первая строка)
	sheet := "Sheet1"
	for col, header := range records[0] {
		cell := fmt.Sprintf("%c1", 'A'+col)
		f.SetCellValue(sheet, cell, header)
		style, _ := f.NewStyle(&excelize.Style{
			NumFmt: 49, // @ - текстовый формат
		})
		f.SetColStyle(sheet, fmt.Sprintf("%c", 'A'+col), style)
	}

	// Заполняем данные (начиная со второй строки)
	for rowIdx, row := range records[1:] {
		for colIdx, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+2)
			f.SetCellValue(sheet, cell, value)
		}
	}

	// Сохраняем Excel файл
	outputFile := strings.TrimSuffix(csvPath, ".csv") + ".xlsx"
	if err := f.SaveAs(outputFile); err != nil {
		statusEntry.SetText(fmt.Sprintf("Ошибка сохранения Excel файла: %v", err))
		return
	}

	statusEntry.SetText(fmt.Sprintf("Файл успешно сохранен как %s", outputFile))
}
