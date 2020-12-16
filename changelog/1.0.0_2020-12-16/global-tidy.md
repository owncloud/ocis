Enhancement: Tidy dependencies

Methodology:

```
go-modules() {
  find . \( -name vendor -o -name '[._].*' -o -name node_modules \) -prune -o -name go.mod -print | sed 's:/go.mod$::'
}
```

```
for m in $(go-modules); do (cd $m && go mod tidy); done
```

https://github.com/owncloud/ocis/pull/845
