package main

import (
	"Inventory_project_with_dyanamodb/admin"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/julienschmidt/httprouter"
)

func SetupRoutes(client *dynamodb.Client) *httprouter.Router {
	router := httprouter.New()

	adminHandler := &admin.NewHandler{
		Client: client,
	}

	router.GET("/admin/:id", adminHandler.GetProductByID)
	router.GET("/products", adminHandler.GetAllProducts)
	router.PUT("/product", adminHandler.UpdateProduct)
	router.DELETE("/product/:id", adminHandler.DeleteProduct)

	return router
}
