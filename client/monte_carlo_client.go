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
	} `graphql:"testBqCredentialsV2(validationName: $validationName, connectionDetails: $connectionDetails, connectionOptions: $connectionOptions)"`
}

type AddConnection struct {
	AddConnection struct {
		Connection struct {
			Uuid      string
			CreatedOn string
			UpdatedOn string
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
			Uuid      string `json:"uuid"`
			Type      string `json:"type"`
			CreatedOn string `json:"createdOn"`
			UpdatedOn string `json:"updatedOn"`
		} `json:"connections"`
		DataCollector struct {
			Uuid string `json:"uuid"`
		} `json:"dataCollector"`
	} `json:"getWarehouse"`
}

const BqConnectionType = "bigquery"
const BqConnectionTypeResponse = "BIGQUERY"
const TrxConnectionType = "transactional-db"
const TrxConnectionTypeResponse = "TRANSACTIONAL_DB"
const GetWarehouseQuery string = "query getWarehouse($uuid: UUID) { getWarehouse(uuid: $uuid) { name,connections{uuid,type,createdOn,updatedOn},dataCollector{uuid} } }"

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
		Success   bool
		UpdatedAt string
	} `graphql:"updateCredentials(changes: $changes, connectionId: $connectionId, shouldReplace: $shouldReplace, shouldValidate: $shouldValidate)"`
}

type UpdateCredentialsV2 struct {
	UpdateCredentialsV2 struct {
		Success   bool
		UpdatedAt string
	} `graphql:"updateCredentialsV2(connectionId: $connectionId, tempCredentialsKey: $tempCredentialsKey)"`
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

const GetDomainQuery string = "query getDomain($uuid: UUID!) { getDomain(uuid: $uuid) { uuid,name,description,tags{name,value},assignments,createdByEmail } }"

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

type GetTablesEdge struct {
	Node struct {
		Mcon        string
		ProjectName string
		Dataset     string
		TableId     string
		Warehouse   struct {
			Uuid    string
			Account struct {
				Uuid string
			}
		}
	}
}

type GetTables struct {
	GetTables struct {
		Edges    []GetTablesEdge
		PageInfo struct {
			StartCursor string
			EndCursor   string
			HasNextPage bool
		}
	} `graphql:"getTables(dwId: $dwId, first: $first, after: $after, isDeleted: $isDeleted, isExcluded: $isExcluded)"`
}

type AuthorizationGroupUser struct {
	CognitoUserId string
	Email         string
	FirstName     string
	LastName      string
	IsSso         bool
}

type AuthorizationGroup struct {
	Name               string
	Label              string
	Description        string
	IsManaged          bool
	Roles              []struct{ Name string }
	DomainRestrictions []struct{ Uuid string }
	SsoGroup           *string
	Users              []AuthorizationGroupUser
}

type CreateOrUpdateAuthorizationGroup struct {
	CreateOrUpdateAuthorizationGroup struct {
		AuthorizationGroup AuthorizationGroup
	} `graphql:"createOrUpdateAuthorizationGroup(name: $name, label: $label, description: $description, roles: $roles, domainRestrictionIds: $domainRestrictionIds, ssoGroup: $ssoGroup)"`
}

type GetAuthorizationGroups struct {
	GetAuthorizationGroups []AuthorizationGroup `graphql:"getAuthorizationGroups"`
}

type DeleteAuthorizationGroup struct {
	DeleteAuthorizationGroup struct {
		Deleted int
	} `graphql:"deleteAuthorizationGroup(name: $name)"`
}

type User struct {
	CognitoUserId string
	Email         string
	FirstName     string
	LastName      string
	IsSso         bool
	Auth          struct {
		Groups []string
	}
}

type GetUsersInAccount struct {
	GetUsersInAccount struct {
		Edges []struct {
			Node User
		}
		PageInfo struct {
			StartCursor string
			EndCursor   string
			HasNextPage bool
		}
	} `graphql:"getUsersInAccount(email: $email, first: $first, after: $after)"`
}

type UpdateUserAuthorizationGroupMembership struct {
	UpdateUserAuthorizationGroupMembership struct {
		AddedToGroups []struct {
			Name        string
			Label       string
			Description string
		}
		RemovedFromGroups []struct {
			Name        string
			Label       string
			Description string
		}
	} `graphql:"updateUserAuthorizationGroupMembership(memberUserId: $memberUserId, groupNames: $groupNames)"`
}

type CreateOrUpdateComparisonRule struct {
	CreateOrUpdateComparisonRule struct {
		CustomRule struct {
			Uuid              string
			AccountUuid       string
			Projects          []string
			Datasets          []string
			Description       string
			Notes             string
			Labels            []string
			IsTemplateManaged bool
			Namespace         string
			Severity          string
			RuleType          string
			WarehouseUuid     string
			Comparisons       []struct {
				ComparisonType string
				FullTableId    string
				FullTableIds   []string
				Field          string
				Metric         string
				Operator       string
				Threshold      float64
			}
		}
	} `graphql:"createOrUpdateComparisonRule(comparisons: $comparisons, customRuleUuid: $customRuleUuid, description: $description, queryResultType: $queryResultType, scheduleConfig: $scheduleConfig, sourceConnectionId: $sourceConnectionId, sourceDwId: $sourceDwId, sourceSqlQuery: $sourceSqlQuery, targetConnectionId: $targetConnectionId, targetDwId: $targetDwId, targetSqlQuery: $targetSqlQuery)"`
}

type CreateOrUpdateServiceApiToken struct {
	CreateOrUpdateServiceApiToken struct {
		AccessToken struct {
			Id    string
			Token string
		}
	} `graphql:"createOrUpdateServiceApiToken(comment: $comment, displayName: $displayName, expirationInDays: $expirationInDays, groups: $groups, tokenId: $tokenId)"`
}

type TokenMetadata struct {
	Id                string
	Comment           string
	CreatedBy         string
	CreationTime      string
	Email             string
	ExpirationTime    string
	FirstName         string
	LastName          string
	Groups            []string
	IsServiceApiToken bool
}

type GetTokenMetadata struct {
	GetTokenMetadata []TokenMetadata `graphql:"getTokenMetadata(index: $index, isServiceApiToken: $isServiceApiToken)"`
}

type DeleteAccessToken struct {
	DeleteAccessToken struct {
		Success bool
	} `graphql:"deleteAccessToken(tokenId: $tokenId)"`
}
