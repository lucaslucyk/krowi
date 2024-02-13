package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ManagerAPI struct {
	// to create
	Domain       string
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`

	// to consume api
	token Token

	// session
	ClientSession *http.Client
}

type Token struct {
	// to consume api
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	ExpiresIn   int64  `json:"expires_in"`
	TokenType   string `json:"token_type"`

	ExpireTime time.Time
}

func (t *Token) SetExpireTime() {
	// decrease 1 minute to expires in for security
	expireOffset := time.Duration(t.ExpiresIn - 60)
	t.ExpireTime = time.Now().Add(expireOffset * time.Second)
}

func New(domain, clientId, clientSecret, audience, grantType string) (*ManagerAPI, error) {

	api := ManagerAPI{
		Domain:       domain,
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Audience:     audience,
		GrantType:    grantType,
	}
	api.StartSession()
	if err := api.Login(); err != nil {
		return &api, err
	}

	// default
	return &api, nil
}

func (m *ManagerAPI) IsTokenExpired() bool {
	if m.token.AccessToken == "" || m.token.ExpiresIn == 0 {
		return true
	}
	return m.token.ExpireTime.After(time.Now())
}

func (m *ManagerAPI) StartSession() {
	m.ClientSession = &http.Client{
		Timeout: 10 * time.Second,
	}
}

func (m *ManagerAPI) Login() error {

	if !m.IsTokenExpired() {
		// nothing to do
		return nil
	}

	url := m.LoginUrl()
	payload := m.LoginPayload()
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Error building login data: %v", err)
	}

	ct := "application/json"
	resp, err := m.ClientSession.Post(url, ct, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Login error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Login error: %v", err)
	}

	// parse response

	var token Token

	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return fmt.Errorf("Bad response on login: %v", err)
	}
	token.SetExpireTime()
	m.token = token

	// set access token to client session
	m.SetAccessToken(token.AccessToken)

	// default return
	return nil
}

type BearerTokenTransport struct {
	AccessToken string
}

// RoundTrip implement http.RoundTripper interface
func (t *BearerTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	//add authorization token
	req.Header.Set("Authorization", "Bearer "+t.AccessToken)

	// send req to next roundtripper on chain
	return http.DefaultTransport.RoundTrip(req)
}

func (m *ManagerAPI) SetAccessToken(accessToken string) {
	if m.ClientSession.Transport == nil {
		m.ClientSession.Transport = &BearerTokenTransport{AccessToken: accessToken}
		return
	}

	if t, ok := m.ClientSession.Transport.(*BearerTokenTransport); ok {
		t.AccessToken = accessToken
	} else {
		// Si el Transport no es del tipo esperado, crear uno nuevo con el token de acceso
		m.ClientSession.Transport = &BearerTokenTransport{AccessToken: accessToken}
	}
}

// payloads
func (m *ManagerAPI) LoginPayload() map[string]interface{} {
	payload := make(map[string]interface{})
	payload["client_id"] = m.ClientId
	payload["client_secret"] = m.ClientSecret
	payload["audience"] = m.Audience
	payload["grant_type"] = m.GrantType

	return payload
}

// urls
func (m *ManagerAPI) LoginUrl() string {
	return fmt.Sprintf("https://%s/oauth/token", m.Domain)
}

func (m *ManagerAPI) UsersUrl() string {
	return fmt.Sprintf("https://%s/api/v2/users", m.Domain)
}
