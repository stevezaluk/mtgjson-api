package api

import (
	"mtgjson/errors"
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

	var valid, invalidCards = new.ValidateCards()
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create deck. Some cards do not exist or are invalid", "uuid": invalidCards})
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

	var mainBoard = _deck.GetMainboard()
	var sideBoard = _deck.GetSideboard()
	var commander = _deck.GetCommander()

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

	type DeckUpdate struct {
		UUID []string
	}

	var updates DeckUpdate
	c.BindJSON(&updates)

	if len(updates.UUID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "UUID is empty. A list of mtgjsonV4 uuid's must be passed to update a deck"})
		return
	}

	valid, invalidCards := _deck.ValidateCards()
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update deck. Some cards do not exist or are invalid", "uuid": invalidCards})
		return
	}

	for i := 0; i < len(updates.UUID); i++ {
		var uuid = updates.UUID[i]
		err = _deck.AddCard(uuid)
		if err == errors.ErrCardAlreadyExist {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
	}

	err = _deck.UpdateDeck()
	if err == errors.ErrDeckUpdateFailed {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Successfully updated deck", "deckCode": code, "count": len(updates.UUID)})
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
