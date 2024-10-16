package api

import (
	"errors"
	models "mtgjson/models/card"
	"mtgjson/server"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

/*
/card - Represents a card that exists either in a deck or in a set. Must have a unique MTGJSON V4 UUID
*/

var ErrNoCard = errors.New("card: failed to find card with specified uuid")
var ErrInvalidUUID = errors.New("card: invalid v5 uuid")

func validateUuid(uuid string) bool {
	var ret = false
	var pattern = `^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`

	re := regexp.MustCompile(pattern)
	if re.MatchString(uuid) {
		ret = true
	}

	return ret
}

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

func CardGET(c *gin.Context) {
	cardId := c.Query("cardId")
	if cardId == "" {
		c.JSON(http.StatusNotImplemented, gin.H{"message": "Fetching all cards is not implemented yet"})
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
