package errors

import "errors"

var ErrNoCard = errors.New("card: failed to find card with specified uuid")
var ErrNoCards = errors.New("card: No card found on index operation")
var ErrInvalidUUID = errors.New("card: invalid v5 uuid")
