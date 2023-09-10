package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/shiba6v/eu"
	"github.com/shiba6v/pprofpage/app/registry"
)

func init() {

}

func run() error {
	e := echo.New()
	s3cli, err := registry.RegisterMinio()
	if err != nil {
		return eu.Wrap(err)
	}
	if err := registry.RegisterServer(e, s3cli); err != nil {
		eu.Wrap(err)
	}
	e.Logger.Fatal(e.Start("0.0.0.0:9000"))
	return nil
}

func main() {
	if err := run(); err != nil {
		panic(fmt.Sprintf("Error:%+v", err))
	}
}
