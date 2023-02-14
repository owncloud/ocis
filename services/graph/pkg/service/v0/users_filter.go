package svc

import (
	"context"
	"strings"

	"github.com/CiscoM31/godata"
	libregraph "github.com/owncloud/libre-graph-api-go"
)

func invalidFilterError() error {
	return godata.BadRequestError("invalid filter")
}

func unsupportedFilterError() error {
	return godata.NotImplementedError("unsupported filter")
}

func (g Graph) applyUserFilter(ctx context.Context, req *godata.GoDataRequest, root *godata.ParseNode) (users []*libregraph.User, err error) {
	logger := g.logger.SubloggerWithRequestID(ctx)

	if root == nil {
		root = req.Query.Filter.Tree
	}

	switch root.Token.Type {
	case godata.ExpressionTokenLambdaNav:
		return g.applyFilterLambda(ctx, req, root.Children)
	case godata.ExpressionTokenLogical:
		return g.applyFilterLogical(ctx, req, root)
	}
	logger.Debug().Str("filter", req.Query.Filter.RawValue).Msg("filter is not supported")
	return users, unsupportedFilterError()
}

func (g Graph) applyFilterLogical(ctx context.Context, req *godata.GoDataRequest, root *godata.ParseNode) (users []*libregraph.User, err error) {
	logger := g.logger.SubloggerWithRequestID(ctx)
	if root.Token.Type != godata.ExpressionTokenLogical {
		return users, invalidFilterError()
	}

	switch root.Token.Value {
	case "and":
		// 'and' needs 2 operands
		if len(root.Children) != 2 {
			return users, invalidFilterError()
		}
		return g.applyFilterLogicalAnd(ctx, req, root.Children[0], root.Children[1])
	}
	logger.Debug().Str("Token", root.Token.Value).Msg("unsupported logical filter")
	return users, unsupportedFilterError()
}

func (g Graph) applyFilterLogicalAnd(ctx context.Context, req *godata.GoDataRequest, operand1 *godata.ParseNode, operand2 *godata.ParseNode) (users []*libregraph.User, err error) {
	logger := g.logger.SubloggerWithRequestID(ctx)
	results := make([][]*libregraph.User, 0, 2)

	for _, node := range []*godata.ParseNode{operand1, operand2} {
		res, err := g.applyUserFilter(ctx, req, node)
		if err != nil {
			return users, err
		}
		logger.Debug().Interface("subfilter", res).Msg("result part")
		results = append(results, res)
	}

	// 'results' contains two slices of libregraph.Users now turn one of them
	// into a map for efficiently getting the intersection of both slices
	userSet := userSliceToMap(results[0])
	var filteredUsers []*libregraph.User
	for _, user := range results[1] {
		if _, found := userSet[user.GetId()]; found {
			filteredUsers = append(filteredUsers, user)
		}
	}
	return filteredUsers, nil
}

func (g Graph) applyFilterLambda(ctx context.Context, req *godata.GoDataRequest, nodes []*godata.ParseNode) (users []*libregraph.User, err error) {
	logger := g.logger.SubloggerWithRequestID(ctx)
	if len(nodes) != 2 {
		return users, invalidFilterError()
	}
	// We only support the "any" operator for lambda queries for now
	if nodes[1].Token.Type != godata.ExpressionTokenLambda || nodes[1].Token.Value != "any" {
		logger.Debug().Str("Token", nodes[1].Token.Value).Msg("unsupported lambda filter")
		return users, unsupportedFilterError()
	}
	if nodes[0].Token.Type != godata.ExpressionTokenLiteral {
		return users, unsupportedFilterError()
	}
	switch nodes[0].Token.Value {
	case "memberOf":
		return g.applyLambdaMemberOfAny(ctx, req, nodes[1].Children)
	case "appRoleAssignments":
		return g.applyLambdaAppRoleAssignmentAny(ctx, req, nodes[1].Children)
	}
	logger.Debug().Str("Token", nodes[0].Token.Value).Msg("unsupported relation for lambda filter")
	return users, unsupportedFilterError()
}

func (g Graph) applyLambdaMemberOfAny(ctx context.Context, req *godata.GoDataRequest, nodes []*godata.ParseNode) (users []*libregraph.User, err error) {
	if len(nodes) != 2 {
		return users, invalidFilterError()
	}

	// First element is the "name" of the lambda function's parameter
	if nodes[0].Token.Type != godata.ExpressionTokenLiteral {
		return users, invalidFilterError()
	}

	// We only support the 'eq' expression for now
	if nodes[1].Token.Type != godata.ExpressionTokenLogical || nodes[1].Token.Value != "eq" {
		return users, unsupportedFilterError()
	}
	return g.applyMemberOfEq(ctx, req, nodes[1].Children)
}

