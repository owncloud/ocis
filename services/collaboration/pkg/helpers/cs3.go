package helpers

import (
	"crypto/sha256"
	"encoding/hex"

	gatewayv1beta1 "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	providerv1beta1 "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/services/collaboration/pkg/config"
)

var commonCS3ApiClient gatewayv1beta1.GatewayAPIClient

// GatewayAPIClient gets an instance based on the provided configuration.
// The instance will be cached and returned if possible, unless the "forceNew"
// parameter is set to true. In this case, the old instance will be replaced
// with the new one if there is no error.
func GetCS3apiClient(cfg *config.Config, forceNew bool) (gatewayv1beta1.GatewayAPIClient, error) {
	// establish a connection to the cs3 api endpoint
	// in this case a REVA gateway, started by oCIS
	if commonCS3ApiClient != nil && !forceNew {
		return commonCS3ApiClient, nil
	}

	client, err := pool.GetGatewayServiceClient(cfg.CS3Api.Gateway.Name)
	if err == nil {
		commonCS3ApiClient = client
	}
	return client, err
}

// HashResourceId builds a urlsafe and stable file reference that can be used for proxy routing,
// so that all sessions on one file end on the same office server
func HashResourceId(resourceId *providerv1beta1.ResourceId) string {
	c := sha256.New()
	c.Write([]byte(storagespace.FormatResourceID(resourceId)))
	return hex.EncodeToString(c.Sum(nil))
}
