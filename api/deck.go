package api

import (
	"mtgjson/errors"
	"mtgjson/models/card"
	"mtgjson/models/deck"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DeckGET(c *gin.Context) {
	code := c.Query("deckCode")
	if code == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := deck.GetDecks(limit)
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

	c.JSON(http.StatusFound, results)
}

func DeckPOST(c *gin.Context) {
	var new deck.Deck

	if c.Bind(&new) != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if new.Name == "" || new.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "The name or deck code is missing from the request"})
		return
	}

	var valid, invalidCards, noExistCards = card.ValidateCards(new.AllCards())
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create deck. Some cards do not exist or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	var err = deck.NewDeck(new)
	if err == errors.ErrDeckAlreadyExists {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Deck already exists under this deck code", "deckCode": new.Code})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully created new deck", "deckCode": new.Code})
}

func DeckDELETE(c *gin.Context) {
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

	result := _deck.DeleteDeck()
	if result == errors.ErrDeckDeleteFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully deleted deck", "deckCode": _deck.Code})
}

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

	var mainBoard = _deck.FetchMainboard()
	var sideBoard = _deck.FetchSideboard()
	var commander = _deck.FetchCommander()

	var resp = gin.H{"mainBoard": mainBoard, "sideBoard": sideBoard, "commander": commander}

	c.JSON(http.StatusFound, resp)
}

func DeckContentPOST(c *gin.Context) {
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

	var updates deck.DeckUpdate
	c.BindJSON(&updates)

	valid, invalidCards, noExistCards := card.ValidateCards(updates.AllCards())
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist or are invalid", "invalid": invalidCards, "noExist": noExistCards})
		return
	}

	_deck.AddCards(updates.MainBoard, deck.MAINBOARD)
	_deck.AddCards(updates.SideBoard, deck.SIDEBOARD)
	_deck.AddCards(updates.Commander, deck.COMMANDER)

	err = _deck.UpdateDeck()
	if err == errors.ErrDeckUpdateFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully updated deck", "deckCode": code, "count": updates.Count()})
}

func DeckContentDELETE(c *gin.Context) {
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

	type DeckUpdate struct {
		UUID []string
	}

	var updates DeckUpdate
	c.BindJSON(&updates)

	for i := range updates.UUID {
		var uuid = updates.UUID[i]

		err = _deck.DeleteCard(uuid)
		if err == errors.ErrNoCard {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		}
	}

	err = _deck.UpdateDeck()
	if err == errors.ErrDeckUpdateFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully removed cards from deck", "deckCode": code, "count": len(updates.UUID)})
}
