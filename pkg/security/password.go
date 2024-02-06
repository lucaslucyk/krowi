package security

import "golang.org/x/crypto/bcrypt"

func PasswordVerify(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword),
		[]byte(password),
	)
	return err != nil
}
