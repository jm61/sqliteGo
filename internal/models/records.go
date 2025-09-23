package models

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Record struct {
	Name         string
	Composer     string
	Milliseconds int
	UnitPrice    float32
}

type RecordModel struct {
	DB *sql.DB
}

func (m *RecordModel) List(albumId int) ([]Record, error) {

	rows, err := m.DB.Query("select name, composer, milliseconds, unitprice from tracks where albumid = ?", albumId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var rec []Record

	for rows.Next() {
		var r Record
		err = rows.Scan(&r.Name, &r.Composer, &r.Milliseconds, &r.UnitPrice)
		if err != nil {
			return nil, err
		}
		rec = append(rec, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rec, nil
}
