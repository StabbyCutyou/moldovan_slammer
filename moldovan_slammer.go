package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	connString    string
	inputFile     string
	pauseInterval time.Duration
}

func main() {
	fmt.Print("Welcome to the Moldovan Slammer")
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", cfg.connString)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Open(cfg.inputFile)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	reader := bufio.NewReader(f)

	for {
		line, err := reader.ReadString('\n')
		// Import it into sql here
		if err == nil || err == io.EOF {
			_, err = db.Exec(line)

			if err != nil {
				fmt.Println(err)
			}
		} else if err != nil {

		}

		time.Sleep(cfg.pauseInterval)
	}

}

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
		inputFile:     os.Getenv("MS_INPUTFILE"),
		pauseInterval: d,
	}

	if cfg.connString == "" {
		return nil, errors.New("MS_CONNSTRING must be set")
	}

	if cfg.inputFile == "" {
		return nil, errors.New("MS_INPUTFILE must be set")
	}

	return cfg, nil
}
