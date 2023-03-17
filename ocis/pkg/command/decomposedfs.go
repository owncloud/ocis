package command

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/lookup"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/metadata"
	"github.com/cs3org/reva/v2/pkg/storage/utils/decomposedfs/options"
	"github.com/cs3org/reva/v2/pkg/storagespace"
	"github.com/owncloud/ocis/v2/ocis-pkg/config"
	"github.com/owncloud/ocis/v2/ocis/pkg/register"
	"github.com/urfave/cli/v2"
)

// DecomposedfsCommand is the entrypoint for the groups command.
func DecomposedfsCommand(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "decomposedfs",
		Usage:       `cli tools to inspect and manipulate a decomposedfs storage.`,
		Category:    "maintenance",
		Subcommands: []*cli.Command{metadataCmd(cfg)},
	}
}

func init() {
	register.AddCommand(DecomposedfsCommand)
}

func metadataCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "metadata",
		Usage: `cli tools to inspect and manipulate node metadata`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "root",
				Aliases:  []string{"r"},
				Required: true,
				Usage:    "Path to the decomposedfs",
			},
			&cli.StringFlag{
				Name:     "node",
				Required: true,
				Aliases:  []string{"n"},
				Usage:    "Path to or ID of the node to inspect",
			},
		},
		Subcommands: []*cli.Command{dumpCmd(cfg), getCmd(cfg), setCmd(cfg)},
	}
}

func dumpCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "dump",
		Usage: `print the metadata of the given node. String attributes will be enclosed in quotes. Binary attributes will be returned encoded as base64 with their value being prefixed with '0s'.`,
		Action: func(c *cli.Context) error {
			lu, backend := getBackend(c)
			path, err := getPath(c, lu)
			if err != nil {
				return err
			}

			attribs, err := backend.All(path)
			if err != nil {
				fmt.Println("Error reading attributes")
				return err
			}
			printAttribs(attribs, c.String("attribute"))
			return nil
		},
	}
}

func getCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: `print a specific attribute of the given node. String attributes will be enclosed in quotes. Binary attributes will be returned encoded as base64 with their value being prefixed with '0s'.`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "attribute",
				Aliases: []string{"a"},
				Usage:   "attribute to inspect",
			},
		},
		Action: func(c *cli.Context) error {
			lu, backend := getBackend(c)
			path, err := getPath(c, lu)
			if err != nil {
				return err
			}

			attribs, err := backend.All(path)
			if err != nil {
				fmt.Println("Error reading attributes")
				return err
			}
			printAttribs(attribs, c.String("attribute"))
			return nil
		},
	}
}

func setCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: `manipulate metadata of the given node. Binary attributes can be given hex encoded (prefix by '0x') or base64 encoded (prefix by '0s').`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "attribute",
				Required: true,
				Aliases:  []string{"a"},
				Usage:    "attribute to inspect",
			},
			&cli.StringFlag{
				Name:     "value",
				Required: true,
				Aliases:  []string{"v"},
				Usage:    "value to set",
			},
		},
		Action: func(c *cli.Context) error {
			lu, backend := getBackend(c)
			path, err := getPath(c, lu)
			if err != nil {
				return err
			}

			v := c.String("value")
			if strings.HasPrefix(v, "0s") {
				b64, err := base64.StdEncoding.DecodeString(v[2:])
				if err == nil {
					v = string(b64)
				}
			} else if strings.HasPrefix(v, "0x") {
				h, err := hex.DecodeString(v)
				if err == nil {
					v = string(h)
				}
			}

			err = backend.Set(path, c.String("attribute"), []byte(v[2:]))
			if err != nil {
				fmt.Println("Error setting attribute")
				return err
			}
			return nil
		},
	}
}

func backend(root, backend string) metadata.Backend {
	switch backend {
	case "xattrs":
		return metadata.XattrsBackend{}
	case "mpk":
		return metadata.NewMessagePackBackend(root, options.CacheOptions{})
	}
	return metadata.NullBackend{}
}

func getBackend(c *cli.Context) (*lookup.Lookup, metadata.Backend) {
	rootFlag := c.String("root")

	bod := lookup.DetectBackendOnDisk(rootFlag)
	backend := backend(rootFlag, bod)
	lu := lookup.New(backend, &options.Options{
		Root:            rootFlag,
		MetadataBackend: bod,
	})
	return lu, backend
}

func getPath(c *cli.Context, lu *lookup.Lookup) (string, error) {
	nodeFlag := c.String("node")

	path := ""
	if strings.HasPrefix(nodeFlag, "/") {
		path = nodeFlag
	} else {
		nId := c.String("node")
		id, err := storagespace.ParseID(nId)
		if err != nil {
			fmt.Println("Invalid node id.")
			return "", err
		}
		n, _ := lu.NodeFromID(context.Background(), &id)
		if err != nil || !n.Exists {
			fmt.Println("Can not find node '" + nId + "'")
			return "", err
		}
		path = n.InternalPath()
	}
	return path, nil
}

func printAttribs(attribs map[string][]byte, onlyAttribute string) {
	if onlyAttribute != "" {
		fmt.Println(onlyAttribute + `=` + attribToString(attribs[onlyAttribute]))
		return
	}

	names := []string{}
	for k, _ := range attribs {
		names = append(names, k)
	}

	sort.Strings(names)

	for _, n := range names {
		fmt.Println(n + `=` + attribToString(attribs[n]))
	}
}

func attribToString(attrib []byte) string {
	for i := 0; i < len(attrib); i++ {
		if attrib[i] < 32 || attrib[i] >= 127 {
			return "0s" + base64.StdEncoding.EncodeToString(attrib)
		}
	}
	return `"` + string(attrib) + `"`
}
