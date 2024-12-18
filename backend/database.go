package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type Provider struct {
	conn *sql.DB
}

func NewProvider(host string, port int, user, password, dbName string) *Provider {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable timezone=Europe/Moscow",
		host, port, user, password, dbName)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка соединения
	err = conn.Ping()
	if err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}

	return &Provider{conn: conn}
}

func (p *Provider) CreateUser(username, email, hashedPassword string) error {
	_, err := p.conn.Exec("INSERT INTO users (name, email, password) VALUES ($1, $2, $3)", username, email, hashedPassword)
	if err != nil {
		log.Printf("Error creating user in database: %v", err)
	}
	return err
}

func (p *Provider) CheckUserByEmail(email string) (bool, error) {
	err := p.conn.QueryRow("SELECT (email) FROM users WHERE email = $1", email).Scan(&email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (p *Provider) GetUsernameAndHashedPassword(email string) (string, string, error) {
	var password_db string
	var name string
	err := p.conn.QueryRow("SELECT name, password FROM users WHERE email = $1", email).Scan(&name, &password_db)
	if err != nil {
		return "", "", err
	}

	return name, password_db, nil
}
func (p *Provider) FetchMenuItems() ([]MenuItem, error) {
	rows, err := p.conn.Query("SELECT id, name, description, price, created_at FROM menu ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var menuItems []MenuItem
	for rows.Next() {
		var item MenuItem
		var createdAt time.Time
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &createdAt)
		if err != nil {
			return nil, err
		}
		item.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		menuItems = append(menuItems, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return menuItems, nil
}
func (p *Provider) DeleteMenuItem(id int) error {
	_, err := p.conn.Exec("DELETE FROM menu WHERE id = $1", id)
	return err
}

func (p *Provider) UpdateMenuItem(item MenuItem) (MenuItem, error) {
	_, err := p.conn.Exec("UPDATE menu SET name = $1, description = $2, price = $3 WHERE id = $4",
		item.Name, item.Description, item.Price, item.ID)
	if err != nil {
		return MenuItem{}, err
	}

	// Возвращаем обновленный элемент
	row := p.conn.QueryRow("SELECT id, name, description, price, created_at FROM menu WHERE id = $1", item.ID)
	var updatedItem MenuItem
	var createdAt time.Time
	err = row.Scan(&updatedItem.ID, &updatedItem.Name, &updatedItem.Description, &updatedItem.Price, &createdAt)
	if err != nil {
		return MenuItem{}, err
	}
	updatedItem.CreatedAt = createdAt.Format("2006-01-02 15:04:05")

	return updatedItem, nil
}
func (p *Provider) AddMenuItem(item MenuItem) (MenuItem, error) {
	var newItem MenuItem
	var createdAt time.Time
	err := p.conn.QueryRow(
		"INSERT INTO menu (name, description, price) VALUES ($1, $2, $3) RETURNING id, name, description, price, created_at",
		item.Name, item.Description, item.Price,
	).Scan(&newItem.ID, &newItem.Name, &newItem.Description, &newItem.Price, &createdAt)
	if err != nil {
		return MenuItem{}, err
	}
	newItem.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
	return newItem, nil
}
func (p *Provider) AddOrder(items []OrderItem) (Order, error) {
	tx, err := p.conn.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return Order{}, err
	}

	var newOrder Order

	// Ensure the transaction is rollbacked in case of error
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Transaction rollback failed: %v", rbErr)
			} else {
				log.Println("Transaction rolled back")
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				log.Printf("Failed to commit transaction: %v", commitErr)
				err = commitErr
			} else {
				log.Println("Transaction committed")
			}
		}
	}()

	// Insert new order with default status "In work" and total 0
	err = tx.QueryRow(
		"INSERT INTO orders (total, status) VALUES ($1, $2) RETURNING id, total, status, created_at",
		0, "В работе",
	).Scan(&newOrder.ID, &newOrder.Total, &newOrder.Status, &newOrder.CreatedAt)
	if err != nil {
		log.Printf("Failed to insert new order: %v", err)
		return Order{}, fmt.Errorf("failed to insert new order: %v", err)
	}
	log.Printf("Inserted new order with ID: %d", newOrder.ID)

	var total float64
	for _, item := range items {
		var price float64
		// Get price of the menu item
		err = tx.QueryRow("SELECT price FROM menu WHERE id = $1", item.MenuItemId).Scan(&price)
		if err != nil {
			log.Printf("Failed to get price for menu item ID %d: %v", item.MenuItemId, err)
			return Order{}, fmt.Errorf("failed to get price for menu item ID %d: %v", item.MenuItemId, err)
		}
		log.Printf("Menu item ID %d has price: %.2f", item.MenuItemId, price)

		total += price * float64(item.Quantity)
		log.Printf("Added %.2f to total. Current total: %.2f", price*float64(item.Quantity), total)

		// Insert into order_items
		_, err = tx.Exec(
			"INSERT INTO order_items (order_id, menu_item_id, quantity) VALUES ($1, $2, $3)",
			newOrder.ID, item.MenuItemId, item.Quantity,
		)
		if err != nil {
			log.Printf("Failed to insert order item (OrderID: %d, MenuItemID: %d, Quantity: %d): %v",
				newOrder.ID, item.MenuItemId, item.Quantity, err)
			return Order{}, fmt.Errorf("failed to insert order item: %v", err)
		}
		log.Printf("Inserted order item (MenuItemID: %d, Quantity: %d)", item.MenuItemId, item.Quantity)
	}

	// Update the total in orders table
	_, err = tx.Exec("UPDATE orders SET total = $1 WHERE id = $2", total, newOrder.ID)
	if err != nil {
		log.Printf("Failed to update order total (OrderID: %d, Total: %.2f): %v", newOrder.ID, total, err)
		return Order{}, fmt.Errorf("failed to update order total: %v", err)
	}
	log.Printf("Updated order total for OrderID: %d to %.2f", newOrder.ID, total)

	newOrder.Total = total
	newOrder.Items = items

	return newOrder, nil
}
func (p *Provider) FetchOrders() ([]Order, error) {
	rows, err := p.conn.Query("SELECT id, total, status, created_at FROM orders ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var order Order
		err := rows.Scan(&order.ID, &order.Total, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
func (p *Provider) UpdateOrderStatus(orderID int, status string) error {
	_, err := p.conn.Exec("UPDATE orders SET status = $1 WHERE id = $2", status, orderID)
	return err
}

func (p *Provider) FetchRevenue(period string) ([]RevenueData, error) {
	var query string

	switch period {
	case "day":
		query = `
            SELECT to_char(created_at, 'HH24:00') AS time_unit, SUM(total) AS total
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 day'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	case "week":
		query = `
            SELECT to_char(created_at, 'YYYY-MM-DD') AS time_unit, SUM(total) AS total
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 week'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	case "month":
		query = `
            SELECT to_char(created_at, 'YYYY-MM-DD') AS time_unit, SUM(total) AS total
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 month'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	case "year":
		query = `
            SELECT to_char(created_at, 'YYYY-MM') AS time_unit, SUM(total) AS total
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 year'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	default:
		return nil, errors.New("invalid period")
	}

	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revenueData []RevenueData
	for rows.Next() {
		var rd RevenueData
		err := rows.Scan(&rd.TimeUnit, &rd.Total)
		if err != nil {
			return nil, err
		}
		revenueData = append(revenueData, rd)
	}

	return revenueData, nil
}

// Добавляем метод для получения количества заказов
func (p *Provider) FetchOrderCounts(period string) ([]OrderCountData, error) {
	var query string

	switch period {
	case "day":
		query = `
            SELECT to_char(created_at, 'HH24:00') AS time_unit, COUNT(*) AS count
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 day'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	case "week":
		query = `
            SELECT to_char(created_at, 'YYYY-MM-DD') AS time_unit, COUNT(*) AS count
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 week'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	case "month":
		query = `
            SELECT to_char(created_at, 'YYYY-MM-DD') AS time_unit, COUNT(*) AS count
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 month'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	case "year":
		query = `
            SELECT to_char(created_at, 'YYYY-MM') AS time_unit, COUNT(*) AS count
            FROM orders
            WHERE created_at >= NOW() - INTERVAL '1 year'
            GROUP BY time_unit
            ORDER BY time_unit ASC
        `
	default:
		return nil, errors.New("invalid period")
	}

	rows, err := p.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderCountData []OrderCountData
	for rows.Next() {
		var ocd OrderCountData
		err := rows.Scan(&ocd.TimeUnit, &ocd.Count)
		if err != nil {
			return nil, err
		}
		orderCountData = append(orderCountData, ocd)
	}

	return orderCountData, nil
}
