package helper

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/demkowo/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	// Call provides access to helper methods.
	Call helperInterface = &helper{}
)

type helperInterface interface {
	// BindJSON calls c.ShouldBindJSON and returns false on error.
	//
	// For easy testing, call StartMock() before using BindJSON.
	//
	// Example:
	//	var acc model.Account
	//	if !helper.Run.BindJSON(c, &acc) {
	//	    return
	//	}
	BindJSON(c *gin.Context, jsonToBind interface{}) bool

	// GetRandomBytes returns the specified number of random bytes.
	//
	// For easy testing, call StartMock() before using GetRandomBytes.
	//
	// Example:
	//	bytes, err := helper.Run.GetRandomBytes(32)
	//	if err != nil {
	//	    return nil, err
	//	}
	GetRandomBytes(bytesNumber int) ([]byte, *resp.Err)

	// HashPassword generates a bcrypt hash for the given password using bcrypt.DefaultCost.
	//
	// For easy testing, call StartMock() before using HashPassword.
	//
	// Example:
	//	hashedPassword, err := helper.Run.HashPassword(acc.Password)
	//	if err != nil {
	//	    return nil, err
	//	}
	HashPassword(password string) (string, *resp.Err)

	// ParseTime parses the given time string based on the specified format.
	// Returns false on error.
	//
	// For easy testing, call StartMock() before using ParseTime.
	//
	// Example:
	//	var until time.Time
	//	if !helper.Run.ParseTime(c, "2006-01-02", req.Until, &until) {
	//	    return
	//	}
	ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool

	// ParseUUID parses the given string into a UUID and stores the result in dest.
	// Returns false on error.
	//
	// For easy testing, call StartMock() before using ParseUUID.
	//
	// Example:
	//	var accountId uuid.UUID
	//	if !helper.Run.ParseUUID(c, "account_id", c.Param("account_id"), &accountId) {
	//	    return
	//	}
	ParseUUID(c *gin.Context, jsonFieldName string, idTxtToParse string, dest *uuid.UUID) bool

	// TokenSignedString signs a JWT token using the given secret.
	//
	// For easy testing, call StartMock() before using TokenSignedString.
	//
	// Example:
	//	tokenString, err := helper.Run.TokenSignedString(token, conf.JWTSecret)
	//	if err != nil {
	//	    return "", err
	//	}
	TokenSignedString(token *jwt.Token, secret []byte) (string, error)
}

type helper struct {
}

func (h *helper) BindJSON(c *gin.Context, jsonToBind interface{}) bool {
	if err := c.ShouldBindJSON(&jsonToBind); err != nil {
		log.Printf("invalid JSON data: %s", err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{err.Error()}).JSON())
		return false
	}
	return true
}

func (h *helper) GetRandomBytes(bytesNumber int) ([]byte, *resp.Err) {
	bytes := make([]byte, bytesNumber)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("failed to generate random bytes:", err.Error())
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{err.Error()})
	}
	return bytes, nil
}

func (h *helper) HashPassword(password string) (string, *resp.Err) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password:", err.Error())
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{err.Error()})
	}
	return string(bytes), nil
}

func (h *helper) ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool {
	res, err := time.Parse(timeFormat, txtToParse)
	if err != nil {
		log.Printf("invalid input: %v", txtToParse)
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err.Error()}).JSON())
		return false
	}

	*dest = res
	return true
}

func (h *helper) ParseUUID(c *gin.Context, jsonFieldName string, idTxtToParse string, dest *uuid.UUID) bool {
	parsedID, err := uuid.Parse(idTxtToParse)
	if err != nil {
		log.Printf("failed to parse %s: %v", jsonFieldName, err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, fmt.Sprintf("invalid UUID: %s", jsonFieldName), []interface{}{err.Error()}).JSON())
		return false
	}

	*dest = parsedID
	return true
}

func (h *helper) TokenSignedString(token *jwt.Token, secret []byte) (string, error) {
	return token.SignedString(secret)
}
