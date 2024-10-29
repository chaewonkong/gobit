package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

// ReadLogFile reads log file with interval of 1 second.
func ReadLogFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("file open error: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// 1 sec interval
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			line, err := reader.ReadString('\n')
			if err != nil {
				if err.Error() == "EOF" {
					continue
				}
				return fmt.Errorf("read line error: %w", err)
			}
			fmt.Printf("line: %s\n", line)
		}
	}
}

// SaveOffset saves offset to file
func SaveOffset(offset int64, offsetFilePath string) error {
	f, err := os.Create(offsetFilePath)
	if err != nil {
		return fmt.Errorf("fail to create offset file: %w", err)
	}

	defer f.Close()

	_, err = f.WriteString(strconv.FormatInt(offset, 10))
	if err != nil {
		return fmt.Errorf("fail to write offset file: %w", err)
	}

	return nil
}

// GetSavedOffset gets saved offset
func GetSavedOffset(offsetFilePath string) (int64, error) {
	f, err := os.Open(offsetFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("fail to get offset file: %w", err)
	}

	defer f.Close()

	sc := bufio.NewScanner(f)
	if sc.Scan() {
		offset, err := strconv.ParseInt(sc.Text(), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("fail to convert offset to int64: %w", err)
		}

		return offset, nil
	}

	return 0, fmt.Errorf("fail to read offset")

}

func main() {
	// read env_var for log file path
	logPath := os.Getenv("LOG_PATH")
	err := ReadLogFile(logPath)
	if err != nil {
		log.Fatal(err)
	}
}
