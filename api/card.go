package api

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	models "mtgjson/models/card"
	"mtgjson/server"
	"net/http"
	"regexp"
)

func CardGET(c *gin.Context) {
	cardId := c.Query("cardId")
	if cardId == "" {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Fetching all cards is not implemented yet"})
		return
	}

	var pattern = `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`
	re := regexp.MustCompile(pattern)
	if !re.MatchString(cardId) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Card ID was not a valid v5 UUID", "cardId": cardId})
		return
	}

	var database = ServerContext.Value("database").(server.Database)

	query := bson.D{{"identifiers.mtgjsonV4Id", cardId}}
	results := database.Find("card", query, &models.CardSet{})
	if results == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Failed to find card with the requested ID", "cardId": cardId})
		return
	}

	c.JSON(http.StatusFound, results)
}
