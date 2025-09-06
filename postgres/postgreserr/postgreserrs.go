package postgreserr

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrUserNotFound = errors.New("user not found")
var ErrBalanceNotFound = errors.New("balance not found")
var ErrTradeNotFound = errors.New("trade not found")
