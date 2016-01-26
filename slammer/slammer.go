package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
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
	start := time.Now()
	errorCount := 0
	for err != io.EOF {
		// Build the line
		line, err = input.ReadString('\n')
		line = strings.TrimRight(line, "\r\n")
		// Import it into sql here
		if err == nil {
			_, err2 := db.Exec(line)
			if err2 != nil {
				errorCount++
				fmt.Println(err2)
			} else {
				time.Sleep(cfg.pauseInterval)
			}
		}
	}
	end := time.Now()
	diff := end.Sub(start)
	fmt.Println("Slammer Status:")
	fmt.Printf("Started at %s , Ended at %s, took %s\n", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), diff.String())
	fmt.Printf("Total errors: %d , Errors per second: %f\n", errorCount, float64(errorCount)/diff.Seconds())
}

// I went with an ENV var based config sheerly out of simplicity sake. I'm considering
// moving to CLI based flags instead but not worth it at the moment
func getConfig() (*config, error) {
	p := flag.String("p", "1s", "The time to pause between each call to the database")
	c := flag.String("c", "", "The connection string to use when connecting to the database")
	flag.Parse()

	if *c == "" {
		return nil, errors.New("You must provide a connection string using the -c option")
	}
	d, err := time.ParseDuration(*p)
	if err != nil {
		return nil, errors.New("You must provide a proper duration value with -p")
	}

	return &config{connString: *c, pauseInterval: d}, nil
}
