package config

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	fmt.Println(GetConfigAbsolutePath("application.yml"))
}
