package manager

import (
	"os"
	"path/filepath"
	"testing"

	"git.oxl.at/acme-client/internal/config"
	"github.com/go-acme/lego/v4/lego"
)

func TestGetOrCreateACMEUser(t *testing.T) {
	tmpDir, _ := os.MkdirTemp("", "account_test")
	defer os.RemoveAll(tmpDir)

	config.Config = &config.ConfigFile{
		PathCerts: tmpDir,
		Email:     "test1@waf.alpenmesh.com",
	}
	config.PathAccountCache = filepath.Join(tmpDir, config.DIR_ACCOUNT)
	os.MkdirAll(config.PathAccountCache, 0755)

	stagingURL := lego.LEDirectoryStaging

	// create user 1
	user1, err := getOrCreateACMEUser(stagingURL)
	if err != nil {
		t.Fatalf("Failed to create user 1: %v", err)
	}

	// user 1 from cache
	user1cached, err := getOrCreateACMEUser(stagingURL)
	if err != nil {
		t.Fatalf("Failed to cache user 1: %v", err)
	}

	if user1.Email != user1cached.Email {
		t.Errorf("Failed to cache user 1 - email mismatch: %s vs %s", user1.Email, user1cached.Email)
	}

	// ensure only one account was cached
	files, _ := os.ReadDir(filepath.Join(tmpDir, config.DIR_ACCOUNT))
	if len(files) != 1 {
		t.Errorf("Expected 1 account file, found %d", len(files))
	}

	// create user 2
	config.Config = &config.ConfigFile{
		PathCerts: tmpDir,
		Email:     "test2@waf.alpenmesh.com",
	}

	user2, err := getOrCreateACMEUser(stagingURL)
	if err != nil {
		t.Fatalf("Failed to create user 2: %v", err)
	}

	// user 2 from cache
	user2cached, err := getOrCreateACMEUser(stagingURL)
	if err != nil {
		t.Fatalf("Failed to cache user 2: %v", err)
	}

	if user2.Email != user2cached.Email {
		t.Errorf("Failed to cache user 2 - email mismatch: %s vs %s", user2.Email, user2cached.Email)
	}

	// ensure two accounts were cached
	files, _ = os.ReadDir(filepath.Join(tmpDir, config.DIR_ACCOUNT))
	if len(files) != 2 {
		t.Errorf("Expected 2 account files after switching providers, found %d", len(files))
	}
}
