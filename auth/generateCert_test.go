package auth

import (
	"gitlab.com/g00g/vk-cli/tools"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	host := "localhost,127.0.0.1"
	validFrom := ""
	validFor := 365 * 24 * time.Hour
	isCA := true
	rsaBits := 2048
	ecdsaCurve := ""
	expectedCertPath := filepath.Join(dataFolder, "cert.pem")
	expectedKeyPath := filepath.Join(dataFolder, "key.pem")

	defer cleanup()

	err := GenerateCert(host, validFrom, validFor, isCA, rsaBits, ecdsaCurve, dataFolder)
	if err != nil {
		t.Fatalf("Certificate generation returned error %s", err)
	}
	certExist := tools.FileExists(expectedCertPath)
	if !certExist {
		t.Errorf("Certificate was not created in %s", expectedCertPath)
	}

	keyExist := tools.FileExists(expectedKeyPath)
	if !keyExist {
		t.Errorf("Key was not created in %s", expectedKeyPath)
	}
}

func cleanup() {
	os.RemoveAll(dataFolder)
}
