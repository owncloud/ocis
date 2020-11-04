---
title: "Build the documentation"
date: 2020-07-27T08:39:38+00:00
weight: 99
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/development
geekdocFilePath: building-docs.md
---

## Buildling the documentation

Following steps can be applied for every oCIS extension repository.

### Setting up

- Install [hugo](https://gohugo.io/getting-started/installing/)
- Run `make docs`

### Viewing the documentation

To view the rendered docs in the browser run:
```bash
cd hugo
hugo -D server
```

Then open "http://localhost:1313/"

When making changes to the docs, run

```bash
rsync -ax --delete ../docs/ content/
```

in the hugo folder and the server will pick up the changes and reload the page automatically.

### Deploying the documentation

The documentation is automatically deployed from the master branch to https://owncloud.github.io

