package main

import (
	"context"
	"fmt"
	_ "ips-lacpass-backend/pkg/docs"
	"os"
	"os/signal"
)

//	@title			Lacpass App API
//	@version		0.1
//	@description	This is the official API for Lacpass mobile app.

//	@host	localhost:9081

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	app := New(LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
