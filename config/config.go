package config

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"io"
)

var configs map[string]string

const configFile = "./config/config.properties"

func init() {
	configs = make(map[string]string)

	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		if len(line) <= 0 && (err != nil || err == io.EOF) {
			break
		}

		//process the config item
		if strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.Trim(line, " \n")
		index := strings.IndexByte(line, '=')
		if len(line) < 3 || index == -1 {
			continue
		}

		configs[line[:index]] = line[index+1:]
	}
}

func GetAsString(key string) string {
	return configs[key]
}

func GetAsInt64(key string) int64 {
	val, _ := configs[key]
	r, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return 0
	}
	return r
}
