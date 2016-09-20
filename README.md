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

### Encrypt File

```bash
go run main.go encrypt --key 4f2f62e1-41e0-49e2-8da4-3a7ec511f498 --file "main.go" --name "testKey"
```

Reference: http://docs.aws.amazon.com/kms/latest/developerguide/workflow.html
