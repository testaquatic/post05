package post05

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Userdata 테이블
type Userdata struct {
	ID          int32
	Username    string
	Name        string
	Surname     string
	Description string
}

// 연결 상세 정보
var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*pgxpool.Pool, error) {
	conn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s", Username, Password, Hostname, Port, Database,
	)
	return pgxpool.New(context.Background(), conn)
}

// 사용자 이름을 받아 ID를 반환한다.
// 사용자가 존재하지 않으면 -1을 반환한다.
func exists(username string) int32 {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := int32(-1)
	statement := fmt.Sprintf(
		`SELECT id FROM users WHERE username = '%s'`, username,
	)
	rows, err := db.Query(context.Background(), statement)

	for rows.Next() {
		var id int32
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}
	defer rows.Close()

	return userID
}

// 데이터베이스에 새로운 사용자를 추가하고 해당 사용자의 User ID를 반환한다.
// 에러가 발생하면 -1을 반환한다.
func AddUser(d Userdata) int32 {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", Username)
		return -1
	}

	insertStatement := `INSERT INTO users (username) VALUES ($1)`
	_, err = db.Exec(context.Background(), insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	insertStatement = `
		INSERT INTO 
			userdata (userid, name, surname, description) 
		VALUES
			($1, $2, $3, $4)`
	_, err = db.Exec(context.Background(), insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return userID
}

// 존재하는 사용자를 지운다.
func DeleteUser(id int32) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	statement := fmt.Sprintf(
		`SELECT username FROM users WHERE id = %d`,
		id,
	)
	rows, err := db.Query(context.Background(), statement)
	if err != nil {
		return err
	}
	defer rows.Close()

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exists", id)
	}
	deleteStatement := `DELETE FROM userdata WHERE userid=$1`
	_, err = db.Exec(context.Background(), deleteStatement, id)
	if err != nil {
		return err
	}

	deleteStatement = `DELETE FROM users WHERE id = $1`
	_, err = db.Exec(context.Background(), deleteStatement, id)

	return err
}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(
		context.Background(),
		`
		SELECT id, username, name, surname, description
		FROM users
		JOIN userdata
		ON users.id = userdata.userid
		`,
	)
	defer rows.Close()

	for rows.Next() {
		var id int32
		var username string
		var name string
		var surname string
		var description string
		err := rows.Scan(&id, &username, &name, &surname, &description)
		if err != nil {
			return Data, err
		}
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: description}
		Data = append(Data, temp)
	}

	return Data, nil
}

// 존재하는 사용자를 업데이트한다.
func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("User does not exists")
	}

	d.ID = userID
	updateStatement := `
		UPDATE userdata
		SET name = $1, surname = $2, description = $3
		WHERE userid = $4
	`
	_, err = db.Exec(context.Background(), updateStatement, d.Name, d.Surname, d.Description, d.ID)

	return err
}
