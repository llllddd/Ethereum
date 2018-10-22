package config

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	err := ParseConfig("conf.json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Cfg.Lock", Cfg.Lock)

}
