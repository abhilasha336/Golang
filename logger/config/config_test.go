package config_test

import (
	"os"
	"testing"
)

func checkFileExistsWithOpen(filename string) error {
	_, err := os.Open(filename)
	return err
}

func checkFileExistsWithStat(filename string) error {
	_, err := os.Stat(filename)
	return err
}
func BenchmarkFileOpen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkFileExistsWithOpen(".env")
	}
}

func BenchmarkFileState(b *testing.B) {
	for i := 0; i < b.N; i++ {
		checkFileExistsWithStat(".env")
	}
}
