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

import (
	"fmt"
)

// KCError is an error as returned by Kopano core.
type KCError uint64

func (err KCError) Error() string {
	return fmt.Sprintf("%s (KC:0x%x)", KCErrorText(err), uint64(err))
}

// Kopano Core error codes as defined in common/include/kopano/kcodes.h.
const (
	KCERR_NONE            = iota
	KCERR_UNKNOWN KCError = (1 << 31) | iota
	KCERR_NOT_FOUND
	KCERR_NO_ACCESS
	KCERR_NETWORK_ERROR
	KCERR_SERVER_NOT_RESPONDING
	KCERR_INVALID_TYPE
	KCERR_DATABASE_ERROR
	KCERR_COLLISION
	KCERR_LOGON_FAILED
	KCERR_HAS_MESSAGES
	KCERR_HAS_FOLDERS
	KCERR_HAS_RECIPIENTS
	KCERR_HAS_ATTACHMENTS
	KCERR_NOT_ENOUGH_MEMORY
	KCERR_TOO_COMPLEX
	KCERR_END_OF_SESSION
	KCWARN_CALL_KEEPALIVE
	KCERR_UNABLE_TO_ABORT
	KCERR_NOT_IN_QUEUE
	KCERR_INVALID_PARAMETER
	KCWARN_PARTIAL_COMPLETION
	KCERR_INVALID_ENTRYID
	KCERR_BAD_VALUE
	KCERR_NO_SUPPORT
	KCERR_TOO_BIG
	KCWARN_POSITION_CHANGED
	KCERR_FOLDER_CYCLE
	KCERR_STORE_FULL
	KCERR_PLUGIN_ERROR
	KCERR_UNKNOWN_OBJECT
	KCERR_NOT_IMPLEMENTED
	KCERR_DATABASE_NOT_FOUND
	KCERR_INVALID_VERSION
	KCERR_UNKNOWN_DATABASE
	KCERR_NOT_INITIALIZED
	KCERR_CALL_FAILED
	KCERR_SSO_CONTINUE
	KCERR_TIMEOUT
	KCERR_INVALID_BOOKMARK
	KCERR_UNABLE_TO_COMPLETE
	KCERR_UNKNOWN_INSTANCE_ID
	KCERR_IGNORE_ME
	KCERR_BUSY
	KCERR_OBJECT_DELETED
	KCERR_USER_CANCEL
	KCERR_UNKNOWN_FLAGS
	KCERR_SUBMITTED
)

// KCSuccess defines success response as returned by Kopano core.
const KCSuccess = KCERR_NONE

// KCErrorTextMap maps the KCErrors to textual description.
var KCErrorTextMap = map[KCError]string{
	KCERR_UNKNOWN:               "Unknown",
	KCERR_NOT_FOUND:             "Not Found",
	KCERR_NO_ACCESS:             "No Access",
	KCERR_NETWORK_ERROR:         "Network Error",
	KCERR_SERVER_NOT_RESPONDING: "Server Not Responding",
	KCERR_INVALID_TYPE:          "Invalid Type",
	KCERR_DATABASE_ERROR:        "Database Error",
	KCERR_COLLISION:             "Object Collision",
	KCERR_LOGON_FAILED:          "Logon Failed",
	KCERR_HAS_MESSAGES:          "Object With Message Children",
	KCERR_HAS_FOLDERS:           "Object With Folder Children",
	KCERR_HAS_RECIPIENTS:        "Object With Recipient Children",
	KCERR_HAS_ATTACHMENTS:       "Object With Attachment Children",
	KCERR_NOT_ENOUGH_MEMORY:     "Not Enough Memory",
	KCERR_TOO_COMPLEX:           "Too Complex To Be Processed",
	KCERR_END_OF_SESSION:        "End Of Session",
	KCWARN_CALL_KEEPALIVE:       "",
	KCERR_UNABLE_TO_ABORT:       "Unable To Abort",
	KCERR_NOT_IN_QUEUE:          "",
	KCERR_INVALID_PARAMETER:     "Invalid Parameter",
	KCWARN_PARTIAL_COMPLETION:   "Partially Completed Operation",
	KCERR_INVALID_ENTRYID:       "Invalid EntryID",
	KCERR_BAD_VALUE:             "Type Error",
	KCERR_NO_SUPPORT:            "Not Supported",
	KCERR_TOO_BIG:               "Request Too Large",
	KCWARN_POSITION_CHANGED:     "",
	KCERR_FOLDER_CYCLE:          "",
	KCERR_STORE_FULL:            "Store Full Quota Reached",
	KCERR_PLUGIN_ERROR:          "Plugin Failed To Start",
	KCERR_UNKNOWN_OBJECT:        "",
	KCERR_NOT_IMPLEMENTED:       "Not Implemented",
	KCERR_DATABASE_NOT_FOUND:    "Database Not Found",
	KCERR_INVALID_VERSION:       "Database Version Unexpected",
	KCERR_UNKNOWN_DATABASE:      "",
	KCERR_NOT_INITIALIZED:       "",
	KCERR_CALL_FAILED:           "",
	KCERR_SSO_CONTINUE:          "SSO Success",
	KCERR_TIMEOUT:               "Timeout",
	KCERR_INVALID_BOOKMARK:      "",
	KCERR_UNABLE_TO_COMPLETE:    "",
	KCERR_UNKNOWN_INSTANCE_ID:   "",
	KCERR_IGNORE_ME:             "",
	KCERR_BUSY:                  "Task Already In Progress",
	KCERR_OBJECT_DELETED:        "Object Deleted",
	KCERR_USER_CANCEL:           "User Canceled Operation",
	KCERR_UNKNOWN_FLAGS:         "",
	KCERR_SUBMITTED:             "",
}

