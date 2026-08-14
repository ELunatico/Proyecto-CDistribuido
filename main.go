package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/lib/pq" // Driver de PostgreSQL
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

var db *sql.DB

func main() {
	connStr := "host=db port=5432 user=postgres password=admin dbname=tareas sslmode=disable"
	
	var err error
	for i := 0; i < 5; i++ {
		db, err = sql.Open("postgres", connStr)
		if err == nil && db.Ping() == nil {
			break
		}
		log.Println("Esperando a la base de datos...")
		time.Sleep(2 * time.Second)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS tasks (id SERIAL PRIMARY KEY, title TEXT)`)
	if err != nil {
		log.Fatal("Error creando tabla:", err)
	}

	http.HandleFunc("/api/tasks", tasksHandler)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	log.Println("Servidor Go corriendo en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		rows, _ := db.Query("SELECT id, title FROM tasks")
		defer rows.Close()

		var tasks []Task
		for rows.Next() {
			var t Task
			rows.Scan(&t.ID, &t.Title)
			tasks = append(tasks, t)
		}
		json.NewEncoder(w).Encode(tasks)

	case "POST":
		var t Task
		json.NewDecoder(r.Body).Decode(&t)
		db.QueryRow("INSERT INTO tasks (title) VALUES ($1) RETURNING id", t.Title).Scan(&t.ID)
		json.NewEncoder(w).Encode(t)

	case "PUT":
		var t Task
		json.NewDecoder(r.Body).Decode(&t)
		db.Exec("UPDATE tasks SET title = $1 WHERE id = $2", t.Title, t.ID)
		json.NewEncoder(w).Encode(t)

	case "DELETE":
		id, _ := strconv.Atoi(r.URL.Query().Get("id"))
		db.Exec("DELETE FROM tasks WHERE id = $1", id)
		w.WriteHeader(http.StatusOK)
	}
}