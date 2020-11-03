package indexer

type queryTree struct {
	token *token
	root  bool
	left  *queryTree
	right *queryTree
}

// token to be resolved by the index
type token struct {
	operator   string // original OData operator. i.e: 'startswith', `or`, `and`.
	filterType string // equivalent operator from OData -> indexer i.e FindByPartial or FindBy.
	operands   []string
}

// newQueryTree constructs a new tree with a root node.
func newQueryTree() queryTree {
	return queryTree{
		root: true,
	}
}

// insert populates first the LHS of the tree first, if this is not possible it fills the RHS.
func (t *queryTree) insert(tkn *token) {
	if t != nil && t.root {
		t.left = &queryTree{token: tkn}
		return
	}

	if t.left == nil {
		t.left = &queryTree{token: tkn}
		return
	}

	if t.left != nil && t.right == nil {
		t.right = &queryTree{token: tkn}
		return
	}
}
