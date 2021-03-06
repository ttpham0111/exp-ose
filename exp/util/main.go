package util

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v8"

	"github.com/ttpham0111/exp-ose/exp/database"
)

func Getenv(name string, fallback string) string {
	if value, ok := os.LookupEnv(name); ok {
		return value
	}
	return fallback
}

func EnsureEnv(name string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	log.Fatal("Missing environmental variable: " + name)
	return ""
}

func FirstError(errors validator.ValidationErrors) *validator.FieldError {
	for _, err := range errors {
		return err
	}
	return nil
}

func HandleBindError(c *gin.Context, err error) {
	var errMessage string

	switch er := err.(type) {
	case validator.ValidationErrors:
		fe := FirstError(er)
		errMessage = "invalid value for " + fe.FieldNamespace
	case database.ValidationError:
		errMessage = er.Error()
	case *json.SyntaxError:
		errMessage = "invalid JSON"
	case *time.ParseError:
		errMessage = fmt.Sprintf(
			"invalid time '%s', expecting format '%s'",
			strings.Trim(er.Value, "\""),
			strings.Trim(er.Layout, "\""),
		)
	default:
		panic(err)
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": errMessage})
}
