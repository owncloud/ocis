package godata

import (
	"encoding/xml"
)

const (
	GoDataString         = "Edm.String"
	GoDataInt16          = "Edm.Int16"
	GoDataInt32          = "Edm.Int32"
	GoDataInt64          = "Edm.Int64"
	GoDataDecimal        = "Edm.Decimal"
	GoDataBinary         = "Edm.Binary"
	GoDataBoolean        = "Edm.Boolean"
	GoDataTimeOfDay      = "Edm.TimeOfDay"
	GoDataDate           = "Edm.Date"
	GoDataDateTimeOffset = "Edm.DateTimeOffset"
)

type GoDataMetadata struct {
	XMLName      xml.Name `xml:"edmx:Edmx"`
	XMLNamespace string   `xml:"xmlns:edmx,attr"`
	Version      string   `xml:"Version,attr"`
	DataServices *GoDataServices
	References   []*GoDataReference
}

func (t *GoDataMetadata) Bytes() ([]byte, error) {
	output, err := xml.MarshalIndent(t, "", "    ")
	if err != nil {
		return nil, err
	}

	return append([]byte(xml.Header), output...), nil
}

func (t *GoDataMetadata) String() string {
	return ""
}

type GoDataReference struct {
	XMLName            xml.Name `xml:"edmx:Reference"`
	Uri                string   `xml:"Uri,attr"`
	Includes           []*GoDataInclude
	IncludeAnnotations []*GoDataIncludeAnnotations
}

type GoDataInclude struct {
	XMLName   xml.Name `xml:"edmx:Include"`
	Namespace string   `xml:"Namespace,attr"`
	Alias     string   `xml:"Alias,attr,omitempty"`
}

type GoDataIncludeAnnotations struct {
	XMLName         xml.Name `xml:"edmx:IncludeAnnotations"`
	TermNamespace   string   `xml:"TermNamespace,attr"`
	Qualifier       string   `xml:"Qualifier,attr,omitempty"`
	TargetNamespace string   `xml:"TargetNamespace,attr,omitempty"`
}

type GoDataServices struct {
	XMLName xml.Name `xml:"edmx:DataServices"`
	Schemas []*GoDataSchema
}

type GoDataSchema struct {
	XMLName          xml.Name `xml:"Schema"`
	Namespace        string   `xml:"Namespace,attr"`
	Alias            string   `xml:"Alias,attr,omitempty"`
	Actions          []*GoDataAction
	Annotations      []*GoDataAnnotations
	Annotation       []*GoDataAnnotation
	ComplexTypes     []*GoDataComplexType
	EntityContainers []*GoDataEntityContainer
	EntityTypes      []*GoDataEntityType
	EnumTypes        []*GoDataEnumType
	Functions        []*GoDataFunction
	Terms            []*GoDataTerm
	TypeDefinitions  []*GoDataTypeDefinition
}

type GoDataAction struct {
	XMLName       xml.Name `xml:"Action"`
	Name          string   `xml:"Name,attr"`
	IsBound       string   `xml:"IsBound,attr,omitempty"`
	EntitySetPath string   `xml:"EntitySetPath,attr,omitempty"`
	Parameters    []*GoDataParameter
	ReturnType    *GoDataReturnType
}

type GoDataAnnotations struct {
	XMLName     xml.Name `xml:"Annotations"`
	Target      string   `xml:"Target,attr"`
	Qualifier   string   `xml:"Qualifier,attr,omitempty"`
	Annotations []*GoDataAnnotation
}

type GoDataAnnotation struct {
	XMLName   xml.Name `xml:"Annotation"`
	Term      string   `xml:"Term,attr"`
	Qualifier string   `xml:"Qualifier,attr,omitempty"`
}

type GoDataComplexType struct {
	XMLName              xml.Name `xml:"ComplexType"`
	Name                 string   `xml:"Name,attr"`
	BaseType             string   `xml:"BaseType,attr,omitempty"`
	Abstract             string   `xml:"Abstract,attr,omitempty"`
	OpenType             string   `xml:"OpenType,attr,omitempty"`
	Properties           []*GoDataProperty
	NavigationProperties []*GoDataNavigationProperty
}

type GoDataEntityContainer struct {
	XMLName         xml.Name `xml:"EntityContainer"`
	Name            string   `xml:"Name,attr"`
	Extends         string   `xml:"Extends,attr,omitempty"`
	EntitySets      []*GoDataEntitySet
	Singletons      []*GoDataSingleton
	ActionImports   []*GoDataActionImport
	FunctionImports []*GoDataFunctionImport
}

type GoDataEntityType struct {
	XMLName              xml.Name `xml:"EntityType"`
	Name                 string   `xml:"Name,attr"`
	BaseType             string   `xml:"BaseType,attr,omitempty"`
	Abstract             string   `xml:"Abstract,attr,omitempty"`
	OpenType             string   `xml:"OpenType,attr,omitempty"`
	HasStream            string   `xml:"HasStream,attr,omitempty"`
	Key                  *GoDataKey
	Properties           []*GoDataProperty
	NavigationProperties []*GoDataNavigationProperty
}

type GoDataEnumType struct {
	XMLName        xml.Name `xml:"EnumType"`
	Name           string   `xml:"Name,attr"`
	UnderlyingType string   `xml:"UnderlyingType,attr,omitempty"`
	IsFlags        string   `xml:"IsFlags,attr,omitempty"`
	Members        []*GoDataMember
}

