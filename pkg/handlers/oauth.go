package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
	"os"

	authenticator "github.com/lucaslucyk/krowi/pkg/authenticators"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func OAuthHandler(auth *authenticator.Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		state, err := generateRandomState()
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		//save the state inside the session
		session := c.Locals("session").(*session.Session)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		url := auth.AuthCodeURL(state, auth.Options...)
		return c.Redirect(url, http.StatusTemporaryRedirect)
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

func CallbackHandler(auth *authenticator.Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Locals("session").(*session.Session)

		if c.Query("state") != session.Get("state") {
			return c.SendStatus(http.StatusBadRequest)
		}
		// Exchange an authorization code for a token.j
		token, err := auth.Exchange(c.Context(), c.Query("code"))
		if err != nil {
			return c.SendStatus(http.StatusUnauthorized)
		}

		idToken, err := auth.VerifyIDToken(c.Context(), token)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		// TODO: remove unnecessary session keys
		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		// fmt.Printf("access_token: %s\n", token.AccessToken)

		if err := session.Save(); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		// return c.Status(fiber.StatusOK).JSON(fiber.Map{
		// 	"access_token": token.AccessToken,
		// 	"id_token":     idToken.AccessTokenHash,
		// })
		return c.Redirect("/oauth/me", http.StatusTemporaryRedirect)
	}
}

func OLogoutHandler(auth *authenticator.Authenticator) fiber.Handler {

	return func(c *fiber.Ctx) error {

		// clear session
		session := c.Locals("session").(*session.Session)
		if err := session.Destroy(); err != nil {
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}

		logoutUrl, err := url.Parse("https://" + os.Getenv("AUTH0_DOMAIN") + "/v2/logout")
		if err != nil {
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}

		scheme := "http"
		if c.Context().IsTLS() {
			scheme = "https"
		}
		returnTo, err := url.Parse(scheme + "://" + string(c.Context().Host()))
		if err != nil {
			return c.SendStatus(fiber.StatusServiceUnavailable)
		}

		parameters := url.Values{}
		parameters.Add("returnTo", returnTo.String())
		parameters.Add("client_id", os.Getenv("AUTH0_CLIENT_ID"))
		logoutUrl.RawQuery = parameters.Encode()

		return c.Redirect(logoutUrl.String(), http.StatusTemporaryRedirect)
	}
}
