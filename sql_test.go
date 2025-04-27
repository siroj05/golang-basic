package golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO customer(id, name) VALUES('ismed', 'Ismed')"
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new customer")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "SELECT id, name FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string
		err = rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("Id :", id)
		fmt.Println("Name :", name)
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "SELECT id, name, email, balance, rating, birth_date, marriage, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var createdAt time.Time
		var birthDate sql.NullTime
		var marriage bool
		err = rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &marriage, &createdAt)
		if err != nil {
			panic(err)
		}
		fmt.Println("===============================")
		fmt.Println("Id :", id)
		fmt.Println("Name :", name)
		if email.Valid {
			fmt.Println("Email :", email.String)
		}
		fmt.Println("Balance :", balance)
		fmt.Println("Rating :", rating)
		if birthDate.Valid {
			fmt.Println("Birth Date :", birthDate.Time)
		}
		fmt.Println("Married :", marriage)
		fmt.Println("Created At :", createdAt)
	}
}

// simulasi sql injection
func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #"
	password := "salah"

	script := "SELECT username FROM user WHERE username = '" + username + "' AND password = '" + password + "' LIMIT 1"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses Login ", username)
	} else {
		fmt.Println("Gagal Login")
	}
}

// simulasi sql injection safe
func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin'; #"
	password := "salah"

	sqlQuery := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"
	rows, err := db.QueryContext(ctx, sqlQuery, username, password)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses Login ", username)
	} else {
		fmt.Println("Gagal Login")
	}
}

// exec context dengan params
func TestExecSqlParams(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	username := "eko'; DROP TABLE user; #"
	password := "Eko"

	ctx := context.Background()
	script := "INSERT INTO user(username, password) VALUES(?, ?)"
	_, err := db.ExecContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new user")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	email := "rojudin123@gmail.com"
	comment := "Test komen 123"

	ctx := context.Background()
	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, script, email, comment)
	if err != nil {
		panic(err)
	}

	insrtId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id", insrtId)
}

// prepare statement (wajib pake ini agar aman dan efisien)
func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}

	defer statement.Close()

	for i := 0; i < 10; i++ {
		email := "Siroj" + strconv.Itoa(i) + "email.com"
		comment := "Komentar ke" + strconv.Itoa(i)

		result, err := statement.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment ID", id)
	}

}

/*
Transaction di database adalah:
Kumpulan operasi SQL yang dijalankan bersama-sama dan dianggap satu paket.

Kalau semua operasi sukses, baru commit.
Kalau salah satu gagal, rollback semua ➔ balik ke kondisi awal.

=====================================================

🔥 Kenapa Butuh Transaction?
Bayangin kamu buat 2 query:

Kurangi saldo rekening A.

Tambah saldo rekening B.

Kalau query (1) sukses tapi query (2) gagal ➔
Tanpa transaction, uang bisa hilang ❗

Nah, transaction memastikan:

Kalau semua berhasil ➔ permanenkan perubahan (commit)

Kalau ada gagal ➔ batalkan semua (rollback)
=====================================================

✍️ Kapan Transaction Digunakan?
Transfer uang 💸

Update stok barang 🛒

# Multi-insert atau multi-update data penting

# CRUD data yang saling bergantung

Menghindari data korup karena sebagian query gagal
*/
func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	// do transaction
	for i := 0; i < 10; i++ {
		email := "Siroj" + strconv.Itoa(i) + "email.com"
		comment := "Komentar ke" + strconv.Itoa(i)

		result, err := tx.ExecContext(ctx, script, email, comment)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("Comment ID", id)
	}

	err = tx.Rollback()
	if err != nil {
		panic(err)
	}
}
