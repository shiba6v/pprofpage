package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (r Controller) Health(c echo.Context) error {
	return c.String(http.StatusOK, "ok")
}
