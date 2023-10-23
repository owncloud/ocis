package config

// Options are the option for the web
type Options struct {
	HomeFolder               string           `json:"homeFolder,omitempty" yaml:"homeFolder" env:"WEB_OPTION_HOME_FOLDER" desc:"Specifies a folder that is used when the user navigates 'home'. Navigating home gets triggered by clicking on the 'All files' menu item. The user will not be jailed in that directory, it simply serves as a default location. A static location can be provided, or variables of the user object to come up with a user specific home path can be used. This uses the twig template variable style and allows picking a value or a substring of a value of the authenticated user. Examples are '/Shares', '/{{.Id}}' and '/{{substr 0 3 .Id}}/{{.Id}'."`
	OpenAppsInTab            bool             `json:"openAppsInTab,omitempty" yaml:"openAppsInTab" env:"WEB_OPTION_OPEN_APPS_IN_TAB" desc:"Configures whether apps and extensions should generally open in a new tab. Defaults to false."`
	DisablePreviews          bool             `json:"disablePreviews,omitempty" yaml:"disablePreviews" env:"OCIS_DISABLE_PREVIEWS;WEB_OPTION_DISABLE_PREVIEWS" desc:"Set this option to 'true' to disable previews in all the different web file listing views. This can speed up file listings in folders with many files. The only list view that is not affected by this setting is the trash bin, as it does not allow previewing at all."`
	PreviewFileMimeTypes     []string         `json:"previewFileMimeTypes,omitempty" yaml:"previewFileMimeTypes" env:"WEB_OPTION_PREVIEW_FILE_MIMETYPES" desc:"Specifies which mimeTypes will be previewed in the UI. For example to only preview jpg and text files, set this option to ['image/jpeg', 'text/plain']."`
	AccountEditLink          *AccountEditLink `json:"accountEditLink,omitempty" yaml:"accountEditLink"`
	DisableFeedbackLink      bool             `json:"disableFeedbackLink,omitempty" yaml:"disableFeedbackLink" env:"WEB_OPTION_DISABLE_FEEDBACK_LINK" desc:"Set this option to 'true' to disable the feedback link in the topbar. Keeping it enabled by setting the value to 'false' or with the absence of the option, allows ownCloud to get feedback from your user base through a dedicated survey website."`
	FeedbackLink             *FeedbackLink    `json:"feedbackLink,omitempty" yaml:"feedbackLink"`
	SharingRecipientsPerPage int              `json:"sharingRecipientsPerPage,omitempty" yaml:"sharingRecipientsPerPage" env:"WEB_OPTION_SHARING_RECIPIENTS_PER_PAGE" desc:"Sets the amount of users shown as recipients in the dropdown menu when sharing resources. Default amount is 200."`
	Sidebar                  Sidebar          `json:"sidebar" yaml:"sidebar"`
	RunningOnEOS             bool             `json:"runningOnEos,omitempty" yaml:"runningOnEos" env:"WEB_OPTION_RUNNING_ON_EOS" desc:"Set this option to 'true' if running on an EOS storage backend (https://eos-web.web.cern.ch/eos-web/) to enable its specific features. Defaults to 'false'."`
	CernFeatures             bool             `json:"cernFeatures,omitempty" yaml:"cernFeatures"`
	HoverableQuickActions    bool             `json:"hoverableQuickActions,omitempty" yaml:"hoverableQuickActions" env:"WEB_OPTION_HOVERABLE_QUICK_ACTIONS" desc:"Set this option to 'true' to hide quick actions (buttons appearing on file rows) and only show them when the user hovers over the row with his mouse. Defaults to 'false'."`
	Routing                  Routing          `json:"routing" yaml:"routing"`
	Upload                   *Upload          `json:"upload,omitempty" yaml:"upload"`
	Editor                   *Editor          `json:"editor,omitempty" yaml:"editor"`
	ContextHelpersReadMore   bool             `json:"contextHelpersReadMore,omitempty" yaml:"contextHelpersReadMore" env:"WEB_OPTION_CONTEXTHELPERS_READ_MORE" desc:"Specifies whether the 'Read more' link should be displayed or not."`
	LogoutURL                string           `json:"logoutUrl,omitempty" yaml:"logoutUrl" env:"WEB_OPTION_LOGOUT_URL" desc:"Adds a link to the user's profile page to point him to an external page, where he can manage his session and devices. This is helpful when an external IdP is used. This option is disabled by default."`
	OpenLinksWithDefaultApp  bool             `json:"openLinksWithDefaultApp,omitempty" yaml:"openLinksWithDefaultApp" env:"WEB_OPTION_OPEN_LINKS_WITH_DEFAULT_APP" desc:"Specifies whether single file link shares should be opened with the default app or not. If not opened by the default app, the Web UI just displays the file details. Defaults to 'true'."`
	ImprintURL               string           `json:"imprintUrl,omitempty" yaml:"imprintUrl" env:"WEB_OPTION_IMPRINT_URL" desc:"Specifies the target URL for the imprint link valid for the ocis instance in the account menu."`
	PrivacyURL               string           `json:"privacyUrl,omitempty" yaml:"privacyUrl" env:"WEB_OPTION_PRIVACY_URL" desc:"Specifies the target URL for the privacy link valid for the ocis instance in the account menu."`
	AccessDeniedHelpURL      string           `json:"accessDeniedHelpUrl,omitempty" yaml:"accessDeniedHelpUrl" env:"WEB_OPTION_ACCESS_DENIED_HELP_URL" desc:"Specifies the target URL valid for the ocis instance for the generic logged out / access denied page."`
	TokenStorageLocal        bool             `json:"tokenStorageLocal" yaml:"tokenStorageLocal" env:"WEB_OPTION_TOKEN_STORAGE_LOCAL" desc:"Specifies whether the access token will be stored in the local storage when set to 'true' or in the session storage when set to 'false''. If stored in the local storage, login state will be persisted across multiple browser tabs, means no additional logins are required. Defaults to 'true'."`
}

