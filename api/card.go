package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	models "mtgjson/models/card"
	"mtgjson/server"
	"net/http"
	"regexp"
	"strconv"
)

/*
utility - Handles any validation or conversion that other functions dont
*/

func validateUuid(uuid string) bool {
	var ret = false
	var pattern = `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

	re := regexp.MustCompile(pattern)
	if re.MatchString(uuid) {
		ret = true
	}

	return ret
}

func limitToInt64(limit string) int64 {
	ret, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return int64(100)
	}

	return ret
}

/*
/card - Represents a card that exists either in a deck or in a set. Must have a unique MTGJSON V4 UUID
*/

var ErrNoCard = errors.New("card: failed to find card with specified uuid")
var ErrNoCards = errors.New("card: No card found on index operation")
var ErrInvalidUUID = errors.New("card: invalid v5 uuid")

func getCard(uuid string) (models.CardSet, error) {
	var result models.CardSet

	if !validateUuid(uuid) {
		return result, ErrInvalidUUID
	}

	var database = ServerContext.Value("database").(server.Database)

	query := bson.M{"identifiers.mtgjsonV4Id": uuid}
	results := database.Find("card", query, &result)
	if results == nil {
		return result, ErrNoCard
	}

	return result, nil
}

func getCards(limit int64) ([]models.CardSet, error) {
	var result []models.CardSet

	var database = ServerContext.Value("database").(server.Database)

	results := database.Index("card", limit, &result)
	if results == nil {
		return result, ErrNoCards
	}

	return result, nil

}

func CardGET(c *gin.Context) {
	cardId := c.Query("cardId")
	if cardId == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := getCards(limit)
		if err == ErrNoCards {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusFound, results)
		return
	}

	results, err := getCard(cardId)
	if err == ErrNoCard {
		c.JSON(404, gin.H{"message": err.Error(), "cardId": cardId})
		return
	} else if err == ErrInvalidUUID {
		c.JSON(400, gin.H{"message": err.Error(), "cardId": cardId})
		return
	}

	c.JSON(http.StatusFound, results)
}
