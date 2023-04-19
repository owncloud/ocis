package godata

type GoDataIdentifier map[string]string

type RequestKind int

const (
	RequestKindUnknown RequestKind = iota
	RequestKindMetadata
	RequestKindService
	RequestKindEntity
	RequestKindCollection
	RequestKindSingleton
	RequestKindProperty
	RequestKindPropertyValue
	RequestKindRef
	RequestKindCount
)

type SemanticType int

const (
	SemanticTypeUnknown SemanticType = iota
	SemanticTypeEntity
	SemanticTypeEntitySet
	SemanticTypeDerivedEntity
	SemanticTypeAction
	SemanticTypeFunction
	SemanticTypeProperty
	SemanticTypePropertyValue
	SemanticTypeRef
	SemanticTypeCount
	SemanticTypeMetadata
)

type GoDataRequest struct {
	FirstSegment *GoDataSegment
	LastSegment  *GoDataSegment
	Query        *GoDataQuery
	RequestKind  RequestKind
}

// Represents a segment (slash-separated) part of the URI path. Each segment
// has a link to the next segment (the last segment precedes nil).
type GoDataSegment struct {
	// The raw segment parsed from the URI
	RawValue string

	// The kind of resource being pointed at by this segment
	SemanticType SemanticType

	// A pointer to the metadata type this object represents, as defined by a
	// particular service
	SemanticReference interface{}

	// The name of the entity, type, collection, etc.
	Name string

	// map[string]string of identifiers passed to this segment. If the identifier
	// is not key/value pair(s), then all values will be nil. If there is no
	// identifier, it will be nil.
	Identifier *GoDataIdentifier

	// The next segment in the path.
	Next *GoDataSegment
	// The previous segment in the path.
	Prev *GoDataSegment
}

type GoDataQuery struct {
	Filter      *GoDataFilterQuery
	At          *GoDataFilterQuery
	Apply       *GoDataApplyQuery
	Expand      *GoDataExpandQuery
	Select      *GoDataSelectQuery
	OrderBy     *GoDataOrderByQuery
	Top         *GoDataTopQuery
	Skip        *GoDataSkipQuery
	Count       *GoDataCountQuery
	InlineCount *GoDataInlineCountQuery
	Search      *GoDataSearchQuery
	Compute     *GoDataComputeQuery
	Format      *GoDataFormatQuery
}

// GoDataExpression encapsulates the tree representation of an expression
// as defined in the OData ABNF grammar.
type GoDataExpression struct {
	Tree *ParseNode
	// The raw string representing an expression
	RawValue string
}

// Stores a parsed version of the filter query string. Can be used by
// providers to apply the filter based on their own implementation. The filter
// is stored as a parse tree that can be traversed.
type GoDataFilterQuery struct {
	Tree *ParseNode
	// The raw filter string
	RawValue string
}

type GoDataApplyQuery string

type GoDataExpandQuery struct {
	ExpandItems []*ExpandItem
}

type GoDataSelectQuery struct {
	SelectItems []*SelectItem
	// The raw select string
	RawValue string
}

type GoDataOrderByQuery struct {
	OrderByItems []*OrderByItem
	// The raw orderby string
	RawValue string
}

type GoDataTopQuery int

type GoDataSkipQuery int

type GoDataCountQuery bool

type GoDataInlineCountQuery string

type GoDataSearchQuery struct {
	Tree *ParseNode
	// The raw search string
	RawValue string
}

type GoDataComputeQuery struct {
	ComputeItems []*ComputeItem
	// The raw compute string
	RawValue string
}

type GoDataFormatQuery struct {
}

// Check if this identifier has more than one key/value pair.
func (id *GoDataIdentifier) HasMultiple() bool {
	count := 0
	for range map[string]string(*id) {
		count++
	}
	return count > 1
}

// Return the first key in the map. This is how you should get the identifier
// for single values, e.g. when the path is Employee(1), etc.
func (id *GoDataIdentifier) Get() string {
	for k := range map[string]string(*id) {
		return k
	}
	return ""
}

// Return a specific value for a specific key.
func (id *GoDataIdentifier) GetKey(key string) (string, bool) {
	v, ok := map[string]string(*id)[key]
	return v, ok
}

