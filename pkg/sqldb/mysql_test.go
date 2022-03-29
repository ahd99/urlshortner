package sqldb

import (
	"fmt"
	"testing"
)

var db_test *psDb

func prepareConnection() error {
	if db_test == nil {
		db_test = NewPSDB()
	}
	if db_test.isConnected() {
		return nil
	}
	err := db_test.connect()
	return err
}

func TestMysqlStationListQuerySuccess(t *testing.T) {
	err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	sts, err := db_test.getStationList(2)
	if err != nil {
		t.Fatalf("Error in query:%v", err)
	}

	if len(sts) != 2 {
		t.Fatalf("sts len error. expected: 2   received:%d", len(sts))
	}

	fmt.Printf("station list. list:%v", sts)
}

func TestMysqlStationQuerySuccess(t *testing.T) {
	stId := "a2a2a2a2a2a2a2a2"
	group := 2

	err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	st, err := db_test.getStation(group, stId)
	if err != nil {
		t.Fatalf("Error in query:%v", err)
	}

	if st.group != group || st.stId != stId {
		t.Fatalf("station data error. received group:%v", st)
	}
	fmt.Printf("station found: %v", st)
}

func TestMysqlStationQueryRowNotFound(t *testing.T) {
	stId := "a3a3a3a3a3a3a3a3"
	group := 33
	err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	_, err = db_test.getStation(group, stId)
	if err != ErrNoRowsFound {
		t.Fatalf("ErrNoRowsFound expected but %v  received.", err)
	}
}

func TestMysqlInsertAndRemoveStationSuccess(t *testing.T) {
	st := station{
		group:     2,
		stId:      "a3a3a3a3a3a3a3a3",
		slotCount: 10,
	}

	err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	id, err := db_test.AddStationToRegistery(st)
	if err != nil {
		t.Fatalf("Error inserting station: %v", err)
	}

	if id <= 0 {
		t.Fatalf("Invalid id after insert station: %d", id)
	}

	fmt.Printf("st inserted. id:%d", id)

	err = db_test.RemoveStationFromRegistery(id)
	if err != nil {
		t.Fatalf("Error removing station: %v", err)
	}
}

func TestMysqlRemoveStationNotFound(t *testing.T) {
	id := 999999

	err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	err = db_test.RemoveStationFromRegistery(id)
	if err != ErrNoRowsFound {
		t.Fatalf("ErrNoRowsFound expected but  %v  received.", err)
	}
}
