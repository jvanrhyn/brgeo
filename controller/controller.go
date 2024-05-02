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

	// Try and find the element in the Cache
	var cg, err = api.GetCacheById(ipaddress)
	if err == nil {
		go slog.Info("Retrieved item from cache for ip", "ipaddress", ipaddress)
		err = c.Status(fiber.StatusOK).JSON(cg)
		if err != nil {
			stack := err.(*errors.Error).ErrorStack()
			go slog.Error("error while retrieving item from cache", "error", err, "stacktrace", stack)
			return err
		}
	}

	geo, retry := api.GetGeoInfo(ipaddress)
	go slog.Info("Retrieval information", "ipaddress", ipaddress, "retries", retry)

	response := model.LookupResponse{}

	// Copy attributes between two structures
	// where structs have the same fields
	err = copier.Copy(&response, &geo)
	if err != nil {
		return err
	}

	req := model.LookupRequest{
		IpAddress:    ipaddress,
		LookupTime:   time.Now(),
		LookupStatus: true,
	}

	err = api.Record(&req)
	if err != nil {
		stack := err.(*errors.Error).ErrorStack()
		go slog.Error("error while recording lookup", "error", err, "stacktrace", stack)
		return err
	}

	// Store the item in the cache
	err = api.AddCacheItem(ipaddress, &response)
	if err != nil {
		stack := err.(*errors.Error).ErrorStack()
		go slog.Error("Error adding item to cache", "error", err, "stacktrace", stack)
	} else {
		go slog.Info("Added item to cache for ip", "ipaddress", ipaddress)
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// clearCache clears the cache and logs an info message.
//
// It takes a fiber.Ctx parameter which represents the context of the HTTP request.
//
// It returns an error if there is an issue with the cache clearing operation.
//
// It sends a status OK response to the client after clearing the cache.
//
// Example:
//
//	err := clearCache(c)
//	if err != nil {
//		log.Fatal(err)
//	}
//	c.SendStatus(fiber.StatusOK)
func clearCache(c *fiber.Ctx) error {
	go slog.Info("Clearing cache")
	api.Cache.Flush()
	return c.SendStatus(fiber.StatusOK)
}
