package conversions_test

import (
	"fmt"
	"testing"

	libregraph "github.com/owncloud/libre-graph-api-go"

	"github.com/owncloud/ocis/v2/ocis-pkg/conversions"
)

func checkIdentical[T any](t *testing.T, p T, want string) {
	t.Helper()
	got := fmt.Sprintf("%T", p)
	if got != want {
		t.Errorf("want:%q got:%q", want, got)
	}
}

func TestToValue(t *testing.T) {
	checkIdentical(t, conversions.ToValue((*int)(nil)), "int")
	checkIdentical(t, conversions.ToValue((*string)(nil)), "string")
	checkIdentical(t, conversions.ToValue((*float64)(nil)), "float64")
	checkIdentical(t, conversions.ToValue((*libregraph.UnifiedRoleDefinition)(nil)), "libregraph.UnifiedRoleDefinition")
	checkIdentical(t, conversions.ToValue((*[]string)(nil)), "[]string")
	checkIdentical(t, conversions.ToValue((*[]libregraph.UnifiedRoleDefinition)(nil)), "[]libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToValue(new("a")), "string")
	checkIdentical(t, conversions.ToValue(new(1)), "int")
	checkIdentical(t, conversions.ToValue(new(-1)), "int")
	checkIdentical(t, conversions.ToValue(new(float64(1))), "float64")
	checkIdentical(t, conversions.ToValue(new(float64(-1))), "float64")
	checkIdentical(t, conversions.ToValue(new(libregraph.UnifiedRoleDefinition{})), "libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToValue(new([]string{"a"})), "[]string")
	checkIdentical(t, conversions.ToValue(new([]int{1})), "[]int")
	checkIdentical(t, conversions.ToValue(new([]float64{1})), "[]float64")
	checkIdentical(t, conversions.ToValue(new([]libregraph.UnifiedRoleDefinition{{}})), "[]libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToValue(new(new("a"))), "*string")
	checkIdentical(t, conversions.ToValue(new(new(1))), "*int")
	checkIdentical(t, conversions.ToValue(new(new(-1))), "*int")
	checkIdentical(t, conversions.ToValue(new(new(float64(1)))), "*float64")
	checkIdentical(t, conversions.ToValue(new(new(float64(-1)))), "*float64")
	checkIdentical(t, conversions.ToValue(new(new(libregraph.UnifiedRoleDefinition{}))), "*libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToValue(new(new([]string{"a"}))), "*[]string")
	checkIdentical(t, conversions.ToValue(new(new([]int{1}))), "*[]int")
	checkIdentical(t, conversions.ToValue(new(new([]float64{1}))), "*[]float64")
	checkIdentical(t, conversions.ToValue(new(new([]libregraph.UnifiedRoleDefinition{{}}))), "*[]libregraph.UnifiedRoleDefinition")
}

func TestToPointerSlice(t *testing.T) {
	checkIdentical(t, conversions.ToPointerSlice([]string{"a"}), "[]*string")
	checkIdentical(t, conversions.ToPointerSlice([]int{1}), "[]*int")
	checkIdentical(t, conversions.ToPointerSlice([]libregraph.UnifiedRoleDefinition{{}}), "[]*libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToPointerSlice(([]string)(nil)), "[]*string")
	checkIdentical(t, conversions.ToPointerSlice(([]int)(nil)), "[]*int")
	checkIdentical(t, conversions.ToPointerSlice(([]libregraph.UnifiedRoleDefinition)(nil)), "[]*libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToPointerSlice([]*string{new("a")}), "[]**string")
	checkIdentical(t, conversions.ToPointerSlice([]*int{new(1)}), "[]**int")
	checkIdentical(t, conversions.ToPointerSlice(([]*libregraph.UnifiedRoleDefinition)(nil)), "[]**libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToPointerSlice(([]*string)(nil)), "[]**string")
	checkIdentical(t, conversions.ToPointerSlice(([]*int)(nil)), "[]**int")
	checkIdentical(t, conversions.ToPointerSlice(([]*libregraph.UnifiedRoleDefinition)(nil)), "[]**libregraph.UnifiedRoleDefinition")
}

func TestToValueSlice(t *testing.T) {
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice([]string{"a"})), "[]string")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice([]int{1})), "[]int")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice([]libregraph.UnifiedRoleDefinition{{}})), "[]libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice(([]string)(nil))), "[]string")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice(([]int)(nil))), "[]int")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice(([]libregraph.UnifiedRoleDefinition)(nil))), "[]libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice([]*string{new("a")})), "[]*string")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice([]*int{new(1)})), "[]*int")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice([]*libregraph.UnifiedRoleDefinition{new(libregraph.UnifiedRoleDefinition{})})), "[]*libregraph.UnifiedRoleDefinition")

	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice(([]*string)(nil))), "[]*string")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice(([]*int)(nil))), "[]*int")
	checkIdentical(t, conversions.ToValueSlice(conversions.ToPointerSlice(([]*libregraph.UnifiedRoleDefinition)(nil))), "[]*libregraph.UnifiedRoleDefinition")
}
