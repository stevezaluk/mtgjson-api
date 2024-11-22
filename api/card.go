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
Convert the limit argument from a string to a 64 bit integer
*/

func limitToInt64(limit string) int64 {
	ret, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return int64(100)
	}

	return ret
}

/*
Gin handler for GET request to the card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardGET(c *gin.Context) {
	cardId := c.Query("cardId")
	if cardId == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := card.IndexCards(limit)
		if err == errors.ErrNoCards {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusFound, results)
		return
	}

	results, err := card.GetCard(cardId)
	if err == errors.ErrNoCard {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error(), "cardId": cardId})
		return
	} else if err == errors.ErrInvalidUUID {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "cardId": cardId})
		return
	}

	c.JSON(http.StatusOK, results)
}

/*
Gin handler for POST request to the card endpoint. This should not be called directly and
should only be passed to the gin router
*/
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

	c.JSON(http.StatusOK, gin.H{"message": "New card created successfully", "mtgjsonV4Id": new.Identifiers.MTGJsonV4Id})
}

/*
Gin handler for DELETE request to the card endpoint. This should not be called directly and
should only be passed to the gin router
*/

func CardDELETE(c *gin.Context) {
	cardId := c.Query("cardId")
	if cardId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "A cardId (mtgjsonV4Id) is required to delete a card"})
		return
	}

	err := card.DeleteCard(cardId)
	if err == errors.ErrNoCard {
		c.JSON(http.StatusNotFound, gin.H{"message": "Failed to find card with the specified id", "mtgjsonV4Id": cardId})
		return
	} else if err == errors.ErrCardDeleteFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete card. Internal server issue"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Card successfully deleted", "mtgjsonV4Id": cardId})
}
