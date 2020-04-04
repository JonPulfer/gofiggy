package main

import (
	"github.com/JonPulfer/gofiggy/pkg/controller"
	"github.com/JonPulfer/gofiggy/pkg/events"
)

func main() {
	var eventHandler = events.NewMockHandler()
	controller.Start("default", eventHandler)
}
