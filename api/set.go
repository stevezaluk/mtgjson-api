package api

import (
	"errors"
	sdkErrors "github.com/stevezaluk/mtgjson-models/errors"
	setModel "github.com/stevezaluk/mtgjson-models/set"
	"github.com/stevezaluk/mtgjson-sdk/card"
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
		results, err := set.IndexSets(limit) // update this to include ownership
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

/*
SetPOST Gin handler for the POST request to the Set Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func SetPOST(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner == "system" {
		if !auth.ValidateScope(ctx, "write:system-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify system or pre-constructed sets", "requiredScope": "write:system-set"})
			return
		}
	}

	if owner != userEmail {
		if !auth.ValidateScope(ctx, "write:user-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify another users sets", "requiredScope": "write:user-set"})
			return
		}
	}

	var newSet *setModel.Set

	if ctx.Bind(&newSet) != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to bind response object. Object structure may be incorrect"})
		return
	}

	if newSet.Name == "" || newSet.Code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Set code or name must not be empty when creating a new set"})
		return
	}

	if newSet.MtgjsonApiMeta != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "The mtgjsonApiMeta field must be null. This will be filled out automatically during deck creation"})
		return
	}

	if newSet.ContentIds != nil || len(newSet.ContentIds) != 0 {
		var valid, invalidCards, noExistCards = card.ValidateCards(newSet.ContentIds)
		if valid != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for set creation", "err": valid.Error()})
			return
		}

		if len(invalidCards) != 0 || len(noExistCards) != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to create set. Some cards do not exist or are invalid", "invalidCards": invalidCards, "noExistCards": noExistCards})
			return
		}
	}

	var err = set.NewSet(newSet, owner)
	if errors.Is(err, sdkErrors.ErrSetAlreadyExists) {
		ctx.JSON(http.StatusConflict, gin.H{"message": "Set already exists under this set code", "code": newSet.Code})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully created new set", "code": newSet.Code})
}
