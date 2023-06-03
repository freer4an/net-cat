package main

import (
	"fmt"
	"os"
)

func logCatcher(msg string) {
	file, err := os.OpenFile("logs.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
