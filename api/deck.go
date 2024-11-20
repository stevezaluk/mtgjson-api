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
func DeckGET(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := deck.IndexDecks(limit)
		if err == errors.ErrNoDecks {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusFound, results)
		return
	}

	results, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	deck.GetDeckContents(&results)

	c.JSON(http.StatusFound, results)
}

/*
Gin handler for POST request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckPOST(c *gin.Context) {
	var new deck_model.Deck

	if c.Bind(&new) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if new.Name == "" || new.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck Code or Name must not be empty when creating a deck"})
		return
	}

	var valid, invalidCards, noExistCards = card.ValidateCards(new.AllCardIds())
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	var err = deck.NewDeck(new)
	if err == errors.ErrDeckAlreadyExists {
		c.JSON(http.StatusConflict, gin.H{"message": "Deck already exists under this deck code", "deckCode": new.Code})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully created new deck", "deckCode": new.Code})
}

/*
Gin handler for DELETE request to the deck endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckDELETE(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	result := deck.DeleteDeck(_deck.Code)
	if result == errors.ErrDeckDeleteFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully deleted deck", "deckCode": _deck.Code})
}

/*
Gin handler for GET request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentGET(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to fetch a deck's contents"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	deck.GetDeckContents(&_deck)

	var resp = gin.H{"mainBoard": _deck.Contents.Mainboard, "sideBoard": _deck.Contents.Sideboard, "commander": _deck.Contents.Commander}

	c.JSON(http.StatusFound, resp)
}

/*
Gin handler for POST request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentPOST(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a POST operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deck_model.DeckUpdate
	c.BindJSON(&updates)

	valid, invalidCards, noExistCards := card.ValidateCards(updates.AllCards())
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist in the database or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	_deck.AddCards(updates.MainBoard, deck_model.MAINBOARD)
	_deck.AddCards(updates.SideBoard, deck_model.SIDEBOARD)
	_deck.AddCards(updates.Commander, deck_model.COMMANDER)

	err = deck.ReplaceDeck(_deck)
	if err == errors.ErrDeckUpdateFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully updated deck", "deckCode": code, "count": updates.Count()})
}

/*
Gin handler for DELETE request to the deck content endpoint. This should not be called directly and
should only be passed to the gin router
*/
func DeckContentDELETE(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform a DELETE operation"})
		return
	}

	_deck, err := deck.GetDeck(code)
	if err == errors.ErrNoDeck {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var updates deck_model.DeckUpdate
	c.BindJSON(&updates)

	_deck.DeleteCards(updates.MainBoard, deck_model.MAINBOARD)
	_deck.DeleteCards(updates.SideBoard, deck_model.SIDEBOARD)
	_deck.DeleteCards(updates.Commander, deck_model.COMMANDER)

	err = deck.ReplaceDeck(_deck)
	if err == errors.ErrDeckUpdateFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully removed cards from deck", "deckCode": code, "count": updates.Count()})
}
