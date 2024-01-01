package db

import (
	"database/sql"
	"fmt"
	"os"
	"path"
)

func Init(filepath string) *sql.DB {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil { // if err is not nil
		panic(err)
	}

	if db == nil { // if db is nil
		panic("db nil")
	}
	return db
}

func Migrate(db *sql.DB) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		panic(err)
	}

	for _, file := range files {
		path := path.Join("migrations", file.Name())
		fmt.Println("Applying migration ", file.Name())
		sql, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec(string(sql))
		// Exit if something goes wrong with our SQL statement above
		if err != nil {
			panic(err)
		}
	}

}
