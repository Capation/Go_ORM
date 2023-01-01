package sql_demo

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // 匿名导入  init()
	"github.com/stretchr/testify/require"
	"log"
	"testing"
	"time"
)

func TestDB(t *testing.T) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/sql_test"
	db, err := sql.Open("mysql", dsn)
	require.NoError(t, err)
	defer db.Close()
	err = db.Ping()
	if err != nil {
		fmt.Printf("connect to db failed, err:%v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// 除了 SELECT 语句，都是使用 ExecContext
	_, err = db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS test_model(
    id INTEGER PRIMARY KEY,
    first_name TEXT NOT NULL,
    age INTEGER,
    last_name TEXT NOT NULL
)
`)
	// 完成了建表
	require.NoError(t, err)

	// 使用 ？ 作为查询的参数的占位符
	//res, err := db.ExecContext(ctx, "INSERT INTO test_model(`id`, `first_name`, `age`, `last_name`) VALUES (?, ?, ?, ?)",
	//	1, "Tom", 18, "Jerry")
	//require.NoError(t, err)
	//
	//affected, err := res.RowsAffected()
	//require.NoError(t, err)
	//fmt.Println("受影响的行数:", affected)
	//lastID, err := res.LastInsertId()
	//require.NoError(t, err)
	//fmt.Println("最后插入的ID:", lastID)

	row := db.QueryRowContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
	require.NoError(t, row.Err())
	tm := TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.NoError(t, err)

	// 查询不到
	row = db.QueryRowContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 2)
	require.NoError(t, row.Err())
	tm = TestModel{}
	err = row.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
	require.Error(t, sql.ErrNoRows, err)

	// 批量查询
	rows, err := db.QueryContext(ctx, "SELECT `id`, `first_name`, `age`, `last_name` FROM `test_model` WHERE `id` = ?", 1)
	require.NoError(t, row.Err())
	for rows.Next() {
		tm = TestModel{}
		err = rows.Scan(&tm.Id, &tm.FirstName, &tm.Age, &tm.LastName)
		require.NoError(t, err)
		log.Println(tm)
	}

	cancel()
}

type TestModel struct {
	Id        int64
	FirstName string
	Age       int8
	LastName  *sql.NullString
}
