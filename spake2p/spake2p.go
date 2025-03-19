package spake2p

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"golang.org/x/crypto/pbkdf2"
	"math/big"
	"crypto/rand"
	"encoding/base64"
)

// VerifierParams contains all parameters for SPAKE2+ verification
type VerifierParams struct {
    Passcode   uint32
    Salt       []byte
    SaltBase64 string
    Iterations uint16
    Verifier   []byte
}

var (
	InvalidPasscodes = map[uint32]bool{
		00000000: true,
		11111111: true,
		22222222: true,
		33333333: true,
		44444444: true,
		55555555: true,
		66666666: true,
		77777777: true,
		88888888: true,
		99999999: true,
		12345678: true,
		87654321: true,
	}

	curve = elliptic.P256()
)

// GenerateVerifier generates a SPAKE2+ verifier
func GenerateVerifier(passcode uint32, salt []byte, iterations uint16) ([]byte, error) {
	passBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(passBytes, passcode)
	
	wsLength := (curve.Params().BitSize + 7) / 8 + 8
	ws := pbkdf2.Key(passBytes, salt, int(iterations), wsLength*2, sha256.New)

	w0 := new(big.Int).SetBytes(ws[:wsLength])
	w1 := new(big.Int).SetBytes(ws[wsLength:])
	w0.Mod(w0, curve.Params().N)
	w1.Mod(w1, curve.Params().N)

	x, y := curve.ScalarBaseMult(w1.Bytes())
	if x == nil {
		return nil, fmt.Errorf("failed to compute L point")
	}

	// Ensure proper byte lengths for each component
	w0Bytes := make([]byte, 32)
	xBytes := make([]byte, 32)
	yBytes := make([]byte, 32)
	
	w0Temp := w0.Bytes()
	copy(w0Bytes[32-len(w0Temp):], w0Temp)
	
	xTemp := x.Bytes()
	copy(xBytes[32-len(xTemp):], xTemp)
	
	yTemp := y.Bytes()
	copy(yBytes[32-len(yTemp):], yTemp)

	result := make([]byte, 0, 97) // 32 + 1 + 32 + 32
	result = append(result, w0Bytes...)
	result = append(result, 0x04) // Uncompressed point format
	result = append(result, xBytes...)
	result = append(result, yBytes...)

	return result, nil
}

// GenerateRandomPasscode generates a valid random passcode between 0-99999999
func GenerateRandomPasscode() (uint32, error) {
	for {
		var b [4]byte
		_, err := rand.Read(b[:])
		if err != nil {
			return 0, err
		}
		
		// Get 27 bits (0-99999999)
		passcode := binary.LittleEndian.Uint32(b[:]) & 0x7FFFFFF
		
		// Check if valid
		if passcode <= 99999999 && !InvalidPasscodes[passcode] {
			return passcode, nil
		}
	}
}

// GenerateRandomSalt generates a 32-byte random salt and returns base64 encoded string
func GenerateRandomSalt() (string, error) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// GenerateRandomVerifier generates all SPAKE2+ parameters randomly
func GenerateRandomVerifier() (*VerifierParams, error) {
    // Generate random passcode
    passcode, err := GenerateRandomPasscode()
    if err != nil {
        return nil, fmt.Errorf("failed to generate passcode: %v", err)
    }

    // Generate random salt
    salt := make([]byte, 32)
    if _, err := rand.Read(salt); err != nil {
        return nil, fmt.Errorf("failed to generate salt: %v", err)
    }

    // Use fixed iterations = 1000
    var iterations uint16 = 1000

    // Generate verifier
    verifier, err := GenerateVerifier(passcode, salt, iterations)
    if err != nil {
        return nil, fmt.Errorf("failed to generate verifier: %v", err)
    }

    return &VerifierParams{
        Passcode:   passcode,
        Salt:       salt,
        SaltBase64: base64.StdEncoding.EncodeToString(salt),
        Iterations: iterations,
        Verifier:   verifier,
    }, nil
}