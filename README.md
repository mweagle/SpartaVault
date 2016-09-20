# SpartaVault
Part of the Sparta toolkit - use AWS KMS to encrypt secrets to Go variables that can be committed to SCM.

## Usage

1. Create an [AWS KMS key](http://docs.aws.amazon.com/kms/latest/developerguide/create-keys.html)
2. Encrypt plaintext (either string or filepath) using the new key ARN or GUID

```
$ go run main.go encrypt help
Error: Provide either --value or --file plaintext input value
Usage:
  SpartaVault encrypt [flags]

Flags:
  -f, --file string    Path to file whose contents should be encrypted
  -k, --key string     AWS KMS Keyname (ARN or GUID) to use for encryption
  -n, --name string    go Property name for encrypted value
  -v, --value string   String value value to encrypt

Provide either --value or --file plaintext input value
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
	Key:             "AQEDAHi8zBTBrgXJ4OyfnaJ8C9B2H/WAF54D9vPaarH9Dob2wwAAAH4wfAYJKoZIhvcNAQcGoG8wbQIBADBoBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDAWYup2u/ZdD4VRV3gIBEIA76z9NVXE3m8AhK6SdT8yEOmu0pXf3CBcUJ4DSAiwYQt4Y3mDePdLfGlkTbratRExo33Zzse8m/G4G6iI=",
	Nonce:           "VDS+3LffkcSUGEpc",
	Value:           "U4RQWOVsYyGiaJ2VhGXeWhO5Gd3+6uhaiqcg",
	Created:         "2016-09-20T05:57:42-07:00",
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
	Key:             "AQEDAHi8zBTBrgXJ4OyfnaJ8C9B2H/WAF54D9vPaarH9Dob2wwAAAH4wfAYJKoZIhvcNAQcGoG8wbQIBADBoBgkqhkiG9w0BBwEwHgYJYIZIAWUDBAEuMBEEDM5EV8Mnf/vCEvVUqQIBEIA7/2QGAOg2VV/AV9+X8Ae9flkraLMek8cOZ5R0zSEPNCGEXnwjqwHkqICK6nYMtmTKGu7qD7rf/nrOtVA=",
	Nonce:           "k8uPquOblLxKmjCO",
	Value:           "VH2QwL43aTf52ztt7lrf2kL2CBwh2cROd0efI7q+/NrP4+FQfXUhKq8uYRHy7mQdds2ZHo7EZG8EQ4Bsy4a4xRq0fa8q/SLdj7aRbzqwjg44hbO7vBl6WnQQGGkqHRM12jdjwK1x0sy0eZ2Nln2sGQcV6+RseDY=",
	Created:         "2016-09-18T03:58:12-07:00",
}

// Usage:
// func main() {
// 	plaintextValue, _ := testKey.Decrypt()
// 	fmt.Printf("Decrypted: %s\n", plaintextValue)
// }


-----  END SNIPPET  -----
```

Reference: http://docs.aws.amazon.com/kms/latest/developerguide/workflow.html

