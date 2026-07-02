# ownCloud Web

ownCloud Web is the next-generation frontend for ownCloud Infinite Scale, built as a single-page application with Vue.js and TypeScript. It provides a modern, accessible and themeable user interface for managing files, sharing, real-time collaboration and administration -- designed to be extensible through a plugin-based architecture that allows third-party developers to build custom apps and extensions.

## Getting Started

Set up a local development environment for ownCloud Web:

### Prerequisites

- [Node.js](https://nodejs.org/) (LTS recommended)
- [pnpm](https://pnpm.io/)
- [Docker Compose](https://docs.docker.com/compose/) (for backend)

For the complete development setup, see the [setup guide](https://owncloud.dev/clients/web/getting-started/).

### Structure

The `web/packages` directory contains the core modules:

- **client** -- Generated TypeScript client for the oCIS Graph API
- **container** -- Static assets and base files
- **extension-sdk** -- Utilities for developing custom extensions
- **pkg** -- Shared logic used across the codebase
- **runtime** -- Authentication, routing, theming and application handling

Built-in apps (also in `web/packages`):

- **files** -- Core file sync-and-share functionality
- **admin-settings** -- User, group and space administration
- **activities** -- Activity stream
- **app-store** -- Browse and install extensions
- **epub-reader** -- Ebook file viewer
- **external** -- WOPI-based document editing (Collabora, ONLYOFFICE)
- **pdf-viewer** -- PDF file viewer
- **preview** -- Audio, video and image previewer
- **text-editor** -- Plain text editor

## Documentation

- [ownCloud Web Documentation](https://owncloud.dev/clients/web)
- [Development Guide](https://owncloud.dev/clients/web/getting-started/)
- [Extension System](https://owncloud.dev/clients/web/extension-system/)
- [Testing Documentation](https://owncloud.dev/clients/web/testing/testing/)
- [Repository Structure](https://owncloud.dev/clients/web/development/repo-structure/)
