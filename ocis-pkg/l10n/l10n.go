// package l10n holds translation mechanics that are used by user facing services (notifications, userlog, graph)
package l10n

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"reflect"

	"github.com/leonelquinteros/gotext"
	"github.com/owncloud/ocis/v2/ocis-pkg/middleware"
	settingssvc "github.com/owncloud/ocis/v2/protogen/gen/ocis/services/settings/v0"
	"github.com/owncloud/ocis/v2/services/settings/pkg/store/defaults"
	micrometadata "go-micro.dev/v4/metadata"
)

// HeaderAcceptLanguage is the header key for the accept-language header
var HeaderAcceptLanguage = "Accept-Language"

// Template marks a string as translatable
func Template(s string) string { return s }

// Translator is able to translate strings
type Translator struct {
	fs            fs.FS
	defaultLocale string
	domain        string
}

// NewTranslator creates a Translator with library path and language code and load default domain
func NewTranslator(defaultLocale string, domain string, fsys fs.FS) Translator {
	return Translator{
		fs:            fsys,
		defaultLocale: defaultLocale,
		domain:        domain,
	}
}

// NewTranslatorFromCommonConfig creates a new Translator from legacy config
func NewTranslatorFromCommonConfig(defaultLocale string, domain string, path string, fsys fs.FS, fsSubPath string) Translator {
	var filesystem fs.FS
	if path == "" {
		filesystem, _ = fs.Sub(fsys, fsSubPath)
	} else { // use custom path instead
		filesystem = os.DirFS(path)
	}
	return NewTranslator(defaultLocale, domain, filesystem)
}

// Translate translates a string to the locale
func (t Translator) Translate(str, locale string) string {
	return t.Locale(locale).Get(str)
}

func TranslateLocation(t Translator, locale string) func(string, ...any) string {
	return t.Locale(locale).Get
}

type structs func() []any
type maps func() []any
type each func() []any
type field func() string

func TranslateField(fieldName string) field {
	return func() string {
		return fieldName
	}
}

func TranslateStruct(args ...any) structs {
	return func() []any {
		return args
	}
}

func TranslateMap(args ...any) maps {
	return func() []any {
		return args
	}
}

func TranslateEach(args ...any) each {
	return func() []any {
		return args
	}
}

// TranslateEntity function provides the generic way to translate the necessary fields in composite entities.
// The function takes TranslateLocation function and entity with fields to translate.
// The function supports nested structs and slices of structs.
// Depending on entity type it should be wrapped to the appropriate function TranslateStruct, TranslateEach, TranslateField.
//
//		type InnerStruct struct {
//			Description string
//			DisplayName *string
//		}
//
//		type WrapperStruct struct {
//	     StructList []*InnerStruct
//		}
//		s:= &WrapperStruct{
//			StructList: []*InnerStruct{
//					{
//						Description: "innerDescription 1",
//						DisplayName: toStrPointer("innerDisplayName 1"),
//					},
//					{
//						Description: "innerDescription 2",
//						DisplayName: toStrPointer("innerDisplayName 2"),
//					},
//				},
//			}
//		tr := l10n_pkg.NewTranslateLocation(loc, "en")
//		err := l10n.TranslateEntity(tr,
//			l10n.TranslateStruct(s,
//				l10n.TranslateEach("StructList",
//					l10n.TranslateField("Description"),
//					l10n.TranslateField("DisplayName"))),
//		)
func TranslateEntity(tr func(string, ...any) string, arg any) error {
	switch a := arg.(type) {
	case structs:
		args := a()
		if len(args) < 2 {
			return fmt.Errorf("the TranslateStruct function expects at least 2 arguments, sructure and fields to translate")
		}
		entity := args[0]
		value := reflect.ValueOf(entity)
		if !isStruct(value) {
			return fmt.Errorf("the root entity must be a struct, got %v", value.Kind())
		}
		rangeOverArgs(tr, value, args[1:]...)
		return nil
	case each:
		args := a()
		if len(args) < 1 {
			return fmt.Errorf("the translateEach function expects at least 1 argument, slice")
		}
		entity := args[0]
		value := reflect.ValueOf(entity)
		if len(args) > 1 {
			translateEach(tr, value, args[1:]...)
		} else {
			translateEach(tr, value)
		}
		return nil
	case maps:
		// TODO implement
	}
	return ErrUnsupportedType
}

