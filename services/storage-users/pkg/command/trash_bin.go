package command

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	userv1beta1 "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/mohae/deepcopy"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/ocis-pkg/tracing"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/event"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/logging"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	SKIP = iota
	REPLACE
	KEEP_BOTH

	retrievingErrorMsg = "trash-bin items retrieving error"
)

var optionFlagTmpl = cli.StringFlag{
	Name:        "option",
	Value:       "skip",
	Aliases:     []string{"o"},
	Usage:       "The restore option defines the behavior for a file to be restored, where the file name already already exists in the target space. Supported values are: 'skip', 'replace' and 'keep-both'.",
	DefaultText: "The default value is 'skip' overwriting an existing file",
}

// TrashBin wraps trash-bin related sub-commands.
func TrashBin(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "trash-bin",
		Usage: "manage trash-bin's",
		Subcommands: []*cli.Command{
			PurgeExpiredResources(cfg),
			listTrashBinItems(cfg),
			restoreAllTrashBinItems(cfg),
			restoreTrashBindItem(cfg),
		},
	}
}

// PurgeExpiredResources cli command removes old trash-bin items.
func PurgeExpiredResources(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "purge-expired",
		Usage: "Purge expired trash-bin items",
		Flags: []cli.Flag{},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			stream, err := event.NewStream(cfg)
			if err != nil {
				return err
			}

			if err := events.Publish(c.Context, stream, event.PurgeTrashBin{ExecutionTime: time.Now()}); err != nil {
				return err
			}

			// go-micro nats implementation uses async publishing,
			// therefore we need to manually wait.
			//
			// FIXME: upstream pr
			//
			// https://github.com/go-micro/plugins/blob/3e77393890683be4bacfb613bc5751867d584692/v4/events/natsjs/nats.go#L115
			time.Sleep(5 * time.Second)

			return nil
		},
	}
}

func listTrashBinItems(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "Print a list of all trash-bin items of a space.",
		ArgsUsage: "['userID' required] ['spaceID' required]",
		Flags:     []cli.Flag{},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			log := logging.Configure(cfg.Service.Name, cfg.Log)
			tp, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			var userID, spaceID string
			if c.NArg() > 1 {
				userID = c.Args().Get(0)
				spaceID = c.Args().Get(1)
			}
			if userID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return fmt.Errorf("userID is requered")
			}
			if spaceID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return fmt.Errorf("spaceID is requered")
			}
			fmt.Printf("Getting trash-bin items for spaceID: '%s' ...\n", spaceID)

			ref, err := storagespace.ParseReference(spaceID)
			if err != nil {
				return err
			}
			client, err := pool.GetGatewayServiceClient(cfg.RevaGatewayGRPCAddr)
			if err != nil {
				log.Error().Err(err).Msg("error selecting next gateway client")
				return err
			}
			ctx, _, err := utils.Impersonate(&userv1beta1.UserId{OpaqueId: userID}, client, cfg.MachineAuthAPIKey)
			if err != nil {
				log.Error().Err(err).Msg("could not impersonate")
				return err
			}

			spanOpts := []trace.SpanStartOption{
				trace.WithSpanKind(trace.SpanKindClient),
				trace.WithAttributes(
					attribute.KeyValue{Key: "userID", Value: attribute.StringValue(userID)},
					attribute.KeyValue{Key: "spaceID", Value: attribute.StringValue(spaceID)},
				),
			}
			ctx, span := tp.Tracer("storage-users trash-bin list").Start(ctx, "serve static asset", spanOpts...)
			defer span.End()

			res, err := client.ListRecycle(ctx, &provider.ListRecycleRequest{Ref: &ref, Key: "/"})
			if err != nil {
				log.Error().Err(err).Msg(retrievingErrorMsg)
				return err
			}
			if res.Status.Code != rpc.Code_CODE_OK {
				return fmt.Errorf("%s %s", retrievingErrorMsg, res.Status.Code)
			}

			if len(res.GetRecycleItems()) > 0 {
				fmt.Println("The list of the trash-bin items. Use an itemID to restore.")
			} else {
				fmt.Println("The list is empty.")
			}

			for _, item := range res.GetRecycleItems() {
				fmt.Printf("itemID: '%s', path: '%s', type: '%s',  delited at :%s\n", item.GetKey(), item.GetRef().GetPath(), itemType(item.GetType()), utils.TSToTime(item.GetDeletionTime()).UTC().Format(time.RFC3339))
			}
			return nil
		},
	}
}

