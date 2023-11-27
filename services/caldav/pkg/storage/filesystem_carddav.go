package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/carddav"
)

func (b *filesystemBackend) AddressbookHomeSetPath(ctx context.Context) (string, error) {
	upPath, err := b.CurrentUserPrincipal(ctx)
	if err != nil {
		return "", err
	}

	return path.Join(upPath, b.carddavPrefix) + "/", nil
}

func (b *filesystemBackend) localCardDAVPath(ctx context.Context, urlPath string) (string, error) {
	homeSetPath, err := b.AddressbookHomeSetPath(ctx)
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

func createDefaultAddressBook(path, localPath string) error {
	// TODO what should the default address book look like?
	defaultAB := carddav.AddressBook{
		Path:                 path,
		Name:                 "My contacts",
		Description:          "Default address book",
		MaxResourceSize:      1024,
		SupportedAddressData: nil,
	}
	blob, err := json.MarshalIndent(defaultAB, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating default address book: %s", err.Error())
	}
	err = os.WriteFile(localPath, blob, 0644)
	if err != nil {
		return fmt.Errorf("error writing default address book: %s", err.Error())
	}
	return nil
}

func (b *filesystemBackend) AddressBook(ctx context.Context) (*carddav.AddressBook, error) {
	log.Debug().Msg("filesystem.AddressBook()")
	localPath, err := b.localCardDAVPath(ctx, "")
	if err != nil {
		return nil, err
	}
	localPath = filepath.Join(localPath, "addressbook.json")

	log.Debug().Str("local_path", localPath).Msg("loading addressbook")

	data, readErr := ioutil.ReadFile(localPath)
	if os.IsNotExist(readErr) {
		urlPath, err := b.AddressbookHomeSetPath(ctx)
		if err != nil {
			return nil, err
		}
		urlPath = path.Join(urlPath, defaultResourceName) + "/"
		log.Debug().Str("local_path", localPath).Str("url_path", urlPath).Msg("creating addressbook")
		err = createDefaultAddressBook(urlPath, localPath)
		if err != nil {
			return nil, err
		}
		data, readErr = ioutil.ReadFile(localPath)
	}
	if readErr != nil {
		return nil, fmt.Errorf("error opening address book: %s", readErr.Error())
	}
	var addressBook carddav.AddressBook
	err = json.Unmarshal(data, &addressBook)
	if err != nil {
		return nil, fmt.Errorf("error reading address book: %s", err.Error())
	}

	return &addressBook, nil
}

func (b *filesystemBackend) GetAddressObject(ctx context.Context, objPath string, req *carddav.AddressDataRequest) (*carddav.AddressObject, error) {
	log.Debug().Str("url_path", objPath).Msg("filesystem.GetAddressObject()")
	localPath, err := b.localCardDAVPath(ctx, objPath)
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

func (b *filesystemBackend) loadAllContacts(ctx context.Context, propFilter []string) ([]carddav.AddressObject, error) {
	var result []carddav.AddressObject

	localPath, err := b.localCardDAVPath(ctx, "")
	if err != nil {
		return result, err
	}

	homeSetPath, err := b.AddressbookHomeSetPath(ctx)
	if err != nil {
		return result, err
	}

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

		obj := carddav.AddressObject{
			Path:          path.Join(homeSetPath, defaultResourceName, filepath.Base(filename)),
			ModTime:       info.ModTime(),
			ContentLength: info.Size(),
			ETag:          etag,
			Card:          card,
		}
		result = append(result, obj)
		return nil
	})

	log.Debug().Int("results", len(result)).Str("path", localPath).Msg("filesystem.loadAllContacts() successful")
	return result, err
}

func (b *filesystemBackend) ListAddressObjects(ctx context.Context, req *carddav.AddressDataRequest) ([]carddav.AddressObject, error) {
	log.Debug().Msg("filesystem.ListAddressObjects()")
	var propFilter []string
	if req != nil && !req.AllProp {
		propFilter = req.Props
	}

	return b.loadAllContacts(ctx, propFilter)
}

func (b *filesystemBackend) QueryAddressObjects(ctx context.Context, query *carddav.AddressBookQuery) ([]carddav.AddressObject, error) {
	log.Debug().Msg("filesystem.QueryAddressObjects()")
	var propFilter []string
	if query != nil && !query.DataRequest.AllProp {
		propFilter = query.DataRequest.Props
	}

	result, err := b.loadAllContacts(ctx, propFilter)
	if err != nil {
		return result, err
	}

	return carddav.Filter(query, result)
}

func (b *filesystemBackend) PutAddressObject(ctx context.Context, objPath string, card vcard.Card, opts *carddav.PutAddressObjectOptions) (loc string, err error) {
	log.Debug().Str("url_path", objPath).Msg("filesystem.PutAddressObject()")

	// Object always get saved as <UID>.vcf
	dirname, _ := path.Split(objPath)
	objPath = path.Join(dirname, card.Value(vcard.FieldUID)+".vcf")

	localPath, err := b.localCardDAVPath(ctx, objPath)
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
	log.Debug().Str("url_path", path).Msg("filesystem.DeleteAddressObject()")

	localPath, err := b.localCardDAVPath(ctx, path)
	if err != nil {
		return err
	}
	err = os.Remove(localPath)
	if err != nil {
		return err
	}
	return nil
}
