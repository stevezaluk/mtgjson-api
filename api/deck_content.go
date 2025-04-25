package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	deckModel "github.com/stevezaluk/mtgjson-models/deck"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"github.com/stevezaluk/mtgjson-sdk/deck"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"mtgjson/auth"
	"net/http"
)

/*
DeckContentGET Gin handler for the GET request to the Deck Content Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func DeckContentGET(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)
		if owner != "system" && owner != userEmail { // caller is trying to read the contents of another users deck
			if !auth.ValidateScope(ctx, "read:deck.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other users decks", "requiredScope": "read:user-deck"})
				return
			}
		}

		code := ctx.Query("deckCode")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to fetch a deck's contents", "err": sdkErrors.ErrDeckMissingId.Error()})
			return
		}

		_deck, err := deck.GetDeck(server.Database(), code, owner)
		if errors.Is(err, sdkErrors.ErrNoDeck) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find deck with the specified deck code", "err": err.Error()})
			return
		}

		contents, err := deck.GetDeckContents(server.Database(), _deck)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Error fetching deck contents", "err": err.Error()})
		}

		ctx.JSON(http.StatusOK, contents)
	}
}

/*
DeckContentPOST Gin handler for the POST request to the Deck Content Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func DeckContentPOST(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)

		if owner == "system" {
			if !auth.ValidateScope(ctx, "write:deck.wotc") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify content of system or pre-constructed decks", "requiredScope": "write:system-deck"})
				return
			}
		}

		if owner != userEmail {
			if !auth.ValidateScope(ctx, "write:deck.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users deck content", "requiredScope": "write:user-deck"})
				return
			}
		}

		code := ctx.Query("deckCode")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a POST operation", "err": sdkErrors.ErrDeckMissingId.Error()})
			return
		}

		requestedDeck, err := deck.GetDeck(server.Database(), code, owner)
		if errors.Is(err, sdkErrors.ErrNoDeck) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find deck with the specified deck code", "err": err.Error()})
			return
		}

		var request deckModel.DeckContentIds

		err = ctx.Bind(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind request. Structure may be invalid", "err": err.Error()})
			return
		}

		allCardIds := deck.AllCardIds(requestedDeck.Contents)
		err, invalidCards, noExistCards := card.ValidateCards(server.Database(), allCardIds)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update deck. Error while validating cards", "err": err.Error()})
			return
		}

		if len(invalidCards) != 0 || len(noExistCards) != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards are invalid or do not exist", "err": sdkErrors.ErrInvalidCards, "invalidCards": invalidCards, "noExistCards": noExistCards})
			return
		}

		err = deck.AddCards(server.Database(), requestedDeck, &request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to add cards to deck", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully updated deck", "deckCode": code})
	}
}

/*
DeckContentDELETE Gin handler for the DELETE request to the Deck Content Endpoint. This function should not be called
directly and should only be passed to the gin router.
*/
func DeckContentDELETE(server *server.Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userEmail := ctx.GetString("userEmail")
		owner := ctx.DefaultQuery("owner", userEmail)

		if owner == "system" {
			if !auth.ValidateScope(ctx, "write:deck.wotc") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify content of system or pre-constructed decks", "requiredScope": "write:system-deck"})
				return
			}
		}

		if owner != userEmail {
			if !auth.ValidateScope(ctx, "write:deck.admin") {
				ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users deck content", "requiredScope": "write:user-deck"})
				return
			}
		}

		code := ctx.Query("deckCode")
		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation", "err": sdkErrors.ErrDeckMissingId.Error()})
			return
		}

		requestedDeck, err := deck.GetDeck(server.Database(), code, owner)
		if errors.Is(err, sdkErrors.ErrNoDeck) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find deck with the specified deck code", "err": err.Error()})
			return
		}

		var request deckModel.DeckContentIds

		err = ctx.Bind(&request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind request. Structure may be invalid", "err": err.Error()})
			return
		}

		allCardIds := deck.AllCardIds(requestedDeck.Contents)
		err, invalidCards, noExistCards := card.ValidateCards(server.Database(), allCardIds)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update deck. Error while validating cards", "err": err.Error()})
			return
		}

		if len(invalidCards) != 0 || len(noExistCards) != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards are invalid or do not exist", "err": sdkErrors.ErrInvalidCards, "invalidCards": invalidCards, "noExistCards": noExistCards})
			return
		}

		err = deck.RemoveCards(server.Database(), requestedDeck, &request)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to remove cards from deck", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully removed cards from deck", "deckCode": code}) // re-add count here
	}
}
