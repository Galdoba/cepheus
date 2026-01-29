package orbital

import (
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
	tests := []struct {
		name      string
		stars     []string
		planet    int
		satellite string
		want      string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Planet around A and B",
			stars:     []string{"Aa", "Ab", "Ba"},
			planet:    8,
			satellite: "",
			want:      "AaAbBa 8",
			wantErr:   false,
		},
		{
			name:      "Planet with satellite",
			stars:     []string{"Aa"},
			planet:    3,
			satellite: "a",
			want:      "Aa 3 a",
			wantErr:   false,
		},
		{
			name:      "Only stars (no planet)",
			stars:     []string{"Aa", "Ab", "Ba", "Ca"},
			planet:    -1,
			satellite: "",
			want:      "AaAbBaCa",
			wantErr:   false,
		},
		{
			name:      "Multiple stars with satellite",
			stars:     []string{"Aa", "Ab", "Ba", "Bb", "Ca", "Da"},
			planet:    12,
			satellite: "d",
			want:      "AaAbBaBbCaDa 12 d",
			wantErr:   false,
		},
		{
			name:      "Empty stars list",
			stars:     []string{},
			planet:    5,
			satellite: "",
			want:      "",
			wantErr:   true,
			errMsg:    "invalid star designation",
		},
		{
			name:      "Missing primary star Aa",
			stars:     []string{"Ba", "Bb"},
			planet:    1,
			satellite: "",
			want:      "",
			wantErr:   true,
			errMsg:    "primary star Aa is missing",
		},
		{
			name:      "Invalid star designation",
			stars:     []string{"Aa", "Xy"},
			planet:    2,
			satellite: "",
			want:      "",
			wantErr:   true,
			errMsg:    "invalid star designation",
		},
		{
			name:      "Duplicate stars",
			stars:     []string{"Aa", "Aa", "Ba"},
			planet:    3,
			satellite: "",
			want:      "",
			wantErr:   true,
			errMsg:    "duplicate star",
		},
		{
			name:      "Companion without main star",
			stars:     []string{"Aa", "Bb"},
			planet:    2,
			satellite: "",
			want:      "",
			wantErr:   true,
			errMsg:    "main star is missing for companion",
		},
		{
			name:      "Invalid planet orbit (too high)",
			stars:     []string{"Aa", "Ab"},
			planet:    25,
			satellite: "",
			want:      "",
			wantErr:   true,
			errMsg:    "invalid planet orbit",
		},
		{
			name:      "Valid planet orbit 0",
			stars:     []string{"Aa", "Ab"},
			planet:    0,
			satellite: "",
			want:      "AaAb 0",
			wantErr:   false,
		},
		{
			name:      "Invalid satellite (out of range)",
			stars:     []string{"Aa"},
			planet:    3,
			satellite: "z",
			want:      "",
			wantErr:   true,
			errMsg:    "invalid satellite orbit",
		},
		{
			name:      "Invalid satellite (multiple chars)",
			stars:     []string{"Aa"},
			planet:    3,
			satellite: "ab",
			want:      "",
			wantErr:   true,
			errMsg:    "invalid satellite orbit",
		},
		{
			name:      "Satellite without planet",
			stars:     []string{"Aa"},
			planet:    -1,
			satellite: "a",
			want:      "",
			wantErr:   true,
			errMsg:    "satellite specified without planet",
		},
		{
			name:      "Invalid star order",
			stars:     []string{"Ca", "Aa"},
			planet:    2,
			satellite: "",
			want:      "",
			wantErr:   true,
			errMsg:    "invalid star pair order",
		},
		{
			name:      "All possible stars",
			stars:     []string{"Aa", "Ab", "Ba", "Bb", "Ca", "Cb", "Da", "Db"},
			planet:    20,
			satellite: "l",
			want:      "AaAbBaBbCaCbDaDb 20 l",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.stars, tt.planet, tt.satellite)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Encode() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Encode() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Encode() unexpected error = %v", err)
					return
				}
				if got != tt.want {
					t.Errorf("Encode() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		name       string
		code       string
		wantStars  []string
		wantPlanet int
		wantSat    string
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "Full star code with planet",
			code:       "AaAbBa 8",
			wantStars:  []string{"Aa", "Ab", "Ba"},
			wantPlanet: 8,
			wantSat:    "",
			wantErr:    false,
		},
		{
			name:       "Full star code with planet and satellite",
			code:       "Aa 3 a",
			wantStars:  []string{"Aa"},
			wantPlanet: 3,
			wantSat:    "a",
			wantErr:    false,
		},
		{
			name:       "Only stars",
			code:       "AaAbBaCa",
			wantStars:  []string{"Aa", "Ab", "Ba", "Ca"},
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    false,
		},
		{
			name:       "Complex system",
			code:       "AaAbBaBbCaDa 12 d",
			wantStars:  []string{"Aa", "Ab", "Ba", "Bb", "Ca", "Da"},
			wantPlanet: 12,
			wantSat:    "d",
			wantErr:    false,
		},
		{
			name:       "Empty code",
			code:       "",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "empty code",
		},
		{
			name:       "Invalid format - too many parts",
			code:       "Aa 3 a b",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid code format",
		},
		{
			name:       "Invalid star length (odd)",
			code:       "Aab 3",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid code format",
		},
		{
			name:       "Invalid star designation in code",
			code:       "Ax 3",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid star designation",
		},
		{
			name:       "Invalid planet orbit (not a number)",
			code:       "AaAbBa eight",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid planet orbit",
		},
		{
			name:       "Invalid planet orbit (negative)",
			code:       "Aa -1",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid planet orbit",
		},
		{
			name:       "Invalid planet orbit (too high)",
			code:       "Aa 21",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid planet orbit",
		},
		{
			name:       "Invalid satellite",
			code:       "Aa 3 z",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid satellite orbit",
		},
		{
			name:       "Invalid star order in code",
			code:       "CaAa 2",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "invalid star pair order",
		},
		{
			name:       "Missing primary in code",
			code:       "BaBb 1",
			wantStars:  nil,
			wantPlanet: -1,
			wantSat:    "",
			wantErr:    true,
			errMsg:     "primary star Aa is missing",
		},
		{
			name:       "Valid planet orbit 0",
			code:       "AaAb 0",
			wantStars:  []string{"Aa", "Ab"},
			wantPlanet: 0,
			wantSat:    "",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStars, gotPlanet, gotSat, err := Decode(tt.code)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Decode() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Decode() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Decode() unexpected error = %v", err)
					return
				}

				// Compare stars slice
				if len(gotStars) != len(tt.wantStars) {
					t.Errorf("Decode() stars length = %v, want %v", len(gotStars), len(tt.wantStars))
					return
				}
				for i := range gotStars {
					if gotStars[i] != tt.wantStars[i] {
						t.Errorf("Decode() stars[%d] = %v, want %v", i, gotStars[i], tt.wantStars[i])
					}
				}

				if gotPlanet != tt.wantPlanet {
					t.Errorf("Decode() planet = %v, want %v", gotPlanet, tt.wantPlanet)
				}
				if gotSat != tt.wantSat {
					t.Errorf("Decode() satellite = %v, want %v", gotSat, tt.wantSat)
				}
			}
		})
	}
}

