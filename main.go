package main

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	cmd := &cobra.Command{
		Use:   "recipes",
		Short: "Recipes API",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	cmd.Execute()
}

func initConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err)
		os.Exit(1)
	}
}

func startServer() {
	initDB()
	r := gin.Default()
	api := r.Group("/api")

	recipe := api.Group("/recipe")
	{
		recipe.POST("/", createRecipe)
		recipe.GET("/:id", getRecipe)
		recipe.PUT("/:id", updateRecipe)
		recipe.DELETE("/:id", deleteRecipe)
		recipe.GET("/", listRecipes)
	}

	r.Run() // Default: ":8080"
}
