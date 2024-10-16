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
