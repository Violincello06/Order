package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"order/internal/models"
)

type MySQLOrderRepository struct {
	db *sql.DB
}

func NewMySQLOrderRepository(db *sql.DB) *MySQLOrderRepository {
	return &MySQLOrderRepository{db: db}
}

// NewOrderRepositoryAutoDetect mencoba koneksi ke MySQL XAMPP, jika offline otomatis fallback ke Memory
func NewOrderRepositoryAutoDetect(user, pass, host, port, dbname string) (OrderRepository, string) {
	db, err := InitMySQLConnection(user, pass, host, port, dbname)
	if err == nil {
		return NewMySQLOrderRepository(db), "MySQL Database (order_db)"
	}
	return NewMemoryOrderRepository(), fmt.Sprintf("In-Memory Mock (MySQL offline: %v)", err)
}

// InitMySQLConnection membuka koneksi ke MySQL XAMPP
func InitMySQLConnection(user, pass, host, port, dbname string) (*sql.DB, error) {
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "3306"
	}
	if user == "" {
		user = "root"
	}
	if dbname == "" {
		dbname = "order_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local", user, pass, host, port, dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (r *MySQLOrderRepository) Save(ctx context.Context, order *models.Order) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Insert ke tabel orders
	orderQuery := `
		INSERT INTO orders (
			id, customer_name, customer_email, total_amount, status,
			payment_id, payment_url, payment_status, payment_method, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, orderQuery,
		order.ID,
		order.CustomerName,
		order.CustomerEmail,
		order.TotalAmount,
		string(order.Status),
		order.PaymentID,
		order.PaymentURL,
		order.PaymentStatus,
		order.PaymentMethod,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("gagal insert order ke mysql: %w", err)
	}

	// 2. Insert ke tabel order_items
	itemQuery := `
		INSERT INTO order_items (
			order_id, product_id, product_name, price, quantity, subtotal, image_url, created_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	for _, item := range order.Items {
		_, err := tx.ExecContext(ctx, itemQuery,
			order.ID,
			item.ProductID,
			item.ProductName,
			item.Price,
			item.Quantity,
			item.Subtotal,
			item.ImageURL,
			order.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("gagal insert item order ke mysql: %w", err)
		}
	}

	return tx.Commit()
}

func (r *MySQLOrderRepository) FindByID(ctx context.Context, id string) (*models.Order, error) {
	orderQuery := `
		SELECT id, customer_name, customer_email, total_amount, status,
		       COALESCE(payment_id, ''), COALESCE(payment_url, ''),
		       COALESCE(payment_status, ''), COALESCE(payment_method, ''),
		       created_at, updated_at
		FROM orders
		WHERE id = ?
	`
	row := r.db.QueryRowContext(ctx, orderQuery, id)

	var o models.Order
	var statusStr string
	err := row.Scan(
		&o.ID,
		&o.CustomerName,
		&o.CustomerEmail,
		&o.TotalAmount,
		&statusStr,
		&o.PaymentID,
		&o.PaymentURL,
		&o.PaymentStatus,
		&o.PaymentMethod,
		&o.CreatedAt,
		&o.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	o.Status = models.OrderStatus(statusStr)

	// Ambil items
	itemsQuery := `
		SELECT product_id, product_name, price, quantity, subtotal, COALESCE(image_url, '')
		FROM order_items
		WHERE order_id = ?
	`
	rows, err := r.db.QueryContext(ctx, itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.OrderItem
		if err := rows.Scan(
			&item.ProductID,
			&item.ProductName,
			&item.Price,
			&item.Quantity,
			&item.Subtotal,
			&item.ImageURL,
		); err != nil {
			return nil, err
		}
		o.Items = append(o.Items, item)
	}

	return &o, nil
}

func (r *MySQLOrderRepository) FindAll(ctx context.Context) ([]*models.Order, error) {
	ordersQuery := `
		SELECT id, customer_name, customer_email, total_amount, status,
		       COALESCE(payment_id, ''), COALESCE(payment_url, ''),
		       COALESCE(payment_status, ''), COALESCE(payment_method, ''),
		       created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, ordersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	orderMap := make(map[string]*models.Order)

	for rows.Next() {
		var o models.Order
		var statusStr string
		if err := rows.Scan(
			&o.ID,
			&o.CustomerName,
			&o.CustomerEmail,
			&o.TotalAmount,
			&statusStr,
			&o.PaymentID,
			&o.PaymentURL,
			&o.PaymentStatus,
			&o.PaymentMethod,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		o.Status = models.OrderStatus(statusStr)
		o.Items = make([]models.OrderItem, 0)
		orders = append(orders, &o)
		orderMap[o.ID] = &o
	}

	if len(orders) == 0 {
		return orders, nil
	}

	// Ambil semua item belanja
	itemsQuery := `
		SELECT order_id, product_id, product_name, price, quantity, subtotal, COALESCE(image_url, '')
		FROM order_items
	`
	itemRows, err := r.db.QueryContext(ctx, itemsQuery)
	if err == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var orderID string
			var item models.OrderItem
			if err := itemRows.Scan(
				&orderID,
				&item.ProductID,
				&item.ProductName,
				&item.Price,
				&item.Quantity,
				&item.Subtotal,
				&item.ImageURL,
			); err == nil {
				if parentOrder, exists := orderMap[orderID]; exists {
					parentOrder.Items = append(parentOrder.Items, item)
				}
			}
		}
	}

	return orders, nil
}

func (r *MySQLOrderRepository) Update(ctx context.Context, order *models.Order) error {
	updateQuery := `
		UPDATE orders
		SET status = ?,
		    payment_id = ?,
		    payment_url = ?,
		    payment_status = ?,
		    payment_method = ?,
		    updated_at = ?
		WHERE id = ?
	`
	result, err := r.db.ExecContext(ctx, updateQuery,
		string(order.Status),
		order.PaymentID,
		order.PaymentURL,
		order.PaymentStatus,
		order.PaymentMethod,
		order.UpdatedAt,
		order.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrOrderNotFound
	}

	return nil
}
