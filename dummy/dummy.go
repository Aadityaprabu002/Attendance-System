package main

import (
	"attsys/connections"
	"database/sql"
	"fmt"
)

func main() {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, "postgres")
	db, _ := sql.Open("postgres", conn)
	defer db.Close()
	query := `insert into employee values(1023,123412.2,'apk') returning numeric`
	result, _ := db.Query(query)
	var sal float64
	for result.Next() {
		result.Scan(&sal)
	}
	fmt.Println(sal)
}
