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

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav"
	"github.com/emersion/go-webdav/caldav"
)

const calendarFileName = "calendar.json"

func (b *filesystemBackend) CalendarHomeSetPath(ctx context.Context) (string, error) {
	user, ok := revaContext.ContextGetUser(ctx)
	if !ok {
		return "", errors.New("no user in context")
	}
	return fmt.Sprintf("/dav/calendars/%s/", user.Username), nil
}

func (b *filesystemBackend) CreateCalendar(ctx context.Context, calendar *caldav.Calendar) error {
	// TODO what should the default calendar look like?
	resourceName := path.Base(calendar.Path)
	localPath, err_ := b.localCalDAVDir(ctx, resourceName)
	if err_ != nil {
		return fmt.Errorf("error creating default calendar: %s", err_.Error())
	}

	log.Debug().Str("local", localPath).Str("url", calendar.Path).Msg("filesystem.CreateCalendar()")

	blob, err := json.MarshalIndent(calendar, "", "  ")
	if err != nil {
		return fmt.Errorf("error creating default calendar: %s", err.Error())
	}
	err = os.WriteFile(path.Join(localPath, calendarFileName), blob, 0644)
	if err != nil {
		return fmt.Errorf("error writing default calendar: %s", err.Error())
	}
	return nil
}

func (b *filesystemBackend) localCalDAVDir(ctx context.Context, components ...string) (string, error) {
	homeSetPath, err := b.CalendarHomeSetPath(ctx)
	if err != nil {
		return "", err
	}

	return b.localDir(homeSetPath, components...)
}

func (b *filesystemBackend) safeLocalCalDAVPath(ctx context.Context, urlPath string) (string, error) {
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

func (b *filesystemBackend) loadAllCalendarObjects(ctx context.Context, urlPath string, propFilter []string) ([]caldav.CalendarObject, error) {
	var result []caldav.CalendarObject

	localPath, err := b.safeLocalCalDAVPath(ctx, urlPath)
	if err != nil {
		return result, err
	}

	log.Debug().Str("path", localPath).Msg("loading calendar objects")

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

		// TODO can this potentially be called on a calendar object resource?
		// Would work (as Walk() includes root), except for the path construction below
		obj := caldav.CalendarObject{
			Path:          path.Join(urlPath, filepath.Base(filename)),
			ModTime:       info.ModTime(),
			ContentLength: info.Size(),
			ETag:          etag,
			Data:          cal,
		}
		result = append(result, obj)
		return nil
	})

	return result, err
}

func (b *filesystemBackend) createDefaultCalendar(ctx context.Context) (*caldav.Calendar, error) {
	homeSetPath, err_ := b.CalendarHomeSetPath(ctx)
	if err_ != nil {
		return nil, fmt.Errorf("error creating default calendar: %s", err_.Error())
	}
	urlPath := path.Join(homeSetPath, defaultResourceName) + "/"

	log.Debug().Str("url", urlPath).Msg("filesystem.CreateCalendar()")

	defaultC := caldav.Calendar{
		Path:            urlPath,
		Name:            "My calendar",
		Description:     "Default calendar",
		MaxResourceSize: 4096,
	}
	err := b.CreateCalendar(ctx, &defaultC)
	if err != nil {
		return nil, err
	}

	return &defaultC, nil
}

func (b *filesystemBackend) ListCalendars(ctx context.Context) ([]caldav.Calendar, error) {
	log.Debug().Msg("filesystem.ListCalendars()")

	localPath, err := b.localCalDAVDir(ctx)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("path", localPath).Msg("looking for calendars")

	var result []caldav.Calendar

	err = filepath.Walk(localPath, func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing %s: %s", filename, err.Error())
		}

		if !info.IsDir() || filename == localPath {
			return nil
		}

		calPath := path.Join(filename, calendarFileName)
		data, err := os.ReadFile(calPath)
		if err != nil {
			if os.IsNotExist(err) {
				return nil // not a calendar dir
			} else {
				return fmt.Errorf("error accessing %s: %s", calPath, err.Error())
			}
		}

		var calendar caldav.Calendar
		err = json.Unmarshal(data, &calendar)
		if err != nil {
			return fmt.Errorf("error reading calendar %s: %s", calPath, err.Error())
		}

		result = append(result, calendar)
		return nil
	})

	if err == nil && len(result) == 0 {
		// Nothing here yet? Create the default calendar.
		log.Debug().Msg("no calendars found, creating default calendar")
		cal, err := b.createDefaultCalendar(ctx)
		if err == nil {
			result = append(result, *cal)
		}
	}
	log.Debug().Int("results", len(result)).Err(err).Msg("filesystem.ListCalendars() done")
	return result, err
}