func (g Graph) applyMemberOfEq(ctx context.Context, req *godata.GoDataRequest, nodes []*godata.ParseNode) (users []*libregraph.User, err error) {
	logger := g.logger.SubloggerWithRequestID(ctx)
	if len(nodes) != 2 {
		return users, invalidFilterError()
	}

	if nodes[0].Token.Type != godata.ExpressionTokenNav {
		return users, invalidFilterError()
	}

	if len(nodes[0].Children) != 2 {
		return users, invalidFilterError()
	}

	switch nodes[0].Children[1].Token.Value {
	case "id":
		var filterValue string
		switch nodes[1].Token.Type {
		case godata.ExpressionTokenGuid:
			filterValue = nodes[1].Token.Value
		case godata.ExpressionTokenString:
			// unquote
			filterValue = strings.Trim(nodes[1].Token.Value, "'")
		default:
			return users, unsupportedFilterError()
		}
		logger.Debug().Str("property", nodes[0].Children[1].Token.Value).Str("value", filterValue).Msg("Filtering memberOf by group id")
		return g.identityBackend.GetGroupMembers(ctx, filterValue, req)
	default:
		return users, unsupportedFilterError()
	}
}

func (g Graph) applyLambdaAppRoleAssignmentAny(ctx context.Context, req *godata.GoDataRequest, nodes []*godata.ParseNode) (users []*libregraph.User, err error) {
	if len(nodes) != 2 {
		return users, invalidFilterError()
	}

	// First element is the "name" of the lambda function's parameter
	if nodes[0].Token.Type != godata.ExpressionTokenLiteral {
		return users, invalidFilterError()
	}

	// We only support the 'eq' expression for now
	if nodes[1].Token.Type != godata.ExpressionTokenLogical || nodes[1].Token.Value != "eq" {
		return users, unsupportedFilterError()
	}
	return g.applyAppRoleAssignmentEq(ctx, req, nodes[1].Children)
}

func (g Graph) applyAppRoleAssignmentEq(ctx context.Context, req *godata.GoDataRequest, nodes []*godata.ParseNode) (users []*libregraph.User, err error) {
	logger := g.logger.SubloggerWithRequestID(ctx)
	if len(nodes) != 2 {
		return users, invalidFilterError()
	}

	if nodes[0].Token.Type != godata.ExpressionTokenNav {
		return users, invalidFilterError()
	}

	if len(nodes[0].Children) != 2 {
		return users, invalidFilterError()
	}

	if nodes[0].Children[1].Token.Value == "appRoleId" {
		var filterValue string
		switch nodes[1].Token.Type {
		case godata.ExpressionTokenGuid:
			filterValue = nodes[1].Token.Value
		case godata.ExpressionTokenString:
			// unquote
			filterValue = strings.Trim(nodes[1].Token.Value, "'")
		default:
			return users, unsupportedFilterError()
		}

		logger.Debug().Str("property", nodes[0].Children[1].Token.Value).Str("value", filterValue).Msg("Filtering appRoleAssignments by appRoleId")
		if users, err = g.identityBackend.GetUsers(ctx, req); err != nil {
			return users, err
		}

		return g.filterUsersByAppRoleId(ctx, filterValue, users)
	}
	return users, unsupportedFilterError()
}

func (g Graph) filterUsersByAppRoleId(ctx context.Context, appRoleId string, users []*libregraph.User) ([]*libregraph.User, error) {
	// We're using a map for the results here, in order to avoid returning
	// a user twice. The settings API, still has an issue that causes it to
	// duplicate some assignments on restart:
	// https://github.com/owncloud/ocis/issues/3432
	resultUsersMap := make(map[string]*libregraph.User, len(users))
	for _, user := range users {
		assignments, err := g.fetchAppRoleAssignments(ctx, user.GetId())
		if err != nil {
			return users, err
		}
		for _, assignment := range assignments {
			if assignment.GetAppRoleId() == appRoleId {
				if _, ok := resultUsersMap[user.GetId()]; !ok {
					resultUsersMap[user.GetId()] = user
				}
			}
		}
	}
	resultUsers := make([]*libregraph.User, 0, len(resultUsersMap))
	for _, user := range resultUsersMap {
		resultUsers = append(resultUsers, user)
	}
	return resultUsers, nil
}

func userSliceToMap(users []*libregraph.User) map[string]*libregraph.User {
	resMap := make(map[string]*libregraph.User, len(users))
	for _, user := range users {
		resMap[user.GetId()] = user
	}
	return resMap
}
