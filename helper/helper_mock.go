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
)

var (
	// Get provides access to the helper mock variables.
	Get = Var.get()
)

// StartMock initializes the helper mock.
func StartMock() {
	Call = &helperMock{
		Error:  make(map[string]error),
		IsMock: make(map[string]bool),
	}
}

// StopMock stops the helper mock.
func StopMock() {
	Call = &helperMock{}
}

type helperMock struct {
	Error    map[string]error
	IsMock   map[string]bool
	Password string
}

func (hm *helperMock) BindJSON(c *gin.Context, jsonToBind interface{}) bool {
	if err := c.ShouldBindJSON(&jsonToBind); err != nil {
		log.Printf("invalid JSON data: %s", err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{err}).JSON())
		return false
	}
	return true
}

func (hm *helperMock) GetRandomBytes(bytesNumber int) ([]byte, *resp.Err) {
	if Get.Error["GetRandomBytes"] != nil {
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{Get.Error["GetRandomBytes"].Error()})
	}

	bytes := make([]byte, bytesNumber)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("failed to generate random bytes:", err)
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{err})
	}
	return bytes, nil
}

func (hm *helperMock) HashPassword(password string) (string, *resp.Err) {
	if err := Get.Error["HashPassword"]; err != nil {
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{err.Error()})
	}

	return "hashedPassword", nil
}

func (hm *helperMock) ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool {
	if err := Get.Error["ParseTime"]; err != nil {
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err.Error()}).JSON())
		return false
	}

	res, err := time.Parse(timeFormat, txtToParse)
	if err != nil {
		log.Printf("invalid input: %v", txtToParse)
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err}).JSON())
		return false
	}

	*dest = res
	return true
}

func (hm *helperMock) ParseUUID(c *gin.Context, jsonFieldName string, idTxtToParse string, dest *uuid.UUID) bool {
	if err := Get.Error["ParseUUID"]; err != nil {
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err.Error()}).JSON())
		return false
	}

	parsedID, err := uuid.Parse(c.Param(idTxtToParse))
	if err != nil {
		log.Printf("failed to parse %s: %v", jsonFieldName, err)
		c.JSON(resp.Error(http.StatusBadRequest, fmt.Sprintf("invalid UUID: %s", jsonFieldName), []interface{}{err}).JSON())
		return false
	}

	*dest = parsedID
	return true
}

func (hm *helperMock) TokenSignedString(token *jwt.Token, secret []byte) (string, error) {
	if err := Get.Error["TokenSignedString"]; err != nil {
		return "", err
	}

	if Get.Error["TokenSignedString"] != nil {
		return "", Get.Error["TokenSignedString"]
	}
	return token.SignedString(secret)
}
