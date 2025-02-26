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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

func main() {
	// Создаем приложение Fyne
	a := app.New()
	w := a.NewWindow("CSV to Excel Converter")
	w.Resize(fyne.NewSize(400, 400))

	// Поле для отображения статуса
	statusEntry := widget.NewEntry()
	statusEntry.SetText("Перетащите CSV-файл сюда или нажмите кнопку ниже")
	statusEntry.Disable()
	statusEntry.MultiLine = true
	statusEntry.Wrapping = fyne.TextWrapWord

	// Радиокнопки для выбора метода парсинга
	parseMethod := widget.NewRadioGroup([]string{"С разделителями", "Фиксированная ширина"}, nil)
	parseMethod.SetSelected("С разделителями") // По умолчанию выбраны разделители
	var currentMethod string = "С разделителями"

	// Контейнеры для динамического содержимого
	delimiterOptions := container.NewVBox()
	fixedWidthOptions := container.NewVBox()
	configContainer := container.NewMax(delimiterOptions) // По умолчанию показываем разделители

	// Чекбоксы для разделителей
	semicolonCheck := widget.NewCheck("Точка с запятой (;)", nil)
	semicolonCheck.Checked = true
	tabCheck := widget.NewCheck("Знак табуляции (\\t)", nil)
	tabCheck.Checked = false
	commaCheck := widget.NewCheck("Запятая (,)", nil)
	commaCheck.Checked = false
	spaceCheck := widget.NewCheck("Пробел ( )", nil)
	spaceCheck.Checked = false
	customCheck := widget.NewCheck("Другой", nil)
	customDelimiter := widget.NewEntry()
	customDelimiter.SetPlaceHolder("Введите разделитель")
	customDelimiter.Disable()
	customDelimiter.OnChanged = func(text string) {
		if customCheck.Checked && text == "" {
			statusEntry.SetText("Ошибка: Укажите пользовательский разделитель")
		}
	}

	customCheck.OnChanged = func(checked bool) {
		if checked {
			customDelimiter.Enable()
		} else {
			customDelimiter.Disable()
		}
	}

	delimiterOptions.Add(semicolonCheck)
	delimiterOptions.Add(tabCheck)
	delimiterOptions.Add(commaCheck)
	delimiterOptions.Add(spaceCheck)
	delimiterOptions.Add(container.NewHBox(customCheck, customDelimiter))

	// Элементы для фиксированной ширины
	fixedWidthEntry := widget.NewEntry()
	fixedWidthEntry.SetPlaceHolder("Введите ширины полей через запятую, например: 10,20,15")
	fixedWidthOptions.Add(widget.NewLabel("Ширины полей (в символах):"))
	fixedWidthOptions.Add(fixedWidthEntry)

	// Переключение интерфейса при выборе радиокнопки
	parseMethod.OnChanged = func(selected string) {
		currentMethod = selected
		if selected == "С разделителями" {
			configContainer.Objects = []fyne.CanvasObject{delimiterOptions}
		} else {
			configContainer.Objects = []fyne.CanvasObject{fixedWidthOptions}
		}
		configContainer.Refresh()
	}

	// Кнопка для выбора файла
	selectButton := widget.NewButton("Выбрать CSV-файл", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				statusEntry.SetText("Ошибка при выборе файла")
				return
			}
			csvPath := reader.URI().Path()
			reader.Close()
			processFile(csvPath, statusEntry, currentMethod, semicolonCheck, tabCheck, commaCheck, spaceCheck, customCheck, customDelimiter, fixedWidthEntry)
		}, w)
	})

	// Устанавливаем содержимое окна
	w.SetContent(container.NewVBox(
		statusEntry,
		parseMethod,
		configContainer,
		layout.NewSpacer(),
		selectButton,
	))

	// Поддержка drag-and-drop
	w.SetOnDropped(func(_ fyne.Position, uris []fyne.URI) {
		if len(uris) > 0 {
			csvPath := uris[0].Path()
			processFile(csvPath, statusEntry, currentMethod, semicolonCheck, tabCheck, commaCheck, spaceCheck, customCheck, customDelimiter, fixedWidthEntry)
		}
	})

	// Запускаем приложение
	w.ShowAndRun()
}

// Обработка файла
func processFile(csvPath string, statusEntry *widget.Entry, method string,
	semicolonCheck, tabCheck, commaCheck, spaceCheck, customCheck *widget.Check,
	customDelimiter, fixedWidthEntry *widget.Entry) {
	file, err := os.Open(csvPath)
	if err != nil {
		statusEntry.SetText(fmt.Sprintf("Ошибка открытия CSV файла: %v", err))
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	if method == "С разделителями" {
		// Проверяем, выбран ли хотя бы один разделитель
		if !semicolonCheck.Checked && !tabCheck.Checked && !commaCheck.Checked && !spaceCheck.Checked && !customCheck.Checked {
			statusEntry.SetText("Ошибка: Выберите хотя бы один разделитель")
			return
		}
		if customCheck.Checked && customDelimiter.Text == "" {
			statusEntry.SetText("Ошибка: Укажите пользовательский разделитель")
			return
		}

		// Устанавливаем первый выбранный разделитель (если несколько, берем первый)
		if semicolonCheck.Checked {
			reader.Comma = ';'
		} else if tabCheck.Checked {
			reader.Comma = '\t'
		} else if commaCheck.Checked {
			reader.Comma = ','
		} else if spaceCheck.Checked {
			reader.Comma = ' '
		} else if customCheck.Checked {
			if len([]rune(customDelimiter.Text)) > 1 {
				statusEntry.SetText("Ошибка: Пользовательский разделитель должен быть одним символом")
				return
			}
			reader.Comma = []rune(customDelimiter.Text)[0]
		}
	} else { // Фиксированная ширина
		if fixedWidthEntry.Text == "" {
			statusEntry.SetText("Ошибка: Укажите ширины полей")
			return
		}
		// Пока не реализуем фиксированную ширину полностью, только читаем как есть
		reader.FieldsPerRecord = -1 // Отключаем проверку количества полей
	}

	records, err := reader.ReadAll()
	if err != nil {
		if parseErr, ok := err.(*csv.ParseError); ok {
			statusEntry.SetText(fmt.Sprintf("Ошибка парсинга CSV на строке %d, колонке %d: %v\nПроверьте формат файла", parseErr.Line, parseErr.Column, err))
		} else {
			statusEntry.SetText(fmt.Sprintf("Ошибка чтения CSV файла: %v", err))
		}
		return
	}

	if len(records) == 0 {
		statusEntry.SetText("CSV файл пуст")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	sheet := "Sheet1"
	for col, header := range records[0] {
		cell := fmt.Sprintf("%c1", 'A'+col)
		f.SetCellValue(sheet, cell, header)
		style, _ := f.NewStyle(&excelize.Style{
			NumFmt: 49, // Текстовый формат
		})
		f.SetColStyle(sheet, fmt.Sprintf("%c", 'A'+col), style)
	}

	for rowIdx, row := range records[1:] {
		for colIdx, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+2)
			f.SetCellValue(sheet, cell, value)
		}
	}

	outputFile := strings.TrimSuffix(csvPath, ".csv") + ".xlsx"
	if err := f.SaveAs(outputFile); err != nil {
		statusEntry.SetText(fmt.Sprintf("Ошибка сохранения Excel файла: %v", err))
		return
	}

	statusEntry.SetText(fmt.Sprintf("Файл успешно сохранен как %s", outputFile))
}