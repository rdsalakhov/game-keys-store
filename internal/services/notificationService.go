package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
	"net/http"
	"time"
)

type NotificationService struct {
	Store interfaces.IStore
}

func (service *NotificationService) NotifySeller(url string, notification *sellerNotification, sessionID int) {
	errCh := make(chan error)
	go service.sendSellerNotification(url, notification, errCh)
	if err := <-errCh; err == nil {
		service.Store.PaymentSession().UpdateNotifiedStatus(sessionID, true)
	}
}

func (service *NotificationService) sendSellerNotification(url string, notification *sellerNotification, errCh chan error) {
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	reqBody, err := json.Marshal(notification)
	if err != nil {
		errCh <- err
		return
	}

	err = nil
	for i := 0; i < 15; i++ {
		var resp *http.Response
		resp, err = client.Post(url, "application/json", bytes.NewBuffer(reqBody))
		if resp.StatusCode/100 == 2 {
			break
		} else if err == nil {
			err = errors.New("fail to notify seller")
		}
	}

	errCh <- err
}
