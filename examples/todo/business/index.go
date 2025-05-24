package business

import(
	"database/sql"

	_ "github.com/lib/pq"
)

var Description string

type Task struct {
    Id int
    Description string
    Done bool
}

var Tasks []Task{}


func main() {
    
}

func Load(){
	query := "SELECT id, description, done FROM task WHERE user_id = $1"
	
    rows,_ = db.Conn.Query(query,user_id)
    defer rows.Close()
	
    // Iterate through the rows
    for rows.Next() {
		var task Task
        err := rows.Scan(&task.Id, &task.Description, &task.Done)
        if err != nil {
            log.Fatal(err)
        }
        Tasks = append(Tasks,task)
    }

    // Check for errors from iterating over rows
    if err = rows.Err(); err != nil {
        log.Fatal(err)
    }
}

func Add(Description string) {
	var id int

	    query := "INSERT INTO task (description) VALUE ($1) RETURNING id"

    db.Conn.QueryRow(query,Description).Scan(&id)
    Description = ""
    Tasks = append(Tasks,Task{Id:id,Description:description, Done: false})
}

func Update(task Task) bool {
    query := "UPDATE task SET done = $1 WHERE id = $2"
    _, err := db.Conn.Exec(query,task.Done,task.Id)
    return err == nil
}

func Delete(id int) bool {
    query = "DELETE FROM task WHERE id = $1"
    _, err := db.Conn.Exec(query, id)

    return err == nil
}
