package middleware

import (
	"Go-PostgreSQL/models" // models package where Stock schema is defined
	"database/sql"
	"encoding/json" // package to encode and decode the json into struct and vice versa
	"fmt"
	"log"
	"net/http" // used to access the request and response object of the api
	"os"       // used to read the environment variable
	"strconv"  // package used to covert string into int type

	"github.com/gorilla/mux" // used to get the params from the route

	"github.com/joho/godotenv" // package used to read the .env file
	_ "github.com/lib/pq"      // postgres golang driver
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message, omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

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
	json.NewEncoder(w).Encode(res) // after we got the res, we need to encode it because at default go lang wont encode JSON
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

	json.NewEncoder(w).Encode(stock) // after we got the stock, we need to encode it because at default go lang wont encode JSON
}

func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := getAllStock()
	if err != nil {
		log.Fatalf("Error unable on getting all stocks. %v", err)
	}
	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Error converting string to int. %v", err)
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock) // whats coming in the body we need to decode it;
	if err != nil {
		log.Fatalf("Error decoding the JSON request body. %v ", err)
	}

	updatedRows := updateStock(int64(id), stock) // we need to pass the id and the stock to update the stock
	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected. %v", updatedRows)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {

	// because we deleted based on ID then we need mux to get the ID
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Error converting string to int. %v", err)
	}

	deletedRow := deletedStock(int64(id))
	msg := fmt.Sprintf("Stock deleted successfully. Total rows/records. %v", deletedRow)

	res := response{
		ID:      int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record %v", id)
	return id
}

func getStock(id int64) (models.Stock, error) {
	db := createConnection()

	defer db.Close()

	var stock models.Stock
	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err {

	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, err
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}
	return stock, nil

}

func getAllStock() ([]models.Stock, error) {
	db := createConnection()

	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stock models.Stock
		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}
		stocks = append(stocks, stock)
	}

	return stocks, nil

}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`

	res, err := db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/records affetced. %v", rowsAffected)
	return rowsAffected

}

func deletedStock(id int64) int64 {

	// create the postgres db connection
	db := createConnection()

	// close the db connection
	defer db.Close()

	// create the delete sql query
	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
