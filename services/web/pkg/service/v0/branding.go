package svc

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"path"
	"path/filepath"

	permissionsapi "github.com/cs3org/go-cs3apis/cs3/permissions/v1beta1"
	rpc "github.com/cs3org/go-cs3apis/cs3/rpc/v1beta1"
	revactx "github.com/cs3org/reva/v2/pkg/ctx"
)

var (
	errInvalidThemeConfig       = errors.New("invalid themes config")
	_themesConfigPath           = filepath.FromSlash("themes/owncloud/theme.json")
	_allowedExtensionMediatypes = map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
	}
)

// UploadLogo implements the endpoint to upload a custom logo for the oCIS instance.
func (p Web) UploadLogo(w http.ResponseWriter, r *http.Request) {
	gatewayClient, err := p.gatewaySelector.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := revactx.ContextMustGetUser(r.Context())
	rsp, err := gatewayClient.CheckPermission(r.Context(), &permissionsapi.CheckPermissionRequest{
		Permission: "Logo.Write",
		SubjectRef: &permissionsapi.SubjectReference{
			Spec: &permissionsapi.SubjectReference_UserId{
				UserId: user.Id,
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rsp.Status.Code != rpc.Code_CODE_OK {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	file, fileHeader, err := r.FormFile("logo")
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			w.WriteHeader(http.StatusBadRequest)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	mediatype := fileHeader.Header.Get("Content-Type")
	if !allowedFiletype(fileHeader.Filename, mediatype) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	fp := filepath.Join("branding", filepath.Join("/", fileHeader.Filename))
	err = p.storeAsset(fp, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = p.updateLogoThemeConfig(fp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ResetLogo implements the endpoint to reset the instance logo.
// The config will be changed back to use the embedded logo asset.
func (p Web) ResetLogo(w http.ResponseWriter, r *http.Request) {
	gatewayClient, err := p.gatewaySelector.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := revactx.ContextMustGetUser(r.Context())
	rsp, err := gatewayClient.CheckPermission(r.Context(), &permissionsapi.CheckPermissionRequest{
		Permission: "Logo.Write",
		SubjectRef: &permissionsapi.SubjectReference{
			Spec: &permissionsapi.SubjectReference_UserId{
				UserId: user.Id,
			},
		},
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if rsp.Status.Code != rpc.Code_CODE_OK {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	f, err := p.fs.OpenEmbedded(_themesConfigPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	originalPath, err := p.getLogoPath(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := p.updateLogoThemeConfig(originalPath); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (p Web) storeAsset(name string, asset io.Reader) error {
	dst, err := p.fs.Create(name)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, asset)
	return err
}

func (p Web) getLogoPath(r io.Reader) (string, error) {
	// This decoding of the themes.json file is not optimal. If we need to decode it for other
	// usecases as well we should consider decoding to a struct.
	var m map[string]interface{}
	_ = json.NewDecoder(r).Decode(&m)

	webCfg, ok := m["web"].(map[string]interface{})
	if !ok {
		return "", errInvalidThemeConfig
	}

	defaultCfg, ok := webCfg["default"].(map[string]interface{})
	if !ok {
		return "", errInvalidThemeConfig
	}

	logoCfg, ok := defaultCfg["logo"].(map[string]interface{})
	if !ok {
		return "", errInvalidThemeConfig
	}

	logoPath, ok := logoCfg["login"].(string)
	if !ok {
		return "", errInvalidThemeConfig
	}

	return logoPath, nil
}

func (p Web) updateLogoThemeConfig(logoPath string) error {
	f, err := p.fs.Open(_themesConfigPath)
	if err == nil {
		defer f.Close()
	}

	// This decoding of the themes.json file is not optimal. If we need to decode it for other
	// usecases as well we should consider decoding to a struct.
	var m map[string]interface{}
	_ = json.NewDecoder(f).Decode(&m)

	// change logo in common part
	commonCfg, ok := m["common"].(map[string]interface{})
	if !ok {
		return errInvalidThemeConfig
	}
	commonCfg["logo"] = logoPath

	webCfg, ok := m["web"].(map[string]interface{})
	if !ok {
		return errInvalidThemeConfig
	}

	// iterate over all possible themes and replace logo
	for theme := range webCfg {
		themeCfg, ok := webCfg[theme].(map[string]interface{})
		if !ok {
			return errInvalidThemeConfig
		}

		logoCfg, ok := themeCfg["logo"].(map[string]interface{})
		if !ok {
			return errInvalidThemeConfig
		}

		logoCfg["login"] = logoPath
		logoCfg["topbar"] = logoPath
	}

	dst, err := p.fs.Create(_themesConfigPath)
	if err != nil {
		return err
	}

	return json.NewEncoder(dst).Encode(m)
}

func allowedFiletype(filename, mediatype string) bool {
	ext := path.Ext(filename)

	// Check if we allow that extension and if the mediatype matches the extension
	mt, ok := _allowedExtensionMediatypes[ext]
	return ok && mt == mediatype
}
