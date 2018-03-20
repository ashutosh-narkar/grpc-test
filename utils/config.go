package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Config struct {
	Listen   string    `json:"listen"`
	Verbose  bool      `json:"verbose"`
	Backends []Backend `json:"backends"`
}

type Backend struct {
	Filter      string `json:"filter"`
	Backend     string `json:"backend"`
	BackendName string `json:"backendName"`
}

func GetConfiguration(file string) Config {
	raw, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var config Config
	if err := json.Unmarshal(raw, &config); err != nil {
		panic(err)
	}

	fmt.Printf("Proxy configuration read from file %q \n%s\n", file, prettyPrintConfig(config))
	return config
}

func prettyPrintConfig(conf interface{}) string {
	bytes, err := json.MarshalIndent(conf, "", "   ")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	return string(bytes)
}
