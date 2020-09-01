package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
	"net/http"
	"os"
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

	checksum, err := checksum(reqBody)
	if err != nil {
		errCh <- err
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		errCh <- err
		return
	}
	req.Header.Add("Checksum", checksum)

	err = nil
	for i := 0; i < 15; i++ {
		var resp *http.Response
		resp, err = client.Do(req)
		if resp == nil {
			continue
		}
		if resp.StatusCode/100 == 2 {
			break
		} else if err == nil {
			err = errors.New("fail to notify seller")
		}
		resp.Body.Close()
	}

	errCh <- err
}

func checksum(data []byte) (string, error) {
	salt := os.Getenv("NOTIFICATION_SALT")
	saltBytes, err := json.Marshal(salt)
	if err != nil {
		return "", err
	}
	saltedData := append(data, saltBytes...)

	sum := md5.Sum(saltedData)
	sumString := hex.EncodeToString(sum[:])
	return sumString, nil
}
