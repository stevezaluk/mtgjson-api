package api

import (
	"mtgjson/errors"
	"mtgjson/models/set"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetGET(c *gin.Context) {
	setCode := c.Query("setCode")
	if setCode == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := set.IndexSets(limit)
		if err == errors.ErrNoSet {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		c.JSON(http.StatusFound, results)
		return
	}

	results, err := set.GetSet(setCode)
	if err == errors.ErrNoSet {
		c.JSON(400, gin.H{"message": err.Error(), "setCode": setCode})
		return
	}

	c.JSON(http.StatusFound, results)
}
