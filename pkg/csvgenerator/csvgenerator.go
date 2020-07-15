package csvgenerator

import (
	"encoding/csv"
	"os"
)

// Generate - takes filename and data and generates csv
func Generate(filename string, data [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, v := range data {
		if err := writer.Write(v); err != nil {
			return err
		}
	}

	return nil
}
