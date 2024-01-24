package command

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"

	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/ocis/v2/ocis-pkg/config/configlog"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/revaconfig"
)

func Uploads(cfg *config.Config) *cli.Command {
	return &cli.Command{

		Name:  "uploads",
		Usage: "manage unfinished uploads",
		Subcommands: []*cli.Command{
			ListUploads(cfg),
			ListUploadSessions(cfg),
			PurgeExpiredUploads(cfg),
		},
	}
}

// ListUploads prints a list of all incomplete uploads
func ListUploads(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "Print a list of all incomplete uploads (deprecated, use sessions)",
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
			fs, err := f(drivers[cfg.Driver].(map[string]interface{}), nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize filesystem driver '%s'\n", cfg.Driver)
				return err
			}

			managingFS, ok := fs.(storage.UploadSessionLister)
			if !ok {
				fmt.Fprintf(os.Stderr, "'%s' storage does not support listing upload sessions\n", cfg.Driver)
				os.Exit(1)
			}
			expired := false
			uploads, err := managingFS.ListUploadSessions(c.Context, storage.UploadSessionFilter{Expired: &expired})
			if err != nil {
				return err
			}

			fmt.Println("Incomplete uploads:")
			for _, u := range uploads {
				ref := u.Reference()
				fmt.Printf(" - %s (Space: %s, Name: %s, Size: %d/%d, Expires: %s, Processing: %t)\n", ref.GetResourceId().GetSpaceId(), u.ID(), u.Filename(), u.Offset(), u.Size(), u.Expires(), u.IsProcessing())
			}
			return nil
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
			&cli.StringFlag{
				Name:  "output",
				Usage: "output format to use (can be 'plain' or 'json', experimental)",
				Value: "plain",
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
			fs, err := f(drivers[cfg.Driver].(map[string]interface{}), nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize filesystem driver '%s'\n", cfg.Driver)
				return err
			}

			managingFS, ok := fs.(storage.UploadSessionLister)
			if !ok {
				fmt.Fprintf(os.Stderr, "'%s' storage does not support listing upload sessions\n", cfg.Driver)
				os.Exit(1)
			}
			var b strings.Builder
			filter := storage.UploadSessionFilter{}
			if c.IsSet("processing") {
				processingValue := c.Bool("processing")
				filter.Processing = &processingValue
				if !processingValue {
					b.WriteString("Not ")
				}
				if b.Len() == 0 {
					b.WriteString("Processing ")
				} else {
					b.WriteString("processing ")
				}
			}
			if c.IsSet("expired") {
				expiredValue := c.Bool("expired")
				filter.Expired = &expiredValue
				if !expiredValue {
					if b.Len() == 0 {
						b.WriteString("Not ")
					} else {
						b.WriteString(", not ")
					}
				}
				if b.Len() == 0 {
					b.WriteString("Expired ")
				} else {
					b.WriteString("expired ")
				}
			}
			if b.Len() == 0 {
				b.WriteString("Sessions")
			} else {
				b.WriteString("sessions")
			}
			if c.IsSet("id") {
				idValue := c.String("id")
				filter.ID = &idValue
				b.WriteString(" with id '" + idValue + "'")
			}
			b.WriteString(":")
			uploads, err := managingFS.ListUploadSessions(c.Context, filter)
			if err != nil {
				return err
			}

			asJson := c.String("output") == "json"
			if !asJson {
				fmt.Println(b.String())
			}
			for _, u := range uploads {
				ref := u.Reference()
				if asJson {
					s := struct {
						ID         string         `json:"id"`
						Space      string         `json:"space"`
						Filename   string         `json:"filename"`
						Offset     int64          `json:"offset"`
						Size       int64          `json:"size"`
						Executant  userpb.UserId  `json:"executant"`
						SpaceOwner *userpb.UserId `json:"spaceowner,omitempty"`
						Expires    time.Time      `json:"expires"`
						Processing bool           `json:"processing"`
					}{
						Space:      ref.GetResourceId().GetSpaceId(),
						ID:         u.ID(),
						Filename:   u.Filename(),
						Offset:     u.Offset(),
						Size:       u.Size(),
						Executant:  u.Executant(),
						SpaceOwner: u.SpaceOwner(),
						Expires:    u.Expires(),
						Processing: u.IsProcessing(),
					}
					j, err := json.Marshal(s)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println(string(j))
				} else {
					fmt.Printf(" - %s (Space: %s, Name: %s, Size: %d/%d, Expires: %s, Processing: %t)\n", ref.GetResourceId().GetSpaceId(), u.ID(), u.Filename(), u.Offset(), u.Size(), u.Expires(), u.IsProcessing())
				}
			}
			return nil
		},
	}
}

// PurgeExpiredUploads is the entry point for the clean command
func PurgeExpiredUploads(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "clean",
		Usage: "Clean up leftovers from expired uploads",
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
			fs, err := f(drivers[cfg.Driver].(map[string]interface{}), nil)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize filesystem driver '%s'\n", cfg.Driver)
				return err
			}

			managingFS, ok := fs.(storage.UploadSessionLister)
			if !ok {
				fmt.Fprintf(os.Stderr, "'%s' storage does not support clean expired uploads\n", cfg.Driver)
				os.Exit(1)
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			processing := false
			expired := true
			uploads, err := managingFS.ListUploadSessions(c.Context, storage.UploadSessionFilter{Expired: &expired, Processing: &processing})
			if err != nil {
				return err
			}

			fmt.Println("purging uploads:")
			go func() {
				for _, u := range uploads {
					ref := u.Reference()
					fmt.Printf(" - %s (Space: %s, Name: %s, Size: %d/%d, Expires: %s, Processing: %t)\n", ref.GetResourceId().GetSpaceId(), u.ID(), u.Filename(), u.Offset(), u.Size(), u.Expires(), u.IsProcessing())
					u.Purge(c.Context)
				}
				wg.Done()
			}()

			wg.Wait()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to clean expired uploads '%s'\n", err)
				return err
			}
			return nil
		},
	}
}
