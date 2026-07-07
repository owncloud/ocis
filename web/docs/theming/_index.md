---
title: 'Theming'
date: 2021-04-01T00:00:00+00:00
weight: 55
geekdocRepo: https://github.com/owncloud/web
geekdocEditPath: edit/master/docs/theming
geekdocFilePath: _index.md
geekdocCollapseSection: true
---

{{< toc >}}

## Introduction

By providing your own theme, you can customize the user experience for your own ownCloud installation. This is being achieved by providing a `json` file that contains text snippets (like brand name and slogan), paths to images (e.g. logos or favicon) and design tokens for various color, sizing and spacing parameters.

This page documents the setup and configuration options, and provides a template for you to get started.

## Ways of providing a theme

Generally, your theming configuration lives inside a `.json` file, e.g. `theme.json`. To load this file, it needs to be correctly referenced inside your `config/config.json` (example configurations can be [found on GitHub](https://github.com/owncloud/web/tree/master/config)).

To reference your theme, you have two options:

- Using a URL, e.g. `"theme": "https://externalurl.example.com/theme-name/theme.json",`. To avoid CORS issues, please make sure that you host the URL on the same URL as your ownCloud web hosting.
- For development and testing purposes, you can store your `theme.json` inside `packages/web-runtime/themes/{theme-name}/` and reference it in the `config.json`. However, this isn't recommended for production use since your changes may get lost when updating oCIS or the `web` app in ownCloud Classic.

**Hint:** If no theme is provided, the loading of your custom theme fails or the theme can't be parsed correctly, the standard ownCloud theme will be loaded as a fallback and an error with further information will be logged on the browser console.

## Configuring a theme

Inside your `theme.json`, there is a `common` key, which is explained in the next section, and a `clients` key: Here, you can find the available ownCloud clients - please note that the documentation below focuses on `web` and check the respective documentation for other clients for details on their themability.

The general top-level structure of a valid `theme.json` is outlined below:

```json
{
  "common": {},
  "clients": {
    "android": {},
    "desktop": {},
    "ios": {},
    "web": {}
  }
}
```

### Common section

The `common` section provides a set of information that is designed to be available for all clients. It gets merged "down" to the final themes and aims to reduce duplication, but can be overwritten by more specific information inside both the clients' defaults and actual themes.

The structure of a valid `common` section is outlined below:

```json
"common": {
  "name": "ownCloud",
  "slogan": "ownCloud – A safe home for all your data",
  "logo": "themes/owncloud/assets/logo.svg",
  "urls": {
    "accessDeniedHelp": "",
    "imprint": "",
    "privacy": ""
  }, 
  "shareRoles": {}
}
```

All of the below parameters are required:
- `name` specifies the publicly visible name
- `slogan` specifies the publicly visible slogan
- `logo` specifies the logo in e.g. the top bar within the web UI
- `accessDeniedHelp` specifies the target URL for the access denied help link
- `imprintUrl` specifies the target URL for the imprint link
- `privacyUrl` specifies the target URL for the privacy link

### Web Theme

The structure of a valid `web` client section is outlined below:

```json
{
  "web": {
    "defaults": {
      "appBanner": {
        // Please see below for details
      },
      "common": {
        // Please see top level "common" section for details
      },
      "logo": {
        // Please see below for details
      },
      "loginPage": {
        // Please see below for details
      },
      "designTokens": {
        // Please see below for details
      }
    },
    "themes": [
      // Your custom web themes go here, see below for details
    ]
  }
}
```

#### The "defaults"

Similar to the top level `common` section, this object contains information that shall be shared among the available themes and can/should be defined only once. The top level `common` section first gets merged into the `defaults`, which then get merged into the individual themes.

##### The "appBanner" options

Configures a app banner that gets shown on mobile devices and suggests downloading the native client from the respective app store. Omitting the key disables the banner.

Example structure:

```json 
{
  "appBanner": {
    "title": "ownCloud",
    "publisher": "ownCloud GmbH",
    "additionalInformation": "",
    "ctaText": "OPEN",
    "icon": "themes/owncloud/assets/owncloud-app-icon.png",
    "appScheme": "owncloud"
  }
}
```

- `title` is usually your app's name as shown in the App Store or Google Play. `publisher` is the app developer's name.
- `additionalInformation` can be used to specify pricing information, such as "FREE" or a catchphrase like "Don't miss out on our awesome app!".
- `ctaText` refers to the text in the call to action button on the right side. The `icon` directive may be used to specify your own app icon.
- `icon` links the icon to be displayed as a preview for the final app icon within the app banner
- `appScheme` is the first part of the URL that is used to tell the mobile OS which app to open, so using `ownCloud` will generate links such as `owncloud://yourdomain.com/f/2b61b822...`.

##### The "logo" options

Here, you can specify the images to be used in the `"topbar"`, for the `"favicon"` and on the `"login"` page. Various formats are supported and it's up to you to decide which one fits your use case best.

```json
"logo": {
  "topbar": "themes/owncloud/assets/logo.svg",
  "favicon": "themes/owncloud/assets/favicon.jpg",
  "login": "themes/owncloud/assets/logo.svg"
},
```

##### The "loginPage" options

You can set the background image for the login page by providing an image file in the `"backgroundImg"` option.

```json
"loginPage": {
  "backgroundImg": "themes/owncloud/assets/loginBackground.jpg"
},
```

##### The "designTokens" options

To further customize your ownCloud instance, you can provide your own styles. To give you an idea of how a working design system looks like, feel free to head over to the [ownCloud design tokens](https://owncloud.design/#/Design%20Tokens) for inspiration.

**Hint:** All the variables are initialized using the [ownCloud design tokens](https://owncloud.design/#/Design%20Tokens) and then overwritten by your theme variables. Therefore, you don't have to provide all the variables and can use the default ownCloud colors as a fallback.

In general, the theme loader looks for a `designTokens` key inside your theme configuration. Inside the `designTokens`, it expects to find a `colorPalette`, `fontSizes` and `spacing` collection. The structure is outlined below:

```json
"designTokens": {
  "breakpoints": {
    // Please see below for details
  },
  "colorPalette": {
    // Please see below for details
  },
  "fontFamily": "", // Please see below for details
  "fontSizes": {
    // Please see below for details
  },
  "sizes": {
    // Please see below for details
  },
  "spacing": {
    // Please see below for details
  }
}
```

###### Breakpoints

If you'd like to set different breakpoints than the default ones in the ownCloud design system, you can set them using theming variables.

Breakpoint variables get prepended with `--oc-breakpoint-`, so e.g. _"large-default"_ creates the custom CSS property `--oc-breakpoint-large-default`.

```json
{
  "breakpoints": {
    "xsmall-max": "",
    "small-default": "",
    "small-max": "",
    "medium-default": "",
    "medium-max": "",
    "large-default": "",
    "large-max": "",
    "xlarge": ""
  }
}
```

###### Colors

For the color values, you can use any valid CSS color format, e.g. **hex** (#fff), **rgb** (rgb(255,255,255)) or **color names** (white).

Color variables get prepended with `--oc-color-`, so e.g. _"background-default"_ creates the custom CSS property `--oc-color-background-default`.

Again, you can use the [ownCloud design tokens](https://owncloud.design/#/Design%20Tokens) as a reference implementation.

```json
{
  "colorPalette": {
    "background-accentuate": "",
    "background-default": "",
    "background-highlight": "",
    "background-muted": "",
    "border": "",
    "input-bg": "",
    "input-border": "",
    "input-text-default": "",
    "input-text-muted": "",
    "swatch-brand-default": "",
    "swatch-brand-hover": "",
    "swatch-brand-muted": "",
    "swatch-brand-contrast": "",
    "swatch-danger-default": "",
    "swatch-danger-hover": "",
    "swatch-danger-muted": "",
    "swatch-danger-contrast": "",
    "swatch-inverse-default": "",
    "swatch-inverse-hover": "",
    "swatch-inverse-muted": "",
    "swatch-passive-default": "",
    "swatch-passive-hover": "",
    "swatch-passive-hover-outline": "",
    "swatch-passive-muted": "",
    "swatch-passive-contrast": "",
    "swatch-primary-default": "",
    "swatch-primary-hover": "",
    "swatch-primary-muted": "",
    "swatch-primary-muted-hover": "",
    "swatch-primary-gradient": "",
    "swatch-primary-gradient-hover": "",
    "swatch-primary-contrast": "",
    "swatch-success-default": "",
    "swatch-success-hover": "",
    "swatch-success-muted": "",
    "swatch-success-contrast": "",
    "swatch-warning-default": "",
    "swatch-warning-hover": "",
    "swatch-warning-muted": "",
    "swatch-warning-contrast": "",
    "text-default": "",
    "text-inverse": "",
    "text-muted": ""
  }
}
```

###### Font sizes

You can change the `default`, `large` and `medium` font sizes according to your needs. If you need more customization options regarding font sizes, please [open an issue on GitHub](https://github.com/owncloud/web/issues/new) with a detailed description.

Font size variables get prepended with `--oc-font-size-`, so e.g. _"default"_ creates the custom CSS property `--oc-font-size-default`.

```json
{
  "fontSizes": {
    "default": "",
    "large": "",
    "medium": ""
  }
}
```

###### Font family

You can change the font family according to your needs. The font family gets written into the `--oc-font-family` CSS variable.

```json
{
  "fontFamily": ""
}
```

Please note that you also need to make the font available as a `font-face` via CSS.

###### Sizes

Use sizing variables to change various UI elements, such as icon and logo appearance, table row or checkbox sizes, according to your needs.
If you need more customization options regarding sizes, please [open an issue on GitHub](https://github.com/owncloud/web/issues/new) with a detailed description.

Size variables get prepended with `--oc-size-`, so e.g. _"icon-default"_ creates the custom CSS property `--oc-size-icon-default`.

```json
{
  "sizes": {
    "form-check-default": "",
    "height-small": "",
    "height-table-row": "",
    "icon-default": "",
    "max-height-logo": "",
    "max-width-logo": "",
    "width-medium": "",
    "tiles-default": "",
    "tiles-resize-step": ""
  }
}
```

###### Spacing

Use the six spacing options (`xsmall | small | medium | large | xlarge | xxlarge`) to create a more (or less) condensed version of the user interface. If you need more customization options regarding sizes, please [open an issue on GitHub](https://github.com/owncloud/web/issues/new) with a detailed description.

Spacing variables get prepended with `--oc-space-`, so e.g. _"xlarge"_ creates the custom CSS property `--oc-space-xlarge`.

```json
{
  "spacing": {
    "xsmall": "",
    "small": "",
    "medium": "",
    "large": "",
    "xlarge": "",
    "xxlarge": ""
  }
}
```

#### Actual Themes

Apart from the `defaults`, you need to provide one or more themes in the `themes` key within the `web`-`clients` in your `theme.json`. As a reminder, the general structure should be:

```json
{
  "common": { ... },
  "clients": {
    ...,
    "web": {
      "defaults": {
        ...
      },
      "themes": [
        {
          "isDark": false,
          "name": "Light Theme",
        }
      ]
    }
  }
}
```

Again, both the global `common` section as well as the `defaults` will get merged into your themes, but locally provided information takes precedence.

Required information
- `name` for the visible name in the theme switcher and to save the current theme to localStorage
- `isDark` to provide the user agent with additional information

Optional information
- `appBanner` see section above
- `common` see section above
- `designTokens` see section above
- `logo` see section above
- `loginPage` see section above

## Extendability

If you define different key-value pairs inside any of the objects (`breakpoints`, `colorPalette`, `fontSizes`, `sizes`, `spacing`) in `"designTokens"`, they will get loaded and initialized as CSS custom properties but don't necessarily take any effect in the user interface. This gives you an opportunity to, for example, customize extensions from within the theme in the web runtime (and not the extension itself).

## Example theme

A full template for your custom theme is provided below, and you can use the instructions above to set it up according to your needs:

```json
{
  "common": {
    "name": "ownCloud",
    "slogan": "ownCloud – A safe home for all your data",
    "logo": "themes/owncloud/assets/logo.svg",
    "urls": {
      "accessDeniedHelp": "",
      "imprint": "",
      "privacy": ""
    },
    "shareRoles": {
      "b1e2218d-eef8-4d4c-b82d-0f1a1b48f3b5": {
        "name": "UnifiedRoleViewer",
        "iconName": "eye"
      },
      "a8d5fe5e-96e3-418d-825b-534dbdf22b99": {
        "label": "UnifiedRoleSpaceViewer",
        "iconName": "eye"
      },
      "2d00ce52-1fc2-4dbc-8b95-a73b73395f5a": {
        "label": "UnifiedRoleFileEditor",
        "iconName": "pencil"
      },
      "fb6c3e19-e378-47e5-b277-9732f9de6e21": {
        "label": "UnifiedRoleEditor",
        "iconName": "pencil"
      },
      "58c63c02-1d89-4572-916a-870abc5a1b7d": {
        "label": "UnifiedRoleSpaceEditor",
        "iconName": "pencil"
      },
      "312c0871-5ef7-4b3a-85b6-0e4074c64049": {
        "label": "UnifiedRoleManager",
        "iconName": "user-star"
      },
      "1c996275-f1c9-4e71-abdf-a42f6495e960": {
        "label": "UnifiedRoleUploader",
        "iconName": "pencil"
      },
      "aa97fe03-7980-45ac-9e50-b325749fd7e6": {
        "label": "UnifiedRoleSecureView",
        "iconName": "shield"
      }
    }
  },
  "clients": {
    "android": {},
    "desktop": {},
    "ios": {},
    "web": {
      "defaults": {
        "logo": {
          "topbar": "themes/owncloud/assets/logo.svg",
          "favicon": "themes/owncloud/assets/favicon.jpg",
          "login": "themes/owncloud/assets/logo.svg"
        },
        "loginPage": {
          "backgroundImg": "themes/owncloud/assets/loginBackground.jpg"
        },
        "designTokens": {
          "breakpoints": {
            "xsmall-max": "",
            "small-default": "",
            "small-max": "",
            "medium-default": "",
            "medium-max": "",
            "large-default": "",
            "large-max": "",
            "xlarge": ""
          },
          "fontSizes": {
            "default": "",
            "large": "",
            "medium": ""
          },
          "sizes": {
            "form-check-default": "",
            "height-small": "",
            "height-table-row": "",
            "icon-default": "",
            "max-height-logo": "",
            "max-width-logo": "",
            "width-medium": "",
            "tiles-default": "",
            "tiles-resize-step": ""
          },
          "spacing": {
            "xsmall": "",
            "small": "",
            "medium": "",
            "large": "",
            "xlarge": "",
            "xxlarge": ""
          }
        }
      },
      "themes": [
        {
          "isDark": false,
          "name": "Light Theme",
          "designTokens": {
            "colorPalette": {
              "background-accentuate": "rgba(255, 255, 5, 0.1)",
              "background-default": "#ffffff",
              "background-highlight": "#edf3fa",
              "background-muted": "#f8f8f8",
              "background-secondary": "#ffffff",
              "background-hover": "rgb(236, 236, 236)",
              "color-components-apptopbar-background": "transparent",
              "color-components-apptopbar-border": "#ceddee",
              "border": "#ecebee",
              "input-bg": "#ffffff",
              "input-border": "#ceddee",
              "input-text-default": "#041e42",
              "input-text-muted": "#4c5f79",
              "swatch-brand-default": "#041e42",
              "swatch-brand-hover": "#223959",
              "swatch-brand-contrast": "#ffffff",
              "swatch-danger-contrast": "#ffffff",
              "swatch-danger-default": "rgb(197, 48, 48)",
              "swatch-danger-hover": "#b12b2b",
              "swatch-danger-muted": "rgb(204, 117, 117)",
              "swatch-inverse-default": "#ffffff",
              "swatch-inverse-hover": "#ffffff",
              "swatch-inverse-muted": "#bfbfbf",
              "swatch-passive-default": "#4c5f79",
              "swatch-passive-hover": "#43536b",
              "swatch-passive-hover-outline": "#f7fafd",
              "swatch-passive-muted": "#283e5d",
              "swatch-passive-contrast": "#ffffff",
              "swatch-primary-default": "#4a76ac",
              "swatch-primary-hover": "#80a7d7",
              "swatch-primary-muted": "#2c588e",
              "swatch-primary-muted-hover": "rgb(36, 75, 119)",
              "swatch-primary-gradient": "#4e85c8",
              "swatch-primary-gradient-hover": "rgb(59, 118, 194)",
              "swatch-primary-contrast": "#ffffff",
              "swatch-success-default": "rgb(3, 84, 63)",
              "swatch-success-hover": "#023b2c",
              "swatch-success-muted": "rgb(83, 150, 10)",
              "swatch-success-contrast": "#ffffff",
              "swatch-warning-default": "rgb(183, 76, 27)",
              "swatch-warning-hover": "#a04318",
              "swatch-warning-muted": "rgba(183, 76, 27, .5)",
              "swatch-warning-contrast": "#ffffff",
              "text-default": "#041e42",
              "text-inverse": "#ffffff",
              "text-muted": "#4c5f79",
              "icon-folder": "#4d7eaf",
              "icon-archive": "#fbbe54",
              "icon-image": "#ee6b3b",
              "icon-spreadsheet": "#15c286",
              "icon-document": "#3b44a6",
              "icon-video": "#045459",
              "icon-audio": "#700460",
              "icon-presentation": "#ee6b3b",
              "icon-pdf": "#ec0d47"
            }
          }
        },
        {
          "isDark": true,
          "name": "Dark Theme",
          "designTokens": {
            "colorPalette": {
              "background-accentuate": "#696969",
              "background-default": "#292929",
              "background-highlight": "#383838",
              "background-muted": "#383838",
              "background-secondary": "#4f4f4f",
              "background-hover": "#383838",
              "color-components-apptopbar-background": "transparent",
              "color-components-apptopbar-border": "#ceddee",
              "border": "#383838",
              "input-bg": "#4f4f4f",
              "input-border": "#696969",
              "input-text-default": "#dadcdf",
              "input-text-muted": "#bdbfc3",
              "swatch-brand-default": "#212121",
              "swatch-brand-hover": "#ffffff",
              "swatch-brand-contrast": "#dadcdf",
              "swatch-inverse-default": "",
              "swatch-inverse-hover": "",
              "swatch-inverse-muted": "#696969",
              "swatch-passive-default": "#c2c2c2",
              "swatch-passive-hover": "",
              "swatch-passive-hover-outline": "#3B3B3B",
              "swatch-passive-muted": "#bdbfc3",
              "swatch-passive-contrast": "#000000",
              "swatch-primary-default": "#73b0f2",
              "swatch-primary-hover": "#7bafef",
              "swatch-primary-muted": "",
              "swatch-primary-muted-hover": "#2282f7",
              "swatch-primary-gradient": "#4e85c8",
              "swatch-primary-gradient-hover": "#76a1d5",
              "swatch-primary-contrast": "#dadcdf",
              "swatch-success-background": "rgba(0, 188, 140, 0)",
              "swatch-success-default": "rgb(0, 188, 140)",
              "swatch-success-hover": "#00f0b4",
              "swatch-success-muted": "rgba(0, 188, 140, .5)",
              "swatch-success-contrast": "#000000",
              "swatch-warning-background": "rgba(0,0,0,0)",
              "swatch-warning-default": "rgb(232, 191, 73)",
              "swatch-warning-hover": "#eed077",
              "swatch-warning-muted": "rgba(232, 178, 19, .5)",
              "swatch-danger-default": "rgb(255, 72, 53)",
              "swatch-danger-hover": "#ff7566",
              "swatch-danger-muted": "rgba(255, 72, 53, .5)",
              "swatch-danger-contrast": "#dadcdf",
              "swatch-warning-contrast": "#000000",
              "text-default": "#dadcdf",
              "text-inverse": "#000000",
              "text-muted": "#c2c2c2",
              "icon-folder": "rgb(44, 101, 255)",
              "icon-archive": "rgb(255, 207, 1)",
              "icon-image": "rgb(255, 111, 0)",
              "icon-spreadsheet": "rgb(0, 182, 87)",
              "icon-document": "rgb(44, 101, 255)",
              "icon-video": "rgb(0, 187, 219)",
              "icon-audio": "rgb(208, 67, 236)",
              "icon-presentation": "rgb(255, 64, 6)",
              "icon-pdf": "rgb(225, 5, 14)"
            }
          }
        }
      ]
    }
  }
}
```
