package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	cardModel "github.com/stevezaluk/mtgjson-models/card"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"mtgjson/auth"
	"net/http"
	"strconv"
)

/*
limitToInt64 Convert the limit argument from a string to a 64-bit integer
*/
func limitToInt64(limit string) int64 {
	ret, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		return int64(100)
	}

	return ret
}

/*
CardGET Gin handler for GET request to the Card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardGET(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner != "system" && owner != userEmail {
		if !auth.ValidateScope(ctx, "read:user-card") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other users cards", "requiredScope": "read:user-card"})
			return
		}
	}

	cardId := ctx.Query("cardId")
	if cardId == "" {
		limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
		results, err := card.IndexCards(limit) // update this function with owner
		if errors.Is(err, sdkErrors.ErrNoCards) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusFound, results)
		return
	}

	results, err := card.GetCard(cardId, owner)
	if errors.Is(err, sdkErrors.ErrNoCard) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error(), "cardId": cardId})
		return
	} else if errors.Is(err, sdkErrors.ErrInvalidUUID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "cardId": cardId})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

/*
CardPOST Gin handler for POST request to the Card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardPOST(ctx *gin.Context) {
	var newCard *cardModel.CardSet

	err := ctx.Bind(&newCard)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect", "err": err.Error()})
		return
	}

	if newCard.Identifiers == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The identifiers field must be filled. At minimum an mtgjsonV4Id is required"})
		return
	}

	if newCard.Name == "" || newCard.Identifiers.MtgjsonV4Id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Either the name or mtgjsonV4Id field is empty. Both of these fields must be filled in"})
		return
	}

	if newCard.MtgjsonApiMeta != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The mtgjsonApiMeta must be null. This will be filled in automatically during card creation"})
		return
	}

	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner == "system" {
		if !auth.ValidateScope(ctx, "write:system-card") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify of system or pre-constructed cards", "requiredScope": "write:system-card"})
			return
		}
	}

	if owner != userEmail {
		if !auth.ValidateScope(ctx, "write:user-card") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users card's", "requiredScope": "write:user-card"})
			return
		}
	}

	err = card.NewCard(newCard, owner)
	if errors.Is(err, sdkErrors.ErrCardAlreadyExist) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Card already exists under this identifier", "mtgjsonV4Id": newCard.Identifiers.MtgjsonV4Id})
		return
	} else if errors.Is(err, sdkErrors.ErrCardMissingId) {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Card name or mtgjsonV4Id must not be empty when creating a card"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "New card created successfully", "mtgjsonV4Id": newCard.Identifiers.MtgjsonV4Id})
}

/*
CardDELETE Gin handler for DELETE request to the Card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardDELETE(ctx *gin.Context) {
	cardId := ctx.Query("cardId")
	if cardId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "A cardId (mtgjsonV4Id) is required to delete a card"})
		return
	}

	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner == "system" {
		if !auth.ValidateScope(ctx, "write:system-card") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify of system or pre-constructed cards", "requiredScope": "write:system-card"})
			return
		}
	}

	if owner != userEmail {
		if !auth.ValidateScope(ctx, "write:user-deck") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users cards", "requiredScope": "write:user-card"})
			return
		}
	}

	err := card.DeleteCard(cardId, owner)
	if errors.Is(err, sdkErrors.ErrNoCard) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find card with the specified id", "mtgjsonV4Id": cardId})
		return
	} else if errors.Is(err, sdkErrors.ErrCardDeleteFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete card. Internal server issue"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Card successfully deleted", "mtgjsonV4Id": cardId})
}
