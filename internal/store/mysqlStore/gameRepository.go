package mysqlStore

import (
	"database/sql"
	"github.com/rdsalakhov/game-keys-store/internal/model"
	"github.com/rdsalakhov/game-keys-store/internal/store"
)

type GameRepository struct {
	store *Store
}

func (repo *GameRepository) Find(ID int) (*model.Game, error) {
	selectQuery := "SELECT id, title, description, price, on_sale, seller_id FROM games WHERE id = ?"
	game := &model.Game{}
	if err := repo.store.db.QueryRow(selectQuery, ID).Scan(
		&game.ID,
		&game.Title,
		&game.Description,
		&game.Price,
		&game.OnSale,
		&game.SellerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return game, nil
}

func (repo *GameRepository) FindByTitle(title string) (*model.Game, error) {
	selectQuery := "SELECT id, title, description, price, on_sale, seller_id FROM games WHERE title = ?"
	game := &model.Game{}
	if err := repo.store.db.QueryRow(selectQuery, title).Scan(
		&game.ID,
		&game.Title,
		&game.Description,
		&game.Price,
		&game.OnSale,
		&game.SellerID); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return game, nil
}

func (repo *GameRepository) FindAll() ([]*model.Game, error) {
	selectQuery := "SELECT id, title, description, price, on_sale, seller_id FROM games"
	games := []*model.Game{}
	rows, err := repo.store.db.Query(selectQuery)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		game := &model.Game{}
		if err := rows.Scan(
			&game.ID,
			&game.Title,
			&game.Description,
			&game.Price,
			&game.OnSale,
			&game.SellerID); err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil

}

func (repo *GameRepository) Create(game *model.Game) error {
	insertQuery := "INSERT INTO games (title, description, price, on_sale, seller_id) VALUES (?, ?, ?, ?, ?);"
	getIdQuery := "select LAST_INSERT_ID();"

	if _, err := repo.store.db.Exec(insertQuery,
		game.Title,
		game.Description,
		game.Price,
		game.OnSale,
		game.SellerID,
	); err != nil {
		return err
	}

	if err := repo.store.db.QueryRow(getIdQuery).Scan(&game.ID); err != nil {
		return err
	}
	return nil
}
