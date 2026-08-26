package theme

import (
	"maps"
	"path"

	"github.com/owncloud/ocis/v2/ocis-pkg/capabilities"
	"github.com/owncloud/ocis/v2/services/graph/pkg/unifiedrole"
)

var (
	_brandingRoot  = "_branding"
	_themeFileName = "theme.json"

	// _brandingLogoKeys are the clients.web.defaults.logo sub-keys the admin
	// logo upload endpoint controls. LogoUpload, LogoReset, and
	// applyLogoOverride all derive their keys from this single list so they
	// can't drift apart — anything added here automatically also overrides a
	// theme's own per-variant logo.
	_brandingLogoKeys = []string{"topbar", "login"}
)

// themeDefaults contains the default values for the theme.
// all rendered themes get the default values from here.
var themeDefaults = KV{
	"common": KV{
		"shareRoles": KV{
			unifiedrole.UnifiedRoleViewerID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleViewerID),
				"iconName": "eye",
			},
			unifiedrole.UnifiedRoleViewerListGrantsID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleViewerListGrantsID),
				"iconName": "eye",
			},
			unifiedrole.UnifiedRoleSpaceViewerID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleSpaceViewerID),
				"iconName": "eye",
			},
			unifiedrole.UnifiedRoleFileEditorID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleFileEditorID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleFileEditorListGrantsID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleFileEditorListGrantsID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleFileEditorListGrantsWithVersionsID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleFileEditorListGrantsWithVersionsID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleEditorID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleEditorID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleEditorListGrantsID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleEditorListGrantsID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleEditorListGrantsWithVersionsID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleEditorListGrantsWithVersionsID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleSpaceEditorID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleSpaceEditorID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleSpaceEditorWithoutVersionsID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleSpaceEditorWithoutVersionsID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleSpaceEditorWithoutTrashbinID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleSpaceEditorWithoutTrashbinID),
				"iconName": "pencil",
			},
			unifiedrole.UnifiedRoleManagerID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleManagerID),
				"iconName": "user-star",
			},
			unifiedrole.UnifiedRoleEditorLiteID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleEditorLiteID),
				"iconName": "upload",
			},
			unifiedrole.UnifiedRoleSecureViewerID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleSecureViewerID),
				"iconName": "shield",
			},
			unifiedrole.UnifiedRoleDeniedID: KV{
				"label":    unifiedrole.GetUnifiedRoleLabel(unifiedrole.UnifiedRoleDeniedID),
				"iconName": "stop-circle",
			},
		},
	},
}

// isFiletypePermitted checks if the given file extension is allowed.
func isFiletypePermitted(filename string, givenMime string) bool {
	// Check if we allow that extension and if the mediatype matches the extension
	extensionMime, ok := capabilities.Default().Theme.Logo.PermittedFileTypes[path.Ext(filename)]
	return ok && extensionMime == givenMime
}

// asMap returns v as a map[string]any, regardless of whether its concrete
// type is KV or map[string]any (both occur depending on whether a KV tree
// was built in code or decoded from JSON).
func asMap(v any) (map[string]any, bool) {
	switch t := v.(type) {
	case map[string]any:
		return t, true
	case KV:
		return map[string]any(t), true
	default:
		return nil, false
	}
}

// getPath looks up a dotted path of keys in a nested KV/map tree.
func getPath(v any, path ...string) (any, bool) {
	cur := v
	for _, p := range path {
		m, ok := asMap(cur)
		if !ok {
			return nil, false
		}
		cur, ok = m[p]
		if !ok {
			return nil, false
		}
	}
	return cur, true
}

// applyLogoOverride forces any admin-uploaded logo (branding) onto every
// theme variant's own logo entries. Without this, a theme that defines
// clients.web.themes[].logo itself would keep showing its own logo instead
// of the one the admin explicitly uploaded, since that per-variant value
// lives outside the clients.web.defaults.logo keys the upload patches.
//
// This mutates merged in place, which is only safe because clients.* is
// per-request data coming from baseTheme/brandingTheme, never from the
// shared, package-level themeDefaults tree that MergeKV also aliases in.
func applyLogoOverride(merged KV, branding KV) {
	overrides := map[string]any{}
	for _, key := range _brandingLogoKeys {
		if val, ok := getPath(branding, "clients", "web", "defaults", "logo", key); ok {
			overrides[key] = val
		}
	}
	if len(overrides) == 0 {
		return
	}

	themes, ok := getPath(merged, "clients", "web", "themes")
	if !ok {
		return
	}
	variants, ok := themes.([]any)
	if !ok {
		return
	}

	for _, item := range variants {
		variant, ok := asMap(item)
		if !ok {
			continue
		}
		logo, ok := asMap(variant["logo"])
		if !ok {
			continue
		}
		maps.Copy(logo, overrides)
	}
}