func restoreAllTrashBinItems(cfg *config.Config) *cli.Command {
	var optionFlagVal string
	var overwriteOption int
	optionFlag := optionFlagTmpl
	optionFlag.Destination = &optionFlagVal
	return &cli.Command{
		Name:      "restore-all",
		Usage:     "Restore all trash-bin items for a space.",
		ArgsUsage: "['userID' required] ['spaceID' required]",
		Flags: []cli.Flag{
			&optionFlag,
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			log := logging.Configure(cfg.Service.Name, cfg.Log)
			tp, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			c.Lineage()
			var userID, spaceID string
			if c.NArg() > 1 {
				userID = c.Args().Get(0)
				spaceID = c.Args().Get(1)
			}
			if userID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return cli.Exit("The userID is required", 1)
			}
			if spaceID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return cli.Exit("The spaceID is required", 1)
			}
			switch optionFlagVal {
			case "skip":
				overwriteOption = SKIP
			case "replace":
				overwriteOption = REPLACE
			case "keep-both":
				overwriteOption = KEEP_BOTH
			default:
				_ = cli.ShowSubcommandHelp(c)
				return cli.Exit("The option flag is invalid", 1)
			}
			fmt.Printf("Restoring trash-bin items for spaceID: '%s' ...\n", spaceID)
			ref, err := storagespace.ParseReference(spaceID)
			if err != nil {
				return err
			}
			client, err := pool.GetGatewayServiceClient(cfg.RevaGatewayGRPCAddr)
			if err != nil {
				log.Error().Err(err).Msg("error selecting next gateway client")
				return err
			}
			ctx, _, err := utils.Impersonate(&userv1beta1.UserId{OpaqueId: userID}, client, cfg.MachineAuthAPIKey)
			if err != nil {
				log.Error().Err(err).Msg("could not impersonate")
				return err
			}

			spanOpts := []trace.SpanStartOption{
				trace.WithSpanKind(trace.SpanKindClient),
				trace.WithAttributes(
					attribute.KeyValue{Key: "option", Value: attribute.StringValue(optionFlagVal)},
					attribute.KeyValue{Key: "userID", Value: attribute.StringValue(userID)},
					attribute.KeyValue{Key: "spaceID", Value: attribute.StringValue(spaceID)},
				),
			}
			ctx, span := tp.Tracer("storage-users trash-bin restore-all").Start(ctx, "serve static asset", spanOpts...)
			defer span.End()

			res, err := client.ListRecycle(ctx, &provider.ListRecycleRequest{Ref: &ref, Key: "/"})
			if err != nil {
				log.Error().Err(err).Msg(retrievingErrorMsg)
				return err
			}
			if res.Status.Code != rpc.Code_CODE_OK {
				return fmt.Errorf("%s %s", retrievingErrorMsg, res.Status.Code)
			}
			if len(res.GetRecycleItems()) == 0 {
				return cli.Exit("The trash-bin is empty. Nothing to restore", 0)
			}

			for {
				fmt.Printf("Foud %d items that could be restored, continue (Y/n), show the items list (s): ", len(res.GetRecycleItems()))
				var i string
				_, err := fmt.Scanf("%s", &i)
				if err != nil {
					log.Err(err).Send()
					continue
				}
				if strings.ToLower(i) == "y" {
					break
				} else if strings.ToLower(i) == "n" {
					return nil
				} else if strings.ToLower(i) == "s" {
					for _, item := range res.GetRecycleItems() {
						fmt.Printf("itemID: '%s', path: '%s', type: '%s', delited at: %s\n", item.GetKey(), item.GetRef().GetPath(), itemType(item.GetType()), utils.TSToTime(item.GetDeletionTime()).UTC().Format(time.RFC3339))
					}
				}
			}

			fmt.Printf("\nRun restoring-all with option=%s\n", optionFlagVal)
			for _, item := range res.GetRecycleItems() {
				fmt.Printf("restoring itemID: '%s', path: '%s', type: '%s'\n", item.GetKey(), item.GetRef().GetPath(), itemType(item.GetType()))
				dstRes, err := restore(ctx, client, ref, item, overwriteOption)
				if err != nil {
					return err
				}
				fmt.Printf("itemID: '%s', path: '%s', restored as '%s'\n", item.GetKey(), item.GetRef().GetPath(), dstRes.GetPath())
			}
			return nil
		},
	}
}

