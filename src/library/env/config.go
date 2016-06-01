package env

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

//data struct to hold a group of config parameters, with key type default as string and value store in string format
type Config struct {
	KV map[string]string
	sync.RWMutex
}

//reg-expression to resolve key value in config file, see test.conf, example:
//	key = value
//line with # will be considered as comments and be ignored
var kvRegexp *regexp.Regexp = regexp.MustCompile(`[\t ]*([0-9A-Za-z_]+)[\t ]*=[\t ]*([^\t\n\f\r# ]+)[\t #]*`)

func NewConfig(path string) *Config {
	c := &Config{}
	if path != "" {
		c.LoadFromFile(path)
	} else {
		c.KV = make(map[string]string, 0)
	}

	return c
}

//all value will be set and store in string formation, ValTYPE method will get the value with TYPE formation
//value of boolean type can be set as 0, 1 or false, true
func (c *Config) Set(key string, val interface{}) {
	value := ""
	b, ok := val.(bool)
	if ok {
		if b {
			value = "true"
		} else {
			value = "false"
		}
	} else {
		value = fmt.Sprint(val)
	}

	c.Lock()
	c.KV[key] = value
	c.Unlock()
}

func (c *Config) ValString(key string) (string, error) {
	var err error = nil
	c.RLock()
	v, ok := c.KV[key]
	c.RUnlock()
	if !ok {
		err = errors.New("Key not exist")
	}

	return v, err
}

func (c *Config) ValBool(key string) (bool, error) {
	v, err := c.ValString(key)
	if err != nil {
		return false, err
	}

	lower := strings.ToLower(v)
	switch lower {
	case "true", "1":
		return true, nil
	case "false", "0":
		return false, nil
	default:
		return false, errors.New("Not boolean value")
	}
}

func (c *Config) ValInt64(key string) (int64, error) {
	v, err := c.ValString(key)
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseInt(v, 10, 64)

	return value, err
}

func (c *Config) ValFloat64(key string) (float64, error) {
	v, err := c.ValString(key)
	if err != nil {
		return float64(0), err
	}

	value, err := strconv.ParseFloat(v, 64)

	return value, err
}

//read all the config key-value pair from config file defined by path
func (c *Config) LoadFromFile(path string) error {
	c.Lock()
	c.KV = make(map[string]string)
	c.Unlock()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	c.Lock()
	for _, line := range strings.Split(string(content), "\n") {
		line := strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		slice := kvRegexp.FindStringSubmatch(line)
		if slice != nil {
			c.KV[slice[1]] = slice[2]
		}
	}
	c.Unlock()

	return nil
}
