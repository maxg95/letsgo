package models

import (
	"database/sql"
	"errors"
	"time"
)

type OrderModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (Order, error)
	Latest() ([]Order, error)
}

type Order struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type OrderModel struct {
	DB *sql.DB
}

func (m *OrderModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO orders (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *OrderModel) Get(id int) (Order, error) {
	stmt := `SELECT id, title, content, created, expires FROM orders 
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	var o Order

	err := row.Scan(&o.ID, &o.Title, &o.Content, &o.Created, &o.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Order{}, ErrNoRecord
		} else {
			return Order{}, err
		}
	}

	return o, nil
}

func (m *OrderModel) Latest() ([]Order, error) {
	stmt := `SELECT id, title, content, created, expires FROM orders
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var orders []Order

	for rows.Next() {
		var o Order
		err = rows.Scan(&o.ID, &o.Title, &o.Content, &o.Created, &o.Expires)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}
