package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hashicorp/vault/helper/pluginutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/logical/plugin"
)

func main() {
	apiClientMeta := &pluginutil.APIClientMeta{}
	flags := apiClientMeta.FlagSet()
	flags.Parse(os.Args[1:])

	tlsConfig := apiClientMeta.GetTLSConfig()
	tlsProviderFunc := pluginutil.VaultPluginTLSProvider(tlsConfig)

	if err := plugin.Serve(&plugin.ServeOpts{
		BackendFactoryFunc: Factory,
		TLSProviderFunc:    tlsProviderFunc,
	}); err != nil {
		log.Fatal(err)
	}
}

func Factory(conf *logical.BackendConfig) (logical.Backend, error) {
	b := Backend()
	if err := b.Setup(conf); err != nil {
		return nil, err
	}
	return b, nil
}

type backend struct {
	*framework.Backend
}

func Backend() *backend {
	var b backend
	b.Backend = &framework.Backend{
		BackendType: logical.TypeLogical,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{"issue"},
		},
		Paths: []*framework.Path{
			&framework.Path{
				Pattern: "issue",
				Fields: map[string]*framework.FieldSchema{
					"csr": &framework.FieldSchema{
						Type: framework.TypeString,
					},
					"url": &framework.FieldSchema{
						Type: framework.TypeString,
					},
				},
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.UpdateOperation: b.pathIssue,
				},
			},
		},
	}

	return &b
}

func (b *backend) pathIssue(req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	csr := d.Get("csr").(string)
	url := d.Get("url").(string)
	jsonData := json.RawMessage(csr)
	jsonValue, err := jsonData.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("Marshalling CSR JSON failed with error %s\n", err)
	}

	cfssl_url := fmt.Sprintf("%s/%s", url, "api/v1/cfssl/newcert")

	response, err := http.Post(cfssl_url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return nil, fmt.Errorf("The HTTP request failed with error %s\n", err)
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("The IOUtil ReadAll request failed with error %s\n", err)

	}
	return &logical.Response{
		Data: map[string]interface{}{"testng": string(data)},
	}, nil
}
