package controller

import (
	"log/slog"
	"os"
	"time"

	"github.com/go-errors/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/jvanrhyn/brgeo/internal/api"
	"github.com/jvanrhyn/brgeo/model"
	slogfiber "github.com/samber/slog-fiber"
)

func StartAndServe() {

	logger := slog.Default()
	port := os.Getenv("PORT")

	go slog.Info("Starting server on port", "port", port)
	app := fiber.New()

	app.Use(slogfiber.New(logger))

	group := app.Group("/api")
	cacheGroup := app.Group("/cache")

	group.Get("/lookup/:ipaddress", getGeoInfo)
	cacheGroup.Post("/clear", clearCache)

	err := app.Listen(":" + port)
	if err != nil {
		go slog.Error("Error starting server", "error", err)
	}
}

func getGeoInfo(c *fiber.Ctx) error {
	ipaddress := c.Params("ipaddress")

	var err error

	// Try and find the element in the Cache
	cg, err := api.GetCacheById(ipaddress)
	if err == nil {
		go slog.Info("Retrieved item from cache for ip", "ipaddress", ipaddress)
		c.Status(fiber.StatusOK).JSON(cg)
	}

	geo, retry := api.GetGeoInfo(ipaddress)
	go slog.Info("Retrieval information", "ipaddress", ipaddress, "retries", retry)

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
	err = api.AddCacheItem(ipaddress, &lresp)
	if err != nil {
		stack := err.(*errors.Error).ErrorStack()
		go slog.Error("Error adding item to cache", "error", err, "stacktrace", stack)
	} else {
		go slog.Info("Added item to cache for ip", "ipaddress", ipaddress)
	}

	return c.Status(fiber.StatusOK).JSON(lresp)
}

func clearCache(c *fiber.Ctx) error {
	go slog.Info("Clearing cache")
	api.Cache.Flush()
	return c.SendStatus(fiber.StatusOK)
}
