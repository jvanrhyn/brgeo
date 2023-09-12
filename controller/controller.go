package controller

import (
	"log/slog"
	"os"

	"brightrock.co.za/brgeo/api"
	"github.com/gofiber/fiber/v2"
)

func StartAndServe() {

	app := fiber.New()
	group := app.Group("/api")

	group.Get("/lookup/:ipaddress", getGeoInfo)

	port := os.Getenv("PORT")

}

func getGeoInfo(c *fiber.Ctx) error {
	ipaddress := c.Params("ipaddress")

	geo, retry := api.GetGeoInfo(ipaddress)
	slog.Info("Retrieval information", "ipaddress", ipaddress, "retries", retry)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"city": geo.City,
	})
}
