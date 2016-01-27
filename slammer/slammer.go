package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	// Load the driver only
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	connString    string
	pauseInterval time.Duration
	workers       int
}

type result struct {
	start     time.Time
	end       time.Time
	workCount int
	errors    int
}

func main() {
	cfg, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := sql.Open("mysql", cfg.connString)
	if err != nil {
		log.Fatal(err)
	}

	// Declare the channel we'll be using
	inputChan := make(chan (string))
	// Declare the channel that will gather results
	outputChan := make(chan (result), cfg.workers)

	// Start the pool of workers up, reading from the channel
	for i := 0; i < cfg.workers; i++ {
		go func(ic <-chan string, oc chan<- result, d *sql.DB, pause time.Duration) {
			r := result{start: time.Now()}
			for line := range ic {
				_, err := db.Exec(line)
				r.workCount++
				if err != nil {
					r.errors++
					//fmt.Println(err)
				} else {
					time.Sleep(pause)
				}
			}
			r.end = time.Now()
			oc <- r
		}(inputChan, outputChan, db, cfg.pauseInterval)
	}

	// Read from STDIN in the main thread
	input := bufio.NewReader(os.Stdin)
	err = nil
	line := ""
	for err != io.EOF {
		line, err = input.ReadString('\n')
		if err == nil {
			fmt.Println(line)
			line = strings.TrimRight(line, "\r\n")

			inputChan <- line
		} else {
			fmt.Println(err)
		}
	}

	fmt.Printf("HEY")

	// Close the channel, since it's done receiving input
	close(inputChan)

	// Collect all results, report them. This will block and wait until all results
	// are in

	fmt.Println("Slammer Status:")
	for i := 0; i < cfg.workers; i++ {
		r := <-outputChan
		diff := r.end.Sub(r.start)
		fmt.Printf("Worker #%d\n", i)
		fmt.Printf("Started at %s , Ended at %s, took %s\n", r.start.Format("2006-01-02 15:04:05"), r.end.Format("2006-01-02 15:04:05"), diff.String())
		fmt.Printf("Total errors: %d , Percentage errors: %f, Average errors per second: %f\n", r.errors, float64(r.errors)/float64(r.workCount), float64(r.errors)/diff.Seconds())
	}

	close(outputChan)
}

// I went with an ENV var based config sheerly out of simplicity sake. I'm considering
// moving to CLI based flags instead but not worth it at the moment
func getConfig() (*config, error) {
	p := flag.String("p", "1s", "The time to pause between each call to the database")
	c := flag.String("c", "", "The connection string to use when connecting to the database")
	w := flag.Int("w", 1, "The number of workers to use. A number greater than 1 will enable statements to be issued concurrently")
	flag.Parse()

	if *c == "" {
		return nil, errors.New("You must provide a connection string using the -c option")
	}
	d, err := time.ParseDuration(*p)
	if err != nil {
		return nil, errors.New("You must provide a proper duration value with -p")
	}

	if *w <= 0 {
		return nil, errors.New("You must provide a worker count > 0 with -w")
	}

	return &config{connString: *c, pauseInterval: d, workers: *w}, nil
}
