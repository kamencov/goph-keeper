package main

import "goph-keeper/internal/api/client"

func main() {

	// запускаем приложение
	if err := client.RunClient(); err != nil {
		panic(err)
	}
}
