Bugfix: Add inotify-tools and bash packages to docker files

We need both packages to make posixfs work. Later, once the golang
package is fixed to not depend on bash any more, bash can be removed
again.

https://github.com/owncloud/ocis/pull/
