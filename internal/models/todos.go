package models

import (
	"database/sql"
	"errors"
	"time"
)

// define a todo type
type Todo struct {
	ID      string
	Body    string `json:"body"`
	Created time.Time
}

// define a todo model type which wraps a sql.DB connection pool
type TodoModel struct {
	DB *sql.DB
}

// insert a new todo into the database
func (m *TodoModel) Insert(newId string, body string) (string, error) {
	// use placeholder parameters instead of interpolating data in the SQL query
	// as this is untrusted user input from a form
	stmt := `INSERT INTO todos (id, body, created) 	VALUES(
	?, ?, UTC_TIMESTAMP())`

	_, err := m.DB.Exec(stmt, newId, body)
	if err != nil {
		return "", err
	}

	// use the LastInserId() method on the result to get the ID of
	// the newly created record in the snippets table
	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }

	// the id returned has the type int64, convert to an int type
	return newId, nil
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

// return the 10 most recently created todos
func (m *TodoModel) All() ([]*Todo, error) {
	// SQL statement we want to execute
	stmt := `SELECT * FROM todos
	ORDER BY id`

	// Use the Query() method on the connection pool to execute the stmt
	// this returns a sql.Rows resultset containing the result of our query
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is always
	// properly closed before the Latest() method returns
	// Defer stmt should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic trying to close a nil resultset
	defer rows.Close()

	// Initialize an empty slice to hold the Snippet structs
	todos := []*Todo{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by
	// the rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		t := &Todo{}
		// Use rows.Scan() to copy the values from each field in the row to
		// the new Snippet object that we created. Again, the arguments to
		// row.Scan() must be pointers to the place you want to copy the data into, and
		// the number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&t.ID, &t.Body, &t.Created)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		todos = append(todos, t)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve
	// any error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return todos, nil
}
