package helper

import (
	"context"
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
	"golang.org/x/oauth2"
)

var (
	helper Helper
)

type Helper interface {
	// BindJSON calls c.ShouldBindJSON and returns false on error.
	// For easy testing, call StartMock() before using BindJSON.
	BindJSON(c *gin.Context, jsonToBind interface{}) bool

	// Exchange exchanges an authorization code for an OAuth2 token using the given config.
	// For easy testing, call StartMock() before using Exchange.
	Exchange(ctx context.Context, code string, cfg *oauth2.Config) (*oauth2.Token, error)

	// GetRandomBytes returns the specified number of random bytes.
	// For easy testing, call StartMock() before using GetRandomBytes.
	GetRandomBytes(bytesNumber int) ([]byte, *resp.Err)

	// HashPassword generates a bcrypt hash for the given password.
	// For easy testing, call StartMock() before using HashPassword.
	HashPassword(password string) (string, *resp.Err)

	// ParseTime parses a time string using the provided format. Returns false on error.
	// For easy testing, call StartMock() before using ParseTime.
	ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool

	// ParseUUID parses a UUID string. Returns false on error.
	// For easy testing, call StartMock() before using ParseUUID.
	ParseUUID(c *gin.Context, jsonFieldName string, idTxtToParse string, dest *uuid.UUID) bool

	// TokenSignedString signs a JWT token using the provided secret.
	// For easy testing, call StartMock() before using TokenSignedString.
	TokenSignedString(token *jwt.Token, secret []byte) (string, error)
}

func NewHelper() Helper {
	if isMock {
		helper = &helperMock{}
	} else {
		helper = &h{}
	}
	return helper
}

type h struct {
}

func (h *h) BindJSON(c *gin.Context, jsonToBind interface{}) bool {
	if err := c.ShouldBindJSON(jsonToBind); err != nil {
		log.Printf("invalid JSON data: %s", err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{err.Error()}).JSON())
		return false
	}
	return true
}

func (h *h) Exchange(ctx context.Context, code string, cfg *oauth2.Config) (*oauth2.Token, error) {
	return cfg.Exchange(ctx, code)
}

func (h *h) GetRandomBytes(bytesNumber int) ([]byte, *resp.Err) {
	bytes := make([]byte, bytesNumber)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("failed to generate random bytes:", err.Error())
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{err.Error()})
	}
	return bytes, nil
}

func (h *h) HashPassword(password string) (string, *resp.Err) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password:", err.Error())
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{err.Error()})
	}
	return string(bytes), nil
}

func (h *h) ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool {
	res, err := time.Parse(timeFormat, txtToParse)
	if err != nil {
		log.Printf("invalid time string '%s': %s", txtToParse, err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err.Error()}).JSON())
		return false
	}
	*dest = res
	return true
}

func (h *h) ParseUUID(c *gin.Context, jsonFieldName string, idTxtToParse string, dest *uuid.UUID) bool {
	parsedID, err := uuid.Parse(idTxtToParse)
	if err != nil {
		log.Printf("failed to parse %s: %v", jsonFieldName, err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, fmt.Sprintf("invalid UUID: %s", jsonFieldName), []interface{}{err.Error()}).JSON())
		return false
	}
	*dest = parsedID
	return true
}

func (h *h) TokenSignedString(token *jwt.Token, secret []byte) (string, error) {
	return token.SignedString(secret)
}
