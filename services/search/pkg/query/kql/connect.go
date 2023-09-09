package kql

import (
	"strings"

	"github.com/owncloud/ocis/v2/services/search/pkg/query/ast"
)

// connectNodes connects given nodes
func connectNodes(c Connector, nodes ...ast.Node) []ast.Node {
	var connectedNodes []ast.Node

	for i := range nodes {
		ri := len(nodes) - 1 - i
		head := nodes[ri]
		pair := []ast.Node{head}

		if connectionNodes := connectNode(c, pair[0], connectedNodes...); len(connectionNodes) >= 1 {
			pair = append(pair, connectionNodes...)
		}

		connectedNodes = append(pair, connectedNodes...)
	}

	return connectedNodes
}

// connectNode connects a tip node with the rest
func connectNode(c Connector, headNode ast.Node, tailNodes ...ast.Node) []ast.Node {
	var nearestNeighborNode ast.Node
	var nearestNeighborOperators []*ast.OperatorNode

l:
	for _, tailNode := range tailNodes {
		switch node := tailNode.(type) {
		case *ast.OperatorNode:
			nearestNeighborOperators = append(nearestNeighborOperators, node)
		default:
			nearestNeighborNode = node
			break l
		}
	}

	if nearestNeighborNode == nil {
		return nil
	}

	return c.Connect(headNode, nearestNeighborNode, nearestNeighborOperators)
}

// Connector is responsible to decide what node connections are needed
type Connector interface {
	Connect(head ast.Node, neighbor ast.Node, connections []*ast.OperatorNode) []ast.Node
}

// DefaultConnector is the default node connector
type DefaultConnector struct {
	sameKeyOPValue string
}

// Connect implements the Connector interface and is used to connect the nodes using
// the default logic defined by the kql spec.
func (c DefaultConnector) Connect(head ast.Node, neighbor ast.Node, connections []*ast.OperatorNode) []ast.Node {
	switch head.(type) {
	case *ast.OperatorNode:
		return nil
	}

	headKey := strings.ToLower(ast.NodeKey(head))
	neighborKey := strings.ToLower(ast.NodeKey(neighbor))

	connection := &ast.OperatorNode{
		Base:  &ast.Base{Loc: &ast.Location{Source: &[]string{"implicitly operator"}[0]}},
		Value: BoolAND,
	}

	// if the current node and the neighbor node have the same key
	// the connection is of type OR, same applies if no keys are in place
	//
	//		"" == ""
	//
	// spec: same
	//		author:"John Smith" author:"Jane Smith"
	//		author:"John Smith" OR author:"Jane Smith"
	//
	// nodes inside of group nodes are handled differently,
	// if no explicit operator give, it uses OR
	//
	// spec: same
	// 		author:"John Smith" AND author:"Jane Smith"
	// 		author:("John Smith" "Jane Smith")
	if headKey == neighborKey {
		connection.Value = c.sameKeyOPValue
	}

	// decisions based on nearest neighbor node
	switch neighbor.(type) {
	// nearest neighbor node type could change the default case
	// docs says, if the next value node:
	//
	//		is a group AND has no key
	//
	// even if the current node has none too, which normal leads to SAME KEY OR
	//
	// 		it should be an AND edge
	//
	// spec: same
	// 		cat (dog OR fox)
	// 		cat AND (dog OR fox)
	//
	// note:
	// 		sounds contradictory to me
	case *ast.GroupNode:
		if headKey == "" && neighborKey == "" {
			connection.Value = BoolAND
		}
	}

	// decisions based on nearest neighbor operators
	for i, node := range connections {
		// consider direct neighbor operator only
		if i == 0 {
			// no connection is necessary here because an `AND` or `OR` edge is already present
			// exit
			for _, skipValue := range []string{BoolOR, BoolAND} {
				if node.Value == skipValue {
					return nil
				}
			}

			// if neighbor node negotiates, AND edge is needed
			//
			// spec: same
			// 		cat -dog
			// 		cat AND NOT dog
			if node.Value == BoolNOT {
				connection.Value = BoolAND
			}
		}
	}

	return []ast.Node{connection}
}
