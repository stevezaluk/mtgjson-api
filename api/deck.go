package api

import (
	"errors"
	deckModel "github.com/stevezaluk/mtgjson-models/deck"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"github.com/stevezaluk/mtgjson-sdk/deck"
	"github.com/stevezaluk/mtgjson-sdk/server"
	"mtgjson/auth"
	"net/http"

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

		err = deck.GetDeckContents(server.Database(), results)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to fetch the deck contents for the specified deck", "err": err.Error(), "deckCode": code})
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

		if newDeck.ContentIds != nil { // user submitted a deck with no content ids. Skip as NewDeck will create this structure regardless
			allCards, allCardErr := deck.AllCardIds(newDeck.ContentIds)
			if allCardErr != nil { // this error check is arbitrary as we validate that the deck is not missing a contentIds field
				ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error deck is missing the contentIds field", "err": sdkErrors.ErrDeckMissingContentIds.Error()})
				return
			}

			if len(allCards) != 0 { // skip this block if an empty ContentIds structure was passed, ensuring we don't waste database calls
				// this function needs to be re-added
				var valid, invalidCards, noExistCards = card.ValidateCards(server.Database(), allCards)
				if valid != nil {
					ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for deck creation", "err": valid.Error()})
					return
				}

				if len(invalidCards) != 0 || len(noExistCards) != 0 {
					ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create deck. Some cards do not exist in the database or are invalid", "err": sdkErrors.ErrInvalidCards.Error(), "invalidCards": invalidCards, "noExistCards": noExistCards})
					return
				}
			}
		}

		// add nil check here for contentIds and skip if they are nil. NewDeck will create this automatically

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

		err = deck.GetDeckContents(server.Database(), _deck)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Error fetching deck contents", "err": err.Error()})
		}

		ctx.JSON(http.StatusOK, _deck.Contents)
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

		var updates deckModel.DeckContentIds
		err := ctx.Bind(&updates)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "err": sdkErrors.ErrInvalidObjectStructure.Error()})
			return
		}

		allCards, allCardErr := deck.AllCardIds(&updates)
		if allCardErr != nil { // this should be a sperate error
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error deck is missing the contentIds field", "err": sdkErrors.ErrDeckMissingContentIds.Error()})
			return
		}

		if len(allCards) == 0 { // this should be a seperate error
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "No cards were passed to update the deck with", "err": sdkErrors.ErrDeckMissingId.Error()})
			return
		}

		var valid, invalidCards, noExistCards = card.ValidateCards(server.Database(), allCards)
		if valid != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for deck creation", "err": valid.Error()})
			return
		}

		if len(invalidCards) != 0 || len(noExistCards) != 0 { // need an error for this
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist in the database or are invalid", "err": sdkErrors.ErrInvalidCards.Error(), "invalid": invalidCards, "noExist": noExistCards})
			return
		}

		_deck, err := deck.GetDeck(server.Database(), code, owner)
		if errors.Is(err, sdkErrors.ErrNoDeck) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find deck under the specified deck code", "err": err.Error(), "deckCode": code})
			return
		}

		deck.AddCards(server.Database(), _deck, &updates)

		err = deck.ReplaceDeck(server.Database(), _deck)
		if errors.Is(err, sdkErrors.ErrDeckUpdateFailed) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update deck", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully updated deck", "deckCode": code})
	}
}

/*
DeckContentDELETE Gin handler for the DELETE request to the Deck Content Endpoint. This function should not be called
directly and should only be passed to the gin router
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

		var updates deckModel.DeckContentIds
		err := ctx.Bind(&updates)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error(), "err": sdkErrors.ErrInvalidObjectStructure.Error()})
			return
		}

		allCards, allCardErr := deck.AllCardIds(&updates)
		if allCardErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error deck is missing the contentIds field", "err": sdkErrors.ErrDeckMissingId.Error()})
			return
		}

		if len(allCards) == 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "No cards were passed to update the deck with", "err": sdkErrors.ErrDeckMissingId.Error()})
			return
		}

		// this function needs to be re-added
		var valid, invalidCards, noExistCards = card.ValidateCards(server.Database(), allCards)
		if valid != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for deck creation", "err": valid.Error()})
			return
		}

		if len(invalidCards) != 0 || len(noExistCards) != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist in the database or are invalid", "err": sdkErrors.ErrInvalidCards.Error(), "invalid": invalidCards, "noExist": noExistCards})
			return
		}

		_deck, err := deck.GetDeck(server.Database(), code, owner)
		if errors.Is(err, sdkErrors.ErrNoDeck) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find deck under the specified deck code", "err": err.Error()})
			return
		}

		// need functions for removing cards from the deck here
		deck.RemoveCards(server.Database(), _deck, &updates)

		err = deck.ReplaceDeck(server.Database(), _deck)
		if errors.Is(err, sdkErrors.ErrDeckUpdateFailed) {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to update deck", "err": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Successfully removed cards from deck", "deckCode": code}) // re-add count here
	}
}
