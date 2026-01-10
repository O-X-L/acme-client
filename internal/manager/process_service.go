package manager

import (
	"fmt"
	"strings"
	"time"

	"git.oxl.at/acme-client/internal/u"
	"git.oxl.at/acme-client/pkg/config"
)

func processDomainBatch(grp config.Group, svc config.Service, batchID int, domains []string) (bool, error) {
	key := fmt.Sprintf("[Group: %d '%s' | Service: %d | Batch: %d]", grp.ID, grp.Name, svc.ID, batchID)
	certBaseName := fmt.Sprintf("%d_%d_%d", grp.ID, svc.ID, batchID)
	u.LogDebugf("%s File %s with domains: %+v", key, certBaseName, domains)

	updateNeeded, reason := NeedsUpdate(certBaseName, domains)
	if !updateNeeded {
		u.Logf("%s Skipping: cert is valid", key)
		return false, nil
	}

	u.Logf("%s Updating: %s", key, reason)

	for attempt := uint(0); attempt <= config.Config.Retries; attempt++ {
		time.Sleep(time.Second * time.Duration(config.Config.CooldownSec))
		u.Logf("%s Obtaining certificate...", key)

		if config.ModeCheck {
			u.Logf("%s Skipping ACME-interaction in check-mode - domains: %v", key, domains)
			return true, nil

		} else {
			err := obtainCert(certBaseName, svc, domains)
			if err == nil {
				return true, nil
			}
			u.LogWarningf("%s Attempt %d failed: \"%v\"", key, attempt+1, strings.Trim(fmt.Sprintf("%v", err), "\n"))
		}
	}

	u.LogWarningf("%s Could not process certificate after %d retries", key, config.Config.Retries)
	return false, fmt.Errorf("failed")
}

func processDomainsInBatches(
	grp config.Group, svc config.Service, domains []string,
	callback func(config.Group, config.Service, int, []string) (bool, error),
) (bool, int, error) {
	anyChanged := false
	var err error
	processedBatches := 0
	for batchID, batchDomains := range u.BuildBatches(domains, int(config.Config.MaxDomains)) {
		processedBatches += 1
		changed, err := callback(grp, svc, batchID+1, batchDomains)
		if changed {
			anyChanged = true
		}
		if err != nil {
			return anyChanged, processedBatches, err
		}
	}
	return anyChanged, processedBatches, err
}

func processService(grp config.Group, svc config.Service) (bool, error) {
	key := fmt.Sprintf("[Group: %d '%s' | Service: %d]", grp.ID, grp.Name, svc.ID)
	u.Logf("%s Processing...", key)

	if len(svc.Domains) == 0 {
		u.Logf("%s Skipping: service has no domains configured", key)
		return false, nil
	}

	providedDomainCount := len(svc.Domains)
	svc.Domains = u.RemoveDuplicates(svc.Domains)
	uniqueDomainCount := len(svc.Domains)
	if uniqueDomainCount != providedDomainCount {
		u.LogWarningf("%s Has %d duplicate domains configured", key, providedDomainCount-uniqueDomainCount)
	}

	changed, _, err := processDomainsInBatches(grp, svc, svc.Domains, processDomainBatch)
	return changed, err
}
