package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	// define common responses and structs for REST
	// interfaces and responses
	InternalServerErrorResponse = gin.H{
		"http_code": http.StatusInternalServerError,
		"message":   "Internal server error",
	}
	InvalidRequestBodyResponse = gin.H{
		"http_code": http.StatusBadRequest,
		"message":   "Invalid request body",
	}
	ForbiddenResponse = gin.H{
		"http_code": http.StatusForbidden,
		"message":   "Forbidden",
	}
	UnauthorizedResponse = gin.H{
		"http_code": http.StatusUnauthorized,
		"message":   "Unauthorized",
	}
	NotImplementedResponse = gin.H{
		"http_code": http.StatusNotImplemented,
		"message":   "Not Implemented",
	}
)
