package controller

import (
	"log/slog"
	"os"

	"brightrock.co.za/brgeo/api"
	"brightrock.co.za/brgeo/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func StartAndServe() {

	port := os.Getenv("PORT")

	slog.Info("Starting server on port", "port", port)
	app := fiber.New()
	group := app.Group("/api")

	group.Get("/lookup/:ipaddress", getGeoInfo)

	err := app.Listen(":" + port)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}

func getGeoInfo(c *fiber.Ctx) error {
	ipaddress := c.Params("ipaddress")

	geo, retry := api.GetGeoInfo(ipaddress)
	slog.Info("Retrieval information", "ipaddress", ipaddress, "retries", retry)

	lresp := model.LookupResponse{}

	// Copy attributes between two structures
	// where structs have the same fields
	err := copier.Copy(&lresp, &geo)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(lresp)
}
