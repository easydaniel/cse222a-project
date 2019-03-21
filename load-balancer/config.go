package main

import (
	"encoding/json"
	"io/ioutil"
)

type LBConfig struct {
	Servers []string
	Weights []int
	Port    int
	Listen  int
}

func loadLBConfig(filename string) (*LBConfig, error) {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := LBConfig{}
	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