type GoDataFunction struct {
	XMLName       xml.Name `xml:"Function"`
	Name          string   `xml:"Name,attr"`
	IsBound       string   `xml:"IsBound,attr,omitempty"`
	IsComposable  string   `xml:"IsComposable,attr,omitempty"`
	EntitySetPath string   `xml:"EntitySetPath,attr,omitempty"`
	Parameters    []*GoDataParameter
	ReturnType    *GoDataReturnType
}

type GoDataTypeDefinition struct {
	XMLName        xml.Name `xml:"TypeDefinition"`
	Name           string   `xml:"Name,attr"`
	UnderlyingType string   `xml:"UnderlyingTypeattr,omitempty"`
	Annotations    []*GoDataAnnotation
}

type GoDataProperty struct {
	XMLName      xml.Name `xml:"Property"`
	Name         string   `xml:"Name,attr"`
	Type         string   `xml:"Type,attr"`
	Nullable     string   `xml:"Nullable,attr,omitempty"`
	MaxLength    int      `xml:"MaxLength,attr,omitempty"`
	Precision    int      `xml:"Precision,attr,omitempty"`
	Scale        int      `xml:"Scale,attr,omitempty"`
	Unicode      string   `xml:"Unicode,attr,omitempty"`
	SRID         string   `xml:"SRID,attr,omitempty"`
	DefaultValue string   `xml:"DefaultValue,attr,omitempty"`
}

type GoDataNavigationProperty struct {
	XMLName                xml.Name `xml:"NavigationProperty"`
	Name                   string   `xml:"Name,attr"`
	Type                   string   `xml:"Type,attr"`
	Nullable               string   `xml:"Nullable,attr,omitempty"`
	Partner                string   `xml:"Partner,attr,omitempty"`
	ContainsTarget         string   `xml:"ContainsTarget,attr,omitempty"`
	ReferentialConstraints []*GoDataReferentialConstraint
}

type GoDataReferentialConstraint struct {
	XMLName            xml.Name        `xml:"ReferentialConstraint"`
	Property           string          `xml:"Property,attr"`
	ReferencedProperty string          `xml:"ReferencedProperty,attr"`
	OnDelete           *GoDataOnDelete `xml:"OnDelete,omitempty"`
}

type GoDataOnDelete struct {
	XMLName xml.Name `xml:"OnDelete"`
	Action  string   `xml:"Action,attr"`
}

type GoDataEntitySet struct {
	XMLName                    xml.Name `xml:"EntitySet"`
	Name                       string   `xml:"Name,attr"`
	EntityType                 string   `xml:"EntityType,attr"`
	IncludeInServiceDocument   string   `xml:"IncludeInServiceDocument,attr,omitempty"`
	NavigationPropertyBindings []*GoDataNavigationPropertyBinding
}

type GoDataSingleton struct {
	XMLName                    xml.Name `xml:"Singleton"`
	Name                       string   `xml:"Name,attr"`
	Type                       string   `xml:"Type,attr"`
	NavigationPropertyBindings []*GoDataNavigationPropertyBinding
}

type GoDataNavigationPropertyBinding struct {
	XMLName xml.Name `xml:"NavigationPropertyBinding"`
	Path    string   `xml:"Path,attr"`
	Target  string   `xml:"Target,attr"`
}

type GoDataActionImport struct {
	XMLName   xml.Name `xml:"ActionImport"`
	Name      string   `xml:"Name,attr"`
	Action    string   `xml:"Action,attr"`
	EntitySet string   `xml:"EntitySet,attr,omitempty"`
}

type GoDataFunctionImport struct {
	XMLName                  xml.Name `xml:"FunctionImport"`
	Name                     string   `xml:"Name,attr"`
	Function                 string   `xml:"Function,attr"`
	EntitySet                string   `xml:"EntitySet,attr,omitempty"`
	IncludeInServiceDocument string   `xml:"IncludeInServiceDocument,attr,omitempty"`
}

type GoDataKey struct {
	XMLName     xml.Name `xml:"Key"`
	PropertyRef *GoDataPropertyRef
}

type GoDataPropertyRef struct {
	XMLName xml.Name `xml:"PropertyRef"`
	Name    string   `xml:"Name,attr"`
}

type GoDataParameter struct {
	XMLName   xml.Name `xml:"Parameter"`
	Name      string   `xml:"Name,attr"`
	Type      string   `xml:"Type,attr"`
	Nullable  string   `xml:"Nullable,attr,omitempty"`
	MaxLength int      `xml:"MaxLength,attr,omitempty"`
	Precision int      `xml:"Precision,attr,omitempty"`
	Scale     int      `xml:"Scale,attr,omitempty"`
	SRID      string   `xml:"SRID,attr,omitempty"`
}

type GoDataReturnType struct {
	XMLName   xml.Name `xml:"ReturnType"`
	Name      string   `xml:"Name,attr"`
	Type      string   `xml:"Type,attr"`
	Nullable  string   `xml:"Nullable,attr,omitempty"`
	MaxLength int      `xml:"MaxLength,attr,omitempty"`
	Precision int      `xml:"Precision,attr,omitempty"`
	Scale     int      `xml:"Scale,attr,omitempty"`
	SRID      string   `xml:"SRID,attr,omitempty"`
}

type GoDataMember struct {
	XMLName xml.Name `xml:"Member"`
	Name    string   `xml:"Name,attr"`
	Value   string   `xml:"Value,attr,omitempty"`
}

type GoDataTerm struct {
	XMLName      xml.Name `xml:"Term"`
	Name         string   `xml:"Name,attr"`
	Type         string   `xml:"Type,attr"`
	BaseTerm     string   `xml:"BaseTerm,attr,omitempty"`
	DefaultValue string   `xml:"DefaultValue,attr,omitempty"`
	AppliesTo    string   `xml:"AppliesTo,attr,omitempty"`
}
