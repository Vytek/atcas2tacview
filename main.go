package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/putridparrot/GoUnits/speed"
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

func IntToString(n int) string {
	return strconv.Itoa(n)
}

func Float64ToString(f float64) string {
	/** converting the f variable into a string */
	/** 5 is the number of decimals */
	/** 64 is for float64 type*/
	return strconv.FormatFloat(f, 'f', 5, 64)
}

// Example: CMD 1136-varie.xls 1002 "DC9 ITAVIA" A1136 A1136
func main() {

	//Load Args
	argsWithoutProg := os.Args[1:]

	//Create data
	dateTimeST, _ := jodaTime.Parse("HHmmss", ST)
	fmt.Println(dateTimeST) //DEBUG

	workbook, err := xls.OpenFile(filepath.Join("data", argsWithoutProg[0]))

	if err != nil {
		log.Panic(err.Error())
	}

	// Кол-во листов в книге
	// Number of sheets in the workbook
	//
	// for i := 0; i <= workbook.GetNumberSheets()-1; i++ {}

	fmt.Println(workbook.GetNumberSheets())
	//For 1136 sheet 4 otherwise 1
	sheet, err := workbook.GetSheet(4)

	if err != nil {
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

	var strTimeToWrite string
	var sumDuration int32
	var s_LAT float64
	var s_LONG float64
	var s_ALTITUDE float64
	var s_VEL float64

	for i := 5; i <= sheet.GetNumberRows(); i++ {
		if row, err := sheet.GetRow(i); err == nil {
			//Time UTC/Zulu
			if cell, err := row.GetCol(0); err == nil {
				//Time Next
				fmt.Println(cell.GetType())
				xfIndex:=cell.GetXFIndex()
				formatIndex:=workbook.GetXFbyIndex(xfIndex)
				format:=workbook.GetFormatByIndex(formatIndex.GetFormatIndex())
				fmt.Println(format.GetFormatString(cell))
				dateTimeNow, _ := jodaTime.Parse("HH:mm:ss", cell.GetString()) //Read TIME from XSL
				fmt.Println(dateTimeNow)
				fmt.Println(cell.GetString())
				fmt.Println(cell.GetFloat64())
				if dateTimeNow.After(dateTimeST) {
					sumDuration = sumDuration + int32(dateTimeNow.Sub(dateTimeST).Seconds()) //TODO: Ricontrollare tutti i tempi!
					strTimeToWrite = fmt.Sprintf("#%s.%s\n", IntToString(int(sumDuration)), "00")
					dateTimeST = dateTimeNow
					//strTimeToWrite = fmt.Sprintf("#%s.%s\n", Float64ToTimeString(dateTimeNow.Sub(dateTimeST).Minutes()), Float64ToTimeString(dateTimeNow.Sub(dateTimeST).Seconds()))
				}
				_, _ = f.WriteString(strTimeToWrite)
			}
			//Latitude
			if cell, err := row.GetCol(8); err == nil {
				s_LAT = cell.GetFloat64()
			}
			//Longitude
			if cell, err := row.GetCol(9); err == nil {
				s_LONG = cell.GetFloat64()
			}
			//Altitude
			if cell, err := row.GetCol(10); err == nil {
				s_ALTITUDE = cell.GetFloat64()
				s_ALTITUDE = feet2meters(s_ALTITUDE)
			}
			//Velocity
			if cell, err := row.GetCol(10); err == nil {
				s_VEL = cell.GetFloat64()
				s_VEL = speed.Knots.ToMetresPerSecond(speed.Knots(s_VEL))
			}
		}
		//Write all to acmi file
		//fmt.Println(argsWithoutProg[1])
		//fmt.Println(argsWithoutProg[2])
		//cdfmt.Println(argsWithoutProg[3])
		//Coodinates
		strToWrite := fmt.Sprintf("%s,T=%s|%s|%s,IAS=%s,Name=%s,Squawk=%s,Label=%s\n",
			argsWithoutProg[1],
			Float64ToString(s_LONG),
			Float64ToString(s_LAT),
			Float64ToString(s_ALTITUDE),
			Float64ToString(s_VEL),
			argsWithoutProg[2],
			argsWithoutProg[3],
			argsWithoutProg[4])
		_, _ = f.WriteString(strToWrite)
	}

	//Write and sync file
	f.Sync()
}
