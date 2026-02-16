// Package masker Provide mask format of Taiwan usually used(Name, Address, Email, ID ...etc.),
package masker

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"
)

const tagName = "mask"

type mtype string

// Maske Types of format string
const (
	MPassword   mtype = "password"
	MName       mtype = "name"
	MAddress    mtype = "addr"
	MEmail      mtype = "email"
	MMobile     mtype = "mobile"
	MTelephone  mtype = "tel"
	MID         mtype = "id"
	MCreditCard mtype = "credit"
	MStruct     mtype = "struct"
	MURL        mtype = "url"
)

// Masker is a instance to marshal masked string
type Masker struct {
	mask string
}

func strLoop(str string, length int) string {
	var mask string
	for i := 1; i <= length; i++ {
		mask += str
	}
	return mask
}

func (m *Masker) overlay(str string, overlay string, start int, end int) (overlayed string) {
	r := []rune(str)
	l := len([]rune(r))

	if l == 0 {
		return ""
	}

	if start < 0 {
		start = 0
	}
	if start > l {
		start = l
	}
	if end < 0 {
		end = 0
	}
	if end > l {
		end = l
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}

	overlayed = ""
	overlayed += string(r[:start])
	overlayed += overlay
	overlayed += string(r[end:])
	return overlayed
}

// Struct must input a interface{}, add tag mask on struct fields, after Struct(), return a pointer interface{} of input type and it will be masked with the tag format type
//
// Example:
//
//   type Foo struct {
//       Name      string `mask:"name"`
//       Email     string `mask:"email"`
//       Password  string `mask:"password"`
//       ID        string `mask:"id"`
//       Address   string `mask:"addr"`
//       Mobile    string `mask:"mobile"`
//       Telephone string `mask:"tel"`
//       Credit    string `mask:"credit"`
//       Foo       *Foo   `mask:"struct"`
//   }
//
//   func main() {
//       s := &Foo{
//           Name: ...,
//           Email: ...,
//           Password: ...,
//           Foo: &{
//               Name: ...,
//               Email: ...,
//               Password: ...,
//           }
//       }
//
//       m := masker.New()
//
//       t, err := m.Struct(s)
//
//       fmt.Println(t.(*Foo))
//   }
func (m *Masker) Struct(s interface{}) (interface{}, error) {
	if s == nil {
		return nil, fmt.Errorf("input is nil")
	}

	var selem, tptr reflect.Value

	st := reflect.TypeOf(s)

	if st.Kind() == reflect.Ptr {
		tptr = reflect.New(st.Elem())
		selem = reflect.ValueOf(s).Elem()
	} else {
		tptr = reflect.New(st)
		selem = reflect.ValueOf(s)
	}

	for i := 0; i < selem.NumField(); i++ {
		if !selem.Type().Field(i).IsExported() {
			continue
		}
		mtag := selem.Type().Field(i).Tag.Get(tagName)
		if len(mtag) == 0 {
			tptr.Elem().Field(i).Set(selem.Field(i))
			continue
		}
		switch selem.Field(i).Type().Kind() {
		default:
			tptr.Elem().Field(i).Set(selem.Field(i))
		case reflect.String:
			tptr.Elem().Field(i).SetString(m.String(mtype(mtag), selem.Field(i).String()))
		case reflect.Struct:
			if mtype(mtag) == MStruct {
				_t, err := m.Struct(selem.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t).Elem())
			}
		case reflect.Ptr:
			if selem.Field(i).IsNil() {
				continue
			}
			if mtype(mtag) == MStruct {
				_t, err := m.Struct(selem.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t))
			}
		case reflect.Slice:
			if selem.Field(i).IsNil() {
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.String {
				orgval := selem.Field(i).Interface().([]string)
				newval := make([]string, len(orgval))
				for i, val := range selem.Field(i).Interface().([]string) {
					newval[i] = m.String(mtype(mtag), val)
				}
				tptr.Elem().Field(i).Set(reflect.ValueOf(newval))
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Struct && mtype(mtag) == MStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(i).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Ptr && mtype(mtag) == MStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(i).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					newval = reflect.Append(newval, reflect.ValueOf(_n))
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
			if selem.Field(i).Type().Elem().Kind() == reflect.Interface && mtype(mtag) == MStruct {
				newval := reflect.MakeSlice(selem.Field(i).Type(), 0, selem.Field(i).Len())
				for j, l := 0, selem.Field(i).Len(); j < l; j++ {
					_n, err := m.Struct(selem.Field(i).Index(j).Interface())
					if err != nil {
						return nil, err
					}
					if reflect.TypeOf(selem.Field(i).Index(j).Interface()).Kind() != reflect.Ptr {
						newval = reflect.Append(newval, reflect.ValueOf(_n).Elem())
					} else {
						newval = reflect.Append(newval, reflect.ValueOf(_n))
					}
				}
				tptr.Elem().Field(i).Set(newval)
				continue
			}
		case reflect.Interface:
			if selem.Field(i).IsNil() {
				continue
			}
			if mtype(mtag) != MStruct {
				continue
			}
			_t, err := m.Struct(selem.Field(i).Interface())
			if err != nil {
				return nil, err
			}
			if reflect.TypeOf(selem.Field(i).Interface()).Kind() != reflect.Ptr {
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t).Elem())
			} else {
				tptr.Elem().Field(i).Set(reflect.ValueOf(_t))
			}
		}
	}

	return tptr.Interface(), nil
}

