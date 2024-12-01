package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/shakinm/xlsReader/xls"
)

// Version
const Version = "0.0.1"

func main() {

	//Load Args
	argsWithoutProg := os.Args[1:]

	workbook, err := xls.OpenFile(filepath.Join("data", argsWithoutProg[0]))

	if err!=nil {
		log.Panic(err.Error())
	}

	// Кол-во листов в книге
	// Number of sheets in the workbook
	//
	// for i := 0; i <= workbook.GetNumberSheets()-1; i++ {}

	fmt.Println(workbook.GetNumberSheets())

	sheet, err := workbook.GetSheet(0)

	if err!=nil {
		log.Panic(err.Error())
	}

	// Имя листа
	// Print sheet name
	println(sheet.GetName())

	// Вывести кол-во строк в листе
	// Print the number of rows in the sheet
	println(sheet.GetNumberRows())

	for i := 0; i <= sheet.GetNumberRows(); i++ {
		if row, err := sheet.GetRow(i); err == nil {
			if cell, err := row.GetCol(1); err == nil {

				// Значение ячейки, тип строка
				// Cell value, string type
				fmt.Println(cell.GetString())

				//fmt.Println(cell.GetInt64())
				//fmt.Println(cell.GetFloat64())

				// Тип ячейки (записи)
				// Cell type (records)
				fmt.Println(cell.GetType())
			}

		}
	}
}