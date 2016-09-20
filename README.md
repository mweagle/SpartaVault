# SpartaVault
Part of the Sparta toolkit - use AWS KMS to securely manage and commit secrets to source control.

## Usage

```
Matts-MBP:SpartaVault mweagle$ go run main.go encrypt help
Error: Provide either string or path plaintext input value
Usage:
  SpartaVault encrypt [flags]

Flags:
  -f, --file string    Path to file whose contents should be encrypted
  -k, --key string     AWS KMS Keyname (ARN) to use for encryption
  -n, --name string    go Property name for encrypted value
  -v, --value string   String value to encrypt

Provide either string or path plaintext input value
exit status 255
```

## Examples

### Encrypt String

```bash
go run main.go encrypt --key 4f2f62e1-41e0-49e2-8da4-3a7ec511f498 --value "Hello World" --name "testKey"
```

Output:

```golang
----- BEGIN SNIPPET -----

package main

import (
	"fmt"
	spartaVault "github.com/mweagle/SpartaVault/encrypt"
)

var testKey = &spartaVault.KMSEncryptedValue{
	KMSKeyARNOrGuid: "4f2f62e1-41e0-49e2-8da4-3a7ec511f498",
	PropertyName:    "testKey",
	Key:             "AQEDAHi8zBTBrgXJ4OyfnaJ8C9B2H/WAF54D9vPaarH9Dob2wwAAAH4wfAYJKoZIhvcNAQcGoG8wbQIBADBoBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDB2uCmIx46f45/7wKgIBEIA71iBmbCr8EYuX8XGeAy1Qpus94Q5HXSwBQoH9A77jzJEnNgu+FpP7wi94qMzBBvAU3+mQbf5S39RxUo0=",
	Nonce:           "lqoVNKQLDlDq8Ij4",
	Value:           "F2DN/7Looc8ajOO8UJdp4B0mSL7UMvfRa9No",
}

// Usage:
// func main() {
// 	plaintextValue, _ := testKey.Decrypt()
// 	fmt.Printf("Decrypted: %s\n", plaintextValue)
// }

-----  END SNIPPET  -----

```

### Encrypt File

```bash
go run main.go encrypt --key 4f2f62e1-41e0-49e2-8da4-3a7ec511f498 --file "main.go" --name "testKey"
```

Reference: http://docs.aws.amazon.com/kms/latest/developerguide/workflow.html
