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

func DeckContentPUT(c *gin.Context) {
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

	for i := 0; i < len(updates.UUID); i++ {
		var uuid = updates.UUID[i]
		err = _deck.AddCard(uuid)
		if err == errors.ErrCardAlreadyExist {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		} else if err == errors.ErrNoCard {
			c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
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
