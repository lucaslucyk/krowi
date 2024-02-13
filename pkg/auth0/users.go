package auth0

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserIdentity struct {
	Connection        string      `json:"connection"`
	UserId            string      `json:"user_id"`
	Provider          string      `json:"provider"`
	IsSocial          bool        `json:"isSocial"`
	AccessToken       string      `json:"access_token"`
	AccessTokenSecret string      `json:"access_token_secret"`
	RefreshToken      string      `json:"refresh_token"`
	ProfileData       interface{} `json:"profileData"`
}

// type UserAppMetadata struct {
// 	// TODO: define itentity fields
// }

type User struct {
	Id            string         `json:"user_id"`
	Email         string         `json:"email"`
	EmailVerified bool           `json:"email_verified"`
	Username      string         `json:"username"`
	PhoneNumber   string         `json:"phone_number"`
	PhoneVerified bool           `json:"phone_verified"`
	CreatedAt     string         `json:"created_at"`
	UpdatedAt     string         `json:"updated_at"`
	Identities    []UserIdentity `json:"identities"`
	AppMetadata   interface{}    `json:"app_metadata"`
	UserMetadata  interface{}    `json:"user_metadata"`
	Picture       string         `json:"picture"`
	Name          string         `json:"name"`
	NickName      string         `json:"nickname"`
	Multifactor   []string       `json:"multifactor"`
	LastIp        string         `json:"last_ip"`
	LastLogin     string         `json:"last_login"`
	LoginsCount   int            `json:"logins_count"`
	Blocked       bool           `json:"blocked"`
	GivenName     string         `json:"given_name"`
	FamilyName    string         `json:"family_name"`
}

func (m *ManagerAPI) GetUsers(api *ManagerAPI) []User {
	// TODO: Get from API
	return []User{{}, {}}
}

func (m *ManagerAPI) GetUser(api *ManagerAPI, userId string, v *User) error {

	url := fmt.Sprintf("%s/%s", m.UsersUrl(), userId)

	resp, err := m.ClientSession.Get(url)
	if err != nil {
		return fmt.Errorf("Error updating user metadata error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("User metadata error: %v", err)
	}

	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return fmt.Errorf("Bad response on login: %v", err)
	}
	return nil
}

func (m *ManagerAPI) FindUser(api *ManagerAPI, email string) User {
	// TODO: Get from API
	return User{}
}

func (m *ManagerAPI) SetUserMetadata(userId string, data map[string]interface{}) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/%s", m.UsersUrl(), userId)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("Error building user metadata: %v", err)
	}

	// Crear la solicitud PATCH utilizando la instancia de *http.Client existente
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error al crear la solicitud PATCH:", err)
		return nil, fmt.Errorf("Error building request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.ClientSession.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error updating user metadata error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("User metadata error: %v", err)
	}

	// parse response

	var userMetadata map[string]interface{}

	if err := json.NewDecoder(resp.Body).Decode(&userMetadata); err != nil {
		return nil, fmt.Errorf("Bad response on login: %v", err)
	}

	// return updated user metadata
	return userMetadata, nil
}
