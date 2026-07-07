---
title: 'Deprecations and Migrations'
date: 2025-06-11T00:00:00+00:00
weight: 70
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/deprecations
geekdocFilePath: deprecations.md
geekdocCollapseSection: true
---

This guide tracks deprecations, planned removals, and migration paths for ownCloud Web. Developers should reference this guide when upgrading between major versions to ensure compatibility and adopt recommended practices.

{{< toc >}}

## Upcoming Deprecations

The following features are deprecated and **will be removed** in the next major version.

| Component/Feature | Deprecated | Migration | Reference | Reasoning |
|---|---|---|---|---|---|
| `owncloud-embed:share` event | November 12, 2025 | Use `owncloud-embed:share-links` event instead | [PR #13296](https://github.com/owncloud/web/pull/13296) | New event provides structured data with both URL and optional password |

## Migration History

The following features **are already removed**.

| Component/Feature | Removed in Version | Deprecated | Migration | Reference | Reasoning |
|---|---|---|---|---|---|
| `CollapsibleOcTable` component | v12 | May 16, 2025 | Use `OcTable` instead | [PR #12567](https://github.com/owncloud/web/pull/12567) | No longer needed by original implementers |
| `ocsUserContext` and `ocsPublicLinkContext` of `ClientService` | v12 | September 26, 2024 | Use `ocs` instead | [PR #11656](https://github.com/owncloud/web/pull/11656) | These methods were only a wrapper around the `ocs` |
| Web Client initialisation | v12 | September 26, 2024 | Use `graph()`, `ocs()` and `webdav()` to initialize and use clients | [PR #11656](https://github.com/owncloud/web/pull/11656) | More transparency to developers |
| `type` prop of `ApplicationInformation` | v12 | July 26, 2024 | N/A | [Commit 5e8ac91](https://github.com/owncloud/web/commit/5e8ac918de780f0b03ad79a85b8c25e4ba1fe4df) | The `type` prop is not used anymore |
| `isFileEditor` prop of `ApplicationInformation` | v12 | July 26, 2024 | N/A | [Commit 67ce21c](https://github.com/owncloud/web/commit/67ce21ce1bcc449d18f53e4930d180198a19faa0) | The `isFileEditor` prop is not used anymore |
| `ApplicationMenuItem` and `applicationMenu` prop of `ApplicationInformation` | v12 | July 24, 2024 | Register app menu items via the `appMenuExtensionPoint` instead | [PR #11258](https://github.com/owncloud/web/pull/11258) | Transitioning to the new extension points system |
| `ApplicationQuickAction` | v12 | December 1, 2023 | Register quick actions as extension instead | [PR #10102](https://github.com/owncloud/web/pull/10102) | Transitioning to the new extension points system |
| Vuex Store | v9 | January 22, 2024 | Migrate to Pinia stores | [PR #10372](https://github.com/owncloud/web/pull/10372) | Better TypeScript support and Composition API integration |
| Home Folder Option | v9 | January 5, 2024 | Use spaces functionality | [PR #10122](https://github.com/owncloud/web/pull/10122) | No longer needed by original implementers |
| OCS User API | v9 | January 3, 2024 | Use Graph user from `@ownclouders/web-client/src/helpers` and `useUserStore` | [PR #10240](https://github.com/owncloud/web/pull/10240) | The OCS API is deprecated and we're slowly transitioning to the Graph API and we want to get rid of the vuex store |
| Vuex Modal Implementation | v9 | December 22, 2023 | Use `useModals()` composable | [PR #10212](https://github.com/owncloud/web/pull/10212) | Remove Vuex store dependency |
| App: `user-management` | v8 | January 4, 2023 | Use the name `admin-settings` when referencing the app (in your configs for example) | [PR #8175](https://github.com/owncloud/web/pull/8175) | The app is supposed to hold a variety of general settings in the future, not just user-related settings |
| `mediaSource` helper & `v-image-source` directive | v6 | July 27, 2022 | Use `loadPreview` from web-pkg | N/A | `loadPreview` is the new default mechanism. We want to get rid of magical everywhere-available helpers |
| `getToken` getter | v6 | July 7, 2022 | Use `runtime/auth/accessToken` getter | N/A | `user` state is not supposed to deal with authentication details. We introduced the namespaced `auth` module for `user` and `public-link` related authentication contexts |

## Getting Help

For questions about specific migrations, refer to the linked pull requests and commits for detailed implementation examples and discussion.
