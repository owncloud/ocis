package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	revaContext "github.com/cs3org/reva/v2/pkg/ctx"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/carddav"
)

const addressBookFileName = "addressbook.json"

func (b *filesystemBackend) AddressBookHomeSetPath(ctx context.Context) (string, error) {
	user, ok := revaContext.ContextGetUser(ctx)
	if !ok {
		return "", errors.New("no user in context")
	}
	return fmt.Sprintf("/dav/contacts/%s/", user.Username), nil
}

func (b *filesystemBackend) CreateAddressBook(ctx context.Context, addressBook *carddav.AddressBook) error {
	//TODO implement me
	panic("implement me")
}

func (b *filesystemBackend) DeleteAddressBook(ctx context.Context, path string) error {
	//TODO implement me
	panic("implement me")
}

func (b *filesystemBackend) localCardDAVDir(ctx context.Context, components ...string) (string, error) {
	homeSetPath, err := b.AddressBookHomeSetPath(ctx)
	if err != nil {
		return "", err
	}

	return b.localDir(homeSetPath, components...)
}

func (b *filesystemBackend) safeLocalCardDAVPath(ctx context.Context, urlPath string) (string, error) {
	homeSetPath, err := b.AddressBookHomeSetPath(ctx)
	if err != nil {
		return "", err
	}

	return b.safeLocalPath(homeSetPath, urlPath)
}

func vcardPropFilter(card vcard.Card, props []string) vcard.Card {
	if card == nil {
		return nil
	}

	if len(props) == 0 {
		return card
	}

	result := make(vcard.Card)
	result["VERSION"] = card["VERSION"]
	for _, prop := range props {
		value, ok := card[prop]
		if ok {
			result[prop] = value
		}
	}

	return result
}

func vcardFromFile(path string, propFilter []string) (vcard.Card, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := vcard.NewDecoder(f)
	card, err := dec.Decode()
	if err != nil {
		return nil, err
	}

	return vcardPropFilter(card, propFilter), nil
}

func (b *filesystemBackend) loadAllAddressObjects(ctx context.Context, urlPath string, propFilter []string) ([]carddav.AddressObject, error) {
	var result []carddav.AddressObject

	localPath, err := b.safeLocalCardDAVPath(ctx, urlPath)
	if err != nil {
		return result, err
	}

	log.Debug().Str("path", localPath).Msg("loading address objects")

	err = filepath.Walk(localPath, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %s", filename, err)
		}

		if !info.Mode().IsRegular() || filepath.Ext(filename) != ".vcf" {
			return nil
		}

		card, err := vcardFromFile(filename, propFilter)
		if err != nil {
			return err
		}

		etag, err := etagForFile(filename)
		if err != nil {
			return err
		}

		// TODO can this potentially be called on an address object resource?
		// would work (as Walk() includes root), except for the path construction below
		obj := carddav.AddressObject{
			Path:          path.Join(urlPath, filepath.Base(filename)),
			ModTime:       info.ModTime(),
			ContentLength: info.Size(),
			ETag:          etag,
			Card:          card,
		}
		result = append(result, obj)
		return nil
	})

	return result, err
}

func (b *filesystemBackend) createDefaultAddressBook(ctx context.Context) (*carddav.AddressBook, error) {
	// TODO what should the default address book look like?
	localPath, err_ := b.localCardDAVDir(ctx, defaultResourceName)
	if err_ != nil {
		return nil, fmt.Errorf("error creating default address book: %s", err_.Error())
	}

	homeSetPath, err_ := b.AddressBookHomeSetPath(ctx)
	if err_ != nil {
		return nil, fmt.Errorf("error creating default address book: %s", err_.Error())
	}

	urlPath := path.Join(homeSetPath, defaultResourceName) + "/"

	log.Debug().Str("local", localPath).Str("url", urlPath).Msg("filesystem.createDefaultAddressBook()")

	defaultAB := carddav.AddressBook{
		Path:                 urlPath,
		Name:                 "My contacts",
		Description:          "Default address book",
		MaxResourceSize:      1024,
		SupportedAddressData: nil,
	}
	blob, err := json.MarshalIndent(defaultAB, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error creating default address book: %s", err.Error())
	}
	err = os.WriteFile(path.Join(localPath, addressBookFileName), blob, 0644)
	if err != nil {
		return nil, fmt.Errorf("error writing default address book: %s", err.Error())
	}
	return &defaultAB, nil
}

func (b *filesystemBackend) ListAddressBooks(ctx context.Context) ([]carddav.AddressBook, error) {
	log.Debug().Msg("filesystem.ListAddressBooks()")

	localPath, err := b.localCardDAVDir(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("path", localPath).Msg("looking for address books")

	var result []carddav.AddressBook

	err = filepath.Walk(localPath, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %s", filename, err.Error())
		}

		if !info.IsDir() || filename == localPath {
			return nil
		}

		abPath := path.Join(filename, addressBookFileName)
		data, err := os.ReadFile(abPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil // not an address book dir
			} else {
				return fmt.Errorf("error accessing %s: %s", abPath, err.Error())
			}
		}

		var addressBook carddav.AddressBook
		err = json.Unmarshal(data, &addressBook)
		if err != nil {
			return fmt.Errorf("error reading address book %s: %s", abPath, err.Error())
		}

		result = append(result, addressBook)
		return nil
	})

	if err == nil && len(result) == 0 {
		// Nothing here yet? Create the default address book
		log.Debug().Msg("no address books found, creating default address book")
		ab, err := b.createDefaultAddressBook(ctx)
		if err == nil {
			result = append(result, *ab)
		}
	}
	log.Debug().Int("results", len(result)).Err(err).Msg("filesystem.ListAddressBooks() done")
	return result, err
}

