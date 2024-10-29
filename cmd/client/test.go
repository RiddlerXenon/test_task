package main

import (
	"fmt"
	"time"

	"github.com/RiddlerXenon/test_task/internal/client"
)

func main() {
	client := client.NewHTTPClient("http://localhost:8080/api")

	err := client.Add("key0", "value1", 10*time.Second)
	if err != nil {
		fmt.Println("Ошибка при вызове Add:", err)
	} else {
		fmt.Println("Ключ успешно добавлен")
	}

	value, err := client.Get("key0")
	if err != nil {
		fmt.Println("Ошибка при вызове Get:", err)
	} else {
		fmt.Println("Полученное значение:", value)
	}

	err = client.Set("key0", "value2", 15*time.Second)
	if err != nil {
		fmt.Println("Ошибка при вызове Set:", err)
	} else {
		fmt.Println("Значение успешно обновлено")
	}

	value, err = client.Get("key0")
	if err != nil {
		fmt.Println("Ошибка при повторном вызове Get:", err)
	} else {
		fmt.Println("Полученное новое значение:", value)
	}

	err = client.Del("key0")
	if err != nil {
		fmt.Println("Ошибка при вызове Del:", err)
	} else {
		fmt.Println("Ключ успешно удален")
	}
}
