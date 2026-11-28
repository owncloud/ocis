package main

// npm package facts that the lockfile does not carry: the declared license and the copyright
// statements. They come from the package's registry metadata and from the license files in its
// tarball. The tarball is checked against the integrity hash of the lockfile, so what is read is
// exactly what the lockfile pins.

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/spdx/tools-golang/spdx/v2/common"
)

// npmInfo is what is known about one npm package (and what the cache stores).
type npmInfo struct {
	Integrity string `json:"integrity"` // hex digest the info was derived from
	License   string `json:"license"`   // SPDX expression or NOASSERTION
	Source    string `json:"source"`    // where the license came from: metadata, license file or readme
	Copyright string `json:"copyright"` // statements or NOASSERTION
}

const maxTextFile = 1 << 20 // license and readme files are small; never read more than this

var readmeFileName = regexp.MustCompile(`(?i)^readme(\.md|\.markdown|\.txt)?$`)
var readmeLicenseHeading = regexp.MustCompile(`(?im)^#{1,6}[ \t]*licen[cs]e[ \t]*$`)

var httpClient = &http.Client{Timeout: 5 * time.Minute}

func getWithRetry(u string) (*http.Response, error) {
	var err error
	for attempt := 0; attempt < 4; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(attempt*attempt) * time.Second)
		}
		var resp *http.Response
		resp, err = httpClient.Get(u)
		if err == nil && resp.StatusCode == 200 {
			return resp, nil
		}
		if err == nil {
			err = fmt.Errorf("%s: %s", u, resp.Status)
			resp.Body.Close()
			if resp.StatusCode == 404 {
				break
			}
		}
	}
	return nil, err
}

type npmMetadata struct {
	License  any `json:"license"`
	Licenses []struct {
		Type string `json:"type"`
	} `json:"licenses"` // legacy form: several entries mean "any of"
	Dist struct {
		Tarball string `json:"tarball"`
	} `json:"dist"`
}

// declared returns the license expression declared in the metadata ("" if there is none).
func (m npmMetadata) declared() string {
	switch l := m.License.(type) {
	case string:
		return l
	case map[string]any:
		if s, _ := l["type"].(string); s != "" {
			return s
		}
	}
	var types []string
	for _, l := range m.Licenses {
		if l.Type != "" {
			types = append(types, l.Type)
		}
	}
	return strings.Join(types, " OR ")
}

func newHash(algo common.ChecksumAlgorithm) hash.Hash {
	switch algo {
	case common.SHA1:
		return sha1.New()
	case common.SHA256:
		return sha256.New()
	case common.SHA384:
		return sha512.New384()
	}
	return sha512.New()
}

