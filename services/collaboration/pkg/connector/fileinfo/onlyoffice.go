package fileinfo

// OnlyOffice fileInfo properties
//
// OnlyOffice WOPI check file info specification:
// https://api.onlyoffice.com/editors/wopi/restapi/checkfileinfo
type OnlyOffice struct {
	//
	// Required response properties
	//

	// copied from MS WOPI
	BaseFileName string `json:"BaseFileName,omitempty"`
	// copied from MS WOPI
	Version string `json:"Version,omitempty"`

	//
	// Breadcrumb properties
	//

	// copied from MS WOPI
	BreadcrumbBrandName string `json:"BreadcrumbBrandName,omitempty"`
	// copied from MS WOPI
	BreadcrumbBrandURL string `json:"BreadcrumbBrandUrl,omitempty"`
	// copied from MS WOPI
	BreadcrumbDocName string `json:"BreadcrumbDocName,omitempty"`
	// copied from MS WOPI
	BreadcrumbFolderName string `json:"BreadcrumbFolderName,omitempty"`
	// copied from MS WOPI
	BreadcrumbFolderURL string `json:"BreadcrumbFolderUrl,omitempty"`

	//
	// PostMessage properties
	//

	// Specifies if the WOPI client should notify the WOPI server in case the user closes the rendering or editing client currently using this file. The host expects to receive the UI_Close PostMessage when the Close UI in the online office is activated.
	ClosePostMessage bool `json:"ClosePostMessage,omitempty"`
	// Specifies if the WOPI client should notify the WOPI server in case the user tries to edit a file. The host expects to receive the UI_Edit PostMessage when the Edit UI in the online office is activated.
	EditModePostMessage bool `json:"EditModePostMessage,omitempty"`
	// Specifies if the WOPI client should notify the WOPI server in case the user tries to edit a file. The host expects to receive the Edit_Notification PostMessage.
	EditNotificationPostMessage bool `json:"EditNotificationPostMessage,omitempty"`
	// Specifies if the WOPI client should notify the WOPI server in case the user tries to share a file. The host expects to receive the UI_Sharing PostMessage when the Share UI in the online office is activated.
	FileSharingPostMessage bool `json:"FileSharingPostMessage,omitempty"`
	// Specifies if the WOPI client will notify the WOPI server in case the user tries to navigate to the previous file version. The host expects to receive the UI_FileVersions PostMessage when the Previous Versions UI in the online office is activated.
	FileVersionPostMessage bool `json:"FileVersionPostMessage,omitempty"`
	// A domain that the WOPI client must use as the targetOrigin parameter when sending messages as described in [W3C-HTML5WEBMSG].
	// copied from collabora WOPI
	PostMessageOrigin string `json:"PostMessageOrigin,omitempty"`

	//
	// File URL properties
	//

	// copied from MS WOPI
	CloseURL string `json:"CloseUrl,omitempty"`
	// copied from MS WOPI
	FileSharingURL string `json:"FileSharingUrl,omitempty"`
	// copied from MS WOPI
	FileVersionURL string `json:"FileVersionUrl,omitempty"`
	// copied from MS WOPI
	HostEditURL string `json:"HostEditUrl,omitempty"`

	//
	// Miscellaneous properties
	//

	// Specifies if the WOPI client must disable the Copy and Paste functionality within the application. By default, all Copy and Paste functionality is enabled, i.e. the setting has no effect. Possible property values:
	// BlockAll - the Copy and Paste functionality is completely disabled within the application;
	// CurrentDocumentOnly - the Copy and Paste functionality is enabled but content can only be copied and pasted within the file currently open in the application.
	// copied from MS WOPI
	CopyPasteRestrictions string `json:"CopyPasteRestrictions,omitempty"`
	// copied from MS WOPI
	DisablePrint bool `json:"DisablePrint"`
	// copied from MS WOPI
	FileExtension string `json:"FileExtension,omitempty"`
	// copied from MS WOPI
	FileNameMaxLength int `json:"FileNameMaxLength,omitempty"`
	// copied from MS WOPI
	LastModifiedTime string `json:"LastModifiedTime,omitempty"`

	//
	// User metadata properties
	//

	// copied from MS WOPI
	IsAnonymousUser bool `json:"IsAnonymousUser,omitempty"`
	// copied from MS WOPI
	UserFriendlyName string `json:"UserFriendlyName,omitempty"`
	// copied from MS WOPI
	UserID string `json:"UserId,omitempty"`

	//
	// User permissions properties
	//

	// copied from MS WOPI
	ReadOnly bool `json:"ReadOnly"`
	// copied from MS WOPI
	UserCanNotWriteRelative bool `json:"UserCanNotWriteRelative"`
	// copied from MS WOPI
	UserCanRename bool `json:"UserCanRename"`
	// Specifies if the user has permissions to review a file.
	UserCanReview bool `json:"UserCanReview,omitempty"`
	// copied from MS WOPI
	UserCanWrite bool `json:"UserCanWrite"`

	//
	// Host capabilities properties
	//

	// copied from MS WOPI
	SupportsLocks bool `json:"SupportsLocks"`
	// copied from MS WOPI
	SupportsRename bool `json:"SupportsRename"`
	// Specifies if the WOPI server supports the review permission.
	SupportsReviewing bool `json:"SupportsReviewing,omitempty"`
	// copied from MS WOPI
	SupportsUpdate bool `json:"SupportsUpdate"` // whether "Putfile" and "PutRelativeFile" work

	//
	// Other properties
	//

	// copied from collabora WOPI
	EnableInsertRemoteImage bool `json:"EnableInsertRemoteImage,omitempty"`
	// copied from collabora WOPI
	HidePrintOption bool `json:"HidePrintOption,omitempty"`
}