// String mask input string of the mask type
//
// Example:
//
//   masker.String(masker.MName, "ggwhite")
//   masker.String(masker.MID, "A123456789")
//   masker.String(masker.MMobile, "0987987987")
func (m *Masker) String(t mtype, i string) string {
	switch t {
	default:
		return i
	case MPassword:
		return m.Password(i)
	case MName:
		return m.Name(i)
	case MAddress:
		return m.Address(i)
	case MEmail:
		return m.Email(i)
	case MMobile:
		return m.Mobile(i)
	case MID:
		return m.ID(i)
	case MTelephone:
		return m.Telephone(i)
	case MCreditCard:
		return m.CreditCard(i)
	case MURL:
		return m.URL(i)
	}
}

// Name mask the second letter and the third letter
//
// Example:
//   input: ABCD
//   output: A**D
func (m *Masker) Name(i string) string {
	l := len([]rune(i))

	if l == 0 {
		return ""
	}

	// if has space
	if strs := strings.Split(i, " "); len(strs) > 1 {
		tmp := make([]string, len(strs))
		for idx, str := range strs {
			tmp[idx] = m.Name(str)
		}
		return strings.Join(tmp, " ")
	}

	if l == 2 || l == 3 {
		return m.overlay(i, strLoop(instance.mask, len("**")), 1, 2)
	}

	if l > 3 {
		return m.overlay(i, strLoop(instance.mask, len("**")), 1, 3)
	}

	return strLoop(instance.mask, len("**"))
}

// ID mask last 4 digits of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func (m *Masker) ID(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return m.overlay(i, strLoop(instance.mask, len("****")), 6, 10)
}

// Address keep first 6 letters, mask the rest
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func (m *Masker) Address(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	if l <= 6 {
		return strLoop(instance.mask, len("******"))
	}
	return m.overlay(i, strLoop(instance.mask, len("******")), 6, math.MaxInt)
}

// CreditCard mask 6 digits from the 7'th digit
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345` (American Express)(len = 15)
//   output2: 123456******345`
func (m *Masker) CreditCard(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return m.overlay(i, strLoop(instance.mask, len("******")), 6, 12)
}

// Email keep domain and the first 3 letters
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func (m *Masker) Email(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	tmp := strings.Split(i, "@")
	if len(tmp) == 1 {
		return m.overlay(i, strLoop(instance.mask, len("****")), 3, 7)
	}

	addr := tmp[0]
	domain := tmp[1]

	addr = m.overlay(addr, strLoop(instance.mask, len("****")), 3, 7)

	return addr + "@" + domain
}

// Mobile mask 3 digits from the 4'th digit
//
// Example:
//   input: 0987654321
//   output: 0987***321
func (m *Masker) Mobile(i string) string {
	if len(i) == 0 {
		return ""
	}
	return m.overlay(i, strLoop(instance.mask, len("***")), 4, 7)
}

