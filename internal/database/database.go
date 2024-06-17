package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

var ErrForeignKey error = errors.New("failed to insert. A foreign key does not reference a valid row in the parent table")
var ErrAlreadyExists error = errors.New("entity already exists")
var ErrNoResult error = errors.New("no results")

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = &service{
		db: db,
	}
	dbInstance.createAllTables()
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := s.db.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	return s.db.Close()
}

func (s *service) createAllTables() {
	queryStr := `CREATE TABLE IF NOT EXISTS "User" (
	"id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"username" text NOT NULL UNIQUE,
	"password" text NOT NULL,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "Categories" (
	"id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"belongs_to" bigint NOT NULL,
	"name" text,
	"order" bigint NOT NULL,
	PRIMARY KEY ("id"),
	FOREIGN KEY ("belongs_to") REFERENCES "User"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "Task" (
	"id" bigint GENERATED ALWAYS AS IDENTITY NOT NULL UNIQUE,
	"title" text,
	"details" text,
	"state" bigint NOT NULL DEFAULT 0,
	"due" text,
	PRIMARY KEY ("id")
);

CREATE TABLE IF NOT EXISTS "CategoryTasks" (
    "category_id" bigint NOT NULL,
    "task_id" bigint NOT NULL,
    "order" bigint NOT NULL,
    PRIMARY KEY ("category_id", "task_id"),
    FOREIGN KEY ("category_id") REFERENCES "Categories"("id") ON DELETE CASCADE,
    FOREIGN KEY ("task_id") REFERENCES "Task"("id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "UserSharesWith" (
	"Sharing" bigint NOT NULL,
	"Receiving" bigint NOT NULL,
	PRIMARY KEY ("Sharing", "Receiving"),
	FOREIGN KEY ("Sharing") REFERENCES "User"("id") ON DELETE CASCADE,
	FOREIGN KEY ("Receiving") REFERENCES "User"("id") ON DELETE CASCADE
);`

	_, err := s.db.Exec(queryStr)
	if err != nil {
		log.Fatal(err)
	}
}

// Returns an empty User instance and ErrNoResult if the user was not found
func GetUserByUsername(username string) (User, error) {
	querystr := `SELECT u.id, u.username, u.password FROM "User" u WHERE "username" = $1`
	var user User
	err := dbInstance.db.QueryRow(querystr, username).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Println("No rows found with username " + username)
			return User{}, ErrNoResult
		}
		log.Fatal(err)
	}

	return user, nil
}

// Returns an empty User instance and ErrNoResult if the user was not found
// Not used
func GetUserByID(userid int64) (User, error) {
	querystr := `SELECT u.id, u.username, u.password FROM "User" u WHERE "id" = $1`
	var user User
	err := dbInstance.db.QueryRow(querystr, userid).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No rows found with userid %d", userid)
			return User{}, ErrNoResult
		}
		log.Fatal(err)
	}

	return user, nil
}

// Returns the username associated with a category id. This is to check whether the user who made the request has the permission to manipulate the task with the specified id
// Modify later for sharing feature
func GetUsernameByCategoryId(category_id int64) (string, error) {
	query := `SELECT u.username FROM "User" u JOIN "Categories" c ON u.id = c.belongs_to WHERE c.id = $1`
	var username string
	err := dbInstance.db.QueryRow(query, category_id).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoResult
		}
		log.Fatalf("Error getting username by Category ID: %v\n", err)
	}
	return username, nil
}

// Returns the username associated with a task id. This is to check whether the user who made the request has the permission to manipulate the task with the specified id
func GetUsernameByTaskId(task_id int64) (string, error) {
	query := `SELECT u.username FROM "User" u JOIN "Categories" c ON u.id = c.belongs_to JOIN "CategoryTasks" a ON c.id = a.category_id WHERE a.task_id = $1`
	var username string
	err := dbInstance.db.QueryRow(query, task_id).Scan(&username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNoResult
		}
		log.Fatalf("Error getting username by Task ID: %v\n", err)
	}
	return username, nil
}

