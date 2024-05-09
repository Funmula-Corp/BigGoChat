package license

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"mmlicense/env"
	"mmlicense/gen"

	"github.com/mattermost/mattermost/server/public/model"
)

var (
	Config     *model.License = &model.License{Customer: &model.Customer{}, Features: &model.Features{}}
	configPath *string        = env.String("MM_LICENSE_CONFIG")
	PrintJSON  bool           = false

	expiresIn     int    = 365
	startsAt      string = ""
	startsAtDate  string = ""
	startsAtTime  string = ""
	expiresAt     string = ""
	expiresAtDate string = ""
	expiresAtTime string = ""

	skuName      *string = env.String("MM_LICENSE_SKU_NAME")
	skuShortName *string = env.String("MM_LICENSE_SKU_SHORT_NAME")

	isTrial  *bool = env.Bool("MM_LICENSE_IS_TRIAL")
	isGovSku *bool = env.Bool("MM_LICENSE_IS_GOVERNMENT")

	customerID      *string = env.String("MM_LICENSE_CUSTOMER_ID")
	customerName    *string = env.String("MM_LICENSE_CUSTOMER_NAME")
	customerEmail   *string = env.String("MM_LICENSE_CUSTOMER_EMAIL")
	customerCompany *string = env.String("MM_LICENSE_CUSTOMER_COMPANY")

	featureUsers                     *int  = env.Int("MM_LICENSE_FEATURE_USERS")
	featureLDAP                      *bool = env.Bool("MM_LICENSE_FEATURE_LDAP")
	featureLDAPGroups                *bool = env.Bool("MM_LICENSE_FEATURE_LDAPGROUPS")
	featureMFA                       *bool = env.Bool("MM_LICENSE_FEATURE_MFA")
	featureGoogleOAuth               *bool = env.Bool("MM_LICENSE_FEATURE_GOOGLEOAUTH")
	featureOffice365OAuth            *bool = env.Bool("MM_LICENSE_FEATURE_OFFICE365OAUTH")
	featureOpenId                    *bool = env.Bool("MM_LICENSE_FEATURE_OPENID")
	featureCompliance                *bool = env.Bool("MM_LICENSE_FEATURE_COMPLIANCE")
	featureCluster                   *bool = env.Bool("MM_LICENSE_FEATURE_CLUSTER")
	featureMetrics                   *bool = env.Bool("MM_LICENSE_FEATURE_METRICS")
	featureMHPNS                     *bool = env.Bool("MM_LICENSE_FEATURE_MHPNS")
	featureSAML                      *bool = env.Bool("MM_LICENSE_FEATURE_SAML")
	featureElasticsearch             *bool = env.Bool("MM_LICENSE_FEATURE_ELASTICSEARCH")
	featureAnnouncement              *bool = env.Bool("MM_LICENSE_FEATURE_ANNOUNCEMENT")
	featureThemeManagement           *bool = env.Bool("MM_LICENSE_FEATURE_THEMEMANAGEMENT")
	featureEmailNotificationContents *bool = env.Bool("MM_LICENSE_FEATURE_EMAILNOTIFICATIONCONTENTS")
	featureDataRetention             *bool = env.Bool("MM_LICENSE_FEATURE_DATARETENTION")
	featureMessageExport             *bool = env.Bool("MM_LICENSE_FEATURE_MESSAGEEXPORT")
	featureCustomPermissionsSchemes  *bool = env.Bool("MM_LICENSE_FEATURE_CUSTOMPERMISSIONSSCHEMES")
	featureCustomTermsOfService      *bool = env.Bool("MM_LICENSE_FEATURE_CUSTOMTERMSOFSERVICE")
	featureGuestAccounts             *bool = env.Bool("MM_LICENSE_FEATURE_GUESTACCOUNTS")
	featureGuestAccountsPermissions  *bool = env.Bool("MM_LICENSE_FEATURE_GUESTACCOUNTSPERMISSIONS")
	featureIDLoadedPushNotifications *bool = env.Bool("MM_LICENSE_FEATURE_IDLOADEDPUSHNOTIFICATIONS")
	featureLockTeammateNameDisplay   *bool = env.Bool("MM_LICENSE_FEATURE_LOCKTEAMMATENAMEDISPLAY")
	featureEnterprisePlugins         *bool = env.Bool("MM_LICENSE_FEATURE_ENTERPRISEPLUGINS")
	featureAdvancedLogging           *bool = env.Bool("MM_LICENSE_FEATURE_ADVANCEDLOGGING")
	featureCloud                     *bool = env.Bool("MM_LICENSE_FEATURE_CLOUD")
	featureSharedChannels            *bool = env.Bool("MM_LICENSE_FEATURE_SHAREDCHANNELS")
	featureRemoteClusterService      *bool = env.Bool("MM_LICENSE_FEATURE_REMOTECLUSTERSERVICE")
	featureOutgoingOAuthConnections  *bool = env.Bool("MM_LICENSE_FEATURE_OUTGOINGOAUTHCONNECTIONS")
	featureFutureFeatures            *bool = env.Bool("MM_LICENSE_FEATURE_FUTUREFEATURES")
)

