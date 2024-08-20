package models

import (
	"database/sql"
	"errors"
	"time"
)

// define a todo type
type Todo struct {
	ID      int
	Body    string
	Created time.Time
}

// define a todo model type which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}

// insert a new todo into the database
func (m *TodoModel) Insert(body string) (int, error) {
	// use placeholder parameters instead of interpolating data in the SQL query
	// as this is untrusted user input from a form
	stmt := `INSERT INTO todos (body, created) 	VALUES(
	?, UTC_TIMESTAMP())`

	result, err := m.DB.Exec(stmt, body)
	if err != nil {
		return 0, err
	}

	// use the LastInserId() method on the result to get the ID of
	// the newly created record in the snippets table
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// the id returned has the type int64, convert to an int type
	return int(id), nil
}

// return a specific snippet based on its id
func (m *TodoModel) Get(id int) (*Todo, error) {
	// Write the SQL statement we want to execute.
	stmt := `SELECT id, body, created FROM todos
	WHERE id = ?`

	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for
	// the placeholder parameter. This returns a pointer to a sql.Row object
	// which holds the result from the database.
	row := m.DB.QueryRow(stmt, id)

	// Initialize a pointer to a new zeroed Snippet struct.
	t := &Todo{}

	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data
	// into, and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&t.ID, &t.Body, &t.Created)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for
		// that error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	// If everything went OK then return the Snippet object.
	return t, nil
}
