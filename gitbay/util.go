package main

import (
	"encoding/csv"
	"os"
)

func AddRowToCSV(filepath string, row []string) error {
	file, err := os.OpenFile(filepath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	err = writer.Write(row)
	if err != nil {
		return err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	return nil
}

func Contains(seq []string, target string) bool {
	for _, item := range seq {
		if target == item {
			return true
		}
	}
	return false
}
