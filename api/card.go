package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	cardModel "github.com/stevezaluk/mtgjson-models/card"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"github.com/stevezaluk/mtgjson-sdk/server"
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
func CardGET(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)

		if owner != "system" && owner != userEmail {
			if !auth.ValidateScope(ctx, "read:card.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other users cards", "requiredScope": "read:user-card"})
				return
			}
		}

		cardId := ctx.Query("cardId")
		if cardId == "" {
			limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
			results, err := card.IndexCards(server.Database(), limit) // update this function with owner
			if errors.Is(err, sdkErrors.ErrNoCards) {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to find cards in the database to index", "err": err.Error()})
				return
			}

			ctx.JSON(http.StatusOK, results)
			return
		}

		results, err := card.GetCard(server.Database(), cardId, owner)
		if errors.Is(err, sdkErrors.ErrNoCard) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find card with specified cardId", "err": err.Error(), "cardId": cardId})
			return
		} else if errors.Is(err, sdkErrors.ErrInvalidUUID) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "cardId is not a valid V5 UUID", "err": err.Error(), "cardId": cardId})
			return
		}

		ctx.JSON(http.StatusOK, results)
	}
}

/*
CardPOST Gin handler for POST request to the Card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardPOST(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)

		if owner == "system" {
			if !auth.ValidateScope(ctx, "write:card.wotc") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify of system or pre-constructed cards", "requiredScope": "write:system-card"})
				return
			}
		}

		if owner != userEmail {
			if !auth.ValidateScope(ctx, "write:card.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users card's", "requiredScope": "write:user-card"})
				return
			}
		}

		var newCard *cardModel.CardSet

		err := ctx.Bind(&newCard)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "err": sdkErrors.ErrInvalidObjectStructure.Error()})
			return
		}

		if newCard.Identifiers == nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "The identifiers field must be filled. At minimum an mtgjsonV4Id is required", "err": sdkErrors.ErrCardMissingId.Error()})
			return
		}

		if newCard.Name == "" || newCard.Identifiers.MtgjsonV4Id == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Either the name or mtgjsonV4Id field is empty. Both of these fields must be filled in", "err": sdkErrors.ErrCardMissingId.Error()})
			return
		}

		if newCard.MtgjsonApiMeta != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "The mtgjsonApiMeta must be null. This will be filled in automatically during card creation", "err": sdkErrors.ErrMetaApiMustBeNull.Error()}) // need seperate error for this
			return
		}

		err = card.NewCard(server.Database(), newCard, owner)
		if errors.Is(err, sdkErrors.ErrCardAlreadyExist) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "Card already exists under this identifier", "err": err.Error(), "cardId": newCard.Identifiers.MtgjsonV4Id})
			return
		} else if errors.Is(err, sdkErrors.ErrCardMissingId) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Card name or mtgjsonV4Id must not be empty when creating a card", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "New card created successfully", "cardId": newCard.Identifiers.MtgjsonV4Id})
	}
}

/*
CardDELETE Gin handler for DELETE request to the Card endpoint. This should not be called directly and
should only be passed to the gin router
*/
func CardDELETE(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)

		if owner == "system" {
			if !auth.ValidateScope(ctx, "write:deck.wotc") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify of system or pre-constructed cards", "requiredScope": "write:system-card"})
				return
			}
		}

		if owner != userEmail {
			if !auth.ValidateScope(ctx, "write:deck.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users cards", "requiredScope": "write:user-card"})
				return
			}
		}

		cardId := ctx.Query("cardId")
		if cardId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "A cardId (mtgjsonV4Id) is required to delete a card", "err": sdkErrors.ErrCardMissingId.Error()})
			return
		}

		err := card.DeleteCard(server.Database(), cardId, owner)
		if errors.Is(err, sdkErrors.ErrNoCard) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find card with the specified id", "err": err.Error(), "cardId": cardId})
			return
		} else if errors.Is(err, sdkErrors.ErrCardDeleteFailed) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete card. Internal server issue", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Card successfully deleted", "cardId": cardId})
	}
}
