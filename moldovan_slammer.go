package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const randomChars string = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

type config struct {
	connString    string
	input         string
	pauseInterval time.Duration
	iterations    int
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

	keepRunning := true
	numLines := 0

	for keepRunning == true {
		// Build the line
		line, err := buildSQL(cfg.input)
		// Import it into sql here
		if err == nil {
			_, err = db.Exec(line)
			numLines++
			if err != nil {
				fmt.Println(err)
			}
		} else if err != nil {
			fmt.Println("Could not generate SQL: " + err.Error())
		}
		if numLines == cfg.iterations {
			keepRunning = false
		} else {
			time.Sleep(cfg.pauseInterval)
		}
	}
}

func newObjectCache() map[string]interface{} {
	return map[string]interface{}{"guid": make([]string, 0), "now": make([]string, 0)}
}

func buildSQL(inputTemplate string) (string, error) {
	// Supports:
	// {guid:ordinal}
	// {int:lower:upper}
	// {now:ordinal}
	// {float:lower:upper}
	// {char:num:case}
	objectCache := newObjectCache()
	var result bytes.Buffer
	var wordBuffer bytes.Buffer
	var foundWord = false
	for _, c := range inputTemplate {
		if c == '{' {
			// We're starting a word to parse
			foundWord = true
		} else if c == '}' {
			// We're closing a word, so eval it and get the data to put in the string
			foundWord = false
			parts := strings.Split(wordBuffer.String(), ":")
			val, err := resolveWord(objectCache, parts...)
			if err != nil {
				return "", err
			}
			result.WriteString(val)
			wordBuffer.Reset()
		} else if foundWord {
			// push it to the wordBuffer
			wordBuffer.WriteRune(c)
		} else {
			// Straight pass through
			result.WriteRune(c)
		}
	}

	return result.String(), nil
}

func resolveWord(objectCache map[string]interface{}, parts ...string) (string, error) {
	switch parts[0] {
	case "guid":
		return guid(objectCache, parts[1:]...)
	case "int":
		return integer(parts[1:]...)
	case "now":
		return now(objectCache, parts[1:]...)
	case "float":
		return float(parts[1:]...)
	case "char":
		return char(parts[1:]...)
	}
	return "", nil
}

func integer(opts ...string) (string, error) {
	lowerBound := 0
	upperBound := 100

	if len(opts) > 1 {
		nu, err := strconv.Atoi(opts[1])
		if err != nil {
			return "", nil
		}
		upperBound = nu
	}

	if len(opts) > 0 {
		nl, err := strconv.Atoi(opts[0])
		if err != nil {
			return "", nil
		}
		lowerBound = nl
	}

	if lowerBound > upperBound {
		return "", errors.New("You cannot generate a random number whose lower bound is greater than it's upper bound. Please check your input string")
	}

	// Incase we need to tell the function to invert the case
	negateResult := false
	// get the difference between them
	diff := upperBound - lowerBound
	// Since this supports negatives, need to handle some special corner cases?
	if lowerBound < 0 && upperBound <= 0 {
		// if the range is entirely negative
		negateResult = true
		// Swap them, so they are still the same relative distance from eachother, but positive - invert the result
		oldLower := lowerBound
		lowerBound = -upperBound
		upperBound = -oldLower
	}
	// neg to pos ranges currently not supported
	// else both are positive
	// get a number from 0 to diff
	n := rand.Intn(diff)
	// add lowerbound to it - now it's between lower and upper
	n += lowerBound
	if negateResult {
		n = -n
	}
	return strconv.Itoa(n), nil
}

func float(opts ...string) (string, error) {
	lowerBound := 0.0
	upperBound := 100.0

	if len(opts) > 1 {
		nu, err := strconv.ParseFloat(opts[1], 64)
		if err != nil {
			return "", nil
		}
		upperBound = nu
	}

	if len(opts) > 0 {
		nl, err := strconv.ParseFloat(opts[0], 64)
		if err != nil {
			return "", nil
		}

		lowerBound = nl
	}

	if lowerBound > upperBound {
		return "", errors.New("You cannot generate a random number whose lower bound is greater than it's upper bound. Please check your input string")
	}

	// Incase we need to tell the function to invert the case
	negateResult := false
	// get the difference between them
	diff := upperBound - lowerBound
	// Since this supports negatives, need to handle some special corner cases?
	if lowerBound < 0 && upperBound <= 0 {
		// if the range is entirely negative
		negateResult = true
		// Swap them, so they are still the same relative distance from eachother, but positive - invert the result
		oldLower := lowerBound
		lowerBound = -upperBound
		upperBound = -oldLower
	}
	// neg to pos ranges currently not supported
	// else both are positive
	// get a number from 0 to diff
	n := rand.Float64()*diff + lowerBound
	// add lowerbound to it - now it's between lower and upper
	n += lowerBound
	if negateResult {
		n = -n
	}
	return fmt.Sprintf("%f", n), nil
}

func char(opts ...string) (string, error) {
	charCase := "down"
	numChars := 2

	if len(opts) > 1 {
		charCase = opts[1]
	}
	if len(opts) > 0 {
		nc, err := strconv.Atoi(opts[0])
		if err != nil {
			return "", err
		}
		if nc <= 0 {
			return "", errors.New("You have specified a number of characters to generate which is not a number greater than zero. Please check your input string")
		}

		numChars = nc
	}

	result := make([]byte, numChars)
	for i := 0; i < numChars; i++ {
		result[i] = randomChars[rand.Intn(len(randomChars))]
	}
	if charCase == "up" {
		return strings.ToUpper(string(result)), nil
	}
	return string(result), nil
}

func now(objectCache map[string]interface{}, opts ...string) (string, error) {
	if len(opts) > 0 {
		// We want to re-use an existing guid
		ordinal, err := strconv.Atoi(opts[0])
		if err != nil {
			return "", err
		}
		c, _ := objectCache["now"]
		cache := c.([]string)
		if len(cache) < ordinal {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for time-now. Please check your input string", ordinal)
		}
		return cache[ordinal], nil
	}
	now := time.Now().Format(time.RFC3339)

	// store it in the cache
	c, _ := objectCache["now"]
	cache := c.([]string)
	objectCache["now"] = append(cache, now)

	return now, nil

}

func guid(objectCache map[string]interface{}, opts ...string) (string, error) {
	if len(opts) > 0 {
		// We want to re-use an existing guid
		ordinal, err := strconv.Atoi(opts[0])
		if err != nil {
			return "", err
		}
		c, _ := objectCache["guid"]
		cache := c.([]string)
		if len(cache) < ordinal {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for guids. Please check your input string", ordinal)
		}
		return cache[ordinal], nil
		// There is apparently no fantastic way to generate guids / uuids in go?
		// There are some libraries, but apparently there is no standard, correct implementation.
		// Thus, I do this instead:
	}
	out, err := exec.Command("uuidgen").Output()
	if err != nil {
		return "", err
	}
	guid := strings.Trim(string(out), "\n")

	// store it in the cache
	c, _ := objectCache["guid"]
	cache := c.([]string)
	objectCache["guid"] = append(cache, guid)

	return guid, nil

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
	iterations := -1
	i := os.Getenv("MS_ITERATIONS")
	if iterations, err = strconv.Atoi(i); err != nil {
		return nil, errors.New("MS_ITERATIONS must be a valid integer (-1 for unlimited)")
	}
	cfg := &config{
		connString:    os.Getenv("MS_CONNSTRING"),
		input:         os.Getenv("MS_INPUT"),
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
