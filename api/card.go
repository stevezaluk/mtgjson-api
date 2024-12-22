package api

import (
	"github.com/gin-gonic/gin"
	cardModel "github.com/stevezaluk/mtgjson-models/card"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"net/http"
	"strconv"
)

/*
limitToInt64 Convert the limit argument from a string to a 64 bit integer
*/
func limitToInt64(limit string) int64 {
	ret, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return int64(100)
	}

	return ret
}

/*
CardGET Gin handler for GET request to the card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardGET(ctx *gin.Context) {
	cardId := ctx.Query("cardId")
	if cardId == "" {
		limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
		results, err := card.IndexCards(limit)
		if err == sdkErrors.ErrNoCards {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusFound, results)
		return
	}

	results, err := card.GetCard(cardId)
	if err == sdkErrors.ErrNoCard {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error(), "cardId": cardId})
		return
	} else if err == sdkErrors.ErrInvalidUUID {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "cardId": cardId})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

/*
CardPOST Gin handler for POST request to the card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardPOST(ctx *gin.Context) {
	var new *cardModel.CardSet

	if ctx.Bind(new) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	err := card.NewCard(new)
	if err == sdkErrors.ErrCardAlreadyExist {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Card already exists under this identifier", "mtgjsonV4Id": new.Identifiers.MtgjsonV4Id})
		return
	} else if err == sdkErrors.ErrCardMissingId {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Card name or mtgjsonV4Id must not be empty when creating a card"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "New card created successfully", "mtgjsonV4Id": new.Identifiers.MtgjsonV4Id})
}

/*
CardDELETE Gin handler for DELETE request to the card endpoint. This should not be called directly and
should only be passed to the gin router
*/

func CardDELETE(ctx *gin.Context) {
	cardId := ctx.Query("cardId")
	if cardId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "A cardId (mtgjsonV4Id) is required to delete a card"})
		return
	}

	err := card.DeleteCard(cardId)
	if err == sdkErrors.ErrNoCard {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find card with the specified id", "mtgjsonV4Id": cardId})
		return
	} else if err == sdkErrors.ErrCardDeleteFailed {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete card. Internal server issue"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Card successfully deleted", "mtgjsonV4Id": cardId})
}
