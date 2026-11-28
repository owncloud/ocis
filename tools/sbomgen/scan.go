package main

// Dependency discovery without Trivy: Go modules come from the Go toolchain, npm packages from
// pnpm-lock.yaml. The result is an SPDX 2.3 document that main.go enriches and converts.

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/spdx/tools-golang/spdx/v2/common"
	"github.com/spdx/tools-golang/spdx/v2/v2_3"
	"gopkg.in/yaml.v3"
)

type stringList []string

func (s *stringList) String() string     { return strings.Join(*s, ",") }
func (s *stringList) Set(v string) error { *s = append(*s, v); return nil }

type scanner struct {
	doc      *v2_3.Document
	root     common.ElementID
	edges    map[string]bool // "from>to" element ids
	goZips   map[string]string
	spdxIDs  map[string]bool
	licenses map[string]fileFacts // module dir -> license and copyright (cache)
	npm      []npmRef
}

// fileFacts is what the license files of one package say.
type fileFacts struct{ License, Copyright string }

func newScanner(version string, spdxIDs map[string]bool) *scanner {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	check(err)
	root := common.ElementID("Package-root")
	doc := &v2_3.Document{
		SPDXVersion:       "SPDX-2.3",
		DataLicense:       "CC0-1.0",
		SPDXIdentifier:    "DOCUMENT",
		DocumentName:      "ocis " + version,
		DocumentNamespace: fmt.Sprintf("https://spdx.org/spdxdocs/ocis-%s-%x", version, uuid),
		CreationInfo: &v2_3.CreationInfo{
			Created:  time.Now().UTC().Format("2006-01-02T15:04:05Z"),
			Creators: []common.Creator{{CreatorType: "Tool", Creator: "sbomgen"}},
		},
	}
	doc.Packages = []*v2_3.Package{{
		PackageName: "ocis", PackageSPDXIdentifier: root, PackageDownloadLocation: noAssertion,
		PrimaryPackagePurpose: "SOURCE", PackageLicenseDeclared: noAssertion, PackageLicenseConcluded: noAssertion,
	}}
	doc.Relationships = []*v2_3.Relationship{{
		RefA: common.DocElementID{ElementRefID: "DOCUMENT"}, RefB: common.DocElementID{ElementRefID: root}, Relationship: "DESCRIBES"}}
	return &scanner{doc: doc, root: root, edges: map[string]bool{}, goZips: map[string]string{},
		spdxIDs: spdxIDs, licenses: map[string]fileFacts{}}
}

func (s *scanner) edge(from, to common.ElementID) {
	key := string(from) + ">" + string(to)
	if from == to || s.edges[key] {
		return
	}
	s.edges[key] = true
	s.doc.Relationships = append(s.doc.Relationships, &v2_3.Relationship{
		RefA: common.DocElementID{ElementRefID: from}, RefB: common.DocElementID{ElementRefID: to}, Relationship: "DEPENDS_ON"})
}

func (s *scanner) addPackage(id common.ElementID, name, version, purl, license string) *v2_3.Package {
	p := &v2_3.Package{
		PackageName: name, PackageSPDXIdentifier: id, PackageVersion: version,
		PackageDownloadLocation: noAssertion, PrimaryPackagePurpose: "LIBRARY", PackageCopyrightText: noAssertion,
		PackageLicenseDeclared: license, PackageLicenseConcluded: noAssertion,
		PackageExternalReferences: []*v2_3.PackageExternalReference{
			{Category: "PACKAGE-MANAGER", RefType: "purl", Locator: purl}},
	}
	s.doc.Packages = append(s.doc.Packages, p)
	return p
}

func shortID(prefix, key string) common.ElementID {
	h := sha1.Sum([]byte(key))
	return common.ElementID(prefix + "-" + hex.EncodeToString(h[:])[:16])
}

// ---- Go ------------------------------------------------------------------------------------------

type goDownload struct{ Path, Version, Zip, Dir, Error string }

func goCmd(dir string, args ...string) []byte {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod", "GOTOOLCHAIN=local")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		fatal("go %s in %s: %v\n%s", strings.Join(args, " "), dir, err, stderr.String())
	}
	return out
}

// addGoModule adds the module in dir (its dependency graph, downloaded zips and licenses).
func (s *scanner) addGoModule(dir string) {
	selected := map[string]string{} // module path -> selected version
	pkgs := map[string]*v2_3.Package{}
	dec := json.NewDecoder(bytes.NewReader(goCmd(dir, "mod", "download", "-json")))
	for dec.More() {
		var m goDownload
		check(dec.Decode(&m))
		if m.Error != "" {
			fatal("go mod download %s@%s: %s", m.Path, m.Version, m.Error)
		}
		selected[m.Path] = m.Version
		key := m.Path + "@" + m.Version
		if _, dup := pkgs[key]; dup {
			continue
		}
		s.goZips[key] = m.Zip
		id := shortID("Package-go", key)
		if p := s.find(id); p != nil {
			pkgs[key] = p // already added by another module of the repository
			continue
		}
		facts := s.goLicenseFiles(m.Dir)
		p := s.addPackage(id, m.Path, m.Version, "pkg:golang/"+strings.ToLower(m.Path)+"@"+m.Version, facts.License)
		p.PackageCopyrightText = facts.Copyright
		pkgs[key] = p
	}
	// `go mod graph` lines: "<module>[@version] <dependency>@version"; edges to a version that MVS did
	// not select are edges to the selected one
	for _, line := range strings.Split(string(goCmd(dir, "mod", "graph")), "\n") {
		a, b, ok := strings.Cut(strings.TrimSpace(line), " ")
		if !ok {
			continue
		}
		bPath, _, _ := strings.Cut(b, "@")
		bVer, bSelected := selected[bPath]
		if !bSelected {
			continue // go, toolchain, modules that are not part of the build list
		}
		to := pkgs[bPath+"@"+bVer].PackageSPDXIdentifier
		if !strings.Contains(a, "@") { // the main module: a direct dependency of the project
			s.edge(s.root, to)
			continue
		}
		aPath, aVer, _ := strings.Cut(a, "@")
		if selected[aPath] == aVer {
			s.edge(pkgs[a].PackageSPDXIdentifier, to)
		}
	}
}

