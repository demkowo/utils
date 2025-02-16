package httpclient

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

var (
	isMock bool
)

type clientMock struct {
	mocks map[string]Mock
}

type Mock struct {
	Error    map[string]error
	Response http.Response
	Test     string
}

func AddMock(mock Mock) {
	if !isMock {
		log.Println("mock server is off, ignoring AddMock")
		return
	}
	mc, ok := client.(*clientMock)
	if !ok {
		log.Fatal("the global client is not a mock; cannot add mock data")
	}
	mc.mocks[mock.Test] = mock
}

func StartMock() {
	isMock = true
}

func StopMock() {
	isMock = false
}

func (cm *clientMock) Get(url string, headers map[string]string) (*http.Response, error) {
	return cm.handleMock(http.MethodGet)
}

func (cm *clientMock) Post(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return cm.handleMock(http.MethodPost)
}

func (cm *clientMock) Put(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return cm.handleMock(http.MethodPut)
}

func (cm *clientMock) Patch(url string, body []byte, headers map[string]string) (*http.Response, error) {
	return cm.handleMock(http.MethodPatch)
}

func (cm *clientMock) Delete(url string, headers map[string]string) (*http.Response, error) {
	return cm.handleMock(http.MethodDelete)
}

func (cm *clientMock) Head(url string, headers map[string]string) (*http.Response, error) {
	return cm.handleMock(http.MethodHead)
}

func (cm *clientMock) Options(url string, headers map[string]string) (*http.Response, error) {
	return cm.handleMock(http.MethodOptions)
}

func (cm *clientMock) handleMock(method string) (*http.Response, error) {
	mock, err := getMock(cm)
	if err != nil {
		return nil, err
	}
	if mockErr := mock.Error[method]; mockErr != nil {
		return nil, mockErr
	}
	return &mock.Response, nil
}

func getMock(cm *clientMock) (*Mock, error) {
	pc, _, _, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	testName := parts[len(parts)-1]

	m, exists := cm.mocks[testName]
	if !exists {
		return nil, fmt.Errorf("mock not found for %s test", testName)
	}
	return &m, nil
}
