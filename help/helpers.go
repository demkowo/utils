package utils

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/demkowo/utils/config"
	"github.com/demkowo/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	H    HInterface   = &h{}
	Var  VarInterface = &hMock{}
	conf              = config.Values.Get()
)

func StartMock() {
	H = &hMock{
		Error: make(map[string]error),
	}
}

func StopMock() {
	H = &h{}
}

type HInterface interface {
	AddJWTToken(jwt.MapClaims) (string, *resp.Err)
	BindJSON(*gin.Context, interface{}) bool
	GetRandomBytes(int) ([]byte, *resp.Err)
	HashPassword(string) (string, *resp.Err)
	ParseTime(*gin.Context, string, string, *time.Time) bool
	ParseUUID(*gin.Context, string, string, *uuid.UUID) bool
}

type h struct {
}

func (h *h) AddJWTToken(claims jwt.MapClaims) (string, *resp.Err) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(conf.JWTSecret)
	if err != nil {
		log.Println("failed to sign JWT token:", err)
		return "", resp.Error(http.StatusInternalServerError, "failed to create token", []interface{}{err})
	}
	return tokenString, nil
}

func (h *h) BindJSON(c *gin.Context, jsonToBind interface{}) bool {
	if err := c.ShouldBindJSON(&jsonToBind); err != nil {
		log.Printf("invalid JSON data: %s", err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{err}).JSON())
		return false
	}
	return true
}

func (h *h) GetRandomBytes(bytesNumber int) ([]byte, *resp.Err) {
	bytes := make([]byte, bytesNumber)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("failed to generate random bytes:", err)
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{err})
	}
	return bytes, nil
}

func (h *h) HashPassword(password string) (string, *resp.Err) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("failed to hash password:", err)
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{err})
	}
	return string(bytes), nil
}

func (h *h) ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool {
	res, err := time.Parse(timeFormat, txtToParse)
	if err != nil {
		log.Printf("invalid input: %v", txtToParse)
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err}).JSON())
		return false
	}

	*dest = res
	return true
}

func (h *h) ParseUUID(c *gin.Context, jsonField string, txtToParse string, dest *uuid.UUID) bool {
	parsedID, err := uuid.Parse(txtToParse)
	if err != nil {
		log.Printf("failed to parse %s: %v", jsonField, err)
		c.JSON(resp.Error(http.StatusBadRequest, fmt.Sprintf("invalid UUID: %s", jsonField), []interface{}{err}).JSON())
		return false
	}

	*dest = parsedID
	return true
}
