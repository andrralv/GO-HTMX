package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func getMangas(context *gin.Context) {
	mangas, err := getAllMangas()
	if err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.IndentedJSON(http.StatusOK, mangas)
}

func mangaById(context *gin.Context) {
	id := context.Param("id")
	m, err := getMangaById(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Manga not found"})
		return
	}

	context.IndentedJSON(http.StatusOK, m)
}

func returnManga(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Missing param"})
		return
	}

	m, err := getMangaById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Manga not found"})
		return
	}

	m.Quantity += 1

	if err := updateManga(m); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.IndentedJSON(http.StatusOK, m)
}

func checkoutManga(context *gin.Context) {
	id, ok := context.GetQuery("id")

	if !ok {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Missing param"})
		return
	}

	m, err := getMangaById(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Manga not found"})
		return
	}

	if m.Quantity <= 0 {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Manga not available"})
		return
	}

	m.Quantity -= 1

	if err := updateManga(m); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.IndentedJSON(http.StatusOK, m)
}

func createManga(context *gin.Context) {
	var newManga manga
	if err := context.BindJSON(&newManga); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := insertManga(newManga); err != nil {
		context.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newManga)
}

func main() {
	initDB()
	defer closeDB()

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = append(config.AllowHeaders, "hx-current-url", "hx-request", "hx-target", "hx-trigger")
	config.AllowMethods = []string{"GET", "POST", "PATCH", "OPTIONS"}
	config.AllowCredentials = true

	router.Use(cors.New(config))

	router.OPTIONS("/*any", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	router.GET("/manga", getMangas)
	router.GET("/manga/:id", mangaById)
	router.POST("/manga", createManga)
	router.PATCH("/checkout", checkoutManga)
	router.PATCH("/return", returnManga)
	router.Run("localhost:8080")
}
