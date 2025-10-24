package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/bananazon/rdbak/pkg/data"
)

func (ac *APIClient) Login(email, pass string) error {
	payload := map[string]any{"email": email, "password": pass}
	payloadStr, _ := json.Marshal(payload)

	loginUrl := url.URL{
		Scheme: "https",
		Host:   apiBase,
		Path:   fmt.Sprintf("%s/auth/email/login", apiVersion),
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", loginUrl.String(), bytes.NewBuffer(payloadStr))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Accept", "application/json")

	resp, err := ac.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status at login: %d: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var loginResult data.LoginResult
	err = json.Unmarshal(body, &loginResult)
	if err != nil {
		return err
	}

	if !loginResult.Result {
		return fmt.Errorf("Login returned false: %s", loginResult.ErrorMessage)
	}

	return nil
}
