// package l10n holds translation mechanics that are used by user facing services (notifications, userlog, graph)
package l10n

import (
	"context"
	"errors"
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

type field func() string
type structs func() (string, []any)
type each func() (string, []any)

func TranslateField(fieldName string) field {
	return func() string {
		return fieldName
	}
}

func TranslateStruct(fieldName string, fn ...any) structs {
	return func() (string, []any) {
		return fieldName, fn
	}
}

func TranslateEach(fieldName string, fn ...any) each {
	return func() (string, []any) {
		return fieldName, fn
	}
}

// TranslateEntity function provides the generic way to translate the necessary fields in composite entities.
// The function takes the entity, translation function and fields to translate
// that are described by the TranslateField function. The function supports nested structs and slices of structs.
//
//				type InnreStruct struct {
//					Description string
//					DisplayName *string
//				}
//
//				type TopLevelStruct struct {
//					Description string
//					DisplayName *string
//					SubStruct   *InnreStruct
//	             StructList []*InnreStruct
//				}
//		    s:=  &TopLevelStruct{
//							Description: "description",
//							DisplayName: toStrPointer("displayName"),
//							SubStruct: &InnreStruct{
//								Description: "inner",
//								DisplayName: toStrPointer("innerDisplayName"),
//							},
//	                     StructList: []*InnreStruct{
//	             			{
//	             				Description: "inner 1",
//	             				DisplayName: toStrPointer("innerDisplayName 1"),
//	             			},
//	             			{
//	             				Description: "inner 2",
//	             				DisplayName: toStrPointer("innerDisplayName 2"),
//	             	 		},
//	             	   	},
//	                  }
//			 TranslateEntity(s, translateFunc(),
//			                    TranslateField("Description"),
//								TranslateField("DisplayName"),
//								TranslateStruct("SubStruct",
//									TranslateField("Description"),
//									TranslateField("DisplayName")),
//			                    TranslateEach("StructList",
//							        TranslateField("Description"),
//							        TranslateField("DisplayName")))
func TranslateEntity(entity any, tr func(string, ...any) string, fields ...any) error {
	value := reflect.ValueOf(entity)
	// Indirect through pointers and interfaces
	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Map:
		for i := 0; i < value.Len(); i++ {
			nextValue := value.Index(i)
			// Indirect through pointers and interfaces
			if nextValue.Kind() == reflect.Ptr || nextValue.Kind() == reflect.Interface {
				if nextValue.IsNil() {
					// treat a nil struct pointer as valid
					continue
				}
				nextValue = value.Index(i).Elem()
			}
			translateInner(nextValue, tr, fields...)
		}
		return nil
	}
	translateInner(value, tr, fields...)
	return nil
}

func translateInner(value reflect.Value, tr func(string, ...any) string, fields ...any) {
	if !value.IsValid() {
		return
	}
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
	for _, fl := range fields {
		switch fl := fl.(type) {
		case field:
			translateStringField(value, tr, fl)
		case each:
			translateEach(value, tr, fl)
		case structs:
			translateStruct(value, tr, fl)
		}
	}
}

func translateEach(value reflect.Value, tr func(string, ...any) string, fl each) {
	fieldName, fields := fl()
	// exported field
	innerValue := value.FieldByName(fieldName)
	if !innerValue.IsValid() {
		return
	}
	switch innerValue.Kind() {
	case reflect.Slice, reflect.Map:
		for i := 0; i < innerValue.Len(); i++ {
			translateInner(innerValue.Index(i), tr, fields...)
		}
	}
}

func translateStruct(value reflect.Value, tr func(string, ...any) string, fl structs) {
	fieldName, fields := fl()
	// exported field
	innerValue := value.FieldByName(fieldName)
	if !innerValue.IsValid() {
		return
	}
	if isStruct(innerValue) {
		translateInner(innerValue, tr, fields...)
		return
	}
}

func translateStringField(value reflect.Value, tr func(string, ...any) string, fl field) {
	fieldName := fl()
	// exported field
	f := value.FieldByName(fieldName)
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
