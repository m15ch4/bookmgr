package database

import (
	"database/sql"
	"fmt"
	"log"

	"bookmgr/config"
	"bookmgr/models"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	conn *sql.DB
}

func New(cfg *config.Config) (*Database, error) {
	// First, connect without database to check/create it
	dsnWithoutDB := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
		cfg.DatabaseUser,
		cfg.DatabasePassword,
		cfg.DatabaseHost,
		cfg.DatabasePort,
	)

	db, err := sql.Open("mysql", dsnWithoutDB)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL server: %w", err)
	}

	if !cfg.SkipBootstrap {
		log.Println("Bootstrapping database...")
		if err := bootstrap(db, cfg.DatabaseName); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to bootstrap database: %w", err)
		}
		log.Println("Database bootstrapped successfully")
	}
	db.Close()

	// Connect to the specific database
	conn, err := sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database successfully")

	return &Database{conn: conn}, nil
}

func bootstrap(db *sql.DB, dbName string) error {
	// Create database if it doesn't exist
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	// Use the database
	_, err = db.Exec(fmt.Sprintf("USE %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to use database: %w", err)
	}

	// Create books table
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		num_pages INT NOT NULL DEFAULT 0,
		author VARCHAR(255) NOT NULL,
		rating DECIMAL(3,2) NOT NULL DEFAULT 0.00,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		INDEX idx_title (title),
		INDEX idx_author (author)
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
	`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	return nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

// Create inserts a new book
func (db *Database) Create(book *models.Book) (*models.Book, error) {
	query := `INSERT INTO books (title, num_pages, author, rating) VALUES (?, ?, ?, ?)`

	result, err := db.conn.Exec(query, book.Title, book.NumPages, book.Author, book.Rating)
	if err != nil {
		return nil, fmt.Errorf("failed to create book: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID: %w", err)
	}

	book.ID = int(id)
	return book, nil
}

// GetAll retrieves all books
func (db *Database) GetAll() ([]*models.Book, error) {
	query := `SELECT id, title, num_pages, author, rating FROM books ORDER BY id`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get books: %w", err)
	}
	defer rows.Close()

	var books []*models.Book
	for rows.Next() {
		book := &models.Book{}
		err := rows.Scan(&book.ID, &book.Title, &book.NumPages, &book.Author, &book.Rating)
		if err != nil {
			return nil, fmt.Errorf("failed to scan book: %w", err)
		}
		books = append(books, book)
	}

	return books, nil
}

// GetByID retrieves a book by ID
func (db *Database) GetByID(id int) (*models.Book, error) {
	query := `SELECT id, title, num_pages, author, rating FROM books WHERE id = ?`

	book := &models.Book{}
	err := db.conn.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.NumPages, &book.Author, &book.Rating)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get book: %w", err)
	}

	return book, nil
}

// Update modifies an existing book
func (db *Database) Update(id int, book *models.Book) (*models.Book, error) {
	query := `UPDATE books SET title = ?, num_pages = ?, author = ?, rating = ? WHERE id = ?`

	result, err := db.conn.Exec(query, book.Title, book.NumPages, book.Author, book.Rating, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, nil
	}

	book.ID = id
	return book, nil
}

// Delete removes a book
func (db *Database) Delete(id int) (bool, error) {
	query := `DELETE FROM books WHERE id = ?`

	result, err := db.conn.Exec(query, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