func translateEach(tr func(string, ...any) string, value reflect.Value, args ...any) {
	// Indirect through pointers and interfaces
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return
	}
	switch value.Kind() {
	case reflect.Array, reflect.Slice:
		for i := 0; i < value.Len(); i++ {
			v := value.Index(i)
			if args != nil {
				rangeOverArgs(tr, v, args...)
				continue
			}
			translateField(tr, v)
		}
	case reflect.Map:
		for _, k := range value.MapKeys() {
			rangeOverArgs(tr, value.MapIndex(k), args...)
		}
	}
}

func rangeOverArgs(tr func(string, ...any) string, value reflect.Value, args ...any) {
	// Indirect through pointers and interfaces
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return
	}
	for _, arg := range args {
		switch a := arg.(type) {
		case field:
			fieldName := a()
			// exported field
			f := value.FieldByName(fieldName)
			translateField(tr, f)
		case structs:
			args := a()
			if len(args) > 2 {
				if fieldName, ok := args[0].(string); ok {
					// exported field
					innerValue := value.FieldByName(fieldName)
					if !innerValue.IsValid() {
						return
					}
					if isStruct(innerValue) {
						rangeOverArgs(tr, innerValue, args[1:]...)
						return
					}
				}
			}
		case each:
			args := a()
			if len(args) > 2 {
				if fieldName, ok := args[0].(string); ok {
					// exported field
					innerValue := value.FieldByName(fieldName)
					if !innerValue.IsValid() {
						return
					}
					switch innerValue.Kind() {
					case reflect.Array, reflect.Slice:
						translateEach(tr, innerValue, args[1:]...)
					}
				}
			}
		}
	}
}

func translateField(tr func(string, ...any) string, f reflect.Value) {
	if f.IsValid() {
		if f.Kind() == reflect.Ptr {
			if f.IsNil() {
				return
			}
			f = f.Elem()
		}
		// A Value can be changed only if it is
		// addressable and was not obtained by
		// the use of unexported struct fields.
		if f.CanSet() {
			// change value
			if f.Kind() == reflect.String {
				val := tr(f.String())
				if val == "" {
					return
				}
				f.SetString(val)
			}
		}
	}
}

func isStruct(r reflect.Value) bool {
	if r.Kind() == reflect.Ptr {
		r = r.Elem()
	}
	return r.Kind() == reflect.Struct
}

var (
	ErrUnsupportedType = errors.New("unsupported type")
)

// Locale returns the gotext.Locale, use `.Get` method to translate strings
func (t Translator) Locale(locale string) *gotext.Locale {
	l := gotext.NewLocaleFS(locale, t.fs)
	l.AddDomain(t.domain) // make domain configurable only if needed
	if locale != "en" && len(l.GetTranslations()) == 0 {
		l = gotext.NewLocaleFS(t.defaultLocale, t.fs)
		l.AddDomain(t.domain) // make domain configurable only if needed
	}
	return l
}

// MustGetUserLocale returns the locale the user wants to use, omitting errors
func MustGetUserLocale(ctx context.Context, userID string, preferedLang string, vc settingssvc.ValueService) string {
	if preferedLang != "" {
		return preferedLang
	}

	locale, _ := GetUserLocale(ctx, userID, vc)
	return locale
}

// GetUserLocale returns the locale of the user
func GetUserLocale(ctx context.Context, userID string, vc settingssvc.ValueService) (string, error) {
	resp, err := vc.GetValueByUniqueIdentifiers(
		micrometadata.Set(ctx, middleware.AccountID, userID),
		&settingssvc.GetValueByUniqueIdentifiersRequest{
			AccountUuid: userID,
			SettingId:   defaults.SettingUUIDProfileLanguage,
		},
	)
	if err != nil {
		return "", err
	}
	val := resp.GetValue().GetValue().GetListValue().GetValues()
	if len(val) == 0 {
		return "", errors.New("no language setting found")
	}
	return val[0].GetStringValue(), nil
}
