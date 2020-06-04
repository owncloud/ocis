package provider

import (
	"fmt"
	"strings"

	"github.com/CiscoM31/godata"
	"github.com/owncloud/ocis-accounts/pkg/config"
	"gopkg.in/ldap.v2"
)

func init() {
	// add (ap)prox filter
	godata.GlobalFilterTokenizer = FilterTokenizer()
	godata.GlobalFilterParser.DefineOperator("ap", 2, godata.OpAssociationLeft, 4, false)
}

// LDAPNodeMap is used to convert query tokens into ldap filters according to https://tools.ietf.org/search/rfc4515
var LDAPNodeMap = map[string]string{
	// 11.2.6.1.1 Built-in Filter Operations according to http://docs.oasis-open.org/odata/odata/v4.01/odata-v4.01-part1-protocol.html#_Toc31358949

	// Comparison Operators
	"eq": "(%s=%s)",    // -> LDAP equal
	"ne": "(!(%s=%s))", // -> LDAP NOT equal
	//"gt": "(&(%s>=%s)(!(%s=%s)))", // -> TODO can be constructed but requires more parameters
	"ge": "(%s>=%s)", // -> LDAP greaterorequal
	//"lt": "(&(%s<=%s)(!(%s=%s)))", // -> TODO can be constructed but requires more parameters
	"le": "(%s<=%s)", // -> LDAP lessorequal
	//"has": "(%s=*)",   // -> TODO LDAP present but in odata has looks like "Style has Sales.Color'Yellow'"
	//"in": "???", // TODO

	// additional native LDAP Search String Filter Definition according to https://tools.ietf.org/search/rfc4515#section-3
	"ap": "(%s~=%s)", // approx, TODO needs token in parser, odata uses $search instead of $filter for fuzzy search

	// Logical Operators
	// While LDAP understands logical filters like (&()()()()) we leave that as an optimization and use at max two params
	"and": "(&%s%s)",
	"or":  "(|%s%s)",
	"not": "(!%s)",

	// Arithmetic operators
	//"add": ""
	//"sub": ""
	//"mul": ""
	//"div": ""
	//"divby": ""
	//"mod": ""

	// Grouping operators

	// 11.2.6.1.2 Built-in Query Functions according to http://docs.oasis-open.org/odata/odata/v4.01/odata-v4.01-part1-protocol.html#sec_BuiltinQueryFunctions
	//String and Collection Functions
	//"concat":           "CONCAT(%s,%s)",
	"contains": "(%s=*%s*)",
	"endswith": "(%s=*%s)",
	//"indexof":    "LOCATE(%s)",
	//"length":     "LENGTH(%s)",
	"startswith": "(%s=%s*)",
	//"substring": "",

	//
	//"tolower":          "LOWER(%s)",
	//"toupper":          "UPPER(%s)",
	//"trim":             "TRIM(%s)",
	//"year":             "YEAR(%s)",
	//"month":            "MONTH(%s)",
	//"day":              "DAY(%s)",
	//"hour":             "HOUR(%s)",
	//"minute":           "MINUTE(%s)",
	//"second":           "SECOND(%s)",
	//"fractionalsecond": "MICROSECOND(%s)",
	//"date":             "DATE(%s)",
	//"time":             "TIME(%s)",
	//"totaloffsetminutes": "",
	//"now": "NOW()",
	//"maxdatetime":"",
	//"mindatetime":"",
	//"totalseconds":"",
	//"round":   "ROUND(%s)",
	//"floor":   "FLOOR(%s)",
	//"ceiling": "CEIL(%s)",
	//"isof": "", // TODO objectclass=
	//"cast": "",
	//"geo.distance": "",
	//"geo.intersects": "",
	//"geo.length": "",
	//"any": "",
	//"all": "",
	//"null": "NULL",
}

// BuildLDAPFilter converts a GoDataFilterQuery into an ldap filter
func BuildLDAPFilter(r *godata.GoDataFilterQuery, c *config.LDAPSchema) (string, error) {
	return recursiveBuildFilter(r.Tree, c)
}

