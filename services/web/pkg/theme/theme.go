package theme

import (
	"path"

	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

var (
	_brandingRoot         = "_branding"
	_themeFileName        = "theme.json"
	_allowedLogoFileTypes = map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}
)

// themeDefaults contains the default values for the theme.
// all rendered themes get the default values from here.
var themeDefaults = KV{
	"common": KV{
		"shareRoles": KV{
			unifiedrole.UnifiedRoleViewerID: KV{
				"name":     "UnifiedRoleViewer",
				"iconName": "eye",
			},
			unifiedrole.UnifiedRoleSpaceViewerID: KV{
				"name":     "UnifiedRoleFileEditor",
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleFileEditorID: KV{
				"name":     "UnifiedRoleEditor",
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleEditorID: KV{
				"name":     "UnifiedRoleSpaceEditor",
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleSpaceEditorID: KV{
				"name":     "UnifiedRoleManager",
				"iconName": "user-star",
			},
			unifiedrole.UnifiedRoleUploaderID: KV{
				"name":     "UnifiedRoleUploader",
				"iconName": "upload",
			},
			// TBD
			"aa97fe03-7980-45ac-9e50-b325749fd7e6": KV{
				"name":     "UnifiedRoleSecureView",
				"iconName": "shield",
			},
		},
	},
}

// isFiletypePermitted checks if the given file extension is allowed and if the given mediatype matches the extension
func isFiletypePermitted(allowed map[string]string, filename string, givenMime string) bool {
	// Check if we allow that extension and if the mediatype matches the extension
	extensionMime, ok := allowed[path.Ext(filename)]
	return ok && extensionMime == givenMime
}
