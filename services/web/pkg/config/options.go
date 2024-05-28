package config

// Options are the option for the web
type Options struct {
	HomeFolder               string              `json:"homeFolder,omitempty" yaml:"homeFolder" env:"WEB_OPTION_HOME_FOLDER" desc:"Specifies a folder that is used when the user navigates 'home'. Navigating home gets triggered by clicking on the 'All files' menu item. The user will not be jailed in that directory, it simply serves as a default location. A static location can be provided, or variables of the user object to come up with a user specific home path can be used. This uses the twig template variable style and allows picking a value or a substring of a value of the authenticated user. Examples are '/Shares', '/{{.Id}}' and '/{{substr 0 3 .Id}}/{{.Id}'." introductionVersion:"pre5.0"`
	OpenAppsInTab            bool                `json:"openAppsInTab,omitempty" yaml:"openAppsInTab" env:"WEB_OPTION_OPEN_APPS_IN_TAB" desc:"Configures whether apps and extensions should generally open in a new tab. Defaults to false." introductionVersion:"pre5.0"`
	DisablePreviews          bool                `json:"disablePreviews,omitempty" yaml:"disablePreviews" env:"OCIS_DISABLE_PREVIEWS;WEB_OPTION_DISABLE_PREVIEWS" desc:"Set this option to 'true' to disable previews in all the different web file listing views. This can speed up file listings in folders with many files. The only list view that is not affected by this setting is the trash bin, as it does not allow previewing at all." introductionVersion:"pre5.0"`
	PreviewFileMimeTypes     []string            `json:"previewFileMimeTypes,omitempty" yaml:"previewFileMimeTypes" env:"WEB_OPTION_PREVIEW_FILE_MIMETYPES" desc:"A list of mimeTypes to specify which ones will be previewed in the UI. For example, to only preview jpg and text files, set this option to 'image/jpeg,text/plain'. See the Environment Variable Types description for more details." introductionVersion:"pre5.0"`
	AccountEditLink          *AccountEditLink    `json:"accountEditLink,omitempty" yaml:"accountEditLink"`
	DisableFeedbackLink      bool                `json:"disableFeedbackLink,omitempty" yaml:"disableFeedbackLink" env:"WEB_OPTION_DISABLE_FEEDBACK_LINK" desc:"Set this option to 'true' to disable the feedback link in the topbar. Keeping it enabled by setting the value to 'false' or with the absence of the option, allows ownCloud to get feedback from your user base through a dedicated survey website." introductionVersion:"pre5.0"`
	FeedbackLink             *FeedbackLink       `json:"feedbackLink,omitempty" yaml:"feedbackLink"`
	SharingRecipientsPerPage int                 `json:"sharingRecipientsPerPage,omitempty" yaml:"sharingRecipientsPerPage" env:"WEB_OPTION_SHARING_RECIPIENTS_PER_PAGE" desc:"Sets the number of users shown as recipients in the dropdown menu when sharing resources." introductionVersion:"pre5.0"`
	Sidebar                  Sidebar             `json:"sidebar" yaml:"sidebar"`
	RunningOnEOS             bool                `json:"runningOnEos,omitempty" yaml:"runningOnEos" env:"WEB_OPTION_RUNNING_ON_EOS" desc:"Set this option to 'true' if running on an EOS storage backend (https://eos-web.web.cern.ch/eos-web/) to enable its specific features. Defaults to 'false'." introductionVersion:"pre5.0"`
	CernFeatures             bool                `json:"cernFeatures,omitempty" yaml:"cernFeatures"`
	HoverableQuickActions    bool                `json:"hoverableQuickActions,omitempty" yaml:"hoverableQuickActions" env:"WEB_OPTION_HOVERABLE_QUICK_ACTIONS" desc:"Set this option to 'true' to hide quick actions (buttons appearing on file rows) and only show them when the user hovers over the row with his mouse. Defaults to 'false'." introductionVersion:"pre5.0"`
	Routing                  Routing             `json:"routing" yaml:"routing"`
	Upload                   *Upload             `json:"upload,omitempty" yaml:"upload"`
	Editor                   *Editor             `json:"editor,omitempty" yaml:"editor"`
	ContextHelpersReadMore   bool                `json:"contextHelpersReadMore,omitempty" yaml:"contextHelpersReadMore" env:"WEB_OPTION_CONTEXTHELPERS_READ_MORE" desc:"Specifies whether the 'Read more' link should be displayed or not." introductionVersion:"pre5.0"`
	LogoutURL                string              `json:"logoutUrl,omitempty" yaml:"logoutUrl" env:"WEB_OPTION_LOGOUT_URL" desc:"Adds a link to the user's profile page to point him to an external page, where he can manage his session and devices. This is helpful when an external IdP is used. This option is disabled by default." introductionVersion:"pre5.0"`
	LoginURL                 string              `json:"loginUrl,omitempty" yaml:"loginUrl" env:"WEB_OPTION_LOGIN_URL" desc:"Specifies the target URL to the login page. This is helpful when an external IdP is used. This option is disabled by default. Example URL like: https://www.myidp.com/login." introductionVersion:"5.0"`
	OpenLinksWithDefaultApp  bool                `json:"openLinksWithDefaultApp,omitempty" yaml:"openLinksWithDefaultApp" env:"WEB_OPTION_OPEN_LINKS_WITH_DEFAULT_APP" desc:"Specifies whether single file link shares should be opened with the default app or not. If not opened by the default app, the Web UI just displays the file details." introductionVersion:"pre5.0"`
	TokenStorageLocal        bool                `json:"tokenStorageLocal" yaml:"tokenStorageLocal" env:"WEB_OPTION_TOKEN_STORAGE_LOCAL" desc:"Specifies whether the access token will be stored in the local storage when set to 'true' or in the session storage when set to 'false'. If stored in the local storage, login state will be persisted across multiple browser tabs, means no additional logins are required." introductionVersion:"pre5.0"`
	DisabledExtensions       []string            `json:"disabledExtensions,omitempty" yaml:"disabledExtensions" env:"WEB_OPTION_DISABLED_EXTENSIONS" desc:"A list to disable specific Web extensions identified by their ID. The ID can e.g. be taken from the 'index.ts' file of the web extension. Example: 'com.github.owncloud.web.files.search,com.github.owncloud.web.files.print'. See the Environment Variable Types description for more details." introductionVersion:"5.0"`
	Embed                    *Embed              `json:"embed,omitempty" yaml:"embed"`
	UserListRequiresFilter   bool                `json:"userListRequiresFilter,omitempty" yaml:"userListRequiresFilter" env:"WEB_OPTION_USER_LIST_REQUIRES_FILTER" desc:"Defines whether one or more filters must be set in order to list users in the Web admin settings. Set this option to 'true' if running in an environment with a lot of users and listing all users could slow down performance. Defaults to 'false'." introductionVersion:"5.0"`
	ConcurrentRequests       *ConcurrentRequests `json:"concurrentRequests,omitempty" yaml:"concurrentRequests"`
}

