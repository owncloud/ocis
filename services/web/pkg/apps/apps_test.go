package apps

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/ocis-pkg/log"
	"github.com/owncloud/ocis/v2/services/web/pkg/config"
)

func TestApplication_ToExternal(t *testing.T) {
	g := gomega.NewWithT(t)
	app := Application{
		ID:         "app",
		Entrypoint: "entrypoint.js",
		Config: map[string]interface{}{
			"foo": "bar",
		},
	}

	externalApp := app.ToExternal("path")

	g.Expect(externalApp.ID).To(gomega.Equal("app"))
	g.Expect(externalApp.Path).To(gomega.Equal("path/entrypoint.js"))
	g.Expect(externalApp.Config).To(gomega.Equal(app.Config))
}

func TestBuild(t *testing.T) {
	g := gomega.NewWithT(t)
	dir := &fstest.MapFile{
		Mode: fs.ModeDir,
	}

	_, err := build(fstest.MapFS{
		"app": &fstest.MapFile{},
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(ErrInvalidApp))

	_, err = build(fstest.MapFS{
		"app": dir,
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(ErrMissingManifest))

	_, err = build(fstest.MapFS{
		"app":               dir,
		"app/manifest.json": dir,
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(ErrInvalidManifest))

	_, err = build(fstest.MapFS{
		"app": dir,
		"app/manifest.json": &fstest.MapFile{
			Data: []byte("{}"),
		},
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(ErrInvalidManifest))

	_, err = build(fstest.MapFS{
		"app":               dir,
		"app/entrypoint.js": &fstest.MapFile{},
		"app/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app", "entrypoint":"entrypoint.js"}`),
		},
	}, "app", map[string]any{})
	g.Expect(err).ToNot(gomega.HaveOccurred())

	_, err = build(fstest.MapFS{
		"app":               dir,
		"app/entrypoint.js": dir,
		"app/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app", "entrypoint":"entrypoint.js"}`),
		},
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(ErrEntrypointDoesNotExist))

	_, err = build(fstest.MapFS{
		"app": dir,
		"app/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app", "entrypoint":"entrypoint.js"}`),
		},
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(ErrEntrypointDoesNotExist))

	application, err := build(fstest.MapFS{
		"app":               dir,
		"app/entrypoint.js": &fstest.MapFile{},
		"app/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app", "entrypoint":"entrypoint.js", "config": {"foo": "1", "bar": "2"}}`),
		},
	}, "app", map[string]any{"foo": "overwritten-1", "baz": "injected-1"})
	g.Expect(err).ToNot(gomega.HaveOccurred())

	g.Expect(application.Entrypoint).To(gomega.Equal("app/entrypoint.js"))
	g.Expect(application.Config).To(gomega.Equal(map[string]interface{}{
		"foo": "overwritten-1", "baz": "injected-1", "bar": "2",
	}))
}

func TestList(t *testing.T) {
	g := gomega.NewWithT(t)

	applications := List(log.NopLogger(), map[string]config.App{})
	g.Expect(applications).To(gomega.BeEmpty())

	applications = List(log.NopLogger(), map[string]config.App{}, nil)
	g.Expect(applications).To(gomega.BeEmpty())

	applications = List(log.NopLogger(), map[string]config.App{}, fstest.MapFS{})
	g.Expect(applications).To(gomega.BeEmpty())

	dir := &fstest.MapFile{
		Mode: fs.ModeDir,
	}

	applications = List(log.NopLogger(), map[string]config.App{
		"app": {
			Disabled: true,
		},
	}, fstest.MapFS{
		"app": dir,
	})
	g.Expect(applications).To(gomega.BeEmpty())

	applications = List(log.NopLogger(), map[string]config.App{
		"app": {},
	}, fstest.MapFS{
		"app": dir,
	})
	g.Expect(applications).To(gomega.BeEmpty())

	applications = List(log.NopLogger(), map[string]config.App{
		"app-3": {
			Config: map[string]any{
				"foo": "local conf 1",
				"bar": "local conf 2",
			},
		},
	}, fstest.MapFS{
		"app-1":               dir,
		"app-1/entrypoint.js": &fstest.MapFile{},
		"app-1/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-1", "entrypoint":"entrypoint.js", "config": {"foo": "fs1"}}`),
		},
		"app-2":               dir,
		"app-2/entrypoint.js": &fstest.MapFile{},
		"app-2/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-2", "entrypoint":"entrypoint.js", "config": {"foo": "fs1"}}`),
		},
	}, fstest.MapFS{
		"app-1":               dir,
		"app-1/entrypoint.js": &fstest.MapFile{},
		"app-1/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-1", "entrypoint":"entrypoint.js", "config": {"foo": "fs2"}}`),
		},
		"app-3":               dir,
		"app-3/entrypoint.js": &fstest.MapFile{},
		"app-3/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-3", "entrypoint":"entrypoint.js", "config": {"foo": "fs2"}}`),
		},
	})
	g.Expect(len(applications)).To(gomega.Equal(3))

	for _, application := range applications {
		switch {
		case application.Entrypoint == "app-1/entrypoint.js":
			g.Expect(application.Config["foo"]).To(gomega.Equal("fs2"))
		case application.Entrypoint == "app-2/entrypoint.js":
			g.Expect(application.Config["foo"]).To(gomega.Equal("fs1"))
		case application.Entrypoint == "app-3/entrypoint.js":
			g.Expect(application.Config["foo"]).To(gomega.Equal("local conf 1"))
			g.Expect(application.Config["bar"]).To(gomega.Equal("local conf 2"))
		default:
			t.Fatalf("unexpected application %s", application.Entrypoint)
		}
	}

}