func (s *scanner) find(id common.ElementID) *v2_3.Package {
	for _, p := range s.doc.Packages {
		if p.PackageSPDXIdentifier == id {
			return p
		}
	}
	return nil
}

// goLicenseFiles reads the license files in the root of a module directory.
func (s *scanner) goLicenseFiles(dir string) fileFacts {
	if f, ok := s.licenses[dir]; ok {
		return f
	}
	ids := map[string]bool{}
	var copyrights [][]string
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		isLicense, isNotice := licenseFileName.MatchString(e.Name()), noticeFileName.MatchString(e.Name())
		if e.IsDir() || !(isLicense || isNotice) {
			continue
		}
		text, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		if isLicense {
			for _, id := range classify(text) {
				ids[id] = true
			}
		}
		copyrights = append(copyrights, copyrightLines(string(text)))
	}
	f := fileFacts{License: joinAnd(ids), Copyright: mergeCopyright(copyrights...)}
	s.licenses[dir] = f
	return f
}

// ---- npm (pnpm-lock.yaml, lockfile v9) ----------------------------------------------------------------

type pnpmDep struct{ Specifier, Version string }

type pnpmLockFile struct {
	Importers map[string]struct {
		Dependencies         map[string]pnpmDep `yaml:"dependencies"`
		OptionalDependencies map[string]pnpmDep `yaml:"optionalDependencies"`
		DevDependencies      map[string]pnpmDep `yaml:"devDependencies"`
	} `yaml:"importers"`
	Snapshots map[string]struct {
		Dependencies         map[string]string `yaml:"dependencies"`
		OptionalDependencies map[string]string `yaml:"optionalDependencies"`
	} `yaml:"snapshots"`
}

// identity strips the peer-dependency suffix: "a@1.0.0(b@2.0.0)" -> "a@1.0.0"
func identity(key string) string {
	if i := strings.Index(key, "("); i >= 0 {
		return key[:i]
	}
	return key
}

func splitNameVersion(id string) (name, version string) {
	i := strings.LastIndex(id, "@")
	return id[:i], id[i+1:]
}

// snapshotKey resolves a dependency entry (version, or "realname@version" for aliases) to a snapshot key.
func snapshotKey(name, value string) (string, bool) {
	if strings.HasPrefix(value, "link:") || strings.HasPrefix(value, "file:") || strings.HasPrefix(value, "workspace:") {
		return "", false
	}
	if strings.LastIndex(identity(value), "@") > 0 { // alias: "string-width-cjs: string-width@4.2.3"
		return value, true
	}
	return name + "@" + value, true
}

func (s *scanner) addPnpmLock(path string, includeDev bool) {
	b, err := os.ReadFile(path)
	check(err)
	var lock pnpmLockFile
	check(yaml.Unmarshal(b, &lock))

	var queue []string
	visited := map[string]bool{}
	pkgs := map[string]*v2_3.Package{} // identity -> package
	pkgFor := func(key string) *v2_3.Package {
		id := identity(key)
		if p, ok := pkgs[id]; ok {
			return p
		}
		name, version := splitNameVersion(id)
		purl := "pkg:npm/" + strings.Replace(strings.ToLower(name), "@", "%40", 1) + "@" + version
		p := s.addPackage(shortID("Package-npm", id), name, version, purl, noAssertion)
		s.npm = append(s.npm, npmRef{Name: name, Version: version})
		pkgs[id] = p
		return p
	}
	enqueue := func(key string) {
		if !visited[key] {
			visited[key] = true
			queue = append(queue, key)
		}
	}

	importers := make([]string, 0, len(lock.Importers))
	for name := range lock.Importers {
		importers = append(importers, name)
	}
	sort.Strings(importers)
	for _, name := range importers {
		imp := lock.Importers[name]
		groups := []map[string]pnpmDep{imp.Dependencies, imp.OptionalDependencies}
		if includeDev {
			groups = append(groups, imp.DevDependencies)
		}
		for _, g := range groups {
			for dep, d := range g {
				if key, ok := snapshotKey(dep, d.Version); ok {
					s.edge(s.root, pkgFor(key).PackageSPDXIdentifier)
					enqueue(key)
				}
			}
		}
	}
	for len(queue) > 0 {
		key := queue[0]
		queue = queue[1:]
		snap, ok := lock.Snapshots[key]
		if !ok {
			fatal("%s: no snapshot for %s", path, key)
		}
		from := pkgFor(key).PackageSPDXIdentifier
		for _, deps := range []map[string]string{snap.Dependencies, snap.OptionalDependencies} {
			for dep, value := range deps {
				if dk, ok := snapshotKey(dep, value); ok {
					s.edge(from, pkgFor(dk).PackageSPDXIdentifier)
					enqueue(dk)
				}
			}
		}
	}
}

// scan builds the SPDX 2.3 document of the project from its Go modules and pnpm lockfile.
func scan(goModules []string, pnpmLock string, includeDev bool, version string, spdxIDs map[string]bool) (*v2_3.Document, map[string]string, []npmRef) {
	s := newScanner(version, spdxIDs)
	for _, dir := range goModules {
		s.addGoModule(dir)
	}
	s.addPnpmLock(pnpmLock, includeDev)
	return s.doc, s.goZips, s.npm
}
