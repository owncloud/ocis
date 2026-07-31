#!/bin/bash

# This script bumps the version in the package.json files of the main app
# and all published packages to the provided version.
#
# NOTE: only run this if all packages are supposed to get the same version!

apps=("design-system" "eslint-config" "extension-sdk" "web-pkg" "web-client" "web-test-helpers")

NEW_VERSION="$1"

if [ -z "$NEW_VERSION" ]; then
	echo "Error: No new version provided."
	echo "Usage: $0 <new_version>"
	exit 1
fi

cd "$(dirname "$0")/../.."
CURRENT_VERSION=$(node -p "require('./package.json').version")
sed -i '' "s/\"version\": \"${CURRENT_VERSION}\"/\"version\": \"${NEW_VERSION}\"/" package.json

for app in "${apps[@]}"; do
	cd "./packages/$app"
	CURRENT_VERSION=$(node -p "require('./package.json').version")
	sed -i '' "s/\"version\": \"${CURRENT_VERSION}\"/\"version\": \"${NEW_VERSION}\"/" package.json
	cd "../.."
done

echo "package.json files have been updated to version $NEW_VERSION"

SONAR_PROPERTIES_FILE="sonar-project.properties"

if [ ! -f "$SONAR_PROPERTIES_FILE" ]; then
    echo "Error: $SONAR_PROPERTIES_FILE file not found!"
    exit 1
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' "s/sonar\.projectVersion=[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/sonar.projectVersion=$NEW_VERSION/" "$SONAR_PROPERTIES_FILE"
else
    sed -i "s/sonar\.projectVersion=[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*/sonar.projectVersion=$NEW_VERSION/" "$SONAR_PROPERTIES_FILE"
fi

if grep -q "sonar.projectVersion=$NEW_VERSION" "$SONAR_PROPERTIES_FILE"; then
    echo "Sonar project version successfully updated to $NEW_VERSION"
else
    echo "Failed to update Sonar project version. Please check the file format."
    exit 1
fi
