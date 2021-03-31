package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

// Database instance
var db *sql.DB

// settingan untuk Database 
const (
	host     = "localhost"
	port     = 5432 // Default port
	user     = "postgres"
	password = "ubuntu"
	dbname   = "testing"
)

// Employee struct
type Employee struct {
	ID     int `json:"id"`
	Name   string `json:"name"`
	Salary int `json:"salary"`
	Age    int `json:"age"`
}

// Employees struct
type Employees struct {
	Employees []Employee `json:"employees"`
}

// function untuk konek ke database
func Connect() error {
	var err error
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname))
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	return nil
}

func main() {
	// konek ke database dengan memamnggil function Connect()
	if err := Connect(); err != nil {
		log.Fatal(err)
	}

	// membuat fiber baru
	app := fiber.New()

	app.Get("/employee", func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, name, salary, age FROM employees order by id")
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		defer rows.Close()
		result := Employees{}

		for rows.Next() {
			employee := Employee{}
			if err := rows.Scan(&employee.ID, &employee.Name, &employee.Salary, &employee.Age); err != nil {
				return err 
			}

			result.Employees = append(result.Employees, employee)
		}
		return c.JSON(result)
	})

	app.Post("/employee", func(c *fiber.Ctx) error {
		u := new(Employee)

		if err := c.BodyParser(u); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		res, err := db.Query("INSERT INTO employees (name, salary, age)VALUES ($1, $2, $3)", u.Name, u.Salary, u.Age)
		if err != nil {
			return err
		}

		log.Println(res)

		return c.JSON(u)
	})

	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		u := new(Employee)
		id := c.Params("id")

		if err := c.BodyParser(u); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		res, err := db.Query("UPDATE employees SET name=$1,salary=$2,age=$3 WHERE id=$5", u.Name, u.Salary, u.Age, u.ID)
		if err != nil {
			return err
		}

		log.Println(res)

		return c.Status(201).JSON(u)
	})

	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		u := new(Employee)
		id := c.Params("id")

		if err := c.BodyParser(u); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		res, err := db.Query("DELETE FROM employees WHERE id = $1", u.ID)
		if err != nil {
			return err
		}

		log.Println(res)

		return c.JSON("Deleted")
	})

	log.Fatal(app.Listen(":3000"))
}