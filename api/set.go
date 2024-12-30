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
SetGET Gin handler for the GET request to the Set Endpoint. This function should not be called
directly and should only be passed to the gin router
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
		return
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

/*
SetDELETE Gin handler for the DELETE request to the Set Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func SetDELETE(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner == "system" {
		if !auth.ValidateScope(ctx, "write:system-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to delete system or pre-constructed set", "requiredScope": "write:system-set"})
			return
		}
	}

	if owner != userEmail {
		if !auth.ValidateScope(ctx, "write:user-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to delete other users sets", "requiredScope": "write:user-set"})
			return
		}
	}

	code := ctx.Query("setCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Set code is required to perform a DELETE operation on a set"})
		return
	}

	_set, err := set.GetSet(code, owner)
	if errors.Is(err, sdkErrors.ErrNoSet) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No set found for DELETE operation"})
		return
	}

	result := set.DeleteSet(_set.Code, owner)
	if errors.Is(result, sdkErrors.ErrSetDeleteFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": result.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully deleted set", "setCode": _set.Code})
}

/*
SetContentGET Gin handler for the GET request to the Set Content Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func SetContentGET(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)
	if owner != "system" && owner != userEmail {
		if !auth.ValidateScope(ctx, "read:user-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to read other users sets", "requiredScope": "read:user-set"})
			return
		}
	}

	code := ctx.Query("setCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Set code is required to fetch a sets contents"})
		return
	}

	_set, err := set.GetSet(code, owner)
	if errors.Is(err, sdkErrors.ErrNoSet) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No set found under the passed set code", "setCode": code})
		return
	}

	err = set.GetSetContents(_set)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, _set.Contents)
}

/*
SetContentPOST Gin handler for the POST request to the Set Content Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func SetContentPOST(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner == "system" {
		if !auth.ValidateScope(ctx, "write:system-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify content of system or pre-constructed set", "requiredScope": "write:system-set"})
			return
		}
	}

	if owner != userEmail {
		if !auth.ValidateScope(ctx, "write:user-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify content of user owned set", "requiredScope": "write:user-set"})
			return
		}
	}

	code := ctx.Query("setCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Set code is required to perform a POST operation on the sets contents"})
		return
	}

	_set, err := set.GetSet(code, owner)
	if errors.Is(err, sdkErrors.ErrNoSet) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No set found under the passed set code", "setCode": code})
		return
	}

	var updates []string
	if ctx.Bind(&updates) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response to object. Object structure may be incorrect"})
		return
	}

	if updates == nil || len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "At least one card must be passed in the body of this request. List cannot be empty"})
		return
	}

	set.AddCards(_set, updates)

	err = set.ReplaceSet(_set)
	if errors.Is(err, sdkErrors.ErrSetUpdateFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully updated set", "setCode": code}) // re-add count here
}

/*
SetContentDELETE Gin handler for the DELETE request to the Set Content Endpoint. This function should not be called
directly and should only be passed to the gin router
*/
func SetContentDELETE(ctx *gin.Context) {
	userEmail := ctx.GetString("userEmail")
	owner := ctx.DefaultQuery("owner", userEmail)

	if owner == "system" {
		if !auth.ValidateScope(ctx, "write:system-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify content of system or pre-constructed sets", "requiredScope": "write:system-set"})
			return
		}
	}

	if owner != userEmail {
		if !auth.ValidateScope(ctx, "write:user-set") {
			ctx.JSON(http.StatusForbidden, gin.H{"message": "Invalid permissions to modify content of user owned sets", "requiredScope": "write:user-set"})
			return
		}
	}

	code := ctx.Query("setCode")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Deck code is required to perform this operation"})
		return
	}

	_set, err := set.GetSet(code, owner)
	if errors.Is(err, sdkErrors.ErrNoSet) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "No set found under the passed set code"})
		return
	}

	var updates []string
	if ctx.Bind(&updates) != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to bind response object. Object structure may be incorrect"})
		return
	}

	if updates == nil || len(updates) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "At least one card must be passed in the body of this request. List cannot be empty"})
		return
	}

	valid, invalidCards, noExistCards := card.ValidateCards(updates)
	if valid != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Error while validating cards for set update", "err": valid.Error()})
		return
	}

	if len(invalidCards) != 0 || len(noExistCards) != 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Failed to update set. Some cards do not exist or are not valid", "invalidCards": invalidCards, "noExistCards": noExistCards})
		return
	}

	set.RemoveCards(_set, updates)

	err = set.ReplaceSet(_set)
	if errors.Is(err, sdkErrors.ErrSetUpdateFailed) {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully removed cards from set", "setCode": code})
}
