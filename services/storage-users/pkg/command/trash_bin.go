package command

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	gateway "github.com/cs3org/go-cs3apis/cs3/gateway/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/cs3org/reva/v2/pkg/events"
	"github.com/cs3org/reva/v2/pkg/rgrpc/todo/pool"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/cs3org/reva/v2/pkg/utils"
	"github.com/mohae/deepcopy"
	tw "github.com/olekukonko/tablewriter"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	zlog "github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/event"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
)

const (
	SKIP = iota
	REPLACE
	KEEP_BOTH
)

var _optionFlagTmpl = cli.StringFlag{
	Name:        "option",
	Value:       "skip",
	Aliases:     []string{"o"},
	Usage:       "The restore option defines the behavior for a file to be restored, where the file name already already exists in the target space. Supported values are: 'skip', 'replace' and 'keep-both'.",
	DefaultText: "The default value is 'skip' overwriting an existing file",
}

var _verboseFlagTmpl = cli.BoolFlag{
	Name:    "verbose",
	Aliases: []string{"v"},
	Usage:   "Get more verbose output",
}

var _applyYesFlagTmpl = cli.BoolFlag{
	Name:    "yes",
	Aliases: []string{"y"},
	Usage:   "Automatic yes to prompts. Assume 'yes' as answer to all prompts and run non-interactively.",
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
	var verboseVal bool
	verboseFlag := _verboseFlagTmpl
	verboseFlag.Destination = &verboseVal
	return &cli.Command{
		Name:      "list",
		Usage:     "Print a list of all trash-bin items of a space.",
		ArgsUsage: "['spaceID' required]",
		Flags: []cli.Flag{
			&verboseFlag,
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			log := cliLogger(verboseVal)
			var spaceID string
			if c.NArg() > 0 {
				spaceID = c.Args().Get(0)
			}
			if spaceID == "" {
				_ = cli.ShowSubcommandHelp(c)
				return fmt.Errorf("spaceID is requered")
			}
			log.Info().Msgf("Getting trash-bin items for spaceID: '%s' ...", spaceID)

			ref, err := storagespace.ParseReference(spaceID)
			if err != nil {
				return err
			}
			client, err := pool.GetGatewayServiceClient(cfg.RevaGatewayGRPCAddr)
			if err != nil {
				return fmt.Errorf("error selecting gateway client %w", err)
			}
			ctx, err := utils.GetServiceUserContext(cfg.ServiceAccount.ServiceAccountID, client, cfg.ServiceAccount.ServiceAccountSecret)
			if err != nil {
				return fmt.Errorf("could not get service user context %w", err)
			}
			res, err := listRecycle(ctx, client, ref)
			if err != nil {
				return err
			}

			table := itemsTable(len(res.GetRecycleItems()))
			for _, item := range res.GetRecycleItems() {
				table.Append([]string{item.GetKey(), item.GetRef().GetPath(), itemType(item.GetType()), utils.TSToTime(item.GetDeletionTime()).UTC().Format(time.RFC3339)})
			}
			table.Render()
			fmt.Println("Use an itemID to restore an item.")
			return nil
		},
	}
}

