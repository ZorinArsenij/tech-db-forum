package database

import (
	"github.com/ZorinArsenij/tech-db-forum/models"
)

const (
	createUser = `INSERT INTO client (email, nickname, fullname, about) 
	VALUES ($1, $2, $3, $4) 
	ON CONFLICT DO NOTHING;`

	getUserByNickname = `SELECT email, nickname, fullname, about 
	FROM client WHERE nickname = $1;`

	updateUser = `UPDATE client 
	SET email = (CASE WHEN $1 = '' THEN email ELSE $1 END),
		fullname = (CASE WHEN $2 = '' THEN fullname ELSE $2 END),
		about = (CASE WHEN $3 = '' THEN about ELSE $3 END)
	WHERE nickname = $4 
	RETURNING email, nickname, fullname, about;`

	findUsersWithEmailAndNickname = `SELECT email, nickname, fullname, about 
	FROM client 
	WHERE email = $1 OR nickname = $2;`
)

func (manager *Manager) CreateUser(client *models.Client) (*models.Clients, error) {
	tx, err := manager.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(createUser, client.Email, client.Nickname, client.Fullname, client.About)
	if err != nil {
		return nil, err
	}

	if result.RowsAffected() == 0 {
		clients := make(models.Clients, 0, 2)

		rows, err := tx.Query(findUsersWithEmailAndNickname, client.Email, client.Nickname)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			client := models.Client{}
			rows.Scan(&client.Email, &client.Nickname, &client.Fullname, &client.About)
			clients = append(clients, client)
		}
		return &clients, nil
	}

	tx.Commit()
	return nil, nil
}

func (manager *Manager) GetUser(nickname string) (*models.Client, error) {
	client := &models.Client{}
	err := manager.conn.QueryRow(getUserByNickname, nickname).Scan(&client.Email, &client.Nickname, &client.Fullname, &client.About)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (manager *Manager) UpdateUser(client *models.Client) error {
	tx, err := manager.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow(updateUser, client.Email, client.Fullname, client.About, client.Nickname).Scan(&client.Email, &client.Nickname, &client.Fullname, &client.About)
	if err != nil {
		return err
	}

	tx.Commit()
	return nil
}
