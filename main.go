package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "godo"
)

type Task struct {
	name string
	body string
	done bool
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	createTable(db)

	choice := menu()
	// check their choice
	switch choice {
	case "a", "A":
		addData(db)
	case "d", "D":
		fmt.Println("Delete!")
	case "m", "M":
		fmt.Println("Modify!")
	default:
		fmt.Println("Invalid input, please try again")
	}
}

func createTable(db *sql.DB) {
	// ID
	// Task name
	// task body
	// is it done
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    task VARCHAR(20) NOT NULL,
    body VARCHAR(120),
    done BOOLEAN
  )
  `
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertTask(db *sql.DB, task Task) int {
	query := `INSERT INTO tasks (task, body, done)
    VALUES ($1, $2, $3) RETURNING id`

	var tk int
	err := db.QueryRow(query, task.name, task.body, task.done)
	if err != nil {
		log.Fatal(err)
	}
	return tk
}

func menu() string {
	fmt.Println("Welcome!")
	fmt.Println("You can (a)dd, (d)elete, or (m)odify tasks")
	fmt.Print("What would you like to do (a,d,m): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	return scanner.Text()
}

func addData(db *sql.DB) {
	fmt.Println("Adding task")
	var new_task Task
	scanner := bufio.NewScanner(os.Stdin)
	err := scanner.Err()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("Task Name: ")
	scanner.Scan()
	new_task.name = scanner.Text()
	fmt.Print("Task Description: ")
	scanner.Scan()
	new_task.body = scanner.Text()

	tk := insertTask(db, new_task)

	fmt.Printf("New Task ID: %d\n", tk)
}
