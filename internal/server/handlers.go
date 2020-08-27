package server

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/services"
	"github.com/rdsalakhov/game-keys-store/internal/store"
	"net/http"
	"os"
	"strconv"
	"time"
)

func (server *Server) handleLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		seller, err := server.store.Seller().FindByEmail(req.Email)
		if err != nil || !seller.ComparePassword(req.Password) {
			server.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		} else {
			tokenDetails, err := CreateToken(seller, os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET"))
			if err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
				return
			}
			if err := SaveAuth(server.redis, seller, tokenDetails); err != nil {
				server.error(w, r, http.StatusInternalServerError, err)
			}
			http.SetCookie(w, &http.Cookie{
				Name:       accessTokenCookie,
				Value:      tokenDetails.AccessToken,
				Path:       "/",
				RawExpires: time.Now().Add(accessTokenLifespan).String(),
			})
			http.SetCookie(w, &http.Cookie{
				Name:       refreshTokenCookie,
				Value:      tokenDetails.RefreshToken,
				Path:       "/",
				RawExpires: time.Now().Add(refreshTokenLifespan).String(),
			})

			server.respond(w, r, http.StatusOK, nil)
			return
		}
	}
}

func (server *Server) handleRegister() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		URL      string `json:"url"`
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}

		seller := &model.Seller{
			Email:    req.Email,
			Password: req.Password,
			URL:      req.URL,
			Account:  req.Account,
		}
		if err := seller.BeforeCreate(); err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
		}

		if err := server.store.Seller().Create(seller); err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		} else {
			seller.HidePassword()
			server.respond(w, r, http.StatusCreated, seller)
		}
	})
}

func (server *Server) handleRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(refreshTokenCookie)
		if err != nil {
			server.error(w, r, http.StatusUnauthorized, err)
			return
		}
		refreshToken := cookie.Value
		token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REFRESH_SECRET")), nil
		})
		if err != nil {
			server.error(w, r, http.StatusUnauthorized, err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		refreshUuid, ok := claims["refresh_uuid"].(string)
		if !ok {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		sellerID, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		//Delete the previous Refresh Token
		deleted, err := DeleteAuth(server.redis, refreshUuid)
		if err != nil || deleted == 0 { //if any goes wrong
			server.error(w, r, http.StatusUnauthorized, errNoAuthenticated)
			return
		}
		user, err := server.store.Seller().Find(int(sellerID))
		//Create new pairs of refresh and access tokens
		newTokenDetails, err := CreateToken(user, os.Getenv("ACCESS_SECRET"), os.Getenv("REFRESH_SECRET"))
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if err = SaveAuth(server.redis, user, newTokenDetails); err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:       accessTokenCookie,
			Value:      newTokenDetails.AccessToken,
			Path:       "/",
			RawExpires: time.Now().Add(accessTokenLifespan).String(),
		})
		http.SetCookie(w, &http.Cookie{
			Name:       refreshTokenCookie,
			Value:      newTokenDetails.RefreshToken,
			Path:       "/",
			RawExpires: time.Now().Add(refreshTokenLifespan).String(),
		})

		server.respond(w, r, http.StatusOK, nil)
		return
	}
}

func (server *Server) handlePostGame() http.HandlerFunc {
	type request struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Price       float32 `json:"price"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		sellerID := r.Context().Value(contextKeyID).(int)
		game := model.Game{
			Title:       req.Title,
			Description: req.Description,
			Price:       req.Price,
			OnSale:      true,
			SellerID:    sellerID,
		}

		if err := game.Validate(); err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		service := &services.GameService{
			Store: server.store,
		}
		if err := service.AddGame(&game); err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, game)
	})
}

func (server *Server) handleFindGameByID() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		service := &services.GameService{
			Store: server.store,
		}

		game, err := service.FindByID(id)
		if err == store.ErrRecordNotFound {

		} else if err != nil && err != store.ErrRecordNotFound {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, game)
	})
}

func (server *Server) handleFindAllGames() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service := services.GameService{Store: server.store}
		games, err := service.FindAll()
		if err != nil && err != store.ErrRecordNotFound {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		server.respond(w, r, http.StatusOK, games)
	})
}
