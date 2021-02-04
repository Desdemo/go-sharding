/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package database

import (
	"encoding/json"
	"github.com/XiaoMi/Gaea/mysql/types"
	"reflect"
	"testing"
	"time"

	"context"
)

func TestReadAllRedo(t *testing.T) {
	// Reuse code from tx_executor_test.
	_, tsv, db := newTestTxExecutor(t)
	defer db.Close()
	defer tsv.StopService()
	tpc := tsv.te.twoPC
	ctx := context.Background()

	conn, err := tsv.qe.conns.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Recycle()

	db.AddQuery(tpc.readAllRedo, &types.Result{})
	prepared, failed, err := tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var want []*tx.PreparedTx
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &types.Result{
		Fields: []*types.Field{
			{Type: types.VarChar},
			{Type: types.Int64},
			{Type: types.Int64},
			{Type: types.VarChar},
		},
		Rows: [][]types.Value{{
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt01"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*tx.PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &types.Result{
		Fields: []*types.Field{
			{Type: types.VarChar},
			{Type: types.Int64},
			{Type: types.Int64},
			{Type: types.VarChar},
		},
		Rows: [][]types.Value{{
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt01"),
		}, {
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt02"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01", "stmt02"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &types.Result{
		Fields: []*types.Field{
			{Type: types.VarChar},
			{Type: types.Int64},
			{Type: types.Int64},
			{Type: types.VarChar},
		},
		Rows: [][]types.Value{{
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt01"),
		}, {
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt02"),
		}, {
			types.NewVarBinary("dtid1"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt11"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01", "stmt02"},
		Time:    time.Unix(0, 1),
	}, {
		Dtid:    "dtid1",
		Queries: []string{"stmt11"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	if len(failed) != 0 {
		t.Errorf("ReadAllRedo (failed): %v, must be empty", jsonStr(failed))
	}

	db.AddQuery(tpc.readAllRedo, &types.Result{
		Fields: []*types.Field{
			{Type: types.VarChar},
			{Type: types.Int64},
			{Type: types.Int64},
			{Type: types.VarChar},
		},
		Rows: [][]types.Value{{
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt01"),
		}, {
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt02"),
		}, {
			types.NewVarBinary("dtid1"),
			types.NewVarBinary("Failed"),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt11"),
		}, {
			types.NewVarBinary("dtid2"),
			types.NewVarBinary("Failed"),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt21"),
		}, {
			types.NewVarBinary("dtid2"),
			types.NewVarBinary("Failed"),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt22"),
		}, {
			types.NewVarBinary("dtid3"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("stmt31"),
		}},
	})
	prepared, failed, err = tpc.ReadAllRedo(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*PreparedTx{{
		Dtid:    "dtid0",
		Queries: []string{"stmt01", "stmt02"},
		Time:    time.Unix(0, 1),
	}, {
		Dtid:    "dtid3",
		Queries: []string{"stmt31"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(prepared, want) {
		t.Errorf("ReadAllRedo: %s, want %s", jsonStr(prepared), jsonStr(want))
	}
	wantFailed := []*PreparedTx{{
		Dtid:    "dtid1",
		Queries: []string{"stmt11"},
		Time:    time.Unix(0, 1),
	}, {
		Dtid:    "dtid2",
		Queries: []string{"stmt21", "stmt22"},
		Time:    time.Unix(0, 1),
	}}
	if !reflect.DeepEqual(failed, wantFailed) {
		t.Errorf("ReadAllRedo failed): %s, want %s", jsonStr(failed), jsonStr(wantFailed))
	}
}

func TestReadAllTransactions(t *testing.T) {
	_, tsv, db := newTestTxExecutor(t)
	defer db.Close()
	defer tsv.StopService()
	tpc := tsv.te.twoPC
	ctx := context.Background()

	conn, err := tsv.qe.conns.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Recycle()

	db.AddQuery(tpc.readAllTransactions, &types.Result{})
	distributed, err := tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	var want []*DistributedTx
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions: %s, want %s", jsonStr(distributed), jsonStr(want))
	}

	db.AddQuery(tpc.readAllTransactions, &types.Result{
		Fields: []*types.Field{
			{Type: types.VarChar},
			{Type: types.Int64},
			{Type: types.Int64},
			{Type: types.VarChar},
			{Type: types.VarChar},
		},
		Rows: [][]types.Value{{
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("ks01"),
			types.NewVarBinary("shard01"),
		}},
	})
	distributed, err = tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*DistributedTx{{
		Dtid:    "dtid0",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []Target{{
			Schema:     "ks01",
			DataSource: "shard01",
		}},
	}}
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions:\n%s, want\n%s", jsonStr(distributed), jsonStr(want))
	}

	db.AddQuery(tpc.readAllTransactions, &types.Result{
		Fields: []*types.Field{
			{Type: types.VarChar},
			{Type: types.Int64},
			{Type: types.Int64},
			{Type: types.VarChar},
			{Type: types.VarChar},
		},
		Rows: [][]types.Value{{
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("ks01"),
			types.NewVarBinary("shard01"),
		}, {
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("ks02"),
			types.NewVarBinary("shard02"),
		}},
	})
	distributed, err = tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*DistributedTx{{
		Dtid:    "dtid0",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []Target{{
			Schema:     "ks01",
			DataSource: "shard01",
		}, {
			Schema:     "ks02",
			DataSource: "shard02",
		}},
	}}
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions:\n%s, want\n%s", jsonStr(distributed), jsonStr(want))
	}

	db.AddQuery(tpc.readAllTransactions, &types.Result{
		Fields: []*types.Field{
			{Type: types.VarChar},
			{Type: types.Int64},
			{Type: types.Int64},
			{Type: types.VarChar},
			{Type: types.VarChar},
		},
		Rows: [][]types.Value{{
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("ks01"),
			types.NewVarBinary("shard01"),
		}, {
			types.NewVarBinary("dtid0"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("ks02"),
			types.NewVarBinary("shard02"),
		}, {
			types.NewVarBinary("dtid1"),
			types.NewInt64(RedoStatePrepared),
			types.NewVarBinary("1"),
			types.NewVarBinary("ks11"),
			types.NewVarBinary("shard11"),
		}},
	})
	distributed, err = tpc.ReadAllTransactions(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want = []*DistributedTx{{
		Dtid:    "dtid0",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []Target{{
			Schema:     "ks01",
			DataSource: "shard01",
		}, {
			Schema:     "ks02",
			DataSource: "shard02",
		}},
	}, {
		Dtid:    "dtid1",
		State:   "PREPARE",
		Created: time.Unix(0, 1),
		Participants: []Target{{
			Schema:     "ks11",
			DataSource: "shard11",
		}},
	}}
	if !reflect.DeepEqual(distributed, want) {
		t.Errorf("ReadAllTransactions:\n%s, want\n%s", jsonStr(distributed), jsonStr(want))
	}
}

func jsonStr(v interface{}) string {
	out, _ := json.Marshal(v)
	return string(out)
}