// AccountEditLink are the AccountEditLink options
type AccountEditLink struct {
	Href string `json:"href,omitempty" yaml:"href" env:"WEB_OPTION_ACCOUNT_EDIT_LINK_HREF" desc:"Set a different target URL for the edit link. Make sure to prepend it with 'http(s)://'." introductionVersion:"pre5.0"`
}

// FeedbackLink are the feedback link options
type FeedbackLink struct {
	Href        string `json:"href,omitempty" yaml:"href" env:"WEB_OPTION_FEEDBACKLINK_HREF" desc:"Set a target URL for the feedback link. Make sure to prepend it with 'http(s)://'. Defaults to 'https://owncloud.com/web-design-feedback'." introductionVersion:"pre5.0"`
	AriaLabel   string `json:"ariaLabel,omitempty" yaml:"ariaLabel" env:"WEB_OPTION_FEEDBACKLINK_ARIALABEL" desc:"Since the feedback link only has an icon, a screen reader accessible label can be set. The text defaults to 'ownCloud feedback survey'." introductionVersion:"pre5.0"`
	Description string `json:"description,omitempty" yaml:"description" env:"WEB_OPTION_FEEDBACKLINK_DESCRIPTION" desc:"For feedbacks, provide any description you want to see as tooltip and as accessible description. Defaults to 'Provide your feedback: We'd like to improve the web design and would be happy to hear your feedback. Thank you! Your ownCloud team'." introductionVersion:"pre5.0"`
}

// Sidebar are the sidebar option
type Sidebar struct {
	Shares SidebarShares `json:"shares" yaml:"shares"`
	Exif   SidebarExif   `json:"exif" yaml:"exif"`
}

// SidebarShares are the options for the shares sidebar
type SidebarShares struct {
	ShowAllOnLoad bool `json:"showAllOnLoad" yaml:"showAllOnLoad" env:"WEB_OPTION_SIDEBAR_SHARES_SHOW_ALL_ON_LOAD" desc:"Sets the list of the (link) shares list in the sidebar to be initially expanded. Default is a collapsed state, only showing the first three shares." introductionVersion:"pre5.0"`
}

type SidebarExif struct {
	ShowLocation bool `json:"showLocation" yaml:"showLocation" env:"WEB_OPTION_SIDEBAR_EXIF_SHOW_LOCATION" desc:"Shows the location data in the EXIF sidebar panel. Default is the location data being shown." introductionVersion:"6.0"`
}

// Routing are the routing options
type Routing struct {
	IDBased bool `json:"idBased" yaml:"idBased" env:"WEB_OPTION_ROUTING_ID_BASED" desc:"Enable or disable fileIds being added to the URL. Defaults to 'true', because otherwise spaces with name clashes cannot be resolved correctly. Note: Only disable this if you can guarantee on the server side, that spaces of the same namespace cannot have name clashes." introductionVersion:"pre5.0"`
}

// Upload are the upload options
type Upload struct {
	XHR          XHR    `json:"xhr,omitempty" yaml:"xhr"`
	CompanionURL string `json:"companionUrl,omitempty" yaml:"companionUrl" env:"WEB_OPTION_UPLOAD_COMPANION_URL" desc:"Sets the URL of Companion which is a service provided by Uppy to import files from external cloud providers. See https://uppy.io/docs/companion/ for instructions on how to set up Companion. This feature is disabled as long as no URL is given." introductionVersion:"pre5.0"`
}

