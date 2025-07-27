package queries

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type GetUsersQueryRow struct {
	ID                 int         `db:"id"`
	Username           string      `db:"username"`
	Email              string      `db:"email"`
	UserType           string      `db:"user_type"`
	Nickname           pgtype.Text `db:"nickname"`
	PermissionBitfield string      `db:"permission_bitfield"`
	MessageCount       int         `db:"message_count"`
}

func GetUsers() ([]GetUsersQueryRow, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	rows, err := conn.Query(context.TODO(), `
		SELECT 
			u.id, 
			u.username, 
			u.email,
			ut.type_key as "user_type",
			u.nickname,
			ut.permission_bitfield::text as "permission_bitfield",
			COUNT(m.id) AS message_count
		FROM
			public.users u
		LEFT JOIN 
			public.messages m ON u.username = m.username
		LEFT JOIN 
			public.user_types ut ON u.user_type = ut.type_key    
		GROUP BY 
			u.id, u.username, u.email, ut.type_key, u.nickname, ut.permission_bitfield
		ORDER BY 
			u.id;
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []GetUsersQueryRow = []GetUsersQueryRow{}
	for rows.Next() {
		var user GetUsersQueryRow
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.UserType,
			&user.Nickname,
			&user.PermissionBitfield,
			&user.MessageCount,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, err
}

type CreateUserParams struct {
	Username string
	Email    string
	UserType string
	Nickname string
}

func CreateUser(params CreateUserParams) (int, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	var userID int
	err := conn.QueryRow(
		context.TODO(),
		`INSERT INTO public.users (username, email, user_type, nickname)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		params.Username,
		params.Email,
		params.UserType,
		params.Nickname,
	).Scan(&userID)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

type InsertMessageParams struct {
	Username string
	Message  string
}

func InsertMessage(params InsertMessageParams) (int, error) {
	conn := GetConnection()
	defer conn.Close(context.TODO())

	var messageID int
	err := conn.QueryRow(
		context.TODO(),
		`INSERT INTO public.messages (username, message)
		 VALUES ($1, $2)
		 RETURNING id`,
		params.Username,
		params.Message,
	).Scan(&messageID)

	if err != nil {
		return 0, err
	}

	return messageID, nil
}
