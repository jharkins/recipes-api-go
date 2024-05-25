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

func getRandomRecipe(c *gin.Context) {
	row := db.QueryRow("SELECT * FROM recipe ORDER BY RAND() LIMIT 1")
	var r Recipe
	err := row.Scan(&r.ID, &r.Name, &r.EnoughFor, &r.Origin, &r.Ingredients, &r.Description, &r.Kind, &r.PrepTime, &r.Difficulty, &r.Notes, &r.CookTime, &r.ServingSize, &r.Rating)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No recipes found"})
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
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	if err != nil || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	offset := (page - 1) * pageSize
	rows, err := db.Query("SELECT * FROM recipe WHERE name <> '' ORDER BY recipe.name LIMIT ? OFFSET ?", pageSize, offset)
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

	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM recipe").Scan(&totalCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"recipes":  recipes,
		"total":    totalCount,
		"page":     page,
		"pageSize": pageSize,
	})
}

func searchRecipesHandler(c *gin.Context) {
	// Extract the search term from the query string
	term := c.Query("term")
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing search term"})
		return
	}

	// Prepare and execute the SQL statement
	stmt := `SELECT * FROM recipe WHERE name LIKE ? OR ingredients LIKE ? OR description LIKE ?`
	rows, err := db.Query(stmt, "%"+term+"%", "%"+term+"%", "%"+term+"%")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var recipes []Recipe
	for rows.Next() {
		var r Recipe
		if err := rows.Scan(&r.ID, &r.Name, &r.EnoughFor, &r.Origin, &r.Ingredients, &r.Description, &r.Kind, &r.PrepTime, &r.Difficulty, &r.Notes, &r.CookTime, &r.ServingSize, &r.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		recipes = append(recipes, r)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send the resulting recipes as JSON
	c.JSON(http.StatusOK, recipes)
}
