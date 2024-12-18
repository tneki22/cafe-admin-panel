package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	// Настройки подключения к базе данных
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=cafe_admin_users sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Количество заказов для вставки
	numOrders := 100

	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	statuses := []string{"Выполнен", "Отменен", "В работе"}

	for i := 0; i < numOrders; i++ {
		// Генерируем случайную сумму заказа от 0 до 2000
		total := rand.Float64() * 2000

		// Выбираем случайный статус заказа
		status := statuses[rand.Intn(len(statuses))]

		// Генерируем случайную дату в пределах последнего месяца
		daysAgo := rand.Intn(30) // От 0 до 29 дней назад
		hoursAgo := rand.Intn(24)
		minutesAgo := rand.Intn(60)
		secondsAgo := rand.Intn(60)
		createdAt := time.Now().AddDate(0, 0, -daysAgo).Add(
			-time.Duration(hoursAgo)*time.Hour -
				time.Duration(minutesAgo)*time.Minute -
				time.Duration(secondsAgo)*time.Second,
		)

		// Выполняем вставку в таблицу orders
		_, err := db.Exec(
			"INSERT INTO orders (total, status, created_at) VALUES ($1, $2, $3)",
			total, status, createdAt,
		)
		if err != nil {
			log.Printf("Ошибка вставки заказа: %v", err)
			continue
		}

		fmt.Printf("Добавлен заказ #%d: сумма=%.2f, статус=%s, дата=%s\n",
			i+1, total, status, createdAt.Format("2006-01-02 15:04:05"))
	}

	fmt.Printf("Добавлено %d заказов в базу данных.\n", numOrders)
}
