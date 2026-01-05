package config

import (
	"os"
	"path/filepath"
	"testing"

	"git.oxl.at/acme-client/internal/u"
	"github.com/go-acme/lego/v4/lego"
)

func setupTestDir(t *testing.T) string {
	tmpDir, _ := os.MkdirTemp("", "manager_test")
	Config = &ConfigFile{
		PathCerts:    tmpDir,
		CreateBundle: u.PtrBool(false),
	}
	return tmpDir
}

func TestConfigValidate(t *testing.T) {
	tmpWebDir := setupTestDir(t)
	os.MkdirAll(filepath.Join(tmpWebDir, DIR_WEB_ACME_CHALLENGE), 0755)
	defer os.RemoveAll(tmpWebDir)

	tests := []struct {
		name    string
		cfg     *ConfigFile
		wantErr bool
	}{
		{
			name:    "Invalid high retries",
			cfg:     &ConfigFile{Retries: 5, FileModeCert: 0644, FileModeKey: 0600, CooldownSec: 1},
			wantErr: true,
		},
		{
			name:    "Invalid file group",
			cfg:     &ConfigFile{Retries: 1, FileModeCert: 0644, FileModeKey: 0600, FileGroup: "non-existent-group-xyz", CooldownSec: 1},
			wantErr: true,
		},
		{
			name: "Valid global config",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Apps: []App{
					{Name: "test1", ID: 1, Certs: []AppCert{{ID: 1}, {ID: 2}}},
					{Name: "test2", ID: 2, Certs: []AppCert{{ID: 1}, {ID: 2}}},
				},
			},
			wantErr: false,
		},
		{
			name: "No cooldown-sec",
			cfg: &ConfigFile{
				PathCerts:    "/tmp",
				FileModeCert: 0600,
				FileModeKey:  0600,
				Retries:      1,
				CooldownSec:  0,
			},
			wantErr: true,
		},
		{
			name: "Unwritable Cert Mode",
			cfg: &ConfigFile{
				PathCerts:    "/tmp",
				FileModeCert: 0400,
				FileModeKey:  0600,
				Retries:      1,
				CooldownSec:  1,
			},
			wantErr: true,
		},
		{
			name: "Duplicate App-ID",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Apps: []App{
					{Name: "test1", ID: 1},
					{Name: "test2", ID: 1},
				},
			},
			wantErr: true,
		},
		{
			name: "Duplicate Cert-ID",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Apps: []App{
					{Name: "test1", ID: 1, Certs: []AppCert{{ID: 1}, {ID: 2}}},
					{Name: "test2", ID: 2, Certs: []AppCert{{ID: 1}, {ID: 1}}},
				},
			},
			wantErr: true,
		},
		{
			name: "Non-existant web-dir",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      "/tmp/non-existant",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Apps: []App{
					{Name: "test1", ID: 1, Certs: []AppCert{{ID: 1}, {ID: 2}}},
					{Name: "test2", ID: 2, Certs: []AppCert{{ID: 1}, {ID: 1}}},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Config = tt.cfg
			err := ValidateConfig(tt.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: got error %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestConfigValidateCertLogic(t *testing.T) {
	Config = &ConfigFile{PathWeb: "/tmp"}

	tests := []struct {
		name    string
		cert    AppCert
		wantErr bool
	}{
		{
			name:    "Invalid challenge type",
			cert:    AppCert{ChallengeType: "dns-05"},
			wantErr: true,
		},
		{
			name:    "Unsupported challenge type",
			cert:    AppCert{ChallengeType: "tls-alpn-01"},
			wantErr: true,
		},
		{
			name:    "Unsupported DNS provider",
			cert:    AppCert{ChallengeType: CHALLENGE_TYPE_DNS, Provider: "unknown-cloud"},
			wantErr: true,
		},
		{
			name:    "HTTP-01 with non-URL provider",
			cert:    AppCert{ChallengeType: CHALLENGE_TYPE_HTTP, Provider: "just-a-string"},
			wantErr: true,
		},
		{
			name:    "Valid HTTP-01",
			cert:    AppCert{ChallengeType: CHALLENGE_TYPE_HTTP, Provider: lego.LEDirectoryStaging},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCertConfig(App{}, tt.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: got error %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}
