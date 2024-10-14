package main

import (
	"context"

	"github.com/Elvilius/go-musthave-diploma-tpl.git/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := app.New()

	app.RunContext(ctx)
}
