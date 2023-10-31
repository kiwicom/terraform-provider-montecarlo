package resources_test

import (
	"fmt"
	"slices"
	"testing"

	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client"
	cmock "github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/client/mock"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/common"
	"github.com/kiwicom/terraform-provider-montecarlo/monte_carlo/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/mock"
)

func TestAccIamMemberResource(t *testing.T) {
	groupName := "group1"
	memberEmail := "user1"
	memberId := "123"
	groupNameUpdated := "group2"

	providerContext := &common.ProviderContext{MonteCarloClient: initIamMemberMonteCarloClient(
		groupName, groupNameUpdated, memberEmail, memberId,
	)}
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"montecarlo": providerserver.NewProtocol6WithError(provider.New("test", providerContext)()),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() {},
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{ // Create and Read testing
				Config: iamMemberConfig(groupName, memberEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "group", fmt.Sprintf("groups/%s", groupName)),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member", fmt.Sprintf("user:%s", memberEmail)),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member_id", memberId),
				),
			},
			{ // ImportState testing
				ResourceName:                         "montecarlo_iam_member.test",
				ImportState:                          true,
				ImportStateVerify:                    true,
				ImportStateId:                        fmt.Sprintf("groups/%[1]s,user:%[2]s", groupName, memberEmail),
				ImportStateVerifyIdentifierAttribute: "group",
			},
			{ // Update and Read testing
				Config: iamMemberConfig(groupNameUpdated, memberEmail),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "group", fmt.Sprintf("groups/%s", groupNameUpdated)),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member", fmt.Sprintf("user:%s", memberEmail)),
					resource.TestCheckResourceAttr("montecarlo_iam_member.test", "member_id", memberId),
				),
			},
		},
	})
}

func iamMemberConfig(groupName, memberEmail string) string {
	return fmt.Sprintf(`
provider "montecarlo" {
	account_service_key = {
		id    = "montecarlo"
		token = "montecarlo"
	}
}

resource "montecarlo_iam_member" "test" {
	group = "groups/%[1]s"
	member = "user:%[2]s"
}
`, groupName, memberEmail)
}

