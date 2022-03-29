package sqldb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

func connect() error {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "admin123",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "powerstation",
	}
	connString := cfg.FormatDSN()
	fmt.Printf("connection string: %s \n", connString)

	var err error
	db, err = sql.Open("mysql", connString)
	if err != nil {
		fmt.Printf("error open db: %s \n", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error ping db: %s \n", err)
		return err
	}

	fmt.Printf("Connected!\n")
	return nil
}

func isConnected() bool {
	if db == nil {
		return false
	}

	if err := db.Ping(); err != nil {
		return false
	}
	return true
}

func close() {
	db.Close()
}

type station struct {
	id        int
	group     int
	stId      string
	slotCount int
}

var ErrNoRowsFound = errors.New("NoRowsFound")

func getStationList(group int) ([]station, error) {
	rows, err := db.Query("SELECT id, stGroup, stId, slotCount FROM psRegistery WHERE stGroup = ?", group)
	if err != nil {
		fmt.Printf("error query stations. group:%d   err:%v \n", group, err)
		return nil, err
	}
	defer rows.Close()

	var stations []station
	for rows.Next() {
		var st station
		err := rows.Scan(&st.id, &st.group, &st.stId, &st.slotCount)
		if err != nil {
			fmt.Printf("error read columns data:%v \n", err)
			return nil, err
		}
		stations = append(stations, st)
	}
	if err := rows.Err(); err != nil {
		fmt.Printf("error in rows.Err() :%v \n", rows.Err())
		return nil, err
	}
	return stations, nil
}

func getStation(group int, stId string) (station, error) {
	row := db.QueryRow("SELECT id, stGroup, stId, slotCount FROM psRegistery WHERE stGroup = ? AND stId = ?", group, stId)
	var st station
	err := row.Scan(&st.id, &st.group, &st.stId, &st.slotCount)
	if err == sql.ErrNoRows {
		fmt.Printf("no station found. group:%d,  tId:%q \n", group, stId)
		return st, ErrNoRowsFound
	} else if err != nil {
		fmt.Printf("Error querying station. group:%d,  tId:%q   err:%v \n", group, stId, err)
		return st, err
	}
	return st, nil
}

func AddStationToRegistery(st station) (int, error) {
	res, err := db.Exec("INSERT INTO psRegistery (stGroup, stId, slotCount) VALUES (?,?,?)", st.group, st.stId, st.slotCount)
	if err != nil {
		fmt.Printf("Error inserting station. st:%v,   err:%v \n", st, err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting inserted row id. st:%v,   err:%v \n", st, err)
		return 0, err
	}
	return int(id), nil
}

func RemoveStationFromRegistery(id int) error {
	res, err := db.Exec("DELETE FROM psRegistery WHERE id = ?", id)
	if err != nil {
		fmt.Printf("Error removing station. id:%v,   err:%v \n", id, err)
		return err
	}
	afftectedRows, err := res.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting affected rows after removing station. id:%v,   err:%v \n", id, err)
		return err
	}
	if afftectedRows == 0 {
		return ErrNoRowsFound
	}
	return nil
}
