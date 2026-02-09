package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"awesomeProject/internal/model"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDB подключает к PostgreSQL и создаёт таблицу с нужными данными
func InitDB() {
	connStr := "host=localhost port=5444 user=postgres password=secret dbname=postgres sslmode=disable"

	var err error

	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic(err)
	}

	err = DB.Ping()

	if err != nil {
		fmt.Println("Ошибка подключения к Postgres:", err)
		panic(err)
	}

	fmt.Println("Успешно подключились к PostgreSQL!")

	createDBTable := `CREATE TABLE IF NOT EXISTS players (
    name TEXT NOT NULL PRIMARY KEY,
    password TEXT,
    gold INTEGER,
    equipment TEXT,
    inventory TEXT
)`

	_, err = DB.Exec(createDBTable)

	if err != nil {
		fmt.Println("Ошибка создания таблицы:", err)
		panic(err)
	}
}

// SaveGame сохраняет данные
func SaveGame(p model.Player) error {
	equipBytes, _ := json.Marshal(p.Equipment)
	invBytes, _ := json.Marshal(p.Inventory)

	equipString := string(equipBytes)
	invString := string(invBytes)

	query := `INSERT INTO players (name, password, gold, equipment, inventory)
 			  VALUES ($1, $2, $3, $4, $5)
 			  ON CONFLICT (name) DO UPDATE SET
 			  password = EXCLUDED.password,
 			  gold = EXCLUDED.gold,
 			  equipment = EXCLUDED.equipment,
 			  inventory = EXCLUDED.inventory`

	_, err := DB.Exec(query, p.Name, p.Password, p.Gold, equipString, invString)

	if err != nil {
		return fmt.Errorf("ошибка сохранения: для игрока %s: %w", p.Name, err)
	}

	return nil
}

// LoadGame загружает данные
func LoadGame(name string) (model.Player, error) {
	var p model.Player
	var equipment string
	var inventory string

	query := `SELECT name, password, gold, equipment, inventory FROM players WHERE name = $1`
	row := DB.QueryRow(query, name)
	err := row.Scan(&p.Name, &p.Password, &p.Gold, &equipment, &inventory)

	if err != nil {
		return model.Player{}, err // Возвращаем ошибку, если игрока нет
	}

	json.Unmarshal([]byte(equipment), &p.Equipment)
	json.Unmarshal([]byte(inventory), &p.Inventory)

	if p.Equipment == nil {
		p.Equipment = make(map[string]string)
	}

	return p, nil
}

// HasName проверяет наличие имени в players
func HasName(name string) bool {
	var p model.Player

	query := `SELECT name FROM players WHERE name = $1`
	row := DB.QueryRow(query, name)
	err := row.Scan(&p.Name)

	return err == nil
}
