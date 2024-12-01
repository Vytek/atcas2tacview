package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/shakinm/xlsReader/xls"
	"github.com/vjeantet/jodaTime"

)

// Version
const Version = "0.0.1"

// Start time
const ST = "180000"

func DDHHMMZmmmYY() string {
	current_time := time.Now().UTC()
	return fmt.Sprintf(current_time.Format("021504ZJan06"))
}

// https://ispycode.com/GO/Math/Metric-Conversions/Distance/Feet-to-meters
func feet2meters(feet float64) float64 {
	return feet * 0.3048
}

// https://siongui.github.io/2018/02/25/go-get-file-name-without-extension/
func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}

func main() {

	//Load Args
	argsWithoutProg := os.Args[1:]

	//Create data
	dateTimeST, _ := jodaTime.Parse("HHmmss", ST)
	fmt.Println(dateTimeST) //DEBUG

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

	//Create and save acmi file (TacView)
	BOF := "FileType=text/acmi/tacview\nFileVersion=2.2\n"
	GIOF := "0,Author=Enrico Speranza\n0,Title=ATCAS Radar activity near ITAVIA I-TIGI IH870 A1136\n0,ReferenceTime=1980-06-27T18:00:00Z\n"
	//Open with name
	f, err := os.Create("out/nearadaractivity19800627180000Z" + FilenameWithoutExtension(argsWithoutProg[0]) + "v" + DDHHMMZmmmYY() + ".acmi")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, _ = f.WriteString(BOF)
	_, _ = f.WriteString(GIOF)

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

	//Write and sync file
	f.Sync()
}