// Returns an empty Categories instance and sql.ErrNoRows if the category was not found
func GetCategoryByID(categoryId int64) (Categories, error) {
	querystr := `SELECT c.id, c.belongs_to, c.name, c.order FROM "Categories" c WHERE "id" = $1`
	var category Categories
	err := dbInstance.db.QueryRow(querystr, categoryId).Scan(&category.Id, &category.Belongs_to, &category.Name, &category.Order)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("No rows found with category ID %d\n", categoryId)
			return Categories{}, err
		}
		log.Fatal(err)
	}

	return category, nil
}

// Returns a slice of Categories belonging to a particular user. Returns an empty slice if the user does not have any categories or if the user does not exist
func GetCategoriesByUsername(username string) []Categories {
	query := `SELECT c.id, c.belongs_to, c.name, c.order FROM "User" u JOIN "Categories" c ON u.id = c.belongs_to WHERE u.username = $1`
	var categories []Categories
	rows, err := dbInstance.db.Query(query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return categories
		}
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var category Categories
		err := rows.Scan(&category.Id, &category.Belongs_to, &category.Name, &category.Order)
		if err != nil {
			log.Fatal(err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return categories
}

// Returns a slice of Tasks belonging to a particular user. Returns an empty slice if the user does not have any tasks or if the user does not exist, nil on error
func GetTasksByUsername(username string) []Task {
	query := `
	SELECT 
		t.id,
		t.title,
		t.details,
		t.state,
		t.due,
		a.order,
		a.category_id

	FROM 
		"User" u JOIN "Categories" c ON c.belongs_to = u.id JOIN "CategoryTasks" a ON c.id = a.category_id JOIN "Task" t ON a.task_id = t.id
	WHERE 
		u.username = $1;
	`
	var tasks []Task
	rows, err := dbInstance.db.Query(query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return tasks
		}
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Title, &task.Details, &task.State, &task.Due, &task.Order, &task.Belongs_to)
		if err != nil {
			log.Fatal(err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return tasks
}

// not used
func GetTasksByUser(user User) ([]Task, error) {
	query := `
	SELECT 
		t.id,
		t.title,
		t.details,
		t.state,
		t.due
	FROM 
		"User" u
	JOIN 
		"Categories" c ON u.id = c.belongs_to
	JOIN 
		"Task" t ON c.id = t.belongs_to
	WHERE 
		u.username = $1;
	`

	rows, err := dbInstance.db.Query(query, user.Username)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.Id, &task.Title, &task.Details, &task.State, &task.Due)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %v", err)
	}

	return tasks, nil
}

func UpdateTask(task Task) Task {
	query := `UPDATE "Task" SET title = $1, details = $2, state = $3, due = $4 WHERE id = $5`
	result, err := dbInstance.db.Exec(query, task.Title, task.Details, task.State, task.Due, task.Id)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}
	return task
}

func DeleteTask(task Task) {

	tx, err := dbInstance.db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	query0 := `SELECT "order", category_id FROM "CategoryTasks" WHERE task_id = $1`
	query := `DELETE FROM "Task" WHERE id = $1`
	query2 := `UPDATE "CategoryTasks" SET "order" = "order" - 1 WHERE "order" >= $1 AND category_id = $2`
	var old_belongs_to int64
	var old_order int64
	err = tx.QueryRow(query0, task.Id).Scan(&old_order, &old_belongs_to)
	if err != nil {
		log.Fatal(err)
	}
	result, err := tx.Exec(query, task.Id)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		tx.Rollback()
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}
	_, err = tx.Exec(query2, old_order, old_belongs_to)
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

