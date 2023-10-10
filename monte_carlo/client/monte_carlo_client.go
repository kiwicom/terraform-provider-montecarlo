package client

import (
	"context"
	"net/http"
	"net/http/httputil"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hasura/go-graphql-client"
)

// client interface
type MonteCarloClient interface {
	Mutate(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error
	Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error
	ExecRaw(ctx context.Context, query string, variables map[string]interface{}, options ...graphql.Option) ([]byte, error)
}

type monteCarloTransport struct {
	API_KEY_ID    string
	API_KEY_TOKEN string
	context       context.Context
}

type monteCarloClient struct {
	client    *graphql.Client
	transport *monteCarloTransport
}

func NewMonteCarloClient(context context.Context, api_key_id string, api_key_token string) (MonteCarloClient, error) {
	transport := monteCarloTransport{api_key_id, api_key_token, context}
	client := graphql.NewClient("https://api.getmontecarlo.com/graphql", &http.Client{Transport: transport})
	return &monteCarloClient{client, &transport}, nil
}

func (mc *monteCarloClient) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return mc.client.Mutate(ctx, m, variables, options...)
}

func (mc *monteCarloClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	return mc.client.Query(ctx, q, variables, options...)
}

func (mc *monteCarloClient) ExecRaw(ctx context.Context, query string, variables map[string]interface{}, options ...graphql.Option) ([]byte, error) {
	return mc.client.ExecRaw(ctx, query, variables, options...)
}

func (transport monteCarloTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	bytes, _ := httputil.DumpRequestOut(req, true)
	req.Header.Set("x-mcd-id", transport.API_KEY_ID)
	req.Header.Set("x-mcd-token", transport.API_KEY_TOKEN)
	resp, err := http.DefaultTransport.RoundTrip(req)
	respBytes, _ := httputil.DumpResponse(resp, true)
	bytes = append(bytes, respBytes...)
	tflog.Debug(transport.context, string(bytes))
	return resp, err
}

type UUID string
type JSONString string

type BqTestDiagnostic struct {
	Cause           string
	FriendlyMessage string
	Resolution      string
}

type BqTestWarnings []BqTestDiagnostic
type BqTestErrors []BqTestDiagnostic

type TestBqCredentialsV2 struct {
	TestBqCredentialsV2 struct {
		Key              string
		ValidationResult struct {
			Success  bool
			Warnings BqTestWarnings
			Errors   BqTestErrors
		}
	} `graphql:"testBqCredentialsV2(validationName: $validationName, connectionDetails: $connectionDetails)"`
}

type AddConnection struct {
	AddConnection struct {
		Connection struct {
			Uuid      string
			Warehouse struct {
				Name string
				Uuid string
			}
		}
	} `graphql:"addConnection(dcId: $dcId, dwId: $dwId, key: $key, jobTypes: $jobTypes, name: $name, connectionType: $connectionType, createWarehouseType: $createWarehouseType)"`
}

type GetWarehouse struct {
	GetWarehouse *struct {
		Name        string `json:"name"`
		Connections []struct {
			Uuid string `json:"uuid"`
			Type string `json:"type"`
		} `json:"connections"`
		DataCollector struct {
			Uuid string `json:"uuid"`
		} `json:"dataCollector"`
	} `json:"getWarehouse"`
}

type RemoveConnection struct {
	RemoveConnection struct {
		Success bool
	} `graphql:"removeConnection(connectionId: $connectionId)"`
}

type SetWarehouseName struct {
	SetWarehouseName struct {
		Warehouse struct {
			Uuid string
			Name string
		}
	} `graphql:"setWarehouseName(dwId: $dwId, name: $name)"`
}

type UpdateCredentials struct {
	UpdateCredentials struct {
		Success bool
	} `graphql:"updateCredentials(changes: $changes, connectionId: $connectionId, shouldReplace: $shouldReplace, shouldValidate: $shouldValidate)"`
}

type TagPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type TagKeyValuePairInput TagPair
type TagKeyValuePairOutput TagPair

type CreateOrUpdateDomain struct {
	CreateOrUpdateDomain struct {
		Domain struct {
			Assignments []string
			Tags        []TagKeyValuePairOutput
			Name        string
			Description string
			Uuid        string
		}
	} `graphql:"createOrUpdateDomain(assignments: $assignments, tags: $tags, name: $name, description: $description, uuid: $uuid)"`
}

type GetDomain struct {
	GetDomain *struct {
		Uuid           string                  `json:"uuid"`
		Name           string                  `json:"name"`
		Description    string                  `json:"description"`
		CreatedByEmail string                  `json:"createdByEmail"`
		Tags           []TagKeyValuePairOutput `json:"tags"`
		Assignments    []string                `json:"assignments"`
	} `json:"getDomain"`
}

type DeleteDomain struct {
	DeleteDomain struct {
		Deleted int
	} `graphql:"deleteDomain(uuid: $uuid)"`
}

type DatabaseTestDiagnostic struct {
	Message string
	Type    string
}

type TestDatabaseCredentials struct {
	TestDatabaseCredentials struct {
		Key         string
		Success     bool
		Warnings    []DatabaseTestDiagnostic
		Validations []DatabaseTestDiagnostic
	} `graphql:"testDatabaseCredentials(connectionType: $connectionType, dbName: $dbName, dbType: $dbType, host: $host, port: $port, user: $user, password: $password)"`
}
