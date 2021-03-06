package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/spf13/cobra"
)

/******************************************************************************/
// CONSTANTS
const encryptedValueCodeTemplate = `
package main

import (
	"fmt"
	spartaVault "github.com/mweagle/SpartaVault/encrypt"
)

var {{ .PropertyName }} = &spartaVault.KMSEncryptedValue{
	KMSKeyARNOrGUID:      "{{ .KMSKeyARNOrGuid}}",
	PropertyName:  				"{{ .PropertyName }}",
	Key:          				"{{ .Key }}",
	Nonce:        				"{{ .Nonce }}",
	Value:       					"{{ .Value }}",
	Created: 							"{{ .Created }}",
}

// Usage:
// func main() {
// 	plaintextValue, _ := testKey.Decrypt(nil)
// 	fmt.Printf("Decrypted: %s\n", plaintextValue)
// }
`

const encryptedValueSnippet = `
----- BEGIN SNIPPET -----

%s

-----  END SNIPPET  -----
`

/******************************************************************************/
// Global options
type encryptOptionsStruct struct {
	KMSKeyName   string `valid:"required,matches(\\w+)"`
	PropertyName string `valid:"required,matches(\\w+)"`
	Value        string `valid:"matches(\\w+)"`
	FilePath     string `valid:"matches(\\w+)"`
}

// OptionsGlobal stores the global command line options
var encryptOptions encryptOptionsStruct

/******************************************************************************/
// Init
func init() {
	encryptCmd.Flags().StringVarP(&encryptOptions.KMSKeyName, "key", "k", "", "AWS KMS Keyname (ARN or GUID) to use for encryption")
	encryptCmd.Flags().StringVarP(&encryptOptions.PropertyName, "name", "n", "", "go Property name for encrypted value")
	encryptCmd.Flags().StringVarP(&encryptOptions.Value, "value", "v", "", "String value to encrypt")
	encryptCmd.Flags().StringVarP(&encryptOptions.FilePath, "file", "f", "", "Path to file whose contents should be encrypted")

	RootCmd.AddCommand(encryptCmd)
}

/******************************************************************************/
// Types

// KMSEncryptedValue represents the encrypted secret value
type KMSEncryptedValue struct {
	KMSKeyARNOrGUID string
	PropertyName    string
	Key             string
	Nonce           string
	Value           string
	Created         string
}

// Decrypt attempts to decrypt the given KMSEncryptedValue using the
// optional awsSession. If the awsSession value is nil, a default
// session will be used
func (kmsValue *KMSEncryptedValue) Decrypt(awsSession *session.Session) ([]byte, error) {
	decryptSession := awsSession
	if nil == decryptSession {
		sess, sessionError := session.NewSession()
		if nil != sessionError {
			return nil, sessionError
		}
		decryptSession = sess
	}

	// Decrypt the one off key
	decodedKey, decodedKeyErr := base64.StdEncoding.DecodeString(kmsValue.Key)
	if nil != decodedKeyErr {
		return nil, decodedKeyErr
	}
	kmsSvc := kms.New(decryptSession)
	params := &kms.DecryptInput{
		CiphertextBlob: decodedKey,
	}
	decryptResp, decryptRespErr := kmsSvc.Decrypt(params)
	if nil != decryptRespErr {
		return nil, decryptRespErr
	}
	// Decode the contents && decrypt them with the key...
	aesBlock, aesBlockErr := aes.NewCipher(decryptResp.Plaintext)
	if nil != aesBlockErr {
		return nil, aesBlockErr
	}

	aesGCM, aesGCMErr := cipher.NewGCM(aesBlock)
	if nil != aesGCMErr {
		return nil, aesGCMErr
	}

	decodedNonce, decodedNonceErr := base64.StdEncoding.DecodeString(kmsValue.Nonce)
	if nil != decodedNonceErr {
		return nil, decodedNonceErr
	}
	decodedValue, decodedValueErr := base64.StdEncoding.DecodeString(kmsValue.Value)
	if nil != decodedValueErr {
		return nil, decodedValueErr
	}

	plaintext, plaintextErr := aesGCM.Open(nil, decodedNonce, decodedValue, nil)
	if plaintextErr != nil {
		return nil, plaintextErr
	}
	return plaintext, nil
}

func newEncryptedValue(keyARN string, PropertyName string, content io.Reader) (*KMSEncryptedValue, error) {
	sess, sessionError := session.NewSession()
	if nil != sessionError {
		return nil, sessionError
	}

	kmsSvc := kms.New(sess)
	params := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(keyARN), // Required
		KeySpec: aws.String("AES_256"),
	}
	generateResp, generateRespErr := kmsSvc.GenerateDataKey(params)
	if nil != generateRespErr {
		return nil, generateRespErr
	}
	// Encrypt some data...
	aesBlock, aesBlockErr := aes.NewCipher(generateResp.Plaintext)
	if nil != aesBlockErr {
		return nil, aesBlockErr
	}

	aesGCM, aesGCMErr := cipher.NewGCM(aesBlock)
	if nil != aesGCMErr {
		return nil, aesGCMErr
	}

	plaintext, plaintextErr := ioutil.ReadAll(content)
	if nil != plaintextErr {
		return nil, plaintextErr
	}
	nonce := make([]byte, 12)
	if _, nonceErr := io.ReadFull(rand.Reader, nonce); nonceErr != nil {
		return nil, nonceErr
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)
	encryptedValue := &KMSEncryptedValue{
		KMSKeyARNOrGUID: keyARN,
		Key:             base64.StdEncoding.EncodeToString(generateResp.CiphertextBlob),
		Nonce:           base64.StdEncoding.EncodeToString(nonce),
		PropertyName:    PropertyName,
		Value:           base64.StdEncoding.EncodeToString(ciphertext),
		Created:         time.Now().Format(time.RFC3339),
	}
	return encryptedValue, nil
}

func outputEncryptedGolang(encryptedKey *KMSEncryptedValue) error {
	template, templateErr := template.New("KMSSnippet").Parse(encryptedValueCodeTemplate)
	if nil != templateErr {
		return templateErr
	}
	var doc bytes.Buffer
	executeErr := template.Execute(&doc, *encryptedKey)
	if nil != executeErr {
		return executeErr
	}
	formatted, formattedErr := format.Source(doc.Bytes())
	if nil != formattedErr {
		return formattedErr
	}
	fmt.Printf(encryptedValueSnippet, string(formatted))
	return nil
}

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "CLI tool to produce AWS KMS secrets",
	Long: `Use AWS KMS to produce envelope encrypted secrets that can be committed to source code.

	The command line tool produces a Go language snippet that can be directly
	used in source code.

	See:
		http://docs.aws.amazon.com/kms/latest/developerguide/workflow.html
	for more details.
	`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if "" == encryptOptions.Value && "" == encryptOptions.FilePath {
			return fmt.Errorf("Provide either --value or --file plaintext input value")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		var input io.Reader
		if "" != encryptOptions.Value {
			input = strings.NewReader(encryptOptions.Value)
		} else {
			fileReader, fileReaderErr := os.Open(encryptOptions.FilePath)
			if nil != fileReaderErr {
				return fileReaderErr
			}
			input = fileReader
		}
		kmsValue, kmsValueErr := newEncryptedValue(encryptOptions.KMSKeyName,
			encryptOptions.PropertyName,
			input)
		if nil != kmsValueErr {
			return kmsValueErr
		}
		// Validate the decryption
		_, decryptedValueErr := kmsValue.Decrypt(nil)
		if nil != decryptedValueErr {
			return decryptedValueErr
		}
		return outputEncryptedGolang(kmsValue)
	},
}
