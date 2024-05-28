package license

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mmlicense/env"
	"mmlicense/gen"

	"git.biggo.com/Funmula/mattermost-funmula/server/public/model"
)

var (
	configPath *string = env.String("MM_LICENSE_CONFIG")

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

	flag.StringVar(&startsAt, "starts_at", startsAt, fmt.Sprintf("set the license start date and time (example: \"%s\")", time.Now().Format(time.DateTime)))
	flag.StringVar(&startsAtDate, "starts_at_date", startsAtDate, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.StringVar(&startsAtTime, "starts_at_time", startsAtTime, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.TimeOnly)))
	flag.StringVar(&expiresAt, "expires_at", expiresAt, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.StringVar(&expiresAtDate, "expires_at_date", expiresAtDate, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.StringVar(&expiresAtTime, "expires_at_time", expiresAtTime, fmt.Sprintf("set the license start date (example: \"%s\")", time.Now().Format(time.DateOnly)))
	flag.IntVar(&expiresIn, "days", expiresIn, "set expiration days for the license")

	flag.StringVar(skuName, "sku_name", *skuName, "set the license sku name (env: MM_LICENSE_SKU_NAME)")
	flag.StringVar(skuShortName, "sku_short_name", *skuShortName, "set the license sku short name (env: MM_LICENSE_SKU_SHORT_NAME)")

	flag.BoolVar(isTrial, "is_trial", *isTrial, "set is trial license flag (env: MM_LICENSE_IS_TRIAL)")
	flag.BoolVar(isGovSku, "is_gov_sku", *isGovSku, "set is government license flag (env: MM_LICENSE_IS_GOVERNMENT)")

	flag.StringVar(customerID, "customer_id", *customerID, "set license customer id (env: MM_LICENSE_CUSTOMER_ID)")
	flag.StringVar(customerName, "customer_name", *customerName, "set license customer name (env: MM_LICENSE_CUSTOMER_NAME)")
	flag.StringVar(customerEmail, "customer_email", *customerEmail, "set license customer email (env: MM_LICENSE_CUSTOMER_EMAIL)")
	flag.StringVar(customerCompany, "customer_company", *customerCompany, "set license customer comany name (env: MM_LICENSE_CUSTOMER_COMPANY)")

	flag.IntVar(featureUsers, "users", *featureUsers, "enable license feature \"Users\" (env: MM_LICENSE_FEATURE_USERS)")
	flag.BoolVar(featureLDAP, "ldap", *featureLDAP, "enable license feature \"LDAP\" (env: MM_LICENSE_FEATURE_LDAP)")
	flag.BoolVar(featureLDAPGroups, "ldap_groups", *featureLDAPGroups, "enable license feature \"LDAPGroups\" (env: MM_LICENSE_FEATURE_LDAPGROUPS)")
	flag.BoolVar(featureMFA, "mfa", *featureMFA, "enable license feature \"MFA\" (env: MM_LICENSE_FEATURE_MFA)")
	flag.BoolVar(featureGoogleOAuth, "google_oauth", *featureGoogleOAuth, "enable license feature \"GoogleOAuth\" (env: MM_LICENSE_FEATURE_GOOGLEOAUTH)")
	flag.BoolVar(featureOffice365OAuth, "office365_oauth", *featureOffice365OAuth, "enable license feature \"Office365OAuth\" (env: MM_LICENSE_FEATURE_OFFICE365OAUTH)")
	flag.BoolVar(featureOpenId, "openid", *featureOpenId, "enable license feature \"OpenId\" (env: MM_LICENSE_FEATURE_OPENID)")
	flag.BoolVar(featureCompliance, "compliance", *featureCompliance, "enable license feature \"Compliance\" (env: MM_LICENSE_FEATURE_COMPLIANCE)")
	flag.BoolVar(featureCluster, "cluster", *featureCluster, "enable license feature \"Cluster\" (env: MM_LICENSE_FEATURE_CLUSTER)")
	flag.BoolVar(featureMetrics, "metrics", *featureMetrics, "enable license feature \"Metrics\" (env: MM_LICENSE_FEATURE_METRICS)")
	flag.BoolVar(featureMHPNS, "mhpns", *featureMHPNS, "enable license feature \"MHPNS\" (env: MM_LICENSE_FEATURE_MHPNS)")
	flag.BoolVar(featureSAML, "saml", *featureSAML, "enable license feature \"SAML\" (env: MM_LICENSE_FEATURE_SAML)")
	flag.BoolVar(featureElasticsearch, "elastic_search", *featureElasticsearch, "enable license feature \"Elasticsearch\" (env: MM_LICENSE_FEATURE_ELASTICSEARCH)")
	flag.BoolVar(featureAnnouncement, "announcement", *featureAnnouncement, "enable license feature \"Announcement\" (env: MM_LICENSE_FEATURE_ANNOUNCEMENT)")
	flag.BoolVar(featureThemeManagement, "theme_management", *featureThemeManagement, "enable license feature \"ThemeManagement\" (env: MM_LICENSE_FEATURE_THEMEMANAGEMENT)")
	flag.BoolVar(featureEmailNotificationContents, "email_notification_contents", *featureEmailNotificationContents, "enable license feature \"EmailNotificationContents\" (env: MM_LICENSE_FEATURE_EMAILNOTIFICATIONCONTENTS)")
	flag.BoolVar(featureDataRetention, "data_retention", *featureDataRetention, "enable license feature \"DataRetention\" (env: MM_LICENSE_FEATURE_DATARETENTION)")
	flag.BoolVar(featureMessageExport, "message_export", *featureMessageExport, "enable license feature \"MessageExport\" (env: MM_LICENSE_FEATURE_MESSAGEEXPORT)")
	flag.BoolVar(featureCustomPermissionsSchemes, "custom_permissions_schemes", *featureCustomPermissionsSchemes, "enable license feature \"CustomPermissionsSchemes\" (env: MM_LICENSE_FEATURE_CUSTOMPERMISSIONSSCHEMES)")
	flag.BoolVar(featureCustomTermsOfService, "custom_terms_of_service", *featureCustomTermsOfService, "enable license feature \"CustomTermsOfService\" (env: MM_LICENSE_FEATURE_CUSTOMTERMSOFSERVICE)")
	flag.BoolVar(featureGuestAccounts, "guest_accounts", *featureGuestAccounts, "enable license feature \"GuestAccounts\" (env: MM_LICENSE_FEATURE_GUESTACCOUNTS)")
	flag.BoolVar(featureGuestAccountsPermissions, "guest_accounts_permissions", *featureGuestAccountsPermissions, "enable license feature \"GuestAccountsPermissions\" (env: MM_LICENSE_FEATURE_GUESTACCOUNTSPERMISSIONS)")
	flag.BoolVar(featureIDLoadedPushNotifications, "id_loaded", *featureIDLoadedPushNotifications, "enable license feature \"IDLoadedPushNotifications\" (env: MM_LICENSE_FEATURE_IDLOADEDPUSHNOTIFICATIONS)")
	flag.BoolVar(featureLockTeammateNameDisplay, "lock_teammate_name_display", *featureLockTeammateNameDisplay, "enable license feature \"LockTeammateNameDisplay\" (env: MM_LICENSE_FEATURE_LOCKTEAMMATENAMEDISPLAY)")
	flag.BoolVar(featureEnterprisePlugins, "enterprise_plugins", *featureEnterprisePlugins, "enable license feature \"EnterprisePlugins\" (env: MM_LICENSE_FEATURE_ENTERPRISEPLUGINS)")
	flag.BoolVar(featureAdvancedLogging, "advanced_logging", *featureAdvancedLogging, "enable license feature \"AdvancedLogging\" (env: MM_LICENSE_FEATURE_ADVANCEDLOGGING)")
	flag.BoolVar(featureCloud, "cloud", *featureCloud, "enable license feature \"Cloud\" (env: MM_LICENSE_FEATURE_CLOUD)")
	flag.BoolVar(featureSharedChannels, "shared_channels", *featureSharedChannels, "enable license feature \"SharedChannels\" (env: MM_LICENSE_FEATURE_SHAREDCHANNELS)")
	flag.BoolVar(featureRemoteClusterService, "remote_cluster_service", *featureRemoteClusterService, "enable license feature \"RemoteClusterService\" (env: MM_LICENSE_FEATURE_REMOTECLUSTERSERVICE)")
	flag.BoolVar(featureOutgoingOAuthConnections, "outgoing_oauth_connections", *featureOutgoingOAuthConnections, "enable license feature \"OutgoingOAuthConnections\" (env: MM_LICENSE_FEATURE_OUTGOINGOAUTHCONNECTIONS)")
	flag.BoolVar(featureFutureFeatures, "future_features", *featureFutureFeatures, "enable license feature \"FutureFeatures\" (env: MM_LICENSE_FEATURE_FUTUREFEATURES)")
}

func New() (config *model.License) {
	config = &model.License{Customer: &model.Customer{}, Features: &model.Features{}}

	if *configPath != "" {
		if buffer, err := os.ReadFile(*configPath); os.IsNotExist(err) {
			log.Fatalln("file not found: ", filepath.Base(*configPath))
		} else if err != nil {
			log.Fatalln("file read error:", err)
		} else {
			if err = json.Unmarshal(buffer, config); err != nil {
				log.Fatalln(err)
			}
		}
	}

	config.Id = gen.NewLicenseID()
	config.IssuedAt = time.Now().UnixMilli()

	config.StartsAt = func() int64 {
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

	config.ExpiresAt = func() int64 {
		startsAt := time.UnixMilli(config.StartsAt)
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

	if isUserValue("sku_name", "MM_LICENSE_SKU_NAME") {
		config.SkuName = *skuName
	}
	if isUserValue("sku_short_name", "MM_LICENSE_SKU_SHORT_NAME") {
		config.SkuShortName = *skuShortName
	}

	if isUserValue("is_trial", "MM_LICENSE_IS_TRIAL") {
		config.IsTrial = *isTrial
	}
	if isUserValue("is_gov_sku", "MM_LICENSE_IS_GOVERNMENT") {
		config.IsGovSku = *isGovSku
	}

	if isUserValue("customer_id", "MM_LICENSE_CUSTOMER_ID") {
		config.Customer.Id = *customerID
	}
	if isUserValue("customer_name", "MM_LICENSE_CUSTOMER_NAME") {
		config.Customer.Name = *customerName
	}
	if isUserValue("customer_email", "MM_LICENSE_CUSTOMER_EMAIL") {
		config.Customer.Email = *customerEmail
	}
	if isUserValue("customer_company", "MM_LICENSE_CUSTOMER_COMPANY") {
		config.Customer.Company = *customerCompany
	}

	if isUserValue("users", "MM_LICENSE_FEATURE_USERS") {
		config.Features.Users = featureUsers
	}
	if isUserValue("ldap", "MM_LICENSE_FEATURE_LDAP") {
		config.Features.LDAP = featureLDAP
	}
	if isUserValue("ldap_groups", "MM_LICENSE_FEATURE_LDAPGROUPS") {
		config.Features.LDAPGroups = featureLDAPGroups
	}
	if isUserValue("mfa", "MM_LICENSE_FEATURE_MFA") {
		config.Features.MFA = featureMFA
	}
	if isUserValue("google_oauth", "MM_LICENSE_FEATURE_GOOGLEOAUTH") {
		config.Features.GoogleOAuth = featureGoogleOAuth
	}
	if isUserValue("office365_oauth", "MM_LICENSE_FEATURE_OFFICE365OAUTH") {
		config.Features.Office365OAuth = featureOffice365OAuth
	}
	if isUserValue("openid", "MM_LICENSE_FEATURE_OPENID") {
		config.Features.OpenId = featureOpenId
	}
	if isUserValue("compliance", "MM_LICENSE_FEATURE_COMPLIANCE") {
		config.Features.Compliance = featureCompliance
	}
	if isUserValue("cluster", "MM_LICENSE_FEATURE_CLUSTER") {
		config.Features.Cluster = featureCluster
	}
	if isUserValue("metrics", "MM_LICENSE_FEATURE_METRICS") {
		config.Features.Metrics = featureMetrics
	}
	if isUserValue("mhpns", "MM_LICENSE_FEATURE_MHPNS") {
		config.Features.MHPNS = featureMHPNS
	}
	if isUserValue("saml", "MM_LICENSE_FEATURE_SAML") {
		config.Features.SAML = featureSAML
	}
	if isUserValue("elastic_search", "MM_LICENSE_FEATURE_ELASTICSEARCH") {
		config.Features.Elasticsearch = featureElasticsearch
	}
	if isUserValue("announcement", "MM_LICENSE_FEATURE_ANNOUNCEMENT") {
		config.Features.Announcement = featureAnnouncement
	}
	if isUserValue("theme_management", "MM_LICENSE_FEATURE_THEMEMANAGEMENT") {
		config.Features.ThemeManagement = featureThemeManagement
	}
	if isUserValue("email_notification_contents", "MM_LICENSE_FEATURE_EMAILNOTIFICATIONCONTENTS") {
		config.Features.EmailNotificationContents = featureEmailNotificationContents
	}
	if isUserValue("data_retention", "MM_LICENSE_FEATURE_DATARETENTION") {
		config.Features.DataRetention = featureDataRetention
	}
	if isUserValue("message_export", "MM_LICENSE_FEATURE_MESSAGEEXPORT") {
		config.Features.MessageExport = featureMessageExport
	}
	if isUserValue("custom_permissions_schemes", "MM_LICENSE_FEATURE_CUSTOMPERMISSIONSSCHEMES") {
		config.Features.CustomPermissionsSchemes = featureCustomPermissionsSchemes
	}
	if isUserValue("custom_terms_of_service", "MM_LICENSE_FEATURE_CUSTOMTERMSOFSERVICE") {
		config.Features.CustomTermsOfService = featureCustomTermsOfService
	}
	if isUserValue("guest_accounts", "MM_LICENSE_FEATURE_GUESTACCOUNTS") {
		config.Features.GuestAccounts = featureGuestAccounts
	}
	if isUserValue("guest_accounts_permissions", "MM_LICENSE_FEATURE_GUESTACCOUNTSPERMISSIONS") {
		config.Features.GuestAccountsPermissions = featureGuestAccountsPermissions
	}
	if isUserValue("id_loaded", "MM_LICENSE_FEATURE_IDLOADEDPUSHNOTIFICATIONS") {
		config.Features.IDLoadedPushNotifications = featureIDLoadedPushNotifications
	}
	if isUserValue("lock_teammate_name_display", "MM_LICENSE_FEATURE_LOCKTEAMMATENAMEDISPLAY") {
		config.Features.LockTeammateNameDisplay = featureLockTeammateNameDisplay
	}
	if isUserValue("enterprise_plugins", "MM_LICENSE_FEATURE_ENTERPRISEPLUGINS") {
		config.Features.EnterprisePlugins = featureEnterprisePlugins
	}
	if isUserValue("advanced_logging", "MM_LICENSE_FEATURE_ADVANCEDLOGGING") {
		config.Features.AdvancedLogging = featureAdvancedLogging
	}
	if isUserValue("cloud", "MM_LICENSE_FEATURE_CLOUD") {
		config.Features.Cloud = featureCloud
	}
	if isUserValue("shared_channels", "MM_LICENSE_FEATURE_SHAREDCHANNELS") {
		config.Features.SharedChannels = featureSharedChannels
	}
	if isUserValue("remote_cluster_service", "MM_LICENSE_FEATURE_REMOTECLUSTERSERVICE") {
		config.Features.RemoteClusterService = featureRemoteClusterService
	}
	if isUserValue("outgoing_oauth_connections", "MM_LICENSE_FEATURE_OUTGOINGOAUTHCONNECTIONS") {
		config.Features.OutgoingOAuthConnections = featureOutgoingOAuthConnections
	}
	if isUserValue("future_features", "MM_LICENSE_FEATURE_FUTUREFEATURES") {
		config.Features.FutureFeatures = featureFutureFeatures
	}

	return
}

var (
	argMap map[string]bool
	envMap map[string]bool
)

func isUserValue(argKey, envKey string) bool {
	if argMap == nil {
		argMap = map[string]bool{}
		for _, arg := range os.Args {
			if strings.HasPrefix(arg, "-") {
				argMap[strings.SplitN(arg, "=", 2)[0][1:]] = true
			}
		}
	}

	if envMap == nil {
		envMap = map[string]bool{}
		for _, env := range os.Environ() {
			envMap[strings.SplitN(env, "=", 2)[0]] = true
		}
	}

	if _, ok := argMap[argKey]; ok {
		return true
	}

	if _, ok := envMap[envKey]; ok {
		return true
	}
	return false
}
