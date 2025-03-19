package manualcode

import (
	"fmt"
	"log"
	"strings"
	"github.com/osamingo/checkdigit"
)

const (
	manualSetupCodeChunk1CharLength = 1
	manualSetupCodeChunk2CharLength = 5
	manualSetupCodeChunk3CharLength = 4
	manualSetupShortCodeCharLength  = 10
	manualSetupLongCodeCharLength   = 20
	manualSetupVendorIdCharLength   = 5
	manualSetupProductIdCharLength  = 5
)

type PayloadContents struct {
	SetUpPINCode      uint32
	Discriminator     uint16
	CommissioningFlow uint8
	VendorID          uint16
	ProductID         uint16
}

type ManualSetupPayloadGenerator struct {
	PayloadContents   PayloadContents
}

func chunk1PayloadRepresentation(payload PayloadContents) uint8 {
	version := uint8(0)
	vidPidPresent := uint8(0)

	isLong := (payload.CommissioningFlow != 0)
	if isLong {
		vidPidPresent = uint8(1)
	}
	disc := uint8(payload.Discriminator >> 10) & 0x0F

 	return uint8(version << 3 | vidPidPresent << 2 | disc)
}

func chunk2PayloadRepresentation(payload PayloadContents) uint16 {
	disc := uint16((payload.Discriminator & 0x300) << 6)
	passcode := uint16(payload.SetUpPINCode & 0x3FFF)
	return disc | passcode
}

func chunk3PayloadRepresentation(payload PayloadContents) uint16 {
	return uint16(payload.SetUpPINCode >> 14)
}

func computeVerhoeffCheckDigit(input string) int {
	p := checkdigit.NewVerhoeff()
	cd, err := p.Generate(input)
	if err != nil {
		log.Fatal(err)
	}
	return cd
}

func (g *ManualSetupPayloadGenerator) GenerateManualcode() (string, error) {
	useLongCode := (g.PayloadContents.CommissioningFlow != 0)

	chunk1 := chunk1PayloadRepresentation(g.PayloadContents)
	chunk2 := chunk2PayloadRepresentation(g.PayloadContents)
	chunk3 := chunk3PayloadRepresentation(g.PayloadContents)

	var result strings.Builder

	// Write chunks with proper padding
	fmt.Fprintf(&result, "%0*d", manualSetupCodeChunk1CharLength, chunk1)
	fmt.Fprintf(&result, "%0*d", manualSetupCodeChunk2CharLength, chunk2)
	fmt.Fprintf(&result, "%0*d", manualSetupCodeChunk3CharLength, chunk3)

	if useLongCode {
		fmt.Fprintf(&result, "%0*d", manualSetupVendorIdCharLength, g.PayloadContents.VendorID)
		fmt.Fprintf(&result, "%0*d", manualSetupProductIdCharLength, g.PayloadContents.ProductID)
	}

	code := result.String()
	checkDigit := computeVerhoeffCheckDigit(code)
	return code + fmt.Sprintf("%d", checkDigit), nil
}