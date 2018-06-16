package vkcli_test

import (
	"testing"
	"time"
	"github.com/SergeyMokhov/vkcli"
	"os"
	"path/filepath"
)

var dataFolder = "./data"

func TestGenerate(t *testing.T) {
	host := "localhost,127.0.0.1"
	validFrom := ""
	validFor := 365 * 24 * time.Hour
	isCA := true
	rsaBits := 2048
	ecdsaCurve := ""
	expectedCertPath := filepath.Join(dataFolder, "cert.pem")
	expectedKeyPath := filepath.Join(dataFolder, "key.pem")

	cleanup()

	vkcli.GenerateCert(host, validFrom, validFor, isCA, rsaBits, ecdsaCurve, dataFolder)

	certExist := fileExists(expectedCertPath)
	if !certExist {
		t.Errorf("Certificate was not created in %s", expectedCertPath)
	}

	keyExist := fileExists(expectedKeyPath)
	if !keyExist {
		t.Errorf("Key was not created in %s", expectedKeyPath)
	}
}

func cleanup() {
	os.RemoveAll(dataFolder)
	os.MkdirAll(dataFolder, os.ModePerm)
}

func fileExists(path string) (bool) {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
