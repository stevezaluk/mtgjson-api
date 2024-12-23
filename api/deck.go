package api

import (
	"errors"
	deckModel "github.com/stevezaluk/mtgjson-models/deck"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"github.com/stevezaluk/mtgjson-sdk/deck"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
DeckGET Gin handler for GET request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckGET(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
		results, err := deck.IndexDecks(limit)
		if errors.Is(err, sdkErrors.ErrNoDecks) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, results)
		return
	}

	results, err := deck.GetDeck(code)
	if errors.Is(err, sdkErrors.ErrNoDeck) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	err = deck.GetDeckContents(results)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	ctx.JSON(http.StatusOK, results)
}

/*
DeckPOST Gin handler for POST request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckPOST(ctx *gin.Context) {
	var newDeck *deckModel.Deck

	if ctx.Bind(&newDeck) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if newDeck.Name == "" || newDeck.Code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck Code or Name must not be empty when creating a deck"})
		return
	}

	allCards, allCardErr := deck.AllCardIds(newDeck.ContentIds)
	if allCardErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error deck is missing the contentIds field"})
		return
	}

	// this function needs to be re-added
	var valid, invalidCards, noExistCards = card.ValidateCards(allCards)
	if valid != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for deck creation", "err": valid.Error()})
		return
	}

	if len(invalidCards) != 0 || len(noExistCards) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	var err = deck.NewDeck(newDeck)
	if errors.Is(err, sdkErrors.ErrDeckAlreadyExists) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Deck already exists under this deck code", "deckCode": newDeck.Code})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully created new deck", "deckCode": newDeck.Code})
}

/*
DeckDELETE Gin handler for DELETE request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckDELETE(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if errors.Is(err, sdkErrors.ErrNoDeck) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	result := deck.DeleteDeck(_deck.Code)
	if errors.Is(result, sdkErrors.ErrDeckDeleteFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted deck", "deckCode": _deck.Code})
}

/*
DeckContentGET Gin handler for GET request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentGET(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to fetch a deck's contents"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if errors.Is(err, sdkErrors.ErrNoDeck) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	err = deck.GetDeckContents(_deck)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Error fetching deck contents", "err": err.Error()})
	}

	ctx.JSON(http.StatusOK, _deck.Contents)
}

/*
DeckContentPOST Gin handler for POST request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentPOST(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a POST operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if errors.Is(err, sdkErrors.ErrNoDeck) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deckModel.DeckContentIds
	if ctx.BindJSON(&updates) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	allCards, allCardErr := deck.AllCardIds(&updates)
	if allCardErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error deck is missing the contentIds field"})
		return
	}

	// this function needs to be re-added
	var valid, invalidCards, noExistCards = card.ValidateCards(allCards)
	if valid != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for deck creation", "err": valid.Error()})
		return
	}

	if len(invalidCards) != 0 || len(noExistCards) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	// need functions for adding to the deck here
	deck.AddCards(_deck, &updates)

	err = deck.ReplaceDeck(_deck)
	if errors.Is(err, sdkErrors.ErrDeckUpdateFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully updated deck", "deckCode": code}) // re-add count here
}

/*
DeckContentDELETE Gin handler for DELETE request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentDELETE(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if errors.Is(err, sdkErrors.ErrNoDeck) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deckModel.DeckContentIds
	if ctx.BindJSON(&updates) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	allCards, allCardErr := deck.AllCardIds(&updates)
	if allCardErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Error deck is missing the contentIds field"})
		return
	}

	// this function needs to be re-added
	var valid, invalidCards, noExistCards = card.ValidateCards(allCards)
	if valid != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for deck creation", "err": valid.Error()})
		return
	}

	if len(invalidCards) != 0 || len(noExistCards) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	// need functions for removing cards from the deck here
	deck.RemoveCards(_deck, &updates)

	err = deck.ReplaceDeck(_deck)
	if errors.Is(err, sdkErrors.ErrDeckUpdateFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully removed cards from deck", "deckCode": code}) // re-add count here
}