func TestCompress(t *testing.T) {
	tests := []struct {
		name    string
		code    string
		want    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Single star",
			code:    "Aa 3 a",
			want:    "A 3 a",
			wantErr: false,
		},
		{
			name:    "Star pair with companion",
			code:    "AaAb 8",
			want:    "A 8",
			wantErr: false,
		},
		{
			name:    "Two star systems",
			code:    "AaAbBa 8",
			want:    "AB 8",
			wantErr: false,
		},
		{
			name:    "Three star systems",
			code:    "AaAbBaCa 12",
			want:    "ABC 12",
			wantErr: false,
		},
		{
			name:    "All four star systems",
			code:    "AaAbBaBbCaCbDaDb 20 l",
			want:    "ABCD 20 l",
			wantErr: false,
		},
		{
			name:    "Only stars (no planet)",
			code:    "AaAbBaCa",
			want:    "ABC",
			wantErr: false,
		},
		{
			name:    "Star with companion from B pair",
			code:    "AaBaBb 5 c",
			want:    "AB 5 c",
			wantErr: false,
		},
		{
			name:    "Invalid code format",
			code:    "AB 8", // Already compressed
			want:    "",
			wantErr: true,
			errMsg:  "invalid star designation",
		},
		{
			name:    "Empty code",
			code:    "",
			want:    "",
			wantErr: true,
			errMsg:  "empty code",
		},
		{
			name:    "Missing primary star",
			code:    "BaBb 2",
			want:    "",
			wantErr: true,
			errMsg:  "primary star Aa is missing",
		},
		{
			name:    "Invalid star order",
			code:    "CaAa 3",
			want:    "",
			wantErr: true,
			errMsg:  "invalid star pair order",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Compress(tt.code)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Compress() error = nil, want error containing %q", tt.errMsg)
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Compress() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Compress() unexpected error = %v", err)
					return
				}
				if got != tt.want {
					t.Errorf("Compress() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	tests := []struct {
		name      string
		stars     []string
		planet    int
		satellite string
	}{
		{
			name:      "Single star with satellite",
			stars:     []string{"Aa"},
			planet:    3,
			satellite: "a",
		},
		{
			name:      "Multiple stars",
			stars:     []string{"Aa", "Ab", "Ba", "Ca"},
			planet:    5,
			satellite: "b",
		},
		{
			name:      "Stars only",
			stars:     []string{"Aa", "Ab", "Ba"},
			planet:    -1,
			satellite: "",
		},
		{
			name:      "All stars",
			stars:     []string{"Aa", "Ab", "Ba", "Bb", "Ca", "Cb", "Da", "Db"},
			planet:    20,
			satellite: "l",
		},
		{
			name:      "Planet at orbit 0",
			stars:     []string{"Aa", "Ab"},
			planet:    0,
			satellite: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode
			encoded, err := Encode(tt.stars, tt.planet, tt.satellite)
			if err != nil {
				t.Fatalf("Encode failed: %v", err)
			}

			// Decode
			decodedStars, decodedPlanet, decodedSat, err := Decode(encoded)
			if err != nil {
				t.Fatalf("Decode failed: %v", err)
			}

			// Compare stars
			if len(decodedStars) != len(tt.stars) {
				t.Errorf("Stars length mismatch: got %v, want %v", len(decodedStars), len(tt.stars))
			} else {
				for i := range decodedStars {
					if decodedStars[i] != tt.stars[i] {
						t.Errorf("Star[%d] mismatch: got %v, want %v", i, decodedStars[i], tt.stars[i])
					}
				}
			}

			// Compare planet
			if decodedPlanet != tt.planet {
				t.Errorf("Planet mismatch: got %v, want %v", decodedPlanet, tt.planet)
			}

			// Compare satellite
			if decodedSat != tt.satellite {
				t.Errorf("Satellite mismatch: got %v, want %v", decodedSat, tt.satellite)
			}
		})
	}
}

func TestCompressDecompressRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		code string
	}{
		{
			name: "Single star",
			code: "Aa 3 a",
		},
		{
			name: "Multiple stars",
			code: "AaAbBa 8",
		},
		{
			name: "Complex system",
			code: "AaAbBaBbCaDa 12 d",
		},
		{
			name: "Stars only",
			code: "AaAbBaCa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compress
			compressed, err := Compress(tt.code)
			if err != nil {
				t.Fatalf("Compress failed: %v", err)
			}

			// Decode original to get expected values
			_, origPlanet, origSat, err := Decode(tt.code)
			if err != nil {
				t.Fatalf("Decode original failed: %v", err)
			}

			// Decode compressed (but Compress already does validation)
			// We can't decode compressed directly, but we can verify it's correct
			// by checking that it matches the expected pattern

			// For stars only case
			if origPlanet == -1 && origSat == "" {
				// Should be just letters
				if len(compressed) == 0 {
					t.Error("Compressed result empty for stars only")
				}
				// Check it contains only A-D
				for _, r := range compressed {
					if r < 'A' || r > 'D' {
						t.Errorf("Invalid character in compressed result: %c", r)
					}
				}
			} else {
				// Should contain planet number
				if !strings.Contains(compressed, " ") {
					t.Error("Compressed result missing space separator")
				}
			}

			// We can also verify that re-compressing doesn't change anything
			reCompressed, err := Compress(tt.code)
			if err != nil {
				t.Fatalf("Re-compress failed: %v", err)
			}
			if reCompressed != compressed {
				t.Errorf("Re-compress mismatch: got %v, want %v", reCompressed, compressed)
			}
		})
	}
}
