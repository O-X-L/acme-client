package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"git.oxl.at/acme-client/pkg/config"
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

func TestProcessDomainsInBatches(t *testing.T) {
	config.Config.MaxDomains = 2
	dummyCallback := func(grp config.Group, cert config.Certificate, id int, input []string) (bool, error) {
		if slices.Contains(input, "eee") {
			return false, fmt.Errorf("y")
		}
		if slices.Contains(input, "ccc") {
			return true, nil
		}
		return false, nil
	}
	dummyGrp := config.Group{}
	dummyCert := config.Certificate{}

	tests := []struct {
		name          string
		input         []string
		wantChanged   bool
		wantErr       bool
		wantProcessed int
	}{
		{
			name:          "No change - single",
			input:         []string{"aaa", "bbb"},
			wantChanged:   false,
			wantErr:       false,
			wantProcessed: 1,
		},
		{
			name:          "No change - multi",
			input:         []string{"aaa", "bbb", "ddd", "hhh", "jjj"},
			wantChanged:   false,
			wantErr:       false,
			wantProcessed: 3,
		},
		{
			name:          "Error on first",
			input:         []string{"aaa", "eee"},
			wantChanged:   false,
			wantErr:       true,
			wantProcessed: 1,
		},
		{
			name:          "Error on second",
			input:         []string{"aaa", "bbb", "eee"},
			wantChanged:   false,
			wantErr:       true,
			wantProcessed: 2,
		},
		{
			name:          "Changed",
			input:         []string{"aaa", "bbb", "ccc", "ddd", "hhh"},
			wantChanged:   true,
			wantErr:       false,
			wantProcessed: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changed, procssedBatches, err := processDomainsInBatches(dummyGrp, dummyCert, tt.input, dummyCallback)

			hasErr := err != nil
			if changed != tt.wantChanged || hasErr != tt.wantErr || procssedBatches != tt.wantProcessed {
				t.Errorf(
					"%s failed: reason: changed %v != %v, processedBatches %d != %d, error %v != %v ('%v')",
					tt.name, changed, tt.wantChanged, procssedBatches, tt.wantProcessed, hasErr, tt.wantErr, err,
				)
			}
		})
	}
}
