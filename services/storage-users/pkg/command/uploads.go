package command

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/urfave/cli/v2"

	"github.com/cs3org/reva/v2/pkg/storage"
	"github.com/cs3org/reva/v2/pkg/storage/fs/registry"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/config/parser"
	"github.com/owncloud/ocis/v2/services/storage-users/pkg/revaconfig"
)

func Uploads(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "uploads",
		Usage: "manage uploads",
		Before: func(c *cli.Context) error {
			if err := parser.ParseConfig(cfg); err != nil {
				fmt.Printf("%v", err)
				return err
			}
			return nil
		},
		Subcommands: []*cli.Command{
			ListUploads(cfg),
			PurgeExpiredUploads(cfg),
		},
	}
}

// ListUploads prints a list of all incomplete uploads
func ListUploads(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "list",
		Usage:    fmt.Sprintf("Print a list of all incomplete uploads"),
		Category: "services",
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			return err
		},
		Action: func(c *cli.Context) error {
			f, ok := registry.NewFuncs[cfg.Driver]
			if !ok {
				fmt.Fprintf(os.Stderr, "Unknown filesystem driver '%s'\n", cfg.Driver)
				os.Exit(1)
			}
			drivers := revaconfig.UserDrivers(cfg)
			fs, err := f(drivers[cfg.Driver].(map[string]interface{}))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize filesystem driver '%s'\n", cfg.Driver)
				return err
			}

			managingFS, ok := fs.(storage.UploadsManager)
			if !ok {
				fmt.Fprintf(os.Stderr, "'%s' storage does not support listing expired uploads\n", cfg.Driver)
				os.Exit(1)
			}

			uploads, err := managingFS.ListUploads()
			if err != nil {
				return err
			}

			fmt.Println("Incomplete uploads:")
			for _, u := range uploads {
				fmt.Printf(" - %s (%s, Size: %d, Expires: %s)\n", u.ID, u.MetaData["filename"], u.Size, expiredString(u.MetaData["expires"]))
			}
			return nil
		},
	}
}

// PurgeExpiredUploads is the entry point for the server command.
func PurgeExpiredUploads(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:     "purge",
		Usage:    fmt.Sprintf("Let %s extension clean up leftovers from expired downloads", cfg.Service.Name),
		Category: "services",
		Before: func(c *cli.Context) error {
			err := parser.ParseConfig(cfg)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			return err
		},
		Action: func(c *cli.Context) error {
			f, ok := registry.NewFuncs[cfg.Driver]
			if !ok {
				fmt.Fprintf(os.Stderr, "Unknown filesystem driver '%s'\n", cfg.Driver)
				os.Exit(1)
			}
			drivers := revaconfig.UserDrivers(cfg)
			fs, err := f(drivers[cfg.Driver].(map[string]interface{}))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to initialize filesystem driver '%s'\n", cfg.Driver)
				return err
			}

			managingFS, ok := fs.(storage.UploadsManager)
			if !ok {
				fmt.Fprintf(os.Stderr, "'%s' storage does not support purging expired uploads\n", cfg.Driver)
				os.Exit(1)
			}

			wg := sync.WaitGroup{}
			wg.Add(1)
			purgedChannel := make(chan tusd.FileInfo)

			go func() {
				for purged := range purgedChannel {
					fmt.Printf("Purging %s (Filename: %s, Size: %d, Expires: %s)\n",
						purged.ID, purged.MetaData["filename"], purged.Size, expiredString(purged.MetaData["expires"]))

				}
				wg.Done()
			}()

			err = managingFS.PurgeExpiredUploads(purgedChannel)
			close(purgedChannel)
			wg.Wait()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to purge expired uploads '%s'\n", err)
				return err
			}
			return nil
		},
	}
}

func expiredString(e string) string {
	expired := "N/A"
	iExpires, err := strconv.Atoi(e)
	if err == nil {
		expired = time.Unix(int64(iExpires), 0).Format(time.RFC3339)
	}
	return expired
}
