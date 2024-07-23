# server license load function

func: LoadLicense()

```bash
server/channels/app/platform/license.go
```

## server license public key info

```text
Algo RSA
Format X.509
```

## funmula server license private key

generated via:

```bash
openssl genrsa -out private.pem 2048
```

## funmula server license public key

generated via:

```bash
openssl rsa -in private.pem -outform PEM -pubout -out public.pem
```
