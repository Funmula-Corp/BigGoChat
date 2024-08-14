// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package utils

import (
	"bytes"
	"encoding/base64"
	"os"
	"testing"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var validTestLicense = []byte("eyJpZCI6InpvZ3c2NW44Z2lmajVkbHJoYThtYnUxcGl3IiwiaXNzdWVkX2F0IjoxNjg0Nzg3MzcxODY5LCJzdGFydHNfYXQiOjE2ODQ3ODczNzE4NjksImV4cGlyZXNfYXQiOjIwMDA0MDY1MzgwMDAsInNrdV9uYW1lIjoiUHJvZmVzc2lvbmFsIiwic2t1X3Nob3J0X25hbWUiOiJwcm9mZXNzaW9uYWwiLCJjdXN0b21lciI6eyJpZCI6InA5dW4zNjlhNjdnaW1qNHlkNmk2aWIzOXdoIiwibmFtZSI6Ik1hdHRlcm1vc3QiLCJlbWFpbCI6ImpvcmFtQG1hdHRlcm1vc3QuY29tIiwiY29tcGFueSI6Ik1hdHRlcm1vc3QifSwiZmVhdHVyZXMiOnsidXNlcnMiOjIwMDAwMCwibGRhcCI6dHJ1ZSwibGRhcF9ncm91cHMiOmZhbHNlLCJtZmEiOnRydWUsImdvb2dsZV9vYXV0aCI6dHJ1ZSwib2ZmaWNlMzY1X29hdXRoIjp0cnVlLCJjb21wbGlhbmNlIjpmYWxzZSwiY2x1c3RlciI6dHJ1ZSwibWV0cmljcyI6dHJ1ZSwibWhwbnMiOnRydWUsInNhbWwiOnRydWUsImVsYXN0aWNfc2VhcmNoIjp0cnVlLCJhbm5vdW5jZW1lbnQiOnRydWUsInRoZW1lX21hbmFnZW1lbnQiOmZhbHNlLCJlbWFpbF9ub3RpZmljYXRpb25fY29udGVudHMiOmZhbHNlLCJkYXRhX3JldGVudGlvbiI6ZmFsc2UsIm1lc3NhZ2VfZXhwb3J0IjpmYWxzZSwiY3VzdG9tX3Blcm1pc3Npb25zX3NjaGVtZXMiOmZhbHNlLCJjdXN0b21fdGVybXNfb2Zfc2VydmljZSI6ZmFsc2UsImd1ZXN0X2FjY291bnRzIjp0cnVlLCJndWVzdF9hY2NvdW50c19wZXJtaXNzaW9ucyI6dHJ1ZSwiaWRfbG9hZGVkIjpmYWxzZSwibG9ja190ZWFtbWF0ZV9uYW1lX2Rpc3BsYXkiOmZhbHNlLCJjbG91ZCI6ZmFsc2UsInNoYXJlZF9jaGFubmVscyI6ZmFsc2UsInJlbW90ZV9jbHVzdGVyX3NlcnZpY2UiOmZhbHNlLCJvcGVuaWQiOnRydWUsImVudGVycHJpc2VfcGx1Z2lucyI6dHJ1ZSwiYWR2YW5jZWRfbG9nZ2luZyI6dHJ1ZSwiZnV0dXJlX2ZlYXR1cmVzIjpmYWxzZX0sImlzX3RyaWFsIjp0cnVlLCJpc19nb3Zfc2t1IjpmYWxzZX0bEOVk2GdE1kSWKJ3dENWnkj0htY6QyXTtNA5hqnQ71Uc6teqXc7htHAxrnT/hV42xu+G24OMrAIsQtX4NjFSX6jvehIMRL5II3RPXYhHKUd2wruQ5ITEh1htFb5DgOJW3tvBdMmXt09nXjLRS1UYJ7ZsX3mU0uQndt7qfMriGAkk71veYuUJgztB3MsV7lRWB+8ZTp6WJ7RH+uWnuDspiA8B85mLnyuoCDokYksF2uIb+CtPGBTUB6qSOgxBBJxu5qftQXISCDAWY4O8lCrN3p5HCA/zf/rSRRNtet06QFobbjUDI4B7ZEAescKBKoHpP6nZPhg4KmhnkUi/o04ox")

func TestValidateLicenseBigGo(t *testing.T) {
	target := []byte("eyJpZCI6IlFxVFQ0RGQ2NVdxNmVuRnJaOHloU0d6SmtXIiwiaXNzdWVkX2F0IjoxNzIzNjI1MzM1MzQ1LCJzdGFydHNfYXQiOjE3MjM1NjQ4MDAwMDAsImV4cGlyZXNfYXQiOjE3NTUxMDA4MDAwMDAsImN1c3RvbWVyIjp7ImlkIjoiYmlnZ29fZGV2X3RlYW0iLCJuYW1lIjoiQiFnR28gRGV2ZWxvcG1lbnQiLCJlbWFpbCI6InJvb3RAbG9jYWxob3N0IiwiY29tcGFueSI6IkZ1bm11bGEifSwiZmVhdHVyZXMiOnsidXNlcnMiOjEwMDAwMDAsImxkYXAiOnRydWUsImxkYXBfZ3JvdXBzIjp0cnVlLCJtZmEiOnRydWUsImdvb2dsZV9vYXV0aCI6dHJ1ZSwib2ZmaWNlMzY1X29hdXRoIjp0cnVlLCJvcGVuaWQiOmZhbHNlLCJjb21wbGlhbmNlIjp0cnVlLCJjbHVzdGVyIjp0cnVlLCJtZXRyaWNzIjp0cnVlLCJtaHBucyI6dHJ1ZSwic2FtbCI6dHJ1ZSwiZWxhc3RpY19zZWFyY2giOnRydWUsImFubm91bmNlbWVudCI6dHJ1ZSwidGhlbWVfbWFuYWdlbWVudCI6dHJ1ZSwiZW1haWxfbm90aWZpY2F0aW9uX2NvbnRlbnRzIjp0cnVlLCJkYXRhX3JldGVudGlvbiI6dHJ1ZSwibWVzc2FnZV9leHBvcnQiOnRydWUsImN1c3RvbV9wZXJtaXNzaW9uc19zY2hlbWVzIjp0cnVlLCJjdXN0b21fdGVybXNfb2Zfc2VydmljZSI6dHJ1ZSwiZ3Vlc3RfYWNjb3VudHMiOnRydWUsImd1ZXN0X2FjY291bnRzX3Blcm1pc3Npb25zIjp0cnVlLCJpZF9sb2FkZWQiOnRydWUsImxvY2tfdGVhbW1hdGVfbmFtZV9kaXNwbGF5Ijp0cnVlLCJlbnRlcnByaXNlX3BsdWdpbnMiOnRydWUsImFkdmFuY2VkX2xvZ2dpbmciOnRydWUsImNsb3VkIjpmYWxzZSwic2hhcmVkX2NoYW5uZWxzIjp0cnVlLCJyZW1vdGVfY2x1c3Rlcl9zZXJ2aWNlIjp0cnVlLCJvdXRnb2luZ19vYXV0aF9jb25uZWN0aW9ucyI6dHJ1ZSwiZnV0dXJlX2ZlYXR1cmVzIjp0cnVlfSwic2t1X25hbWUiOiJCIWdHbyBDaGF0IExpY2Vuc2UiLCJza3Vfc2hvcnRfbmFtZSI6IkIhZ0dvIENoYXQgTGljZW5zZSIsImlzX3RyaWFsIjpmYWxzZSwiaXNfZ292X3NrdSI6ZmFsc2UsInNpZ251cF9qd3QiOm51bGx9guz0Sxp3PcxWngKqWRLB9Ve1HCzyvzLExWRafZovMU6LxcW8wbT7iWvJ5ByOOLI/efA7pA7wo1vIGF3vc2aiCjxa0eWWbu2hdyAfU7HD7Ob+WrpgfCcaMSiEhEGbg1p6uXGK0C6DAjCxOlHnPFs5zwZQ4smUPml47aVdjRsBUTt6hvRmk7kXvxiYyhptPnj218HSHV8Ka6muQEPvz++oHGRVq+4HWi7MivWpAp5MWczS3hDsCUJs3+i3IoqHwimimnFyh0eXhVsyFu0Cw0DHyYOhRPq7OTESZd0rNa2mfcRtwCDfUuYTbE6UnAjiBOVuWd1M1zqp5ig0q2Um7ykq6Q==")
	_, err := LicenseValidator.ValidateLicense(target)
	require.NoError(t, err)
}

func TestValidateLicense(t *testing.T) {
	t.Run("should fail with junk data", func(t *testing.T) {
		b1 := []byte("junk")
		_, err := LicenseValidator.ValidateLicense(b1)
		require.Error(t, err, "should have failed - bad license")

		b2 := []byte("junkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunkjunk")
		_, err = LicenseValidator.ValidateLicense(b2)
		require.Error(t, err, "should have failed - bad license")
	})

	t.Run("should not panic on shorter than expected input", func(t *testing.T) {
		var licenseData bytes.Buffer
		var inputData []byte

		for i := 0; i < 255; i++ {
			inputData = append(inputData, 'A')
		}
		inputData = append(inputData, 0x00)

		encoder := base64.NewEncoder(base64.StdEncoding, &licenseData)
		_, err := encoder.Write(inputData)
		require.NoError(t, err)
		err = encoder.Close()
		require.NoError(t, err)

		str, err := LicenseValidator.ValidateLicense(licenseData.Bytes())
		require.Error(t, err)
		require.Empty(t, str)
	})

	t.Run("should not panic with input filled of null terminators", func(t *testing.T) {
		var licenseData bytes.Buffer
		var inputData []byte

		for i := 0; i < 256; i++ {
			inputData = append(inputData, 0x00)
		}

		encoder := base64.NewEncoder(base64.StdEncoding, &licenseData)
		_, err := encoder.Write(inputData)
		require.NoError(t, err)
		err = encoder.Close()
		require.NoError(t, err)

		str, err := LicenseValidator.ValidateLicense(licenseData.Bytes())
		require.Error(t, err)
		require.Empty(t, str)
	})

	t.Run("should reject invalid license in test service environment", func(t *testing.T) {
		os.Setenv("MM_SERVICEENVIRONMENT", model.ServiceEnvironmentTest)
		defer os.Unsetenv("MM_SERVICEENVIRONMENT")

		str, err := LicenseValidator.ValidateLicense(nil)
		require.Error(t, err)
		require.Empty(t, str)
	})

	t.Run("should validate valid test license in test service environment", func(t *testing.T) {
		os.Setenv("MM_SERVICEENVIRONMENT", model.ServiceEnvironmentTest)
		defer os.Unsetenv("MM_SERVICEENVIRONMENT")

		str, err := LicenseValidator.ValidateLicense(validTestLicense)
		require.NoError(t, err)
		require.NotEmpty(t, str)
	})

	t.Run("should reject valid test license in production service environment", func(t *testing.T) {
		os.Setenv("MM_SERVICEENVIRONMENT", model.ServiceEnvironmentProduction)
		defer os.Unsetenv("MM_SERVICEENVIRONMENT")

		str, err := LicenseValidator.ValidateLicense(validTestLicense)
		require.Error(t, err)
		require.Empty(t, str)
	})
}

func TestGetLicenseFileLocation(t *testing.T) {
	fileName := GetLicenseFileLocation("")
	require.NotEmpty(t, fileName, "invalid default file name")

	fileName = GetLicenseFileLocation("mattermost.mattermost-license")
	require.Equal(t, fileName, "mattermost.mattermost-license", "invalid file name")
}

func TestGetLicenseFileFromDisk(t *testing.T) {
	t.Run("missing file", func(t *testing.T) {
		fileBytes := GetLicenseFileFromDisk("thisfileshouldnotexist.mattermost-license")
		assert.Empty(t, fileBytes, "invalid bytes")
	})

	t.Run("not a license file", func(t *testing.T) {
		f, err := os.CreateTemp("", "TestGetLicenseFileFromDisk")
		require.NoError(t, err)
		defer os.Remove(f.Name())
		os.WriteFile(f.Name(), []byte("not a license"), 0777)

		fileBytes := GetLicenseFileFromDisk(f.Name())
		require.NotEmpty(t, fileBytes, "should have read the file")

		_, err = LicenseValidator.ValidateLicense(fileBytes)
		assert.Error(t, err, "should have been an invalid file")
	})
}
