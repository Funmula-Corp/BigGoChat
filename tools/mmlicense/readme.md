# funmula license generator for biggo-chat (mattermost)

## use config file

```bash
./mmlicense -config ./license.json
```

sample **license.json**

```json
{
  "customer": {
    "id": "biggo_dev_team",
    "name": "B!gGo Development",
    "email": "root@localhost",
    "company": "Funmula"
  },
  "features": {
    "advanced_logging": false,
    "announcement": false,
    "cloud": false,
    "cluster": true,
    "compliance": false,
    "custom_permissions_schemes": false,
    "custom_terms_of_service": false,
    "data_retention": false,
    "elastic_search": true,
    "email_notification_contents": false,
    "enterprise_plugins": false,
    "future_features": false,
    "google_oauth": false,
    "guest_accounts_permissions": false,
    "guest_accounts": false,
    "id_loaded": false,
    "ldap_groups": false,
    "ldap": false,
    "lock_teammate_name_display": false,
    "message_export": false,
    "metrics": false,
    "mfa": false,
    "mhpns": false,
    "office365_oauth": false,
    "openid": false,
    "outgoing_oauth_connections": false,
    "remote_cluster_service": true,
    "saml": false,
    "shared_channels": false,
    "theme_management": false,
    "users": null
  },
  "sku_name": "B!gGo Chat License",
  "sku_short_name": "B!gGo Chat License",
  "is_trial": false,
  "is_gov_sku": true,
  "signup_jwt": null
}
```

## use environment variables

```bash
# set config path
export MM_LICENSE_CONFIG="./license.json"

# set general information
export MM_LICENSE_SKU_NAME="B!gGo Chat License"
export MM_LICENSE_SKU_SHORT_NAME="B!gGo Chat License"
export MM_LICENSE_IS_TRIAL=false
export MM_LICENSE_IS_GOVERNMENT=false

# set customer information
export MM_LICENSE_CUSTOMER_ID="biggo_dev_team"
export MM_LICENSE_CUSTOMER_NAME="B!gGo Development"
export MM_LICENSE_CUSTOMER_EMAIL="root@localhost"
export MM_LICENSE_CUSTOMER_COMPANY="Funmula"

# set features
export MM_LICENSE_FEATURE_USERS=12
export MM_LICENSE_FEATURE_LDAP=false
export MM_LICENSE_FEATURE_LDAPGROUPS=false
export MM_LICENSE_FEATURE_MFA=false
export MM_LICENSE_FEATURE_GOOGLEOAUTH=false
export MM_LICENSE_FEATURE_OFFICE365OAUTH=false
export MM_LICENSE_FEATURE_OPENID=false
export MM_LICENSE_FEATURE_COMPLIANCE=false
export MM_LICENSE_FEATURE_CLUSTER=false
export MM_LICENSE_FEATURE_METRICS=false
export MM_LICENSE_FEATURE_MHPNS=false
export MM_LICENSE_FEATURE_SAML=false
export MM_LICENSE_FEATURE_ELASTICSEARCH=false
export MM_LICENSE_FEATURE_ANNOUNCEMENT=false
export MM_LICENSE_FEATURE_THEMEMANAGEMENT=false
export MM_LICENSE_FEATURE_EMAILNOTIFICATIONCONTENTS=false
export MM_LICENSE_FEATURE_DATARETENTION=false
export MM_LICENSE_FEATURE_MESSAGEEXPORT=false
export MM_LICENSE_FEATURE_CUSTOMPERMISSIONSSCHEMES=false
export MM_LICENSE_FEATURE_CUSTOMTERMSOFSERVICE=false
export MM_LICENSE_FEATURE_GUESTACCOUNTS=false
export MM_LICENSE_FEATURE_GUESTACCOUNTSPERMISSIONS=false
export MM_LICENSE_FEATURE_IDLOADEDPUSHNOTIFICATIONS=false
export MM_LICENSE_FEATURE_LOCKTEAMMATENAMEDISPLAY=false
export MM_LICENSE_FEATURE_ENTERPRISEPLUGINS=false
export MM_LICENSE_FEATURE_ADVANCEDLOGGING=false
export MM_LICENSE_FEATURE_CLOUD=false
export MM_LICENSE_FEATURE_SHAREDCHANNELS=false
export MM_LICENSE_FEATURE_REMOTECLUSTERSERVICE=false
export MM_LICENSE_FEATURE_OUTGOINGOAUTHCONNECTIONS=false
export MM_LICENSE_FEATURE_FUTUREFEATURES=false

# run the license generator
./mmlicense
```

## use cli arguments

```bash
./mmlicense \
-config="./license.json" \
-sku_name="B!gGo Chat License" \
-sku_short_name="B!gGo Chat License" \
-is_trial \
-is_gov_sku \
-customer_id="biggo_dev_team" \
-customer_name="B!gGo Development" \
-customer_email="root@localhost" \
-customer_company="Funmula" \
-users=12 \
-ldap \
-ldap_groups \
-mfa \
-google_oauth \
-office365_oauth \
-openid \
-compliance \
-cluster \
-metrics \
-mhpns \
-saml \
-elastic_search \
-announcement \
-theme_management \
-email_notification_contents \
-data_retention \
-message_export \
-custom_permissions_schemes \
-custom_terms_of_service \
-guest_accounts \
-guest_accounts_permissions \
-id_loaded \
-lock_teammate_name_display \
-enterprise_plugins \
-advanced_logging \
-cloud \
-shared_channels \
-remote_cluster_service \
-outgoing_oauth_connections \
-future_features
```

### Optional CLI args

- **-q** quiet mode - to print only the encoded license

- **-h** print help (also invoked if unknown args are passed in)

- **-starts_at** set the license start date-time (example: **-starts_at="2000-01-01 00:00:00"**)
- **-starts_at_date** set the license start date (example: **-starts_at_date="2000-01-01"**)
- **-starts_at_time** set the license start time (example: **-starts_at_time="00:00:00"**)
- **-expires_at** set the license expiration date-time (example: **-expires_at="2000-01-01 00:00:00"**)
- **-expires_at_date** set the license expiration date (example: **-expires_at_date="2000-01-01"**)
- **-expires_at_time** set the license expiration time (example: **-expires_at_time="00:00:00"**)

- **-days** set the expiration time in days (takes starts_at settings into account)

- **-insert** insert the newly generated license into the datebase
- **-activate** activate the generated license or the license with the specified license id
- **-license_id** set the license id and skips license generation (use together with activate license)

- **-pg_username** configure the postgres client connection and set the username
- **-pg_password** configure the postgres client connection and set the password
- **-pg_host** configure the postgres client connection and set the host
- **-pg_db** configure the postgres client connection and set the db
- **-pg_port** configure the postgres client connection and set the port

### NOTES

- **starts_at** and **expires_at** settings are only available via CLI!

- default duration of a license is 365 days

- set environment variables will override values from the configuration file
- set cli arguments will override values from environment variables and the configuration file

- all settings in the config file are optional
- all environment variables are optional
- all cli arguments are optional
- cli flag args (boolean) have a default value of true when set - to set them to false use assignment notation (example: **-is_trial=false**)
- string and integer cli arguments support indirect assignment notation (example: **-users 12**)
