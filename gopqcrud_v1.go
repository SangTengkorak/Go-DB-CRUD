package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

var db *sql.DB

const (
    DB_USER     = "postgres"
    DB_PASSWORD = "tengkorak123"
    DB_NAME     = "generic_db"
)

// User struct
type User struct {
    ID        int    `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    Age       int    `json:"age"`
    City      string `json:"city"`
    CreatedAt string `json:"created_at"`
}

// Initialize database connection
func initDB() {
    var err error
    psqlInfo := fmt.Sprintf("host=localhost port=5432 user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Database connection successful!")
}

// Create a new user
func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `INSERT INTO users (name, email, age, city) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
    err = db.QueryRow(sqlStatement, user.Name, user.Email, user.Age, user.City).Scan(&user.ID, &user.CreatedAt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

// View all users
func getAllUsers(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT id, name, email, age, city, created_at FROM users")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    users := []User{}
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.City, &user.CreatedAt)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        users = append(users, user)
    }

    json.NewEncoder(w).Encode(users)
}

// Get a specific user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user User
    sqlStatement := `SELECT id, name, email, age, city, created_at FROM users WHERE id=$1`
    err = db.QueryRow(sqlStatement, userID).Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.City, &user.CreatedAt)
    if err == sql.ErrNoRows {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    } else if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

// Update a user's information
func updateUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user User
    err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `UPDATE users SET name=$1, email=$2, age=$3, city=$4 WHERE id=$5`
    _, err = db.Exec(sqlStatement, user.Name, user.Email, user.Age, user.City, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

// Update a specific column (email) of a user
func updateUserEmail(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user User
    err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `UPDATE users SET email=$1 WHERE id=$2`
    _, err = db.Exec(sqlStatement, user.Email, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "User's email updated successfully"})
}

// Update a specific column (city) of a user
func updateUserCity(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    var user User
    err = json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `UPDATE users SET city=$1 WHERE id=$2`
    _, err = db.Exec(sqlStatement, user.City, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "User's city updated successfully"})
}

// Delete a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userID, err := strconv.Atoi(vars["id"])
    if err != nil {
        http.Error(w, "Invalid user ID", http.StatusBadRequest)
        return
    }

    sqlStatement := `DELETE FROM users WHERE id=$1`
    _, err = db.Exec(sqlStatement, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}

func main() {
    // Initialize the database connection
    initDB()

    // Initialize the router
    r := mux.NewRouter()

    // Define API routes
    r.HandleFunc("/users", createUser).Methods("POST")          // Insert a new user
    r.HandleFunc("/users", getAllUsers).Methods("GET")          // View all users
    r.HandleFunc("/users/{id}", getUser).Methods("GET")         // Select a specific user
    r.HandleFunc("/users/{id}", updateUser).Methods("PUT")      // Update a user's information
    r.HandleFunc("/users/{id}/city", updateUserCity).Methods("PATCH") // Update a specific column (city)
	r.HandleFunc("/users/{id}/email", updateUserEmail).Methods("PATCH") // Update a specific column (email)
    r.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")   // Delete a user

    // Start the server
    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}