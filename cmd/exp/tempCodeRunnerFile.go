	userID := 1
	for i := 1; i <= 5; i++ {
		amount := i * 100
		desc := fmt.Sprintf("Fake order #%d", i)

		_, err := db.Exec(`
			INSERT INTO orders(user_id, amount, description)
			VALUES ($1, $2, $3)`, userID, amount, desc)

		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Created fake orders.")