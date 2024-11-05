package api

import (
	"github.com/gin-gonic/gin"
	card_model "github.com/stevezaluk/mtgjson-models/card"
	"github.com/stevezaluk/mtgjson-models/errors"
	card "github.com/stevezaluk/mtgjson-sdk/card"
	"net/http"
	"strconv"
)

/*
limitToInt64 - Convert the limit query arg to an Integer

Parameters:
limit (string) - The limit as a string

Returns:
ret (int64) - The limit as an integer
*/

func limitToInt64(limit string) int64 {
	ret, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return int64(100)
	}

	return ret
}

/*
CardGET - All functionality for fetching card metadata

Parameters:
c (gin.Context) - The request context

Returns:
Nothing
*/
func CardGET(c *gin.Context) {
	cardId := c.Query("cardId")
	if cardId == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := card.IndexCards(limit)
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
	var new card_model.Card

	if c.Bind(&new) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	err := card.NewCard(new)
	if err == errors.ErrCardAlreadyExist {
		c.JSON(http.StatusConflict, gin.H{"message": "Card already exists under this identifier", "mtgjsonV4Id": new.Identifiers.MTGJsonV4Id})
		return
	} else if err == errors.ErrCardMissingId {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Card name or mtgjsonV4Id must not be empty when creating a card"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "New card created successfully", "mtgjsonV4Id": new.Identifiers.MTGJsonV4Id})
}
