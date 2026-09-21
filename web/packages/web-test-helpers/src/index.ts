// Published entrypoint for everything outside `design-system` and `web-pkg`. The helpers
// themselves live in those two packages, so that neither has to depend on a package above it:
//
//   design-system/testing  mount, shallowMount, prop types, design system + gettext plugins
//   web-pkg/src/testing    the above + casl abilities, mocked pinia stores, component mocks
//   web-test-helpers       this re-export
export * from '@ownclouders/web-pkg/src/testing'
