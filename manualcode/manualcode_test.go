package manualcode

import (
	"testing"
)

func TestGenerateManualcode(t *testing.T) {
	tests := []struct {
		name          string
		payload       PayloadContents
		manualcode    string
	}{
		{
			name: "Test 1",
			payload: PayloadContents {
				SetUpPINCode:      49910688, 
				Discriminator:     3431, 
				VendorID:      	   0x1317,
				ProductID:         0x0002, 
				CommissioningFlow: 0,
			},
			manualcode: "32140830464",
		},
		{
			name:"Test 2", 
			payload: PayloadContents {
				SetUpPINCode:      54757432,
				Discriminator:     80, 
				VendorID: 		   0, 
				ProductID: 		   0, 
				CommissioningFlow: 0,
			},
			manualcode: "00210433428",
		},
		{ 
			name:"Test 3", 
			payload: PayloadContents {
				SetUpPINCode:      43338551, 
				Discriminator:     3091, 
				VendorID:          0x1123, 
				ProductID:         0x0012, 
				CommissioningFlow: 1,
			},
			manualcode: "702871264504387000187",
		},
		{
			name:"Test 4",
			payload: PayloadContents {
				SetUpPINCode:      20202021, 
				Discriminator: 	   3840, 
				VendorID: 		   0, 
				ProductID: 		   0, 
				CommissioningFlow: 0,
			},
			manualcode: "34970112332",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &ManualSetupPayloadGenerator{
				PayloadContents: tt.payload,
			}
			manualcode, err := g.GenerateManualcode()
			if err != nil {
				t.Errorf("GenerateManualcode() error = %v", err)
			}
			if manualcode != tt.manualcode {
				t.Errorf("GenerateManualcode() = %v, want %v", manualcode, tt.manualcode)
			}
		})
	}
}