// fetchNPM reads the registry metadata and the tarball of name@version.
func fetchNPM(name, version string, pinned *common.Checksum, spdxIDs map[string]bool) (npmInfo, error) {
	info := npmInfo{License: noAssertion, Copyright: noAssertion}
	resp, err := getWithRetry("https://registry.npmjs.org/" + strings.Replace(url.PathEscape(name), "%40", "@", 1) + "/" + version)
	if err != nil {
		return info, err
	}
	var meta npmMetadata
	err = json.NewDecoder(resp.Body).Decode(&meta)
	resp.Body.Close()
	if err != nil {
		return info, fmt.Errorf("%s@%s metadata: %v", name, version, err)
	}
	if l := meta.declared(); l != "" && validExpression(l, spdxIDs) {
		info.License, info.Source = l, "metadata"
	}
	if meta.Dist.Tarball == "" {
		return info, fmt.Errorf("%s@%s: no tarball in metadata", name, version)
	}

	tgz, err := getWithRetry(meta.Dist.Tarball)
	if err != nil {
		return info, err
	}
	defer tgz.Body.Close()
	var h hash.Hash
	var body io.Reader = tgz.Body
	if pinned != nil {
		h = newHash(pinned.Algorithm)
		body = io.TeeReader(tgz.Body, h)
	}
	gz, err := gzip.NewReader(body)
	if err != nil {
		return info, fmt.Errorf("%s@%s tarball: %v", name, version, err)
	}
	ids := map[string]bool{}
	var copyrights [][]string
	var readme string
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return info, fmt.Errorf("%s@%s tarball: %v", name, version, err)
		}
		// only files in the package root: "package/LICENSE", not "package/node_modules/x/LICENSE"
		dir, file := path.Split(strings.TrimPrefix(hdr.Name, "./"))
		if hdr.Typeflag != tar.TypeReg || strings.Count(strings.TrimSuffix(dir, "/"), "/") != 0 || dir == "" {
			continue
		}
		switch {
		case licenseFileName.MatchString(file), noticeFileName.MatchString(file):
			text, _ := io.ReadAll(io.LimitReader(tr, maxTextFile))
			if licenseFileName.MatchString(file) {
				for _, id := range classify(text) {
					ids[id] = true
				}
			}
			copyrights = append(copyrights, copyrightLines(string(text)))
		case readmeFileName.MatchString(file):
			text, _ := io.ReadAll(io.LimitReader(tr, maxTextFile))
			readme = string(text)
		}
	}
	_, _ = io.Copy(io.Discard, gz)
	_, _ = io.Copy(io.Discard, body) // the whole stream must be read for the hash
	if h != nil && hex.EncodeToString(h.Sum(nil)) != pinned.Value {
		return info, fmt.Errorf("%s@%s: tarball does not match the integrity hash of the lockfile", name, version)
	}
	if pinned != nil {
		info.Integrity = pinned.Value
	}

	if info.License == noAssertion {
		if l := joinAnd(ids); validExpression(l, spdxIDs) && l != noAssertion {
			info.License, info.Source = l, "license file"
		} else if l := readmeLicense(readme, spdxIDs); l != "" {
			info.License, info.Source = l, "readme"
		}
	}
	info.Copyright = mergeCopyright(copyrights...)
	return info, nil
}

// readmeLicense returns the license a README names in its "License" section, if that is a single SPDX id.
func readmeLicense(readme string, spdxIDs map[string]bool) string {
	loc := readmeLicenseHeading.FindStringIndex(readme)
	if loc == nil {
		return ""
	}
	for _, line := range strings.Split(readme[loc[1]:], "\n") {
		line = strings.Trim(strings.TrimSpace(line), "*_`.")
		if line == "" {
			continue
		}
		if spdxIDs[line] {
			return line
		}
		return ""
	}
	return ""
}

// loadNPMInfo returns license and copyright information for every npm package, using the cache file
// for packages whose pinned hash is unchanged.
func loadNPMInfo(pkgs []npmRef, cachePath string, spdxIDs map[string]bool) map[string]npmInfo {
	cache := map[string]npmInfo{}
	if cachePath != "" {
		if b, err := os.ReadFile(cachePath); err == nil {
			_ = json.Unmarshal(b, &cache)
		}
	}
	result := map[string]npmInfo{}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var firstErr error
	sem := make(chan struct{}, 8)
	fetched := 0
	for _, p := range pkgs {
		key := p.Name + "@" + p.Version
		if c, ok := cache[key]; ok && p.Pinned != nil && c.Integrity == p.Pinned.Value {
			result[key] = c
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(p npmRef, key string) {
			defer wg.Done()
			defer func() { <-sem }()
			info, err := fetchNPM(p.Name, p.Version, p.Pinned, spdxIDs)
			mu.Lock()
			defer mu.Unlock()
			if err != nil && firstErr == nil {
				firstErr = err
			}
			result[key] = info
			fetched++
		}(p, key)
	}
	wg.Wait()
	if firstErr != nil {
		fatal("%v", firstErr)
	}
	fmt.Fprintf(os.Stderr, "npm: %d packages, %d fetched, %d from cache\n", len(pkgs), fetched, len(pkgs)-fetched)
	if cachePath != "" {
		b, err := json.MarshalIndent(result, "", " ")
		check(err)
		check(os.WriteFile(cachePath, b, 0o644))
	}
	return result
}

type npmRef struct {
	Name, Version string
	Pinned        *common.Checksum // integrity hash from the lockfile
}
