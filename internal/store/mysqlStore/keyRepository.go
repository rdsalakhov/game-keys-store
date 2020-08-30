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
	selectQuery := "SELECT id, key_string, game_id, status FROM `keys` WHERE id = ?"
	key := &model.Key{}
	if err := repo.store.db.QueryRow(selectQuery, ID).Scan(
		&key.ID,
		&key.KeyString,
		&key.GameID,
		&key.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return key, nil
}

func (repo *KeyRepository) FindSingleAvailable(gameID int) (*model.Key, error) {
	selectQuery := "SELECT id, key_string, game_id, status FROM `keys` WHERE game_id = ? && status = 'available' ORDER BY id LIMIT 1"
	key := &model.Key{}
	if err := repo.store.db.QueryRow(selectQuery, gameID).Scan(
		&key.ID,
		&key.KeyString,
		&key.GameID,
		&key.Status); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return key, nil
}

func (repo *KeyRepository) FindByGameID(gameID int) (*[]model.Key, error) {
	selectQuery := "SELECT id, key_string, game_id, status FROM `keys` WHERE game_id = ?"

	rows, err := repo.store.db.Query(selectQuery, gameID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	defer rows.Close()
	keys := []model.Key{}

	for rows.Next() {
		key := model.Key{}
		if err := rows.Scan(&key.ID,
			&key.KeyString,
			&key.GameID,
			&key.Status); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}
	return &keys, nil
}

func (repo *KeyRepository) UpdateStatus(keyID int, newStatus string) error {
	updateQuery := "UPDATE `keys` SET status = ? WHERE id = ?"
	result, err := repo.store.db.Exec(updateQuery,
		newStatus,
		keyID)
	if err != nil {
		return err
	}
	if rows, err := result.RowsAffected(); err != nil {
		return err
	} else if rows == 0 {
		return store.ErrRecordNotFound
	}

	return nil
}

func (repo *KeyRepository) Create(key *model.Key) error {
	insertQuery := "INSERT INTO `keys` (key_string, game_id) VALUES (?, ?);"
	getIdQuery := "select LAST_INSERT_ID();"

	if _, err := repo.store.db.Exec(insertQuery,
		key.KeyString,
		key.GameID,
	); err != nil {
		return err
	}

	if err := repo.store.db.QueryRow(getIdQuery).Scan(&key.ID); err != nil {
		return err
	}
	return nil
}
