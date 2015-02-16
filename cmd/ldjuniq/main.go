package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

const Version = "0.1.0"

// StringValue returns the value for a given key in dot notation
func StringValue(key string, doc map[string]interface{}) (string, error) {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return "", fmt.Errorf("keys exhausted")
	}
	head := keys[0]
	val, ok := doc[head]
	if !ok {
		return "", fmt.Errorf("key %s not found", head)
	}
	switch t := val.(type) {
	case string:
		return val.(string), nil
	case map[string]interface{}:
		if len(keys) < 2 {
			return "", fmt.Errorf("no value found")
		}
		return StringValue(keys[1], val.(map[string]interface{}))
	case float64:
		return strconv.FormatFloat(val.(float64), 'f', 6, 64), nil
	case int:
		return strconv.Itoa(val.(int)), nil
	case []interface{}:
		return "", fmt.Errorf("not supported yet: %+v\n", t)
	default:
		return "", fmt.Errorf("unknown type: %+v, %+v\n", t, reflect.TypeOf(t))
	}
	return "", fmt.Errorf("no value found")
}

func main() {
	key := flag.String("key", "", "key in dot notation")
	version := flag.Bool("v", false, "show version and exit")

	flag.Parse()

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if *key == "" {
		log.Fatal("key is required")
	}

	var reader *bufio.Reader
	if flag.NArg() < 1 {
		reader = bufio.NewReader(os.Stdin)
	} else {
		file, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		reader = bufio.NewReader(file)
	}

	seen := make(map[string]struct{})

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		doc := make(map[string]interface{})
		err = json.Unmarshal([]byte(line), &doc)
		if err != nil {
			log.Fatal(err)
		}
		val, err := StringValue(*key, doc)
		if err != nil {
			log.Fatal(err)
		}
		_, ok := seen[val]
		if ok {
			continue
		}
		b, err := json.Marshal(doc)
		if err != nil {
			log.Fatal(err)
		}
		seen[val] = struct{}{}
		fmt.Println(string(b))
	}
}
