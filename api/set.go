package api

import (
	"github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/set"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Gin handler for GET request to the set endpoint. This should not be called directly and
should only be passed to the gin router
*/
func SetGET(c *gin.Context) {
	setCode := c.Query("setCode")
	if setCode == "" {
		limit := limitToInt64(c.DefaultQuery("limit", "100"))
		results, err := set.IndexSets(limit)
		if err == errors.ErrNoSet {
			c.JSON(http.StatusBadRequest, gin.H{"message": "No sets available to index"})
			return
		}

		c.JSON(http.StatusOK, results)
		return
	}

	results, err := set.GetSet(setCode)
	if err == errors.ErrNoSet {
		c.JSON(http.StatusNotFound, gin.H{"message": "Failed to find set under the requested Set Code", "setCode": setCode})
		return
	}

	c.JSON(http.StatusOK, results)
}
