package app_response

import "github.com/gin-gonic/gin"

// FieldError represents a validation error for a single field
type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Param string `json:"param,omitempty"`
}

// Response represents a standard API app_response body
type Response struct {
	Code    int          `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"`
}

// Success writes a standard success app_response
func Success(c *gin.Context, code int, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: "success",
		Data:    data,
	})
}

// Error writes a standard error app_response
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// ValidationError writes a standard validation error app_response with field details
func ValidationError(c *gin.Context, code int, message string, errors []FieldError) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Errors:  errors,
	})
}
