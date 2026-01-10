package config

import (
	"log"
	"os"
	"time"

	"git.oxl.at/acme-client/internal/u"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

var (
	Config                 *ConfigFile
	PathCertsPublic        string
	PathCertsPrivate       string
	PathAccountCache       string
	PathCertsBundlePublic  string
	PathCertsBundlePrivate string
	RenewalDays            time.Duration
	CheckMode              bool
)

const (
	CHALLENGE_TYPE_DNS     = "dns-01"
	CHALLENGE_TYPE_HTTP    = "http-01"
	DIR_BUNDLE_PUBLIC      = "bundle_certs"
	DIR_BUNDLE_PRIVATE     = "bundle_private"
	DIR_ACCOUNT            = "account"
	DIR_CERTS_PUBLIC       = "certs"
	DIR_CERTS_PRIVATE      = "private"
	DIR_WEB_ACME_CHALLENGE = ".well-known/acme-challenge"
	VERSION                = "1.1.0"
)

type Certificate struct {
	ID             uint              `yaml:"id" required:"true"`
	Provider       string            `yaml:"provider" required:"true"` // URL if HTTP-01 else one of the listed DNS-providers
	ProviderConfig map[string]string `yaml:"provider_config"`          // env-vars for DNS-01
	ChallengeType  string            `yaml:"challenge_type" validate_regex:"^(http-01|dns-01)$" required:"true"`
	Domains        []string          `yaml:"domains" validate:"domain_wildcard"`
}

type Group struct {
	Name  string        `yaml:"name" required:"true"`
	ID    uint          `yaml:"id" required:"true"`
	Certs []Certificate `yaml:"certs"`
}

type ConfigFile struct {
	Email        string      `yaml:"email" validate:"email" required:"true"`
	Groups       []Group     `yaml:"groups"`
	Retries      uint        `yaml:"retries" default:"0"`
	CooldownSec  uint        `yaml:"cooldown_sec" default:"2"`        // do not overwhelm the ACME service with requests - speed is not that important for requesting certs
	PathWeb      string      `yaml:"path_web" validate:"path_simple"` // only required if certs use http-01
	PathCerts    string      `yaml:"path_certs" required:"true" validate:"path_simple"`
	CreateBundle *bool       `yaml:"create_bundle" default:"false"`
	FileModeCert os.FileMode `yaml:"file_mode_cert" default:"0644"`
	FileModeKey  os.FileMode `yaml:"file_mode_key" default:"0600"`
	FileGroup    string      `yaml:"file_group"`
	HookCmd      string      `yaml:"hook_cmd"`
	RenewalDays  uint        `yaml:"renewal_days" default:"14"`
	MaxDomains   uint        `yaml:"max_domains" default:"50"`
}

func LoadConfig(path string) (*ConfigFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cnf ConfigFile
	if err := yaml.Unmarshal(data, &cnf); err != nil {
		return nil, err
	}

	if err := defaults.Set(&cnf); err != nil {
		u.LogErrorf("Failed to set config defaults: \"%v\"", err)
		log.Fatalln("failed to load config")
	}

	return &cnf, nil
}
