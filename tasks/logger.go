package tasks

import (
	"fmt"
	"os"
	"time"
)

func WriteToFile(id string, status string) {
	fileName := "task_logs.txt"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	logEntry := fmt.Sprintf("Task ID: %s, Date: %s, Status: %s\n", id, time.Now().Format(time.RFC3339), status)
	_, err = file.WriteString(logEntry)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}

}
