package tools

import (
	"encoding/csv"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"tunaiku/models"
)

func ExportStock(header []string, stock []models.Stock, filename string, extension string) string {
	if filename == "" {
		filename = "result-" + time.Now().String() + extension
	}

	if stat, err := os.Stat("/csv-data"); err == nil && stat.IsDir() {
		os.Mkdir("csv-data", os.ModePerm)
	}

	basepath := path.Base("/csv-data")
	filepath := path.Base(filename)

	file, err := os.Create(strings.Join([]string{basepath, filepath}, "/"))
	checkError("Cannot create file", err)
	defer file.Close()

	// Initialize csv writer
	writer := csv.NewWriter(file)

	// Write header into file
	writer.Write(header)
	for _, item := range stock {
		record := []string{
			item.ID.Hex(),
			item.DateStr,
			strconv.Itoa(item.Open),
			strconv.Itoa(item.High),
			strconv.Itoa(item.Low),
			strconv.Itoa(item.Close),
		}
		err = writer.Write(record)
		checkError("Cannot write to file", err)
	}
	defer writer.Flush()
	return filename
}

// Private methods
func checkError(message string, err error) {
	if err != nil {
		log.Println(message, err)
	}
}
