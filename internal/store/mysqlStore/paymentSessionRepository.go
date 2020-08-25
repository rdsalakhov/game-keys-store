package mysqlStore

import (
	"database/sql"
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store"
)

type PaymentSessionRepository struct {
	store *Store
}

func (repo *PaymentSessionRepository) Find(ID int) (*model.PaymentSession, error) {
	selectQuery := "SELECT id, key_id, price, date, customer_name, customer_email, customer_address FROM payment_sessions WHERE id = ?"
	session := &model.PaymentSession{}
	if err := repo.store.db.QueryRow(selectQuery, ID).Scan(
		&session.ID,
		&session.KeyID,
		&session.Price,
		&session.Date,
		&session.CustomerName,
		&session.CustomerEmail,
		&session.CustomerAddress); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return session, nil
}

func (repo *PaymentSessionRepository) Create(session *model.PaymentSession) error {
	insertQuery := "INSERT INTO payment_sessions (key_id, price, date, customer_name, customer_email, customer_address) VALUES (?, ?, ?, ?, ?, ?);"
	getIdQuery := "select LAST_INSERT_ID();"

	if _, err := repo.store.db.Exec(insertQuery,
		session.KeyID,
		session.Price,
		session.Date,
		session.CustomerName,
		session.CustomerEmail,
		session.CustomerAddress); err != nil {
		return err
	}

	if err := repo.store.db.QueryRow(getIdQuery).Scan(&session.ID); err != nil {
		return err
	}
	return nil
}
