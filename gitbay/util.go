package main

import (
	"encoding/csv"
	"os"
	"time"

	"github.com/gabriel-vasile/mimetype"
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

func ReadCSV(filepath string) ([][]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	r := csv.NewReader(file)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func Contains(seq []string, target string) bool {
	for _, item := range seq {
		if target == item {
			return true
		}
	}
	return false
}

func GetContentType(data []byte) (string, error) {
	mime := mimetype.Detect(data)
	return mime.String(), nil
}

func LogCrash(log ...string) error {
	msg := "crashed"
	if len(log) > 0 {
		msg = log[0]
	}
	row := []string{
		time.Now().Format("2006-01-02 15:04:05"),
		msg,
	}
	err := AddRowToCSV(DATA_STORE_PATH+"node_crash.csv", row)
	if err != nil {
		return err
	}
	return nil
}
