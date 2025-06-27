package jwt

import (
	"fmt"
	"log"
	"maps"

	"github.com/fachrunwira/basic-go-api-template/config"

	"github.com/golang-jwt/jwt/v5"
)

var jwt_cfg = config.LoadJWTConfig()

func GenerateToken(claims map[string]interface{}) string {
	defaultClaims := map[string]interface{}{
		"iat": jwt_cfg.Iat,
		"exp": jwt_cfg.Exp,
	}

	mergedClaims := mergeClaims(defaultClaims, claims)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(mergedClaims))
	signToken, err := token.SignedString(jwt_cfg.Key)
	if err != nil {
		log.Fatalf("failed to generate token: %v", "err")
	}

	return signToken
}

func mergeClaims(map1, map2 map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})

	maps.Copy(merged, map1)
	maps.Copy(merged, map2)

	return merged
}

func ValidateToken(token string) (*jwt.Token, error) {
	parseToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwt_cfg.Key, nil
	})

	return parseToken, err
}
