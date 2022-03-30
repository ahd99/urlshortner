package sqldb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

type station struct {
	id        int
	group     int
	stId      string
	slotCount int
}

var ErrNoRowsFound = errors.New("NoRowsFound")

const (
	Query_Timeout  time.Duration = 5 * time.Second
	Update_Timeout time.Duration = 10 * time.Second
)

type psDb struct {
	db            *sql.DB
	ps_getStation *sql.Stmt
}

func NewPSDB() *psDb {
	return &psDb{}
}

func (p *psDb) connect() error {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "admin123",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "powerstation",
	}
	connString := cfg.FormatDSN()
	fmt.Printf("connection string: %s \n", connString)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		fmt.Printf("error open db: %s \n", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		fmt.Printf("error ping db: %s \n", err)
		return err
	}

	p.db = db
	fmt.Printf("Connected!\n")
	return nil
}

func (p *psDb) isConnected() bool {
	if p.db == nil {
		return false
	}

	if err := p.db.Ping(); err != nil {
		return false
	}
	return true
}

func (p *psDb) close() {
	p.ps_getStation.Close()
	p.db.Close()
}

func (p *psDb) getStationList(group int) ([]station, error) {
	rows, err := p.db.Query("SELECT id, stGroup, stId, slotCount FROM psRegistery WHERE stGroup = ?", group)
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

func (p *psDb) getStation(group int, stId string, ctx context.Context) (station, error) {
	var st station

	//TODO : this prepared-statement initialization is not thread-safe
	if p.ps_getStation == nil {
		stm, err := p.db.Prepare("SELECT id, stGroup, stId, slotCount FROM psRegistery WHERE stGroup = ? AND stId = ?")
		if err != nil {
			fmt.Printf("Error preparing statement: %v \n", err)
			return st, err
		}
		p.ps_getStation = stm
	}

	ctx, cancelFunc := context.WithTimeout(ctx, Query_Timeout)
	defer cancelFunc()

	row := p.ps_getStation.QueryRowContext(ctx, group, stId)
	err := row.Scan(&st.id, &st.group, &st.stId, &st.slotCount)
	if err == sql.ErrNoRows {
		fmt.Printf("no station found. group:%d,  tId:%q \n", group, stId)
		return st, ErrNoRowsFound
	} else if err == context.DeadlineExceeded {
		fmt.Printf("context.DeadlineExceeded. group:%d,  tId:%q   err:%v \n", group, stId, err)
		return st, err
	} else if err == context.Canceled {
		fmt.Printf("context.Canceled. group:%d,  tId:%q   err:%v \n", group, stId, err)
		return st, err
	} else if err != nil {
		fmt.Printf("Error querying station. group:%d,  tId:%q   err:%v \n", group, stId, err)
		return st, err
	}
	return st, nil
}

func (p *psDb) getStation_noPS(group int, stId string) (station, error) {
	row := p.db.QueryRow("SELECT id, stGroup, stId, slotCount FROM psRegistery WHERE stGroup = ? AND stId = ?", group, stId)
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

func (p *psDb) AddStationToRegistery(st station) (int, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), Update_Timeout)
	defer cancelFunc()

	res, err := p.db.ExecContext(ctx, "INSERT INTO psRegistery (stGroup, stId, slotCount) VALUES (?,?,?)", st.group, st.stId, st.slotCount)
	if err != nil { // includes context.DeadlineExceeded
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

func (p *psDb) RemoveStationFromRegistery(id int) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), Update_Timeout)
	defer cancelFunc()

	res, err := p.db.ExecContext(ctx, "DELETE FROM psRegistery WHERE id = ?", id)
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

func (p *psDb) getStation_contextTest(group int, stId string, ctx context.Context, sleepTime int) (station, error) {
	var st station

	ctx, cancelFunc := context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	row := p.db.QueryRowContext(ctx, "SELECT id, stGroup, stId, slotCount, sleep(?) FROM psRegistery WHERE stGroup = ? AND stId = ?", sleepTime, group, stId)
	var sl int
	err := row.Scan(&st.id, &st.group, &st.stId, &st.slotCount, &sl)
	fmt.Printf("ctx.Err():%v\n", ctx.Err())
	if err == sql.ErrNoRows {
		fmt.Printf("no station found. group:%d,  tId:%q \n", group, stId)
		return st, ErrNoRowsFound
	} else if err == context.DeadlineExceeded {
		fmt.Printf("context.DeadlineExceeded. group:%d,  tId:%q   err:%v \n", group, stId, err)
		return st, err
	} else if err == context.Canceled {
		fmt.Printf("context.Canceled. group:%d,  tId:%q   err:%v \n", group, stId, err)
		return st, err
	} else if err != nil {
		fmt.Printf("Error querying station. group:%d,  tId:%q   err:%v \n", group, stId, err)
		return st, err
	}
	return st, nil
}
