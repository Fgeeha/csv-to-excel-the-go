package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Запрашиваем путь к CSV файлу
	fmt.Println("Введите путь к CSV файлу:")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	csvPath := scanner.Text()

	// Открываем CSV файл
	file, err := os.Open(csvPath)
	if err != nil {
		fmt.Printf("Ошибка открытия CSV файла: %v\n", err)
		return
	}
	defer file.Close()

	// Читаем CSV с разделителем ";"
	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Ошибка чтения CSV файла: %v\n", err)
		return
	}

	if len(records) == 0 {
		fmt.Println("CSV файл пуст")
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
		// Устанавливаем текстовый формат для всего столбца
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
		fmt.Printf("Ошибка сохранения Excel файла: %v\n", err)
		return
	}

	fmt.Printf("Файл успешно сохранен как %s\n", outputFile)
}