// Hashes the given password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Adds a new user to the "User" table. Will return ErrAlreadyExists if the user already exists in the database
func AddUser(user User) error {

	//checks if the user already exists
	_, alrExErr := GetUserByUsername(user.Username)
	if alrExErr == nil {
		return ErrAlreadyExists
	}

	// Hash the user's password
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	query := `
	INSERT INTO "User" (username, password)
	VALUES ($1, $2)
	RETURNING id
	`
	var userID int64
	err = dbInstance.db.QueryRow(query, user.Username, hashedPassword).Scan(&userID)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	fmt.Printf("User added with ID: %d\n", userID)
	return nil
}

func AddCategory(category Categories) int64 {

	tx, err := dbInstance.db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	query := `INSERT INTO "Categories" ("belongs_to", "name", "order") VALUES ($1, $2, $3) RETURNING id`
	query2 := `UPDATE "Categories" SET "order" = "order" + 1 WHERE "order" >= $1 AND NOT "id" = $2`
	var categoryID int64
	err = tx.QueryRow(query, category.Belongs_to, category.Name, category.Order).Scan(&categoryID)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	_, err = tx.Exec(query2, category.Order, categoryID)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return categoryID
}

func UpdateCategory(category Categories) Categories {

	query := `UPDATE "Categories" SET name = $1 WHERE id = $2`
	result, err := dbInstance.db.Exec(query, category.Name, category.Id)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}
	return category
}

func DeleteCategory(category Categories) {

	tx, err := dbInstance.db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	query := `UPDATE "Categories" SET "order" = "order" - 1 WHERE "order" > (SELECT "order" FROM "Categories" WHERE "id" = $1)`
	query2 := `DELETE FROM "Categories" WHERE id = $1`

	_, err = tx.Exec(query, category.Id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	result, err := tx.Exec(query2, category.Id)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}
	if rows != 1 {
		tx.Rollback()
		log.Fatalf("expected to affect 1 row, affected %d", rows)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}

func ChangeCategoryOrder(category_id int64, to int64, belongs_to int64) {
	var oldCategoryOrder int64
	tx, err := dbInstance.db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	query := `SELECT "order" FROM "Categories" WHERE id = $1`
	err = tx.QueryRow(query, category_id).Scan(&oldCategoryOrder)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	if oldCategoryOrder == to {
		return
	}
	if oldCategoryOrder > to {
		query = `UPDATE "Categories" SET "order" = "order" + 1 WHERE "order" >= $1 AND "order" < $2 AND belongs_to = $3`
		_, err = tx.Exec(query, to, oldCategoryOrder, belongs_to)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		query = `UPDATE "Categories" SET "order" = $1 WHERE id = $2`
		_, err = tx.Exec(query, to, category_id)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		query = `UPDATE "Categories" SET "order" = "order" - 1 WHERE "order" > $1 AND "order" <= $2`
		_, err = tx.Exec(query, oldCategoryOrder, to)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		query = `UPDATE "Categories" SET "order" = $1 WHERE id = $2`
		_, err = tx.Exec(query, to, category_id)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func AddTask(task Task) Task {

	tx, err := dbInstance.db.BeginTx(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	query := `
	WITH rows AS (
        INSERT INTO "Task" (title, details, state, due)
        VALUES ($1, $2, $3, $4)
        RETURNING id
        )
        INSERT INTO "CategoryTasks" ("category_id", "order", "task_id")
        SELECT $5, $6, id FROM rows
		RETURNING task_id;
	`
	query2 := `UPDATE "CategoryTasks" SET "order" = "order" + 1 WHERE "order" >= $1 AND category_id = $2 AND NOT task_id = $3`

	err = tx.QueryRow(query, task.Title, task.Details, task.State, task.Due, task.Belongs_to, task.Order).Scan(&task.Id)
	if err != nil {
		tx.Rollback()
		log.Println("query1")
		log.Fatal(err)
	}
	_, err = tx.Exec(query2, task.Order, task.Belongs_to, task.Id)
	if err != nil {
		log.Println("query2")
		tx.Rollback()
		log.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task %s added with ID: %d\n", task.Title, task.Id)
	return task
}
