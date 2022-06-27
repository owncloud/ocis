/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package kcc

// A LogonResponse holds tthe returned data of a SOAP logon request.
type LogonResponse struct {
	Er         KCError     `xml:"er" json:"-"`
	SessionID  KCSessionID `xml:"ulSessionId" json:"ulSessionId"`
	ServerGUID string      `xml:"sServerGuid" json:"sServerGuid"`
}

// A LogoffResponse holds the returned data of a SOAP logoff request.
type LogoffResponse struct {
	Er KCError `xml:"er"`
}

// A ResolveUserResponse holds the returned data of a SOAP request which
// retruns a user's ID details.
type ResolveUserResponse struct {
	Er          KCError `xml:"er"`
	ID          uint64  `xml:"ulUserId"`
	UserEntryID string  `xml:"sUserId"`
}

// A GetUserResponse holds the returned data of a SOAP request which fetches
// user detail meta data.
type GetUserResponse struct {
	Er   KCError `xml:"er"`
	User *User   `xml:"lpsUser"`
}

// ABResolveNamesResponse holds the returned data of a SOAP request which
// resolves names.
type ABResolveNamesResponse struct {
	Er     KCError          `xml:"er"`
	RowSet []*PropTagRowSet `xml:"sRowSet>item"`
	Flags  []ABFlag         `xml:"aFlags>item"`
}

// A User represents the meta data of a user as stored by Kopano server.
type User struct {
	ID          uint64     `xml:"ulUserId" json:"ulUserID"`
	Username    string     `xml:"lpszUsername" json:"lpszUsername"`
	MailAddress string     `xml:"lpszMailAddress" json:"lpszMailAddress"`
	FullName    string     `xml:"lpszFullName" json:"lpszFullName"`
	IsAdmin     uint64     `xml:"ulIsAdmin" json:"ulIsAdmin"`
	IsNonActive uint64     `xml:"ulIsNonActive" json:"ulIsNonActive"`
	UserEntryID string     `xml:"sUserId" json:"sUserId"`
	Props       *PropMap   `xml:"lpsPropmap>item" json:"lpsPropmap"`
	MVProps     *MVPropMap `xml:"lpsMVPropmap>item" json:"lpsMVPropmap"`
}

// A PropMap is a mapping of property IDs to a value.
type PropMap []*PropMapValue

// Get returns the accociaged PropMap's value for the provided id. When the
// property is not found, an empty string and false is returned.
func (pm PropMap) Get(id PT) (string, bool) {
	for _, value := range pm {
		if id == value.ID {
			return value.StringValue, true
		}
	}

	return "", false
}

// A PropMapValue represents a single string Value with an ID.
type PropMapValue struct {
	ID          PT     `xml:"ulPropId" json:"ulPropId"`
	StringValue string `xml:"lpszValue" json:"lpszValue"`
}

// A MVPropMap is a mapping of properties to a array of values.
type MVPropMap []*MVPropMapValue

// Get returns the accociaged MVPropMap's value for the provided id. When the
// property is not found, nil and false is returned.
func (pm MVPropMap) Get(id PT) ([]string, bool) {
	for _, value := range pm {
		if id == value.ID {
			return value.StringValues, true
		}
	}

	return nil, false
}

// A MVPropMapValue represents a set of string values with an ID.
type MVPropMapValue struct {
	ID           PT       `xml:"ulPropId" json:"ulPropId"`
	StringValues []string `xml:"sValues>item" json:"sValues"`
}

// A PropTagRowSet represents a row set of array type with prop tag items.
type PropTagRowSet struct {
	PropTagValues []*PropTagRowSetValue `xml:"item,omitempty" json:"items"`
}

// A PropTagRowSetValue represents a prop tag row set value item.
type PropTagRowSetValue struct {
	PropTag      PT         `xml:"ulPropTag" json:"ulPropTag"`
	AStringValue string     `xml:"lpszA" json:"lpszA,omitempty"`
	ULValue      uint64     `xml:"ul" json:"ul,omitempty"`
	BinValue     []byte     `xml:"bin" json:"bin,omitempty"`
	BinValues    [][][]byte `xml:"mvbin>item" json:"mvbin,omitempty"`
}
