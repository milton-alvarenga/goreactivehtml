package business

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var Description string

type Task struct {
	Id          int
	Description string
	Done        bool
}

var Tasks []Task

var db *sql.DB

// ConnectDB establishes a connection to the PostgreSQL database.
func ConnectDB() error {
	var err error
	// Database connection string
	connStr := "user=yourusername password=yourpassword dbname=yourdb sslmode=disable" // Update with your DB credentials

	// Open the connection
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening DB connection: %v", err)
	}

	// Ping the database to ensure the connection works
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error connecting to DB: %v", err)
	}

	log.Println("Successfully connected to the database")
	return nil
}

func main() {
	ConnectDB()
}

func Load(user_id int) {
	query := "SELECT id, description, done FROM task WHERE user_id = $1"

	rows, _ := db.Conn.Query(query, user_id)
	defer rows.Close()

	// Iterate through the rows
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Description, &task.Done)
		if err != nil {
			log.Fatal(err)
		}
		Tasks = append(Tasks, task)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func Add(description string) {
	var id int

	query := "INSERT INTO task (description) VALUE ($1) RETURNING id"

	db.Conn.QueryRow(query, description).Scan(&id)
	Description = ""
	Tasks = append(Tasks, Task{Id: id, Description: description, Done: false})
}

func Update(task Task) bool {
	query := "UPDATE task SET done = $1 WHERE id = $2"
	_, err := db.Conn.Exec(query, task.Done, task.Id)
	return err == nil
}

func Delete(id int) bool {
	query := "DELETE FROM task WHERE id = $1"
	_, err := db.Conn.Exec(query, id)

	return err == nil
}
