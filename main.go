package main

import "github.com/izaakdale/service-ticket/internal/app"

func main() {
	app.Run()
	// qrcode.WriteFile("hello", qrcode.Medium, 256, "test.png")
}
