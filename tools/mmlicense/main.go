package main

/*
TODO:
+ connect to postgres
+ store license in database
*/

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"mmlicense/cert"
	"mmlicense/license"
)

func main() {
	if licenseBuffer, err := json.Marshal(license.Config); err != nil {
		log.Fatalln("error marshalling license model")
	} else {
		if license.PrintJSON {
			fmt.Println("License Info")
			fmt.Println("StartsAt:", time.UnixMilli(license.Config.StartsAt))
			fmt.Println("ExpiresAt:", time.UnixMilli(license.Config.ExpiresAt))
			fmt.Println("License JSON:")
			fmt.Println(string(licenseBuffer))
		} else {
			signedLicense := cert.SignLicense(licenseBuffer)
			signedLicenseBuffer := make([]byte, base64.StdEncoding.EncodedLen(len(signedLicense)))
			base64.StdEncoding.Encode(signedLicenseBuffer, signedLicense)
			cert.ValidateLicense(signedLicenseBuffer)
			fmt.Println(string(signedLicenseBuffer))
		}
	}
}
