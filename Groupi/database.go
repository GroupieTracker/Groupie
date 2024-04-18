package Groupi

import (
	"database/sql"
	"log"
    "fmt"

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



func GetRoomCreatorID(db *sql.DB , roomID string) (int, error) {
    var creatorID int

    err := db.QueryRow("SELECT created_by FROM ROOMS WHERE id = ?", roomID).Scan(&creatorID)
    if err != nil {
        fmt.Println("paric")
        return 0, err
    }

    return creatorID, nil
}

func GetUsernameByID(db *sql.DB , userID int) (string, error) {
    var username string

    err := db.QueryRow("SELECT pseudo FROM USER WHERE id = ?", userID).Scan(&username)
    if err != nil {
        return "", err
    }

    return username, nil
}

func AddRoomUser(db *sql.DB ,roomID int, userID int) error {
    _, err := db.Exec("INSERT INTO ROOM_USERS (id_room, id_user, score) VALUES (?, ?, ?)", roomID, userID, 0)
    if err != nil {
        return err
    }

    return nil
}