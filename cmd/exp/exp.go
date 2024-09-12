package main

import (
	"fmt"

	"github.com/SuWh1/WebDevGo/models"
)

func main() {
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close() // make sure that database will close

	err = db.Ping() // that make request, so we have to up the docker first (make sure it is up and running)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected!")

	us := models.UserService{
		DB: db,
	}
	user, err := us.Create("alNUR@gmail.com", "alnurloh")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// // create a table
	// _, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id SERIAL PRIMARY KEY,
	// 		name TEXT,
	// 		email TEXT UNIQUE NOT NULL
	// 	);

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 		id SERIAL PRIMARY KEY,
	// 		user_id INT NOT NULL,
	// 		amount INT,
	// 		description TEXT
	// 	);
	// `)

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Table created!")

	// name := "Dauren5"
	// email := "dauren5@gmail.com"

	// row := db.QueryRow(`
	// 	INSERT INTO users(name, email)
	// 	VALUES ($1, $2) RETURNING id;`, name, email) // queryRow output one value

	// row.Err() // checks if error occur in QueryRow execution || if there is error it will return it
	// var id int
	// err = row.Scan(&id)
	// if err != nil {
	// // 	panic(err)
	// // }
	// // fmt.Println("User created. id = ", id)

	// id := 6

	// row := db.QueryRow(`
	// 	SELECT name, email
	// 	FROM users
	// 	WHERE id=$1;`, id)

	// var name, email string
	// err = row.Scan(&name, &email)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("User information: name=%s, email=%s\n", name, email)

	// userID := 1
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i)

	// 	_, err := db.Exec(`
	// 		INSERT INTO orders(user_id, amount, description)
	// 		VALUES ($1, $2, $3)`, userID, amount, desc)

	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fmt.Println("Created fake orders.")

	// type Order struct {
	// 	ID          int
	// 	UserID      int
	// 	Amount      int
	// 	Description string
	// }

	// var orders []Order
	// userID := 1
	// // now we using just query which is slight different, it take many values and do not return error by itself
	// rows, err := db.Query(`
	// 	SELECT id, amount, description
	// 	FROM orders
	// 	WHERE user_id=$1`, userID)

	// if err != nil {
	// 	panic(err)
	// }
	// defer rows.Close()

	// for rows.Next() { // rows has nothing at the start, only next is value
	// 	var order Order
	// 	order.UserID = userID
	// 	err := rows.Scan(&order.ID, &order.Amount, &order.Description)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	orders = append(orders, order)
	// }

	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Orders:", orders)
}
