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

var (
	// HeaderAcceptLanguage is the header key for the accept-language header
	HeaderAcceptLanguage = "Accept-Language"

	// ErrUnsupportedType is returned when the type is not supported
	ErrUnsupportedType = errors.New("unsupported type")
)

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

// TranslateEntity translates a slice, array or struct
func (t Translator) TranslateEntity(locale string, entity any, opts ...TranslateOption) error {
	return TranslateEntity(t.Locale(locale).Get, entity, opts...)
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

// TranslateOption is used to specify fields in structs or slices to translate
type TranslateOption func() (string, FieldType, []TranslateOption)

// FieldType is used to specify the type of field to translate
type FieldType int

const (
	// FieldTypeString is a string field
	FieldTypeString FieldType = iota
	// FieldTypeStruct is a struct field
	FieldTypeStruct
	// FieldTypeIterable is a slice or array field
	FieldTypeIterable
	// FieldTypeMap is a map field
	FieldTypeMap
)

// TranslateField function provides the generic way to translate the necessary field in composite entities.
func TranslateField(fieldName string) TranslateOption {
	return func() (string, FieldType, []TranslateOption) {
		return fieldName, FieldTypeString, nil
	}
}

// TranslateStruct function provides the generic way to translate the nested fields in composite entities.
func TranslateStruct(fieldName string, args ...TranslateOption) TranslateOption {
	return func() (string, FieldType, []TranslateOption) {
		return fieldName, FieldTypeStruct, args
	}
}

// TranslateEach function provides the generic way to translate the necessary fields in slices or nested entities.
func TranslateEach(fieldName string, args ...TranslateOption) TranslateOption {
	return func() (string, FieldType, []TranslateOption) {
		return fieldName, FieldTypeIterable, args
	}
}

// TranslateMap function provides the generic way to translate the necessary fields in maps.
// It's not implemented yet.
func TranslateMap(fieldName string, args ...TranslateOption) TranslateOption {
	return func() (string, FieldType, []TranslateOption) {
		return fieldName, FieldTypeMap, args
	}
}

// TranslateEntity function provides the generic way to translate the necessary fields in composite entities.
// The function takes a translation function that has the locale already set, see Translator.TranslateEntity
// The function also takes the entity with fields to translate.
// The function supports nested structs and slices of structs.
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
//		err := l10n.TranslateEntity(tr, s,
//					l10n.TranslateEach("StructList",
//						l10n.TranslateField("Description"),
//						l10n.TranslateField("DisplayName"),
//				))
func TranslateEntity(tr func(string, ...any) string, entity any, opts ...TranslateOption) error {
	value := reflect.ValueOf(entity)

	value, ok := cleanValue(value)
	if !ok {
		return errors.New("entity is not valid")
	}

	switch value.Kind() {
	case reflect.Struct:
		if !isStruct(value) {
			return fmt.Errorf("the root entity must be a struct, got %v", value.Kind())
		}
		rangeOverArgs(tr, value, opts...)
		return nil
	case reflect.Slice, reflect.Array:
		if len(opts) > 0 {
			translateEach(tr, value, opts...)
		} else {
			translateEach(tr, value)
		}
		return nil
	case reflect.Map:
		// TODO implement
	case reflect.String:
		translateField(tr, value)
		return nil
	}
	return ErrUnsupportedType
}

func translateEach(tr func(string, ...any) string, value reflect.Value, args ...TranslateOption) {
	value, ok := cleanValue(value)
	if !ok {
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

func rangeOverArgs(tr func(string, ...any) string, value reflect.Value, args ...TranslateOption) {
	value, ok := cleanValue(value)
	if !ok {
		return
	}

	for _, arg := range args {
		fieldName, fieldType, opts := arg()

		switch fieldType {
		case FieldTypeString:
			// exported field
			f := value.FieldByName(fieldName)
			translateField(tr, f)
		case FieldTypeStruct:
			// exported field
			innerValue := value.FieldByName(fieldName)
			if !innerValue.IsValid() || !isStruct(innerValue) {
				return
			}
			rangeOverArgs(tr, innerValue, opts...)
		case FieldTypeIterable:
			// exported field
			innerValue := value.FieldByName(fieldName)
			if !innerValue.IsValid() {
				return
			}
			if kind := innerValue.Kind(); kind != reflect.Array && kind != reflect.Slice {
				return
			}
			translateEach(tr, innerValue, opts...)
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

func cleanValue(v reflect.Value) (reflect.Value, bool) {
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return v, false
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return v, false
	}
	return v, true
}
