package config

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"git.oxl.at/acme-client/internal/acme"
	"git.oxl.at/acme-client/internal/u"
	oxl_validate "git.oxl.at/go-validator/pkg/validate"
	oxl_validate_regex "git.oxl.at/go-validator/pkg/validate/regex"
)

func ValidateConfig(cnf *ConfigFile) error {
	Config.PathCerts = strings.TrimSpace(Config.PathCerts)
	Config.PathWeb = strings.TrimSpace(Config.PathWeb)

	if Config.PathWeb != "" {
		pathWebAcmeChallenge := filepath.Join(Config.PathWeb, DIR_WEB_ACME_CHALLENGE)
		if _, err := os.Stat(pathWebAcmeChallenge); err != nil {
			err := fmt.Errorf("ACME-challenge directory does not exist: %s", pathWebAcmeChallenge)
			if CheckMode {
				u.LogWarning(fmt.Sprintf("%v", err))
			} else {
				return err
			}
		}
	}

	if cnf.Retries >= 4 {
		return fmt.Errorf("retries must be < 4")
	}

	if cnf.CooldownSec < 1 {
		return fmt.Errorf("cooldown_sec must be >= 1")
	}

	if cnf.FileGroup != "" {
		if err := validateGroup(cnf.FileGroup); err != nil {
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

	groupIDs := []uint{}
	for _, grp := range cnf.Groups {
		if slices.Contains(groupIDs, grp.ID) {
			return fmt.Errorf("got duplicate Group-ID: '%s' (%d)", grp.Name, grp.ID)
		}
		groupIDs = append(groupIDs, grp.ID)

		certIDs := []uint{}
		for _, cert := range grp.Certs {
			if slices.Contains(certIDs, cert.ID) {
				return fmt.Errorf("got duplicate Certificate-ID: group '%s' (%d - %d)", grp.Name, grp.ID, cert.ID)
			}
			certIDs = append(certIDs, cert.ID)
			if err := validateCertConfig(cert); err != nil {
				return fmt.Errorf("got invalid certificate config: group '%s' (%d - %d) %v", grp.Name, grp.ID, cert.ID, err)
			}
		}
	}

	return nil
}

func validateGroup(grp string) error {
	_, err := u.GetGroupID(grp)
	if err != nil {
		return err
	}
	return nil
}

func validateCertConfig(cert GroupCert) error {
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
			errMsg := fmt.Errorf("webroot dir for http-01 does not exist: %s (/.well-known/acme-challenge)", Config.PathWeb)
			if CheckMode {
				u.LogWarning(fmt.Sprintf("%v", errMsg))
			} else {
				return errMsg
			}
		}
		if !u.RegexMatch(oxl_validate_regex.REGEX_URL_SIMPLE, cert.Provider) {
			return fmt.Errorf("provider for http-01 should be a valid URL: %s", cert.Provider)
		}

	default:
		return fmt.Errorf("unsupported challenge-type: %s", cert.ChallengeType)
	}

	return nil
}

func validateDomainWildcard(value interface{}) bool {
	s, ok := value.(string)
	if !ok {
		return false
	}
	s = strings.TrimPrefix(s, "*.")
	return value == "localhost" || oxl_validate_regex.ValidateRegex(oxl_validate_regex.REGEX_DOMAINS_SIMPLE, s)
}

func ValidateSchema(cnf *ConfigFile) bool {
	v := &oxl_validate.StructValidator{}
	v.ValidatorsCustom = oxl_validate.GetDefaultCustomValidators()
	v.ValidatorsCustom["domain_wildcard"] = validateDomainWildcard

	validationErrors := v.Validate(cnf)
	if len(validationErrors) > 0 {
		u.LogError(fmt.Sprintf("Got invalid config (schema): %+v", validationErrors))
		return false
	}
	return true
}
