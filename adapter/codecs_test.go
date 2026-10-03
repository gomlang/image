package adapter

import (
	"encoding/binary"
	"os"
	"strings"
	"testing"
)

func readFixture(t *testing.T, name string) string {
	t.Helper()
	data, err := os.ReadFile("../tests/data/" + name)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestIndependentFixtures(t *testing.T) {
	for _, name := range []string{"solid.jpg", "progressive.jpg"} {
		raw, w, h, err := DecodeJPEG(readFixture(t, name), 8388608, 4096, 1048576, 4194304)
		if err != "" || w != 8 || h != 8 || len(raw) != 256 {
			t.Fatalf("%s: %d x %d: %s", name, w, h, err)
		}
		for i := 0; i < len(raw); i += 4 {
			for channel, expected := range []byte{32, 96, 160, 255} {
				delta := int(raw[i+channel]) - int(expected)
				if delta < -2 || delta > 2 {
					t.Fatalf("%s pixel %d channel %d: %d", name, i/4, channel, raw[i+channel])
				}
			}
		}
	}
	raw, w, h, err := DecodeWebP(readFixture(t, "alpha.webp"), 8388608, 4096, 1048576, 4194304)
	expected := string([]byte{255, 0, 0, 255, 0, 128, 0, 128, 0, 0, 64, 64, 0, 0, 0, 0})
	if err != "" || w != 2 || h != 2 || raw != expected {
		t.Fatalf("lossless fixture: %d x %d %x: %s", w, h, raw, err)
	}
}

func TestHeaderBoundsAndCanvasMismatch(t *testing.T) {
	jpg := []byte(readFixture(t, "solid.jpg"))
	changed := false
	for i := 0; i+9 < len(jpg); i++ {
		if jpg[i] == 255 && jpg[i+1] == 0xc0 {
			jpg[i+5], jpg[i+6], jpg[i+7], jpg[i+8] = 255, 255, 255, 255
			changed = true
			break
		}
	}
	if !changed {
		t.Fatal("fixture lacks baseline frame header")
	}
	if raw, _, _, err := DecodeJPEG(string(jpg), 8388608, 4096, 1048576, 4194304); err == "" || raw != "" {
		t.Fatal("oversized JPEG header accepted")
	}
	wp := []byte(readFixture(t, "alpha.webp"))
	payload := make([]byte, 0, len(wp)+18)
	payload = append(payload, wp[:12]...)
	payload = append(payload, []byte{'V', 'P', '8', 'X', 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}...)
	payload = append(payload, wp[12:]...)
	binary.LittleEndian.PutUint32(payload[4:8], uint32(len(payload)-8))
	if raw, _, _, err := DecodeWebP(string(payload), 8388608, 4096, 1048576, 4194304); err == "" || raw != "" {
		t.Fatal("mismatched WebP canvas accepted")
	}
	payload[20] = 2
	if _, _, _, err := DecodeWebP(string(payload), 8388608, 4096, 1048576, 4194304); err == "" {
		t.Fatal("animated WebP accepted")
	}
	if _, _, _, err := DecodeWebP(string(wp), 8388608, 4096, 1048576, 15); err == "" {
		t.Fatal("decoded byte limit ignored")
	}
}

func TestMalformedPrefixesRemainRecoverable(t *testing.T) {
	for _, name := range []string{"solid.jpg", "progressive.jpg", "alpha.webp", "gray.webp"} {
		data := readFixture(t, name)
		for length := 0; length < len(data)/2; length++ {
			var raw, err string
			if strings.HasSuffix(name, ".jpg") {
				raw, _, _, err = DecodeJPEG(data[:length], 8388608, 4096, 1048576, 4194304)
			} else {
				raw, _, _, err = DecodeWebP(data[:length], 8388608, 4096, 1048576, 4194304)
			}
			if err == "" || raw != "" {
				t.Fatalf("%s truncated to %d accepted", name, length)
			}
		}
	}
}

func TestEncoderBoundsAndValidation(t *testing.T) {
	pixel := string([]byte{128, 0, 0, 128})
	if _, err := EncodeJPEG(pixel, 1, 1, 4096, 1048576, 4194304, 10, 100, 0xffffffff); err == "" {
		t.Fatal("JPEG writer limit ignored")
	}
	if _, err := EncodeWebP(pixel, 1, 1, 4096, 1048576, 4194304, 10, 0); err == "" {
		t.Fatal("WebP writer limit ignored")
	}
	for _, raw := range []string{"", "123", string([]byte{255, 0, 0, 128})} {
		if _, err := EncodeJPEG(raw, 1, 1, 4096, 1048576, 4194304, 8388608, 90, 0xffffffff); err == "" {
			t.Fatal("invalid JPEG source accepted")
		}
		if _, err := EncodeWebP(raw, 1, 1, 4096, 1048576, 4194304, 8388608, 0); err == "" {
			t.Fatal("invalid WebP source accepted")
		}
	}
}

func TestLosslessPremultipliedChannels(t *testing.T) {
	raw := make([]byte, 256*256*4)
	for a := 0; a < 256; a++ {
		for r := 0; r < 256; r++ {
			i := (a*256 + r) * 4
			raw[i] = byte(min(r, a))
			raw[i+1] = byte(a / 2)
			raw[i+2] = byte(a)
			raw[i+3] = byte(a)
		}
	}
	encoded, err := EncodeWebP(string(raw), 256, 256, 4096, 1048576, 4194304, 8388608, 0)
	if err != "" {
		t.Fatal(err)
	}
	decoded, w, h, err := DecodeWebP(encoded, 8388608, 4096, 1048576, 4194304)
	if err != "" || w != 256 || h != 256 {
		t.Fatalf("encoded WebP rejected: %s", err)
	}
	for i := range raw {
		if decoded[i] != raw[i] {
			t.Fatalf("channel %d got %d want %d", i, decoded[i], raw[i])
		}
	}
}

func TestWebPColorRangeMatchesIndependentLibwebp(t *testing.T) {
	cases := []struct {
		name  string
		pixel [4]byte
	}{
		{"color.webp", [4]byte{31, 96, 160, 255}},
		{"color-alpha.webp", [4]byte{16, 48, 80, 128}},
		{"black.webp", [4]byte{0, 0, 0, 255}},
		{"white.webp", [4]byte{255, 255, 255, 255}},
	}
	for _, tc := range cases {
		raw, w, h, err := DecodeWebP(readFixture(t, tc.name), 8388608, 4096, 1048576, 4194304)
		if err != "" || w != 8 || h != 8 {
			t.Fatalf("%s: %d x %d: %s", tc.name, w, h, err)
		}
		for i := range raw {
			delta := int(raw[i]) - int(tc.pixel[i%4])
			if delta < -1 || delta > 1 {
				t.Fatalf("%s channel %d: got %d want %d", tc.name, i, raw[i], tc.pixel[i%4])
			}
		}
	}
}
