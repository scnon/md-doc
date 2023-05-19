package utils

import (
	"log"

	"github.com/labstack/echo/v4"
)

func Resp404(c echo.Context) error {
	return c.HTML(404, "404 Not Found")
}

func Resp500(c echo.Context, err error) error {
	log.Println("500:", err)
	return c.HTML(500, "500 Internal Server Error")
}
