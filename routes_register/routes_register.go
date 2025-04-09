package routes_register

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Route struct {
	ID      uuid.UUID `json:"id"`
	Method  string    `json:"method"`
	Path    string    `json:"path"`
	Service string    `json:"service"`
	Active  bool      `json:"active"`
}

func RegisterRoutesWithRBAC(router *gin.Engine, service string) {
	routes := router.Routes()
	var payload []Route

	for _, route := range routes {
		payload = append(payload, Route{
			ID:      uuid.New(),
			Method:  route.Method,
			Path:    route.Path,
			Service: service,
			Active:  true,
		})
	}

	go func() {
		time.Sleep(2 * time.Second)

		rbacURL := os.Getenv("RBAC_REGISTER_URL")
		if rbacURL == "" {
			rbacURL = "http://rbac:5001/api/v1/routes"
		}

		body, err := json.Marshal(payload)
		if err != nil {
			log.Printf("marshal error during route register: %v", err)
			return
		}

		resp, err := http.Post(rbacURL, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("failed to register routes to RBAC: %v", err)
			return
		}
		defer resp.Body.Close()

		log.Printf("[%s] registered %d routes with RBAC", service, len(payload))
	}()
}
