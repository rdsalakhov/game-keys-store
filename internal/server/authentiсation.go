package server

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	accessTokenCookie    = "AccessToken"
	refreshTokenCookie   = "RefreshToken"
	contextKeyID         = "SellerId"
	accessTokenLifespan  = time.Minute * 15
	refreshTokenLifespan = time.Hour * 24 * 7
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

func (server *Server) authenticateSeller(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenCookie, err := r.Cookie(accessTokenCookie)
		if err != nil {
			server.error(w, r, http.StatusUnauthorized, errNoAuthenticated)
			return
		}
		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(tokenCookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("ACCESS_SECRET")), nil
		})
		if err != nil {
			server.error(w, r, http.StatusUnauthorized, errNoAuthenticated)
			return
		}
		sellerID := int(claims["user_id"].(float64))
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKeyID, sellerID)))
	})
}

func CreateToken(user *model.Seller, accessSecret string, refreshSecret string) (*TokenDetails, error) {
	tokenDetails := &TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(accessTokenLifespan).Unix()
	tokenDetails.AccessUuid = uuid.New().String()
	tokenDetails.RtExpires = time.Now().Add(refreshTokenLifespan).Unix()
	tokenDetails.RefreshUuid = uuid.New().String()

	atClaims := jwt.MapClaims{}
	atClaims["access_uuid"] = tokenDetails.AccessUuid
	atClaims["user_id"] = user.ID
	atClaims["exp"] = tokenDetails.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	tokenDetails.AccessToken, err = at.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, err
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = tokenDetails.RefreshUuid
	rtClaims["user_id"] = user.ID
	rtClaims["exp"] = tokenDetails.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	tokenDetails.RefreshToken, err = rt.SignedString([]byte(refreshSecret))
	if err != nil {
		return nil, err
	}

	return tokenDetails, nil
}

func SaveAuth(client *redis.Client, seller *model.Seller, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := client.Set(context.Background(), td.AccessUuid, strconv.Itoa(seller.ID), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := client.Set(context.Background(), td.RefreshUuid, strconv.Itoa(seller.ID), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func DeleteAuth(client *redis.Client, uuid string) (int64, error) {
	deleted, err := client.Del(context.Background(), uuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
