---
title: Theming
weight: 20
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/extensions/theming
geekdocFilePath: theming.md
---

{{< toc >}}

## Intro
Our default IDP UI is built with the [LibreGraph Connect](https://github.com/libregraph/lico) React app. Even though this app comes already with a simple theming options, we have compiled our own edited version of the app with more advanced changes than the default theming offers. Because of that, it is not possible at the moment to do any kind of easy theming and including custom theme means again compiling custom assets.

## Customizing assets
Depending on what changes you wish to do with the theme, there are several files you can edit. All of them are located in the `idp/ui` folder.

### Static assets
If you wish to add static assets like images, CSS, etc., you can add them to `idp/ui/public/static`. The `public` folder also contains the `index.html` file which can be adjusted to your needs.

### CSS
LibreGraph Connect is built with [kpop](https://github.com/Kopano-dev/kpop), a collection of React UI components. To include any custom styles on top of that collection, you can define them in the `idp/ui/src/app.css` file. These rules will take precedence over the kpop.

### Containers
Layouts of all pages are located in the `idp/ui/src/containers` folder. By editing any of files in that folder, you can do any kind of changes in the layout and create advanced themes. It is, however, important to be careful when touching this code as it imports also actions which are responsible for the login flow.

#### What pages to theme
- Login
  - Login - login form used to authenticate the users
  - Consent - consent page used to authorise apps for already signed in users
  - Chooseaccount - page with a list of accounts to choose from
- Goodbye
  - Goodbyescreen - goodbye message displayed to users after they signed out
- Welcome
  - Welcomescreen - welcome message displayed to users after they signed in

### Components
`idp/ui/src/components` folder contains all custom components which are then imported into containers.

### Images
Every image placed in `idp/ui/src/images` can be directly import into components or containers and will be optimized when compiling assets.

### Locales
If you need to edit or add new locales, you can do so with json files in the `idp/ui/src/locales` folder. If adding new locale, make sure to add it also in the `index.js` file in the same folder.

## Building assets
In order to build all assets, run `yarn build` in the `idp` folder. This script will compile all assets and output them into `idp/assets` folder.

At this point, you have two possible ways how to deploy your new theme:
- run `make generate` in the root folder of your oCIS clone and generate the new assets
- start the IDP service directly with custom assets by specifying the env var `IDP_ASSET_PATH`
