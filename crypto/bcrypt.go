package crypto

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// ValidatePassword Validate the password during login
func ValidatePassword(hashedPassword, password string) error {
	// Compare the hashed password with the entered password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err
}
