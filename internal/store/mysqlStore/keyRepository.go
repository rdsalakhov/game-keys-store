package mysqlStore

import (
	"database/sql"
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store"
)

type KeyRepository struct {
	store *Store
}

func (repo *KeyRepository) Find(ID int) (*model.Key, error) {
	selectQuery := "SELECT id, key_string, game_id, seller_id, status FROM `keys` WHERE id = ?"
	key := &model.Key{}
	if err := repo.store.db.QueryRow(selectQuery, ID).Scan(
		&key.ID,
		&key.KeyString,
		&key.GameID,
		&key.SellerID,
		&key.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return key, nil
}

func (repo *KeyRepository) Create(key *model.Key) error {
	insertQuery := "INSERT INTO `keys` (key_string, game_id, seller_id) VALUES (?, ?, ?);"
	getIdQuery := "select LAST_INSERT_ID();"

	if _, err := repo.store.db.Exec(insertQuery,
		key.KeyString,
		key.GameID,
		key.SellerID,
	); err != nil {
		return err
	}

	if err := repo.store.db.QueryRow(getIdQuery).Scan(&key.ID); err != nil {
		return err
	}
	return nil
}
