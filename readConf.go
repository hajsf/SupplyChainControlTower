package main

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func readConf(filename string) (*conf, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &conf{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
