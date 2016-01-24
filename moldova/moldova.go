package moldova

import (
	"bytes"
	crand "crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func newObjectCache() map[string]interface{} {
	return map[string]interface{}{"guid": make([]string, 0), "now": make([]string, 0), "country": make([]string, 0)}
}

// ParseTemplate will take an input string of text, and replace any recongized
// tokens with a random value that is determined for each type of token
func ParseTemplate(inputTemplate string) (string, error) {
	// Supports:
	// {guid:ordinal}
	// {int:lower:upper}
	// {now:ordinal}
	// {float:lower:upper}
	// {char:num:case}
	// {country:case:ordinal}
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

// This function was borrowed with permission from the following location
// https://github.com/dgryski/trifles/blob/master/uuid/uuid.go
// All credit / lawsuits can be forwarded to Damian Gryski and Russ Cox
func uuidv4() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(crand.Reader, b)
	if err != nil {
		// probably "shouldn't happen"
		log.Fatal(err)
	}
	b[6] = (b[6] & 0x0F) | 0x40
	b[8] = (b[8] &^ 0x40) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[:4], b[4:6], b[6:8], b[8:10], b[10:])
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
	case "country":
		return country(objectCache, parts[1:]...)
	}
	return "", nil
}

// TODO All the below functions need way better commenting and parameter annotations
// It's described in the readme, but I should probably make these public and then
// give them proper comments, so that GoDoc can also document them

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
	if lowerBound < 0.0 && upperBound <= 0.0 {
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
	n := (rand.Float64() * diff) + lowerBound

	if negateResult {
		n = -n
	}
	return fmt.Sprintf("%f", n), nil
}

func country(objectCache map[string]interface{}, opts ...string) (string, error) {
	charCase := "up"

	if len(opts) > 0 {
		charCase = opts[0]
	}

	if len(opts) > 1 {
		// We want to re-use an existing country
		ordinal, err := strconv.Atoi(opts[1])
		if err != nil {
			return "", err
		}
		c, _ := objectCache["country"]
		cache := c.([]string)
		if len(cache) < ordinal {
			return "", fmt.Errorf("Ordinal %d has not yet been encountered for countries. Please check your input string", ordinal)
		}
		country := cache[ordinal]
		// Countries go into the cache upper case, only check for lowering it
		if charCase == "down" {
			return strings.ToLower(country), nil
		}
		return country, nil
	}
	// Generate a new one
	n := rand.Intn(len(CountryCodes))
	country := CountryCodes[n]
	// store it in the cache
	c, _ := objectCache["country"]
	cache := c.([]string)
	objectCache["country"] = append(cache, country)

	if charCase == "down" {
		return strings.ToLower(country), nil
	}

	return country, nil

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

	result := generateRandomString(numChars)

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
	}

	guid := uuidv4()
	// store it in the cache
	c, _ := objectCache["guid"]
	cache := c.([]string)
	objectCache["guid"] = append(cache, guid)

	return guid, nil

}
