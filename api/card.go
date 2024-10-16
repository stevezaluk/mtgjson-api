package api

import (
	"github.com/gin-gonic/gin"
	"mtgjson/errors"
	"mtgjson/models/card"
	"net/http"
	"strconv"
)

/*
utility - Handles any validation or conversion that other functions dont
*/

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

func CardGET(c *gin.Context) {
	cardId := c.Query("cardId")
	if cardId == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := card.GetCards(limit)
		if err == errors.ErrNoCards {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusFound, results)
		return
	}

	results, err := card.GetCard(cardId)
	if err == errors.ErrNoCard {
		c.JSON(404, gin.H{"message": err.Error(), "cardId": cardId})
		return
	} else if err == errors.ErrInvalidUUID {
		c.JSON(400, gin.H{"message": err.Error(), "cardId": cardId})
		return
	}

	c.JSON(http.StatusFound, results)
}

func CardPOST(c *gin.Context) {

}
