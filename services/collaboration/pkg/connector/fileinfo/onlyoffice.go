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
	for key, value := range props {
		switch key {
		case "BaseFileName":
			oinfo.BaseFileName = value.(string)
		case "Version":
			oinfo.Version = value.(string)

		case "BreadcrumbBrandName":
			oinfo.BreadcrumbBrandName = value.(string)
		case "BreadcrumbBrandURL":
			oinfo.BreadcrumbBrandURL = value.(string)
		case "BreadcrumbDocName":
			oinfo.BreadcrumbDocName = value.(string)
		case "BreadcrumbFolderName":
			oinfo.BreadcrumbFolderName = value.(string)
		case "BreadcrumbFolderURL":
			oinfo.BreadcrumbFolderURL = value.(string)

		case "ClosePostMessage":
			oinfo.ClosePostMessage = value.(bool)
		case "EditModePostMessage":
			oinfo.EditModePostMessage = value.(bool)
		case "EditNotificationPostMessage":
			oinfo.EditNotificationPostMessage = value.(bool)
		case "FileSharingPostMessage":
			oinfo.FileSharingPostMessage = value.(bool)
		case "FileVersionPostMessage":
			oinfo.FileVersionPostMessage = value.(bool)
		case "PostMessageOrigin":
			oinfo.PostMessageOrigin = value.(string)

		case "CloseURL":
			oinfo.CloseURL = value.(string)
		case "FileSharingURL":
			oinfo.FileSharingURL = value.(string)
		case "FileVersionURL":
			oinfo.FileVersionURL = value.(string)
		case "HostEditURL":
			oinfo.HostEditURL = value.(string)

		case "CopyPasteRestrictions":
			oinfo.CopyPasteRestrictions = value.(string)
		case "DisablePrint":
			oinfo.DisablePrint = value.(bool)
		case "FileExtension":
			oinfo.FileExtension = value.(string)
		case "FileNameMaxLength":
			oinfo.FileNameMaxLength = value.(int)
		case "LastModifiedTime":
			oinfo.LastModifiedTime = value.(string)

		case "IsAnonymousUser":
			oinfo.IsAnonymousUser = value.(bool)
		case "UserFriendlyName":
			oinfo.UserFriendlyName = value.(string)
		case "UserID":
			oinfo.UserID = value.(string)

		case "ReadOnly":
			oinfo.ReadOnly = value.(bool)
		case "UserCanNotWriteRelative":
			oinfo.UserCanNotWriteRelative = value.(bool)
		case "UserCanRename":
			oinfo.UserCanRename = value.(bool)
		case "UserCanReview":
			oinfo.UserCanReview = value.(bool)
		case "UserCanWrite":
			oinfo.UserCanWrite = value.(bool)

		case "SupportsLocks":
			oinfo.SupportsLocks = value.(bool)
		case "SupportsRename":
			oinfo.SupportsRename = value.(bool)
		case "SupportsReviewing":
			oinfo.SupportsReviewing = value.(bool)
		case "SupportsUpdate":
			oinfo.SupportsUpdate = value.(bool)

		case "EnableInsertRemoteImage":
			oinfo.EnableInsertRemoteImage = value.(bool)
		case "HidePrintOption":
			oinfo.HidePrintOption = value.(bool)
		}
	}
}

// GetTarget will always return "OnlyOffice"
func (oinfo *OnlyOffice) GetTarget() string {
	return "OnlyOffice"
}
