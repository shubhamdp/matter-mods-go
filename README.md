# matter-mods-go

- Can generate spake2p verifier required for PASE
- Can generate the manual pairing code

#### Get the package
```bash
go get github.com/shubhamdp/matter-setup-payload-go/spake2p
go get github.com/shubhamdp/matter-setup-payload-go/manualcode
```

## Usage in code

```go
package main

import (
    "fmt"
    "encoding/base64"
    "github.com/shubhamdp/matter-mods-go/spake2p"
    "github.com/shubhamdp/matter-mods-go/manualcode"
)

func main() {
    // Generate spake2p verifier
	params, err := spake2p.GenerateRandomVerifier()
	if err != nil {
		fmt.Printf("Failed to generate random verifier: %v\n", err)
		return
	}
	fmt.Printf("Passcode: %d\n", params.Passcode)
	fmt.Printf("Salt: %s\n", params.SaltBase64)
	fmt.Printf("Iterations: %d\n", params.Iterations)
	fmt.Printf("Verifier: %s\n", base64.StdEncoding.EncodeToString(params.Verifier))

	fmt.Println("-------------------")

	// lets use the same parameters to generate the manual code
	payload := manualcode.PayloadContents{
		SetUpPINCode:      params.Passcode,
		Discriminator:     3431,
		VendorID:          0x1317,
		ProductID:         0x0002,
		CommissioningFlow: 0,
	}

  	g := &manualcode.ManualSetupPayloadGenerator{
		PayloadContents: payload,
	}

  	code, err := g.GenerateManualcode()
	if err != nil {
   		fmt.Printf("Failed to generate manual code: %v\n", err)
	}
	fmt.Println(code)
}

```

#### Running unit tests
```bash
go test -v ./...
```

#### invidvidually running in main.go

- copy the code above in main.go
- the try these
```
go mod init main
go mod tidy
go run main.go
```
