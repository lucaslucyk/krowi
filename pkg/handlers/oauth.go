package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	authenticator "github.com/lucaslucyk/krowi/pkg/authenticators"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func OAuthHandler(auth *authenticator.Authenticator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		state, err := generateRandomState()
		if err != nil {
			c.SendStatus(http.StatusInternalServerError)
		}

		//save the state inside the session
		session := c.Locals("session").(*session.Session)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		c.Redirect(auth.AuthCodeURL(state), http.StatusTemporaryRedirect)
		return nil
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
		session.Set("profile", profile)

		if err := session.Save(); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		return c.Redirect("/oauth/me", http.StatusTemporaryRedirect)
	}
}
