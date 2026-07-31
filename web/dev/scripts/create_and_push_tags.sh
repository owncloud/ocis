#!/bin/bash

# This script creates and pushes tags for the main app and all published packages.

APPS=("design-system" "eslint-config" "extension-sdk" "web-pkg" "web-client" "web-test-helpers")

cd "$(dirname "$0")/../.."
VERSION=$(node -p "require('./package.json').version")
TAG="v${VERSION}"

git tag -s -a "$TAG" -m "$TAG"
git push origin "$TAG"
echo "v$VERSION has been created and pushed"

for app in "${APPS[@]}"; do
	cd "./packages/$app"

	VERSION=$(node -p "require('./package.json').version")
	TAG="${app}-v${VERSION}"

	git tag -s -a "$TAG" -m "$TAG"
	git push origin "$TAG"

	echo "$app-v$VERSION has been created and pushed"
	cd "../.."
done
