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
	selectQuery := "SELECT id, key_id, price, customer_name, customer_email, customer_address, is_performed, is_notified FROM payment_sessions WHERE id = ?"
	session := &model.PaymentSession{}
	if err := repo.store.db.QueryRow(selectQuery, ID).Scan(
		&session.ID,
		&session.KeyID,
		&session.Price,
		&session.CustomerName,
		&session.CustomerEmail,
		&session.CustomerAddress,
		&session.IsPerformed,
		&session.IsNotified); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return session, nil
}

func (repo *PaymentSessionRepository) Create(session *model.PaymentSession) error {
	insertQuery := "INSERT INTO payment_sessions (key_id, price, date, customer_name, customer_email, customer_address) VALUES (?, ?, ?, ?, ?, ?);"
	selectPriceQuery := "SELECT price FROM games JOIN `keys` k ON games.id = k.game_id WHERE k.id = ?"
	getIdQuery := "select LAST_INSERT_ID();"

	if err := repo.store.db.QueryRow(selectPriceQuery, session.KeyID).Scan(&session.Price); err != nil {
		return err
	}

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

func (repo *PaymentSessionRepository) DeleteByID(id int) error {
	deleteQuery := "DELETE FROM payment_sessions WHERE id = ?"
	result, err := repo.store.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return store.ErrRecordNotFound
	}

	return nil
}

func (repo *PaymentSessionRepository) GetPaymentInfo(sessionID int) (*model.PaymentInfo, error) {
	info := &model.PaymentInfo{}
	selectQuery :=
		"SELECT s.account, s.url, g.title, g.price, k.id, k.key_string, p.customer_name, p.customer_email, p.customer_address " +
			"FROM sellers s " +
			"JOIN games g on s.id = g.seller_id " +
			"JOIN `keys` k on g.id = k.game_id " +
			"JOIN payment_sessions p on k.id = p.key_id " +
			"WHERE p.id = ?"

	if err := repo.store.db.QueryRow(selectQuery, sessionID).Scan(
		&info.SellerAccount,
		&info.SellerURL,
		&info.GameTitle,
		&info.TotalAmount,
		&info.KeyID,
		&info.Key,
		&info.CustomerName,
		&info.CustomerEmail,
		&info.CustomerAddress); err != nil {
		return nil, err
	}
	return info, nil
}

func (repo *PaymentSessionRepository) UpdateNotifiedStatus(sessionID int, newStatus bool) error {
	updateQuery := "UPDATE payment_sessions SET is_notified = TRUE WHERE id = ?"
	result, err := repo.store.db.Exec(updateQuery, sessionID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return store.ErrRecordNotFound
	}
	if err != nil {
		return err
	}
	return nil
}

func (repo *PaymentSessionRepository) UpdatePerformedStatus(sessionID int, newStatus bool) error {
	updateQuery := "UPDATE payment_sessions SET is_performed = TRUE WHERE id = ?"
	result, err := repo.store.db.Exec(updateQuery, sessionID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if rows == 0 {
		return store.ErrRecordNotFound
	}
	if err != nil {
		return err
	}
	return nil
}
