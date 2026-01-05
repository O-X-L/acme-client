package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"

	"git.oxl.at/acme-client/internal/acme"
	"git.oxl.at/acme-client/internal/u"
	oxl_validate_regex "git.oxl.at/go-validator/pkg/validate/regex"
)

func ValidateConfig(cnf *ConfigFile) error {
	Config.PathCerts = strings.TrimSpace(Config.PathCerts)
	Config.PathWeb = strings.TrimSpace(Config.PathWeb)

	if Config.PathWeb != "" {
		pathWebAcmeChallenge := filepath.Join(Config.PathWeb, DIR_WEB_ACME_CHALLENGE)
		if _, err := os.Stat(pathWebAcmeChallenge); err != nil {
			return fmt.Errorf("ACME-challenge directory does not exist: %s", pathWebAcmeChallenge)
		}
	}

	if cnf.Retries >= 4 {
		return fmt.Errorf("retries must be < 4")
	}

	if cnf.CooldownSec < 1 {
		return fmt.Errorf("cooldown_sec must be >= 1")
	}

	if cnf.FileGroup != "" {
		if _, err := user.LookupGroup(cnf.FileGroup); err != nil {
			return fmt.Errorf("group does not exist")
		}
	}

	if cnf.FileModeCert < os.FileMode(0600) {
		return fmt.Errorf("invalid file_mode_cert - must be >= 600 (writable for service-user)")
	}
	if cnf.FileModeKey < os.FileMode(0600) {
		return fmt.Errorf("invalid file_mode_key - must be >= 600 (writable for service-user)")
	}
	if cnf.FileModeKey > os.FileMode(0640) {
		u.LogWarning("The file_mode_key is too permissive! It's recommended to change it to 0640 or lower!")
	}

	appIDs := []uint{}
	for _, app := range cnf.Apps {
		if slices.Contains(appIDs, app.ID) {
			return fmt.Errorf("got duplicate App-ID: '%s' (%d)", app.Name, app.ID)
		}
		appIDs = append(appIDs, app.ID)

		certIDs := []uint{}
		for _, cert := range app.Certs {
			if slices.Contains(certIDs, cert.ID) {
				return fmt.Errorf("got duplicate Certificate-ID: app '%s' (%d - %d)", app.Name, app.ID, cert.ID)
			}
			certIDs = append(certIDs, cert.ID)
		}

	}

	return nil
}

func ValidateCertConfig(app App, cert AppCert) error {
	switch cert.ChallengeType {
	case CHALLENGE_TYPE_DNS:
		if !acme.IsSupportedProvider(cert.Provider) {
			return fmt.Errorf("unsupported dns-01 provider: %s", cert.Provider)
		}

	case CHALLENGE_TYPE_HTTP:
		if Config.PathWeb == "" {
			return fmt.Errorf("webroot path required for http-01 (/.well-known/acme-challenge)")
		}
		if _, err := os.Stat(Config.PathWeb); os.IsNotExist(err) {
			return fmt.Errorf("webroot dir for http-01 does not exist: %s (/.well-known/acme-challenge)", Config.PathWeb)
		}
		if !u.RegexMatch(oxl_validate_regex.REGEX_URL_SIMPLE, cert.Provider) {
			return fmt.Errorf("provider for http-01 should be a valid URL: %s", cert.Provider)
		}

	default:
		return fmt.Errorf("unsupported challenge-type: %s", cert.ChallengeType)
	}

	return nil
}
