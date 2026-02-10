package common

import "testing"

func TestDecToHex(t *testing.T) {
	tests := map[int]string{
		17:  "11",
		23:  "17",
		34:  "22",
		41:  "29",
		58:  "3A",
		67:  "43",
		79:  "4F",
		92:  "5C",
		104: "68",
		117: "75",
		131: "83",
		146: "92",
		159: "9F",
		173: "AD",
		198: "C6",
	}

	for dec, expectedHex := range tests {
		got, err := ConvertDecToHex(dec)
		if err != nil {
			t.Fatalf("unexpected error for %d: %v", dec, err)
		}

		if got != expectedHex {
			t.Errorf("convertDecToHex(%d) = %s, want %s", dec, got, expectedHex)
		}
	}
}
