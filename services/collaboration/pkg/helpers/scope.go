package helpers

import (
	"fmt"
	"strings"

	auth "github.com/cs3org/go-cs3apis/cs3/auth/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/golang/protobuf/proto"
)

// GetScopeByKeyPrefix returns the scope from the AccessToken Scope map by key prefix
func GetScopeByKeyPrefix(scopes map[string]*auth.Scope, keyPrefix string, m proto.Message) error {
	for k, v := range scopes {
		if strings.HasPrefix(k, keyPrefix) && v.Resource.Decoder == "json" {
			err := utils.UnmarshalJSONToProtoV1(v.Resource.Value, m)
			if err != nil {
				return fmt.Errorf("can't unmarshal public share from scope: %w", err)
			}
			return nil
		}
	}
	return fmt.Errorf("scope %s not found", keyPrefix)
}
