package sqldb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var db_ps_test *psDb

func prepareConnection() (*psDb, error) {
	if db_ps_test == nil {
		db_ps_test = NewPSDB()
	}
	if db_ps_test.isConnected() {
		return db_ps_test, nil
	}
	err := db_ps_test.connect()
	return db_ps_test, err
}

func TestMysqlStationListQuerySuccess(t *testing.T) {
	db_ps, err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	sts, err := db_ps.getStationList(2)
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

	db_ps, err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	ctx := context.Background()
	st, err := db_ps.getStation(group, stId, ctx)
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
	db_ps, err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	ctx := context.Background()
	_, err = db_ps.getStation(group, stId, ctx)
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

	db_ps, err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	id, err := db_ps.AddStationToRegistery(st)
	if err != nil {
		t.Fatalf("Error inserting station: %v", err)
	}

	if id <= 0 {
		t.Fatalf("Invalid id after insert station: %d", id)
	}

	fmt.Printf("st inserted. id:%d", id)

	err = db_ps.RemoveStationFromRegistery(id)
	if err != nil {
		t.Fatalf("Error removing station: %v", err)
	}
}

func TestMysqlRemoveStationNotFound(t *testing.T) {
	id := 999999

	db_ps, err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	err = db_ps.RemoveStationFromRegistery(id)
	if err != ErrNoRowsFound {
		t.Fatalf("ErrNoRowsFound expected but  %v  received.", err)
	}
}

func TestMysqlStationQuery_Timeout(t *testing.T) {
	stId := "a2a2a2a2a2a2a2a2"
	group := 2

	db_ps, err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	ctx := context.Background()
	_, err = db_ps.getStation_contextTest(group, stId, ctx, 10)
	if err != context.DeadlineExceeded {
		t.Fatalf("context.DeadlineExceeded expected but %v  received", err)
	}
}

func TestMysqlStationQuery_Cancel(t *testing.T) {
	stId := "a2a2a2a2a2a2a2a2"
	group := 2

	db_ps, err := prepareConnection()
	if err != nil {
		t.Fatalf("Error in connection: %v", err)
	}

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go func() {
		time.Sleep(3 * time.Second)
		cancelFunc()
	}()
	_, err = db_ps.getStation_contextTest(group, stId, ctx, 10)
	if err != context.Canceled {
		t.Fatalf("context.DeadlineExceeded expected but %v  received", err)
	}
}