func restoreTrashBindItem(cfg *config.Config) *cli.Command {
	var optionFlagVal string
	var overwriteOption int
	optionFlag := optionFlagTmpl
	optionFlag.Destination = &optionFlagVal
	return &cli.Command{
		Name:      "restore",
		Usage:     "Restore a trash-bin item by ID.",
		ArgsUsage: "['userId' required] ['spaceID' required] ['itemID' required]",
		Flags: []cli.Flag{
			&optionFlag,
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			log := logging.Configure(cfg.Service.Name, cfg.Log)
			tp, err := tracing.GetServiceTraceProvider(cfg.Tracing, cfg.Service.Name)
			if err != nil {
				return err
			}
			c.Lineage()
			var userID, spaceID, itemID string
			if c.NArg() > 2 {
				userID = c.Args().Get(0)
				spaceID = c.Args().Get(1)
				itemID = c.Args().Get(2)
			}
			if userID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return fmt.Errorf("userID is requered")
			}
			if spaceID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return fmt.Errorf("spaceID is requered")
			}
			if itemID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return fmt.Errorf("itemID is requered")
			}
			switch optionFlagVal {
			case "skip":
				overwriteOption = SKIP
			case "replace":
				overwriteOption = REPLACE
			case "keep-both":
				overwriteOption = KEEP_BOTH
			default:
				_ = cli.ShowSubcommandHelp(c)
				return cli.Exit("The option flag is invalid", 1)
			}

			ref, err := storagespace.ParseReference(spaceID)
			if err != nil {
				return err
			}
			client, err := pool.GetGatewayServiceClient(cfg.RevaGatewayGRPCAddr)
			if err != nil {
				log.Error().Err(err).Msg("error selecting gateway client")
				return err
			}
			ctx, _, err := utils.Impersonate(&userv1beta1.UserId{OpaqueId: userID}, client, cfg.MachineAuthAPIKey)
			if err != nil {
				log.Error().Err(err).Msg("could not impersonate")
				return err
			}

			spanOpts := []trace.SpanStartOption{
				trace.WithSpanKind(trace.SpanKindClient),
				trace.WithAttributes(
					attribute.KeyValue{Key: "option", Value: attribute.StringValue(optionFlagVal)},
					attribute.KeyValue{Key: "userID", Value: attribute.StringValue(userID)},
					attribute.KeyValue{Key: "spaceID", Value: attribute.StringValue(spaceID)},
					attribute.KeyValue{Key: "itemID", Value: attribute.StringValue(itemID)},
				),
			}
			ctx, span := tp.Tracer("storage-users trash-bin restore").Start(ctx, "serve static asset", spanOpts...)
			defer span.End()

			res, err := client.ListRecycle(ctx, &provider.ListRecycleRequest{Ref: &ref, Key: "/"})
			if err != nil {
				log.Error().Err(err).Msg(retrievingErrorMsg)
				return err
			}
			if res.Status.Code != rpc.Code_CODE_OK {
				return fmt.Errorf("%s %s", retrievingErrorMsg, res.Status.Code)
			}

			var found bool
			var itemRef *provider.RecycleItem
			for _, item := range res.GetRecycleItems() {
				if item.GetKey() == itemID {
					itemRef = item
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("itemID '%s' not found", itemID)
			}
			fmt.Printf("\nRun restoring with option=%s\n", optionFlagVal)
			fmt.Printf("restoring itemID: '%s', path: '%s', type: '%s'\n", itemRef.GetKey(), itemRef.GetRef().GetPath(), itemType(itemRef.GetType()))
			dstRes, err := restore(ctx, client, ref, itemRef, overwriteOption)
			if err != nil {
				return err
			}
			fmt.Printf("itemID: '%s', path: '%s', restored as '%s'\n", itemRef.GetKey(), itemRef.GetRef().GetPath(), dstRes.GetPath())
			return nil
		},
	}
}

