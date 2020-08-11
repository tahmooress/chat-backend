package dbf

import (
	"database/sql"
	"fmt"
)

//RunDB open a postgres database
func RunDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", SetDbInfo())
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "900587101"
	dbname   = "test"
)

//SetDbInfo sets the required information for connecting to database
func SetDbInfo() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

//required statment for working with DataBase
const (
	RegisterNewUserStatement = `INSERT INTO users (username, email, password) VALUES ($1, $2, $3)`
	UpdateProfileStatement   = `UPDATE users SET bio = $1, image= $2 WHERE email= $3`
	NewMessageStatement      = `INSERT INTO messages (sender_email, reciever_email, message, created_at) VALUES ($1, $2, $3, $4)`
	SelectUserID             = `SELECT user_id FROM users WHERE email= $1 `
	SelectAllRecMes          = `SELECT message FROM messages WHERE sender_id= $1 AND reciever_id= $2`
	CheckExist               = `SELECT email, password FROM users WHERE email=$1`
)

//Env is enviorment variable for connecting to db
type Env struct {
	DB *sql.DB
}
