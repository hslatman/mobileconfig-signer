# MobileConfig Signer

An example `.mobileconfig` signing utility

## Usage

With password prompt:

```console
go run main.go -file example.mobileconfig -cert signing.crt -key private.key
Please enter the password to decrypt "private.key":
2022/10/13 20:43:26 Written signed mobileconfig to "example.signed.mobileconfig"
```

Without password prompt:

```console
go run main.go -file example.mobileconfig -cert signing.crt -key private.key -password password
2022/10/13 20:43:26 Written signed mobileconfig to "example.signed.mobileconfig"
```