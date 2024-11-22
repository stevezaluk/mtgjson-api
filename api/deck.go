package api

import (
	deck_model "github.com/stevezaluk/mtgjson-models/deck"
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/card"
	"github.com/stevezaluk/mtgjson-sdk/deck"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Gin handler for GET request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckGET(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
		results, err := deck.IndexDecks(limit)
		if err == errors.ErrNoDecks {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, results)
		return
	}

	results, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	deck.GetDeckContents(&results)

	ctx.JSON(http.StatusOK, results)
}

/*
Gin handler for POST request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckPOST(ctx *gin.Context) {
	var new deck_model.Deck

	if ctx.Bind(&new) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if new.Name == "" || new.Code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck Code or Name must not be empty when creating a deck"})
		return
	}

	var valid, invalidCards, noExistCards = card.ValidateCards(new.AllCardIds())
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	var err = deck.NewDeck(new)
	if err == errors.ErrDeckAlreadyExists {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Deck already exists under this deck code", "deckCode": new.Code})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully created new deck", "deckCode": new.Code})
}

/*
Gin handler for DELETE request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckDELETE(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	result := deck.DeleteDeck(_deck.Code)
	if result == errors.ErrDeckDeleteFailed {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted deck", "deckCode": _deck.Code})
}

/*
Gin handler for GET request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentGET(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to fetch a deck's contents"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	deck.GetDeckContents(&_deck)

	var resp = gin.H{"mainBoard": _deck.Contents.Mainboard, "sideBoard": _deck.Contents.Sideboard, "commander": _deck.Contents.Commander}

	ctx.JSON(http.StatusOK, resp)
}

/*
Gin handler for POST request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentPOST(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a POST operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deck_model.DeckUpdate
	ctx.BindJSON(&updates)

	valid, invalidCards, noExistCards := card.ValidateCards(updates.AllCards())
	if !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	_deck.AddCards(updates.MainBoard, deck_model.MAINBOARD)
	_deck.AddCards(updates.SideBoard, deck_model.SIDEBOARD)
	_deck.AddCards(updates.Commander, deck_model.COMMANDER)

	err = deck.ReplaceDeck(_deck)
	if err == errors.ErrDeckUpdateFailed {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully updated deck", "deckCode": code, "count": updates.Count()})
}

/*
Gin handler for DELETE request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentDELETE(ctx *gin.Context) {
	code := ctx.Query("deckCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deck_model.DeckUpdate
	ctx.BindJSON(&updates)

	_deck.DeleteCards(updates.MainBoard, deck_model.MAINBOARD)
	_deck.DeleteCards(updates.SideBoard, deck_model.SIDEBOARD)
	_deck.DeleteCards(updates.Commander, deck_model.COMMANDER)

	err = deck.ReplaceDeck(_deck)
	if err == errors.ErrDeckUpdateFailed {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully removed cards from deck", "deckCode": code, "count": updates.Count()})
}
