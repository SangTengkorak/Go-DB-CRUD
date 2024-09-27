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
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
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

// Handler to insert a new user
func createUser(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    sqlStatement := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`
    err = db.QueryRow(sqlStatement, user.Name, user.Email).Scan(&user.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

// Handler to update a user's name by ID
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

    sqlStatement := `UPDATE users SET name=$1 WHERE id=$2`
    _, err = db.Exec(sqlStatement, user.Name, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
}

func main() {
    // Initialize the database
    initDB()

    // Initialize router
    r := mux.NewRouter()

    // Define API routes
    r.HandleFunc("/users", createUser).Methods("POST")
    r.HandleFunc("/users/{id}", updateUser).Methods("PUT")

    // Start server
    fmt.Println("Server is running on port 8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}
