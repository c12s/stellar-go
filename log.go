package main

import (
	"fmt"
	"os"
)

func writeBytes(file *os.File, bytes []byte) error {
	_, err := file.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}

func Log(data []byte, traceid, spanid string) error {
	file, err := os.Create(fmt.Sprintf("logs/%s_%s.log", traceid, spanid))
	defer file.Close()
	if err != nil {
		return err
	}

	err = writeBytes(file, data)
	if err != nil {
		return err
	}

	return nil
}
