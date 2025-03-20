// Copyright 2025 Espressif Systems (Shanghai) PTE LTD
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spake2p

import (
	"testing"
	"encoding/base64"
	"bytes"
)

// test data generated using esp-matter-mfg-tool -n 5 --vendor-id 0xfff1 --product-id 0x8001
func TestGenerateVerifier(t *testing.T) {
	tests := []struct {
		name       string
		passcode   uint32
		salt       string // base64 encoded
		verifier   string // base64 encoded
			}{
		{"Test 1", 38411287, "CTNvIPK4s0cPo0EtSLbIUMAB1j5xQjAA4P+ock5LaHI=", "X+3EKg3pcZ+h2nzmRQQt58vjB5jcEpbHJ9oGHwEfbMUEbAdOxOq2JfcQi2okuZ81F1PHKhI2xhPGuDdHJr85ZEG7JZ6GjQTiJ2ZkuFrGwxW+F9GI5Q59LvbE9LFOxo7CAg=="},
		{"Test 2", 90640905, "4PM0ZFe6COhZWnoS+6ka2qc3sj+XaCywL/tRRmFK3s8=", "XwJ+pUUnAbcikIoPV+FNB4eY7ogohmYc/NRUAjcd6tAE64Va5awqZ1TDJjA+qu4nQJ5ETKVX6tXYHKEx1Y24MyY8oJ4Am4lS21spAMhQ+Mod3HI+BM4RN0h1ESL2CCGsvw=="},
		{"Test 3", 46806472, "XYVhFHnKtZAvmrQuo3Usmqmsn7YOyYV1RSUJ538ayUw=", "h6kV9sEmfjgCMEraoPhJqAsDn624H7qjl05UAfdBhU0EYgsi7z6XdMvzGdyycdIA/j2PXnwz/Q7GmQ6qmIOMXMuHQJ9UCL5vzFWSIUwdjBM6zLgMOtQCkBPsBCBomxiRhw=="},
		{"Test 4", 3796423, "0jMvjaY7pvUUK3a+46YqeVo/jsNxdH5B0sAoF9brOeo=", "rrRZfDpO3fNyhiLLjZSyP1iQ3K2UauvSKYU+dbANrA4EpFCn0U02/HErG2TdSWt9VumRfzLiOfJ61XfjLUCYscoGeQRan8r106UHcKnSJbJGeAElAK7TNg/j/YiVKrThzw=="},
		{"Test 5", 68120576, "PyC3/Qs7Et25lscGmx85Frd4LCl1cAEKRrKC3HBdlLU=", "UWC2EtDfhPR5dTbzEEV+MnT4X7wycUe3kWRxumN0MaIErmrqzaHOouOKCOOiPU+bW7ojvAXDDxPsVzVJ7iVTtM+dQfhnzkMPbKFHVsrf6emzkyOMfKau+Ve0cLrRaiKdlQ=="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			salt, err := base64.StdEncoding.DecodeString(tt.salt)
			if err != nil {
				t.Fatalf("Failed to decode salt: %v", err)
			}

			expectedVerifier, err := base64.StdEncoding.DecodeString(tt.verifier)
			if err != nil {
				t.Fatalf("Failed to decode verifier: %v", err)
			}

			// iterations are fixed to 1000
			verifier, err := GenerateVerifier(tt.passcode, salt, uint16(10000))
			if err != nil {
				t.Errorf("GenerateVerifier failed: %v", err)
			}

			if !bytes.Equal(verifier, expectedVerifier) {
				t.Errorf("Verifier mismatch\nexpected: %x\ngot:      %x", expectedVerifier, verifier)
			}
		})
	}
}