// KCErrorNameMap maps the KCErrors to their string names.
var KCErrorNameMap = map[KCError]string{
	KCERR_UNKNOWN:               "KCERR_UNKNOWN",
	KCERR_NOT_FOUND:             "KCERR_NOT_FOUND:",
	KCERR_NO_ACCESS:             "KCERR_NO_ACCESS",
	KCERR_NETWORK_ERROR:         "KCERR_NETWORK_ERROR",
	KCERR_SERVER_NOT_RESPONDING: "KCERR_SERVER_NOT_RESPONDING",
	KCERR_INVALID_TYPE:          "KCERR_INVALID_TYPE",
	KCERR_DATABASE_ERROR:        "KCERR_DATABASE_ERROR",
	KCERR_COLLISION:             "KCERR_COLLISION",
	KCERR_LOGON_FAILED:          "KCERR_LOGON_FAILED",
	KCERR_HAS_MESSAGES:          "KCERR_HAS_MESSAGE",
	KCERR_HAS_FOLDERS:           "KCERR_HAS_FOLDERS",
	KCERR_HAS_RECIPIENTS:        "KCERR_HAS_RECIPIENTS",
	KCERR_HAS_ATTACHMENTS:       "KCERR_HAS_ATTACHMENTS",
	KCERR_NOT_ENOUGH_MEMORY:     "KCERR_NOT_ENOUGH_MEMORY",
	KCERR_TOO_COMPLEX:           "KCERR_TOO_COMPLEX",
	KCERR_END_OF_SESSION:        "KCERR_END_OF_SESSION",
	KCWARN_CALL_KEEPALIVE:       "KCWARN_CALL_KEEPALIVE",
	KCERR_UNABLE_TO_ABORT:       "KCERR_UNABLE_TO_ABORT",
	KCERR_NOT_IN_QUEUE:          "KCERR_NOT_IN_QUEUE",
	KCERR_INVALID_PARAMETER:     "KCERR_INVALID_PARAMETER",
	KCWARN_PARTIAL_COMPLETION:   "KCWARN_PARTIAL_COMPLETION",
	KCERR_INVALID_ENTRYID:       "KCERR_INVALID_ENTRYID",
	KCERR_BAD_VALUE:             "KCERR_BAD_VALUE",
	KCERR_NO_SUPPORT:            "KCERR_NO_SUPPORT",
	KCERR_TOO_BIG:               "KCERR_TOO_BIG",
	KCWARN_POSITION_CHANGED:     "KCWARN_POSITION_CHANGED",
	KCERR_FOLDER_CYCLE:          "KCERR_FOLDER_CYCLE",
	KCERR_STORE_FULL:            "KCERR_STORE_FULL",
	KCERR_PLUGIN_ERROR:          "KCERR_PLUGIN_ERROR",
	KCERR_UNKNOWN_OBJECT:        "KCERR_UNKNOWN_OBJECT",
	KCERR_NOT_IMPLEMENTED:       "KCERR_NOT_IMPLEMENTED",
	KCERR_DATABASE_NOT_FOUND:    "KCERR_DATABASE_NOT_FOUND",
	KCERR_INVALID_VERSION:       "KCERR_INVALID_VERSION",
	KCERR_UNKNOWN_DATABASE:      "KCERR_UNKNOWN_DATABASE",
	KCERR_NOT_INITIALIZED:       "KCERR_NOT_INITIALIZED",
	KCERR_CALL_FAILED:           "KCERR_CALL_FAILED",
	KCERR_SSO_CONTINUE:          "KCERR_SSO_CONTINUE",
	KCERR_TIMEOUT:               "KCERR_TIMEOUT",
	KCERR_INVALID_BOOKMARK:      "KCERR_INVALID_BOOKMARK",
	KCERR_UNABLE_TO_COMPLETE:    "KCERR_UNABLE_TO_COMPLETE",
	KCERR_UNKNOWN_INSTANCE_ID:   "KCERR_UNKNOWN_INSTANCE_ID",
	KCERR_IGNORE_ME:             "KCERR_IGNORE_ME",
	KCERR_BUSY:                  "KCERR_BUSY",
	KCERR_OBJECT_DELETED:        "KCERR_OBJECT_DELETED",
	KCERR_USER_CANCEL:           "KCERR_USER_CANCEL",
	KCERR_UNKNOWN_FLAGS:         "KCERR_UNKNOWN_FLAGS",
	KCERR_SUBMITTED:             "KCERR_SUBMITTED",
}

// KCErrorText returns a text for the KC error. It returns the empty string if
// the code is unknown.
func KCErrorText(code KCError) string {
	text := KCErrorNameMap[code]
	description := KCErrorTextMap[code]
	if description != "" {
		return text + " (" + description + ")"
	}
	return text
}
