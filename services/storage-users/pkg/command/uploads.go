package command

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"github.com/shamaton/msgpack/v2"
	"github.com/urfave/cli/v2"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	provider "github.com/cs3org/go-cs3apis/cs3/storage/provider/v1beta1"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/event"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/revaconfig"
	"github.com/owncloud/reva/v2/pkg/events"
	"github.com/owncloud/reva/v2/pkg/storage"
	"github.com/owncloud/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/owncloud/reva/v2/pkg/storage/utils/decomposedfs/node"
	"github.com/owncloud/reva/v2/pkg/utils"
)

const (
	MSGPACK_KEY_USER_OCIS_NODESTATUS = "user.ocis.nodestatus"

	// Log indentation levels
	LOG_INDENT_L1 = "  " // 2 spaces
	LOG_INDENT_L2 = LOG_INDENT_L1 + LOG_INDENT_L1
	LOG_INDENT_L3 = LOG_INDENT_L2 + LOG_INDENT_L1
)

// Session contains the information of an upload session
type Session struct {
	ID         string         `json:"id"`
	Space      string         `json:"space"`
	Filename   string         `json:"filename"`
	Offset     int64          `json:"offset"`
	Size       int64          `json:"size"`
	Executant  userpb.UserId  `json:"executant"`
	SpaceOwner *userpb.UserId `json:"spaceowner,omitempty"`
	Expires    time.Time      `json:"expires"`
	Processing bool           `json:"processing"`
	ScanDate   time.Time      `json:"virus_scan_date"`
	ScanResult string         `json:"virus_scan_result"`
}

// Uploads is the entry point for the uploads command
func Uploads(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "uploads",
		Usage: "manage unfinished uploads",
		Subcommands: []*cli.Command{
			ListUploadSessions(cfg),
			DeleteStaleProcessingNodes(cfg),
			Consistency(cfg),
		},
	}
}

