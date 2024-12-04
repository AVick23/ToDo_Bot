package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Подключаем драйвер для PostgreSQL
)

type DB struct {
	*sql.DB
}

// ConnectDB подключается к PostgreSQL
func CreateDB() (*sql.DB, error) {
	connStr := "postgres://avick123:super123@127.0.0.1:5432/todo_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("не удалось установить соединение: %w", err)
	}

	createTableDB := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT NOT NULL UNIQUE
        )`,
		`CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            user_id INTEGER NOT NULL,
            tasks TEXT,
            date TEXT,
            list_name TEXT,
            notification TEXT,
            FOREIGN KEY(user_id) REFERENCES users(id)
        )`,
	}

	for _, query := range createTableDB {
		_, err := db.Exec(query)
		if err != nil {
			return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
		}
	}
	return db, nil
}

func ConnectDB() (*sql.DB, error) {
	connStr := "postgres://avick123:super123@127.0.0.1:5432/todo_db?sslmode=disable"
	return sql.Open("postgres", connStr)
}

func SaveUser(db *sql.DB, username string) (int, error) {
	var id int
	err := db.QueryRow("INSERT INTO users (username) VALUES ($1) ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username RETURNING id", username).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("ошибка сохранения пользователя: %w", err)
	}
	return id, nil
}

func SaveTasks(db *sql.DB, userID int, taskText string) error {
	_, err := db.Exec("INSERT INTO tasks (user_id, tasks) VALUES ($1, $2)", userID, taskText)
	if err != nil {
		return fmt.Errorf("ошибка сохранения задачи: %w", err)
	}
	return nil
}

func GetTasks(db *sql.DB, username string) ([]string, error) {
	var tasks []string
	var id int
	err := db.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения userID: %w", err)
	}
	rows, err := db.Query("SELECT tasks FROM tasks WHERE user_id = $1", id)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнение задачи %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task string
		err := rows.Scan(&task)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования задач %v", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации %v", err)
	}
	return tasks, nil
}
