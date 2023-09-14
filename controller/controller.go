package controller

import (
	"log/slog"
	"os"
	"time"

	"brightrock.co.za/brgeo/api"
	"brightrock.co.za/brgeo/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func StartAndServe() {

	//logger := slog.Default()
	port := os.Getenv("PORT")

	slog.Info("Starting server on port", "port", port)
	app := fiber.New()

	group := app.Group("/api")
	// group.Use(slogfiber.New(logger))

	group.Get("/lookup/:ipaddress", getGeoInfo)

	err := app.Listen(":" + port)
	if err != nil {
		slog.Error("Error starting server", "error", err)
	}
}

func getGeoInfo(c *fiber.Ctx) error {
	ipaddress := c.Params("ipaddress")

	var err error

	// Try and find the element in the Cache
	cg, err := api.GetCacheById(ipaddress)
	if err == nil {
		slog.Info("Retrieved item from cache for ip", "ipaddress", ipaddress)
		return c.Status(fiber.StatusOK).JSON(cg)
	}

	geo, retry := api.GetGeoInfo(ipaddress)
	slog.Info("Retrieval information", "ipaddress", ipaddress, "retries", retry)

	lresp := model.LookupResponse{}

	// Copy attributes between two structures
	// where structs have the same fields
	err = copier.Copy(&lresp, &geo)
	if err != nil {
		return err
	}

	req := model.LookupRequest{
		IpAddress:    ipaddress,
		LookupTime:   time.Now(),
		LookupStatus: true,
	}

	api.Record(&req)

	// Store the item in the cache
	api.AddCacheItem(ipaddress, &lresp)
	slog.Info("Added item to cache for ip", "ipaddress", ipaddress)

	return c.Status(fiber.StatusOK).JSON(lresp)
}

func clearCache(c *fiber.Ctx) error {
	slog.Info("Clearing cache")
	api.Cache.Flush()
	return c.SendStatus(fiber.StatusOK)
}
