package mysqlStore

import (
	"database/sql"
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store"
)

type SellerRepository struct {
	store *Store
}

func (repo *SellerRepository) Find(ID int) (*model.Seller, error) {
	selectQuery := "SELECT id, url, account, email, encrypted_password FROM sellers WHERE id = ?"
	seller := &model.Seller{}
	if err := repo.store.db.QueryRow(selectQuery, ID).Scan(
		&seller.ID,
		&seller.URL,
		&seller.Account,
		&seller.Email,
		&seller.EncryptedPassword); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return seller, nil
}

func (repo *SellerRepository) Create(seller *model.Seller) error {
	insertQuery := "INSERT INTO sellers (url, account, email, encrypted_password) VALUES (?, ?, ?, ?);"
	getIdQuery := "select LAST_INSERT_ID();"

	if _, err := repo.store.db.Exec(insertQuery,
		seller.URL,
		seller.Account,
		seller.Email,
		seller.EncryptedPassword); err != nil {
		return err
	}

	if err := repo.store.db.QueryRow(getIdQuery).Scan(&seller.ID); err != nil {
		return err
	}
	return nil
}
