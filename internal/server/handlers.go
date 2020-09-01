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
			server.error(w, r, http.StatusNotFound, errItemNotFound)
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

func (server *Server) handleDeleteGameByID() http.HandlerFunc {
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

		err = service.DeleteByID(id)
		if err == store.ErrRecordNotFound {
			server.respond(w, r, http.StatusNotFound, nil)
		} else if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusNoContent, nil)
	})
}

func (server *Server) handlePostKeys() http.HandlerFunc {
	type request struct {
		GameID int      `json:"game_id"`
		Keys   []string `json:"keys"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := server.checkGameID(req.GameID); err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, errGameAccessDenied)
			return
		}
		sellerID := r.Context().Value(contextKeyID).(int)
		if !server.checkGameOwner(req.GameID, sellerID) {
			server.error(w, r, http.StatusForbidden, errGameAccessDenied)
			return
		}
		keys := []model.Key{}
		for _, keyString := range req.Keys {
			key := model.Key{KeyString: keyString}
			if err := key.Validate(); err != nil {
				server.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
			keys = append(keys, key)
		}
		service := services.KeyService{Store: server.store}
		if err := service.AddKeysToGame(req.GameID, &keys); err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, req)
	})
}

func (server *Server) handleBuyGame() http.HandlerFunc {
	type request struct {
		Name    string `json:"name"`
		Email   string `json:"email"`
		Address string `json:"address"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gameID, _ := strconv.Atoi(mux.Vars(r)["id"])
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		session := &model.PaymentSession{
			CustomerName:    req.Name,
			CustomerEmail:   req.Email,
			CustomerAddress: req.Address,
		}
		if err := session.Validate(); err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		if err := server.checkGameID(gameID); err != nil {
			server.error(w, r, http.StatusNotFound, errItemNotFound)
			return
		}
		keyService := services.KeyService{Store: server.store}
		key, err := keyService.FindAvailableKey(gameID)
		if err != nil {
			server.error(w, r, http.StatusNotFound, errNoKeys)
			return
		}

		paymentService := services.PaymentService{Store: server.store}
		sessionID, err := paymentService.CreateSession(key.ID, req.Name, req.Email, req.Address)
		if err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := keyService.MarkOnHold(key.ID); err != nil {
			paymentService.DeleteSession(sessionID)
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusCreated, sessionID)
	})
}

func (server *Server) handlePostPurchase() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, _ := strconv.Atoi(mux.Vars(r)["id"])
		if !server.checkAvailableSession(sessionID) {
			server.error(w, r, http.StatusBadRequest, errPerformedSession)
			return
		}
		cardInfo := &model.CardInfo{}
		if err := json.NewDecoder(r.Body).Decode(cardInfo); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
		if !checkCardNumber(cardInfo.Number) {
			server.error(w, r, http.StatusBadRequest, errInvalidCardNumber)
			return
		}

		service := services.PaymentService{Store: server.store}
		if err := service.PerformPurchase(sessionID, cardInfo); err != nil {
			server.error(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusNoContent, nil)
	})
}

func (server *Server) handleDeletePurchase() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, _ := strconv.Atoi(mux.Vars(r)["id"])
		if !server.checkAvailableSession(sessionID) {
			server.error(w, r, http.StatusBadRequest, errPerformedSession)
			return
		}
		service := services.PaymentService{Store: server.store}
		if err := service.DeleteSession(sessionID); err != nil {
			server.error(w, r, http.StatusBadRequest, err)
			return
		}
	})
}

func (server *Server) handleGetGameKeys() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gameID, _ := strconv.Atoi(mux.Vars(r)["id"])
		if err := server.checkGameID(gameID); err != nil {
			server.error(w, r, http.StatusUnprocessableEntity, errItemNotFound)
			return
		}
		sellerID := r.Context().Value(contextKeyID).(int)
		if !server.checkGameOwner(gameID, sellerID) {
			server.error(w, r, http.StatusForbidden, errGameAccessDenied)
			return
		}

		service := services.KeyService{Store: server.store}
		keys, err := service.FindByGameID(gameID)
		if err != nil {
			server.respond(w, r, http.StatusInternalServerError, err)
			return
		}
		server.respond(w, r, http.StatusOK, keys)
	})
}

func (server *Server) checkGameID(id int) error {
	service := services.GameService{Store: server.store}
	_, err := service.FindByID(id)
	return err
}

func (server *Server) checkGameOwner(gameID int, sellerID int) bool {
	service := services.GameService{Store: server.store}
	game, err := service.FindByID(gameID)
	if err != nil {
		return false
	}

	return game.SellerID == sellerID

}

func (server *Server) checkAvailableSession(sessionID int) bool {
	service := services.PaymentService{Store: server.store}
	session, err := service.FindByID(sessionID)
	if err != nil {
		return false
	}
	return !session.IsPerformed
}

func checkCardNumber(stringNumber string) bool {
	var luhn int64
	intNumber, err := strconv.ParseInt(stringNumber, 10, 64)
	if err != nil {
		return false
	}
	for i := 0; intNumber > 0; i++ {
		cur := intNumber % 10

		if i%2 == 0 { // even
			cur = cur * 2
			if cur > 9 {
				cur = cur%10 + cur/10
			}
		}

		luhn += cur
		intNumber = intNumber / 10
	}
	return luhn%10 == 0
}