// ListUploadSessions prints a list of upload sessiens
func ListUploadSessions(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "sessions",
		Usage: "Print a list of upload sessions",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "id",
				DefaultText: "unset",
				Usage:       "filter sessions by upload session id",
			},
			&cli.BoolFlag{
				Name:        "processing",
				DefaultText: "unset",
				Usage:       "filter sessions by processing status",
			},
			&cli.BoolFlag{
				Name:        "expired",
				DefaultText: "unset",
				Usage:       "filter sessions by expired status",
			},
			&cli.BoolFlag{
				Name:        "has-virus",
				DefaultText: "unset",
				Usage:       "filter sessions by virus scan result",
			},
			&cli.BoolFlag{
				Name:  "json",
				Usage: "output as json",
			},
			&cli.BoolFlag{
				Name:  "restart",
				Usage: "send restart event for all listed sessions. Only one of resume/restart/clean can be set.",
			},
			&cli.BoolFlag{
				Name:  "resume",
				Usage: "send resume event for all listed sessions. Only one of resume/restart/clean can be set.",
			},
			&cli.BoolFlag{
				Name:  "clean",
				Usage: "remove uploads for all listed sessions. Only one of resume/restart/clean can be set.",
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(c *cli.Context) error {
			f, ok := registry.NewFuncs[cfg.Driver]
			if !ok {
				fmt.Fprintf(os.Stderr, "Unknown filesystem driver '%s'\n", cfg.Driver)
				os.Exit(1)
			}
			drivers := revaconfig.StorageProviderDrivers(cfg)
			fs, err := f(drivers[cfg.Driver].(map[string]interface{}), nil, nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize filesystem driver '%s'\n", cfg.Driver)
				return err
			}

			managingFS, ok := fs.(storage.UploadSessionLister)
			if !ok {
				fmt.Fprintf(os.Stderr, "'%s' storage does not support listing upload sessions\n", cfg.Driver)
				os.Exit(1)
			}

			var stream events.Stream
			if c.Bool("restart") || c.Bool("resume") || c.Bool("clean") {
				stream, err = event.NewStream(cfg)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to create event stream: %v\n", err)
					os.Exit(1)
				}
			}

			filter := buildFilter(c)
			uploads, err := managingFS.ListUploadSessions(c.Context, filter)
			if err != nil {
				return err
			}

			var (
				table *tablewriter.Table
				raw   []*Session
			)

			if !c.Bool("json") {
				fmt.Println(buildInfo(filter))

				table = tablewriter.NewTable(os.Stdout)
				table.Header("Space", "Upload Id", "Name", "Offset", "Size", "Executant", "Owner", "Expires", "Processing", "Scan Date", "Scan Result")
			}

			for _, u := range uploads {
				ref := u.Reference()
				sr, sd := u.ScanData()

				session := Session{
					Space:      ref.GetResourceId().GetSpaceId(),
					ID:         u.ID(),
					Filename:   u.Filename(),
					Offset:     u.Offset(),
					Size:       u.Size(),
					Executant:  u.Executant(),
					SpaceOwner: u.SpaceOwner(),
					Expires:    u.Expires(),
					Processing: u.IsProcessing(),
					ScanDate:   sd,
					ScanResult: sr,
				}

				if c.Bool("json") {
					raw = append(raw, &session)
				} else {
					table.Append([]string{
						session.Space,
						session.ID,
						session.Filename,
						strconv.FormatInt(session.Offset, 10),
						strconv.FormatInt(session.Size, 10),
						session.Executant.OpaqueId,
						session.SpaceOwner.GetOpaqueId(),
						session.Expires.Format(time.RFC3339),
						strconv.FormatBool(session.Processing),
						session.ScanDate.Format(time.RFC3339),
						session.ScanResult,
					})
				}

				switch {
				case c.Bool("restart"):
					if err := events.Publish(context.Background(), stream, events.RestartPostprocessing{
						UploadID:  u.ID(),
						Timestamp: utils.TSNow(),
					}); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to send restart event for upload session '%s'\n", u.ID())
						// if publishing fails there is no need to try publishing other events - they will fail too.
						os.Exit(1)
					}

				case c.Bool("resume"):
					if err := events.Publish(context.Background(), stream, events.ResumePostprocessing{
						UploadID:  u.ID(),
						Timestamp: utils.TSNow(),
					}); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to send resume event for upload session '%s'\n", u.ID())
						// if publishing fails there is no need to try publishing other events - they will fail too.
						os.Exit(1)
					}

				case c.Bool("clean"):
					if err := events.Publish(context.Background(), stream, events.CleanUpload{
						UploadID:  u.ID(),
						Timestamp: utils.TSNow(),
					}); err != nil {
						fmt.Fprintf(os.Stderr, "Failed to send clean upload event for upload session '%s'\n", u.ID())
						// if publishing fails there is no need to try publishing other events - they will fail too.
						os.Exit(1)
					}
				}

			}

			if !c.Bool("json") {
				table.Render()
				return nil
			}

			j, err := json.Marshal(raw)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println(string(j))
			return nil
		},
	}
}

func buildFilter(c *cli.Context) storage.UploadSessionFilter {
	filter := storage.UploadSessionFilter{}
	if c.IsSet("processing") {
		processingValue := c.Bool("processing")
		filter.Processing = &processingValue
	}
	if c.IsSet("expired") {
		expiredValue := c.Bool("expired")
		filter.Expired = &expiredValue
	}
	if c.IsSet("has-virus") {
		infectedValue := c.Bool("has-virus")
		filter.HasVirus = &infectedValue
	}
	if c.IsSet("id") {
		idValue := c.String("id")
		filter.ID = &idValue
	}
	return filter
}

func buildInfo(filter storage.UploadSessionFilter) string {
	var b strings.Builder
	if filter.Processing != nil {
		if !*filter.Processing {
			b.WriteString("Not ")
		}
		if b.Len() == 0 {
			b.WriteString("Processing")
		} else {
			b.WriteString("processing")
		}
	}

	if filter.Expired != nil {
		if b.Len() != 0 {
			b.WriteString(", ")
		}
		if !*filter.Expired {
			if b.Len() == 0 {
				b.WriteString("Not ")
			} else {
				b.WriteString("not ")
			}
		}
		if b.Len() == 0 {
			b.WriteString("Expired")
		} else {
			b.WriteString("expired")
		}
	}

	if filter.HasVirus != nil {
		if b.Len() != 0 {
			b.WriteString(", ")
		}
		if !*filter.HasVirus {
			if b.Len() == 0 {
				b.WriteString("Not ")
			} else {
				b.WriteString("not ")
			}
		}
		if b.Len() == 0 {
			b.WriteString("Virusinfected")
		} else {
			b.WriteString("virusinfected")
		}
	}

	if b.Len() == 0 {
		b.WriteString("Session")
	} else {
		b.WriteString(" session")
	}

	if filter.ID != nil {
		b.WriteString(" with id '" + *filter.ID + "'")
	} else {
		// to make `session` plural
		b.WriteString("s")
	}

	b.WriteString(":")
	return b.String()
}