// Telephone remove "(", ")", " ", "-" chart, and mask last 4 digits of telephone number, format to "(??)????-????"
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func (m *Masker) Telephone(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}

	i = strings.Replace(i, " ", "", -1)
	i = strings.Replace(i, "(", "", -1)
	i = strings.Replace(i, ")", "", -1)
	i = strings.Replace(i, "-", "", -1)

	l = len([]rune(i))

	if l != 10 && l != 8 {
		return i
	}

	ans := ""

	if l == 10 {
		ans += "("
		ans += i[:2]
		ans += ")"
		i = i[2:]
	}

	ans += i[:4]
	ans += "-"
	ans += "****"

	return ans
}

// Password always return "************"
func (m *Masker) Password(i string) string {
	l := len([]rune(i))
	if l == 0 {
		return ""
	}
	return strLoop(instance.mask, len("************"))
}

// URL mask the password part of the URL if exists
//
// Example:
//   input: http://admin:mysecretpassword@localhost:1234/uri
//   output:http://admin:xxxxx@localhost:1234/uri
func (m *Masker) URL(i string) string {
	u, err := url.Parse(i)
	if err != nil {
		return i
	}
	return u.Redacted()
}

// New create Masker
func New() *Masker {
	return &Masker{
		mask: "*",
	}
}

var instance *Masker

func init() {
	instance = New()
}

// Struct must input a interface{}, add tag mask on struct fields, after Struct(), return a pointer interface{} of input type and it will be masked with the tag format type
//
// Example:
//
//   type Foo struct {
//       Name      string `mask:"name"`
//       Email     string `mask:"email"`
//       Password  string `mask:"password"`
//       ID        string `mask:"id"`
//       Address   string `mask:"addr"`
//       Mobile    string `mask:"mobile"`
//       Telephone string `mask:"tel"`
//       Credit    string `mask:"credit"`
//       Foo       *Foo   `mask:"struct"`
//   }
//
//   func main() {
//       s := &Foo{
//           Name: ...,
//           Email: ...,
//           Password: ...,
//           Foo: &{
//               Name: ...,
//               Email: ...,
//               Password: ...,
//           }
//       }
//
//       t, err := masker.Struct(s)
//
//       fmt.Println(t.(*Foo))
//   }
func Struct(s interface{}) (interface{}, error) {
	return instance.Struct(s)
}

// String mask input string of the mask type
//
// Example:
//
//   masker.String(masker.MName, "ggwhite")
//   masker.String(masker.MID, "A123456789")
//   masker.String(masker.MMobile, "0987987987")
func String(t mtype, i string) string {
	return instance.String(t, i)
}

// Name mask the second letter and the third letter
//
// Example:
//   input: ABCD
//   output: A**D
func Name(i string) string {
	return instance.Name(i)
}

// ID mask last 4 digits of ID number
//
// Example:
//   input: A123456789
//   output: A12345****
func ID(i string) string {
	return instance.ID(i)
}

// Address keep first 6 letters, mask the rest
//
// Example:
//   input: 台北市內湖區內湖路一段737巷1號1樓
//   output: 台北市內湖區******
func Address(i string) string {
	return instance.Address(i)
}

// CreditCard mask 6 digits from the 7'th digit
//
// Example:
//   input1: 1234567890123456 (VISA, JCB, MasterCard)(len = 16)
//   output1: 123456******3456
//   input2: 123456789012345 (American Express)(len = 15)
//   output2: 123456******345
func CreditCard(i string) string {
	return instance.CreditCard(i)
}

// Email keep domain and the first 3 letters
//
// Example:
//   input: ggw.chang@gmail.com
//   output: ggw****@gmail.com
func Email(i string) string {
	return instance.Email(i)
}

// Mobile mask 3 digits from the 4'th digit
//
// Example:
//   input: 0987654321
//   output: 0987***321
func Mobile(i string) string {
	return instance.Mobile(i)
}

// Telephone remove "(", ")", " ", "-" chart, and mask last 4 digits of telephone number, format to "(??)????-????"
//
// Example:
//   input: 0227993078
//   output: (02)2799-****
func Telephone(i string) string {
	return instance.Telephone(i)
}

// Password always return "************"
func Password(i string) string {
	return instance.Password(i)
}

func SetMask(mask string) {
	instance.mask = mask
}
