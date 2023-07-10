package opa

import (
	"github.com/open-policy-agent/opa/topdown/print"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
)

type logPrinter struct {
	logger log.Logger
}

func (lp logPrinter) Print(_ print.Context, msg string) error {
	lp.logger.Info().Msg(msg)
	return nil
}
