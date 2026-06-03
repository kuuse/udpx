package targets

import (
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		target    string
		wantCount int
		wantFirst string
		wantErr   bool
	}{
		// Bare IPs
		{"bare IP", "192.0.2.1", 1, "192.0.2.1", false},
		{"bare IP v6", "::1", 1, "::1", false},

		// CIDR
		{"CIDR /30", "192.0.2.0/30", 4, "192.0.2.0", false},
		{"CIDR /31", "192.0.2.0/31", 2, "192.0.2.0", false},

		// Octet ranges
		{"octet range simple", "192.0.2.1-3", 3, "192.0.2.1", false},
		{"octet range reversed", "192.0.2.3-1", 3, "192.0.2.1", false},
		{"octet range multi", "192.0-1.2.1-2", 4, "192.0.2.1", false},

		// Edge cases
		{"empty", "", 0, "", false},
		{"whitespace", "  192.0.2.1  ", 1, "192.0.2.1", false},

		// Errors
		{"invalid CIDR", "192.0.2.0/33", 0, "", true},
		{"invalid octet range", "192.0.2.1-300", 0, "", true},
		{"invalid format", "not-an-ip", 0, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.target, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("Parse(%q) count = %d, want %d", tt.target, len(got), tt.wantCount)
			}
			if tt.wantCount > 0 && got[0] != tt.wantFirst {
				t.Errorf("Parse(%q) first = %q, want %q", tt.target, got[0], tt.wantFirst)
			}
		})
	}
}

// BenchmarkParseOctetRange measures the cost of expanding multi-octet ranges.
func BenchmarkParseOctetRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("192.0-255.0-255.1-10")
	}
}

// sanityCheckOctetExpansion ensures multi-octet ranges are correctly ordered.
func TestOctetRangeOrder(t *testing.T) {
	got, err := Parse("192.0-1.0-1.1-2")
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	// Should expand to: 192.0.0.1, 192.0.0.2, 192.0.1.1, 192.0.1.2,
	//                   192.1.0.1, 192.1.0.2, 192.1.1.1, 192.1.1.2
	wantOrder := []string{
		"192.0.0.1", "192.0.0.2",
		"192.0.1.1", "192.0.1.2",
		"192.1.0.1", "192.1.0.2",
		"192.1.1.1", "192.1.1.2",
	}
	if len(got) != len(wantOrder) {
		t.Errorf("Got %d IPs, want %d", len(got), len(wantOrder))
		return
	}
	for i, want := range wantOrder {
		if got[i] != want {
			t.Errorf("Position %d: got %q, want %q", i, got[i], want)
		}
	}
}
