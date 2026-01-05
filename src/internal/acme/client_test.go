package acme

import "testing"

func TestIsSupportedProvider(t *testing.T) {
	tests := []struct {
		name     string
		provider string
		want     bool
	}{
		{"Valid Cloudflare", "cloudflare", true},
		{"Valid Route53", "route53", true},
		{"Invalid Provider", "not-a-real-provider", false},
		{"Empty Provider", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupportedProvider(tt.provider); got != tt.want {
				t.Errorf("IsSupportedProvider(%s) = %v, want %v", tt.provider, got, tt.want)
			}
		})
	}
}
