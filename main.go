package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"go.mozilla.org/pkcs7"
	"go.step.sm/cli-utils/ui"
	"go.step.sm/crypto/pemutil"
	"howett.net/plist"
)

var (
	certPath         string
	keyPath          string
	mobileConfigPath string
	password         string
)

func init() {
	flag.StringVar(&certPath, "cert", "", "Full path to the certificate to sign with")
	flag.StringVar(&keyPath, "key", "", "Full path to the private key to sign with")
	flag.StringVar(&mobileConfigPath, "file", "", "Full path to .mobileconfig to sign")
	flag.StringVar(&password, "password", "", "Password for the private key (optional)")
}

func main() {

	processFlags()

	f, err := os.Open(mobileConfigPath)
	fatalIf(err)
	defer f.Close()

	// decode the plist contents
	var r map[string]interface{}
	decoder := plist.NewDecoder(f)
	decoder.Format = plist.XMLFormat
	err = decoder.Decode(&r)
	fatalIf(err)

	encoder := plist.NewEncoder(os.Stdout)
	encoder.Indent("\t")
	//encoder.Encode(r) // encodes the mobileconfig contents and writes to stdout
	//fatalIf(err)

	// reencode the plist contents
	var buf bytes.Buffer
	encoder = plist.NewEncoder(&buf)
	encoder.Indent("\t")
	err = encoder.Encode(r) // encodes the mobileconfig contents and writes to the buffer
	fatalIf(err)

	// prepare signing operations
	cert, err := pemutil.ReadCertificate(certPath)
	fatalIf(err)

	options := []pemutil.Options{
		pemutil.WithPasswordPrompt(
			fmt.Sprintf("Please enter the password to decrypt %q", keyPath),
			func(s string) ([]byte, error) {
				// prompt for password if not provided as flag
				return ui.PromptPassword(s, ui.WithValue(password))
			}),
	}

	key, err := pemutil.Read(keyPath, options...)
	fatalIf(err)

	// sign the plist bytes
	unsigned := buf.Bytes()
	sd, err := pkcs7.NewSignedData(unsigned)
	fatalIf(err)

	err = sd.AddSigner(cert, key, pkcs7.SignerInfoConfig{})
	fatalIf(err)

	signed, err := sd.Finish()
	fatalIf(err)

	// write the signed file next to the original file
	ext := path.Ext(mobileConfigPath)
	signedMobileConfigPath := mobileConfigPath[:len(mobileConfigPath)-len(ext)] + ".signed" + ext
	err = os.WriteFile(signedMobileConfigPath, signed, 0644)
	fatalIf(err)

	log.Printf("Written signed mobileconfig to %q", signedMobileConfigPath)
}

func processFlags() {
	flag.Parse()
	if certPath == "" {
		log.Fatal("-cert flag is required")
	}
	if keyPath == "" {
		log.Fatal("-key flag is required")
	}
	if mobileConfigPath == "" {
		log.Fatal("-file flag is required")
	}
}

func fatalIf(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
