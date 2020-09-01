package services

import (
	"fmt"
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store/interfaces"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

type sellerNotification struct {
	SessionID       int     `json:"session_id"`
	GameTitle       string  `json:";game_title"`
	CustomerName    string  `json:"customer_name"`
	CustomerEmail   string  `json:"customer_email"`
	CustomerAddress string  `json:"customer_address"`
	SellerAmount    float64 `json:"seller_amount"`
	PlatformAmount  float64 `json:"platform_fee"`
	Signature       string  `json:"signature"`
}

type PaymentService struct {
	Store interfaces.IStore
}

func (service *PaymentService) CreateSession(keyID int, name string, email string, address string) (int, error) {
	session := &model.PaymentSession{
		KeyID:           keyID,
		Date:            time.Now(),
		CustomerName:    name,
		CustomerEmail:   email,
		CustomerAddress: address,
	}
	if err := service.Store.PaymentSession().Create(session); err != nil {
		return 0, err
	}
	return session.ID, nil
}

func (service *PaymentService) DeleteSession(id int) error {
	paymentInfo, err := service.Store.PaymentSession().GetPaymentInfo(id)
	if err != nil {
		return err
	}
	if err := service.Store.PaymentSession().DeleteByID(id); err != nil {
		return err
	}
	if err := service.Store.Key().UpdateStatus(paymentInfo.KeyID, model.KeyStatusAvailable); err != nil {
		return err
	}

	return nil
}

func (service *PaymentService) FindByID(id int) (*model.PaymentSession, error) {
	session, err := service.Store.PaymentSession().Find(id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (service *PaymentService) PerformPurchase(sessionID int, cardInfo *model.CardInfo) error {
	platformFeeShare, err := strconv.ParseFloat(os.Getenv("PLATFORM_FEE_SHARE"), 64)
	if err != nil {
		return err
	}

	paymentInfo, err := service.Store.PaymentSession().GetPaymentInfo(sessionID)
	if err != nil {
		return err
	}

	platformAmount := platformFeeShare * paymentInfo.TotalAmount
	sellerAmount := paymentInfo.TotalAmount - platformAmount

	if err := service.paySellerShare(sellerAmount, paymentInfo.SellerAccount, cardInfo); err != nil {
		return err
	}
	if err := service.payPlatformFee(paymentInfo.PlatformAmount, cardInfo); err != nil {
		return err
	}
	if err := service.sendKey(paymentInfo.CustomerEmail, paymentInfo.GameTitle, paymentInfo.CustomerName, paymentInfo.Key); err != nil {
		return err
	}

	if err := service.Store.Key().UpdateStatus(paymentInfo.KeyID, model.KeyStatusSold); err != nil {
		return err
	}

	sellerNotification := &sellerNotification{
		SellerAmount:    sellerAmount,
		PlatformAmount:  platformAmount,
		SessionID:       sessionID,
		CustomerName:    paymentInfo.CustomerName,
		CustomerEmail:   paymentInfo.CustomerEmail,
		CustomerAddress: paymentInfo.CustomerAddress,
	}
	noteService := &NotificationService{Store: service.Store}
	go noteService.NotifySeller(paymentInfo.SellerURL, sellerNotification, sessionID)

	if err := service.Store.PaymentSession().UpdatePerformedStatus(sessionID, true); err != nil {
		logrus.Printf("failed to mark session %d as performed", sessionID)
	}
	return nil
}

func (service *PaymentService) sendKey(email string, game string, name string, key string) error {
	from := os.Getenv("PLATFORM_EMAIL")
	auth := smtp.PlainAuth("", from, os.Getenv("PLATFORM_EMAIL_PASSWORD"), "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{email}
	msg := []byte(fmt.Sprintf(
		"To: %s \r\n"+
			"Subject: Your key\r\n"+
			"\r\n"+
			"Hello %s!\r\n"+
			"Thank you for your purchase!\r\n"+
			"Here's your key for %s: %s\r\n"+
			"Have fun!", email, name, game, key))
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}

func (service *PaymentService) payPlatformFee(amount float64, cardInfo *model.CardInfo) error {
	os.Getenv("PLATFORM_ACCOUNT")
	// perform payment
	return nil
}

func (service *PaymentService) paySellerShare(amount float64, account string, cardInfo *model.CardInfo) error {
	// perform payment
	return nil
}