// Builds the filter recursively using DFS
func recursiveBuildFilter(n *godata.ParseNode, c *config.LDAPSchema) (string, error) {
	if n.Token.Type == godata.FilterTokenLiteral {
		switch n.Token.Value {
		case "accountid":
			return c.AccountID, nil
		case "displayname":
			return c.DisplayName, nil
		case "username":
			return c.Username, nil
		case "mail":
			return c.Mail, nil
		case "groups":
			// TODO groups
			return "", godata.NotImplementedError(n.Token.Value + " is not implemented.")
		case "identities":
			// TODO identities
			return "", godata.NotImplementedError(n.Token.Value + " is not implemented.")
		}
		return "", godata.BadRequestError("unknown property " + n.Token.Value)
	}
	if n.Token.Type == godata.FilterTokenString {
		// without leading and ending ' required by odata
		// encode LDAP safe
		return strings.TrimSuffix(strings.TrimPrefix(ldap.EscapeFilter(n.Token.Value), "'"), "'"), nil
	}
	if n.Token.Type == godata.FilterTokenInteger {
		return n.Token.Value, nil
	}
	if n.Token.Type == godata.FilterTokenFloat {
		return n.Token.Value, nil
	}

	if v, ok := LDAPNodeMap[n.Token.Value]; ok {
		children := []interface{}{}
		// build each child first using DFS
		for _, child := range n.Children {
			f, err := recursiveBuildFilter(child, c)
			if err != nil {
				return "", err
			}
			children = append(children, f)
		}
		// merge together the children and the current node
		result := fmt.Sprintf(v, children...)
		return result, nil
	}
	return "", godata.NotImplementedError(n.Token.Value + " is not implemented.")
}

// FilterTokenizer creates a tokenizer capable of tokenizing filter statements
// TODO disable tokens we don't handle anyway
func FilterTokenizer() *godata.Tokenizer {
	t := godata.Tokenizer{}
	t.Add("^[0-9]{4,4}-[0-9]{2,2}-[0-9]{2,2}T[0-9]{2,2}:[0-9]{2,2}(:[0-9]{2,2}(.[0-9]+)?)?(Z|[+-][0-9]{2,2}:[0-9]{2,2})", godata.FilterTokenDateTime)
	t.Add("^-?[0-9]{4,4}-[0-9]{2,2}-[0-9]{2,2}", godata.FilterTokenDate)
	t.Add("^[0-9]{2,2}:[0-9]{2,2}(:[0-9]{2,2}(.[0-9]+)?)?", godata.FilterTokenTime)
	t.Add("^\\(", godata.FilterTokenOpenParen)
	t.Add("^\\)", godata.FilterTokenCloseParen)
	t.Add("^/", godata.FilterTokenNav)
	t.Add("^:", godata.FilterTokenColon)
	t.Add("^,", godata.FilterTokenComma)
	t.Add("^(geo.distance|geo.intersects|geo.length)", godata.FilterTokenFunc)
	t.Add("^(substringof|substring|length|indexof)", godata.FilterTokenFunc)
	// only change from the global tokenizer is the added ap
	t.Add("^(eq|ne|gt|ge|lt|le|and|or|not|has|in|ap)", godata.FilterTokenLogical)
	t.Add("^(add|sub|mul|divby|div|mod)", godata.FilterTokenOp)
	t.Add("^(contains|endswith|startswith|tolower|toupper|"+
		"trim|concat|year|month|day|hour|minute|second|fractionalseconds|date|"+
		"time|totaloffsetminutes|now|maxdatetime|mindatetime|totalseconds|round|"+
		"floor|ceiling|isof|cast)", godata.FilterTokenFunc)
	t.Add("^(any|all)", godata.FilterTokenLambda)
	t.Add("^null", godata.FilterTokenNull)
	t.Add("^\\$it", godata.FilterTokenIt)
	t.Add("^\\$root", godata.FilterTokenRoot)
	t.Add("^-?[0-9]+\\.[0-9]+", godata.FilterTokenFloat)
	t.Add("^-?[0-9]+", godata.FilterTokenInteger)
	t.Add("^'(''|[^'])*'", godata.FilterTokenString)
	t.Add("^(true|false)", godata.FilterTokenBoolean)
	t.Add("^@*[a-zA-Z][a-zA-Z0-9_.]*", godata.FilterTokenLiteral) // The optional '@' character is used to identify parameter aliases
	t.Ignore("^ ", godata.FilterTokenWhitespace)

	return &t
}
