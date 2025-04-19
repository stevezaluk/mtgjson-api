package api

import (
	"errors"
	deckModel "github.com/stevezaluk/mtgjson-models/deck"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"github.com/stevezaluk/mtgjson-sdk/deck"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"maps"
	"mtgjson/auth"
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
)

/*
DeckGET Gin handler for the GET request to the Deck Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func DeckGET(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)

		if owner != "system" && owner != userEmail { // caller is trying to read another users deck
			if !auth.ValidateScope(ctx, "read:deck.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other users decks", "requiredScope": "read:user-deck"})
				return
			}
		}

		code := ctx.Query("deckCode")
		if code == "" {
			limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
			results, err := deck.IndexDecks(server.Database(), limit) // update this function with owner
			if errors.Is(err, sdkErrors.ErrNoDecks) {
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to find decks in the database to index", "err": err.Error()})
				return
			}

			ctx.JSON(http.StatusOK, results)
			return
		}

		results, err := deck.GetDeck(server.Database(), code, owner)
		if errors.Is(err, sdkErrors.ErrNoDeck) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find deck under the specified deck code", "err": err.Error(), "deckCode": code})
			return
		}

		ctx.JSON(http.StatusOK, results)
	}
}

/*
DeckPOST Gin handler for the POST request to the Deck Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func DeckPOST(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)

		if owner == "system" {
			if !auth.ValidateScope(ctx, "write:deck.wotc") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify of system or pre-constructed decks", "requiredScope": "write:system-deck"})
				return
			}
		}

		if owner != userEmail {
			if !auth.ValidateScope(ctx, "write:deck.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users deck content", "requiredScope": "write:user-deck"})
				return
			}
		}

		var newDeck *deckModel.Deck

		if ctx.Bind(&newDeck) != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
			return
		}

		if newDeck.MtgjsonApiMeta != nil { // need an error for this
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "The mtgjsonApiMeta field must be null. This will be filled out automatically during deck creation", "err": sdkErrors.ErrMetaApiMustBeNull.Error()})
			return
		}

		var cardValidate []string

		if newDeck.MainBoard != nil {
			cardValidate = append(cardValidate, slices.Collect(maps.Keys(newDeck.MainBoard))...)
		}

		if newDeck.SideBoard != nil {
			cardValidate = append(cardValidate, slices.Collect(maps.Keys(newDeck.SideBoard))...)
		}

		if newDeck.Commander != nil {
			cardValidate = append(cardValidate, slices.Collect(maps.Keys(newDeck.Commander))...)
		}

		isValid, invalidCards, noExistCards := card.ValidateCards(server.Database(), cardValidate)
		if isValid != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to validate the cards of the deck", "err": isValid.Error()})
			return
		}

		if len(invalidCards) != 0 || len(noExistCards) != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid cards were detected in the deck model", "err": sdkErrors.ErrInvalidCards, "invalidCards": invalidCards, "noExistCards": noExistCards})
			return
		}

		// handle deck contents here

		var err = deck.NewDeck(server.Database(), newDeck, owner)
		if errors.Is(err, sdkErrors.ErrDeckMissingId) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck is missing a name and/or a deck code. Both of these values must be filled", "err": err.Error()})
			return
		} else if errors.Is(err, sdkErrors.ErrDeckAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"message": "Deck already exists under this deck code", "err": err.Error(), "deckCode": newDeck.Code})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully created new deck", "deckCode": newDeck.Code})
	}
}

/*
DeckDELETE Gin handler for the DELETE request to the Deck Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func DeckDELETE(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)
		if owner == "system" { // caller is trying to delete a system created (pre-constructed) deck
			if !auth.ValidateScope(ctx, "write:deck.wotc") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to delete pre-constructed decks", "requiredScope": "write:system-deck"})
				return
			}
		}

		if owner != userEmail { // caller is trying to delete a different users deck
			if !auth.ValidateScope(ctx, "write:deck.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to delete other users decks", "requiredScope": "write:user-deck"})
				return
			}
		}

		code := ctx.Query("deckCode")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation", "err": sdkErrors.ErrDeckMissingId.Error()})
			return
		}

		_deck, err := deck.GetDeck(server.Database(), code, owner)
		if errors.Is(err, sdkErrors.ErrNoDeck) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find deck under the specified deck code", "err": err.Error(), "deckCode": code})
			return
		}

		result := deck.DeleteDeck(server.Database(), _deck.Code, owner)
		if errors.Is(result, sdkErrors.ErrDeckDeleteFailed) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Delete deck operation has failed", "err": err.Error(), "deckCode": code})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted deck", "deckCode": _deck.Code})
	}
}
