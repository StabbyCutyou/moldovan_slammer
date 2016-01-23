package slammer

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/StabbyCutyou/moldovan_slammer/moldova"
)

type config struct {
	connString    string
	input         string
	pauseInterval time.Duration
	iterations    int
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

	for i := 0; i < cfg.iterations; i++ {
		// Build the line
		line, err := moldova.ParseTemplate(cfg.input)
		// Import it into sql here
		if err == nil {
			_, err2 := db.Exec(line)
			if err2 != nil {
				fmt.Println(err2)
			}
		} else if err != nil {
			fmt.Println("Could not generate SQL: " + err.Error())
		}
		time.Sleep(cfg.pauseInterval)
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
	iterations := -1
	i := os.Getenv("MS_ITERATIONS")
	if iterations, err = strconv.Atoi(i); err != nil {
		return nil, errors.New("MS_ITERATIONS must be a valid integer (-1 for unlimited)")
	}

	cfg := &config{
		connString:    os.Getenv("MS_CONNSTRING"),
		pauseInterval: d,
		iterations:    iterations,
	}

	if cfg.connString == "" {
		return nil, errors.New("MS_CONNSTRING must be set")
	}

	if cfg.input == "" {
		return nil, errors.New("MS_INPUT must be set")
	}

	return cfg, nil
}
