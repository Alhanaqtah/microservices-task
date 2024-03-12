package jwt

/*
import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"


	"github.com/go-chi/jwtauth/v5"
	"github.com/golang-jwt/jwt/v5"
)

func NewToken(user models.User, duration time.Duration, secret string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["exp"] = time.Now().Add(duration).Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CheckClaim(ctx context.Context, claim, expectedClaim string) (bool, error) {
	const op = "CheckClaim"

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	c, ok := claims[claim]
	if !ok {
		return false, fmt.Errorf("%s: claim not found", op)
	}

	switch c.(type) {
	case float64:
		claim, ok := c.(float64)
		if !ok {
			return false, fmt.Errorf("%s: %w", op, errors.New("type not found"))
		}

		expClaim, err := strconv.ParseFloat(expectedClaim, 64)
		if err != nil {
			return false, fmt.Errorf("%s: %w", op, err)
		}

		if claim != expClaim {
			return false, nil
		}
	case string:
		claim, ok := c.(string)
		if !ok {
			return false, fmt.Errorf("%s: %w", op, errors.New("type not found"))
		}

		if claim != expectedClaim {
			return false, nil
		}
	}

	return true, nil
}
*/
