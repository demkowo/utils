package serviceauth

import (
	"bytes"
	"fmt"
	"net/http"
)

func RegisterAtAuth(cfg Config) error {
	url := fmt.Sprintf("%s/internal/register", cfg.AuthServiceURL)

	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader(nil))
	req.Header.Set("X-Service-Name", cfg.ServiceName)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to register at auth: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth returned status %d", resp.StatusCode)
	}

	return nil
}
