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

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

func (b *filesystemBackend) CalendarHomeSetPath(ctx context.Context) (string, error) {
	upPath, err := b.CurrentUserPrincipal(ctx)
	if err != nil {
		return "", err
	}

	return path.Join(upPath, b.caldavPrefix) + "/", nil
}

func (b *filesystemBackend) localCalDAVPath(ctx context.Context, urlPath string) (string, error) {
	homeSetPath, err := b.CalendarHomeSetPath(ctx)
	if err != nil {
		return "", err
	}

	return b.safeLocalPath(homeSetPath, urlPath)
}

func calendarFromFile(path string, propFilter []string) (*ical.Calendar, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	dec := ical.NewDecoder(f)
	cal, err := dec.Decode()
	if err != nil {
		return nil, err
	}

	return cal, nil
	// TODO implement
	//return icalPropFilter(cal, propFilter), nil
}

func (b *filesystemBackend) loadAllCalendars(ctx context.Context, propFilter []string) ([]caldav.CalendarObject, error) {
	var result []caldav.CalendarObject

	localPath, err := b.localCalDAVPath(ctx, "")
	if err != nil {
		return result, err
	}

	homeSetPath, err := b.CalendarHomeSetPath(ctx)
	if err != nil {
		return result, err
	}

	err = filepath.Walk(localPath, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %s", filename, err)
		}

		// Skip address book meta data files
		if !info.Mode().IsRegular() || filepath.Ext(filename) != ".ics" {
			return nil
		}

		cal, err := calendarFromFile(filename, propFilter)
		if err != nil {
			fmt.Printf("load calendar error for %s: %v\n", filename, err)
			return err
		}

		etag, err := etagForFile(filename)
		if err != nil {
			return err
		}

		obj := caldav.CalendarObject{
			Path:          path.Join(homeSetPath, defaultResourceName, filepath.Base(filename)),
			ModTime:       info.ModTime(),
			ContentLength: info.Size(),
			ETag:          etag,
			Data:          cal,
		}
		result = append(result, obj)
		return nil
	})

	log.Debug().Int("results", len(result)).Str("path", localPath).Msg("filesystem.loadAllCalendars() successful")
	return result, err
}

func createDefaultCalendar(path, localPath string) error {
	// TODO what should the default calendar look like?
	defaultC := caldav.Calendar{
		Path:            path,
		Name:            "My calendar",
		Description:     "Default calendar",
		MaxResourceSize: 4096,
	}
	blob, err := json.MarshalIndent(defaultC, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating default calendar: %s", err.Error())
	}
	err = os.WriteFile(localPath, blob, 0644)
	if err != nil {
		return fmt.Errorf("error writing default calendar: %s", err.Error())
	}
	return nil
}

func (b *filesystemBackend) Calendar(ctx context.Context) (*caldav.Calendar, error) {
	log.Debug().Msg("filesystem.Calendar()")

	localPath, err := b.localCalDAVPath(ctx, "")
	if err != nil {
		return nil, err
	}
	localPath = filepath.Join(localPath, "calendar.json")

	log.Debug().Str("local_path", localPath).Msg("loading calendar")

	data, readErr := ioutil.ReadFile(localPath)
	if os.IsNotExist(readErr) {
		urlPath, err := b.CalendarHomeSetPath(ctx)
		if err != nil {
			return nil, err
		}
		urlPath = path.Join(urlPath, defaultResourceName) + "/"
		log.Debug().Str("local_path", localPath).Str("url_path", urlPath).Msg("creating calendar")
		err = createDefaultCalendar(urlPath, localPath)
		if err != nil {
			return nil, err
		}
		data, readErr = ioutil.ReadFile(localPath)
	}
	if readErr != nil {
		return nil, fmt.Errorf("error opening calendar: %s", readErr.Error())
	}
	var calendar caldav.Calendar
	err = json.Unmarshal(data, &calendar)
	if err != nil {
		return nil, fmt.Errorf("error reading calendar: %s", err.Error())
	}

	return &calendar, nil
}

func (b *filesystemBackend) GetCalendarObject(ctx context.Context, objPath string, req *caldav.CalendarCompRequest) (*caldav.CalendarObject, error) {
	log.Debug().Str("url_path", objPath).Msg("filesystem.GetCalendarObject()")

	localPath, err := b.localCalDAVPath(ctx, objPath)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(localPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			log.Debug().Str("local_path", localPath).Msg("object not found")
			return nil, webdav.NewHTTPError(404, err)
		}
		return nil, err
	}

	var propFilter []string
	if req != nil && !req.AllProps {
		propFilter = req.Props
	}

	calendar, err := calendarFromFile(localPath, propFilter)
	if err != nil {
		log.Debug().Err(err).Msg("error reading calendar")
		return nil, err
	}

	etag, err := etagForFile(localPath)
	if err != nil {
		return nil, err
	}

	obj := caldav.CalendarObject{
		Path:          objPath,
		ModTime:       info.ModTime(),
		ContentLength: info.Size(),
		ETag:          etag,
		Data:          calendar,
	}
	return &obj, nil
}

func (b *filesystemBackend) ListCalendarObjects(ctx context.Context, req *caldav.CalendarCompRequest) ([]caldav.CalendarObject, error) {
	log.Debug().Msg("filesystem.ListCalendarObjects()")

	var propFilter []string
	if req != nil && !req.AllProps {
		propFilter = req.Props
	}

	return b.loadAllCalendars(ctx, propFilter)
}

func (b *filesystemBackend) QueryCalendarObjects(ctx context.Context, query *caldav.CalendarQuery) ([]caldav.CalendarObject, error) {
	log.Debug().Msg("filesystem.QueryCalendarObjects()")

	var propFilter []string
	if query != nil && !query.CompRequest.AllProps {
		propFilter = query.CompRequest.Props
	}

	result, err := b.loadAllCalendars(ctx, propFilter)
	if err != nil {
		return result, err
	}

	return caldav.Filter(query, result)
}

func (b *filesystemBackend) PutCalendarObject(ctx context.Context, objPath string, calendar *ical.Calendar, opts *caldav.PutCalendarObjectOptions) (loc string, err error) {
	log.Debug().Str("url_path", objPath).Msg("filesystem.PutCalendarObject()")

	_, uid, err := caldav.ValidateCalendarObject(calendar)
	if err != nil {
		return "", caldav.NewPreconditionError(caldav.PreconditionValidCalendarObjectResource)
	}

	// Object always get saved as <UID>.ics
	dirname, _ := path.Split(objPath)
	objPath = path.Join(dirname, uid+".ics")

	localPath, err := b.localCalDAVPath(ctx, objPath)
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
		return "", caldav.NewPreconditionError(caldav.PreconditionNoUIDConflict)
	} else if err != nil {
		return "", err
	}
	defer f.Close()

	enc := ical.NewEncoder(f)
	err = enc.Encode(calendar)
	if err != nil {
		return "", err
	}

	return objPath, nil
}

func (b *filesystemBackend) DeleteCalendarObject(ctx context.Context, path string) error {
	log.Debug().Str("url_path", path).Msg("filesystem.DeleteCalendarObject()")

	localPath, err := b.localCalDAVPath(ctx, path)
	if err != nil {
		return err
	}
	err = os.Remove(localPath)
	if err != nil {
		return err
	}
	return nil
}
