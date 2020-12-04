package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	v0proto "github.com/owncloud/ocis/accounts/pkg/proto/v0"
	"go.opencensus.io/trace"
)

// NewTracing returns a service that instruments traces.
func NewTracing(next v0proto.AccountsServiceHandler) v0proto.AccountsServiceHandler {
	return tracing{
		next: next,
	}
}

type tracing struct {
	next v0proto.AccountsServiceHandler
}

func (t tracing) ListAccounts(ctx context.Context, req *v0proto.ListAccountsRequest, rsp *v0proto.ListAccountsResponse) error {
	ctx, span := trace.StartSpan(ctx, "Accounts.ListAccounts")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("page_size", int64(req.PageSize)),
		trace.StringAttribute("page_token", req.PageToken),
		trace.StringAttribute("field_mask", req.FieldMask.String()),
	}, "Execute Accounts.ListAccount handler")

	return t.next.ListAccounts(ctx, req, rsp)
}

func (t tracing) GetAccount(ctx context.Context, req *v0proto.GetAccountRequest, acc *v0proto.Account) error {
	ctx, span := trace.StartSpan(ctx, "Accounts.GetAccount")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("id", req.Id),
	}, "Execute Accounts.GetAccount handler")

	return t.next.GetAccount(ctx, req, acc)
}

func (t tracing) CreateAccount(ctx context.Context, req *v0proto.CreateAccountRequest, acc *v0proto.Account) error {
	ctx, span := trace.StartSpan(ctx, "Accounts.CreateAccount")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("account", req.Account.String()),
	}, "Execute Accounts.CreateAccount handler")

	return t.next.CreateAccount(ctx, req, acc)
}

func (t tracing) UpdateAccount(ctx context.Context, req *v0proto.UpdateAccountRequest, acc *v0proto.Account) error {
	ctx, span := trace.StartSpan(ctx, "Accounts.UpdateAccount")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("account", req.Account.String()),
	}, "Execute Accounts.UpdateAccount handler")

	return t.next.UpdateAccount(ctx, req, acc)
}

func (t tracing) DeleteAccount(ctx context.Context, req *v0proto.DeleteAccountRequest, e *empty.Empty) error {
	ctx, span := trace.StartSpan(ctx, "Accounts.DeleteAccount")
	defer span.End()

	span.Annotate([]trace.Attribute{
		trace.StringAttribute("id", req.Id),
	}, "Execute Accounts.DeleteAccout handler")
	return t.next.DeleteAccount(ctx, req, e)
}
