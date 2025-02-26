package kql

import (
	"strings"

	"github.com/owncloud/ocis/v2/ocis-pkg/ast"
)

// connectNodes connects given nodes
func connectNodes(c Connector, nodes ...ast.Node) []ast.Node {
	var connectedNodes []ast.Node

	for i := range nodes {
		ri := len(nodes) - 1 - i
		head := nodes[ri]

		if connectionNodes := connectNode(c, head, connectedNodes...); len(connectionNodes) > 0 {
			connectedNodes = append(connectionNodes, connectedNodes...)
		}

		connectedNodes = append([]ast.Node{head}, connectedNodes...)
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
	// the connection is of type OR
	//
	// spec: same
	//		author:"John Smith" author:"Jane Smith"
	//		author:"John Smith" OR author:"Jane Smith"
	//
	// if the nodes have NO key, the edge is a AND connection
	//
	// spec: same
	//		cat dog
	//		cat AND dog
	// from the spec:
	// 		To construct complex queries, you can combine multiple
	// 		free-text expressions with KQL query operators.
	// 		If there are multiple free-text expressions without any
	// 		operators in between them, the query behavior is the same
	// 		as using the AND operator.
	//
	// nodes inside of group node are handled differently,
	// if no explicit operator given, it uses AND
	//
	// spec: same
	// 		author:"John Smith" AND author:"Jane Smith"
	// 		author:("John Smith" "Jane Smith")
	if headKey == neighborKey && headKey != "" && neighborKey != "" {
		connection.Value = c.sameKeyOPValue
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

			// if neighbor node negotiates, an AND edge is needed
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
