package main

/*
TODO:
+ connect to postgres
+ store license in database
*/

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"mmlicense/cert"
	"mmlicense/license"
	"mmlicense/storage"

	"github.com/mattermost/mattermost/server/public/model"
)

var (
	Quiet bool = false

	Activate  bool   = false
	Insert    bool   = false
	LicenseId string = ""
)

func init() {
	flag.BoolVar(&Quiet, "q", false, "quiet mode - print only encoded license")
	flag.BoolVar(&Activate, "activate", false, "quiet mode - print only encoded license")
	flag.BoolVar(&Insert, "insert", false, "insert the created license into the database")
	flag.StringVar(&LicenseId, "license_id", LicenseId, "set license id (!SKIPS LICENSE GENERATION!)")
	flag.Parse()
}

func main() {
	if LicenseId == "" {
		NewLicense()
	}

	if Activate {
		storage.ActivateLicense(LicenseId)
	}
}

func NewLicense() {
	licenseConfig := license.New()

	if licenseBuffer, err := json.Marshal(licenseConfig); err != nil {
		log.Fatalln("error marshalling license model")
	} else {
		if !Quiet {
			PrintDetails(licenseConfig)
		}

		signedLicense := cert.SignLicense(licenseBuffer)
		signedLicenseString := base64.StdEncoding.EncodeToString(signedLicense)
		cert.ValidateLicense([]byte(signedLicenseString))
		fmt.Println(signedLicenseString)

		LicenseId = licenseConfig.Id
		if Insert {
			storage.InsertLicense(licenseConfig.Id, licenseConfig.IssuedAt, []byte(signedLicenseString))
		}
	}
}

func PrintDetails(licenceConfig *model.License) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("\t", "  ")

	fmt.Println("")
	fmt.Println("==============================================")
	fmt.Println("==   Funmula Mattermost Licence Generator   ==")
	fmt.Println("==============================================")
	fmt.Println("")

	fmt.Println("License Details:")
	fmt.Println("")
	fmt.Println("\tID           :", licenceConfig.Id)
	fmt.Println("\tIsGovSku     :", licenceConfig.IsGovSku)
	fmt.Println("\tIsTrial      :", licenceConfig.IsTrial)
	fmt.Println("\tSkuName      :", licenceConfig.SkuName)
	fmt.Println("\tSkuShortName :", licenceConfig.SkuShortName)
	fmt.Println("\tStartsAt     :", time.UnixMilli(licenceConfig.StartsAt))
	fmt.Println("\tExpiresAt    :", time.UnixMilli(licenceConfig.ExpiresAt))
	fmt.Println("")
	fmt.Println("License Customer:")
	fmt.Printf("\r\n\t")
	encoder.Encode(licenceConfig.Customer)
	fmt.Println("")
	fmt.Println("License Features:")
	fmt.Printf("\r\n\t")
	encoder.Encode(licenceConfig.Features)
	fmt.Println("")
	fmt.Println("Encoded License:")
}