func (b *filesystemBackend) GetAddressBook(ctx context.Context, urlPath string) (*carddav.AddressBook, error) {
	log.Debug().Str("path", urlPath).Msg("filesystem.AddressBook()")

	localPath, err := b.safeLocalCardDAVPath(ctx, urlPath)
	if err != nil {
		return nil, err
	}
	localPath = filepath.Join(localPath, addressBookFileName)

	log.Debug().Str("path", localPath).Msg("loading addressbook")

	data, err := os.ReadFile(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, webdav.NewHTTPError(404, err)
		}
		return nil, fmt.Errorf("error opening address book: %s", err.Error())
	}
	var addressBook carddav.AddressBook
	err = json.Unmarshal(data, &addressBook)
	if err != nil {
		return nil, fmt.Errorf("error reading address book: %s", err.Error())
	}

	return &addressBook, nil
}

func (b *filesystemBackend) GetAddressObject(ctx context.Context, objPath string, req *carddav.AddressDataRequest) (*carddav.AddressObject, error) {
	log.Debug().Str("path", objPath).Msg("filesystem.GetAddressObject()")

	localPath, err := b.safeLocalCardDAVPath(ctx, objPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, webdav.NewHTTPError(404, err)
		}
		return nil, err
	}

	var propFilter []string
	if req != nil && !req.AllProp {
		propFilter = req.Props
	}

	card, err := vcardFromFile(localPath, propFilter)
	if err != nil {
		log.Debug().Str("path", localPath).Err(err).Msg("error reading calendar")
		return nil, err
	}

	etag, err := etagForFile(localPath)
	if err != nil {
		return nil, err
	}

	obj := carddav.AddressObject{
		Path:          objPath,
		ModTime:       info.ModTime(),
		ContentLength: info.Size(),
		ETag:          etag,
		Card:          card,
	}
	return &obj, nil
}

func (b *filesystemBackend) ListAddressObjects(ctx context.Context, urlPath string, req *carddav.AddressDataRequest) ([]carddav.AddressObject, error) {
	log.Debug().Str("path", urlPath).Msg("filesystem.ListAddressObjects()")

	var propFilter []string
	if req != nil && !req.AllProp {
		propFilter = req.Props
	}

	result, err := b.loadAllAddressObjects(ctx, urlPath, propFilter)
	log.Debug().Int("results", len(result)).Err(err).Msg("filesystem.ListAddressObjects() done")
	return result, err
}

func (b *filesystemBackend) QueryAddressObjects(ctx context.Context, urlPath string, query *carddav.AddressBookQuery) ([]carddav.AddressObject, error) {
	log.Debug().Str("path", urlPath).Msg("filesystem.QueryAddressObjects()")

	var propFilter []string
	if query != nil && !query.DataRequest.AllProp {
		propFilter = query.DataRequest.Props
	}

	result, err := b.loadAllAddressObjects(ctx, urlPath, propFilter)
	log.Debug().Int("results", len(result)).Err(err).Msg("filesystem.QueryAddressObjects() load done")
	if err != nil {
		return result, err
	}

	filtered, err := carddav.Filter(query, result)
	log.Debug().Int("results", len(filtered)).Err(err).Msg("filesystem.QueryAddressObjects() filter done")
	return filtered, err
}

func (b *filesystemBackend) PutAddressObject(ctx context.Context, objPath string, card vcard.Card, opts *carddav.PutAddressObjectOptions) (loc string, err error) {
	log.Debug().Str("path", objPath).Msg("filesystem.PutAddressObject()")

	// Object always get saved as <UID>.vcf
	dirname, _ := path.Split(objPath)
	objPath = path.Join(dirname, card.Value(vcard.FieldUID)+".vcf")

	localPath, err := b.safeLocalCardDAVPath(ctx, objPath)
	if err != nil {
		return "", err
	}

	flags := os.O_RDWR | os.O_CREATE | os.O_TRUNC
	// TODO handle IfNoneMatch == ETag
	if opts.IfNoneMatch.IsWildcard() {
		// Make sure we're not overwriting an existing file
		flags |= os.O_EXCL
	} else if opts.IfMatch.IsWildcard() {
		// Make sure we _are_ overwriting an existing file
		flags &= ^os.O_CREATE
	} else if opts.IfMatch.IsSet() {
		// Make sure we overwrite the _right_ file
		etag, err := etagForFile(localPath)
		if err != nil {
			return "", webdav.NewHTTPError(http.StatusPreconditionFailed, err)
		}
		want, err := opts.IfMatch.ETag()
		if err != nil {
			return "", webdav.NewHTTPError(http.StatusBadRequest, err)
		}
		if want != etag {
			err = fmt.Errorf("If-Match does not match current ETag (%s/%s)", want, etag)
			return "", webdav.NewHTTPError(http.StatusPreconditionFailed, err)
		}
	}

	f, err := os.OpenFile(localPath, flags, 0666)
	if os.IsExist(err) {
		return "", carddav.NewPreconditionError(carddav.PreconditionNoUIDConflict)
	} else if err != nil {
		return "", err
	}
	defer f.Close()

	enc := vcard.NewEncoder(f)
	err = enc.Encode(card)
	if err != nil {
		return "", err
	}

	return objPath, nil
}

func (b *filesystemBackend) DeleteAddressObject(ctx context.Context, path string) error {
	log.Debug().Str("path", path).Msg("filesystem.DeleteAddressObject()")

	localPath, err := b.safeLocalCardDAVPath(ctx, path)
	if err != nil {
		return err
	}
	err = os.Remove(localPath)
	if err != nil {
		return err
	}
	return nil
}
