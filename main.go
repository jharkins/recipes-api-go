package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
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
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// It's not a "file not found" error, so it's a genuine issue with the config file
			fmt.Printf("Error reading config file: %s\n", err)
			os.Exit(1)
		}
		// If it's a "file not found" error, we just continue and rely on environment variables
	}

	viper.SetEnvPrefix("RECIPES")
	viper.AutomaticEnv() // Automatically read environment variables

	// Bind the environment variables to the Viper keys
	viper.BindEnv("mysql.user", "RECIPES_MYSQL_USER")
	viper.BindEnv("mysql.password", "RECIPES_MYSQL_PASSWORD")
	viper.BindEnv("mysql.host", "RECIPES_MYSQL_HOST")
	viper.BindEnv("mysql.port", "RECIPES_MYSQL_PORT")
	viper.BindEnv("mysql.database", "RECIPES_MYSQL_DATABASE")
}

func startServer() {
	initDB()
	r := gin.Default()
	r.Use(cors.Default())
	api := r.Group("/api")

	public := api.Group("/recipe")
	{
		public.GET("/", listRecipes)
		public.GET("/:id", getRecipe)
		public.GET("/random", getRandomRecipe)
	}

	private := api.Group("/recipe")
	private.Use(authMiddleware())
	{
		private.POST("/", createRecipe)
		private.PUT("/:id", updateRecipe)
		private.DELETE("/:id", deleteRecipe)
	}

	r.Run() // Default: ":8080"
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the user is authenticated
		if isUserAuthenticated(c.Request) {
			// User is authenticated, call the next middleware
			c.Next()
		} else {
			// User is not authenticated, return a 401 Unauthorized response
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		}
	}
}

func isUserAuthenticated(r *http.Request) bool {
	// Check if the user is authenticated
	return false
}
