package user

const (
	createUser = `INSERT INTO client (email, nickname, fullname, about) 
	VALUES ($1, $2, $3, $4) 
	ON CONFLICT DO NOTHING;`

	getUserIdAndNicknameByNickname = `SELECT id, nickname
	FROM client
	WHERE nickname = $1`

	getUserByNickname = `SELECT email, nickname, fullname, about 
	FROM client WHERE nickname = $1;`

	updateUser = `UPDATE client 
	SET email = (CASE WHEN $1 = '' THEN email ELSE $1 END),
		fullname = (CASE WHEN $2 = '' THEN fullname ELSE $2 END),
		about = (CASE WHEN $3 = '' THEN about ELSE $3 END)
	WHERE nickname = $4 
	RETURNING email, nickname, fullname, about;`

	getUsersWithEmailAndNickname = `SELECT email, nickname, fullname, about 
	FROM client 
	WHERE email = $1 OR nickname = $2;`
)
