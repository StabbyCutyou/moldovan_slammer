package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	// Load the driver only
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	connString    string
	pauseInterval time.Duration
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	fmt.Print("Welcome to the Moldovan Slammer\n")
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", cfg.connString)
	if err != nil {
		log.Fatal(err)
	}
	input := bufio.NewReader(os.Stdin)
	err = nil
	line := ""
	for err != io.EOF {
		// Build the line
		line, err = input.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		// Import it into sql here
		if err == nil {
			_, err2 := db.Exec(line)
			if err2 != nil {
				fmt.Println(err2)
			} else {
				time.Sleep(cfg.pauseInterval)
			}
		}
	}
}

// I went with an ENV var based config sheerly out of simplicity sake. I'm considering
// moving to CLI based flags instead but not worth it at the moment
func getConfig() (*config, error) {
	duration := os.Getenv("MS_PAUSEINTERVAL")
	if duration == "" {
		return nil, errors.New("MS_PAUSEINTERVAL must be set")
	}
	d, err := time.ParseDuration(duration)
	if err != nil {
		return nil, err
	}

	cfg := &config{
		connString:    os.Getenv("MS_CONNSTRING"),
		pauseInterval: d,
	}

	if cfg.connString == "" {
		return nil, errors.New("MS_CONNSTRING must be set")
	}

	return cfg, nil
}
