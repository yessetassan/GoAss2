package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Qasaqayj7"
	dbname   = "fix"
)

type Task struct {
	ID        int
	Name      string
	Completed bool
}

func main() {
	db := connectDB()
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)
FOOR_LOOP_FOR_DB:
	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. Create a new task")
		fmt.Println("2. Read all tasks")
		fmt.Println("3. Update a task")
		fmt.Println("4. Delete a task")
		fmt.Println("5. Exit")

		option, _ := reader.ReadString('\n')
		option = strings.TrimSpace(option)

		switch option {
		case "1":
			fmt.Println("Enter task name:")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			err := createTask(db, name)
			if err != nil {
				log.Printf("Error creating task: %v\n", err)
			} else {
				fmt.Println("Task created successfully!")
			}
		case "2":
			tasks, err := readTasks(db)
			if err != nil {
				log.Printf("Error reading tasks: %v\n", err)
			} else {
				fmt.Println("Tasks:")
				for _, task := range tasks {
					fmt.Printf("ID: %d, Name: %s, Completed: %t\n", task.ID, task.Name, task.Completed)
				}
			}
		case "3":
			fmt.Println("Enter task ID to mark as completed:")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("Invalid input for ID")
				continue
			}
			err = updateTask(db, id)
			if err != nil {
				log.Printf("Error updating task: %v\n", err)
			} else {
				fmt.Println("Task updated successfully!")
			}
		case "4":
			fmt.Println("Enter task ID to delete:")
			idStr, _ := reader.ReadString('\n')
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				log.Println("Invalid input for ID")
				continue
			}
			err = deleteTask(db, id)
			if err != nil {
				log.Printf("Error deleting task: %v\n", err)
			} else {
				fmt.Println("Task deleted successfully!")
			}
		case "5":
			fmt.Println("Exiting...")
			break FOOR_LOOP_FOR_DB
		default:
			fmt.Println("Invalid option, please choose 1-5.")
		}
	}
}

func connectDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected!")
	return db
}

func createTask(db *sql.DB, name string) error {
	_, err := db.Exec("INSERT INTO tasks (name) VALUES ($1) ON CONFLICT DO NOTHING", name)
	return err
}

func readTasks(db *sql.DB) ([]Task, error) {
	rows, err := db.Query("SELECT id, name, completed FROM tasks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.Name, &t.Completed)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func updateTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE tasks SET completed = TRUE WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func deleteTask(db *sql.DB, id int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
