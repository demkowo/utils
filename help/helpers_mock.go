package utils

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

type hMock struct {
	Password string
	Error    map[string]error
}

type VarInterface interface {
	Get() *hMock
	SetExpectedPassword(string)
	SetExpectedError(map[string]error)
}

func (v *hMock) SetExpectedPassword(password string) {
	v.Password = password
}

func (v *hMock) SetExpectedError(err map[string]error) {
	v.Error = err
}

func (v *hMock) Get() *hMock {
	return v
}

func (h *hMock) AddJWTToken(claims jwt.MapClaims) (string, *resp.Err) {
	err := Var.Get().Error["AddJWTToken"]

	if err != nil {
		return "", resp.Error(http.StatusInternalServerError, "failed to create token", []interface{}{err.Error()})
	}

	return "validTokenString", nil
}

func (h *hMock) BindJSON(c *gin.Context, jsonToBind interface{}) bool {
	if err := c.ShouldBindJSON(&jsonToBind); err != nil {
		log.Printf("invalid JSON data: %s", err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{err}).JSON())
		return false
	}
	return true
}

func (h *hMock) GetRandomBytes(bytesNumber int) ([]byte, *resp.Err) {
	if Var.Get().Error["GetRandomBytes"] != nil {
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{Var.Get().Error["GetRandomBytes"].Error()})
	}

	bytes := make([]byte, bytesNumber)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("failed to generate random bytes:", err)
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{err})
	}
	return bytes, nil
}

func (h *hMock) HashPassword(password string) (string, *resp.Err) {
	if err := Var.Get().Error["HashPassword"]; err != nil {
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{err})
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{err})
	}

	if Var.Get().Password != "" {
		if err := bcrypt.CompareHashAndPassword(bytes, []byte(Var.Get().Password)); err != nil {
			return "", resp.Error(http.StatusInternalServerError, "password mismatch", []interface{}{err.Error()})
		}
	}

	return password, nil
}

func (h *hMock) ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool {
	res, err := time.Parse(timeFormat, txtToParse)
	if err != nil {
		log.Printf("invalid input: %v", txtToParse)
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err}).JSON())
		return false
	}

	*dest = res
	return true
}

func (h *hMock) ParseUUID(c *gin.Context, jsonField string, txtToParse string, dest *uuid.UUID) bool {
	parsedID, err := uuid.Parse(txtToParse)
	if err != nil {
		log.Printf("failed to parse %s: %v", jsonField, err)
		c.JSON(resp.Error(http.StatusBadRequest, fmt.Sprintf("invalid UUID: %s", jsonField), []interface{}{err}).JSON())
		return false
	}

	*dest = parsedID
	return true
}