// GoDataCommonStructure represents either a GoDataQuery or ExpandItem in a uniform manner
// as a Go interface. This allows the writing of functional logic that can work on either type,
// such as a provider implementation which starts at the GoDataQuery and walks any nested ExpandItem
// in an identical manner.
type GoDataCommonStructure interface {
	GetFilter() *GoDataFilterQuery
	GetAt() *GoDataFilterQuery
	GetApply() *GoDataApplyQuery
	GetExpand() *GoDataExpandQuery
	GetSelect() *GoDataSelectQuery
	GetOrderBy() *GoDataOrderByQuery
	GetTop() *GoDataTopQuery
	GetSkip() *GoDataSkipQuery
	GetCount() *GoDataCountQuery
	GetInlineCount() *GoDataInlineCountQuery
	GetSearch() *GoDataSearchQuery
	GetCompute() *GoDataComputeQuery
	GetFormat() *GoDataFormatQuery
	// AddExpandItem adds an item to the list of expand clauses in the underlying GoDataQuery/ExpandItem
	// structure.
	// AddExpandItem may be used to add items based on the requirements of the application using godata.
	// For example applications may support the introduction of dynamic navigational fields using $compute.
	// A possible implementation is to parse the request url using godata and then during semantic
	// post-processing identify dynamic navigation properties and call AddExpandItem to add them to the
	// list of expanded fields.
	AddExpandItem(*ExpandItem)
}

// GoDataQuery implementation of GoDataCommonStructure interface
func (o *GoDataQuery) GetFilter() *GoDataFilterQuery           { return o.Filter }
func (o *GoDataQuery) GetAt() *GoDataFilterQuery               { return o.At }
func (o *GoDataQuery) GetApply() *GoDataApplyQuery             { return o.Apply }
func (o *GoDataQuery) GetExpand() *GoDataExpandQuery           { return o.Expand }
func (o *GoDataQuery) GetSelect() *GoDataSelectQuery           { return o.Select }
func (o *GoDataQuery) GetOrderBy() *GoDataOrderByQuery         { return o.OrderBy }
func (o *GoDataQuery) GetTop() *GoDataTopQuery                 { return o.Top }
func (o *GoDataQuery) GetSkip() *GoDataSkipQuery               { return o.Skip }
func (o *GoDataQuery) GetCount() *GoDataCountQuery             { return o.Count }
func (o *GoDataQuery) GetInlineCount() *GoDataInlineCountQuery { return o.InlineCount }
func (o *GoDataQuery) GetSearch() *GoDataSearchQuery           { return o.Search }
func (o *GoDataQuery) GetCompute() *GoDataComputeQuery         { return o.Compute }
func (o *GoDataQuery) GetFormat() *GoDataFormatQuery           { return o.Format }

// AddExpandItem adds an expand clause to the toplevel odata request structure 'o'.
func (o *GoDataQuery) AddExpandItem(item *ExpandItem) {
	if o.Expand == nil {
		o.Expand = &GoDataExpandQuery{}
	}
	o.Expand.ExpandItems = append(o.Expand.ExpandItems, item)
}

// ExpandItem implementation of GoDataCommonStructure interface
func (o *ExpandItem) GetFilter() *GoDataFilterQuery           { return o.Filter }
func (o *ExpandItem) GetAt() *GoDataFilterQuery               { return o.At }
func (o *ExpandItem) GetApply() *GoDataApplyQuery             { return nil }
func (o *ExpandItem) GetExpand() *GoDataExpandQuery           { return o.Expand }
func (o *ExpandItem) GetSelect() *GoDataSelectQuery           { return o.Select }
func (o *ExpandItem) GetOrderBy() *GoDataOrderByQuery         { return o.OrderBy }
func (o *ExpandItem) GetTop() *GoDataTopQuery                 { return o.Top }
func (o *ExpandItem) GetSkip() *GoDataSkipQuery               { return o.Skip }
func (o *ExpandItem) GetCount() *GoDataCountQuery             { return nil }
func (o *ExpandItem) GetInlineCount() *GoDataInlineCountQuery { return nil }
func (o *ExpandItem) GetSearch() *GoDataSearchQuery           { return o.Search }
func (o *ExpandItem) GetCompute() *GoDataComputeQuery         { return o.Compute }
func (o *ExpandItem) GetFormat() *GoDataFormatQuery           { return nil }

// AddExpandItem adds an expand clause to 'o' creating a nested expand, ie $expand 'item'
// nested within $expand 'o'.
func (o *ExpandItem) AddExpandItem(item *ExpandItem) {
	if o.Expand == nil {
		o.Expand = &GoDataExpandQuery{}
	}
	o.Expand.ExpandItems = append(o.Expand.ExpandItems, item)
}