// AccountEditLink are the AccountEditLink options
type AccountEditLink struct {
	Href string `json:"href,omitempty" yaml:"href" env:"WEB_OPTION_ACCOUNT_EDIT_LINK_HREF" desc:"Set a different target URL for the edit link. Make sure to prepend it with 'http(s)://'."`
}

// FeedbackLink are the feedback link options
type FeedbackLink struct {
	Href        string `json:"href,omitempty" yaml:"href" env:"WEB_OPTION_FEEDBACKLINK_HREF" desc:"Set a target URL for the feedback link. Make sure to prepend it with 'http(s)://'. Defaults to 'https://owncloud.com/web-design-feedback'."`
	AriaLabel   string `json:"ariaLabel,omitempty" yaml:"ariaLabel" env:"WEB_OPTION_FEEDBACKLINK_ARIALABEL" desc:"Since the feedback link only has an icon, a screen reader accessible label can be set. The text defaults to 'ownCloud feedback survey'."`
	Description string `json:"description,omitempty" yaml:"description" env:"WEB_OPTION_FEEDBACKLINK_DESCRIPTION" desc:"For feedbacks, provide any description you want to see as tooltip and as accessible description. Defaults to 'Provide your feedback: We'd like to improve the web design and would be happy to hear your feedback. Thank you! Your ownCloud team'."`
}

// Sidebar are the sidebar option
type Sidebar struct {
	Shares SidebarShares `json:"shares" yaml:"shares"`
}

// SidebarShares are the options for the shares sidebar
type SidebarShares struct {
	ShowAllOnLoad bool `json:"showAllOnLoad" yaml:"showAllOnLoad" env:"WEB_OPTION_SIDEBAR_SHARES_SHOW_ALL_ON_LOAD" desc:"Sets the list of the (link) shares list in the sidebar to be initially expanded. Default is a collapsed state, only showing the first three shares."`
}

// Routing are the routing options
type Routing struct {
	IDBased bool `json:"idBased" yaml:"idBased" env:"WEB_OPTION_ROUTING_ID_BASED" desc:"Enable or disable fileIds being added to the URL. Defaults to 'true', because otherwise spaces with name clashes cannot be resolved correctly. Note: Only disable this if you can guarantee on the server side, that spaces of the same namespace cannot have name clashes."`
}

// Upload are the upload options
type Upload struct {
	XHR          XHR    `json:"xhr,omitempty" yaml:"xhr"`
	CompanionURL string `json:"companionUrl,omitempty" yaml:"companionUrl" env:"WEB_OPTION_UPLOAD_COMPANION_URL" desc:"Sets the URL of Companion which is a service provided by Uppy to import files from external cloud providers. See https://uppy.io/docs/companion/ for instructions on how to set up Companion. This feature is disabled as long as no URL is given."`
}

// XHR are the XHR options
type XHR struct {
	Timeout int `json:"timeout,omitempty" yaml:"timeout" env:"WEB_OPTION_UPLOAD_XHR_TIMEOUT" desc:"Specifies the timeout for XHR uploads in milliseconds."`
}

// Editor are the web editor options
type Editor struct {
	AutosaveEnabled  bool `json:"autosaveEnabled,omitempty" yaml:"autosaveEnabled" env:"WEB_OPTION_EDITOR_AUTOSAVE_ENABLED" desc:"Specifies if the autosave for the editor apps is enabled."`
	AutosaveInterval int  `json:"autosaveInterval,omitempty" yaml:"autosaveInterval" env:"WEB_OPTION_EDITOR_AUTOSAVE_INTERVAL" desc:"Specifies the time interval for the autosave of editor apps in seconds. Has no effect when WEB_OPTION_EDITOR_AUTOSAVE_ENABLED is set to 'false'."`
}
