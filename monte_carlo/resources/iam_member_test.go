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
	readUser := mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetUsersInAccount"), mock.MatchedBy(func(in map[string]interface{}) bool {
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
		}, {
			Name:               groupNameUpdated,
			Label:              groupNameUpdated,
			Description:        "",
			Roles:              []struct{ Name string }{{Name: "mcd/owner"}},
			DomainRestrictions: []struct{ Uuid string }{},
			SsoGroup:           nil,
		}}
	})

	// create operation
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.UpdateUserAuthorizationGroupMembership"), mock.MatchedBy(func(in map[string]interface{}) bool {
		groupNames, groupNamesOk := in["groupNames"].([]string)
		return in["memberUserId"] == memberId && groupNamesOk && len(groupNames) == 1 && groupNames[0] == groupName
	})).Return(nil).Run(func(args mock.Arguments) {
		mcClient.ExpectedCalls = slices.DeleteFunc(mcClient.ExpectedCalls, func(call *mock.Call) bool { return call.Arguments.Is(readUser.Arguments...) })
		readUser = mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetUsersInAccount"), mock.MatchedBy(func(in map[string]interface{}) bool {
			return in["email"] == memberEmail && in["after"] == (*string)(nil)
		})).Return(nil).Run(func(args mock.Arguments) {
			getResult := args.Get(1).(*client.GetUsersInAccount)
			node := client.User{CognitoUserId: memberId, IsSso: false, Email: memberEmail}
			node.Auth.Groups = []string{groupName}
			getResult.GetUsersInAccount.Edges = []struct{ Node client.User }{{Node: node}}
		})
	})

	// delete operation
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.UpdateUserAuthorizationGroupMembership"), mock.MatchedBy(func(in map[string]interface{}) bool {
		groupNames, groupNamesOk := in["groupNames"].([]string)
		return in["memberUserId"] == memberId && groupNamesOk && len(groupNames) == 0
	})).Return(nil)

	// update
	mcClient.On("Mutate", mock.Anything, mock.AnythingOfType("*client.UpdateUserAuthorizationGroupMembership"), mock.MatchedBy(func(in map[string]interface{}) bool {
		groupNames, groupNamesOk := in["groupNames"].([]string)
		return in["memberUserId"] == memberId && groupNamesOk && slices.Contains(groupNames, groupNameUpdated)
	})).Return(nil).Run(func(args mock.Arguments) {
		mcClient.ExpectedCalls = slices.DeleteFunc(mcClient.ExpectedCalls, func(call *mock.Call) bool { return call.Arguments.Is(readUser.Arguments...) })
		readUser = mcClient.On("Query", mock.Anything, mock.AnythingOfType("*client.GetUsersInAccount"), mock.MatchedBy(func(in map[string]interface{}) bool {
			return in["email"] == memberEmail && in["after"] == (*string)(nil)
		})).Return(nil).Run(func(args mock.Arguments) {
			getResult := args.Get(1).(*client.GetUsersInAccount)
			node := client.User{CognitoUserId: memberId, IsSso: false, Email: memberEmail}
			node.Auth.Groups = []string{groupNameUpdated}
			getResult.GetUsersInAccount.Edges = []struct{ Node client.User }{{Node: node}}
		})
	})
	return &mcClient
}
