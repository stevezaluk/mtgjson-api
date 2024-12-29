package api

import (
	"errors"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	"github.com/stevezaluk/mtgjson-sdk/set"
	"mtgjson/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
SetGET Gin handler for GET request to the set endpoint. This should not be called directly and
should only be passed to the gin router
*/
func SetGET(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner != "system" && owner != userEmail {
		if !auth.ValidateScope(ctx, "read:user-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other users sets", "requiredScope": "read:user-set"})
			return
		}
	}

	setCode := ctx.Query("setCode")
	if setCode == "" {
		limit := limitToInt64(ctx.DefaultQuery("limit", "100"))
		results, err := set.IndexSets(limit)
		if errors.Is(err, sdkErrors.ErrNoSet) {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "No sets available to index"})
			return
		}

		ctx.JSON(http.StatusOK, results)
		return
	}

	results, err := set.GetSet(setCode, owner)
	if errors.Is(err, sdkErrors.ErrNoSet) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Failed to find set under the requested Set Code", "setCode": setCode})
		return
	}

	err = set.GetSetContents(results)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	ctx.JSON(http.StatusOK, results)
}
