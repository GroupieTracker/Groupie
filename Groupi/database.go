package Groupi

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Game struct {
	ID   int
	Name string
}

type Roomms struct {
	ID         int
	CreatedBy  int
	MaxPlayers int
	Name       string
	GameID     int
	Category   []string
}

func InitializeDatabase() {
	db, err := sql.Open("sqlite3", "BDD.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS USER (
        id INTEGER PRIMARY KEY,
        pseudo TEXT NOT NULL,
        email TEXT NOT NULL,
        password TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Tables créées avec succès dans la base de données.")
}

func CreateRoomAndGetID(db *sql.DB, room Roomms) (int, error) {
	query := "INSERT INTO ROOMS (created_by, max_player, name, id_game) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, room.CreatedBy, room.MaxPlayers, room.Name, room.GameID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func CreateGameAndGetID(db *sql.DB, game Game) (int, error) {
	query := "INSERT INTO GAMES (name) VALUES (?)"
	result, err := db.Exec(query, game.Name)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func GetUserIDByUsername(db *sql.DB, username string) (int, error) {
	var userID int
	query := "SELECT id FROM USER WHERE pseudo = ?"
	err := db.QueryRow(query, username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("l'utilisateur avec le pseudo %s n'existe pas", username)
		}
		return 0, err
	}

	return userID, nil
}

func DeleteRoomAndRoomUsersByID(db *sql.DB, roomID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	queryDeleteRoom := "DELETE FROM ROOMS WHERE id = ?"
	_, err = tx.Exec(queryDeleteRoom, roomID)
	if err != nil {
		return err
	}
	queryDeleteRoomUsers := "DELETE FROM ROOM_USERS WHERE id_room = ?"
	_, err = tx.Exec(queryDeleteRoomUsers, roomID)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func GetUsersInRoom(db *sql.DB, roomID string) ([]int, error) {
	var userIDs []int

	rows, err := db.Query("SELECT id_user FROM ROOM_USERS WHERE id_room = ?", roomID)
	if err != nil {

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {

			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	if err := rows.Err(); err != nil {

		return nil, err
	}

	return userIDs, nil
}

func GetRoomCreatorID(db *sql.DB, roomID string) (int, error) {
	var creatorID int

	err := db.QueryRow("SELECT created_by FROM ROOMS WHERE id = ?", roomID).Scan(&creatorID)
	if err != nil {
		fmt.Println("paric")
		return 0, err
	}

	return creatorID, nil
}

func GetUsernameByID(db *sql.DB, userID int) (string, error) {
	var username string

	err := db.QueryRow("SELECT pseudo FROM USER WHERE id = ?", userID).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func AddRoomUser(db *sql.DB, roomID int, userID int) error {
	_, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)", roomID, userID, 0)
	if err != nil {
		return err
	}

	return nil
}

func UpdateRoomUserScore(db *sql.DB, roomID, userID, scoreToAdd int) error {

	var currentScore int
	err := db.QueryRow("SELECT score FROM ROOM_USERS WHERE id_room = ? AND id_user = ?", roomID, userID).Scan(&currentScore)
	if err != nil {
		return err
	}
	newScore := currentScore + scoreToAdd
	_, err = db.Exec("UPDATE ROOM_USERS SET score = ? WHERE id_room = ? AND id_user = ?", newScore, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}
func GetUserScoresForRoom(db *sql.DB, userIDs []int, roomID int) ([][]string, error) {
	var userScores [][]string
	var userIDsStr string
	for i, userID := range userIDs {
		if i > 0 {
			userIDsStr += ","
		}
		userIDsStr += strconv.Itoa(userID)
	}
	query := `
        SELECT u.pseudo, ru.score
        FROM USER u
        INNER JOIN ROOM_USERS ru ON u.id = ru.id_user
        WHERE ru.id_user IN (` + userIDsStr + `) AND ru.id_room = ?`
	rows, err := db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var score int
		if err := rows.Scan(&username, &score); err != nil {
			return nil, err
		}
		userScores = append(userScores, []string{username, strconv.Itoa(score)})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return userScores, nil
}

func GetMaxPlayersForRoom(db *sql.DB, roomID int) (int, error) {
	var maxPlayers int

	err := db.QueryRow("SELECT max_player FROM ROOMS WHERE id = ?", roomID).Scan(&maxPlayers)
	if err != nil {
		return 0, err
	}

	return maxPlayers, nil
}

func GetUsernameByEmailOrUsername(db *sql.DB, identifier string) (string, error) {
	var username string
	query := `
        SELECT pseudo FROM USER
        WHERE email = ? OR pseudo = ?`
	err := db.QueryRow(query, identifier, identifier).Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func ClearDatabase(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM ROOM_USERS")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM ROOMS")
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM GAMES")
	if err != nil {
		return err
	}
	// _, err = db.Exec("DELETE FROM USER")
	// if err != nil {
	//     return err
	// }
	return nil
}

func GetRoomsByGameCategory(db *sql.DB, categoryName string) ([][]string, error) {
	var rooms [][]string
	query := `
        SELECT r.name, r.id
        FROM ROOMS r
        INNER JOIN GAMES g ON r.id_game = g.id
        WHERE g.name = ?`
	rows, err := db.Query(query, categoryName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var roomName string
		var roomID int
		if err := rows.Scan(&roomName, &roomID); err != nil {
			return nil, err
		}
		rooms = append(rooms, []string{roomName, strconv.Itoa(roomID)})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

func GetUserScoresForRoomID(db *sql.DB, roomID int) ([][]string, error) {
	var userScores [][]string

	query := `
        SELECT u.pseudo, ru.score
        FROM USER u
        INNER JOIN ROOM_USERS ru ON u.id = ru.id_user
        WHERE ru.id_room = ?`
	rows, err := db.Query(query, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var score int
		if err := rows.Scan(&username, &score); err != nil {
			return nil, err
		}
		userScores = append(userScores, []string{username, strconv.Itoa(score)})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userScores, nil
}

func DeleteRoomData(db *sql.DB, roomID int) error {
	_, err := db.Exec("DELETE FROM ROOM_USERS WHERE id_room = ?", roomID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM ROOMS WHERE id = ?", roomID)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM GAMES WHERE id IN (SELECT id_game FROM ROOMS WHERE id = ?)", roomID)
	if err != nil {
		return err
	}

	return nil
}
