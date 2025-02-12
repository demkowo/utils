// helper_mock.go
package helper

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/demkowo/utils/resp"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	isMock bool
)

type helperMock struct {
	mocks map[string]Mock
}

type Mock struct {
	Test     string
	Error    map[string]error
	IsMock   map[string]bool
	Password string
}

// AddMock allows to add variables to the mock when it's on
//
// Test variable is mandatory and the value should be equal with the name
// of function from which it's called
func AddMock(mock Mock) {
	if !isMock {
		log.Println("mock server is off, therefore ignoring AddMock")
		return
	}

	h, ok := helper.(*helperMock)
	if !ok {
		log.Fatal("invalid type of helperMock")
	}

	h.mocks = make(map[string]Mock)

	h.mocks[mock.Test] = mock
}

// StartMock replaces the real helper with a new mock instance.
func StartMock() {
	isMock = true
}

// StopMock reverts the Call pointer back to the real helper implementation
func StopMock() {
	isMock = false
}

func (hm *helperMock) BindJSON(c *gin.Context, jsonToBind interface{}) bool {
	mock, err := getMock(hm)
	if err != nil {
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{err.Error()}).JSON())
		return false
	}

	if mock.Error["BindJSON"] != nil {
		log.Printf("invalid JSON data: %s", mock.Error["BindJSON"].Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{mock.Error["BindJSON"].Error()}).JSON())
		return false
	}

	if err := c.ShouldBindJSON(jsonToBind); err != nil {
		log.Printf("invalid JSON data: %s", err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid JSON", []interface{}{err.Error()}).JSON())
		return false
	}
	return true
}

func (hm *helperMock) GetRandomBytes(bytesNumber int) ([]byte, *resp.Err) {
	mock, err := getMock(hm)
	if err != nil {
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{err.Error()})
	}

	if mock.Error["GetRandomBytes"] != nil {
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{mock.Error["GetRandomBytes"].Error()})
	}

	bytes := make([]byte, bytesNumber)
	if _, err := rand.Read(bytes); err != nil {
		log.Println("failed to generate random bytes:", err)
		return nil, resp.Error(http.StatusInternalServerError, "failed to create API Key", []interface{}{err.Error()})
	}

	return bytes, nil
}

func (hm *helperMock) HashPassword(password string) (string, *resp.Err) {
	mock, err := getMock(hm)
	if err != nil {
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{err.Error()})
	}

	if mock.Error["HashPassword"] != nil {
		return "", resp.Error(http.StatusInternalServerError, "failed to hash password", []interface{}{mock.Error["HashPassword"].Error()})
	}

	return mock.Password, nil

}

func (hm *helperMock) ParseTime(c *gin.Context, timeFormat string, txtToParse string, dest *time.Time) bool {
	mock, err := getMock(hm)
	if err != nil {
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err.Error()}).JSON())
		return false
	}

	if mock.Error["ParseTime"] != nil {
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{mock.Error["ParseTime"].Error()}).JSON())
		return false
	}

	res, err := time.Parse(timeFormat, txtToParse)
	if err != nil {
		log.Printf("invalid time string '%s': %s", txtToParse, err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, "invalid time format", []interface{}{err.Error()}).JSON())
		return false
	}

	*dest = res
	return true
}

func (hm *helperMock) ParseUUID(c *gin.Context, jsonFieldName string, idTxtToParse string, dest *uuid.UUID) bool {
	mock, err := getMock(hm)
	if err != nil {
		c.JSON(resp.Error(http.StatusBadRequest, "invalid UUID format", []interface{}{err.Error()}).JSON())
		return false
	}

	if mock.Error["ParseUUID"] != nil {
		c.JSON(resp.Error(http.StatusBadRequest, "invalid UUID format", []interface{}{mock.Error["ParseUUID"].Error()}).JSON())
		return false
	}

	parsedID, err := uuid.Parse(idTxtToParse)
	if err != nil {
		log.Printf("failed to parse %s: %v", jsonFieldName, err.Error())
		c.JSON(resp.Error(http.StatusBadRequest, fmt.Sprintf("invalid UUID: %s", jsonFieldName), []interface{}{err.Error()}).JSON())
		return false
	}
	*dest = parsedID
	return true
}

func (hm *helperMock) TokenSignedString(token *jwt.Token, secret []byte) (string, error) {
	mock, err := getMock(hm)
	if err != nil {
		return "", err
	}

	if mock.Error["TokenSignedString"] != nil {
		return "", mock.Error["TokenSignedString"]
	}
	return token.SignedString(secret)
}

func getMock(hm *helperMock) (*Mock, error) {
	pc, _, _, _ := runtime.Caller(3)
	test := strings.Split(runtime.FuncForPC(pc).Name(), ".")[len(strings.Split(runtime.FuncForPC(pc).Name(), "."))-1]
	mock, exists := hm.mocks[test]
	if !exists {
		return nil, fmt.Errorf("mock not found for %s test", test)
	}
	return &mock, nil
}
