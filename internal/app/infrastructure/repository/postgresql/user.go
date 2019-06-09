package postgresql

import (
	"errors"

	"github.com/ZorinArsenij/tech-db-forum/internal/app/domain/user"
	"github.com/jackc/pgx"
)

const (
	getUserByNickname              = "getUserByNickname"
	updateUser                     = "updateUser"
	getUsersWithEmailAndNickname   = "getUsersWithEmailAndNickname"
	createUser                     = "createUser"
	getUserIdAndNicknameByNickname = "getUserIdAndNicknameByNickname"
)

var userQueries = map[string]string{
	getUserByNickname: `SELECT email, nickname, fullname, about 
	FROM client WHERE nickname = $1;`,

	updateUser: `UPDATE client 
	SET email = COALESCE($1, email),
		fullname = COALESCE($2, fullname),
		about = COALESCE($3, about)
	WHERE nickname = $4 
	RETURNING email, nickname, fullname, about;`,

	getUsersWithEmailAndNickname: `SELECT email, nickname, fullname, about 
	FROM client 
	WHERE email = $1 OR nickname = $2;`,

	createUser: `INSERT INTO client (email, nickname, fullname, about)
	VALUES ($1, $2, $3, $4)
	RETURNING email, nickname, fullname, about;`,

	getUserIdAndNicknameByNickname: `SELECT id, nickname
	FROM client
	WHERE nickname = $1;`,
}

func NewUserRepo(conn *pgx.ConnPool) *User {
	return &User{
		conn: conn,
	}
}

type User struct {
	conn *pgx.ConnPool
}

func (u *User) GetUserByNickname(nickname string) (*user.User, error) {
	received := &user.User{}
	if err := u.conn.QueryRow(getUserByNickname, nickname).Scan(&received.Email, &received.Nickname, &received.Fullname, &received.About); err != nil {
		return nil, err
	}

	return received, nil
}

func (u *User) UpdateUser(data *user.Update, nickname string) (*user.User, error) {
	updated := &user.User{}

	tx, err := u.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err = tx.QueryRow(updateUser, data.Email, data.Fullname, data.About, nickname).Scan(&updated.Email, &updated.Nickname, &updated.Fullname, &updated.About); err != nil {
		return nil, err
	}

	tx.Commit()
	return updated, nil
}

func (u *User) CreateUser(data *user.User) (*user.Users, error) {
	tx, err := u.conn.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	users := make(user.Users, 0, 2)

	rows, _ := tx.Query(getUsersWithEmailAndNickname, data.Email, data.Nickname)

	for rows.Next() {
		var row user.User
		rows.Scan(&row.Email, &row.Nickname, &row.Fullname, &row.About)
		users = append(users, row)
	}
	rows.Close()

	if len(users) == 0 {
		var created user.User
		if err := tx.QueryRow(createUser, data.Email, data.Nickname, data.Fullname, data.About).Scan(&created.Email, &created.Nickname, &created.Fullname, &created.About); err != nil {
			return nil, err
		}
		users = append(users, created)

		tx.Commit()
		return &users, nil
	}

	return &users, errors.New("conflict")
}
