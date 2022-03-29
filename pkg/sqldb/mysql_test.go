package sqldb

import (
	"fmt"
	"testing"
)

func TestMysqlStationListQuerySuccess(t *testing.T) {
	if !isConnected() {
		connect()
	}
	sts, err := getStationList(2)
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
	if !isConnected() {
		connect()
	}
	st, err := getStation(2, stId)
	if err != nil {
		t.Fatalf("Error in query:%v", err)
	}

	if st.group != 2 || st.stId != stId {
		t.Fatalf("station data error. received group:%v", st)
	}
	fmt.Printf("station found: %v", st)
}

func TestMysqlStationQueryRowNotFound(t *testing.T) {
	stId := "a3a3a3a3a3a3a3a3"
	if !isConnected() {
		connect()
	}
	_, err := getStation(3, stId)
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

	if !isConnected() {
		connect()
	}
	id, err := AddStationToRegistery(st)
	if err != nil {
		t.Fatalf("Error inserting station: %v", err)
	}

	if id <= 0 {
		t.Fatalf("Invalid id after insert station: %d", id)
	}

	fmt.Printf("st inserted. id:%d", id)

	err = RemoveStationFromRegistery(id)
	if err != nil {
		t.Fatalf("Error removing station: %v", err)
	}
}

func TestMysqlRemoveStationNotFound(t *testing.T) {
	id := 999999

	if !isConnected() {
		connect()
	}
	err := RemoveStationFromRegistery(id)
	if err != ErrNoRowsFound {
		t.Fatalf("ErrNoRowsFound expected but  %v  received.", err)
	}
}