func initIamMemberMonteCarloClient(groupName, groupNameUpdated, memberEmail, memberId string) client.MonteCarloClient {
	mcClient := cmock.MonteCarloClient{}
	mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetUsersInAccount"), mock.MatchedBy(func(in map[string]interface{}) bool {
		return in["email"] == memberEmail && in["after"] == (*string)(nil)
	})).Return(nil).Run(func(args mock.Arguments) {
		getResult := args.Get(1).(*client.GetUsersInAccount)
		getResult.GetUsersInAccount.Edges = []struct{ Node client.User }{
			{Node: client.User{CognitoUserId: memberId, IsSso: false, Email: memberEmail}},
		}
	})

	mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		getResult := args.Get(1).(*client.GetAuthorizationGroups)
		getResult.GetAuthorizationGroups = []client.AuthorizationGroup{{
			Name:               groupName,
			Label:              groupName,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
			DomainRestrictions: []struct{ Uuid string }{},
			SsoGroup:           nil,
		}}
	})

	// create operation
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateAuthorizationGroup"), mock.MatchedBy(func(in map[string]interface{}) bool {
		roles, rolesOk := in["roles"].([]string)
		domainRestrictions, domainRestrictionsOk := in["domainRestrictionIds"].([]client.UUID)
		memberUserIds, memberUserIdsOk := in["memberUserIds"].([]string)
		return in["name"] == groupName &&
			in["label"] == groupName &&
			in["description"] == "" &&
			memberUserIdsOk && len(memberUserIds) == 1 && memberUserIds[0] == memberId &&
			rolesOk && len(roles) == 1 && roles[0] == "mcd/owner" &&
			domainRestrictionsOk && len(domainRestrictions) == 0 &&
			in["ssoGroup"] == (*string)(nil)
	})).Return(nil).Run(func(args mock.Arguments) {
		createResult := args.Get(1).(*client.CreateOrUpdateAuthorizationGroup)
		createResult.CreateOrUpdateAuthorizationGroup.AuthorizationGroup = client.AuthorizationGroup{
			Name:               groupName,
			Label:              groupName,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
			DomainRestrictions: []struct{ Uuid string }{},
			SsoGroup:           nil,
			Users: []client.User{{
				CognitoUserId: memberId,
				Email:         memberEmail,
			}},
		}

		mcClient.ExpectedCalls = slices.DeleteFunc(mcClient.ExpectedCalls, func(call *mock.Call) bool {
			return call.Arguments.Is(mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything)
		})
		mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			getResult := args.Get(1).(*client.GetAuthorizationGroups)
			getResult.GetAuthorizationGroups = []client.AuthorizationGroup{{
				Name:               groupName,
				Label:              groupName,
				Description:        "",
				Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
				DomainRestrictions: []struct{ Uuid string }{},
				SsoGroup:           nil,
				Users: []client.User{{
					CognitoUserId: memberId,
					Email:         memberEmail,
				}},
			}}
		})
	})

	// delete operation
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateAuthorizationGroup"), mock.MatchedBy(func(in map[string]interface{}) bool {
		roles, rolesOk := in["roles"].([]string)
		domainRestrictions, domainRestrictionsOk := in["domainRestrictionIds"].([]client.UUID)
		memberUserIds, memberUserIdsOk := in["memberUserIds"].([]string)
		return in["name"] == groupName &&
			in["label"] == groupName &&
			in["description"] == "" &&
			memberUserIdsOk && len(memberUserIds) == 0 &&
			rolesOk && len(roles) == 1 && roles[0] == "mcd/owner" &&
			domainRestrictionsOk && len(domainRestrictions) == 0 &&
			in["ssoGroup"] == (*string)(nil)
	})).Return(nil).Run(func(args mock.Arguments) {
		createResult := args.Get(1).(*client.CreateOrUpdateAuthorizationGroup)
		createResult.CreateOrUpdateAuthorizationGroup.AuthorizationGroup = client.AuthorizationGroup{
			Name:               groupName,
			Label:              groupName,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
			DomainRestrictions: []struct{ Uuid string }{},
			SsoGroup:           nil,
			Users:              []client.User{},
		}

		mcClient.ExpectedCalls = slices.DeleteFunc(mcClient.ExpectedCalls, func(call *mock.Call) bool {
			return call.Arguments.Is(mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything)
		})
		mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			getResult := args.Get(1).(*client.GetAuthorizationGroups)
			getResult.GetAuthorizationGroups = []client.AuthorizationGroup{{
				Name:               groupName,
				Label:              groupName,
				Description:        "",
				Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
				DomainRestrictions: []struct{ Uuid string }{},
				SsoGroup:           nil,
				Users:              []client.User{},
			}, {
				Name:               groupNameUpdated,
				Label:              groupNameUpdated,
				Description:        "",
				Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
				DomainRestrictions: []struct{ Uuid string }{},
				SsoGroup:           nil,
				Users:              []client.User{},
			}}
		})
	})

	// update
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateAuthorizationGroup"), mock.MatchedBy(func(in map[string]interface{}) bool {
		roles, rolesOk := in["roles"].([]string)
		domainRestrictions, domainRestrictionsOk := in["domainRestrictionIds"].([]client.UUID)
		memberUserIds, memberUserIdsOk := in["memberUserIds"].([]string)
		return in["name"] == groupNameUpdated &&
			in["label"] == groupNameUpdated &&
			in["description"] == "" &&
			memberUserIdsOk && len(memberUserIds) == 1 && memberUserIds[0] == memberId &&
			rolesOk && len(roles) == 1 && roles[0] == "mcd/owner" &&
			domainRestrictionsOk && len(domainRestrictions) == 0 &&
			in["ssoGroup"] == (*string)(nil)
	})).Return(nil).Run(func(args mock.Arguments) {
		createResult := args.Get(1).(*client.CreateOrUpdateAuthorizationGroup)
		createResult.CreateOrUpdateAuthorizationGroup.AuthorizationGroup = client.AuthorizationGroup{
			Name:               groupNameUpdated,
			Label:              groupNameUpdated,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
			DomainRestrictions: []struct{ Uuid string }{},
			SsoGroup:           nil,
			Users: []client.User{{
				CognitoUserId: memberId,
				Email:         memberEmail,
			}},
		}

		mcClient.ExpectedCalls = slices.DeleteFunc(mcClient.ExpectedCalls, func(call *mock.Call) bool {
			return call.Arguments.Is(mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything)
		})
		mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			getResult := args.Get(1).(*client.GetAuthorizationGroups)
			getResult.GetAuthorizationGroups = []client.AuthorizationGroup{{
				Name:               groupNameUpdated,
				Label:              groupNameUpdated,
				Description:        "",
				Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
				DomainRestrictions: []struct{ Uuid string }{},
				SsoGroup:           nil,
				Users: []client.User{{
					CognitoUserId: memberId,
					Email:         memberEmail,
				}},
			}}
		})
	})

	// delete after update operation
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.CreateOrUpdateAuthorizationGroup"), mock.MatchedBy(func(in map[string]interface{}) bool {
		roles, rolesOk := in["roles"].([]string)
		domainRestrictions, domainRestrictionsOk := in["domainRestrictionIds"].([]client.UUID)
		memberUserIds, memberUserIdsOk := in["memberUserIds"].([]string)
		return in["name"] == groupNameUpdated &&
			in["label"] == groupNameUpdated &&
			in["description"] == "" &&
			memberUserIdsOk && len(memberUserIds) == 0 &&
			rolesOk && len(roles) == 1 && roles[0] == "mcd/owner" &&
			domainRestrictionsOk && len(domainRestrictions) == 0 &&
			in["ssoGroup"] == (*string)(nil)
	})).Return(nil).Run(func(args mock.Arguments) {
		createResult := args.Get(1).(*client.CreateOrUpdateAuthorizationGroup)
		createResult.CreateOrUpdateAuthorizationGroup.AuthorizationGroup = client.AuthorizationGroup{
			Name:               groupNameUpdated,
			Label:              groupNameUpdated,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
			DomainRestrictions: []struct{ Uuid string }{},
			SsoGroup:           nil,
			Users:              []client.User{},
		}

		mcClient.ExpectedCalls = slices.DeleteFunc(mcClient.ExpectedCalls, func(call *mock.Call) bool {
			return call.Arguments.Is(mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything)
		})
		mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetAuthorizationGroups"), mock.Anything).Return(nil).Run(func(args mock.Arguments) {
			getResult := args.Get(1).(*client.GetAuthorizationGroups)
			getResult.GetAuthorizationGroups = []client.AuthorizationGroup{{
				Name:               groupNameUpdated,
				Label:              groupNameUpdated,
				Description:        "",
				Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
				DomainRestrictions: []struct{ Uuid string }{},
				SsoGroup:           nil,
				Users:              []client.User{},
			}}
		})
	})
	return &mcClient
}
