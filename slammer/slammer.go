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
	"sync"
	"time"

	// Load the driver only
	_ "github.com/go-sql-driver/mysql"
)

type config struct {
	connString    string
	pauseInterval time.Duration
	workers       int
	debugMode     bool
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
	// Declare a waitgroup to help prevent log interleaving - I technically do not
	// need one, but without it, I find there are stray log messages creeping into
	// the final report. Setting sync() on STDOUT didn't seem to fix it
	var wg sync.WaitGroup
	wg.Add(cfg.workers)
	// Start the pool of workers up, reading from the channel
	for i := 0; i < cfg.workers; i++ {
		go func(workerNum int, ic <-chan string, oc chan<- result, d *sql.DB, done *sync.WaitGroup, pause time.Duration, debugMode bool) {
			r := result{start: time.Now()}
			for line := range ic {
				if debugMode {
					log.Printf("Worker #%d: About to run %s", workerNum, line)
				}
				_, err := db.Exec(line)
				r.workCount++
				if err != nil {
					r.errors++
					if debugMode {
						log.Printf("Worker #%d: %s", workerNum, err.Error())
					}
				} else {
					time.Sleep(pause)
				}
			}
			r.end = time.Now()
			oc <- r
			done.Done()
		}(i, inputChan, outputChan, db, &wg, cfg.pauseInterval, cfg.debugMode)
	}

	// Read from STDIN in the main thread
	input := bufio.NewReader(os.Stdin)
	err = nil
	line := ""
	for err != io.EOF {
		line, err = input.ReadString('\n')
		if err == nil {
			line = strings.TrimRight(line, "\r\n")

			inputChan <- line
		} else if cfg.debugMode {
			log.Println(err)
		}
	}

	// Close the channel, since it's done receiving input
	close(inputChan)
	wg.Wait()
	// Collect all results, report them. This will block and wait until all results
	// are in
	fmt.Println("Slammer Status:")
	for i := 0; i < cfg.workers; i++ {
		r := <-outputChan
		diff := r.end.Sub(r.start)
		fmt.Printf("---- Worker #%d ----\n", i)
		fmt.Printf("  Started at %s , Ended at %s, took %s\n", r.start.Format("2006-01-02 15:04:05"), r.end.Format("2006-01-02 15:04:05"), diff.String())
		fmt.Printf("  Total errors: %d , Percentage errors: %f, Average errors per second: %f\n", r.errors, float64(r.errors)/float64(r.workCount), float64(r.errors)/diff.Seconds())
	}

	// Lets just be nice and tidy
	close(outputChan)
}

// I went with an ENV var based config sheerly out of simplicity sake. I'm considering
// moving to CLI based flags instead but not worth it at the moment
func getConfig() (*config, error) {
	p := flag.String("p", "1s", "The time to pause between each call to the database")
	c := flag.String("c", "", "The connection string to use when connecting to the database")
	w := flag.Int("w", 1, "The number of workers to use. A number greater than 1 will enable statements to be issued concurrently")
	d := flag.Bool("d", false, "Debug mode - turn this on to have errors printed to the terminal")
	flag.Parse()

	if *c == "" {
		return nil, errors.New("You must provide a connection string using the -c option")
	}
	pi, err := time.ParseDuration(*p)
	if err != nil {
		return nil, errors.New("You must provide a proper duration value with -p")
	}

	if *w <= 0 {
		return nil, errors.New("You must provide a worker count > 0 with -w")
	}

	return &config{connString: *c, pauseInterval: pi, workers: *w, debugMode: *d}, nil
}
