package main

import (
	"github.com/JonPulfer/gofiggy/pkg/controller"
	"github.com/JonPulfer/gofiggy/pkg/handlers"
)

func main() {
	var eventHandler = handlers.NewWebsiteFetchHandler()
	controller.Start("default", eventHandler)
}
