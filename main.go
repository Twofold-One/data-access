package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/joho/godotenv"
)

var db *sql.DB

type Album struct {
	ID int64
	Title string
	Artist string
	Price float32
}

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	var err error
	db, err = sql.Open("pgx", os.Getenv("DBURL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	fmt.Println("DB URL:", os.Getenv("DBURL"))

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	allAlbums, err := allAlbums()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All albums found: %v\n", allAlbums)

	albums, err := albumsByArtist("Metallica")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found by name: %v\n", albums)

	singleAlbum, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	} 
	fmt.Printf("Album found by id: %v\n", singleAlbum)

	albumID, albumTitle, err := addAlbum(Album{
		Title: "Reign in Blood",
		Artist: "Slayer",
		Price: 50.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID and title of added album: %v, %v\n", albumID, albumTitle)
}

// allAlbums queries for all available albums.
func allAlbums() ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album")
	if err != nil {
		return nil, fmt.Errorf("err %v", err)
	}
	defer rows.Close()
	
	for rows.Next() {
	var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist: %v", err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist: %v", err)
	}
	return albums, nil
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID and Title of the new entry
func addAlbum(alb Album) (int64, string, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, "", fmt.Errorf("addAlbum: %v", err)
	}
	fmt.Println(result)
	return alb.ID, alb.Title, nil
}


