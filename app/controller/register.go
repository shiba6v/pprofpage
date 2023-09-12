package controller

import (
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/shiba6v/eu"
)

type RegisterPProfResponse struct {
	Path string
}

func (r Controller) RegisterPProf(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return eu.Wrap(err)
	}
	f, err := file.Open()
	if err != nil {
		return eu.Wrap(err)
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return eu.Wrap(err)
	}
	id := uuid.NewString()
	if err := r.storage.UploadObject(c.Request().Context(), id, b); err != nil {
		return eu.Wrap(err)
	}
	return c.JSON(http.StatusOK, RegisterPProfResponse{
		Path: fmt.Sprintf("/pprof/%s", id),
	})
}
