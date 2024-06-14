package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(username string) string {
	var (
		key []byte
		t   *jwt.Token
		s   string
	)

	key = []byte(os.Getenv("JWT_SECRET"))
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 10)),
		Subject:   username,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	s, err = t.SignedString(key)
	if err != nil {
		log.Fatalf("Error signing jwt: %v", err)
	}
	return s
}

func parseToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		hmacSampleSecret := []byte(os.Getenv("JWT_SECRET"))
		return hmacSampleSecret, nil
	})
	return token, err
}

// Returns true, nil, username if the validation was successful. Returns false, err, "" otherwise
func validateToken(tokenString string) error {

	token, err := parseToken(tokenString)

	switch {
	case token.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenExpired):
		return err
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return err
	default:
		log.Fatalf("Error when validating token. %v", err)
		return err
	}
}

type UnsignedResponse struct {
	Message interface{} `json:"message"`
}

func JwtTokenCheck(c *gin.Context) {
	jwtToken, err := extractBearerToken(c.GetHeader("Authorization"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	err = validateToken(jwtToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, UnsignedResponse{
			Message: err.Error(),
		})
		return
	}

	c.Next()
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}

func GetUsernameFromCtx(ctx *gin.Context) (string, error) {
	username := ""
	jwtToken, _ := extractBearerToken(ctx.GetHeader("Authorization"))
	token, _ := parseToken(jwtToken)

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		username, err := claims.GetSubject()
		if err != nil {
			return username, errors.New("could not extract username from jwt claim")
		}
		return username, nil
	} else {
		return username, errors.New("could not extract username from jwt claim")
	}
}
