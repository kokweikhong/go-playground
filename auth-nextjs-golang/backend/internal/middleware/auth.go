package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// cookie, err := r.Cookie("token")
		// if err != nil {
		// 	// authentication failed
		// 	http.Error(w, err.Error(), http.StatusUnauthorized)
		// 	return
		// }

		// tokenString := cookie.Value

		// get token from header authorization bearer
		authHeader := r.Header.Get("Authorization")

		// split bearer and token
		bearer := strings.Split(authHeader, "Bearer")
		if len(bearer) != 2 {
			// authentication failed
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimSpace(bearer[1])

		isTokenValid, err := ValidateJWTToken(tokenString)

		if err != nil || !isTokenValid {
			// authentication failed
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// authentication success
		next.ServeHTTP(w, r)

	})
}

const JWT_SECRET = "my_jwt_secret"

type JWTCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type JWTPayload struct {
	Username string `json:"username"`
	Issuer   string `json:"iss"`
	Expiry   int64  `json:"exp"`
}

func ValidateJWTToken(tokenString string) (bool, error) {
	// parse with JWTCustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %s", token.Header["alg"])
		}
		// return secret key
		return []byte(JWT_SECRET), nil
	})
	if err != nil {
		return false, err
	}

	// validate token
	if _, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return true, nil
	}

	// return false if token is invalid
	return false, nil
}

// return payload and token string
func NewJWTClaims(username string) (string, *JWTPayload, error) {
	// create new token with JWTCustomClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTCustomClaims{
		username,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			Issuer:    "myapp",
		},
	})

	// sign token with secret key
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", nil, err
	}

	// get payload
	payload, err := json.Marshal(token.Claims)
	if err != nil {
		return "", nil, err
	}

	jwtPayload := new(JWTPayload)
	if err := json.Unmarshal(payload, jwtPayload); err != nil {
		return "", nil, err
	}

	return tokenString, jwtPayload, nil
}

func NewRefreshToken(username string) (string, *JWTPayload, error) {
	// create new token with JWTCustomClaims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, JWTCustomClaims{
		username,
		jwt.RegisteredClaims{
			// 1 minute
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 10)),
			Issuer:    "myapp",
		},
	})

	// sign token with secret key
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", nil, err
	}

	// get payload
	payload, err := json.Marshal(token.Claims)
	if err != nil {
		return "", nil, err
	}

	jwtPayload := new(JWTPayload)
	if err := json.Unmarshal(payload, jwtPayload); err != nil {
		return "", nil, err
	}

	return tokenString, jwtPayload, nil
}

func ExtractJWTClaims(tokenString string) (jwt.MapClaims, error) {
	// parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %s", token.Header["alg"])
		}
		// return secret key
		return []byte(JWT_SECRET), nil
	})

	// validate token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	// return false if token is invalid
	return nil, err
}
