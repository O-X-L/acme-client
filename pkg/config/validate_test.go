package config

import (
	"os"
	"path/filepath"
	"testing"

	"git.oxl.at/acme-client/internal/u"
	"github.com/go-acme/lego/v4/lego"
)

const (
	DUMMY_URL = "https://acme.example.net"
	TYPE_HTTP = "http-01"
	TYPE_DNS  = "dns-01"
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
			name: "Valid global config",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      tmpWebDir,
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
							{ID: 2, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
					{
						Name: "test1",
						ID:   2,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
							{ID: 2, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid high retries",
			cfg: &ConfigFile{
				Retries:      10,
				PathCerts:    "/tmp",
				PathWeb:      tmpWebDir,
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid file group",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      tmpWebDir,
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				FileGroup:    "does-not-exist",
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "No cooldown-sec",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      tmpWebDir,
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  0,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Unwritable Cert Mode",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      tmpWebDir,
				FileModeCert: 0400,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Duplicate Group-ID",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      tmpWebDir,
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Duplicate Cert-ID",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      tmpWebDir,
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Non-existant web-dir",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      "/tmp/does/not/exist",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid global config with dns-01",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_DNS, Provider: "cloudflare", Domains: []string{"waf.alpenmesh.com"}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid global config with wildcard domain",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_DNS, Provider: "cloudflare", Domains: []string{"*.waf.alpenmesh.com"}},
						},
					},
				},
			},
			wantErr: false,
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
		cert    GroupCert
		wantErr bool
	}{
		{
			name:    "Invalid challenge type",
			cert:    GroupCert{ChallengeType: "dns-05"},
			wantErr: true,
		},
		{
			name:    "Unsupported challenge type",
			cert:    GroupCert{ChallengeType: "tls-alpn-01"},
			wantErr: true,
		},
		{
			name:    "Unsupported DNS provider",
			cert:    GroupCert{ChallengeType: CHALLENGE_TYPE_DNS, Provider: "unknown-cloud"},
			wantErr: true,
		},
		{
			name:    "HTTP-01 with non-URL provider",
			cert:    GroupCert{ChallengeType: CHALLENGE_TYPE_HTTP, Provider: "just-a-string"},
			wantErr: true,
		},
		{
			name:    "Valid HTTP-01",
			cert:    GroupCert{ChallengeType: CHALLENGE_TYPE_HTTP, Provider: lego.LEDirectoryStaging},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCertConfig(tt.cert)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: got error %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestConfigValidateGroup(t *testing.T) {
	tests := []struct {
		name    string
		group   string
		wantErr bool
	}{
		{
			name:    "Valid group name",
			group:   "root",
			wantErr: false,
		},
		{
			name:    "Valid GID",
			group:   "1001",
			wantErr: false,
		},
		{
			name:    "Non-existent group name",
			group:   "non-exist",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateGroup(tt.group)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s: got error %v, wantErr %v", tt.name, err, tt.wantErr)
			}
		})
	}
}

func TestConfigValidateSchema(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *ConfigFile
		wantErr bool
	}{
		{
			name: "Valid global config",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				PathWeb:      "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Email:        "test@alpenmesh.com",
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
							{ID: 2, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
					{
						Name: "test1",
						ID:   2,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
							{ID: 2, ChallengeType: TYPE_HTTP, Provider: DUMMY_URL},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid global config with dns-01",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Email:        "test@alpenmesh.com",
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_DNS, Provider: "cloudflare", Domains: []string{"waf.alpenmesh.com"}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid global config with wildcard domain",
			cfg: &ConfigFile{
				Retries:      2,
				PathCerts:    "/tmp",
				FileModeCert: 0644,
				FileModeKey:  0600,
				CooldownSec:  1,
				Email:        "test@alpenmesh.com",
				Groups: []Group{
					{
						Name: "test1",
						ID:   1,
						Certs: []GroupCert{
							{ID: 1, ChallengeType: TYPE_DNS, Provider: "cloudflare", Domains: []string{"*.waf.alpenmesh.com"}},
						},
					},
				},
			},
			wantErr: false,
		},
		// todo: extend
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Config = tt.cfg
			isValid := ValidateSchema(tt.cfg)
			if (!isValid) != tt.wantErr {
				t.Errorf("%s: wantErr %v", tt.name, tt.wantErr)
			}
		})
	}
}
