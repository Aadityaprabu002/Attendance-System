package classroom

import (
	connections "attsys/connections"
	"database/sql"
	"fmt"
	"strings"
)

func IsClassRoomExist(ClassroomId string) bool {
	conn := fmt.Sprintf("host = %s port = %d user = %s password = %d dbname = %s sslmode = disable", connections.Host, connections.Port, connections.User, connections.Password, connections.DBname)
	db, err := sql.Open("postgres", conn)
	if err != nil {
		panic("failed to establish connection with sql")
	}
	defer db.Close()
	query := fmt.Sprintf(`select exists(select 1 from classrooms where classroom_id = '%s')`, ClassroomId)
	query = strings.TrimSpace(query)

	result, err := db.Query(query)

	if err != nil {
		panic(err)
	}

	var exist bool
	for result.Next() {
		result.Scan(&exist)
	}
	return exist

}
