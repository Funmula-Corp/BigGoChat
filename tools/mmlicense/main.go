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
	"mmlicense/cert"
	"mmlicense/gen"
	"os"
	"path/filepath"
	"time"

	"github.com/mattermost/mattermost/server/public/model"
)

var (
	fileName   = "license.json"
	workDir, _ = os.Getwd()

	configPath = filepath.Join(workDir, fileName)
	license    = &model.License{}
)

func init() {
	flag.StringVar(&configPath, "config", configPath, "path to the license config file")
	flag.Parse()

	if buffer, err := os.ReadFile(configPath); os.IsNotExist(err) {
		log.Fatalf("%s file not found", filepath.Base(configPath))
	} else if err != nil {
		log.Fatalln(err)
	} else {
		if err = json.Unmarshal(buffer, license); err != nil {
			log.Fatalln(err)
		}
	}

	license.Id = gen.NewLicenseID()
	license.IssuedAt = time.Now().UnixMilli()
	license.StartsAt = time.Now().UnixMilli()
	license.ExpiresAt = time.Now().Add(time.Hour * 24 * 365).UnixMilli()
}

func main() {
	var (
		err error
		lic []byte
		sig []byte
		buf []byte
	)

	if lic, err = json.Marshal(license); err != nil {
		log.Fatalln("error marshalling license model")
	}
	sig = cert.SignBuffer(lic)
	buf = append(lic, sig...)

	dstBuffer := make([]byte, base64.StdEncoding.EncodedLen(len(buf)))
	base64.StdEncoding.Encode(dstBuffer, buf)

	cert.ValidateLicense(dstBuffer)
	fmt.Println(string(dstBuffer))
}
