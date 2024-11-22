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
func SetGET(ctx *gin.Context) {
	setCode := ctx.Query("setCode")
	if setCode == "" {
		limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
		results, err := set.IndexSets(limit)
		if err == errors.ErrNoSet {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "No sets available to index"})
			return
		}

		ctx.JSON(http.StatusOK, results)
		return
	}

	results, err := set.GetSet(setCode)
	if err == errors.ErrNoSet {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find set under the requested Set Code", "setCode": setCode})
		return
	}

	ctx.JSON(http.StatusOK, results)
}
