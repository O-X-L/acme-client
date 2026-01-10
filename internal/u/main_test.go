package u

import (
	"reflect"
	"testing"
)

func TestSortDomainsExceptFirst(t *testing.T) {
	data := []string{"oxl.at", "www.oxl.at", "host-svc.com", "xyz.oxl.at"}
	want := []string{"oxl.at", "host-svc.com", "www.oxl.at", "xyz.oxl.at"}

	SortDomainsExceptFirst(data)
	if !reflect.DeepEqual(data, want) {
		t.Errorf("SortDomainsExceptFirst failed, got data=%v, want=%s", data, want)
	}
}

func TestRemoveDuplicates(t *testing.T) {
	tests := []struct {
		name string
		data []string
		want []string
	}{
		{
			name: "One duplicate",
			data: []string{"oxl.at", "www.oxl.at", "abc.oxl.at", "www.oxl.at"},
			want: []string{"oxl.at", "www.oxl.at", "abc.oxl.at", "www.oxl.at"},
		},
		{
			name: "Two duplicates",
			data: []string{"abc.oxl.at", "www.oxl.at", "abc.oxl.at", "www.oxl.at"},
			want: []string{"abc.oxl.at", "www.oxl.at", "abc.oxl.at", "www.oxl.at"},
		},
		{
			name: "Three duplicates",
			data: []string{"abc.oxl.at", "www.oxl.at", "abc.oxl.at", "abc.oxl.at"},
			want: []string{"abc.oxl.at", "www.oxl.at", "abc.oxl.at", "abc.oxl.at"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RemoveDuplicates(tt.data)
			if reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s failed: reason: %v != %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestBuildDiffMap(t *testing.T) {
	tests := []struct {
		name   string
		before []string
		after  []string
		want   map[string][]string
	}{
		{
			name:   "Creation",
			before: []string{},
			after:  []string{"oxl.at", "www.oxl.at"},
			want: map[string][]string{
				"+": {"oxl.at", "www.oxl.at"},
			},
		},
		{
			name:   "Change",
			before: []string{"oxl.at", "www.oxl.at"},
			after:  []string{"zuo.oxl.at", "oxl.at", "host-svc.com"},
			want: map[string][]string{
				"+": {"host-svc.com", "zuo.oxl.at"},
				"-": {"www.oxl.at"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildDiffMap(tt.before, tt.after)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s failed: reason: %v != %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestBuildDiffString(t *testing.T) {
	tests := []struct {
		name string
		diff map[string][]string
		want string
	}{
		{
			name: "Creation",
			diff: map[string][]string{
				"+": {"oxl.at", "www.oxl.at"},
				"-": {},
			},
			want: "+[oxl.at, www.oxl.at]",
		},
		{
			name: "Change",
			diff: map[string][]string{
				"+": {"host-svc.com", "zuo.oxl.at"},
				"-": {"www.oxl.at"},
			},
			want: "+[host-svc.com, zuo.oxl.at] -[www.oxl.at]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildDiffString(tt.diff)
			if got != tt.want {
				t.Errorf("%s failed: reason: '%v' != '%v'", tt.name, got, tt.want)
			}
		})
	}
}
