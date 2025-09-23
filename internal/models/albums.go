package models

import (
	"database/sql"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Album struct {
	AlbumId  int
	Title    string
	ArtistId int
	Artist   string
}

type AlbumModel struct {
	DB *sql.DB
}

func (m *AlbumModel) ListAlbums(artistId, artist string) ([]Album, error) {
	aId, _ := strconv.Atoi(artistId)
	rows, err := m.DB.Query("select AlbumId, Title, ArtistId from albums where ArtistId = ?", aId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var al []Album

	for rows.Next() {
		var a Album
		a.Artist = artist
		err = rows.Scan(&a.AlbumId, &a.Title, &a.ArtistId)
		if err != nil {
			return nil, err
		}
		al = append(al, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return al, nil
}
