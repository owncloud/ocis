/*
 * Copyright 2017-2019 Kopano and its licensors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package ldap

// Define some known LDAP attribute descriptors.
const (
	AttributeDN         = "dn"
	AttributeLogin      = "uid"
	AttributeEmail      = "mail"
	AttributeName       = "cn"
	AttributeFamilyName = "sn"
	AttributeGivenName  = "givenName"
	AttributeUUID       = "uuid"
)

// Additional mappable virtual attributes.
const (
	AttributeNumericUID = "konnectNumericID"
)

// Define our known LDAP attribute value types.
const (
	AttributeValueTypeText   = "text"
	AttributeValueTypeBinary = "binary"
	AttributeValueTypeUUID   = "uuid"
)
