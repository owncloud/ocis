// Copyright 2018-2022 CERN
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// In applying this license, CERN does not waive the privileges and immunities
// granted to it by virtue of its status as an Intergovernmental Organization
// or submit itself to any jurisdiction.

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
