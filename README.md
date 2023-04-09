# Recipes API

A simple API for storing and retrieving recipes.

## Prerequisites

- Go
- MySQL

## Getting Started

1. Clone the repository
2. Copy `config.sample.yaml` to `config.yaml` and fill in your MySQL database details
3. Run the following command to start the server:

```bash
go run .
```

## Endpoints

### GET /api/recipe/

Lists all recipes.

```bash
curl --location --request GET 'http://localhost:8080/api/recipe'
```

### POST /api/recipe/

Creates a new recipe.

```bash
curl --location --request POST 'http://localhost:8080/api/recipe' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Spaghetti Carbonara",
    "enough_for": "4 servings",
    "origin": "Italy",
    "ingredients": "spaghetti, eggs, bacon, parmesan cheese, salt, black pepper",
    "description": "A classic Italian pasta dish",
    "kind": "pasta",
    "prep_time": "10 minutes",
    "difficulty": "easy",
    "notes": "",
    "cook_time": "15 minutes",
    "serving_size": "",
    "rating": "4.8"
}'
```

### GET /api/recipe/:id

Gets a single recipe by ID.

```bash
curl --location --request GET 'http://localhost:8080/api/recipe/1'
```

### PUT /api/recipe/:id

Updates a recipe by ID.

```bash
curl --location --request PUT 'http://localhost:8080/api/recipe/1' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name": "Spaghetti Carbonara",
    "enough_for": "4 servings",
    "origin": "Italy",
    "ingredients": "spaghetti, eggs, bacon, parmesan cheese, salt, black pepper",
    "description": "A classic Italian pasta dish with a creamy sauce",
    "kind": "pasta",
    "prep_time": "10 minutes",
    "difficulty": "easy",
    "notes": "",
    "cook_time": "15 minutes",
    "serving_size": "",
    "rating": "4.8"
}'
```

### DELETE /api/recipe/:id

Deletes a recipe by ID.

```bash
curl --location --request DELETE 'http://localhost:8080/api/recipe/1'
```

## Credits

Created by Joe Harkins ([jharkins](https://github.com/jharkins)) and [OpenAI's ChatGPT](https://www.openai.com/).
