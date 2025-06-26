package api

import "net/http"

var BadRequest = http.StatusText(http.StatusBadRequest)
var Conflict = http.StatusText(http.StatusConflict)
var InternalServerError = http.StatusText(http.StatusInternalServerError)
var MethodNotAllowed = http.StatusText(http.StatusMethodNotAllowed)
var NotFound = http.StatusText(http.StatusNotFound)
var RequestEntityTooLarge = http.StatusText(http.StatusRequestEntityTooLarge)
var Unauthorized = http.StatusText(http.StatusUnauthorized)
var UnsupportedMediaType = http.StatusText(http.StatusUnsupportedMediaType)
var Accepted = http.StatusText(http.StatusAccepted)