func (b *filesystemBackend) GetCalendar(ctx context.Context, urlPath string) (*caldav.Calendar, error) {
	log.Debug().Str("path", urlPath).Msg("filesystem.GetCalendar()")

	localPath, err := b.safeLocalCalDAVPath(ctx, urlPath)
	if err != nil {
		return nil, err
	}
	localPath = filepath.Join(localPath, calendarFileName)

	log.Debug().Str("path", localPath).Msg("loading calendar")

	data, err := os.ReadFile(localPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, webdav.NewHTTPError(404, err)
		}
		return nil, fmt.Errorf("error opening calendar: %s", err.Error())
	}
	var calendar caldav.Calendar
	err = json.Unmarshal(data, &calendar)
	if err != nil {
		return nil, fmt.Errorf("error reading calendar: %s", err.Error())
	}

	return &calendar, nil
}

func (b *filesystemBackend) GetCalendarObject(ctx context.Context, objPath string, req *caldav.CalendarCompRequest) (*caldav.CalendarObject, error) {
	log.Debug().Str("path", objPath).Msg("filesystem.GetCalendarObject()")

	localPath, err := b.safeLocalCalDAVPath(ctx, objPath)
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
	if req != nil && !req.AllProps {
		propFilter = req.Props
	}

	calendar, err := calendarFromFile(localPath, propFilter)
	if err != nil {
		log.Debug().Str("path", localPath).Err(err).Msg("error reading calendar")
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

func (b *filesystemBackend) ListCalendarObjects(ctx context.Context, urlPath string, req *caldav.CalendarCompRequest) ([]caldav.CalendarObject, error) {
	log.Debug().Str("path", urlPath).Msg("filesystem.ListCalendarObjects()")

	var propFilter []string
	if req != nil && !req.AllProps {
		propFilter = req.Props
	}

	result, err := b.loadAllCalendarObjects(ctx, urlPath, propFilter)
	log.Debug().Int("results", len(result)).Err(err).Msg("filesystem.ListCalendarObjects() done")
	return result, err
}

func (b *filesystemBackend) QueryCalendarObjects(ctx context.Context, urlPath string, query *caldav.CalendarQuery) ([]caldav.CalendarObject, error) {
	log.Debug().Str("path", urlPath).Msg("filesystem.QueryCalendarObjects()")

	var propFilter []string
	if query != nil && !query.CompRequest.AllProps {
		propFilter = query.CompRequest.Props
	}

	result, err := b.loadAllCalendarObjects(ctx, urlPath, propFilter)
	log.Debug().Int("results", len(result)).Err(err).Msg("filesystem.QueryCalendarObjects() load done")
	if err != nil {
		return result, err
	}

	filtered, err := caldav.Filter(query, result)
	log.Debug().Int("results", len(filtered)).Err(err).Msg("filesystem.QueryCalendarObjects() filter done")
	return filtered, err
}

func (b *filesystemBackend) PutCalendarObject(ctx context.Context, objPath string, calendar *ical.Calendar, opts *caldav.PutCalendarObjectOptions) (loc string, err error) {
	log.Debug().Str("path", objPath).Msg("filesystem.PutCalendarObject()")

	_, _, err = caldav.ValidateCalendarObject(calendar)
	if err != nil {
		return "", caldav.NewPreconditionError(caldav.PreconditionValidCalendarObjectResource)
	}

	localPath, err := b.safeLocalCalDAVPath(ctx, objPath)
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
	log.Debug().Str("path", path).Msg("filesystem.DeleteCalendarObject()")

	localPath, err := b.safeLocalCalDAVPath(ctx, path)
	if err != nil {
		return err
	}
	err = os.Remove(localPath)
	if err != nil {
		return err
	}
	return nil
}