func init() {
	flag.StringVar(configPath, "config", *configPath, "path to the license config file (env: MM_LICENSE_CONFIG)")

	flag.StringVar(&startsAt, "starts-at", startsAt, fmt.Sprintf("set the license start date and time (example: \"%s\")", time.Now().Format(time.DateTime)))
	flag.StringVar(&startsAtDate, "starts-at-date", startsAtDate, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.StringVar(&startsAtTime, "starts-at-time", startsAtTime, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.TimeOnly)))
	flag.StringVar(&expiresAt, "expires-at", expiresAt, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.StringVar(&expiresAtDate, "expires-at-date", expiresAtDate, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.StringVar(&expiresAtTime, "expires-at-time", expiresAtTime, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.IntVar(&expiresIn, "days", expiresIn, "set expiration days for the license")

	flag.StringVar(skuName, "skuName", *skuName, "set the license sku name (env: MM_LICENSE_SKU_NAME)")
	flag.StringVar(skuShortName, "skuShortName", *skuShortName, "set the license sku short name (env: MM_LICENSE_SKU_SHORT_NAME)")

	flag.BoolVar(isTrial, "IsTrial", *isTrial, "set is trial license flag (env: MM_LICENSE_IS_TRIAL)")
	flag.BoolVar(isGovSku, "IsGovSku", *isGovSku, "set is government license flag (env: MM_LICENSE_IS_GOVERNMENT)")

	flag.StringVar(customerID, "ID", *customerID, "set license customer id (env: MM_LICENSE_CUSTOMER_ID)")
	flag.StringVar(customerName, "Name", *customerName, "set license customer name (env: MM_LICENSE_CUSTOMER_NAME)")
	flag.StringVar(customerEmail, "Email", *customerEmail, "set license customer email (env: MM_LICENSE_CUSTOMER_EMAIL)")
	flag.StringVar(customerCompany, "Company", *customerCompany, "set license customer comany name (env: MM_LICENSE_CUSTOMER_COMPANY)")

	flag.IntVar(featureUsers, "Users", *featureUsers, "enable license feature \"Users\" (env: MM_LICENSE_FEATURE_USERS)")
	flag.BoolVar(featureLDAP, "LDAP", *featureLDAP, "enable license feature \"LDAP\" (env: MM_LICENSE_FEATURE_LDAP)")
	flag.BoolVar(featureLDAPGroups, "LDAPGroups", *featureLDAPGroups, "enable license feature \"LDAPGroups\" (env: MM_LICENSE_FEATURE_LDAPGROUPS)")
	flag.BoolVar(featureMFA, "MFA", *featureMFA, "enable license feature \"MFA\" (env: MM_LICENSE_FEATURE_MFA)")
	flag.BoolVar(featureGoogleOAuth, "GoogleOAuth", *featureGoogleOAuth, "enable license feature \"GoogleOAuth\" (env: MM_LICENSE_FEATURE_GOOGLEOAUTH)")
	flag.BoolVar(featureOffice365OAuth, "Office365OAuth", *featureOffice365OAuth, "enable license feature \"Office365OAuth\" (env: MM_LICENSE_FEATURE_OFFICE365OAUTH)")
	flag.BoolVar(featureOpenId, "OpenId", *featureOpenId, "enable license feature \"OpenId\" (env: MM_LICENSE_FEATURE_OPENID)")
	flag.BoolVar(featureCompliance, "Compliance", *featureCompliance, "enable license feature \"Compliance\" (env: MM_LICENSE_FEATURE_COMPLIANCE)")
	flag.BoolVar(featureCluster, "Cluster", *featureCluster, "enable license feature \"Cluster\" (env: MM_LICENSE_FEATURE_CLUSTER)")
	flag.BoolVar(featureMetrics, "Metrics", *featureMetrics, "enable license feature \"Metrics\" (env: MM_LICENSE_FEATURE_METRICS)")
	flag.BoolVar(featureMHPNS, "MHPNS", *featureMHPNS, "enable license feature \"MHPNS\" (env: MM_LICENSE_FEATURE_MHPNS)")
	flag.BoolVar(featureSAML, "SAML", *featureSAML, "enable license feature \"SAML\" (env: MM_LICENSE_FEATURE_SAML)")
	flag.BoolVar(featureElasticsearch, "Elasticsearch", *featureElasticsearch, "enable license feature \"Elasticsearch\" (env: MM_LICENSE_FEATURE_ELASTICSEARCH)")
	flag.BoolVar(featureAnnouncement, "Announcement", *featureAnnouncement, "enable license feature \"Announcement\" (env: MM_LICENSE_FEATURE_ANNOUNCEMENT)")
	flag.BoolVar(featureThemeManagement, "ThemeManagement", *featureThemeManagement, "enable license feature \"ThemeManagement\" (env: MM_LICENSE_FEATURE_THEMEMANAGEMENT)")
	flag.BoolVar(featureEmailNotificationContents, "EmailNotificationContents", *featureEmailNotificationContents, "enable license feature \"EmailNotificationContents\" (env: MM_LICENSE_FEATURE_EMAILNOTIFICATIONCONTENTS)")
	flag.BoolVar(featureDataRetention, "DataRetention", *featureDataRetention, "enable license feature \"DataRetention\" (env: MM_LICENSE_FEATURE_DATARETENTION)")
	flag.BoolVar(featureMessageExport, "MessageExport", *featureMessageExport, "enable license feature \"MessageExport\" (env: MM_LICENSE_FEATURE_MESSAGEEXPORT)")
	flag.BoolVar(featureCustomPermissionsSchemes, "CustomPermissionsSchemes", *featureCustomPermissionsSchemes, "enable license feature \"CustomPermissionsSchemes\" (env: MM_LICENSE_FEATURE_CUSTOMPERMISSIONSSCHEMES)")
	flag.BoolVar(featureCustomTermsOfService, "CustomTermsOfService", *featureCustomTermsOfService, "enable license feature \"CustomTermsOfService\" (env: MM_LICENSE_FEATURE_CUSTOMTERMSOFSERVICE)")
	flag.BoolVar(featureGuestAccounts, "GuestAccounts", *featureGuestAccounts, "enable license feature \"GuestAccounts\" (env: MM_LICENSE_FEATURE_GUESTACCOUNTS)")
	flag.BoolVar(featureGuestAccountsPermissions, "GuestAccountsPermissions", *featureGuestAccountsPermissions, "enable license feature \"GuestAccountsPermissions\" (env: MM_LICENSE_FEATURE_GUESTACCOUNTSPERMISSIONS)")
	flag.BoolVar(featureIDLoadedPushNotifications, "IDLoadedPushNotifications", *featureIDLoadedPushNotifications, "enable license feature \"IDLoadedPushNotifications\" (env: MM_LICENSE_FEATURE_IDLOADEDPUSHNOTIFICATIONS)")
	flag.BoolVar(featureLockTeammateNameDisplay, "LockTeammateNameDisplay", *featureLockTeammateNameDisplay, "enable license feature \"LockTeammateNameDisplay\" (env: MM_LICENSE_FEATURE_LOCKTEAMMATENAMEDISPLAY)")
	flag.BoolVar(featureEnterprisePlugins, "EnterprisePlugins", *featureEnterprisePlugins, "enable license feature \"EnterprisePlugins\" (env: MM_LICENSE_FEATURE_ENTERPRISEPLUGINS)")
	flag.BoolVar(featureAdvancedLogging, "AdvancedLogging", *featureAdvancedLogging, "enable license feature \"AdvancedLogging\" (env: MM_LICENSE_FEATURE_ADVANCEDLOGGING)")
	flag.BoolVar(featureCloud, "Cloud", *featureCloud, "enable license feature \"Cloud\" (env: MM_LICENSE_FEATURE_CLOUD)")
	flag.BoolVar(featureSharedChannels, "SharedChannels", *featureSharedChannels, "enable license feature \"SharedChannels\" (env: MM_LICENSE_FEATURE_SHAREDCHANNELS)")
	flag.BoolVar(featureRemoteClusterService, "RemoteClusterService", *featureRemoteClusterService, "enable license feature \"RemoteClusterService\" (env: MM_LICENSE_FEATURE_REMOTECLUSTERSERVICE)")
	flag.BoolVar(featureOutgoingOAuthConnections, "OutgoingOAuthConnections", *featureOutgoingOAuthConnections, "enable license feature \"OutgoingOAuthConnections\" (env: MM_LICENSE_FEATURE_OUTGOINGOAUTHCONNECTIONS)")
	flag.BoolVar(featureFutureFeatures, "FutureFeatures", *featureFutureFeatures, "enable license feature \"FutureFeatures\" (env: MM_LICENSE_FEATURE_FUTUREFEATURES)")

	flag.BoolVar(&PrintJSON, "print", false, "print license as JSON to stdout")
	flag.Parse()

	if *configPath != "" {
		if buffer, err := os.ReadFile(*configPath); os.IsNotExist(err) {
			log.Fatalln("file not found: ", filepath.Base(*configPath))
		} else if err != nil {
			log.Fatalln("file read error:", err)
		} else {
			if err = json.Unmarshal(buffer, Config); err != nil {
				log.Fatalln(err)
			}
		}
	}

	if *skuName != "" {
		Config.SkuName = *skuName
	}
	if *skuShortName != "" {
		Config.SkuShortName = *skuShortName
	}

	if *isTrial {
		Config.IsTrial = *isTrial
	}

	if *isGovSku {
		Config.IsGovSku = *isGovSku
	}

	if *customerID != "" {
		Config.Customer.Id = *customerID
	}
	if *customerName != "" {
		Config.Customer.Name = *customerName
	}
	if *customerEmail != "" {
		Config.Customer.Email = *customerEmail
	}
	if *customerCompany != "" {
		Config.Customer.Company = *customerCompany
	}

	if *featureUsers > 0 {
		Config.Features.Users = featureUsers
	}
	if *featureLDAP {
		Config.Features.LDAP = featureLDAP
	}
	if *featureLDAPGroups {
		Config.Features.LDAPGroups = featureLDAPGroups
	}
	if *featureMFA {
		Config.Features.MFA = featureMFA
	}
	if *featureGoogleOAuth {
		Config.Features.GoogleOAuth = featureGoogleOAuth
	}
	if *featureOffice365OAuth {
		Config.Features.Office365OAuth = featureOffice365OAuth
	}
	if *featureOpenId {
		Config.Features.OpenId = featureOpenId
	}
	if *featureCompliance {
		Config.Features.Compliance = featureCompliance
	}
	if *featureCluster {
		Config.Features.Cluster = featureCluster
	}
	if *featureMetrics {
		Config.Features.Metrics = featureMetrics
	}
	if *featureMHPNS {
		Config.Features.MHPNS = featureMHPNS
	}
	if *featureSAML {
		Config.Features.SAML = featureSAML
	}
	if *featureElasticsearch {
		Config.Features.Elasticsearch = featureElasticsearch
	}
	if *featureAnnouncement {
		Config.Features.Announcement = featureAnnouncement
	}
	if *featureThemeManagement {
		Config.Features.ThemeManagement = featureThemeManagement
	}
	if *featureEmailNotificationContents {
		Config.Features.EmailNotificationContents = featureEmailNotificationContents
	}
	if *featureDataRetention {
		Config.Features.DataRetention = featureDataRetention
	}
	if *featureMessageExport {
		Config.Features.MessageExport = featureMessageExport
	}
	if *featureCustomPermissionsSchemes {
		Config.Features.CustomPermissionsSchemes = featureCustomPermissionsSchemes
	}
	if *featureCustomTermsOfService {
		Config.Features.CustomTermsOfService = featureCustomTermsOfService
	}
	if *featureGuestAccounts {
		Config.Features.GuestAccounts = featureGuestAccounts
	}
	if *featureGuestAccountsPermissions {
		Config.Features.GuestAccountsPermissions = featureGuestAccountsPermissions
	}
	if *featureIDLoadedPushNotifications {
		Config.Features.IDLoadedPushNotifications = featureIDLoadedPushNotifications
	}
	if *featureLockTeammateNameDisplay {
		Config.Features.LockTeammateNameDisplay = featureLockTeammateNameDisplay
	}
	if *featureEnterprisePlugins {
		Config.Features.EnterprisePlugins = featureEnterprisePlugins
	}
	if *featureAdvancedLogging {
		Config.Features.AdvancedLogging = featureAdvancedLogging
	}
	if *featureCloud {
		Config.Features.Cloud = featureCloud
	}
	if *featureSharedChannels {
		Config.Features.SharedChannels = featureSharedChannels
	}
	if *featureRemoteClusterService {
		Config.Features.RemoteClusterService = featureRemoteClusterService
	}
	if *featureOutgoingOAuthConnections {
		Config.Features.OutgoingOAuthConnections = featureOutgoingOAuthConnections
	}
	if *featureFutureFeatures {
		Config.Features.FutureFeatures = featureFutureFeatures
	}

	Config.Id = gen.NewLicenseID()
	Config.IssuedAt = time.Now().UnixMilli()

	Config.StartsAt = func() int64 {
		now := time.Now()
		if startsAt != "" {
			if ts, err := time.Parse(time.DateTime, startsAt); err != nil {
				log.Fatalln(err)
			} else {
				return ts.UnixMilli()
			}
		} else if startsAtDate != "" || startsAtTime != "" {
			var (
				tsDate int64 = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).UnixMilli()
				tsTime int64 = 0
			)
			if startsAtDate != "" {
				if ts, err := time.Parse(time.DateTime, startsAtDate); err != nil {
					log.Fatalln(err)
				} else {
					tsDate = ts.UnixMilli()
				}
			}
			if startsAtTime != "" {
				if ts, err := time.Parse(time.DateTime, startsAtTime); err != nil {
					log.Fatalln(err)
				} else {
					tsTime = ts.UnixMilli()
				}
			}
			return tsDate + tsTime
		}
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).UnixMilli()
	}()

	Config.ExpiresAt = func() int64 {
		startsAt := time.UnixMilli(Config.StartsAt)
		startsAt = time.Date(startsAt.Year(), startsAt.Month(), startsAt.Day(), 0, 0, 0, 0, time.Local)

		if expiresAt != "" {
			if ts, err := time.Parse(time.DateTime, expiresAt); err != nil {
				log.Fatalln(err)
			} else {
				return ts.UnixMilli()
			}
		} else if expiresAtDate != "" || expiresAtTime != "" {
			var (
				tsDate int64 = startsAt.UnixMilli()
				tsTime int64 = 0
			)
			if expiresAtDate != "" {
				if ts, err := time.Parse(time.DateTime, expiresAtDate); err != nil {
					log.Fatalln(err)
				} else {
					tsDate = ts.UnixMilli()
				}
			}
			if expiresAtTime != "" {
				if ts, err := time.Parse(time.DateTime, expiresAtTime); err != nil {
					log.Fatalln(err)
				} else {
					tsTime = ts.UnixMilli()
				}
			}
			return tsDate + tsTime
		}

		expiresAt := startsAt.Add(time.Hour * 24 * time.Duration(expiresIn))
		return time.Date(expiresAt.Year(), expiresAt.Month(), expiresAt.Day(), 0, 0, 0, 0, time.Local).UnixMilli()
	}()
}
