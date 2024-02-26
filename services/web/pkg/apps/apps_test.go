package apps_test

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/onsi/gomega"

	"github.com/owncloud/ocis/v2/services/web/pkg/apps"
	"github.com/owncloud/ocis/v2/services/web/pkg/config"
)

func TestApplication_ToExternal(t *testing.T) {
	g := gomega.NewWithT(t)
	app := apps.Application{
		ID:         "id",
		Entrypoint: "entrypoint.js",
		Config: map[string]interface{}{
			"foo": "bar",
		},
	}

	externalApp := app.ToExternal("path")

	g.Expect(externalApp.ID).To(gomega.Equal("id"))
	g.Expect(externalApp.Path).To(gomega.Equal("path/id/entrypoint.js"))
	g.Expect(externalApp.Config).To(gomega.Equal(app.Config))
}

func TestBuild(t *testing.T) {
	g := gomega.NewWithT(t)
	appContainer := &fstest.MapFile{
		Mode: fs.ModeDir,
	}

	_, err := apps.Build(fstest.MapFS{
		"app": &fstest.MapFile{},
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(apps.ErrInvalidApp))

	_, err = apps.Build(fstest.MapFS{
		"app": appContainer,
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(apps.ErrMissingManifest))

	_, err = apps.Build(fstest.MapFS{
		"app":               appContainer,
		"app/manifest.json": appContainer,
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(apps.ErrInvalidManifest))

	_, err = apps.Build(fstest.MapFS{
		"app": appContainer,
		"app/manifest.json": &fstest.MapFile{
			Data: []byte("{}"),
		},
	}, "app", map[string]any{})
	g.Expect(err).To(gomega.MatchError(apps.ErrInvalidManifest))

	_, err = apps.Build(fstest.MapFS{
		"app": appContainer,
		"app/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app", "entrypoint":"entrypoint.js"}`),
		},
	}, "app", map[string]any{})
	g.Expect(err).ToNot(gomega.HaveOccurred())

	application, err := apps.Build(fstest.MapFS{
		"app": appContainer,
		"app/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app", "entrypoint":"entrypoint.js", "config": {"foo": "1", "bar": "2"}}`),
		},
	}, "app", map[string]any{"foo": "overwritten-1", "baz": "injected-1"})
	g.Expect(err).ToNot(gomega.HaveOccurred())

	g.Expect(application.ID).To(gomega.Equal("app"))
	g.Expect(application.Entrypoint).To(gomega.Equal("entrypoint.js"))
	g.Expect(application.Config).To(gomega.Equal(map[string]interface{}{
		"foo": "overwritten-1", "baz": "injected-1", "bar": "2",
	}))
}

func TestList(t *testing.T) {
	g := gomega.NewWithT(t)

	applications := apps.List(map[string]config.App{})
	g.Expect(applications).To(gomega.BeEmpty())

	applications = apps.List(map[string]config.App{}, nil)
	g.Expect(applications).To(gomega.BeEmpty())

	applications = apps.List(map[string]config.App{}, fstest.MapFS{})
	g.Expect(applications).To(gomega.BeEmpty())

	appContainer := &fstest.MapFile{
		Mode: fs.ModeDir,
	}

	applications = apps.List(map[string]config.App{
		"app": {
			Disabled: true,
		},
	}, fstest.MapFS{
		"app": appContainer,
	})
	g.Expect(applications).To(gomega.BeEmpty())

	applications = apps.List(map[string]config.App{
		"app": {},
	}, fstest.MapFS{
		"app": appContainer,
	})
	g.Expect(applications).To(gomega.BeEmpty())

	applications = apps.List(map[string]config.App{
		"app-3": {
			Config: map[string]any{
				"foo": "local conf 1",
				"bar": "local conf 2",
			},
		},
	}, fstest.MapFS{
		"app-1": appContainer,
		"app-1/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-1", "entrypoint":"entrypoint.js", "config": {"foo": "fs1"}}`),
		},
		"app-2": appContainer,
		"app-2/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-2", "entrypoint":"entrypoint.js", "config": {"foo": "fs1"}}`),
		},
	}, fstest.MapFS{
		"app-1": appContainer,
		"app-1/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-1", "entrypoint":"entrypoint.js", "config": {"foo": "fs2"}}`),
		},
		"app-3": appContainer,
		"app-3/manifest.json": &fstest.MapFile{
			Data: []byte(`{"id":"app-3", "entrypoint":"entrypoint.js", "config": {"foo": "fs2"}}`),
		},
	})
	g.Expect(len(applications)).To(gomega.Equal(3))

	for _, application := range applications {
		switch {
		case application.ID == "app-1":
			g.Expect(application.Config["foo"]).To(gomega.Equal("fs2"))
		case application.ID == "app-2":
			g.Expect(application.Config["foo"]).To(gomega.Equal("fs1"))
		case application.ID == "app-3":
			g.Expect(application.Config["foo"]).To(gomega.Equal("local conf 1"))
			g.Expect(application.Config["bar"]).To(gomega.Equal("local conf 2"))
		default:
			t.Fatalf("unexpected application %s", application.ID)
		}
	}

}
