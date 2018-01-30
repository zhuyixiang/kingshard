// Copyright 2016 The kingshard Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package sqlparser

import (
	"testing"
	"fmt"
)

func testParse(t *testing.T, sql string) Statement {
	stat, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	return stat
}

func TestSet(t *testing.T) {
	sql := "set names gbk"
	testParse(t, sql)
}

func TestSimpleSelect(t *testing.T) {
	sql := "select last_insert_id() as a"
	testParse(t, sql)
}

func TestMixer(t *testing.T) {
	sql := `admin upnode("node1", "master", "127.0.0.1")`
	testParse(t, sql)

	sql = "show databases"
	testParse(t, sql)

	sql = "show tables from abc"
	testParse(t, sql)

	sql = "show tables from abc like a"
	testParse(t, sql)

	sql = "show tables from abc where a = 1"
	testParse(t, sql)

	sql = "show proxy abc"
	testParse(t, sql)
}

func TestSelect(t *testing.T) {
	sql := "select * from table1 where id = ? and name =? and type= ?"
	statement := testParse(t, sql)
	stat1 := statement.(*Select)
	fmt.Println(String(stat1.From))
	if stat1.Where == nil {
		fmt.Println("no Where")
	} else {
		testFillOne(stat1.Where.Expr)
		fmt.Println(String(stat1))
	}

	//sql = "update table2 set col1 = ? , col2 = ? where id = ?"
	//statement = testParse(t, sql)
	//stat2 := statement.(*Update)
	//fmt.Println(String(stat2.Table))
	//fmt.Println(String(stat1.Where.Expr))

	//sql := "insert into table1(col1, col2) values(?, ?)"
	//statement := testParse(t, sql)
	//stat3 := statement.(*Insert)
	//fmt.Println(String(stat3.Table))
	//fmt.Println(String(stat3.Rows))
}

func testFillOne(expr BoolExpr) {
	switch v:= expr.(type) {
	case *ComparisonExpr:
		comp := expr.(*ComparisonExpr)
		fmt.Println(comp.Left)
		bytes := []byte{'C','D'}
		comp.Right = StrVal(bytes)
	case *AndExpr:
		and := expr.(*AndExpr)
		testFillOne(and.Left)
		testFillOne(and.Right)
	default:
		fmt.Println(v)
	}
}