func restoreAllTrashBinItems(cfg *config.Config) *cli.Command {
	var optionFlagVal string
	var overwriteOption int
	optionFlag := _optionFlagTmpl
	optionFlag.Destination = &optionFlagVal
	var verboseVal bool
	verboseFlag := _verboseFlagTmpl
	verboseFlag.Destination = &verboseVal
	var applyYesVal bool
	applyYesFlag := _applyYesFlagTmpl
	applyYesFlag.Destination = &applyYesVal
	return &cli.Command{
		Name:      "restore-all",
		Usage:     "Restore all trash-bin items for a space.",
		ArgsUsage: "['spaceID' required]",
		Flags: []cli.Flag{
			&optionFlag,
			&verboseFlag,
			&applyYesFlag,
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			log := cliLogger(verboseVal)
			var spaceID string
			if c.NArg() > 0 {
				spaceID = c.Args().Get(0)
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
			log.Info().Msgf("Restoring trash-bin items for spaceID: '%s' ...", spaceID)

			ref, err := storagespace.ParseReference(spaceID)
			if err != nil {
				return err
			}
			client, err := pool.GetGatewayServiceClient(cfg.RevaGatewayGRPCAddr)
			if err != nil {
				return fmt.Errorf("error selecting gateway client %w", err)
			}
			ctx, err := utils.GetServiceUserContext(cfg.ServiceAccount.ServiceAccountID, client, cfg.ServiceAccount.ServiceAccountSecret)
			if err != nil {
				return fmt.Errorf("could not get service user context %w", err)
			}
			res, err := listRecycle(ctx, client, ref)
			if err != nil {
				return err
			}

			if !applyYesVal {
				for {
					fmt.Printf("Found %d items that could be restored, continue (Y/n), show the items list (s): ", len(res.GetRecycleItems()))
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
						table := itemsTable(len(res.GetRecycleItems()))
						for _, item := range res.GetRecycleItems() {
							table.Append([]string{item.GetKey(), item.GetRef().GetPath(), itemType(item.GetType()), utils.TSToTime(item.GetDeletionTime()).UTC().Format(time.RFC3339)})
						}
						table.Render()
					}
				}
			}

			log.Info().Msgf("Run restoring-all with option=%s", optionFlagVal)
			for _, item := range res.GetRecycleItems() {
				log.Info().Msgf("restoring itemID: '%s', path: '%s', type: '%s'", item.GetKey(), item.GetRef().GetPath(), itemType(item.GetType()))
				dstRes, err := restore(ctx, client, ref, item, overwriteOption, cfg.CliMaxAttemptsRenameFile, log)
				if err != nil {
					log.Err(err).Msg("trash-bin item restoring error")
					continue
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
	optionFlag := _optionFlagTmpl
	optionFlag.Destination = &optionFlagVal
	var verboseVal bool
	verboseFlag := _verboseFlagTmpl
	verboseFlag.Destination = &verboseVal
	return &cli.Command{
		Name:      "restore",
		Usage:     "Restore a trash-bin item by ID.",
		ArgsUsage: "['spaceID' required] ['itemID' required]",
		Flags: []cli.Flag{
			&optionFlag,
			&verboseFlag,
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			log := cliLogger(verboseVal)
			var spaceID, itemID string
			if c.NArg() > 1 {
				spaceID = c.Args().Get(0)
				itemID = c.Args().Get(1)
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
			log.Info().Msgf("Restoring trash-bin item for spaceID: '%s' itemID: '%s' ...", spaceID, itemID)

			ref, err := storagespace.ParseReference(spaceID)
			if err != nil {
				return err
			}
			client, err := pool.GetGatewayServiceClient(cfg.RevaGatewayGRPCAddr)
			if err != nil {
				return fmt.Errorf("error selecting gateway client %w", err)
			}
			ctx, err := utils.GetServiceUserContext(cfg.ServiceAccount.ServiceAccountID, client, cfg.ServiceAccount.ServiceAccountSecret)
			if err != nil {
				return fmt.Errorf("could not get service user context %w", err)
			}
			res, err := listRecycle(ctx, client, ref)
			if err != nil {
				return err
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
			log.Info().Msgf("Run restoring with option=%s", optionFlagVal)
			log.Info().Msgf("restoring itemID: '%s', path: '%s', type: '%s", itemRef.GetKey(), itemRef.GetRef().GetPath(), itemType(itemRef.GetType()))
			dstRes, err := restore(ctx, client, ref, itemRef, overwriteOption, cfg.CliMaxAttemptsRenameFile, log)
			if err != nil {
				return err
			}
			fmt.Printf("itemID: '%s', path: '%s', restored as '%s'\n", itemRef.GetKey(), itemRef.GetRef().GetPath(), dstRes.GetPath())
			return nil
		},
	}
}

func listRecycle(ctx context.Context, client gateway.GatewayAPIClient, ref provider.Reference) (*provider.ListRecycleResponse, error) {
	_retrievingErrorMsg := "trash-bin items retrieving error"
	res, err := client.ListRecycle(ctx, &provider.ListRecycleRequest{Ref: &ref, Key: ""})
	if err != nil {
		return nil, fmt.Errorf("%s %w", _retrievingErrorMsg, err)
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return nil, fmt.Errorf("%s %s", _retrievingErrorMsg, res.Status.Code)
	}
	if len(res.GetRecycleItems()) == 0 {
		return res, cli.Exit("The trash-bin is empty. Nothing to restore", 0)
	}
	return res, nil
}

func restore(ctx context.Context, client gateway.GatewayAPIClient, ref provider.Reference, item *provider.RecycleItem, overwriteOption int, maxRenameAttempt int, log zlog.Logger) (*provider.Reference, error) {
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
		log.Info().Msgf("destination '%s' exists.", dstStatRes.GetInfo().GetPath())
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
			req.RestoreRef, err = resolveDestination(ctx, client, dst, maxRenameAttempt)
			if err != nil {
				return &dst, err
			}
		}
	}

	res, err := client.RestoreRecycleItem(ctx, req)
	if err != nil {
		return req.RestoreRef, fmt.Errorf("restoring error  %w", err)
	}
	if res.Status.Code != rpc.Code_CODE_OK {
		return req.RestoreRef, fmt.Errorf("can not restore %s", res.Status.Code)
	}
	return req.RestoreRef, nil
}

func resolveDestination(ctx context.Context, client gateway.GatewayAPIClient, dstRef provider.Reference, maxRenameAttempt int) (*provider.Reference, error) {
	dst := dstRef
	if maxRenameAttempt < 100 {
		maxRenameAttempt = 100
	}
	for i := 1; i < maxRenameAttempt; i++ {
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

func itemsTable(total int) *tw.Table {
	table := tw.NewWriter(os.Stdout)
	table.SetHeader([]string{"itemID", "path", "type", "delete at"})
	table.SetAutoFormatHeaders(false)
	table.SetFooter([]string{"", "", "", "total count: " + strconv.Itoa(total)})
	return table
}

func cliLogger(verbose bool) zlog.Logger {
	logLvl := zerolog.ErrorLevel
	if verbose {
		logLvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true}
	return zlog.Logger{zerolog.New(output).With().Timestamp().Logger().Level(logLvl)}
}
