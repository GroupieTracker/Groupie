package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Game struct {
    ID   int
    Name string
}

type Room struct {
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

func CreateRoomAndGetID(db *sql.DB, room Room) (int, error) {
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
