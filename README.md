[![Build Status](https://semaphoreci.com/api/v1/learnsecurely/vault-plugin-cfssl/branches/gen-cfssl-cert/badge.svg)](https://semaphoreci.com/learnsecurely/vault-plugin-cfssl)

# HashiCorp Vault Plugin: CFSSL

This project provides a plugin for HashiCorp Vault that enables it to offer
certificates on behalf of a CFSSL issuing authority. This is a primarily
used as a proof-of-concept on the way towards providing a similar plugin
for Microsoft ADCS.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

You must have a running CFSSL service and already have initialized it as a CA.
```
# Start CFSSL server
$ cfssl serve -ca-key ca-key.pem \
              -ca ca.pem \
              -config config_ca.json
```
```
# Test CFSSL cert gen w/ curl
$ curl -X POST \
       -H "Content-Type: application/json" \
       -d @/path/to/example-csr.json \
       127.0.0.1:8888/api/v1/cfssl/newkey
```

### Building

Build using the go build command

```
$ dep ensure
```
```
$ go build -o /path/to/vault-plugins/vault-plugin-cfssl
```

### Tests

Test using the go test command
```
$ go test -v ./...
```

## Deployment

If you do not have a running Vault server, follow this step to start a
Vault server in dev mode:
```
$ echo 'plugin_directory = "/path/to/vault-plugins"' \
  | tee /tmp/vault.hcl
$ vault server -dev \
               -dev-root-token-id="root" \
               -config=/tmp/vault.hcl
$ export VAULT_ADDR=http://127.0.0.1:8200
```  
Install & Mount the Plugin
```
$ SHASUM=$(shasum -a 256 "/path/to/vault-plugins/vault-plugin-cfssl" | cut -d " " -f1)
$ vault write sys/plugins/catalog/cfssl \
              sha_256="$SHASUM" \
              command="vault-plugin-cfssl"
$ vault mount -path=cfssl-example \
              -plugin-name=cfssl \
              plugin
```
Pass CSR & Get New Cert
```
$ vault write cfssl-example/issue
              csr=@csr.json-example \
              url=http://127.0.0.1:8888
```
Write Cert to File and Verify
```
$ vault write cfssl-example/issue \
              csr=@csr.json-example \
              url=http://127.0.0.1:8888 \
  | tail -n +3 | sed 's/^testng//' \
  | jq -r .result.certificate \
  | tee /tmp/cert-example.pem
```
Dump the contents of the cert to verify
```
# Using ceritgo
$ certigo -v dump /tmp/cert-example.pem

# Using openssl
$ openssl x509 -in /tmp/cert-example.pem -text
```

## Built With

* [vault](https://github.com/hashicorp/vault) - Secrets management library
* [dep](https://github.com/golang/dep) - Dependency management

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct, and the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/learnsecurely/vault-plugin-cfssl/tags). 

## Authors

* **Jeremy Pruitt** - *Initial work* - [jeremypruitt](https://github.com/jeremypruitt)
* **Jim Smyth** - *Initial work* - [jsmyth](https://github.com/jsmyth)

See also the list of [contributors](https://github.com/learnsecurely/vault-plugin-cfssl/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Thanks to Seth Vargo for his [blog post](https://www.hashicorp.com/blog/building-a-vault-secure-plugin) and [example plugin](https://github.com/hashicorp/vault-auth-plugin-example) on how to create a secure vault plugin
* The Vault [PKI logcial backend](https://github.com/hashicorp/vault/tree/master/builtin/logical/pki) was also a useful resource while developing this plugin