// SetProperties will set the file properties for the OnlyOffice implementation.
func (oinfo *OnlyOffice) SetProperties(props map[string]interface{}) {
	setters := map[string]func(value interface{}){
		"BaseFileName": assignStringTo(&oinfo.BaseFileName),
		"Version":      assignStringTo(&oinfo.Version),

		"BreadcrumbBrandName":  assignStringTo(&oinfo.BreadcrumbBrandName),
		"BreadcrumbBrandURL":   assignStringTo(&oinfo.BreadcrumbBrandURL),
		"BreadcrumbDocName":    assignStringTo(&oinfo.BreadcrumbDocName),
		"BreadcrumbFolderName": assignStringTo(&oinfo.BreadcrumbFolderName),
		"BreadcrumbFolderURL":  assignStringTo(&oinfo.BreadcrumbFolderURL),

		"ClosePostMessage":            assignBoolTo(&oinfo.ClosePostMessage),
		"EditModePostMessage":         assignBoolTo(&oinfo.EditModePostMessage),
		"EditNotificationPostMessage": assignBoolTo(&oinfo.EditNotificationPostMessage),
		"FileSharingPostMessage":      assignBoolTo(&oinfo.FileSharingPostMessage),
		"FileVersionPostMessage":      assignBoolTo(&oinfo.FileVersionPostMessage),
		"PostMessageOrigin":           assignStringTo(&oinfo.PostMessageOrigin),

		"CloseURL":       assignStringTo(&oinfo.CloseURL),
		"FileSharingURL": assignStringTo(&oinfo.FileSharingURL),
		"FileVersionURL": assignStringTo(&oinfo.FileVersionURL),
		"HostEditURL":    assignStringTo(&oinfo.HostEditURL),

		"CopyPasteRestrictions": assignStringTo(&oinfo.CopyPasteRestrictions),
		"DisablePrint":          assignBoolTo(&oinfo.DisablePrint),
		"FileExtension":         assignStringTo(&oinfo.FileExtension),
		"FileNameMaxLength":     assignIntTo(&oinfo.FileNameMaxLength),
		"LastModifiedTime":      assignStringTo(&oinfo.LastModifiedTime),

		"IsAnonymousUser":  assignBoolTo(&oinfo.IsAnonymousUser),
		"UserFriendlyName": assignStringTo(&oinfo.UserFriendlyName),
		"UserID":           assignStringTo(&oinfo.UserID),

		"ReadOnly":                assignBoolTo(&oinfo.ReadOnly),
		"UserCanNotWriteRelative": assignBoolTo(&oinfo.UserCanNotWriteRelative),
		"UserCanRename":           assignBoolTo(&oinfo.UserCanRename),
		"UserCanReview":           assignBoolTo(&oinfo.UserCanReview),
		"UserCanWrite":            assignBoolTo(&oinfo.UserCanWrite),

		"SupportsLocks":     assignBoolTo(&oinfo.SupportsLocks),
		"SupportsRename":    assignBoolTo(&oinfo.SupportsRename),
		"SupportsReviewing": assignBoolTo(&oinfo.SupportsReviewing),
		"SupportsUpdate":    assignBoolTo(&oinfo.SupportsUpdate),

		"EnableInsertRemoteImage": assignBoolTo(&oinfo.EnableInsertRemoteImage),
		"HidePrintOption":         assignBoolTo(&oinfo.HidePrintOption),
	}

	for key, value := range props {
		setterFn := setters[key]
		if setterFn != nil {
			setterFn(value)
		}
	}
}

// GetTarget will always return "OnlyOffice"
func (oinfo *OnlyOffice) GetTarget() string {
	return "OnlyOffice"
}