// DeleteStaleProcessingNodes is the entry point for the delete-stale-nodes command
func DeleteStaleProcessingNodes(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "delete-stale-nodes",
		Usage: "Delete (or revert) all nodes in processing state that are not referenced by any upload session",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "spaceid",
				Usage:    "Space ID to check for processing nodes (omit to check all spaces)",
				Required: false,
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "Only show what would be deleted without actually deleting",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Usage: "Enable verbose logging",
				Value: false,
			},
		},
		Before: func(c *cli.Context) error {
			return configlog.ReturnFatal(parser.ParseConfig(cfg))
		},
		Action: func(ctx *cli.Context) error {
			spaceIDs := []string{}
			dryRun := ctx.Bool("dry-run")
			verbose := ctx.Bool("verbose")
			start := time.Now()

			// Check if specific space ID provided
			if ctx.IsSet("spaceid") {
				spaceIDs = append(spaceIDs, ctx.String("spaceid"))
			} else {
				fmt.Println("Scanning all spaces for stale processing nodes...")
				spaceIDs = globSpaceIDs(cfg)
			}

			if verbose {
				fmt.Printf("Spaces to cleanup: %d\n", len(spaceIDs))
				for _, spaceID := range spaceIDs {
					fmt.Printf("  - %s\n", spaceID)
				}
			}

			var stream events.Stream
			if !dryRun {
				s, err := event.NewStream(cfg)
				if err != nil {
					log.Fatalf("Failed to create event stream: %v", err)
				}
				stream = s
			}

			staleCount := 0
			for _, spaceID := range spaceIDs {
				staleCount += deleteStaleUploads(cfg, spaceID, dryRun, verbose, stream)
			}

			if verbose {
				fmt.Printf("Took %ds\n", int(time.Since(start).Seconds()))
			}
			fmt.Printf("Total stale nodes: %d\n", staleCount)

			return nil
		},
	}
}

// globSpaceIDs returns a list of all space IDs in the storage root
func globSpaceIDs(cfg *config.Config) []string {
	fsys := os.DirFS(cfg.Drivers.OCIS.Root)
	dirs, err := fs.Glob(fsys, "spaces/*/*/nodes")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error globbing spaces root directory %s: %v\n", cfg.Drivers.OCIS.Root, err)
		return []string{}
	}

	spaceIDs := []string{}
	for _, dir := range dirs {
		// For dir i.e. spaces/9d/408cec-8f0a-4d33-8715-89df1217a10c/nodes
		// spaceID is 9d408cec-8f0a-4d33-8715-89df1217a10c
		spaceIDs = append(spaceIDs, strings.ReplaceAll(strings.TrimSuffix(strings.TrimPrefix(dir, "spaces/"), "/nodes"), "/", ""))
	}
	return spaceIDs
}

