package main

import (
	"log"

	"example.com/rest-api/db"
	"example.com/rest-api/routes"
	"example.com/rest-api/telemetry" // Import your utils
	"github.com/Cyprinus12138/otelgin"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	
	db.InitDB()

	
	meterProvider, err := telemetry.InitMetrics()
	if err != nil {
		log.Fatal(err)
	}
	
	defer meterProvider.Shutdown(nil)

	server := gin.Default()


	server.Use(otelgin.Middleware("my-rest-api"))

	
	server.GET("/metrics", gin.WrapH(promhttp.Handler()))

	routes.RegisterRoutes(server)

	server.Run(":8081")
}
