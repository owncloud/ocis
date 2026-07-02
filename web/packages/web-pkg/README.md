# web-pkg

`web-pkg` is a package that provides utilities, most importantly a variety of components and composables, that can be useful when developing apps and extensions for ownCloud Web.

## Installation

Depending on your package manager, run one of the following commands:

```
$ npm install @ownclouders/web-pkg --save-dev

$ pnpm add -D @ownclouders/web-pkg

$ yarn add @ownclouders/web-pkg --dev
```

It's recommended to install this package as a dev dependency because it's only really needed for providing autocompletion in your IDE and unit tests. In a runtime context, the ownCloud Web runtime provides the actual implementation.
