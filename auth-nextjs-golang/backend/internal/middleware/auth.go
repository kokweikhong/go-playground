package middleware

import (
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

func ValidateJWTToken(tokenString string) (bool, error) {
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
	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	// return false if token is invalid
	return false, err
}

func NewJWTClaims(user string) (string, error) {
	// // get user service
	// userService := service.NewUser()
	// // get user by username
	// user, err := userService.GetUserByUsername(username)
	// if err != nil {
	// 	return "", err
	// }

	// // compare password with bcrypt and error handling
	// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	// if err != nil {
	// 	return "", err
	// }
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = user
	claims["exp"] = time.Now().Add(10 * time.Second).Unix()

	// generate encoded token and send it as response
	// the signing string should be secret (a generated UUID works too)
	t, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}

	return t, nil
}

func NewRefreshToken() (string, error) {
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := refreshToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = 1
	rtClaims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	rt, err := refreshToken.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}
	return rt, nil
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
