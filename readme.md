<h1>Todo List Backend</h1>
<h3>Project Overview</h3>
<p>This repository contains a backend server for a todo list application. It provides RESTful CRUD operations for the todo model and handles user authentication (signup, login, logout)</p>
<h3>Architecture</h3>
<p>The project follows a REST architecture.</p>

<h3>Features</h3>
<table>
  <tr>
    <th>User Authentication</th>
    <th>Todo Management</th>
  </tr>
  <tr>
    <td>Signup: Allows new users to create an account.</td>
    <td>Create: Adds a new todo item.</td>
  </tr>
  <tr>
    <td>Login: Authenticates users and establishes a session.</td>
    <td>Read: Retrieves todo items.</td>
  </tr>
  <tr>
    <td>Logout: Logs users out, ending their session.</td>
    <td>Update: Updates existing todo items.</td>
  </tr>
  <tr>
    <td></td>
    <td>Delete: Removes todo items.</td>
  </tr>
  <tr>
    <td></td>
    <td>Toggle Status: Marks todo items as completed or incomplete.</td>
  </tr>
</table>
<hr>

<h3>Project Structure</h3>
<code>

    .
    ├── .air.toml
    ├── .gitignore
    ├── bin
    │   └── main.exe
    ├── cmd
    │   └── web
    │       ├── context.go
    │       ├── helpers.go
    │       ├── main.go
    │       ├── middleware.go
    │       ├── routes.go
    │       ├── todo_handlers.go
    │       └── user_handlers.go
    ├── go.mod
    ├── go.sum
    ├── internal
    │   ├── models
    │   │   ├── errors.go
    │   │   ├── todos.go
    │   │   └── users.go
    │   └── validator
    │       └── validator.go
    └── tmp
        └── build-errors.log

</code>

<h3>Dependencies</h3>
<table>
<tr>
<td>github.com/go-sql-driver/mysql</td>
<td>This is a MySQL driver for Go's database/sql package. It allows the Go application to interact with a MySQL database, enabling performance of CRUD operations on users and todos tables.</td>
</tr>
<tr>
<td>github.com/joho/godotenv</td>
<td>This library is used to load environment variables from a .env file into the application. It helps manage configuration settings, such as database credentials and API keys, without hardcoding them into a codebase.</td>
</tr>
<tr>
<td>github.com/julienschmidt/httprouter</td>
<td>A lightweight, high-performance HTTP request router for Go. It helps in defining and managing routes in the application, enabling mapping specific HTTP methods and paths to handler functions. It's known for its efficiency and minimal memory overhead.</td>
</tr>
<tr>
<td>github.com/justinas/alice</td>
<td>This package simplifies chaining of HTTP middleware in your Go application. It allows to compose middleware in a clean, readable way, which is particularly useful for applying multiple middlewares (e.g., for authentication, logging, CSRF protection) to specific routes.</td>
</tr>
<tr>
<td>github.com/justinas/nosurf</td>
<td>nosurf provides Cross-Site Request Forgery (CSRF) protection for your Go web application. It helps secure forms and other POST requests by ensuring that they include a valid CSRF token, preventing unauthorized commands from being executed.</td>
</tr>
<tr>
<td>github.com/alexedwards/scs/v2</td>
<td>This package is a session management library for Go. It simplifies the handling of user sessions, providing functionality for storing session data in various backends, including MySQL (through the mysqlstore package). It's crucial for managing user authentication and maintaining session state across requests.</td>
</tr>
<tr>
<td>github.com/alexedwards/scs/mysqlstore</td>
<td>This package is an extension for scs that allows session data to be stored in a MySQL database. It's particularly useful for distributed applications where session data needs to be persisted and shared across multiple instances of your application.</td>
</tr>
<tr>
<td>github.com/google/uuid</td>
<td> A package that provides functionality for generating UUIDs (Universally Unique Identifiers). It's used for creating unique identifiers for resources, such as user IDs or todo IDs, ensuring that they are unique across the system.</td>
</tr>
<tr>
<td>golang.org/x/crypto</td>
<td> A collection of cryptographic packages for Go, providing various cryptographic algorithms and utilities. Used for secure password hashing (e.g., via bcrypt) and other cryptographic operations related to user authentication.</td>
</tr>
</table>

<h3>Environment Variables</h3>
<table>
<tr>
<td>DB_USER</td>
<td>The MySQL database username.</td>
</tr>
<tr>
<td>DB_NAME</td>
<td>The MySQL database name</td>
</tr>
<tr>
<td>DB_PASSWORD</td>
<td>The MySQL database password.</td>
</tr>
<tr>
<td>REACT_ADDRESS</td>
<td>The address of the React frontend.</td>
</tr>
</table>

<h3>Database Schema</h3>
<table>
  <tr>
    <th>Users Table</th>
    <th>Todos Table</th>
  </tr>
  <tr>
    <td>Describes the structure of the users table, including fields like Uuid, Name, Email, HashedPassword, and Created.</td>
    <td>Describes the structure of the todos table, including fields like ID, Body, Status, and Created.</td>
  </tr>
</table>
<hr>

<h3>Middleware</h3>
<table>
  <tr>
    <td>secureHeaders</td>
    <td>Adds security headers to HTTP responses.</td>
  </tr>
  <tr>
    <td>requireAuthentication</td>
    <td>Ensures that only authenticated users can access certain routes.</td>
  </tr>
  <tr>
    <td>authenticate</td>
    <td>Checks if a user is authenticated and adds relevant context to the request.</td>
  </tr>
  <tr>
    <td>noSurf</td>
    <td>Adds CSRF protection using a CSRF token.</td>
  </tr>
  <tr>
    <td>logRequest</td>
    <td>Logs details about each incoming request.</td>
  </tr>
  <tr>
    <td>recoverPanic</td>
    <td>Recovers from panics and returns a 500 Internal Server Error.</td>
  </tr>
</table>

<table>
  <caption>Routes</caption>
  <tr>
    <th>Route</th>
    <th>HTTP Method</th>
    <th>Description</th>
  </tr>
  <tr>
    <td>/signup</td>
    <td>POST</td>
    <td>Registers a new user and returns a JWT token.</td>
  </tr>
  <tr>
    <td>/login</td>
    <td>POST</td>
    <td>Authenticates a user and returns a JWT token.</td>
  </tr>
  <tr>
    <td>/logout</td>
    <td>POST</td>
    <td>Logs out a user and invalidates their JWT token.</td>
  </tr>
  <tr>
    <td>/todos</td>
    <td>GET</td>
    <td>Retrieves all todos for the authenticated user.</td>
  </tr>
  <tr>
    <td>/todos</td>
    <td>POST</td>
    <td>Creates a new todo item.</td>
  </tr>
  <tr>
    <td>/todos/{id}</td>
    <td>PUT</td>
    <td>Updates an existing todo item by ID.</td>
  </tr>
  <tr>
    <td>/todos/{id}</td>
    <td>DELETE</td>
    <td>Deletes a todo item by ID.</td>
  </tr>
  <tr>
    <td>/todos/{id}/toggle</td>
    <td>PUT</td>
    <td>Toggles the completion status of a todo item by ID.</td>
  </tr>
</table>