func restore(ctx context.Context, client gateway.GatewayAPIClient, ref provider.Reference, item *provider.RecycleItem, overwriteOption int) (*provider.Reference, error) {
	dst, _ := deepcopy.Copy(ref).(provider.Reference)
	dst.Path = utils.MakeRelativePath(item.GetRef().GetPath())
	// Restore request
	req := &provider.RestoreRecycleItemRequest{
		Ref:        &ref,
		Key:        path.Join(item.GetKey(), "/"),
		RestoreRef: &dst,
	}

	exists, dstStatRes, err := isDestinationExists(ctx, client, dst)
	if err != nil {
		return &dst, err
	}

	if exists {
		fmt.Printf("destination '%s' exists.\n", dstStatRes.GetInfo().GetPath())
		switch overwriteOption {
		case SKIP:
			return &dst, nil
		case REPLACE:
			// delete existing tree
			delReq := &provider.DeleteRequest{Ref: &dst}
			delRes, err := client.Delete(ctx, delReq)
			if err != nil {
				return &dst, fmt.Errorf("error sending grpc delete request %w", err)
			}
			if delRes.Status.Code != rpc.Code_CODE_OK && delRes.Status.Code != rpc.Code_CODE_NOT_FOUND {
				return &dst, fmt.Errorf("deleting error %w", err)
			}
		case KEEP_BOTH:
			// modify the file name
			req.RestoreRef, err = resolveDestination(ctx, client, dst)
			if err != nil {
				return &dst, fmt.Errorf("trash-bin item restoring error %w", err)
			}
		}
	}

	res, err := client.RestoreRecycleItem(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("trash-bin item restoring error")
		return req.RestoreRef, err
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return req.RestoreRef, fmt.Errorf("trash-bin item restoring error %s", res.Status.Code)
	}
	return req.RestoreRef, nil
}

func resolveDestination(ctx context.Context, client gateway.GatewayAPIClient, dstRef provider.Reference) (*provider.Reference, error) {
	dst := dstRef
	for i := 1; i < 100; i++ {
		dst.Path = modifyFilename(dstRef.Path, i)
		exists, _, err := isDestinationExists(ctx, client, dst)
		if err != nil {
			return nil, err
		}
		if exists {
			continue
		}
		return &dst, nil
	}
	return nil, fmt.Errorf("too many attempts to resolve the destination")
}

func isDestinationExists(ctx context.Context, client gateway.GatewayAPIClient, dst provider.Reference) (bool, *provider.StatResponse, error) {
	dstStatReq := &provider.StatRequest{Ref: &dst}
	dstStatRes, err := client.Stat(ctx, dstStatReq)
	if err != nil {
		return false, nil, fmt.Errorf("error sending grpc stat request %w", err)
	}
	if dstStatRes.GetStatus().GetCode() == rpc.Code_CODE_OK {
		return true, dstStatRes, nil
	}
	if dstStatRes.GetStatus().GetCode() == rpc.Code_CODE_NOT_FOUND {
		return false, dstStatRes, nil
	}
	return false, dstStatRes, fmt.Errorf("stat request failed %s", dstStatRes.GetStatus())
}

// modify the file name like UI do
func modifyFilename(filename string, mod int) string {
	var extension string
	var found bool
	expected := []string{".tar.gz", ".tar.bz", ".tar.bz2"}
	for _, s := range expected {
		var prefix string
		prefix, found = strings.CutSuffix(strings.ToLower(filename), s)
		if found {
			extension = strings.TrimPrefix(filename, prefix)
			break
		}
	}
	if !found {
		extension = filepath.Ext(filename)
	}
	name := filename[0 : len(filename)-len(extension)]
	return fmt.Sprintf("%s (%d)%s", name, mod, extension)
}

func itemType(it provider.ResourceType) string {
	var itemType = "file"
	if it == provider.ResourceType_RESOURCE_TYPE_CONTAINER {
		itemType = "folder"
	}
	return itemType
}
