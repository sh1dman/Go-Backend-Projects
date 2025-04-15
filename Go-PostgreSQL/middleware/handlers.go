package middleware

import (
	"Go-PostgreSQL/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"encoding/json"

	"github.com/joho/gotdotenv"

	"github.com/gorilla/mux"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message, omitempty"`
}

func CreateConnection() *sql.DB {
	err := gotdotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("PSTGRES_URL"))

	if err != nil { //this checks if the connection to the database is successful
		// if not, it will panic and print the error
		panic(err)
	}

	err = db.Ping()

	if err != nil { //this checks if the PING to the database is successful if not it goes panic
		panic(err)
	}

	fmt.Println("Successfuly connected to postgres database")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatal("Error decoding the JSON request body. %v ", err)
	}

	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}
	json.Encoder(w).Encode(res) // after we got the res, we need to encode it because at default go lang wont encode JSON
}

func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"]) // converting the string to int

	if err != nil {
		log.Fatalf("Error converting string to int. %v", err)
	}
	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Error getting stock. %v", err)
	}

	json.Encoder(w).Encode(stock) // after we got the stock, we need to encode it because at default go lang wont encode JSON
}

func GetAllStock() {

}

func UpdateStock() {

}

func DeleteStock() {

}
