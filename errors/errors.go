package errors

import "errors"

/*
Card Errors - Holds all errors that could arise from fetching or inserting cards
*/
var ErrNoCard = errors.New("card: failed to find card with specified uuid")
var ErrNoCards = errors.New("card: No card found on index operation")
var ErrInvalidUUID = errors.New("card: invalid v5 uuid")

/*
Deck Errors - Holds all errors that could arise from fetching or inserting decks
*/
var ErrNoDeck = errors.New("deck: failed to find deck with the specified code")