// XHR are the XHR options
type XHR struct {
	Timeout int `json:"timeout,omitempty" yaml:"timeout" env:"WEB_OPTION_UPLOAD_XHR_TIMEOUT" desc:"Specifies the timeout for XHR uploads in milliseconds." introductionVersion:"pre5.0"`
}

// Editor are the web editor options
type Editor struct {
	AutosaveEnabled  bool `json:"autosaveEnabled,omitempty" yaml:"autosaveEnabled" env:"WEB_OPTION_EDITOR_AUTOSAVE_ENABLED" desc:"Specifies if the autosave for the editor apps is enabled." introductionVersion:"pre5.0"`
	AutosaveInterval int  `json:"autosaveInterval,omitempty" yaml:"autosaveInterval" env:"WEB_OPTION_EDITOR_AUTOSAVE_INTERVAL" desc:"Specifies the time interval for the autosave of editor apps in seconds. Has no effect when WEB_OPTION_EDITOR_AUTOSAVE_ENABLED is set to 'false'." introductionVersion:"pre5.0"`
}

// Embed are the Embed options
type Embed struct {
	Enabled                      string `json:"enabled,omitempty" yaml:"enabled" env:"WEB_OPTION_EMBED_ENABLED" desc:"Defines whether Web should be running in 'embed' mode. Setting this to 'true' will enable a stripped down version of Web with reduced functionality used to integrate Web into other applications like via iFrame. Setting it to 'false' or not setting it (default) will run Web as usual with all functionality enabled. See the text description for more details." introductionVersion:"5.0"`
	Target                       string `json:"target,omitempty" yaml:"target" env:"WEB_OPTION_EMBED_TARGET" desc:"Defines how Web is being integrated when running in 'embed' mode. Currently, the only supported options are '' (empty) and 'location'. With '' which is the default, Web will run regular as defined via the 'embed.enabled' config option. With 'location', Web will run embedded as location picker. Resource selection will be disabled and the selected resources array always includes the current folder as the only item. See the text description for more details." introductionVersion:"5.0"`
	MessagesOrigin               string `json:"messagesOrigin,omitempty" yaml:"messagesOrigin" env:"WEB_OPTION_EMBED_MESSAGES_ORIGIN" desc:"Defines a URL under which Web can be integrated via iFrame in 'embed' mode. Note that setting this is mandatory when running Web in 'embed' mode. Use '*' as value to allow running the iFrame under any URL, although this is not recommended for security reasons. See the text description for more details." introductionVersion:"5.0"`
	DelegateAuthentication       string `json:"delegateAuthentication,omitempty" yaml:"delegateAuthentication" env:"WEB_OPTION_EMBED_DELEGATE_AUTHENTICATION" desc:"Defines whether Web should require authentication to be done by the parent application when running in 'embed' mode. If set to 'true' Web will not try to authenticate the user on its own but will require an access token coming from the parent application. Defaults to being unset." introductionVersion:"5.0"`
	DelegateAuthenticationOrigin string `json:"delegateAuthenticationOrigin,omitempty" yaml:"delegateAuthenticationOrigin" env:"WEB_OPTION_EMBED_DELEGATE_AUTHENTICATION_ORIGIN" desc:"Defines the host to validate the message event origin against when running Web in 'embed' mode with delegated authentication. Defaults to event message origin validation being omitted, which is only recommended for development setups." introductionVersion:"5.0"`
}

// ConcurrentRequests are the ConcurrentRequests options
type ConcurrentRequests struct {
	ResourceBatchActions int                       `json:"resourceBatchActions,omitempty" yaml:"resourceBatchActions" env:"WEB_OPTION_CONCURRENT_REQUESTS_RESOURCE_BATCH_ACTIONS" desc:"Defines the maximum number of concurrent requests per file/folder/space batch action. Defaults to 4." introductionVersion:"5.0"`
	SSE                  int                       `json:"sse,omitempty" yaml:"sse" env:"WEB_OPTION_CONCURRENT_REQUESTS_SSE" desc:"Defines the maximum number of concurrent requests in SSE event handlers. Defaults to 4." introductionVersion:"5.0"`
	Shares               *ConcurrentRequestsShares `json:"shares,omitempty" yaml:"shares"`
}

// ConcurrentRequestsShares are the Shares options inside the ConcurrentRequests options
type ConcurrentRequestsShares struct {
	Create int `json:"create,omitempty" yaml:"create" env:"WEB_OPTION_CONCURRENT_REQUESTS_SHARES_CREATE" desc:"Defines the maximum number of concurrent requests per sharing invite batch. Defaults to 4." introductionVersion:"5.0"`
	List   int `json:"list,omitempty" yaml:"list" env:"WEB_OPTION_CONCURRENT_REQUESTS_SHARES_LIST" desc:"Defines the maximum number of concurrent requests when loading individual share information inside listings. Defaults to 2." introductionVersion:"5.0"`
}
