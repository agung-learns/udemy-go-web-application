package models

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}

type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type Order struct {
	ID            int       `json:"id"`
	WidgetID      int       `json:"widget_id"`
	TransactionID int       `json:"transaction_id"`
	CustomerID    int       `json:"customer_id"`
	StatusID      int       `json:"status_id"`
	Quantity      int       `json:"quantity"`
	Amount        int       `json:"amount"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var widget Widget
	row := m.DB.QueryRowContext(ctx, ""+
		`SELECT 
    		id, name, description, inventory_level, price, coalesce(image, ''), 
    		is_recurring, plan_id, created_at, updated_at 
		FROM widgets 
		WHERE id = $1`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		return widget, err
	}
	return widget, nil
}

func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO transactions 
		    (amount, currency, last_four, bank_return_code, transaction_status_id, expiry_month, expiry_year, payment_intent, payment_method, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt, txn.Amount, txn.Currency, txn.LastFour, txn.BankReturnCode, txn.TransactionStatusID, txn.ExpiryMonth, txn.ExpiryYear, txn.PaymentIntent, txn.PaymentMethod, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO orders 
		    (widget_id, transaction_id, status_id, customer_id, quantity, amount, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt, order.WidgetID, order.TransactionID, order.StatusID, order.CustomerID, order.Quantity, order.Amount, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `
		INSERT INTO customers 
		    (first_name, last_name, email, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt, c.FirstName, c.LastName, c.Email, c.CreatedAt, c.UpdatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qs := `SELECT id, first_name, last_name, email, password, created_at, updated_at FROM users WHERE email = $1`
	row := m.DB.QueryRowContext(ctx, qs, strings.ToLower(email))

	var user User
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	err := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", email).Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return 0, errors.New("wrong password")
		default:
			return 0, err
		}
	}

	return id, nil
}