// delete stale processing nodes for a given spaceID
func deleteStaleUploads(cfg *config.Config, spaceID string, dryRun bool, verbose bool, stream events.Stream) int {
	if verbose {
		fmt.Printf("\nDeleting stale processing nodes for space: %s\n", spaceID)
	}

	// Find .mpk files in space directory
	spaceRoot := filepath.Join(cfg.Drivers.OCIS.Root, "spaces", lookup.Pathify(spaceID, 1, 2))
	mpkFiles := []string{}
	err := filepath.Walk(spaceRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing path %s: %s\n", path, err)
			return filepath.SkipDir
		}
		if !info.IsDir() && strings.HasSuffix(path, ".mpk") {
			mpkFiles = append(mpkFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error walking space directory %s: %s\n", spaceRoot, err)
		return 0
	}

	if verbose {
		fmt.Printf("%sFound total %d .mpk files\n", LOG_INDENT_L1, len(mpkFiles))
	}

	staleCount := 0
	for _, path := range mpkFiles {
		staleCount += deleteStaleNode(cfg, path, dryRun, verbose, stream)
	}

	if verbose {
		fmt.Printf("%sFound total %d stale nodes\n", LOG_INDENT_L1, staleCount)
	}

	return staleCount
}

// deleteStaleNode deletes a stale node: if it is not referenced by any upload session
// returns 1 if the node stale node was detected for deletion, 0 otherwise, for counting purposes
func deleteStaleNode(cfg *config.Config, path string, dryRun bool, verbose bool, stream events.Stream) int {
	nodeDir := filepath.Dir(path)

	// Read .mpk file to get processing info
	b, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file %s: %s\n", path, err)
		return 0
	}
	var mpkData map[string]interface{}
	if err := msgpack.Unmarshal(b, &mpkData); err != nil {
		fmt.Fprintf(os.Stderr, "Error unmarshaling file %s: %s\n", path, err)
		return 0
	}

	processingID := extractProcessingID(mpkData)
	if processingID == "" {
		return 0
	}

	// Construct path to upload info file:
	// i.e. ~/.ocis/storage/users/uploads/5329c14b-b786-4b27-8f7d-7429f03009d7.info
	// And pass only the .info file not exists: err is ErrNotExist
	pathUploadInfo := filepath.Join(cfg.Drivers.OCIS.Root, "uploads", processingID) + ".info"
	_, infoStatErr := os.Stat(pathUploadInfo)
	if infoStatErr == nil {
		return 0
	}
	if !os.IsNotExist(infoStatErr) {
		// Tere was an error other than file not existing, log and return
		fmt.Fprintf(os.Stderr, "Error checking upload info %s: %s\n", pathUploadInfo, infoStatErr)
		return 0
	}

	if verbose {
		fmt.Printf("%sFound stale upload at %s (Processing ID: %s)\n", LOG_INDENT_L1, path, processingID)
		fmt.Printf("%sUpload info missing at: %s\n", LOG_INDENT_L2, pathUploadInfo)
	}

	if dryRun {
		return 1
	}

	rid := extractResourceID(strings.TrimSuffix(path, ".mpk"))
	if rid == nil {
		fmt.Fprintf(os.Stderr, "Failed to extract resource ID from path %s\n", path)
		return 0
	}

	if err := events.Publish(context.Background(), stream, events.RevertRevision{
		ResourceID: rid,
		Timestamp:  utils.TSNow(),
	}); err != nil {
		// if publishing fails there is no need to try publishing other events - they will fail too.
		log.Fatalf("Failed to send revert revision event for node '%s'\n", path)
	}

	if verbose {
		fmt.Printf("%sDeleted stale node: %s\n", LOG_INDENT_L2, nodeDir)
	}

	return 1
}

func extractProcessingID(mpkData map[string]interface{}) string {
	processingID := ""
	for k, v := range mpkData {
		vStr := string(v.([]byte))
		if k == MSGPACK_KEY_USER_OCIS_NODESTATUS && strings.Contains(vStr, node.ProcessingStatus) {
			processingID = strings.Split(vStr, ":")[1]
			break
		}
	}
	return processingID
}

func extractResourceID(path string) *provider.ResourceId {
	// path looks like /.../storage/users/spaces/f2/06bccf-0f10-4070-9e63-40943f060667/nodes/5b/ba/1e/a7/-f185-4f31-8342-ed4b5743f096
	parts := strings.Split(path, "spaces")
	if len(parts) < 2 {
		return nil
	}

	spaceParts := strings.Split(parts[1], "nodes")
	if len(spaceParts) < 2 {
		return nil
	}

	return &provider.ResourceId{
		SpaceId:  strings.ReplaceAll(spaceParts[0], "/", ""),
		OpaqueId: strings.ReplaceAll(spaceParts[1], "/", ""),
	}
}
