package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func createRecipe(c *gin.Context) {
	var r Recipe
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec(`INSERT INTO recipe (name, enough_for, origin, ingredients, description, kind, prep_time, difficulty, notes, cook_time, serving_size, rating) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		r.Name, r.EnoughFor, r.Origin, r.Ingredients, r.Description, r.Kind, r.PrepTime, r.Difficulty, r.Notes, r.CookTime, r.ServingSize, r.Rating)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe created"})
}

func getRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	row := db.QueryRow("SELECT * FROM recipe WHERE id=?", id)
	var r Recipe
	err = row.Scan(&r.ID, &r.Name, &r.EnoughFor, &r.Origin, &r.Ingredients, &r.Description, &r.Kind, &r.PrepTime, &r.Difficulty, &r.Notes, &r.CookTime, &r.ServingSize, &r.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, r)
}

func updateRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var r Recipe
	if err := c.BindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec(`UPDATE recipe SET name=?, enough_for=?, origin=?, ingredients=?, description=?, kind=?, prep_time=?, difficulty=?, notes=?, cook_time=?, serving_size=?, rating=? WHERE id=?`,
		r.Name, r.EnoughFor, r.Origin, r.Ingredients, r.Description, r.Kind, r.PrepTime, r.Difficulty, r.Notes, r.CookTime, r.ServingSize, r.Rating, id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe updated"})
}

func deleteRecipe(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	_, err = db.Exec("DELETE FROM recipe WHERE id=?", id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe deleted"})
}

func listRecipes(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM recipe")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		err = rows.Scan(&r.ID, &r.Name, &r.EnoughFor, &r.Origin, &r.Ingredients, &r.Description, &r.Kind, &r.PrepTime, &r.Difficulty, &r.Notes, &r.CookTime, &r.ServingSize, &r.Rating)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		recipes = append(recipes, r)
	}

	err = rows.Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, recipes)
}
