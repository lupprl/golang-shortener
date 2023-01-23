package main

//go:generate go test -c -o=../../bin/shortenertest

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Основной тест, запускает все остальные тесты
	os.Exit(m.Run())
}
