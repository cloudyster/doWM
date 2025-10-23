// Package main.
package main

import (
	"log/slog"

	"github.com/BobdaProgrammer/doWM/wm"
)

func main() {
	WM, err := wm.Create()
	if err != nil {
		slog.Error("Couldn't initialize window manager", "error:", err)
		return
	}
	defer WM.Close()

	WM.Run()
}
