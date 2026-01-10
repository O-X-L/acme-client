package manager

import (
	"fmt"
	"slices"
	"testing"

	"git.oxl.at/acme-client/pkg/config"
)

func TestProcessDomainsInBatches(t *testing.T) {
	config.Config.MaxDomains = 2
	dummyCallback := func(grp config.Group, svc config.Service, id int, input []string) (bool, error) {
		if slices.Contains(input, "eee") {
			return false, fmt.Errorf("y")
		}
		if slices.Contains(input, "ccc") {
			return true, nil
		}
		return false, nil
	}
	dummyGrp := config.Group{}
	dummySvc := config.Service{}

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
			changed, procssedBatches, err := processDomainsInBatches(dummyGrp, dummySvc, tt.input, dummyCallback)

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
