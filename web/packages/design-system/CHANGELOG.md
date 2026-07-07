# Changelog for [14.0.0] (2022-11-03)

The following sections list the changes in ownCloud Design System 14.0.0.

[14.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v13.1.0...v14.0.0

## Summary

* Bugfix - Omit special characters in user avatar initials: [#2070](https://github.com/owncloud/owncloud-design-system/issues/2070)
* Bugfix - Avatar link icon: [#2269](https://github.com/owncloud/owncloud-design-system/pull/2269)
* Bugfix - Firefox drag & drop move of folders not possible: [#7495](https://github.com/owncloud/web/issues/7495)
* Bugfix - Lazy loading render performance: [#2260](https://github.com/owncloud/owncloud-design-system/pull/2260)
* Bugfix - Modal input message overlays with buttons: [#2343](https://github.com/owncloud/owncloud-design-system/pull/2343)
* Bugfix - Remove width shrinking of the ocAvatarItem: [#2241](https://github.com/owncloud/owncloud-design-system/issues/2241)
* Bugfix - Remove click event on OcIcon: [#2216](https://github.com/owncloud/owncloud-design-system/pull/2216)
* Bugfix - Translate contextual helpers: [#2334](https://github.com/owncloud/owncloud-design-system/pull/2334)
* Change - Redesign contextual helper: [#2271](https://github.com/owncloud/owncloud-design-system/pull/2271)
* Change - Remove OcAlert component: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)
* Change - Remove transition animations: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)
* Change - Revamp animations: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)
* Change - OcTable emit event data on row click: [#2218](https://github.com/owncloud/owncloud-design-system/pull/2218)
* Enhancement - Give OcModal the option to use OcContextualHelper: [#2325](https://github.com/owncloud/owncloud-design-system/pull/2325)
* Enhancement - Add nestedd drop functionality: [#2238](https://github.com/owncloud/owncloud-design-system/issues/2238)
* Enhancement - Add OcInfoDrop: [#2286](https://github.com/owncloud/owncloud-design-system/pull/2286)
* Enhancement - Add rounded prop to OcTag: [#2284](https://github.com/owncloud/owncloud-design-system/pull/2284)
* Enhancement - Adjust avatar font weight from bold to normal: [#2275](https://github.com/owncloud/owncloud-design-system/pull/2275)
* Enhancement - Align breadcrumb context menu with regular context menu: [#2296](https://github.com/owncloud/owncloud-design-system/pull/2296)
* Enhancement - Adjust breadcrumb spacing: [#7676](https://github.com/owncloud/web/issues/7676)
* Enhancement - Button text align left: [#7619](https://github.com/owncloud/web/issues/7619)
* Enhancement - OcCheckbox add outline: [#2218](https://github.com/owncloud/owncloud-design-system/pull/2218)
* Enhancement - Add offset property to the drop component: [#7335](https://github.com/owncloud/web/issues/7335)
* Enhancement - Input background color: [#7353](https://github.com/owncloud/web/issues/7353)
* Enhancement - Make UI smaller: [#2270](https://github.com/owncloud/owncloud-design-system/pull/2270)
* Enhancement - Oc-card style: [#2306](https://github.com/owncloud/owncloud-design-system/pull/2306)
* Enhancement - OcSelect dark mode improvements: [#2262](https://github.com/owncloud/owncloud-design-system/pull/2262)
* Enhancement - Progress bar indeterminate state: [#2200](https://github.com/owncloud/owncloud-design-system/pull/2200)
* Enhancement - Redesign notifications: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)
* Enhancement - Remove border on buttons: [#7353](https://github.com/owncloud/web/issues/7353)
* Enhancement - "Chancel"-button and -handler in OcSearchBar: [#7617](https://github.com/owncloud/web/issues/7617)
* Enhancement - Use Inter font: [#2270](https://github.com/owncloud/owncloud-design-system/pull/2270)

## Details

* Bugfix - Omit special characters in user avatar initials: [#2070](https://github.com/owncloud/owncloud-design-system/issues/2070)

   We now make sure that user avatars without a picture (i.e. rendering user initials) only show
   unicode letters and numbers.

   https://github.com/owncloud/owncloud-design-system/issues/2070
   https://github.com/owncloud/owncloud-design-system/pull/2267


* Bugfix - Avatar link icon: [#2269](https://github.com/owncloud/owncloud-design-system/pull/2269)

   We've updated the avatar link icon to the current icon used in Web.

   https://github.com/owncloud/web/issues/7345
   https://github.com/owncloud/owncloud-design-system/pull/2269


* Bugfix - Firefox drag & drop move of folders not possible: [#7495](https://github.com/owncloud/web/issues/7495)

   We've fixed a bug in firefox which caused drag & drop move to redirect the page.

   https://github.com/owncloud/web/issues/7495
   https://github.com/owncloud/owncloud-design-system/pull/2302


* Bugfix - Lazy loading render performance: [#2260](https://github.com/owncloud/owncloud-design-system/pull/2260)

   The render performance of the lazy loading option in tables (OcTable, OcTableSimple) has been
   improved by removing the debounce option and by moving the lazy loading visualization from the
   OcTd to the OcTr component. For lazy loading, the colspan property has to be provided now.

   https://github.com/owncloud/web/issues/7038
   https://github.com/owncloud/owncloud-design-system/pull/2260
   https://github.com/owncloud/owncloud-design-system/pull/2266


* Bugfix - Modal input message overlays with buttons: [#2343](https://github.com/owncloud/owncloud-design-system/pull/2343)

   We've fixed a bug where the modal input message eventually overlays with the confirm and cancel
   buttons.

   https://github.com/owncloud/web/issues/7807
   https://github.com/owncloud/owncloud-design-system/pull/2343


* Bugfix - Remove width shrinking of the ocAvatarItem: [#2241](https://github.com/owncloud/owncloud-design-system/issues/2241)

   We fixed an issue that the width of ocAvatarItem is shrinking in the sidebar of web by longer
   group names

   https://github.com/owncloud/owncloud-design-system/issues/2241
   https://github.com/owncloud/owncloud-design-system/pull/2242


* Bugfix - Remove click event on OcIcon: [#2216](https://github.com/owncloud/owncloud-design-system/pull/2216)

   We have removed an unnecessary default click handler on the OcIcon component, expecting it to
   increase performance of the UI.

   https://github.com/owncloud/owncloud-design-system/pull/2216


* Bugfix - Translate contextual helpers: [#2334](https://github.com/owncloud/owncloud-design-system/pull/2334)

   We've fixed a bug where contextual helpers were not translated.

   https://github.com/owncloud/web/issues/7716
   https://github.com/owncloud/owncloud-design-system/pull/2334


* Change - Redesign contextual helper: [#2271](https://github.com/owncloud/owncloud-design-system/pull/2271)

   We've redesigned the contextual helper, which accepts now a title property and is able to
   display a description list.

   https://github.com/owncloud/web/issues/7331
   https://github.com/owncloud/owncloud-design-system/pull/2271
   https://github.com/owncloud/owncloud-design-system/pull/2273


* Change - Remove OcAlert component: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)

   https://github.com/owncloud/owncloud-design-system/pull/2210


* Change - Remove transition animations: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)

   https://github.com/owncloud/owncloud-design-system/pull/2210


* Change - Revamp animations: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)

   We have removed the old animation classes and will subsequently add new ones, respecting the
   `prefers-reduced-motion` browser setting.

   https://github.com/owncloud/owncloud-design-system/pull/2210


* Change - OcTable emit event data on row click: [#2218](https://github.com/owncloud/owncloud-design-system/pull/2218)

   We've extended the emit data on row click to now also include the event data

   https://github.com/owncloud/web/issues/6029
   https://github.com/owncloud/owncloud-design-system/pull/2218


* Enhancement - Give OcModal the option to use OcContextualHelper: [#2325](https://github.com/owncloud/owncloud-design-system/pull/2325)

   We've added the option for OcModal to use OcContextualHelper

   https://github.com/owncloud/web/issues/6892
   https://github.com/owncloud/owncloud-design-system/pull/2325


* Enhancement - Add nestedd drop functionality: [#2238](https://github.com/owncloud/owncloud-design-system/issues/2238)

   We've added the property "isNested" to ocDrop that prevents the parent drop from hiding by not
   firing hideAll() in the child drop

   https://github.com/owncloud/owncloud-design-system/issues/2238
   https://github.com/owncloud/owncloud-design-system/pull/2239


* Enhancement - Add OcInfoDrop: [#2286](https://github.com/owncloud/owncloud-design-system/pull/2286)

   We've added the new component OcInfoDrop, which will be consumed by the OcContextualHelper
   component.

   https://github.com/owncloud/owncloud-design-system/pull/2286


* Enhancement - Add rounded prop to OcTag: [#2284](https://github.com/owncloud/owncloud-design-system/pull/2284)

   We've added a rounded property to the OcTag component

   https://github.com/owncloud/owncloud-design-system/pull/2284


* Enhancement - Adjust avatar font weight from bold to normal: [#2275](https://github.com/owncloud/owncloud-design-system/pull/2275)

   https://github.com/owncloud/owncloud-design-system/pull/2275


* Enhancement - Align breadcrumb context menu with regular context menu: [#2296](https://github.com/owncloud/owncloud-design-system/pull/2296)

   We've aligned the breadcrumb context menu visually to match with the regular context menu in
   the files table.

   https://github.com/owncloud/web/issues/7493
   https://github.com/owncloud/owncloud-design-system/pull/2296


* Enhancement - Adjust breadcrumb spacing: [#7676](https://github.com/owncloud/web/issues/7676)

   We've adjusted some spacing in the breadcrumbs to improve the overall look.

   https://github.com/owncloud/web/issues/7676
   https://github.com/owncloud/web/issues/7525
   https://github.com/owncloud/owncloud-design-system/pull/2329


* Enhancement - Button text align left: [#7619](https://github.com/owncloud/web/issues/7619)

   We've changed the text alignment of buttons to left.

   https://github.com/owncloud/web/issues/7619
   https://github.com/owncloud/owncloud-design-system/pull/2323


* Enhancement - OcCheckbox add outline: [#2218](https://github.com/owncloud/owncloud-design-system/pull/2218)

   We've added an optional outline to be able to highlight the checkbox

   https://github.com/owncloud/web/issues/6029
   https://github.com/owncloud/owncloud-design-system/pull/2218


* Enhancement - Add offset property to the drop component: [#7335](https://github.com/owncloud/web/issues/7335)

   We've added an offset property to the drop component to define a custom offset. Also, the max
   width of drop menus was increased to 400px.

   https://github.com/owncloud/web/issues/7335
   https://github.com/owncloud/owncloud-design-system/pull/2276


* Enhancement - Input background color: [#7353](https://github.com/owncloud/web/issues/7353)

   The background color for input fields has been adjusted to better match with the overall
   design.

   https://github.com/owncloud/web/issues/7353
   https://github.com/owncloud/web/issues/7373
   https://github.com/owncloud/owncloud-design-system/pull/2345
   https://github.com/owncloud/owncloud-design-system/pull/2352


* Enhancement - Make UI smaller: [#2270](https://github.com/owncloud/owncloud-design-system/pull/2270)

   We've adjusted several values to make the UI appear less big.

   https://github.com/owncloud/owncloud-design-system/pull/2270


* Enhancement - Oc-card style: [#2306](https://github.com/owncloud/owncloud-design-system/pull/2306)

   We've enhanced the oc-card style classes, to fit better in the corporate design

   https://github.com/owncloud/web/issues/7537
   https://github.com/owncloud/owncloud-design-system/pull/2306
   https://github.com/owncloud/owncloud-design-system/pull/2321


* Enhancement - OcSelect dark mode improvements: [#2262](https://github.com/owncloud/owncloud-design-system/pull/2262)

   We've improved the visual appearance of the OcSelect component in dark mode, now the selected
   items have an adjusted color and are easier to read.

   https://github.com/owncloud/web/issues/7269
   https://github.com/owncloud/owncloud-design-system/pull/2262


* Enhancement - Progress bar indeterminate state: [#2200](https://github.com/owncloud/owncloud-design-system/pull/2200)

   We've added an indeterminate state to the progress bar.

   https://github.com/owncloud/web/issues/7105
   https://github.com/owncloud/owncloud-design-system/pull/2200


* Enhancement - Redesign notifications: [#2210](https://github.com/owncloud/owncloud-design-system/pull/2210)

   We have redesigned the notifications component to fit the overal new look of the web frontend,
   e.g. adding shadow and rounded corners. It can now also be rendered "unpositioned" instead of
   having it always stick to the top of the screen.

   https://github.com/owncloud/web/issues/7082
   https://github.com/owncloud/owncloud-design-system/pull/2210
   https://github.com/owncloud/owncloud-design-system/pull/2216


* Enhancement - Remove border on buttons: [#7353](https://github.com/owncloud/web/issues/7353)

   The outer border for buttons has been removed. Outline buttons now have an inner outline
   instead.

   https://github.com/owncloud/web/issues/7353
   https://github.com/owncloud/web/issues/7373
   https://github.com/owncloud/owncloud-design-system/pull/2345
   https://github.com/owncloud/owncloud-design-system/pull/2352
   https://github.com/owncloud/owncloud-design-system/pull/2354


* Enhancement - "Chancel"-button and -handler in OcSearchBar: [#7617](https://github.com/owncloud/web/issues/7617)

   We've added to possibility to have a "cancel"-button and -handler in the `OcSearchBar`
   component.

   https://github.com/owncloud/web/issues/7617
   https://github.com/owncloud/owncloud-design-system/pull/2328


* Enhancement - Use Inter font: [#2270](https://github.com/owncloud/owncloud-design-system/pull/2270)

   We've switched the default font from Roboto to Inter.

   https://github.com/owncloud/owncloud-design-system/pull/2270
   https://github.com/owncloud/owncloud-design-system/pull/2299

# Changelog for [13.1.0] (2022-06-07)

The following sections list the changes in ownCloud Design System 13.1.0.

[13.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v13.0.0...v13.1.0

## Summary

* Bugfix - Disabled textarea color contrast in darkmode: [#2055](https://github.com/owncloud/owncloud-design-system/pull/2055)
* Bugfix - Fix event handler for submit-via-enter in OcModal input: [#2159](https://github.com/owncloud/owncloud-design-system/pull/2159)
* Bugfix - OcTextInput: Fix event handlers in loops: [#2054](https://github.com/owncloud/owncloud-design-system/pull/2054)
* Bugfix - Text color of buttons and avatars in hovered table rows: [#2139](https://github.com/owncloud/owncloud-design-system/pull/2139)
* Bugfix - Add word breaking in tooltips: [#2137](https://github.com/owncloud/owncloud-design-system/pull/2137)
* Enhancement - Add OcContextualHelper: [#6590](https://github.com/owncloud/web/issues/6590)
* Enhancement - Add ROOT file icon: [#2158](https://github.com/owncloud/owncloud-design-system/pull/2158)
* Enhancement - Add selection range for OcModal and OcTextInput: [#6729](https://github.com/owncloud/web/issues/6729)
* Enhancement - Export package members: [#2048](https://github.com/owncloud/owncloud-design-system/pull/2048)
* Enhancement - Hover in ocDrop menues: [#2069](https://github.com/owncloud/owncloud-design-system/pull/2069)
* Enhancement - OcModal add checkbox and secondary button: [#6994](https://github.com/owncloud/web/pull/6994)
* Enhancement - OcModal input type: [#2077](https://github.com/owncloud/owncloud-design-system/pull/2077)
* Enhancement - Make OcResource inline-flex: [#2041](https://github.com/owncloud/owncloud-design-system/pull/2041)
* Enhancement - Add `isFileExtensionDisplayed` property: [#2087](https://github.com/owncloud/owncloud-design-system/pull/2087)
* Enhancement - Redesign OcGhostElement: [#2049](https://github.com/owncloud/owncloud-design-system/pull/2049)
* Enhancement - Replace deprecated String.prototype.substr(): [#2059](https://github.com/owncloud/owncloud-design-system/pull/2059)
* Enhancement - Add HTML title to the resourceName component: [#2164](https://github.com/owncloud/owncloud-design-system/pull/2164)
* Enhancement - Add option to not truncate the resource name: [#2157](https://github.com/owncloud/owncloud-design-system/pull/2157)

## Details

* Bugfix - Disabled textarea color contrast in darkmode: [#2055](https://github.com/owncloud/owncloud-design-system/pull/2055)

   We fixed an issue that made text on disabled textarea fields unreadable since it was the same
   color as the background.

   https://github.com/owncloud/owncloud-design-system/issues/2053
   https://github.com/owncloud/owncloud-design-system/pull/2055


* Bugfix - Fix event handler for submit-via-enter in OcModal input: [#2159](https://github.com/owncloud/owncloud-design-system/pull/2159)

   We fixed a small difference between clicking the confirm button and hitting the enter key when
   using an OcModal with an input field, which lead to unexpected behaviour.

   https://github.com/owncloud/owncloud-design-system/pull/2159
   https://github.com/owncloud/web/pull/6961#pullrequestreview-979854434


* Bugfix - OcTextInput: Fix event handlers in loops: [#2054](https://github.com/owncloud/owncloud-design-system/pull/2054)

   We pass all event handlers specified on `OcTextInput` to the underlying `input` element
   except for `input`, `change` and `focus` event handlers. We fixed an issue in this exclusion
   code that made `change`, `input` and `focus` handlers be re-registered on rerenders,
   particularly in loop rerenders, so they were called multiple times for a single event.

   https://github.com/owncloud/owncloud-design-system/pull/2054


* Bugfix - Text color of buttons and avatars in hovered table rows: [#2139](https://github.com/owncloud/owncloud-design-system/pull/2139)

   We fixed an issue that made text of buttons and avatars inside hovered rows bad readable for
   light mode.

   https://github.com/owncloud/owncloud-design-system/issues/2138
   https://github.com/owncloud/owncloud-design-system/pull/2139


* Bugfix - Add word breaking in tooltips: [#2137](https://github.com/owncloud/owncloud-design-system/pull/2137)

   We've added word wrapping to the tippy tooltips so they handle very long paths properly.

   https://github.com/owncloud/owncloud-design-system/pull/2137


* Enhancement - Add OcContextualHelper: [#6590](https://github.com/owncloud/web/issues/6590)

   We've added a contextual helper component to provide more information based on the context

   https://github.com/owncloud/web/issues/6590
   https://github.com/owncloud/owncloud-design-system/pull/2064


* Enhancement - Add ROOT file icon: [#2158](https://github.com/owncloud/owncloud-design-system/pull/2158)

   We've added an icon for files of ROOT type. ROOT is a software suite designed for data analysis in
   particle physics, astronomy and other sciences.

   https://github.com/owncloud/owncloud-design-system/pull/2158


* Enhancement - Add selection range for OcModal and OcTextInput: [#6729](https://github.com/owncloud/web/issues/6729)

   We've added the possibility to set a selection range for the initial focus selection in OcModal
   and OcTextinput.

   https://github.com/owncloud/web/issues/6729
   https://github.com/owncloud/owncloud-design-system/pull/2061


* Enhancement - Export package members: [#2048](https://github.com/owncloud/owncloud-design-system/pull/2048)

   Add exports for `composables`, `utils`, `components`, `directives`, `helpers` and
   `mixins`. Start using them via `import { composables, utils, ... } from
   'owncloud-design-system'`.

   https://github.com/owncloud/owncloud-design-system/pull/2048


* Enhancement - Hover in ocDrop menues: [#2069](https://github.com/owncloud/owncloud-design-system/pull/2069)

   We've added the "oc-menu-item-hover" class for <li> elements inside ocDrop, to add the hover
   effect on buttons and links.

   https://github.com/owncloud/owncloud-design-system/pull/2069


* Enhancement - OcModal add checkbox and secondary button: [#6994](https://github.com/owncloud/web/pull/6994)

   We've added an optional checkbox and secondary button to the OcModal

   https://github.com/owncloud/web/issues/6996
   https://github.com/owncloud/web/pull/6994


* Enhancement - OcModal input type: [#2077](https://github.com/owncloud/owncloud-design-system/pull/2077)

   We've added an option to set the input type for input fields in the OcModal component.

   https://github.com/owncloud/owncloud-design-system/pull/2077


* Enhancement - Make OcResource inline-flex: [#2041](https://github.com/owncloud/owncloud-design-system/pull/2041)

   We've changed OcResource's display CSS attribute to inline-flex to prevent a line break

   https://github.com/owncloud/owncloud-design-system/pull/2041


* Enhancement - Add `isFileExtensionDisplayed` property: [#2087](https://github.com/owncloud/owncloud-design-system/pull/2087)

   We've added the `isFileExtensionDisplayed` property to the `OcResource` and
   `OcResourceName` components, to determine whether the file extension should be displayed or
   not.

   https://github.com/owncloud/web/issues/6730
   https://github.com/owncloud/owncloud-design-system/pull/2087


* Enhancement - Redesign OcGhostElement: [#2049](https://github.com/owncloud/owncloud-design-system/pull/2049)

   We've redesigned OcGhostElement to use OcResourceIcon and to display a better preview of the
   items that have been dragged

   https://github.com/owncloud/owncloud-design-system/pull/2049


* Enhancement - Replace deprecated String.prototype.substr(): [#2059](https://github.com/owncloud/owncloud-design-system/pull/2059)

   We've replaced all occurrences of the deprecated String.prototype.substr() function with
   String.prototype.slice() which works similarly but isn't deprecated.

   https://github.com/owncloud/owncloud-design-system/pull/2059


* Enhancement - Add HTML title to the resourceName component: [#2164](https://github.com/owncloud/owncloud-design-system/pull/2164)

   We've added the HTML title attribute to the resourceName component in case no tooltip is being
   displayed.

   https://github.com/owncloud/owncloud-design-system/pull/2164


* Enhancement - Add option to not truncate the resource name: [#2157](https://github.com/owncloud/owncloud-design-system/pull/2157)

   We've added a new property to the resourceName component that indicates whether a resource
   name should be truncated or not.

   https://github.com/owncloud/owncloud-design-system/pull/2157

# Changelog for [13.0.0] (2022-03-23)

The following sections list the changes in ownCloud Design System 13.0.0.

[13.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v12.2.2...v13.0.0

## Summary

* Bugfix - Disabled OcSelect background: [#2008](https://github.com/owncloud/owncloud-design-system/pull/2008)
* Bugfix - Icons/Thumbnails were only visible for clickable resources: [#2007](https://github.com/owncloud/owncloud-design-system/pull/2007)
* Bugfix - OcSelect transparent background: [#2036](https://github.com/owncloud/owncloud-design-system/pull/2036)
* Change - Default type of OcButton: [#2009](https://github.com/owncloud/owncloud-design-system/pull/2009)
* Change - Remove OcStatusIndicators from OcResource: [#2014](https://github.com/owncloud/owncloud-design-system/pull/2014)
* Enhancement - Configurable OcResource parentfolder name: [#2029](https://github.com/owncloud/owncloud-design-system/pull/2029)
* Enhancement - Icons for drawio, ifc and odg resource types: [#2005](https://github.com/owncloud/owncloud-design-system/pull/2005)
* Enhancement - Polish OcSwitch: [#2018](https://github.com/owncloud/owncloud-design-system/pull/2018)
* Enhancement - Make filled primary OcButton use gradient background: [#2036](https://github.com/owncloud/owncloud-design-system/pull/2036)
* Enhancement - Redesign OcStatusIndicators: [#2014](https://github.com/owncloud/owncloud-design-system/pull/2014)
* Enhancement - Underline OcResourceName: [#2019](https://github.com/owncloud/owncloud-design-system/pull/2019)
* Enhancement - Apply size property to oc-tag: [#2011](https://github.com/owncloud/owncloud-design-system/pull/2011)

## Details

* Bugfix - Disabled OcSelect background: [#2008](https://github.com/owncloud/owncloud-design-system/pull/2008)

   We've fixed a bug that prevented the background of disabled OcSelect components from using
   theming colors.

   https://github.com/owncloud/owncloud-design-system/pull/2008


* Bugfix - Icons/Thumbnails were only visible for clickable resources: [#2007](https://github.com/owncloud/owncloud-design-system/pull/2007)

   We fixed that only clickable resources had icons/thumbnails in `OcResource`. It was fixed by
   introducing an `OcResourceLink` component that reduces code complexity and duplication
   when linking resources.

   https://github.com/owncloud/owncloud-design-system/pull/2007


* Bugfix - OcSelect transparent background: [#2036](https://github.com/owncloud/owncloud-design-system/pull/2036)

   We fixed a non-transparent background in the OcSelect button, leading to visual glitches.

   https://github.com/owncloud/owncloud-design-system/issues/2030
   https://github.com/owncloud/owncloud-design-system/pull/2036


* Change - Default type of OcButton: [#2009](https://github.com/owncloud/owncloud-design-system/pull/2009)

   We've changed the default type of buttons rendered by `OcButton` to `button`. Browsers
   otherwise assume they are of type `submit` which leads to very unexpected behavior in forms,
   especially as we use `OcButton` in a lot of (not so obvious) places for a11y reasons.

   https://github.com/owncloud/owncloud-design-system/pull/2009


* Change - Remove OcStatusIndicators from OcResource: [#2014](https://github.com/owncloud/owncloud-design-system/pull/2014)

   We've removed OcStatusIndicators from OcResource since it will be moved in a separate column

   https://github.com/owncloud/web/issues/5976
   https://github.com/owncloud/owncloud-design-system/pull/2014
   https://github.com/owncloud/web/pull/6552


* Enhancement - Configurable OcResource parentfolder name: [#2029](https://github.com/owncloud/owncloud-design-system/pull/2029)

   We've added a `parent-folder-name-default` property to the OcResource component. Before,
   an empty parent resulted in a hardcoded "All files and folders" which becomes misleading with
   the introduction of spaces in oCIS.

   https://github.com/owncloud/owncloud-design-system/pull/2029


* Enhancement - Icons for drawio, ifc and odg resource types: [#2005](https://github.com/owncloud/owncloud-design-system/pull/2005)

   We've added resource type extension mapping and icons for the drawio, ifc, ipynb and odg file
   extensions.

   https://github.com/owncloud/web/issues/6416
   https://github.com/owncloud/owncloud-design-system/pull/2005


* Enhancement - Polish OcSwitch: [#2018](https://github.com/owncloud/owncloud-design-system/pull/2018)

   We've adjusted the OcSwitch to fit the redesign

   https://github.com/owncloud/web/issues/6492
   https://github.com/owncloud/owncloud-design-system/pull/2018


* Enhancement - Make filled primary OcButton use gradient background: [#2036](https://github.com/owncloud/owncloud-design-system/pull/2036)

   We've updated the OcButton to use the gradient background color when used in its `filled`
   appearance.

   https://github.com/owncloud/owncloud-design-system/issues/1952
   https://github.com/owncloud/owncloud-design-system/pull/2036
   https://github.com/owncloud/owncloud-design-system/pull/2038


* Enhancement - Redesign OcStatusIndicators: [#2014](https://github.com/owncloud/owncloud-design-system/pull/2014)

   We've redesigned the share/status indicators to fit the new design in web.

   https://github.com/owncloud/web/issues/5976
   https://github.com/owncloud/owncloud-design-system/pull/2014
   https://github.com/owncloud/web/pull/6552


* Enhancement - Underline OcResourceName: [#2019](https://github.com/owncloud/owncloud-design-system/pull/2019)

   We've added an underline on hover effect to OcResourceName

   https://github.com/owncloud/web/issues/6492
   https://github.com/owncloud/owncloud-design-system/pull/2019


* Enhancement - Apply size property to oc-tag: [#2011](https://github.com/owncloud/owncloud-design-system/pull/2011)

   We've added a size property to oc-tag

   https://github.com/owncloud/owncloud-design-system/pull/2011

# Changelog for [12.2.2] (2022-03-03)

The following sections list the changes in ownCloud Design System 12.2.2.

[12.2.2]: https://github.com/owncloud/owncloud-design-system/compare/v12.2.1...v12.2.2

## Summary

* Bugfix - Initial focus in OcModal: [#2001](https://github.com/owncloud/owncloud-design-system/pull/2001)
* Bugfix - Hidden overflow on OcModal: [#2000](https://github.com/owncloud/owncloud-design-system/pull/2000)

## Details

* Bugfix - Initial focus in OcModal: [#2001](https://github.com/owncloud/owncloud-design-system/pull/2001)

   We've fixed a bug that was introduced in the last version, where the initial focus element
   provided via the focusTrapInitial property was broken.

   https://github.com/owncloud/owncloud-design-system/pull/2001


* Bugfix - Hidden overflow on OcModal: [#2000](https://github.com/owncloud/owncloud-design-system/pull/2000)

   We've fixed a bug that prevented overflow within the OcModal to e.g. use dropdowns inside the
   modal body.

   https://github.com/owncloud/owncloud-design-system/pull/2000

# Changelog for [12.2.1] (2022-03-02)

The following sections list the changes in ownCloud Design System 12.2.1.

[12.2.1]: https://github.com/owncloud/owncloud-design-system/compare/v12.2.0...v12.2.1

## Summary

* Bugfix - OcTable Sort: [#1996](https://github.com/owncloud/owncloud-design-system/pull/1996)
* Bugfix - Icon & background color for mobile OcBreadcrumb: [#1980](https://github.com/owncloud/owncloud-design-system/issues/1980)
* Bugfix - Initial focus in OcModal: [#1995](https://github.com/owncloud/owncloud-design-system/pull/1995)

## Details

* Bugfix - OcTable Sort: [#1996](https://github.com/owncloud/owncloud-design-system/pull/1996)

   We've fixed a visual bug that caused the field sort arrow to be always visible.

   https://github.com/owncloud/owncloud-design-system/pull/1996


* Bugfix - Icon & background color for mobile OcBreadcrumb: [#1980](https://github.com/owncloud/owncloud-design-system/issues/1980)

   The icon in the mobile OcBreadcrumb drop was missing, also the background color made reading
   the content impossible. This has been adressed.

   https://github.com/owncloud/owncloud-design-system/issues/1980
   https://github.com/owncloud/owncloud-design-system/pull/1994


* Bugfix - Initial focus in OcModal: [#1995](https://github.com/owncloud/owncloud-design-system/pull/1995)

   We've fixed a bug that was introduced in the last version, where the initial focus of modals with
   text fields was broken.

   https://github.com/owncloud/owncloud-design-system/pull/1995

# Changelog for [12.2.0] (2022-02-28)

The following sections list the changes in ownCloud Design System 12.2.0.

[12.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v12.1.0...v12.2.0

## Summary

* Enhancement - Apply outstanding background color to oc-card: [#1974](https://github.com/owncloud/owncloud-design-system/pull/1974)
* Enhancement - Extend modal focus-trap functionality: [#1987](https://github.com/owncloud/owncloud-design-system/pull/1987)
* Enhancement - Redesign OcBreadcrumb: [#6218](https://github.com/owncloud/web/issues/6218)
* Enhancement - Redesign files table related components: [#1958](https://github.com/owncloud/owncloud-design-system/pull/1958)

## Details

* Enhancement - Apply outstanding background color to oc-card: [#1974](https://github.com/owncloud/owncloud-design-system/pull/1974)

   We've adjusted he background color to oc-card to have an outstanding look

   https://github.com/owncloud/owncloud-design-system/pull/1974


* Enhancement - Extend modal focus-trap functionality: [#1987](https://github.com/owncloud/owncloud-design-system/pull/1987)

   We've added the option to define which child element should be focused initially.

   https://github.com/owncloud/owncloud-design-system/pull/1987


* Enhancement - Redesign OcBreadcrumb: [#6218](https://github.com/owncloud/web/issues/6218)

   We've adjustet the look of the OcBreadcrumb to fit the Redesign

   https://github.com/owncloud/web/issues/6218
   https://github.com/owncloud/owncloud-design-system/pull/1975
   https://github.com/owncloud/owncloud-design-system/pull/1982


* Enhancement - Redesign files table related components: [#1958](https://github.com/owncloud/owncloud-design-system/pull/1958)

   We've adjusted OcTable, OcResource, OcDrop and OcCheckbox to fit the redesign.

   https://github.com/owncloud/web/issues/6207
   https://github.com/owncloud/owncloud-design-system/pull/1958
   https://github.com/owncloud/owncloud-design-system/pull/1978
   https://github.com/owncloud/owncloud-design-system/pull/1988

# Changelog for [12.1.0] (2022-02-10)

The following sections list the changes in ownCloud Design System 12.1.0.

[12.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v12.0.0...v12.1.0

## Summary

* Bugfix - Background-primary-gradient border: [#6383](https://github.com/owncloud/web/issues/6383)
* Enhancement - Redesign OcModal: [#1953](https://github.com/owncloud/owncloud-design-system/pull/1953)

## Details

* Bugfix - Background-primary-gradient border: [#6383](https://github.com/owncloud/web/issues/6383)

   The `.oc-background-primary-gradient` class was setting the CSS border property instead of
   only the border-color, hereby removing any border-width the target tag might already carry.
   This lead to some very small differences when rendering buttons and has been resolved now.

   https://github.com/owncloud/web/issues/6383
   https://github.com/owncloud/owncloud-design-system/pull/1945


* Enhancement - Redesign OcModal: [#1953](https://github.com/owncloud/owncloud-design-system/pull/1953)

   We have redesigned the OcModal to suit better to dark mode UIs, to have a nicer background
   overlay, to always have a primary confirm button and fixed some icon bugs in its code examples.

   https://github.com/owncloud/owncloud-design-system/pull/1953

# Changelog for [12.0.0] (2022-02-07)

The following sections list the changes in ownCloud Design System 12.0.0.

[12.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v11.3.0...v12.0.0

## Summary

* Bugfix - Missing OcDrop shadow: [#1926](https://github.com/owncloud/owncloud-design-system/pull/1926)
* Bugfix - OcNotification positioning: [#1658](https://github.com/owncloud/owncloud-design-system/pull/1658)
* Bugfix - Rename GhostElement: [#1845](https://github.com/owncloud/owncloud-design-system/pull/1845)
* Bugfix - OcTooltip isn't reactive: [#1863](https://github.com/owncloud/owncloud-design-system/pull/1863)
* Change - Drop Internet Explorer support: [#1909](https://github.com/owncloud/owncloud-design-system/pull/1909)
* Change - Do not sort in OcTable: [#1825](https://github.com/owncloud/owncloud-design-system/pull/1825)
* Change - Pass folderLink to OcResource component: [#1913](https://github.com/owncloud/owncloud-design-system/pull/1913)
* Change - Remove OcAppSideBar component: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810)
* Change - Remove OcAppBar component: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810)
* Change - Remove implicit ODS registration: [#1848](https://github.com/owncloud/owncloud-design-system/pull/1848)
* Change - Remove oc-table-files from ods: [#1817](https://github.com/owncloud/owncloud-design-system/pull/1817)
* Change - Remove OcGrid options: [#1658](https://github.com/owncloud/owncloud-design-system/pull/1658)
* Change - Move OcSidebarNav and OcSidebarNavItem to web: [#6036](https://github.com/owncloud/web/issues/6036)
* Change - Remove UiKit: [#1658](https://github.com/owncloud/owncloud-design-system/pull/1658)
* Change - Remove unused props for unstyled components: [#1795](https://github.com/owncloud/owncloud-design-system/pull/1795)
* Change - Use remixicons for redesign: [#1826](https://github.com/owncloud/owncloud-design-system/pull/1826)
* Enhancement - Make Vue-Composition-API available: [#1848](https://github.com/owncloud/owncloud-design-system/pull/1848)
* Enhancement - Export mappings of types, icons and colors of resources: [#1920](https://github.com/owncloud/owncloud-design-system/pull/1920)
* Enhancement - Fix OcAvatar line-height: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810)
* Enhancement - Add option to render table cells lazy: [#1848](https://github.com/owncloud/owncloud-design-system/pull/1848)
* Enhancement - Make OcDrop rounded: [#1881](https://github.com/owncloud/owncloud-design-system/pull/1881)
* Enhancement - Change background color of OcDrop: [#1919](https://github.com/owncloud/owncloud-design-system/pull/1919)
* Enhancement - Improve OcList: [#1881](https://github.com/owncloud/owncloud-design-system/pull/1881)
* Enhancement - Show path / parent folder to distinguish files: [#5953](https://github.com/owncloud/web/issues/5953)
* Enhancement - Redesign Filetype icons: [#6278](https://github.com/owncloud/web/issues/6278)
* Enhancement - Adjust OcSearchBar to new design: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810/)
* Enhancement - Sizes: [#1858](https://github.com/owncloud/owncloud-design-system/pull/1858)
* Enhancement - Add svg icon for spaces: [#1846](https://github.com/owncloud/owncloud-design-system/pull/1846)
* Enhancement - OcTable header alignment: [#1922](https://github.com/owncloud/owncloud-design-system/pull/1922)
* Enhancement - Use Roboto font: [#1876](https://github.com/owncloud/owncloud-design-system/pull/1876)

## Details

* Bugfix - Missing OcDrop shadow: [#1926](https://github.com/owncloud/owncloud-design-system/pull/1926)

   In certain situations, other DOM elements made the OcDrop shadow invisible. This has been
   resolved.

   https://github.com/owncloud/owncloud-design-system/pull/1926
   https://github.com/owncloud/owncloud-design-system/pull/1931


* Bugfix - OcNotification positioning: [#1658](https://github.com/owncloud/owncloud-design-system/pull/1658)

   We have taken care of the positioning in the OcNotification component, which didn't work as
   expected. Notifications can now be displayed on the left or right side or centered.

   https://github.com/owncloud/owncloud-design-system/pull/1658


* Bugfix - Rename GhostElement: [#1845](https://github.com/owncloud/owncloud-design-system/pull/1845)

   We've renamed the GhostElement component to OcGhostElement to follow the atomic design
   principles.

   https://github.com/owncloud/owncloud-design-system/pull/1845


* Bugfix - OcTooltip isn't reactive: [#1863](https://github.com/owncloud/owncloud-design-system/pull/1863)

   We've fixed an issue with the OcTooltip when the tooltip gets hidden by reactivity and then
   reactivated with some value.

   https://github.com/owncloud/owncloud-design-system/pull/1863


* Change - Drop Internet Explorer support: [#1909](https://github.com/owncloud/owncloud-design-system/pull/1909)

   Since it's nearing its end-of-life, we've dropped polyfills for IE in favor of a smaller bundle
   size.

   https://github.com/owncloud/owncloud-design-system/pull/1909


* Change - Do not sort in OcTable: [#1825](https://github.com/owncloud/owncloud-design-system/pull/1825)

   We removed sorting from OcTable and added a `sort` event instead, which can be used to sort the
   data externally. This is crucial for server side sorting and handling pagination correctly.

   https://github.com/owncloud/web/issues/5687
   https://github.com/owncloud/owncloud-design-system/pull/1825
   https://github.com/owncloud/owncloud-design-system/pull/1839
   https://github.com/owncloud/web/pull/6136


* Change - Pass folderLink to OcResource component: [#1913](https://github.com/owncloud/owncloud-design-system/pull/1913)

   For more flexibility, the folderLink needs to be passed to the OcResource component

   https://github.com/owncloud/owncloud-design-system/pull/1913


* Change - Remove OcAppSideBar component: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810)

   We have removed the OcAppSideBar component since it's not actively used anywhere.

   https://github.com/owncloud/owncloud-design-system/pull/1810


* Change - Remove OcAppBar component: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810)

   We have removed the OcAppBar component since it's not actively used anywhere and broken.

   https://github.com/owncloud/owncloud-design-system/pull/1810


* Change - Remove implicit ODS registration: [#1848](https://github.com/owncloud/owncloud-design-system/pull/1848)

   Remove implicit registration of ODS, from now on applications using ODS must register it
   explicit via `Vue.use`.

   https://github.com/owncloud/owncloud-design-system/pull/1848


* Change - Remove oc-table-files from ods: [#1817](https://github.com/owncloud/owncloud-design-system/pull/1817)

   Ods oc-table-files always contained concrete web-app-files logic, to make development more
   agile and keep things close oc-table-files was removed from ODS and relocated to live in
   web-app-files from now on.

   https://github.com/owncloud/owncloud-design-system/pull/1817
   https://github.com/owncloud/web/pull/6106


* Change - Remove OcGrid options: [#1658](https://github.com/owncloud/owncloud-design-system/pull/1658)

   We have removed the `match` and `childWidth` option in the OcGrid component since they were
   unused and relied heavily on the removed UiKit library.

   https://github.com/owncloud/owncloud-design-system/pull/1658


* Change - Move OcSidebarNav and OcSidebarNavItem to web: [#6036](https://github.com/owncloud/web/issues/6036)

   We've moved OcSidebarNav and OcSidebarNavItem to web and renamed it to SidebarNav and
   SidebarNavItem.

   https://github.com/owncloud/web/issues/6036
   https://github.com/owncloud/owncloud-design-system/pull/1810


* Change - Remove UiKit: [#1658](https://github.com/owncloud/owncloud-design-system/pull/1658)

   We have removed the UiKit library this design system originally was built on. The necessary
   style rules for the design system itself and our web repository have been internalized, and
   everything else got dropped to greatly reduce bundle size and build times.

   Please note that with this change, we have also dropped and/or refactored a lot of CSS classes
   and correlated styling.

   https://github.com/owncloud/owncloud-design-system/issues/538
   https://github.com/owncloud/owncloud-design-system/pull/1658
   https://github.com/owncloud/owncloud-design-system/pull/1882


* Change - Remove unused props for unstyled components: [#1795](https://github.com/owncloud/owncloud-design-system/pull/1795)

   We removed the `stopClassProgragation` property in two components, which resulted in
   unstyled components before but was unused.

   https://github.com/owncloud/owncloud-design-system/pull/1795


* Change - Use remixicons for redesign: [#1826](https://github.com/owncloud/owncloud-design-system/pull/1826)

   We've switched the iconset to remixicons to fit the new design.

   https://github.com/owncloud/web/issues/6100
   https://github.com/owncloud/owncloud-design-system/pull/1826
   https://github.com/owncloud/owncloud-design-system/pull/1853


* Enhancement - Make Vue-Composition-API available: [#1848](https://github.com/owncloud/owncloud-design-system/pull/1848)

   To support upcoming Vue composition-api we`ve added the compatibility layer from the
   creators. From now on all features described here
   `https://github.com/vuejs/composition-api` can be used.

   https://github.com/owncloud/owncloud-design-system/pull/1848


* Enhancement - Export mappings of types, icons and colors of resources: [#1920](https://github.com/owncloud/owncloud-design-system/pull/1920)

   The bundled design system now contains two json files that map file extensions to their
   dedicated resource icon and the corresponding color token.

   https://github.com/owncloud/owncloud-design-system/pull/1920


* Enhancement - Fix OcAvatar line-height: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810)

   We've fixed an visual bug that caused the OcAvatar to be positioned 1px too high

   https://github.com/owncloud/owncloud-design-system/pull/1810


* Enhancement - Add option to render table cells lazy: [#1848](https://github.com/owncloud/owncloud-design-system/pull/1848)

   In cases where an instance of `OcTable` has many child rows and cells, it can be a bottleneck to
   render all of them immediately. With this in mind we've added the lazy option to the table fields
   object where the consuming app can decide how lazy rendering should behave.

   By default lazy cell rendering is disabled, to enable it add a lazy object to the field.

   Following options are available: * `delay: 250` - when the cell visibility on screen is below
   the given milliseconds value rendering gets skipped. * `mode: show` - cell gets rendered and
   stays painted, no de-rendering happens. * `mode: showHide` - cell gets rendered when it enters
   the screen and de-rendered when its off. * `rootMargin: 100px` - given value will be added to the
   outer area of the element which then increases the visibility detection radius

   https://github.com/owncloud/owncloud-design-system/pull/1848


* Enhancement - Make OcDrop rounded: [#1881](https://github.com/owncloud/owncloud-design-system/pull/1881)

   We've added a border-radius for the OcDrop to make it rounded.

   https://github.com/owncloud/owncloud-design-system/pull/1881


* Enhancement - Change background color of OcDrop: [#1919](https://github.com/owncloud/owncloud-design-system/pull/1919)

   The OcDrop component now has a different background color than the default color to improve the
   contrast.

   https://github.com/owncloud/owncloud-design-system/pull/1919
   https://github.com/owncloud/owncloud-design-system/pull/1923


* Enhancement - Improve OcList: [#1881](https://github.com/owncloud/owncloud-design-system/pull/1881)

   We've fixed some styling and added a 'raw' property in OcList to disable list styling.

   https://github.com/owncloud/owncloud-design-system/pull/1881


* Enhancement - Show path / parent folder to distinguish files: [#5953](https://github.com/owncloud/web/issues/5953)

   We've added an option to show the path / parent folder under the resource name in order to
   distinguish files better

   https://github.com/owncloud/web/issues/5953
   https://github.com/owncloud/owncloud-design-system/pull/1860
   https://github.com/owncloud/owncloud-design-system/pull/1871


* Enhancement - Redesign Filetype icons: [#6278](https://github.com/owncloud/web/issues/6278)

   We've adjusted the resource icons to fit the redesign. We've added a new OcResourceIcon
   component that is themable.

   https://github.com/owncloud/web/issues/6278
   https://github.com/owncloud/owncloud-design-system/pull/1900
   https://github.com/owncloud/owncloud-design-system/pull/1924
   https://github.com/owncloud/owncloud-design-system/pull/1925
   https://github.com/owncloud/owncloud-design-system/pull/1934


* Enhancement - Adjust OcSearchBar to new design: [#1810](https://github.com/owncloud/owncloud-design-system/pull/1810/)

   We've redesigned the OcSearchBar to fit the new design.

   https://github.com/owncloud/web/issues/6036
   https://github.com/owncloud/owncloud-design-system/pull/1810/


* Enhancement - Sizes: [#1858](https://github.com/owncloud/owncloud-design-system/pull/1858)

   The size variables which define margins and paddings have been changed to use multiples of 8
   instead of 10.

   https://github.com/owncloud/owncloud-design-system/pull/1858


* Enhancement - Add svg icon for spaces: [#1846](https://github.com/owncloud/owncloud-design-system/pull/1846)

   https://github.com/owncloud/owncloud-design-system/pull/1846


* Enhancement - OcTable header alignment: [#1922](https://github.com/owncloud/owncloud-design-system/pull/1922)

   We've applied the full row height to the table header and made it vertically centered to give a
   more fluffy experience for the eye.

   https://github.com/owncloud/owncloud-design-system/pull/1922


* Enhancement - Use Roboto font: [#1876](https://github.com/owncloud/owncloud-design-system/pull/1876)

   We've switched the default font from Fira Sans to Roboto.

   https://github.com/owncloud/web/issues/6100/
   https://github.com/owncloud/owncloud-design-system/pull/1876

# Changelog for [11.3.0] (2021-12-03)

The following sections list the changes in ownCloud Design System 11.3.0.

[11.3.0]: https://github.com/owncloud/owncloud-design-system/compare/v11.3.1...v11.3.0

## Summary

* Bugfix - Set language for date formatting: [#1806](https://github.com/owncloud/owncloud-design-system/pull/1806)
* Enhancement - Relative date tooltips in the OcTableFiles component: [#1787](https://github.com/owncloud/owncloud-design-system/pull/1787)
* Enhancement - Breadcrumb contextmenu: [#6030](https://github.com/owncloud/web/issues/6030)
* Enhancement - Optional padding size for OcDrop: [#1798](https://github.com/owncloud/owncloud-design-system/pull/1798)
* Enhancement - Truncate file names while preserving file extensions: [#1758](https://github.com/owncloud/owncloud-design-system/issues/1758)

## Details

* Bugfix - Set language for date formatting: [#1806](https://github.com/owncloud/owncloud-design-system/pull/1806)

   We're now setting the language when formatting dates, so that localized parts of the date/time
   get shown according to the respective language.

   https://github.com/owncloud/owncloud-design-system/pull/1806


* Enhancement - Relative date tooltips in the OcTableFiles component: [#1787](https://github.com/owncloud/owncloud-design-system/pull/1787)

   Relative dates like "1 day ago" now have a tooltip that shows the absolute date. The following
   dates in the files table are affected:

   * modify date * share date * delete date

   https://github.com/owncloud/web/issues/5672
   https://github.com/owncloud/owncloud-design-system/pull/1787


* Enhancement - Breadcrumb contextmenu: [#6030](https://github.com/owncloud/web/issues/6030)

   We've added a button to the last item of the OcBreadcrumb component that triggers a dropdown
   which can be customizably filled.

   https://github.com/owncloud/web/issues/6030
   https://github.com/owncloud/owncloud-design-system/pull/1786


* Enhancement - Optional padding size for OcDrop: [#1798](https://github.com/owncloud/owncloud-design-system/pull/1798)

   We've added a `paddingSize` property to the OcDrop component for specifying a dedicated
   padding of the div that wraps the content slot. It defaults to "medium" and also allows to remove
   the padding entirely.

   https://github.com/owncloud/owncloud-design-system/pull/1798


* Enhancement - Truncate file names while preserving file extensions: [#1758](https://github.com/owncloud/owncloud-design-system/issues/1758)

   File names that exceed the horizontal space of a file list now get truncated while preserving
   their extensions. This way, the user can quickly identify the type of a file.

   https://github.com/owncloud/owncloud-design-system/issues/1758
   https://github.com/owncloud/owncloud-design-system/pull/1782

# Changelog for [11.3.1] (2021-12-03)

The following sections list the changes in ownCloud Design System 11.3.1.

[11.3.1]: https://github.com/owncloud/owncloud-design-system/compare/v11.2.2...v11.3.1

## Summary

* Bugfix - Padding in breadcrumb context menu: [#1813](https://github.com/owncloud/owncloud-design-system/pull/1813)

## Details

* Bugfix - Padding in breadcrumb context menu: [#1813](https://github.com/owncloud/owncloud-design-system/pull/1813)

   We've removed the padding from the context menu in the breadcrumbs to align its visual
   appearance with the context menu in the files table.

   https://github.com/owncloud/owncloud-design-system/pull/1813

# Changelog for [11.2.2] (2021-11-11)

The following sections list the changes in ownCloud Design System 11.2.2.

[11.2.2]: https://github.com/owncloud/owncloud-design-system/compare/v11.2.1...v11.2.2

## Summary

* Bugfix - Fix extension icon rendering: [#1779](https://github.com/owncloud/owncloud-design-system/pull/1779)

## Details

* Bugfix - Fix extension icon rendering: [#1779](https://github.com/owncloud/owncloud-design-system/pull/1779)

   The `extension` icon wouldn't render in the ownCloud web. We fixed the SVG file to prevent this
   problem.

   https://github.com/owncloud/owncloud-design-system/pull/1779

# Changelog for [11.2.1] (2021-11-10)

The following sections list the changes in ownCloud Design System 11.2.1.

[11.2.1]: https://github.com/owncloud/owncloud-design-system/compare/v11.2.0...v11.2.1

## Summary

* Bugfix - Fix files table event: [#1777](https://github.com/owncloud/owncloud-design-system/pull/1777)

## Details

* Bugfix - Fix files table event: [#1777](https://github.com/owncloud/owncloud-design-system/pull/1777)

   In v11.2.0 we refactored the naming of the table drag and drop events. This introduced a naming
   clash in the files table events which is fixed now.

   https://github.com/owncloud/owncloud-design-system/pull/1777

# Changelog for [11.2.0] (2021-11-09)

The following sections list the changes in ownCloud Design System 11.2.0.

[11.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v11.1.0...v11.2.0

## Summary

* Enhancement - Add accentuated class for OcTable: [#5967](https://github.com/owncloud/web/issues/5967)
* Enhancement - Add Ghost Element for Drag & Drop: [#5788](https://github.com/owncloud/web/issues/5788)
* Enhancement - Add "extension" svg icon: [#1771](https://github.com/owncloud/owncloud-design-system/pull/1771)
* Enhancement - Add closure to mutate resource dom selector: [#1766](https://github.com/owncloud/owncloud-design-system/pull/1766)
* Enhancement - Reduce filename text weight: [#1759](https://github.com/owncloud/owncloud-design-system/pull/1759)

## Details

* Enhancement - Add accentuated class for OcTable: [#5967](https://github.com/owncloud/web/issues/5967)

   We've added an accentuated class for the OcTable for accentuating certain files e.g. after
   uploading, in a different way than already highlighted resources.

   https://github.com/owncloud/web/issues/5967
   https://github.com/owncloud/owncloud-design-system/pull/1763


* Enhancement - Add Ghost Element for Drag & Drop: [#5788](https://github.com/owncloud/web/issues/5788)

   We've added a custom ghost element for the drag&drop functionality in the OcTableFiles to
   better visualize the action to the user.

   https://github.com/owncloud/web/issues/5788
   https://github.com/owncloud/owncloud-design-system/pull/1754


* Enhancement - Add "extension" svg icon: [#1771](https://github.com/owncloud/owncloud-design-system/pull/1771)

   https://github.com/owncloud/owncloud-design-system/pull/1771


* Enhancement - Add closure to mutate resource dom selector: [#1766](https://github.com/owncloud/owncloud-design-system/pull/1766)

   In some cases the resource or item id is not a valid dom selector, we've introduced an optional
   closure property for oc-table and oc-table-files to customize the selector to its need.

   https://github.com/owncloud/owncloud-design-system/pull/1766


* Enhancement - Reduce filename text weight: [#1759](https://github.com/owncloud/owncloud-design-system/pull/1759)

   We reduced the font weight for filenames to make the files list easier for the eye.

   https://github.com/owncloud/owncloud-design-system/pull/1759

# Changelog for [11.1.0] (2021-11-03)

The following sections list the changes in ownCloud Design System 11.1.0.

[11.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v11.0.0...v11.1.0

## Summary

* Bugfix - Limit select event in OcTableFiles: [#1753](https://github.com/owncloud/owncloud-design-system/pull/1753)
* Bugfix - Add word-break rule to OcNotificationMessage component: [#1712](https://github.com/owncloud/owncloud-design-system/issues/1712)
* Bugfix - OcTable sorting case sensitivity: [#1698](https://github.com/owncloud/owncloud-design-system/issues/1698)
* Bugfix - Drag and Drop triggers wrong actions: [#5808](https://github.com/owncloud/web/issues/5808)
* Enhancement - Make OcDatepicker themable: [#1679](https://github.com/owncloud/owncloud-design-system/issues/1679)
* Enhancement - Streamline OcTextInput: [#1636](https://github.com/owncloud/owncloud-design-system/pull/1636)

## Details

* Bugfix - Limit select event in OcTableFiles: [#1753](https://github.com/owncloud/owncloud-design-system/pull/1753)

   We've updated the behaviour in the OcTableFiles component so select events only get fired if
   the target row is not already selected. This allows for batch actions in the contextmenu which
   was limited to only one resource before.

   https://github.com/owncloud/owncloud-design-system/pull/1753


* Bugfix - Add word-break rule to OcNotificationMessage component: [#1712](https://github.com/owncloud/owncloud-design-system/issues/1712)

   We had issues with long filenames overflowing the OcNotificationMessage body.

   https://github.com/owncloud/owncloud-design-system/issues/1712
   https://github.com/owncloud/owncloud-design-system/pulls/1736


* Bugfix - OcTable sorting case sensitivity: [#1698](https://github.com/owncloud/owncloud-design-system/issues/1698)

   We've addressed the problem when sorting upper and lowercase which resulted in both being
   sorted separately.

   https://github.com/owncloud/owncloud-design-system/issues/1698
   https://github.com/owncloud/owncloud-design-system/pull/1735


* Bugfix - Drag and Drop triggers wrong actions: [#5808](https://github.com/owncloud/web/issues/5808)

   We addressed the problem when a user tries to upload a file via drag and drop which falsely
   triggered the drag and drop move in the files table.

   https://github.com/owncloud/web/issues/5808
   https://github.com/owncloud/owncloud-design-system/pull/1732


* Enhancement - Make OcDatepicker themable: [#1679](https://github.com/owncloud/owncloud-design-system/issues/1679)

   We've added styling to the OcDatepicker to use the colors set by the used theme.

   https://github.com/owncloud/owncloud-design-system/issues/1679
   https://github.com/owncloud/owncloud-design-system/pull/1740


* Enhancement - Streamline OcTextInput: [#1636](https://github.com/owncloud/owncloud-design-system/pull/1636)

   We have updated the OcTextInput component to be in line with other form components in the design
   system.

   - Fixed the clear button being visible even when element is disabled - Made clear button emit
   null instead of empty string - Made input/change event handling a bit more consistent - Added
   defaultValue prop that is for now passed as placeholder to the input field

   https://github.com/owncloud/owncloud-design-system/pull/1636

# Changelog for [11.0.0] (2021-10-04)

The following sections list the changes in ownCloud Design System 11.0.0.

[11.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v10.0.0...v11.0.0

## Summary

* Bugfix - Prevent contextmenu display issues within OcTableFiles: [#1691](https://github.com/owncloud/owncloud-design-system/pull/1691)
* Bugfix - Prevent hover style on footer of <oc-table>: [#1667](https://github.com/owncloud/owncloud-design-system/pull/1667)
* Change - Replace vue-datetime with v-calendar in our datepicker component: [#1661](https://github.com/owncloud/owncloud-design-system/pull/1661)
* Enhancement - Allow hover option in OcTableFiles: [#1632](https://github.com/owncloud/owncloud-design-system/pull/1632)

## Details

* Bugfix - Prevent contextmenu display issues within OcTableFiles: [#1691](https://github.com/owncloud/owncloud-design-system/pull/1691)

   Context menu for files table now detects if available space is enough to show all items. If not,
   it automatically calculates the height and adds scroll bars to the menu.

   https://github.com/owncloud/web/issues/5845
   https://github.com/owncloud/owncloud-design-system/pull/1691


* Bugfix - Prevent hover style on footer of <oc-table>: [#1667](https://github.com/owncloud/owncloud-design-system/pull/1667)

   https://github.com/owncloud/owncloud-design-system/pull/1667


* Change - Replace vue-datetime with v-calendar in our datepicker component: [#1661](https://github.com/owncloud/owncloud-design-system/pull/1661)

   We've added `v-calendar` dependency so that we can use it as our datepicker component. By doing
   so, we removed the old datepicker library that we use, `vue-datetime`.

   https://github.com/owncloud/owncloud-design-system/pull/1661


* Enhancement - Allow hover option in OcTableFiles: [#1632](https://github.com/owncloud/owncloud-design-system/pull/1632)

   We've added the possibility to use the `hover` option from OcTable also in OcTableFiles.

   https://github.com/owncloud/owncloud-design-system/pull/1632
   https://github.com/owncloud/owncloud-design-system/pull/1680

# Changelog for [10.0.0] (2021-09-13)

The following sections list the changes in ownCloud Design System 10.0.0.

[10.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v9.3.0...v10.0.0

## Summary

* Bugfix - Reset droptarget background color in OcTableFiles: [#1625](https://github.com/owncloud/owncloud-design-system/pull/1625)
* Change - Label for OcSelect: [#1633](https://github.com/owncloud/owncloud-design-system/pull/1633)
* Change - Use route query to store active page: [#1626](https://github.com/owncloud/owncloud-design-system/pull/1626)
* Change - Refactor OcAvatarGroup and rename to OcAvatars: [#5736](https://github.com/owncloud/web/issues/5736)
* Enhancement - Make OcAvatarX sizeable: [#1645](https://github.com/owncloud/owncloud-design-system/pull/1645)
* Enhancement - Switch filesize calculation base: [#1598](https://github.com/owncloud/owncloud-design-system/pull/1598)

## Details

* Bugfix - Reset droptarget background color in OcTableFiles: [#1625](https://github.com/owncloud/owncloud-design-system/pull/1625)

   The targeted table row in the OcTableFiles gets highlighted background when another resource
   is dragged over it for visual user feedback, but the background wasn't reset upon a successful
   drop event, which has been fixed.

   https://github.com/owncloud/owncloud-design-system/pull/1625


* Change - Label for OcSelect: [#1633](https://github.com/owncloud/owncloud-design-system/pull/1633)

   We've added a configurable `<label>` for OcSelect accessible via the `label` property. This
   shadows the `label` property of `vue-select`. Hence we introduced the `optionLabel`
   property on OcSelect which maps to the `label` property of `vue-select`.

   https://github.com/owncloud/owncloud-design-system/issues/1503
   https://github.com/owncloud/owncloud-design-system/pull/1633
   https://github.com/owncloud/owncloud-design-system/pull/1570


* Change - Use route query to store active page: [#1626](https://github.com/owncloud/owncloud-design-system/pull/1626)

   We've switched the storage of current paginated page to use route query object instead of
   param.

   https://github.com/owncloud/owncloud-design-system/pull/1626


* Change - Refactor OcAvatarGroup and rename to OcAvatars: [#5736](https://github.com/owncloud/web/issues/5736)

   We've added more share types for the OcAvatars and refactored the components used, also we
   removed the expand animation when focused or hovered.

   https://github.com/owncloud/web/issues/5736
   https://github.com/owncloud/owncloud-design-system/pull/1614
   https://github.com/owncloud/owncloud-design-system/pull/1639


* Enhancement - Make OcAvatarX sizeable: [#1645](https://github.com/owncloud/owncloud-design-system/pull/1645)

   We've added the possibility to size Avatars and gave e.g. OcAvatarGroup a width and iconSize.

   https://github.com/owncloud/web/issues/5763
   https://github.com/owncloud/owncloud-design-system/pull/1645


* Enhancement - Switch filesize calculation base: [#1598](https://github.com/owncloud/owncloud-design-system/pull/1598)

   We've switched from base-2 to base-10 when calculating the displayed file-size to align it
   better with user expectations.

   https://github.com/owncloud/owncloud-design-system/pull/1598
   https://github.com/owncloud/owncloud-design-system/pull/1615

# Changelog for [9.3.0] (2021-08-23)

The following sections list the changes in ownCloud Design System 9.3.0.

[9.3.0]: https://github.com/owncloud/owncloud-design-system/compare/v9.2.0...v9.3.0

## Summary

* Bugfix - Fix search for options provided as objects: [#1602](https://github.com/owncloud/owncloud-design-system/pull/1602)
* Bugfix - Contextmenu button triggered wrong event: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)
* Bugfix - Use pointer cursor for OcSelect actions: [#1604](https://github.com/owncloud/owncloud-design-system/pull/1604)
* Enhancement - OcTableFiles Contextmenu Tooltip: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)
* Enhancement - Highlight droptarget in OcTableFiles: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)
* Enhancement - Remove "Showdetails" button in OcTableFiles: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)

## Details

* Bugfix - Fix search for options provided as objects: [#1602](https://github.com/owncloud/owncloud-design-system/pull/1602)

   We fixed a regression that was introduced in
   https://github.com/owncloud/owncloud-design-system/pull/1521. `vue-select`
   automatically uses the property specified in `label` for filtering. When custom filtering
   based on Fuse.js was introduced that functionality got lost. Hence it was not possible to
   filter objects at all.

   https://github.com/owncloud/owncloud-design-system/pull/1602


* Bugfix - Contextmenu button triggered wrong event: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)

   The contextmenu button in the OcTableFiles triggered a `showDetails` event, which lead to the
   sidebar appearing instead of the (expected) contextmenu. This has been fixed.

   https://github.com/owncloud/owncloud-design-system/issues/1608
   https://github.com/owncloud/owncloud-design-system/pull/1610


* Bugfix - Use pointer cursor for OcSelect actions: [#1604](https://github.com/owncloud/owncloud-design-system/pull/1604)

   We changed the cursor for the actions (down/up arrows) on `OcSelect` to `pointer`. It used to be
   a `text` cursor.

   https://github.com/owncloud/owncloud-design-system/pull/1604


* Enhancement - OcTableFiles Contextmenu Tooltip: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)

   Since most of the quickactions in the OcTableFiles have a tooltip the contextmenu button
   should get one, too. It also replaces the (removed) Showdetails button and leads to better
   discoverability of the contextmenu (and therefore the sidebar).

   https://github.com/owncloud/owncloud-design-system/pull/1610


* Enhancement - Highlight droptarget in OcTableFiles: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)

   The targeted table row in the OcTableFiles now gets a highlighted background when another
   resource is dragged over it for visual user feedback.

   https://github.com/owncloud/web/issues/5705
   https://github.com/owncloud/owncloud-design-system/pull/1610


* Enhancement - Remove "Showdetails" button in OcTableFiles: [#1610](https://github.com/owncloud/owncloud-design-system/pull/1610)

   We removed the Showdetails button in the OcTableFiles quickactions to de-clutter the UI.
   Opening the sidebar is supposed to happen from the contextmenu which gets triggered by the
   three dots button.

   https://github.com/owncloud/owncloud-design-system/pull/1610

# Changelog for [9.2.0] (2021-08-18)

The following sections list the changes in ownCloud Design System 9.2.0.

[9.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v9.1.0...v9.2.0

## Summary

* Enhancement - Add select items icon: [#1597](https://github.com/owncloud/owncloud-design-system/pull/1597)

## Details

* Enhancement - Add select items icon: [#1597](https://github.com/owncloud/owncloud-design-system/pull/1597)

   We added new icon that suggests users to select items in order to enable interactions.

   https://github.com/owncloud/owncloud-design-system/pull/1597

# Changelog for [9.1.0] (2021-08-17)

The following sections list the changes in ownCloud Design System 9.1.0.

[9.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v9.0.1...v9.1.0

## Summary

* Bugfix - Initial sorting for OcTableFiles: [#1588](https://github.com/owncloud/owncloud-design-system/pull/1588)
* Bugfix - Sorting by date: [#1552](https://github.com/owncloud/owncloud-design-system/issues/1552)
* Enhancement - Add sidebar toggle icons: [#5165](https://github.com/owncloud/web/issues/5165)
* Enhancement - Show compile errors and enforce node.js version in package.json: [#1579](https://github.com/owncloud/owncloud-design-system/pull/1579)
* Enhancement - Downgrade sass version: [#1583](https://github.com/owncloud/owncloud-design-system/pull/1583)

## Details

* Bugfix - Initial sorting for OcTableFiles: [#1588](https://github.com/owncloud/owncloud-design-system/pull/1588)

   The OcTableFiles component didn't apply initial sorting. This was set to the first sortable
   field in the list of fields (which is "name").

   https://github.com/owncloud/owncloud-design-system/pull/1588


* Bugfix - Sorting by date: [#1552](https://github.com/owncloud/owncloud-design-system/issues/1552)

   The OcTableFiles component sorted rows lexicographically by relative dates (e.g. "1 hour
   ago"). This was fixed to sort by unix timestamps instead.

   https://github.com/owncloud/owncloud-design-system/issues/1552
   https://github.com/owncloud/owncloud-design-system/pull/1588


* Enhancement - Add sidebar toggle icons: [#5165](https://github.com/owncloud/web/issues/5165)

   We added new chevron icons for the toggle button to switch between sidebar open and close

   https://github.com/owncloud/web/issues/5165
   https://github.com/owncloud/owncloud-design-system/pull/1587
   https://github.com/owncloud/owncloud-design-system/pull/1592


* Enhancement - Show compile errors and enforce node.js version in package.json: [#1579](https://github.com/owncloud/owncloud-design-system/pull/1579)

   Ensure we show compile-time errors on the command line when present and enforce Node.js
   v14.0.0 or greater to permit optional chaining to be used

   https://github.com/owncloud/owncloud-design-system/pull/1579


* Enhancement - Downgrade sass version: [#1583](https://github.com/owncloud/owncloud-design-system/pull/1583)

   Decrease the version of sass in order to prevent emitting of deprecation warnings in the ui-kit
   libray.

   https://github.com/owncloud/owncloud-design-system/pull/1583

# Changelog for [9.0.1] (2021-08-11)

The following sections list the changes in ownCloud Design System 9.0.1.

[9.0.1]: https://github.com/owncloud/owncloud-design-system/compare/v9.0.0...v9.0.1

## Summary

* Bugfix - Contextmenu should change selection model: [#1580](https://github.com/owncloud/owncloud-design-system/pull/1580)

## Details

* Bugfix - Contextmenu should change selection model: [#1580](https://github.com/owncloud/owncloud-design-system/pull/1580)

   We changed the rightclick logic so the selection model changes to the file the rightclick was
   performed on.

   https://github.com/owncloud/owncloud-design-system/pull/1580

# Changelog for [9.0.0] (2021-08-09)

The following sections list the changes in ownCloud Design System 9.0.0.

[9.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v8.3.1...v9.0.0

## Summary

* Change - Remove deprecated components: [#1545](https://github.com/owncloud/owncloud-design-system/pull/1545)
* Change - Enable multiple highlighting for OcTableFiles: [#5164](https://github.com/owncloud/web/issues/5164)
* Change - Production Bundle Content: [#1553](https://github.com/owncloud/owncloud-design-system/pull/1553)
* Enhancement - Add drop target: [#1575](https://github.com/owncloud/owncloud-design-system/pull/1575)
* Enhancement - Group files and folder when sorting by name: [#5613](https://github.com/owncloud/web/issues/5613)
* Enhancement - Add manual mode to oc-drop: [#1575](https://github.com/owncloud/owncloud-design-system/pull/1575)
* Enhancement - Add sign-out icon: [#5590](https://github.com/owncloud/web/issues/5590)
* Enhancement - Added drag-drop property: [#5592](https://github.com/owncloud/web/issues/5592)

## Details

* Change - Remove deprecated components: [#1545](https://github.com/owncloud/owncloud-design-system/pull/1545)

   To focus on the quality of the currently relevant components and to reduce bundle size, we are
   dropping a bunch of deprecated components from our design system.

   List of removed components, for reference:

   - OcActionDrop.vue - OcAppLayout.vue - OcAutocomplete.vue - OcDatepicker.vue -
   OcDisclosureDrop.vue - OcNav.vue - OcNavItem.vue - OcNavbar.vue - OcTabItem.vue -
   OcTabbed.vue - OcTabbedPanel.vue - OcTabbedTab.vue - OcTabs.vue - OcTopBar.vue -
   OcUser.vue - _OcNavbarItem.vue - _OcSidebarNavDivider.vue - _OcSidebarNavHeader.vue -
   _OcTopBarItem.vue - _OcTopBarLogo.vue

   https://github.com/owncloud/owncloud-design-system/pull/1545


* Change - Enable multiple highlighting for OcTableFiles: [#5164](https://github.com/owncloud/web/issues/5164)

   We changed the highlighting in a way that now every selected file is highlighted
   automatically. The `highlighted` prop has been removed as it's not used anymore.

   https://github.com/owncloud/web/issues/5164
   https://github.com/owncloud/owncloud-design-system/pull/1568


* Change - Production Bundle Content: [#1553](https://github.com/owncloud/owncloud-design-system/pull/1553)

   In the past, we shipped the docs `.scss` file and some example images which now have been removed
   from the production bundle to reduce size.

   https://github.com/owncloud/owncloud-design-system/pull/1553


* Enhancement - Add drop target: [#1575](https://github.com/owncloud/owncloud-design-system/pull/1575)

   We've added a target prop to the `oc-drop` component so that the drop can be opened at another
   element than the trigger.

   https://github.com/owncloud/owncloud-design-system/pull/1575


* Enhancement - Group files and folder when sorting by name: [#5613](https://github.com/owncloud/web/issues/5613)

   Like in oc10 when sorting by name, files and folders should be listed separately

   https://github.com/owncloud/web/issues/5613
   https://github.com/owncloud/owncloud-design-system/pull/1559


* Enhancement - Add manual mode to oc-drop: [#1575](https://github.com/owncloud/owncloud-design-system/pull/1575)

   We've added the manual mode to `oc-drop` component. This mode enables showing/hiding the drop
   **only** via provided methods ([see
   docs](https://owncloud.design/#/oC%20Components/OcDrop)).

   https://github.com/owncloud/owncloud-design-system/pull/1575


* Enhancement - Add sign-out icon: [#5590](https://github.com/owncloud/web/issues/5590)

   There has been confusion in user experience about the current usage of the `exit_to_app` icon
   as sign-out icon. We have added a dedicated `sign-out` icon to avoid this in the future.

   https://github.com/owncloud/web/issues/5590
   https://github.com/owncloud/owncloud-design-system/pull/1551


* Enhancement - Added drag-drop property: [#5592](https://github.com/owncloud/web/issues/5592)

   In order to enable file moving via drag & drop we added the property drag-drop on OcTable and
   OcTableFiles

   https://github.com/owncloud/web/issues/5592
   https://github.com/owncloud/owncloud-design-system/pull/1539
   https://github.com/owncloud/owncloud-design-system/pull/1562
   https://github.com/owncloud/owncloud-design-system/pull/1574

# Changelog for [8.3.1] (2021-08-04)

The following sections list the changes in ownCloud Design System 8.3.1.

[8.3.1]: https://github.com/owncloud/owncloud-design-system/compare/v8.3.0...v8.3.1

## Summary

* Bugfix - Unnecessary context menu events: [#1564](https://github.com/owncloud/owncloud-design-system/pull/1564)

## Details

* Bugfix - Unnecessary context menu events: [#1564](https://github.com/owncloud/owncloud-design-system/pull/1564)

   Clicking on the context menu in OcTableFiles was emitting unnecessary showDetails events.
   This has been fixed.

   https://github.com/owncloud/owncloud-design-system/pull/1564

# Changelog for [8.3.0] (2021-07-28)

The following sections list the changes in ownCloud Design System 8.3.0.

[8.3.0]: https://github.com/owncloud/owncloud-design-system/compare/v8.2.0...v8.3.0

## Summary

* Enhancement - Update vue select: [#1536](https://github.com/owncloud/owncloud-design-system/pull/1536)

## Details

* Enhancement - Update vue select: [#1536](https://github.com/owncloud/owncloud-design-system/pull/1536)

   We've updated vue select to version 3.12.0

   https://github.com/owncloud/owncloud-design-system/pull/1536

# Changelog for [8.2.0] (2021-07-23)

The following sections list the changes in ownCloud Design System 8.2.0.

[8.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v8.1.0...v8.2.0

## Summary

* Bugfix - Button as inline element: [#1529](https://github.com/owncloud/owncloud-design-system/pull/1529)
* Bugfix - OcTableFiles Contextmenu ID: [#1525](https://github.com/owncloud/owncloud-design-system/pull/1525)
* Enhancement - Deprecate OcAutocomplete: [#1521](https://github.com/owncloud/owncloud-design-system/pull/1521)
* Enhancement - Add OcRecipient component: [#1521](https://github.com/owncloud/owncloud-design-system/pull/1521)
* Enhancement - Add `search:input` event to OcSelect: [#1521](https://github.com/owncloud/owncloud-design-system/pull/1521)

## Details

* Bugfix - Button as inline element: [#1529](https://github.com/owncloud/owncloud-design-system/pull/1529)

   We made our OcButton an inline-flex instead of a flex element, so that it behaves correctly
   regarding alignment inside containers.

   https://github.com/owncloud/owncloud-design-system/pull/1529


* Bugfix - OcTableFiles Contextmenu ID: [#1525](https://github.com/owncloud/owncloud-design-system/pull/1525)

   Handling of `=`s in resource IDs in the OcTableFiles contextmenu was broken and lead to
   undesired behavior. This has been fixed now.

   https://github.com/owncloud/owncloud-design-system/pull/1525
   https://github.com/owncloud/owncloud-design-system/pull/1528


* Enhancement - Deprecate OcAutocomplete: [#1521](https://github.com/owncloud/owncloud-design-system/pull/1521)

   We've deprecated the OcAutocomplete component. OcSelect with `:multiple="true"` prop can
   be used to achieve the same behaviour.

   https://github.com/owncloud/owncloud-design-system/pull/1521


* Enhancement - Add OcRecipient component: [#1521](https://github.com/owncloud/owncloud-design-system/pull/1521)

   We've added OcRecipient component.

   https://github.com/owncloud/owncloud-design-system/pull/1521


* Enhancement - Add `search:input` event to OcSelect: [#1521](https://github.com/owncloud/owncloud-design-system/pull/1521)

   We've added `search:input` event to the OcSelect component which is triggered when a search
   query is typed into the select.

   https://github.com/owncloud/owncloud-design-system/pull/1521

# Changelog for [8.1.0] (2021-07-22)

The following sections list the changes in ownCloud Design System 8.1.0.

[8.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v8.0.0...v8.1.0

## Summary

* Bugfix - No hover on current breadcrumb item: [#1416](https://github.com/owncloud/owncloud-design-system/issues/1416)
* Bugfix - Inverse button hover: [#1401](https://github.com/owncloud/owncloud-design-system/issues/1401)
* Bugfix - OcSwitch Off State is not pixel perfect: [#1458](https://github.com/owncloud/owncloud-design-system/issues/1458)
* Enhancement - OcSwitch off state color: [#1457](https://github.com/owncloud/owncloud-design-system/issues/1457)
* Enhancement - OcTableFiles dropdown menu: [#1420](https://github.com/owncloud/owncloud-design-system/pull/1420)
* Enhancement - OcTextarea configurable Enter/Linebreak: [#1422](https://github.com/owncloud/owncloud-design-system/issues/1422)

## Details

* Bugfix - No hover on current breadcrumb item: [#1416](https://github.com/owncloud/owncloud-design-system/issues/1416)

   The current page's breadcrumb item is not interactive and shouldn't have a focus on hover, so we
   removed it.

   https://github.com/owncloud/owncloud-design-system/issues/1416
   https://github.com/owncloud/owncloud-design-system/pull/1511


* Bugfix - Inverse button hover: [#1401](https://github.com/owncloud/owncloud-design-system/issues/1401)

   Icons inside inverse buttons were hidden on hover, which shouldn't happen anymore. Also added
   an icon and some padding to the OcButton docs.

   https://github.com/owncloud/owncloud-design-system/issues/1401
   https://github.com/owncloud/owncloud-design-system/pull/1506


* Bugfix - OcSwitch Off State is not pixel perfect: [#1458](https://github.com/owncloud/owncloud-design-system/issues/1458)

   We've addressed the fact that the OcSwitch wasn't aligned properly in it's off state. Both
   states are now aligned equaly.

   https://github.com/owncloud/owncloud-design-system/issues/1458
   https://github.com/owncloud/owncloud-design-system/pull/1512


* Enhancement - OcSwitch off state color: [#1457](https://github.com/owncloud/owncloud-design-system/issues/1457)

   We've addressed the fact that the OcSwitch 'off' state color doesn't differ enough from 'on'
   state

   https://github.com/owncloud/owncloud-design-system/issues/1457
   https://github.com/owncloud/owncloud-design-system/pull/1522


* Enhancement - OcTableFiles dropdown menu: [#1420](https://github.com/owncloud/owncloud-design-system/pull/1420)

   The OcTableFiles component now has a slot that will be displayed inside an OcDrop upon either
   right-clicking a table row or the three-dot button on the right end of the row.

   https://github.com/owncloud/web/issues/5102
   https://github.com/owncloud/web/issues/5103
   https://github.com/owncloud/owncloud-design-system/pull/1420


* Enhancement - OcTextarea configurable Enter/Linebreak: [#1422](https://github.com/owncloud/owncloud-design-system/issues/1422)

   OcTextArea has now an property 'submitOnEnter'. This prop controls how the textarea should
   react to ENTER.

   https://github.com/owncloud/owncloud-design-system/issues/1422
   https://github.com/owncloud/owncloud-design-system/pull/1517

# Changelog for [8.0.0] (2021-07-08)

The following sections list the changes in ownCloud Design System 8.0.0.

[8.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v7.5.0...v8.0.0

## Summary

* Change - Rename previews to thumbnails: [#1429](https://github.com/owncloud/owncloud-design-system/pull/1429)
* Change - Remove custom model from OcSelect: [#1478](https://github.com/owncloud/owncloud-design-system/pull/1478)
* Enhancement - Add page size component: [#1478](https://github.com/owncloud/owncloud-design-system/pull/1478)
* Enhancement - Initial focus in OcModal: [#1453](https://github.com/owncloud/owncloud-design-system/pull/1453)
* Enhancement - Table Sorting by String or Function: [#1377](https://github.com/owncloud/owncloud-design-system/pull/1377)

## Details

* Change - Rename previews to thumbnails: [#1429](https://github.com/owncloud/owncloud-design-system/pull/1429)

   Until now, we have been referring to preview images in OcFilesTable and OcResource as
   `preview`. The technically more correct term for our use case would be `thumbnail` and since
   we're introducing proper previews in the web repo now this change renames `previews` in ODS to
   `thumbnails` to avoid naming conflicts down the road.

   https://github.com/owncloud/owncloud-design-system/pull/1429


* Change - Remove custom model from OcSelect: [#1478](https://github.com/owncloud/owncloud-design-system/pull/1478)

   We've removed the custom model defined in OcSelect component. It is no longer possible to pass
   the value via `model` prop. You can use `value` instead. `v-model` is still available.

   https://github.com/owncloud/owncloud-design-system/pull/1478


* Enhancement - Add page size component: [#1478](https://github.com/owncloud/owncloud-design-system/pull/1478)

   We've added page size component. This component can be used for specifying the pagination
   limit of items per page.

   https://github.com/owncloud/owncloud-design-system/pull/1478


* Enhancement - Initial focus in OcModal: [#1453](https://github.com/owncloud/owncloud-design-system/pull/1453)

   The OcModal component sets the initial focus to the input element now, if one exists. If the
   input element doesn't exist the focus remains on the modal itself, like we had it before.

   https://github.com/owncloud/web/issues/3684
   https://github.com/owncloud/owncloud-design-system/pull/1453


* Enhancement - Table Sorting by String or Function: [#1377](https://github.com/owncloud/owncloud-design-system/pull/1377)

   Sorting inside the OcTable now can be done by passing a string or a function, so you can e.g. sort
   objects inside Array by their properties.

   https://github.com/owncloud/owncloud-design-system/pull/1377

# Changelog for [7.5.0] (2021-06-25)

The following sections list the changes in ownCloud Design System 7.5.0.

[7.5.0]: https://github.com/owncloud/owncloud-design-system/compare/v7.4.2...v7.5.0

## Summary

* Bugfix - Show `--` as file-size when the folder size is not yet computed: [#1402](https://github.com/owncloud/owncloud-design-system/issues/1402)
* Enhancement - Improve `OcSwitch` UI and accessibility: [#1421](https://github.com/owncloud/owncloud-design-system/pull/1421)

## Details

* Bugfix - Show `--` as file-size when the folder size is not yet computed: [#1402](https://github.com/owncloud/owncloud-design-system/issues/1402)

   When the folder size is unknown (not yet calculated or cannot be calculated) the server returns
   `-1` as the size. In this case we now show `--` and not only an empty string

   https://github.com/owncloud/owncloud-design-system/issues/1402
   https://github.com/owncloud/owncloud-design-system/pull/1393


* Enhancement - Improve `OcSwitch` UI and accessibility: [#1421](https://github.com/owncloud/owncloud-design-system/pull/1421)

   We've added round corners to the `OcSwitch` component so that it is more consistent with other
   elements in our UI. We've also made the component accessible, added visible label and enabled
   toggling its state with the spacebar.

   https://github.com/owncloud/owncloud-design-system/pull/1421

# Changelog for [7.4.2] (2021-06-21)

The following sections list the changes in ownCloud Design System 7.4.2.

[7.4.2]: https://github.com/owncloud/owncloud-design-system/compare/v7.4.0...v7.4.2

## Summary

* Bugfix - 0.5px separator line between OcSidebarNav items: [#1402](https://github.com/owncloud/owncloud-design-system/pull/1402)
* Bugfix - OcIcon crashes if icon gets updated: [#1407](https://github.com/owncloud/owncloud-design-system/pull/1407)
* Bugfix - Pagination renders unnecessary ... skip label: [#1406](https://github.com/owncloud/owncloud-design-system/pull/1406)

## Details

* Bugfix - 0.5px separator line between OcSidebarNav items: [#1402](https://github.com/owncloud/owncloud-design-system/pull/1402)

   The small line between OcSidebarNav items was displayed with different widths depending on
   the used browser, since some can't handle half-pixel values. We've changed that to create the
   same look across browsers.

   https://github.com/owncloud/owncloud-design-system/pull/1402


* Bugfix - OcIcon crashes if icon gets updated: [#1407](https://github.com/owncloud/owncloud-design-system/pull/1407)

   Before this bugfix, updating `oc-icon` name prop crashed `vue-inline-svg`. We had to
   overwrite `vue-inline-svg` download method and forgot to implement its `isPending`
   property on the returned promise.

   This is fixed now and tested

   https://github.com/owncloud/owncloud-design-system/pull/1407


* Bugfix - Pagination renders unnecessary ... skip label: [#1406](https://github.com/owncloud/owncloud-design-system/pull/1406)

   In cases where the pagination should only render 4 pages at once and 4 are available, it rendered
   a skip label `< 1 2 ... 4 >` even if it's not required.

   Now this is fixed and it renders `< 1 2 3 4 >` instead.

   https://github.com/owncloud/owncloud-design-system/pull/1406

# Changelog for [7.4.0] (2021-06-17)

The following sections list the changes in ownCloud Design System 7.4.0.

[7.4.0]: https://github.com/owncloud/owncloud-design-system/compare/v7.4.1...v7.4.0

## Summary

* Bugfix - Sticky OcAppSideBar Header: [#1384](https://github.com/owncloud/owncloud-design-system/pull/1384)
* Enhancement - Theme-able breakpoint variables: [#5281](https://github.com/owncloud/web/pull/5281)

## Details

* Bugfix - Sticky OcAppSideBar Header: [#1384](https://github.com/owncloud/owncloud-design-system/pull/1384)

   The OcAppSideBar Header is now sticky so the user always sees the filename and most important
   info. UX is also improved since he now can always close the appsidebar on mobile without having
   to scroll back to the top.

   https://github.com/owncloud/owncloud-design-system/pull/1384


* Enhancement - Theme-able breakpoint variables: [#5281](https://github.com/owncloud/web/pull/5281)

   We've added custom CSS props for breakpoints a while ago but missed adding them to the theming
   initialization, so they weren't theme-able until now.

   https://github.com/owncloud/web/pull/5281
   https://github.com/owncloud/owncloud-design-system/pull/1388

# Changelog for [7.4.1] (2021-06-17)

The following sections list the changes in ownCloud Design System 7.4.1.

[7.4.1]: https://github.com/owncloud/owncloud-design-system/compare/v7.3.0...v7.4.1

## Summary

* Bugfix - Remove pagination list padding: [#1398](https://github.com/owncloud/owncloud-design-system/pull/1398)
* Bugfix - Visible separator between OcSidebarNav items: [#1387](https://github.com/owncloud/owncloud-design-system/issues/1387)

## Details

* Bugfix - Remove pagination list padding: [#1398](https://github.com/owncloud/owncloud-design-system/pull/1398)

   The pagination list had a small left padding which caused it to be visually off from the desired
   horizontal center.

   https://github.com/owncloud/owncloud-design-system/pull/1398


* Bugfix - Visible separator between OcSidebarNav items: [#1387](https://github.com/owncloud/owncloud-design-system/issues/1387)

   We have added a small line between OcSidebarNav items in the `active` and `hover` state to
   visually make them better differentiable.

   https://github.com/owncloud/owncloud-design-system/issues/1387
   https://github.com/owncloud/owncloud-design-system/pull/1390

# Changelog for [7.3.0] (2021-06-14)

The following sections list the changes in ownCloud Design System 7.3.0.

[7.3.0]: https://github.com/owncloud/owncloud-design-system/compare/v7.2.0...v7.3.0

## Summary

* Enhancement - Logo sizing: [#795](https://github.com/owncloud/owncloud-design-system/issues/795)

## Details

* Enhancement - Logo sizing: [#795](https://github.com/owncloud/owncloud-design-system/issues/795)

   We've added theming variables for the max-width and max-height of the OcLogo component and
   increased the default (max-)height.

   https://github.com/owncloud/owncloud-design-system/issues/795
   https://github.com/owncloud/web/issues/5128
   https://github.com/owncloud/owncloud-design-system/pull/1380

# Changelog for [7.2.0] (2021-06-11)

The following sections list the changes in ownCloud Design System 7.2.0.

[7.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v7.1.0...v7.2.0

## Summary

* Bugfix - Use correct selected background colour: [#1368](https://github.com/owncloud/owncloud-design-system/pull/1368)
* Enhancement - Add pagination component: [#1055](https://github.com/owncloud/owncloud-design-system/pull/1055)
* Enhancement - Table Row Mounted Event: [#1371](https://github.com/owncloud/owncloud-design-system/pull/1371)

## Details

* Bugfix - Use correct selected background colour: [#1368](https://github.com/owncloud/owncloud-design-system/pull/1368)

   We've fixed the css custom property in `.oc-background-selected` helper class so that the
   background highlighted colour is used there instead.

   https://github.com/owncloud/owncloud-design-system/pull/1368


* Enhancement - Add pagination component: [#1055](https://github.com/owncloud/owncloud-design-system/pull/1055)

   We've added `OcPagination` component.

   https://github.com/owncloud/owncloud-design-system/pull/1055


* Enhancement - Table Row Mounted Event: [#1371](https://github.com/owncloud/owncloud-design-system/pull/1371)

   The OcTable now emits an event if a row is mounted.

   https://github.com/owncloud/owncloud-design-system/pull/1371

# Changelog for [7.1.0] (2021-06-02)

The following sections list the changes in ownCloud Design System 7.1.0.

[7.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v7.0.1...v7.1.0

## Summary

* Enhancement - Accessible label for all landmarks in OcSidebarNav: [#1345](https://github.com/owncloud/owncloud-design-system/pull/1345)
* Enhancement - Add check icon without circle: [#1341](https://github.com/owncloud/owncloud-design-system/issues/1341)

## Details

* Enhancement - Accessible label for all landmarks in OcSidebarNav: [#1345](https://github.com/owncloud/owncloud-design-system/pull/1345)

   The OcSidebarNav component now has properties for providing accessible labels for all
   landmarks of the sidebar.

   https://github.com/owncloud/owncloud-design-system/pull/1345


* Enhancement - Add check icon without circle: [#1341](https://github.com/owncloud/owncloud-design-system/issues/1341)

   According to community requests, we've added a check/tick/approve icon without a background
   shape.

   https://github.com/owncloud/owncloud-design-system/issues/1341
   https://github.com/owncloud/owncloud-design-system/pull/1346

# Changelog for [7.0.1] (2021-06-01)

The following sections list the changes in ownCloud Design System 7.0.1.

[7.0.1]: https://github.com/owncloud/owncloud-design-system/compare/v7.0.0...v7.0.1

## Summary

* Bugfix - OcStatusIndicator ID fixes: [#1342](https://github.com/owncloud/owncloud-design-system/pull/1342)

## Details

* Bugfix - OcStatusIndicator ID fixes: [#1342](https://github.com/owncloud/owncloud-design-system/pull/1342)

   Inside the OcStatusIndicator, components were rendering the same IDs which resulted in
   invalid HTML. Since the clickhandler on the OcStatusIndicator also was refering to these IDs
   and they now are dynamically generated on the frontend, it has now been replaced with a computed
   property which receives the click target.

   Also, the accessible description (only applicable to screen readers) has been changed from
   `<span>` to `<p>` which sounds better to the listener.

   https://github.com/owncloud/owncloud-design-system/pull/1342

# Changelog for [7.0.0] (2021-05-31)

The following sections list the changes in ownCloud Design System 7.0.0.

[7.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v6.4.0...v7.0.0

## Summary

* Bugfix - Interactive texts in OcTooltips: [#1288](https://github.com/owncloud/owncloud-design-system/pull/1288)
* Bugfix - Make OcSidebarNavitem text bold: [#1308](https://github.com/owncloud/owncloud-design-system/pull/1308)
* Bugfix - Move arialabel in OcSideBarNav to nav element: [#1325](https://github.com/owncloud/owncloud-design-system/pull/1325)
* Bugfix - Do not define table header height via row height var: [#1309](https://github.com/owncloud/owncloud-design-system/pull/1309)
* Bugfix - Non-interactive tooltips: [#1330](https://github.com/owncloud/owncloud-design-system/pull/1330)
* Change - Rename OcSidebar to OcSidebarNav: [#1300](https://github.com/owncloud/owncloud-design-system/pull/1300)
* Change - Use slots in OcSidebarNav: [#1300](https://github.com/owncloud/owncloud-design-system/pull/1300)
* Enhancement - A11y color contrasts: [#1302](https://github.com/owncloud/owncloud-design-system/pull/1302)
* Enhancement - Color contrast for OcSidebarNav: [#1310](https://github.com/owncloud/owncloud-design-system/pull/1310)
* Enhancement - Avatar link & count a11y improvements: [#1298](https://github.com/owncloud/owncloud-design-system/pull/1298)
* Enhancement - Add inverse button variant: [#1300](https://github.com/owncloud/owncloud-design-system/pull/1300)
* Enhancement - Modals listen to escape key: [#1324](https://github.com/owncloud/owncloud-design-system/pull/1324)
* Enhancement - Apply CSS custom properties to `:host`: [#1333](https://github.com/owncloud/owncloud-design-system/pull/1333)
* Enhancement - Aria-labelledby through heading for modals: [#1327](https://github.com/owncloud/owncloud-design-system/pull/1327)
* Enhancement - Reduce the height of the filelist rows: [#1319](https://github.com/owncloud/owncloud-design-system/pull/1319)
* Enhancement - Add tabindex on active breadcrumb element: [#1334](https://github.com/owncloud/owncloud-design-system/pull/1334)
* Enhancement - Clear button for OcTextInput: [#1329](https://github.com/owncloud/owncloud-design-system/pull/1329)
* Enhancement - Themable Table Row Height: [#1291](https://github.com/owncloud/owncloud-design-system/pull/1291)
* Enhancement - Use button to trigger table sort: [#1309](https://github.com/owncloud/owncloud-design-system/pull/1309)

## Details

* Bugfix - Interactive texts in OcTooltips: [#1288](https://github.com/owncloud/owncloud-design-system/pull/1288)

   The new tooltip texts would only be initialized once and then didn't offer the possibility of
   chaging reactively. By updating the directive, we've changed that.

   https://github.com/owncloud/owncloud-design-system/pull/1288


* Bugfix - Make OcSidebarNavitem text bold: [#1308](https://github.com/owncloud/owncloud-design-system/pull/1308)

   We've made the text of `OcSidebarNavItem` component bold again.

   https://github.com/owncloud/owncloud-design-system/pull/1308


* Bugfix - Move arialabel in OcSideBarNav to nav element: [#1325](https://github.com/owncloud/owncloud-design-system/pull/1325)

   The accessibleLabel property describing the navigation inside the OcSideBarNav was set to
   the whole component, yet was meant to describe the `<nav>` element so we changed its target.

   https://github.com/owncloud/owncloud-design-system/pull/1325


* Bugfix - Do not define table header height via row height var: [#1309](https://github.com/owncloud/owncloud-design-system/pull/1309)

   We've separated the height styles of table header and table row so that the row height variable
   is not affecting both of them.

   https://github.com/owncloud/owncloud-design-system/pull/1309


* Bugfix - Non-interactive tooltips: [#1330](https://github.com/owncloud/owncloud-design-system/pull/1330)

   Interactive tooltips were breaking correct z-index behaviour so we set them to
   non-interactive.

   https://github.com/owncloud/owncloud-design-system/pull/1330


* Change - Rename OcSidebar to OcSidebarNav: [#1300](https://github.com/owncloud/owncloud-design-system/pull/1300)

   We've changed the name of the `OcSidebar` component to `OcSidebarNav` to better reflect that
   it shouldn't be used for anything else than the navigation.

   https://github.com/owncloud/owncloud-design-system/pull/1300


* Change - Use slots in OcSidebarNav: [#1300](https://github.com/owncloud/owncloud-design-system/pull/1300)

   We've moved away from defining the nav items in the `OcSidebarNav` via props to including them
   through slots. The component now contains 3 slots - `header`, `nav` and `footer`.

   https://github.com/owncloud/owncloud-design-system/pull/1300


* Enhancement - A11y color contrasts: [#1302](https://github.com/owncloud/owncloud-design-system/pull/1302)

   The color contrast checker now runs when generating tokens via the yarn token command and
   reports if any colors don't match the minimum color contrast

   https://github.com/owncloud/owncloud-design-system/pull/1302
   https://github.com/owncloud/owncloud-design-system/pull/1316


* Enhancement - Color contrast for OcSidebarNav: [#1310](https://github.com/owncloud/owncloud-design-system/pull/1310)

   We've set different colors for the active and hover state of nav items in the OcSidebarNav
   component. Those colors now fulfill required color contrast ratios.

   https://github.com/owncloud/owncloud-design-system/pull/1310
   https://github.com/owncloud/owncloud-design-system/pull/1315


* Enhancement - Avatar link & count a11y improvements: [#1298](https://github.com/owncloud/owncloud-design-system/pull/1298)

   - wrappend OcAvatarLink in span (was div) - wrappend OcAvatarCount in span (was div) - changed
   the way the OcAvatarLink background color was picked, which improves the color contrast to
   >4.5 (a11y requirement)

   https://github.com/owncloud/owncloud-design-system/pull/1298


* Enhancement - Add inverse button variant: [#1300](https://github.com/owncloud/owncloud-design-system/pull/1300)

   We've added inverse button variant which should be used on a dark background.

   https://github.com/owncloud/owncloud-design-system/pull/1300


* Enhancement - Modals listen to escape key: [#1324](https://github.com/owncloud/owncloud-design-system/pull/1324)

   The modal component now also emits its "cancel" event when the user presses the escape key. This
   can for example be used to close the modal via keyboard.

   https://github.com/owncloud/owncloud-design-system/pull/1324


* Enhancement - Apply CSS custom properties to `:host`: [#1333](https://github.com/owncloud/owncloud-design-system/pull/1333)

   We've applied CSS custom properties to `:host` so that we can access them also in web
   components.

   https://github.com/owncloud/owncloud-design-system/pull/1333


* Enhancement - Aria-labelledby through heading for modals: [#1327](https://github.com/owncloud/owncloud-design-system/pull/1327)

   The modal component now has an `aria-labelledby` attribute linked to the `<h2>` tag inside,
   for improved accessibility.

   https://github.com/owncloud/owncloud-design-system/pull/1327


* Enhancement - Reduce the height of the filelist rows: [#1319](https://github.com/owncloud/owncloud-design-system/pull/1319)

   We've reduced the height of the filelist rows to the possible minimum to provide a more
   condensed view especially for desktop users. The possible minimum of the filelist row is
   currently determined by the height of the filename and the sharing-indicators below the
   filename.

   https://github.com/owncloud/owncloud-design-system/pull/1319


* Enhancement - Add tabindex on active breadcrumb element: [#1334](https://github.com/owncloud/owncloud-design-system/pull/1334)

   We've added a tabindex to the span that can represent the latest breadcrumb to make it possible
   to set the focus there on page navigation for accessibility reasons.

   https://github.com/owncloud/owncloud-design-system/pull/1334


* Enhancement - Clear button for OcTextInput: [#1329](https://github.com/owncloud/owncloud-design-system/pull/1329)

   The OcTextInput component now has an optional button for clearing the input.

   https://github.com/owncloud/owncloud-design-system/pull/1329


* Enhancement - Themable Table Row Height: [#1291](https://github.com/owncloud/owncloud-design-system/pull/1291)

   By moving the table row height from a component property to a themable variable, we give users of
   the web frontend a way to customize the appearance of their UI (instead of only giving the
   freedom to users of the ODS).

   https://github.com/owncloud/owncloud-design-system/pull/1291


* Enhancement - Use button to trigger table sort: [#1309](https://github.com/owncloud/owncloud-design-system/pull/1309)

   We've added a button to the table head cell which can trigger sort for its column so that it is
   keyboard accessible and is aligned directly next to the cell text.

   https://github.com/owncloud/owncloud-design-system/pull/1309

# Changelog for [6.4.0] (2021-05-06)

The following sections list the changes in ownCloud Design System 6.4.0.

[6.4.0]: https://github.com/owncloud/owncloud-design-system/compare/v6.3.0...v6.4.0

## Summary

* Bugfix - OcSpinner make ariaLabel prop optional: [#1251](https://github.com/owncloud/owncloud-design-system/pull/1251)
* Enhancement - Lazy img loading & accessibility for OcAvatar: [#1282](https://github.com/owncloud/owncloud-design-system/pull/1282)
* Enhancement - Add description for avatar groups: [#1250](https://github.com/owncloud/owncloud-design-system/pull/1250)
* Enhancement - Make files table headings more descriptive: [#1252](https://github.com/owncloud/owncloud-design-system/pull/1252)
* Enhancement - Remove uk-drop from OcDrop: [#1269](https://github.com/owncloud/owncloud-design-system/pull/1269)
* Enhancement - Accessibility for OcSelect: [#1268](https://github.com/owncloud/owncloud-design-system/pull/1268)
* Enhancement - Remove the notification close button: [#1247](https://github.com/owncloud/owncloud-design-system/pull/1247)
* Enhancement - Show the notification close button via property: [#1247](https://github.com/owncloud/owncloud-design-system/pull/1247)
* Enhancement - Tooltips: [#1263](https://github.com/owncloud/owncloud-design-system/pull/1263)

## Details

* Bugfix - OcSpinner make ariaLabel prop optional: [#1251](https://github.com/owncloud/owncloud-design-system/pull/1251)

   The OcSpinner component had a required property that caused console errors. Since we're not
   always using it in a case where a label is required, we've made it optional but strongly
   recommend to use it unless the element is wrapped in/equiped with a `aria-hidden="true"`
   attribute.

   https://github.com/owncloud/owncloud-design-system/pull/1251


* Enhancement - Lazy img loading & accessibility for OcAvatar: [#1282](https://github.com/owncloud/owncloud-design-system/pull/1282)

   - Add lazy loading to OcImg component - Internalize former dependency vue-avatar into
   OcAvatar - Make OcAvatar use OcImg component, using lazy loading - Change OcAvatar to be a11y
   compliant (color contrasts, DOM structure)

   https://github.com/owncloud/owncloud-design-system/pull/1282


* Enhancement - Add description for avatar groups: [#1250](https://github.com/owncloud/owncloud-design-system/pull/1250)

   The description is mandatory for avatar groups because the avatar group element is hidden for
   screen readers.

   https://github.com/owncloud/owncloud-design-system/pull/1250


* Enhancement - Make files table headings more descriptive: [#1252](https://github.com/owncloud/owncloud-design-system/pull/1252)

   This makes it easier to understand for people using a screen reader.

   https://github.com/owncloud/owncloud-design-system/pull/1252


* Enhancement - Remove uk-drop from OcDrop: [#1269](https://github.com/owncloud/owncloud-design-system/pull/1269)

   We've used uikit to manage oc-drop before, now tippy.js acts as a drop in replacement. The api
   stays the same.

   https://github.com/owncloud/owncloud-design-system/pull/1269


* Enhancement - Accessibility for OcSelect: [#1268](https://github.com/owncloud/owncloud-design-system/pull/1268)

   - Add l10n for default no options text - Add l10n for vue-select combobox aria-label

   https://github.com/owncloud/owncloud-design-system/pull/1268


* Enhancement - Remove the notification close button: [#1247](https://github.com/owncloud/owncloud-design-system/pull/1247)

   The notification will now automatically close after a certain amount of time which can be
   defined via property. Also fixed some links in the docs in addition to that.

   https://github.com/owncloud/owncloud-design-system/pull/1247


* Enhancement - Show the notification close button via property: [#1247](https://github.com/owncloud/owncloud-design-system/pull/1247)

   Close button can now be enabled via property dismissible, additionally the notification will
   now automatically close after a certain amount of time. The duration a notification shows can
   be defined via a newly added property. Also fixed some links in the docs in addition to that.

   https://github.com/owncloud/owncloud-design-system/pull/1247


* Enhancement - Tooltips: [#1263](https://github.com/owncloud/owncloud-design-system/pull/1263)

   We've used uikit to display tooltips before, now we added a own `v-oc-tooltip` directive and
   use this instead. This improves accessibility and possibilities for later re-use.

   https://github.com/owncloud/owncloud-design-system/pull/1263
   https://github.com/owncloud/owncloud-design-system/pull/1271
   https://github.com/owncloud/owncloud-design-system/pull/1272

# Changelog for [6.3.0] (2021-04-29)

The following sections list the changes in ownCloud Design System 6.3.0.

[6.3.0]: https://github.com/owncloud/owncloud-design-system/compare/v6.2.0...v6.3.0

## Summary

* Bugfix - Remove login paragraph meta styling: [#1246](https://github.com/owncloud/owncloud-design-system/pull/1246)
* Bugfix - Set correct href to oc-button if type is router-link: [#1245](https://github.com/owncloud/owncloud-design-system/pull/1245)
* Bugfix - Translateable default close button label in OcSidebar: [#1243](https://github.com/owncloud/owncloud-design-system/pull/1243)
* Bugfix - Remove unnecessary role attribute from oc-icon: [#1241](https://github.com/owncloud/owncloud-design-system/pull/1241)
* Enhancement - Improved accessibility for oc-accordion: [#1241](https://github.com/owncloud/owncloud-design-system/pull/1241)
* Enhancement - Aria-hide images if needed: [#1244](https://github.com/owncloud/owncloud-design-system/pull/1244)
* Enhancement - Use bigger font size for breadcrumbs: [#1239](https://github.com/owncloud/owncloud-design-system/pull/1239)
* Enhancement - Modal focus trap: [#1237](https://github.com/owncloud/owncloud-design-system/pull/1237)
* Enhancement - Add prop to define table padding: [#1240](https://github.com/owncloud/owncloud-design-system/pull/1240)

## Details

* Bugfix - Remove login paragraph meta styling: [#1246](https://github.com/owncloud/owncloud-design-system/pull/1246)

   The login screen styling paragraph tag was expanding uikit's text-meta which doesn't have
   enough contrast for accessibility. By removing it we fall back to the themed text color on the
   login/redirect pages.

   https://github.com/owncloud/owncloud-design-system/pull/1246


* Bugfix - Set correct href to oc-button if type is router-link: [#1245](https://github.com/owncloud/owncloud-design-system/pull/1245)

   While setting the type property to 'router-link' the empty href property has overwritten the
   href target with null. Now the href property will be set correctly.

   https://github.com/owncloud/owncloud-design-system/pull/1245


* Bugfix - Translateable default close button label in OcSidebar: [#1243](https://github.com/owncloud/owncloud-design-system/pull/1243)

   The sidebar component's close button label was hardcoded to be in English. Since the label is
   pretty unambiguous we can move its translations to the ODS.

   https://github.com/owncloud/owncloud-design-system/pull/1243


* Bugfix - Remove unnecessary role attribute from oc-icon: [#1241](https://github.com/owncloud/owncloud-design-system/pull/1241)

   The oc-icon component had `role="presentation"` on the svg if there is no accessible label.
   Since we already have `aria-hidden="true"` the role is not needed. Removing it for less code
   complexity.

   https://github.com/owncloud/owncloud-design-system/pull/1241


* Enhancement - Improved accessibility for oc-accordion: [#1241](https://github.com/owncloud/owncloud-design-system/pull/1241)

   We improved the accessibility of the oc-accordion component: - Switch from ul+li to just using
   divs. A list of accordions is not useful. - Make sure that collapsed accordions have
   aria-expanded="false" (attribute vanished from the html before, needs to be a string). -
   Remove aria-label from the button as the button already contains all the accessibility hints
   it needs.

   https://github.com/owncloud/owncloud-design-system/pull/1241


* Enhancement - Aria-hide images if needed: [#1244](https://github.com/owncloud/owncloud-design-system/pull/1244)

   When the `alt` property of the oc-image is empty we now set `aria-hidden="true"`.

   https://github.com/owncloud/owncloud-design-system/pull/1244


* Enhancement - Use bigger font size for breadcrumbs: [#1239](https://github.com/owncloud/owncloud-design-system/pull/1239)

   We've increased the font size of breadcrumbs to match the one from `oc-table` items.

   https://github.com/owncloud/owncloud-design-system/pull/1239


* Enhancement - Modal focus trap: [#1237](https://github.com/owncloud/owncloud-design-system/pull/1237)

   We've added [Vue focus trap library](https://github.com/posva/focus-trap-vue) to trap
   the keyboard navigation inside of the modal. To enable the focus trap, use prop
   `focusTrapActive`.

   https://github.com/owncloud/owncloud-design-system/pull/1237


* Enhancement - Add prop to define table padding: [#1240](https://github.com/owncloud/owncloud-design-system/pull/1240)

   We've added new prop called `paddingX` which can be used to set the padding along x axis of the
   `oc-table` or `oc-table-files`.

   https://github.com/owncloud/owncloud-design-system/pull/1240

# Changelog for [6.2.0] (2021-04-22)

The following sections list the changes in ownCloud Design System 6.2.0.

[6.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v6.1.0...v6.2.0

## Summary

* Bugfix - Remove required prop from sidebar & obsolete role: [#1233](https://github.com/owncloud/owncloud-design-system/pull/1233)
* Enhancement - Add tabindex to table rows: [#1233](https://github.com/owncloud/owncloud-design-system/pull/1233)

## Details

* Bugfix - Remove required prop from sidebar & obsolete role: [#1233](https://github.com/owncloud/owncloud-design-system/pull/1233)

   - property productName shouldn't be required anymore in OcSidebar after logoAlt has been
   introduced - role="progressbar" is obsolete on `<progress>` HTML elements

   https://github.com/owncloud/owncloud-design-system/pull/1233


* Enhancement - Add tabindex to table rows: [#1233](https://github.com/owncloud/owncloud-design-system/pull/1233)

   By adding a negative tabindex the table rows are now focusable which is an important aspect of
   accessibility/keyboard navigation.

   https://github.com/owncloud/owncloud-design-system/pull/1233

# Changelog for [6.1.0] (2021-04-22)

The following sections list the changes in ownCloud Design System 6.1.0.

[6.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v6.0.1...v6.1.0

## Summary

* Enhancement - Files table accessibility: [#1229](https://github.com/owncloud/owncloud-design-system/pull/1229)
* Enhancement - Improve modal component: [#1214](https://github.com/owncloud/owncloud-design-system/pull/1214)
* Enhancement - Improve accessibility of oc-breadcrumb component: [#1228](https://github.com/owncloud/owncloud-design-system/pull/1228)
* Enhancement - OCDrop accessibility: [#1230](https://github.com/owncloud/owncloud-design-system/pull/1230)
* Enhancement - Accessibility improvements on the sidebar component: [#1231](https://github.com/owncloud/owncloud-design-system/pull/1231)

## Details

* Enhancement - Files table accessibility: [#1229](https://github.com/owncloud/owncloud-design-system/pull/1229)

   * Add accessible description in case a resource link opens in a new window * Add accessible
   description for status indicators * Adjust some of the labels to make them more descriptive *
   Fix broken outline on directories

   https://github.com/owncloud/owncloud-design-system/pull/1229


* Enhancement - Improve modal component: [#1214](https://github.com/owncloud/owncloud-design-system/pull/1214)

   We've made the OcModal component more accessible: - It now features `role="dialog"` and
   `aria-modal="true"` - The modal title is now a `<h2>` - Component styles have been moved from an
   individual stylesheet to the component file

   https://github.com/owncloud/owncloud-design-system/pull/1214


* Enhancement - Improve accessibility of oc-breadcrumb component: [#1228](https://github.com/owncloud/owncloud-design-system/pull/1228)

   In order to enhance accessibility oc breadcrumb has been wrapped into a <nav> tag. The <ul>
   element changed to <ol> element, furthermore aria-current=page applies to the last item and
   the default home breadcrumb has been removed. Keyboard navigation has been fixed.

   https://github.com/owncloud/owncloud-design-system/pull/1228


* Enhancement - OCDrop accessibility: [#1230](https://github.com/owncloud/owncloud-design-system/pull/1230)

   Changed DOM order of example data to allow proper tab interactions, added focus border to a
   elements in list.

   https://github.com/owncloud/owncloud-design-system/pull/1230


* Enhancement - Accessibility improvements on the sidebar component: [#1231](https://github.com/owncloud/owncloud-design-system/pull/1231)

   Added a current page indicator and an overall aria label to the sidebar. Also implemented an alt
   property for the logo. It used the product name previously which is problematic because alt
   needs to describe the link rather than the image.

   https://github.com/owncloud/owncloud-design-system/pull/1231

# Changelog for [6.0.1] (2021-04-19)

The following sections list the changes in ownCloud Design System 6.0.1.

[6.0.1]: https://github.com/owncloud/owncloud-design-system/compare/v6.0.0...v6.0.1

## Summary

* Bugfix - Swap background colors: [#1227](https://github.com/owncloud/owncloud-design-system/pull/1227)

## Details

* Bugfix - Swap background colors: [#1227](https://github.com/owncloud/owncloud-design-system/pull/1227)

   In the `v6.0.0` release, the color values for `background-muted` and
   `background-hightlighted` got swapped by accident. This produced unwanted results in the
   results and gets reverted to the original and working version with this change.

   https://github.com/owncloud/owncloud-design-system/pull/1227

# Changelog for [6.0.0] (2021-04-19)

The following sections list the changes in ownCloud Design System 6.0.0.

[6.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v5.1.0...v6.0.0

## Summary

* Change - Use style-dictionary instead of theo: [#1217](https://github.com/owncloud/owncloud-design-system/pull/1217)
* Enhancement - More themable variables: [#1224](https://github.com/owncloud/owncloud-design-system/pull/1224)
* Enhancement - Table sorting: [#1219](https://github.com/owncloud/owncloud-design-system/pull/1219)

## Details

* Change - Use style-dictionary instead of theo: [#1217](https://github.com/owncloud/owncloud-design-system/pull/1217)

   We've used theo to generate tokens, it was a bit cumbersome to manage global aliases and hard to
   define property based filters.

   From now on we use amazon style-dictionary to generate the tokens which is more flexible. In the
   same process we renamed the tokens to get a more meaningful naming convention for css // scss
   variables.

   https://github.com/owncloud/owncloud-design-system/pull/1217


* Enhancement - More themable variables: [#1224](https://github.com/owncloud/owncloud-design-system/pull/1224)

   Due to recent changes, more variables other than colors have been moved from SASS to CSS
   variables and are therefore customizable through theming. This change loads them correctly
   when providing a `theme.json` as a Vue plugin option.

   https://github.com/owncloud/owncloud-design-system/pull/1224


* Enhancement - Table sorting: [#1219](https://github.com/owncloud/owncloud-design-system/pull/1219)

   The latest version of OcTable and OcTableFiles had no more sorting. We brought it back. It can be
   enabled by simply adding a `sortable: true` to table fields.

   https://github.com/owncloud/owncloud-design-system/pull/1219

# Changelog for [5.1.0] (2021-04-15)

The following sections list the changes in ownCloud Design System 5.1.0.

[5.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v5.0.1...v5.1.0

## Summary

* Bugfix - Remove fixed login logo width: [#1207](https://github.com/owncloud/owncloud-design-system/pull/1207)
* Bugfix - Muted color class: [#1210](https://github.com/owncloud/owncloud-design-system/pull/1210)
* Enhancement - OcCheckbox remove aria checked: [#1212](https://github.com/owncloud/owncloud-design-system/pull/1212)
* Enhancement - Text-inital helper class: [#1211](https://github.com/owncloud/owncloud-design-system/pull/1211)
* Enhancement - Update tag border color: [#1213](https://github.com/owncloud/owncloud-design-system/pull/1213)

## Details

* Bugfix - Remove fixed login logo width: [#1207](https://github.com/owncloud/owncloud-design-system/pull/1207)

   The logo displayed on login/error pages had a fixed width and height, wich was fine for the
   provided ownCloud logo but produced suboptimal results when using files with other
   width/height ratios. By changing this to max-width & max-height we allow for other files while
   keeping their proper aspect ratio.

   https://github.com/owncloud/owncloud-design-system/pull/1207


* Bugfix - Muted color class: [#1210](https://github.com/owncloud/owncloud-design-system/pull/1210)

   The `uk-text-muted` class caused some inconsistancy in ownCloud-web, sometimes being
   overwritten by default uikit styles. Changing it to `oc-text-muted` and updating the
   references in ownCloud-web should prevent this.

   https://github.com/owncloud/owncloud-design-system/pull/1210


* Enhancement - OcCheckbox remove aria checked: [#1212](https://github.com/owncloud/owncloud-design-system/pull/1212)

   We've introduced `aria-checked` on the checkbox component but since then learned it's
   redundant information, so we're removing it again.

   https://github.com/owncloud/owncloud-design-system/pull/1212


* Enhancement - Text-inital helper class: [#1211](https://github.com/owncloud/owncloud-design-system/pull/1211)

   Introducing a `oc-text-inital` class to set font-size, e.g. in headings, back to `1rem`.

   https://github.com/owncloud/owncloud-design-system/pull/1211


* Enhancement - Update tag border color: [#1213](https://github.com/owncloud/owncloud-design-system/pull/1213)

   The OcTag border color didn't offer enough contrast, so it's been updated to match the font
   color. Also, the OcTag styles have been moved to the component level.

   https://github.com/owncloud/owncloud-design-system/pull/1213

# Changelog for [5.0.1] (2021-04-08)

The following sections list the changes in ownCloud Design System 5.0.1.

[5.0.1]: https://github.com/owncloud/owncloud-design-system/compare/v5.0.0...v5.0.1

## Summary

* Bugfix - Add missing peerDependency: [#1205](https://github.com/owncloud/owncloud-design-system/pull/1205)

## Details

* Bugfix - Add missing peerDependency: [#1205](https://github.com/owncloud/owncloud-design-system/pull/1205)

   In the 5.0.0 release, we missed to add the dependency for `vue-inline-svg` to the
   peerDependencies.

   https://github.com/owncloud/owncloud-design-system/pull/1205

# Changelog for [5.0.0] (2021-04-08)

The following sections list the changes in ownCloud Design System 5.0.0.

[5.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v4.3.0...v5.0.0

## Summary

* Change - Use peerDependencies instead of dependencies: [#1202](https://github.com/owncloud/owncloud-design-system/pull/1202)
* Enhancement - Export translations: [#1201](https://github.com/owncloud/owncloud-design-system/pull/1201)

## Details

* Change - Use peerDependencies instead of dependencies: [#1202](https://github.com/owncloud/owncloud-design-system/pull/1202)

   In the past we used dependencies in package.json which can blow up the bundle size a lot. Expect
   this, it is also possible that the same package with 2 versions is part of the bundle.

   From now on dependencies that are required to use ODS are added to the peerDependencies section
   in package.json. Then the consuming application has to add the dependency on it's own and can
   decide which minor or bugfix version to use.

   https://github.com/owncloud/owncloud-design-system/pull/1202


* Enhancement - Export translations: [#1201](https://github.com/owncloud/owncloud-design-system/pull/1201)

   Some ODS components depend on translations and they correctly get pulled from Transifex into
   `l10n/translations.json`, yet we never exported them for other projects to use. Now, they get
   copied into the `dist` folder and can be imported and used alongside the styles and components.

   https://github.com/owncloud/owncloud-design-system/pull/1201

# Changelog for [4.3.0] (2021-04-07)

The following sections list the changes in ownCloud Design System 4.3.0.

[4.3.0]: https://github.com/owncloud/owncloud-design-system/compare/v4.2.1...v4.3.0

## Summary

* Bugfix - Embed svg images: [#1198](https://github.com/owncloud/owncloud-design-system/pull/1198)
* Enhancement - Create `oc-select` component: [#1198](https://github.com/owncloud/owncloud-design-system/pull/1198)

## Details

* Bugfix - Embed svg images: [#1198](https://github.com/owncloud/owncloud-design-system/pull/1198)

   We've changed the way how svg images get loaded in oc-icon. Before this, svg images where
   transformed to data urls by webpack. Now, we use require to inline them, which prevents the
   browser from doing a XMLHttpRequest that can lead to problems with the
   content-security-police.

   https://github.com/owncloud/owncloud-design-system/pull/1198


* Enhancement - Create `oc-select` component: [#1198](https://github.com/owncloud/owncloud-design-system/pull/1198)

   We've created a new component called `oc-select`. This component is based on
   [vue-select](https://vue-select.org/) and can be used to selecting values in a dropdown.

   https://github.com/owncloud/owncloud-design-system/pull/1198

# Changelog for [4.2.1] (2021-04-01)

The following sections list the changes in ownCloud Design System 4.2.1.

[4.2.1]: https://github.com/owncloud/owncloud-design-system/compare/v4.1.2...v4.2.1

## Summary

* Bugfix - Use primary text color for the dropzone: [#1192](https://github.com/owncloud/owncloud-design-system/pull/1192)
* Bugfix - Progressbar background color: [#1192](https://github.com/owncloud/owncloud-design-system/pull/1192)

## Details

* Bugfix - Use primary text color for the dropzone: [#1192](https://github.com/owncloud/owncloud-design-system/pull/1192)

   We've changed the color of the dropzone hint message to the `oc-color` so that the text is
   visible even when brand and background colors are both dark.

   https://github.com/owncloud/owncloud-design-system/pull/1192


* Bugfix - Progressbar background color: [#1192](https://github.com/owncloud/owncloud-design-system/pull/1192)

   The progressbar background color wasn't theme-able and therefore inherited uikit-styles
   that provided undesired results on dark backgrounds. It now uses the text-muted color for the
   background.

   https://github.com/owncloud/owncloud-design-system/pull/1192

# Changelog for [4.1.2] (2021-03-31)

The following sections list the changes in ownCloud Design System 4.1.2.

[4.1.2]: https://github.com/owncloud/owncloud-design-system/compare/v4.2.0...v4.1.2

## Summary

* Bugfix - Breadcrumb font size: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)
* Bugfix - Links in oc-avatar-group: [#1184](https://github.com/owncloud/owncloud-design-system/pull/1184)
* Bugfix - Modal title variation: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)
* Bugfix - Muted text uikit class: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)
* Bugfix - No underline on button hover: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)

## Details

* Bugfix - Breadcrumb font size: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)

   The `lead` class in breadcrumbs didn't resize the font to match `<h1>` tags. This fix brings
   back the old behavior.

   https://github.com/owncloud/owncloud-design-system/pull/1185


* Bugfix - Links in oc-avatar-group: [#1184](https://github.com/owncloud/owncloud-design-system/pull/1184)

   The oc-avatar-group was showing all link, ignoring the `maxDisplayed` property. We fixed
   that by properly cutting off the items used for rendering avatars.

   https://github.com/owncloud/owncloud-design-system/pull/1184


* Bugfix - Modal title variation: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)

   The modal variation was only reflected in the title color in the `danger` variation. Now, the
   title colors gets correctly adapted to any of the 5 possible variations

   https://github.com/owncloud/owncloud-design-system/pull/1185


* Bugfix - Muted text uikit class: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)

   The `uk-text-muted` was overwritten to correctly use CI colors.

   https://github.com/owncloud/owncloud-design-system/pull/1185


* Bugfix - No underline on button hover: [#1185](https://github.com/owncloud/owncloud-design-system/pull/1185)

   When hovering `<button>`s we don't want to display an underline to differentiate them
   properly from `<a>`s and `<router-link>`s.

   https://github.com/owncloud/owncloud-design-system/pull/1185

# Changelog for [4.2.0] (2021-03-31)

The following sections list the changes in ownCloud Design System 4.2.0.

[4.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v4.1.1...v4.2.0

## Summary

* Bugfix - Input border variable: [#1187](https://github.com/owncloud/owncloud-design-system/pull/1187)
* Enhancement - Use oc-color for breadcrumbs: [#1189](https://github.com/owncloud/owncloud-design-system/pull/1189)
* Enhancement - Add `oc-text-lead` class: [#1189](https://github.com/owncloud/owncloud-design-system/pull/1189)
* Enhancement - Unify input colors: [#1190](https://github.com/owncloud/owncloud-design-system/pull/1190)

## Details

* Bugfix - Input border variable: [#1187](https://github.com/owncloud/owncloud-design-system/pull/1187)

   The custom CSS prop for `input-border` had a duplicate and wrong dash and was therefore not
   rendered correctly.

   https://github.com/owncloud/owncloud-design-system/pull/1187


* Enhancement - Use oc-color for breadcrumbs: [#1189](https://github.com/owncloud/owncloud-design-system/pull/1189)

   We've changed the colour of breadcrumbs to use the `oc-color` instead of the brand color.

   https://github.com/owncloud/owncloud-design-system/pull/1189


* Enhancement - Add `oc-text-lead` class: [#1189](https://github.com/owncloud/owncloud-design-system/pull/1189)

   We've added a utility class called `oc-text-lead` which is increasing the font size of the
   text.

   https://github.com/owncloud/owncloud-design-system/pull/1189


* Enhancement - Unify input colors: [#1190](https://github.com/owncloud/owncloud-design-system/pull/1190)

   `OcAutocomplete`, `OcTextInput`, `OcTextarea` and `OcSearchBar` had all slightly
   different ways of using variables to defined border- and text-colors. This change introduces
   two dedicated color varaibles for text inside input fields (for theme-ability) and unifies
   their usage. It also moves the styles for these components from a stylesheet each into a
   `<style>` tag in each Vue component.

   https://github.com/owncloud/owncloud-design-system/pull/1190

# Changelog for [4.1.1] (2021-03-30)

The following sections list the changes in ownCloud Design System 4.1.1.

[4.1.1]: https://github.com/owncloud/owncloud-design-system/compare/v4.1.0...v4.1.1

## Summary

* Bugfix - OcTextInput Coloring: [#1182](https://github.com/owncloud/owncloud-design-system/pull/1182)

## Details

* Bugfix - OcTextInput Coloring: [#1182](https://github.com/owncloud/owncloud-design-system/pull/1182)

   The OcTextInput input validation didn't correctly change to the `danger` or `warning` color
   variations (discovered this when checking the modals in `web`). This change fixes it,
   correctly changing colors when the validation gets triggered.

   https://github.com/owncloud/owncloud-design-system/pull/1182

# Changelog for [4.1.0] (2021-03-26)

The following sections list the changes in ownCloud Design System 4.1.0.

[4.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v4.0.0...v4.1.0

## Summary

* Enhancement - Add messages to autocomplete component: [#1174](https://github.com/owncloud/owncloud-design-system/pull/1174)
* Enhancement - Datepicker input style: [#1176](https://github.com/owncloud/owncloud-design-system/pull/1176)
* Enhancement - Turn placeholders in OcModal into descriptionMessages: [#1172](https://github.com/owncloud/owncloud-design-system/pull/1172)

## Details

* Enhancement - Add messages to autocomplete component: [#1174](https://github.com/owncloud/owncloud-design-system/pull/1174)

   The component supports 3 kinds of messages: description, warning and error. The description
   message will basically replace the informative text which was displayed via placeholder.
   While it's still possible to use a placeholder if you need it, we encourage you to not use it
   anymore. A label is always required, a description message can be added if more context is
   needed.

   https://github.com/owncloud/owncloud-design-system/pull/1174


* Enhancement - Datepicker input style: [#1176](https://github.com/owncloud/owncloud-design-system/pull/1176)

   The input style of the datepicker now matches with the ones we use in other components (like text
   input e.g.).

   https://github.com/owncloud/owncloud-design-system/pull/1176


* Enhancement - Turn placeholders in OcModal into descriptionMessages: [#1172](https://github.com/owncloud/owncloud-design-system/pull/1172)

   We have dropped the support for placeholders in the `<oc-text-input>` component. As this
   component is used inside OcModal, this change also drops placeholder support and adds the
   possibilty to instead make use of a descriptionMessage (which is only being display if no
   inputError is present).

   https://github.com/owncloud/owncloud-design-system/pull/1172

# Changelog for [4.0.0] (2021-03-25)

The following sections list the changes in ownCloud Design System 4.0.0.

[4.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v3.2.2...v4.0.0

## Summary

* Bugfix - Increase the datepicker z-index: [#1147](https://github.com/owncloud/owncloud-design-system/pull/1147)
* Change - Add label to the autocomplete component: [#1156](https://github.com/owncloud/owncloud-design-system/pull/1156)
* Change - Add properties to make icons and images a11y-compliant: [#1161](https://github.com/owncloud/owncloud-design-system/pull/1161)
* Change - Remove URL prop from icons: [#1082](https://github.com/owncloud/owncloud-design-system/pull/1082)
* Change - Refactor variations & color palette: [#1140](https://github.com/owncloud/owncloud-design-system/pull/1140)
* Enhancement - Add focus border to accordion items: [#1158](https://github.com/owncloud/owncloud-design-system/pull/1158)
* Enhancement - Add label and messages to datepicker: [#1149](https://github.com/owncloud/owncloud-design-system/pull/1149)
* Enhancement - Add label and messages to textarea: [#1144](https://github.com/owncloud/owncloud-design-system/pull/1144)
* Enhancement - Make alert component a11y-compliant: [#1166](https://github.com/owncloud/owncloud-design-system/pull/1166)
* Enhancement - Add aria label to the home icon within the breadcrumb: [#1152](https://github.com/owncloud/owncloud-design-system/pull/1152)
* Enhancement - Add aria properties to checkbox: [#1145](https://github.com/owncloud/owncloud-design-system/pull/1145)
* Enhancement - Element nesting: [#1155](https://github.com/owncloud/owncloud-design-system/pull/1155)
* Enhancement - Add aria properties to radio input: [#1148](https://github.com/owncloud/owncloud-design-system/pull/1148)
* Enhancement - Remove the OcTextInput component from the OcSearchBar component: [#1151](https://github.com/owncloud/owncloud-design-system/pull/1151)
* Enhancement - Make searchbar component a11y-compliant: [#1164](https://github.com/owncloud/owncloud-design-system/pull/1164)
* Enhancement - Implement labels and descriptions for text input fields: [#1141](https://github.com/owncloud/owncloud-design-system/pull/1141)
* Enhancement - Initialize the design system with themable colors: [#1168](https://github.com/owncloud/owncloud-design-system/pull/1168)

## Details

* Bugfix - Increase the datepicker z-index: [#1147](https://github.com/owncloud/owncloud-design-system/pull/1147)

   The z-index of the datepicker modal was set to 1, which is a pretty small value for a modal window.
   As a result the date picker modal has been displayed behind the code example in the design
   system.

   https://github.com/owncloud/owncloud-design-system/pull/1147


* Change - Add label to the autocomplete component: [#1156](https://github.com/owncloud/owncloud-design-system/pull/1156)

   Adding the label "from the outside" is not a11y compliant as you cannot get the id of the actual
   input element which therefore cannot be referenced. This change also removes all the
   unnecessary prefixes which were used in the component.

   https://github.com/owncloud/web/issues/4329
   https://github.com/owncloud/owncloud-design-system/pull/1156


* Change - Add properties to make icons and images a11y-compliant: [#1161](https://github.com/owncloud/owncloud-design-system/pull/1161)

   https://github.com/owncloud/owncloud-design-system/pull/1161


* Change - Remove URL prop from icons: [#1082](https://github.com/owncloud/owncloud-design-system/pull/1082)

   We've removed the URL prop from icons which was replacing the internal icons with an image
   loading the icon from arbitrary URLs. As of now, only internal icons can be used in the `oc-icon`
   component.

   https://github.com/owncloud/owncloud-design-system/pull/1082


* Change - Refactor variations & color palette: [#1140](https://github.com/owncloud/owncloud-design-system/pull/1140)

   We have updated the ownCloud Design System colors, removing duplicates, introducing CI
   colors and explicitly adding colors that formerly were calculated or came from uikit.

   We have also unified the usage of "variations" that are used to give visual clues about
   different usages of the same components. Icons, buttons, modals and notifications (and
   perhaps others) can now have the variations `passive, primary, success, danger, warning`.

   While doing that, we also replaced the SASS variables in `src/` with custom CSS properties
   which are overwrite-able at runtime, so theming the ODS is possible now (at least from a colors
   perspective, font sizes and spacing will eventually follow).

   https://github.com/owncloud/owncloud-design-system/pull/1140
   https://github.com/owncloud/owncloud-design-system/pull/1169


* Enhancement - Add focus border to accordion items: [#1158](https://github.com/owncloud/owncloud-design-system/pull/1158)

   https://github.com/owncloud/owncloud-design-system/pull/1158


* Enhancement - Add label and messages to datepicker: [#1149](https://github.com/owncloud/owncloud-design-system/pull/1149)

   This also implies all the necessary accessibility changes.

   https://github.com/owncloud/owncloud-design-system/pull/1149


* Enhancement - Add label and messages to textarea: [#1144](https://github.com/owncloud/owncloud-design-system/pull/1144)

   This also implies all the necessary accessibility changes.

   https://github.com/owncloud/owncloud-design-system/pull/1144


* Enhancement - Make alert component a11y-compliant: [#1166](https://github.com/owncloud/owncloud-design-system/pull/1166)

   It is now possible to reach and submit the "close"-button of an alert via keyboard. This change
   also removes the "uk-close" icon as it is not a11y-compliant.

   https://github.com/owncloud/owncloud-design-system/pull/1166


* Enhancement - Add aria label to the home icon within the breadcrumb: [#1152](https://github.com/owncloud/owncloud-design-system/pull/1152)

   https://github.com/owncloud/web/issues/4329
   https://github.com/owncloud/owncloud-design-system/pull/1152


* Enhancement - Add aria properties to checkbox: [#1145](https://github.com/owncloud/owncloud-design-system/pull/1145)

   This also implies all the necessary accessibility changes.

   https://github.com/owncloud/owncloud-design-system/pull/1145


* Enhancement - Element nesting: [#1155](https://github.com/owncloud/owncloud-design-system/pull/1155)

   Fix the nesting for some html elements to make them w3c compliant.

   https://github.com/owncloud/owncloud-design-system/pull/1155


* Enhancement - Add aria properties to radio input: [#1148](https://github.com/owncloud/owncloud-design-system/pull/1148)

   This also implies all the necessary accessibility changes.

   https://github.com/owncloud/owncloud-design-system/pull/1148


* Enhancement - Remove the OcTextInput component from the OcSearchBar component: [#1151](https://github.com/owncloud/owncloud-design-system/pull/1151)

   Also replace the loading spinner with the OcSpinner component in the course of this.

   https://github.com/owncloud/owncloud-design-system/pull/1151


* Enhancement - Make searchbar component a11y-compliant: [#1164](https://github.com/owncloud/owncloud-design-system/pull/1164)

   It is now possible to reach and submit the "clear"-button within the search bar via keyboard.
   After a search query has been cleared, the input will be focused again. This change also removes
   the "uk-close" icon as it is not a11y-compliant.

   https://github.com/owncloud/owncloud-design-system/issues/1160
   https://github.com/owncloud/owncloud-design-system/pull/1164


* Enhancement - Implement labels and descriptions for text input fields: [#1141](https://github.com/owncloud/owncloud-design-system/pull/1141)

   This also implies all the required accessibility changes.

   https://github.com/owncloud/owncloud-design-system/pull/1141


* Enhancement - Initialize the design system with themable colors: [#1168](https://github.com/owncloud/owncloud-design-system/pull/1168)

   By passing suitable plugin options, you can overwrite the CSS color variables within the ODS to
   adjust it to your likings. The default color values are generated from within the design system
   and act as a fallback if no (or not all) options are passed when initializing the design system
   plugin.

   https://github.com/owncloud/owncloud-design-system/pull/1168

# Changelog for [3.2.2] (2021-03-22)

The following sections list the changes in ownCloud Design System 3.2.2.

[3.2.2]: https://github.com/owncloud/owncloud-design-system/compare/v3.2.1...v3.2.2

## Summary

* Bugfix - Remove leading comma from link avatar group tooltip: [#1165](https://github.com/owncloud/owncloud-design-system/pull/1165)

## Details

* Bugfix - Remove leading comma from link avatar group tooltip: [#1165](https://github.com/owncloud/owncloud-design-system/pull/1165)

   We fixed a bug in the oc-avatar-group component, which showed a leading ", " in its tooltip in
   case it was only about links.

   https://github.com/owncloud/owncloud-design-system/pull/1165

# Changelog for [3.2.1] (2021-03-17)

The following sections list the changes in ownCloud Design System 3.2.1.

[3.2.1]: https://github.com/owncloud/owncloud-design-system/compare/v3.2.0...v3.2.1

## Summary

* Bugfix - Folder names with dots: [#1153](https://github.com/owncloud/owncloud-design-system/pull/1153)

## Details

* Bugfix - Folder names with dots: [#1153](https://github.com/owncloud/owncloud-design-system/pull/1153)

   Folder names with dots showed were separated into a basename and extension part in the
   oc-resource-name component, which has been fixed. Additionally, the oc-resource-name
   component now as a `resource-type` attribute, which reflects the `type` property of the
   displayed resource. This makes testing easier.

   https://github.com/owncloud/owncloud-design-system/pull/1153

# Changelog for [3.2.0] (2021-03-12)

The following sections list the changes in ownCloud Design System 3.2.0.

[3.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v3.0.0...v3.2.0

## Summary

* Bugfix - Use correct color for the accordion title: [#1130](https://github.com/owncloud/owncloud-design-system/pull/1130)
* Bugfix - Update table component in OcIcon docs: [#1127](https://github.com/owncloud/owncloud-design-system/pull/1127)
* Bugfix - Update OcSwitch docs: [#1135](https://github.com/owncloud/owncloud-design-system/pull/1135)
* Enhancement - Use accessible labels for avatars in the avatars group component: [#1126](https://github.com/owncloud/owncloud-design-system/pull/1126)
* Enhancement - Extract SAAS functionality from `src/` to `docs/``: [#1134](https://github.com/owncloud/owncloud-design-system/pull/1134)
* Enhancement - Improve documentation for image and user component: [#1133](https://github.com/owncloud/owncloud-design-system/pull/1133)
* Enhancement - Remove filetype icon variation: [#1139](https://github.com/owncloud/owncloud-design-system/pull/1139)

## Details

* Bugfix - Use correct color for the accordion title: [#1130](https://github.com/owncloud/owncloud-design-system/pull/1130)

   In the last release, we introduced a regression in the accordion component. The custom color
   set for its title had not been respected and the button color was used instead. We've fixed this
   to use the custom color again.

   https://github.com/owncloud/owncloud-design-system/pull/1130


* Bugfix - Update table component in OcIcon docs: [#1127](https://github.com/owncloud/owncloud-design-system/pull/1127)

   The `<oc-table>` component had a breaking change in the *2.1.0* release, which broke the
   [oc-icon component docs](https://owncloud.design/#/oC%20Components/OcIcon). To fix it
   we needed to use the new `<oc-table-simple>` component instead, which mimics the behaviour of
   the old `<oc-table>` component.

   https://github.com/owncloud/owncloud-design-system/pull/1127


* Bugfix - Update OcSwitch docs: [#1135](https://github.com/owncloud/owncloud-design-system/pull/1135)

   The OcSwitch documentation was featuring the `<oc-star>` component which was dropped a while
   ago. This change deletes an obsolete `oc-star.scss` file and updates the OcSwitch
   documentation by displaying its reactive functionality using other HTML elements.

   https://github.com/owncloud/owncloud-design-system/pull/1135


* Enhancement - Use accessible labels for avatars in the avatars group component: [#1126](https://github.com/owncloud/owncloud-design-system/pull/1126)

   We've added the accessible labels for all avatars in the avatars group component because they
   are shown there without any supporting label thus rendering them hidden from screen readers.

   https://github.com/owncloud/owncloud-design-system/pull/1126


* Enhancement - Extract SAAS functionality from `src/` to `docs/``: [#1134](https://github.com/owncloud/owncloud-design-system/pull/1134)

   Some SCSS-files were only used to build the ODS documentation, but lived within the design
   system styles. We moved them to the documentation folder to avoid shipping them with the
   official ODS releases.

   https://github.com/owncloud/owncloud-design-system/pull/1134


* Enhancement - Improve documentation for image and user component: [#1133](https://github.com/owncloud/owncloud-design-system/pull/1133)

   Changed the documentation markup for more semantic HTML that also looks better.

   https://github.com/owncloud/owncloud-design-system/pull/1133


* Enhancement - Remove filetype icon variation: [#1139](https://github.com/owncloud/owncloud-design-system/pull/1139)

   This icon variation wasn't really used anywhere but added duplicate CSS rules and could
   therefore be removed.

   https://github.com/owncloud/owncloud-design-system/pull/1139

# Changelog for [3.0.0] (2021-02-24)

The following sections list the changes in ownCloud Design System 3.0.0.

[3.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v3.1.0...v3.0.0

## Summary

* Bugfix - Button hover color: [#1091](https://github.com/owncloud/owncloud-design-system/pull/1091)
* Bugfix - Missing logo in OcSidebar: [#1090](https://github.com/owncloud/owncloud-design-system/pull/1090)
* Change - Merge resource target path with folder path: [#1085](https://github.com/owncloud/owncloud-design-system/pull/1085)
* Change - Make dropzone component public: [#1100](https://github.com/owncloud/owncloud-design-system/pull/1100)
* Change - Remove basic styles: [#1084](https://github.com/owncloud/owncloud-design-system/pull/1084)
* Change - Add resource path attribute to oc-resource-name: [#1103](https://github.com/owncloud/owncloud-design-system/pull/1103)
* Change - Pass object as the target route: [#1102](https://github.com/owncloud/owncloud-design-system/pull/1102)
* Change - Use values as keys for table cells: [#1094](https://github.com/owncloud/owncloud-design-system/pull/1094)
* Change - Use resources instead of ids in selection model: [#1095](https://github.com/owncloud/owncloud-design-system/pull/1095)
* Change - Text utility classes: [#1115](https://github.com/owncloud/owncloud-design-system/pull/1115)
* Enhancement - Add link avatar component: [#1101](https://github.com/owncloud/owncloud-design-system/pull/1101)
* Enhancement - Add prop for the table header position: [#1092](https://github.com/owncloud/owncloud-design-system/pull/1092)
* Enhancement - Remove the trailing slash from resource target path: [#1087](https://github.com/owncloud/owncloud-design-system/pull/1087)
* Enhancement - Add `disabled` prop to table components: [#1093](https://github.com/owncloud/owncloud-design-system/pull/1093)
* Enhancement - Add table footer: [#1099](https://github.com/owncloud/owncloud-design-system/pull/1099)
* Enhancement - Add `isSelectable` prop to files table component: [#1093](https://github.com/owncloud/owncloud-design-system/pull/1093)

## Details

* Bugfix - Button hover color: [#1091](https://github.com/owncloud/owncloud-design-system/pull/1091)

   When hovering the oc-button (*primary*) component had different styling, depending on if it
   was used as `a` or `button` tag. This was caused by the underlying `ui-kit` styles, since no
   specific component-based styling was provided. By making the hover styling explicit in the
   component stylesheet, this is fixed.

   https://github.com/owncloud/owncloud-design-system/pull/1091


* Bugfix - Missing logo in OcSidebar: [#1090](https://github.com/owncloud/owncloud-design-system/pull/1090)

   After to changes in the formerly referenced structure of the ownCloud.com website, no logo was
   displayed in the sidebar anymore. Due to license changes in this repository, it was decided not
   to include an ownCloud-branded logo anymore, so this change adds a placeholder svg file in
   `assets/example/`.

   https://github.com/owncloud/owncloud-design-system/pull/1090


* Change - Merge resource target path with folder path: [#1085](https://github.com/owncloud/owncloud-design-system/pull/1085)

   We've merged the resource target route path with the folder path in the `router-link`. When
   passing the path as a parameter, it wasn't being recognised and the navigation just pointed to
   the target route path. By explicitly merging them, we make sure that the navigation behaves as
   expected.

   https://github.com/owncloud/owncloud-design-system/pull/1085


* Change - Make dropzone component public: [#1100](https://github.com/owncloud/owncloud-design-system/pull/1100)

   We've made the dropzone component public so that it appears in the documentation.

   https://github.com/owncloud/owncloud-design-system/pull/1100


* Change - Remove basic styles: [#1084](https://github.com/owncloud/owncloud-design-system/pull/1084)

   We've removed styles from the body element and from the #oc-content element. We were forcing
   styles on a global level which were too specific and because of that limited the usage of the
   design system. Any specific styling like that should be done in the consuming app.

   https://github.com/owncloud/owncloud-design-system/pull/1084


* Change - Add resource path attribute to oc-resource-name: [#1103](https://github.com/owncloud/owncloud-design-system/pull/1103)

   We added a custom attribute to the oc-resource-name component which contains the full
   resource path. While this is not needed for rendering the component, this makes end to end
   testing much easier.

   https://github.com/owncloud/owncloud-design-system/pull/1103


* Change - Pass object as the target route: [#1102](https://github.com/owncloud/owncloud-design-system/pull/1102)

   We've changed the `targetRoute` prop in the `oc-resource` component to accept an object
   instead of a string. This object is now accepting keys `name`, `params` and `query`. Only
   `name` is required. All keys are then passed to the router link which enables using more complex
   routes.

   https://github.com/owncloud/owncloud-design-system/pull/1102


* Change - Use values as keys for table cells: [#1094](https://github.com/owncloud/owncloud-design-system/pull/1094)

   We've started using values as keys for the table cells so that they are reactive to the changes in
   values.

   https://github.com/owncloud/owncloud-design-system/pull/1094


* Change - Use resources instead of ids in selection model: [#1095](https://github.com/owncloud/owncloud-design-system/pull/1095)

   We changed the selection model of the `oc-table-files` component to use resources instead of
   their ids.

   https://github.com/owncloud/owncloud-design-system/pull/1095


* Change - Text utility classes: [#1115](https://github.com/owncloud/owncloud-design-system/pull/1115)

   We introduced text utility classes for deciding if a text element should break into multiple
   lines, prevent line breaks or truncate text with an ellipsis.

   https://github.com/owncloud/owncloud-design-system/pull/1115


* Enhancement - Add link avatar component: [#1101](https://github.com/owncloud/owncloud-design-system/pull/1101)

   We've created a link avatar component which mimics styles of the avatar component.

   https://github.com/owncloud/owncloud-design-system/pull/1101


* Enhancement - Add prop for the table header position: [#1092](https://github.com/owncloud/owncloud-design-system/pull/1092)

   We've added a prop which is defining the top position of the table header if the `sticky` prop is
   set to `true`.

   https://github.com/owncloud/owncloud-design-system/pull/1092


* Enhancement - Remove the trailing slash from resource target path: [#1087](https://github.com/owncloud/owncloud-design-system/pull/1087)

   We've removed the trailing slash from the resource target path. We're merging target path with
   the folder path into one and place a slash in between them. Having trailing slash in the target
   path resulted in two slashes in the final path.

   https://github.com/owncloud/owncloud-design-system/pull/1087


* Enhancement - Add `disabled` prop to table components: [#1093](https://github.com/owncloud/owncloud-design-system/pull/1093)

   We've added `disabled` prop to `oc-table` and `oc-table-files` components. This prop is used
   to pass IDs of disabled resources so that they are clearly distinguishable from the other
   resources.

   https://github.com/owncloud/owncloud-design-system/pull/1093


* Enhancement - Add table footer: [#1099](https://github.com/owncloud/owncloud-design-system/pull/1099)

   We've added a footer to the table and table files components. To add any content to the footer,
   use `footer` slot.

   https://github.com/owncloud/owncloud-design-system/pull/1099


* Enhancement - Add `isSelectable` prop to files table component: [#1093](https://github.com/owncloud/owncloud-design-system/pull/1093)

   We've added `isSelectable` prop to `oc-table-files` component. This prop is used to assert
   whether resources in the table can be selected.

   https://github.com/owncloud/owncloud-design-system/pull/1093

# Changelog for [3.1.0] (2021-02-24)

The following sections list the changes in ownCloud Design System 3.1.0.

[3.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v2.1.2...v3.1.0

## Summary

* Enhancement - Add name attribute in the resource name component: [#1119](https://github.com/owncloud/owncloud-design-system/pull/1119)

## Details

* Enhancement - Add name attribute in the resource name component: [#1119](https://github.com/owncloud/owncloud-design-system/pull/1119)

   We've added a `resource-name` data attribute in the `oc-resource-name` component which
   equals the concatenated resource path, name and extension.

   https://github.com/owncloud/owncloud-design-system/pull/1119

# Changelog for [2.1.2] (2021-01-21)

The following sections list the changes in ownCloud Design System 2.1.2.

[2.1.2]: https://github.com/owncloud/owncloud-design-system/compare/v2.1.1...v2.1.2

## Summary

* Change - Pin UIkit version to 3.5.16: [#1064](https://github.com/owncloud/owncloud-design-system/pull/1064)

## Details

* Change - Pin UIkit version to 3.5.16: [#1064](https://github.com/owncloud/owncloud-design-system/pull/1064)

   We've pinned the version of UIkit to 3.5.16. The reason for this is that the version 3.5.17
   caused a regression to our autocomplete component. Since we want to remove UIkit from the ODS in
   the near future, we are not taking the time to come up with a fix. We do not plan to upgrade UIkit
   anymore.

   https://github.com/owncloud/owncloud-design-system/pull/1064

# Changelog for [2.1.1] (2021-01-21)

The following sections list the changes in ownCloud Design System 2.1.1.

[2.1.1]: https://github.com/owncloud/owncloud-design-system/compare/v2.1.0...v2.1.1

## Summary

* Bugfix - Fix uniqueId: [#1060](https://github.com/owncloud/owncloud-design-system/pull/1060)

## Details

* Bugfix - Fix uniqueId: [#1060](https://github.com/owncloud/owncloud-design-system/pull/1060)

   The uniqueId helper function returned a callback instead of a string.

   https://github.com/owncloud/owncloud-design-system/pull/1060

# Changelog for [2.1.0] (2021-01-19)

The following sections list the changes in ownCloud Design System 2.1.0.

[2.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v2.0.4...v2.1.0

## Summary

* Change - Create new table component: [#1177](https://github.com/owncloud/ocis/issues/1177)
* Enhancement - Introduce breakpoints tokens and utility classes: [#1045](https://github.com/owncloud/owncloud-design-system/pull/1045)
* Enhancement - Create new files table component: [#1177](https://github.com/owncloud/ocis/issues/1177)

## Details

* Change - Create new table component: [#1177](https://github.com/owncloud/ocis/issues/1177)

   We've dropped the table component and created a new one from scratch. The new component is
   defining its fields and data via props. To still use more "traditional" table, take a look at
   component `OcTableSimple`.

   https://github.com/owncloud/ocis/issues/1177
   https://github.com/owncloud/owncloud-design-system/pull/1034
   https://github.com/owncloud/owncloud-design-system/pull/1050


* Enhancement - Introduce breakpoints tokens and utility classes: [#1045](https://github.com/owncloud/owncloud-design-system/pull/1045)

   We've added tokens and utility classes to dynamically display/hide content depending on the
   breakpoints.

   https://github.com/owncloud/owncloud-design-system/pull/1045


* Enhancement - Create new files table component: [#1177](https://github.com/owncloud/ocis/issues/1177)

   We've built a new files table component which has a set of pre-defined fields. Which fields will
   be displayed depends on the provided resources. This component also wraps a few other
   components:

   - OcResource - OcAvatarGroup

   https://github.com/owncloud/ocis/issues/1177
   https://github.com/owncloud/owncloud-design-system/pull/1034
   https://github.com/owncloud/owncloud-design-system/pull/1050

# Changelog for [2.0.4] (2020-12-15)

The following sections list the changes in ownCloud Design System 2.0.4.

[2.0.4]: https://github.com/owncloud/owncloud-design-system/compare/v2.0.3...v2.0.4

## Summary

* Bugfix - Positioning of dropzone: [#1052](https://github.com/owncloud/ocis/issues/1052)

## Details

* Bugfix - Positioning of dropzone: [#1052](https://github.com/owncloud/ocis/issues/1052)

   Switch oc-dropzone positioning from fixed to absolute, fixed always orientates by the
   viewport which means it always covers the entire screen. This is something we can't know and the
   decision should be left to the consuming app.

   Instead we use position absolute to just cover the next parent which ist positioned relative to
   the viewport. In situations where the body is higher or wider than the viewport fixed works
   better because absolute won't cover the scrollable parts of the content.

   Since oc apps always do scrolling on their own this is no problem

   https://github.com/owncloud/ocis/issues/1052

# Changelog for [2.0.3] (2020-12-12)

The following sections list the changes in ownCloud Design System 2.0.3.

[2.0.3]: https://github.com/owncloud/owncloud-design-system/compare/v2.0.2...v2.0.3

## Summary

* Bugfix - Breadcrumb icon floating: [#1053](https://github.com/owncloud/ocis/issues/1053)

## Details

* Bugfix - Breadcrumb icon floating: [#1053](https://github.com/owncloud/ocis/issues/1053)

   Home icon in breadcrumb don't break into second line anymore

   https://github.com/owncloud/ocis/issues/1053
   https://github.com/owncloud/owncloud-design-system/pull/997

# Changelog for [2.0.2] (2020-12-07)

The following sections list the changes in ownCloud Design System 2.0.2.

[2.0.2]: https://github.com/owncloud/owncloud-design-system/compare/v2.0.1...v2.0.2

## Summary

* Bugfix - Add `to` prop to `oc-tag` component: [#975](https://github.com/owncloud/owncloud-design-system/pull/975)
* Enhancement - Add new icons: [#975](https://github.com/owncloud/owncloud-design-system/pull/975)

## Details

* Bugfix - Add `to` prop to `oc-tag` component: [#975](https://github.com/owncloud/owncloud-design-system/pull/975)

   We've added an optional `to` prop to the `oc-tag` component so that it is properly assigned to
   the element in case of `router-link` type.

   https://github.com/owncloud/owncloud-design-system/pull/975


* Enhancement - Add new icons: [#975](https://github.com/owncloud/owncloud-design-system/pull/975)

   We've added a `key` and `checklist` icons.

   https://github.com/owncloud/owncloud-design-system/pull/975
   https://fontawesome.com/icons/key?style=solid
   https://fontawesome.com/icons/tasks?style=solid

# Changelog for [2.0.1] (2020-12-02)

The following sections list the changes in ownCloud Design System 2.0.1.

[2.0.1]: https://github.com/owncloud/owncloud-design-system/compare/v2.0.0...v2.0.1

## Summary

* Bugfix - Fix component docs: [#937](https://github.com/owncloud/owncloud-design-system/pull/937)

## Details

* Bugfix - Fix component docs: [#937](https://github.com/owncloud/owncloud-design-system/pull/937)

   The docs were broken after merging Elements and Patterns into a Components section.
   Apparently `Components` is not allowed as section name for technical reasons, so we renamed it
   to `oC Components`.

   https://github.com/owncloud/owncloud-design-system/pull/937

# Changelog for [2.0.0] (2020-11-23)

The following sections list the changes in ownCloud Design System 2.0.0.

[2.0.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.15.0...v2.0.0

## Summary

* Change - Change license to Apache 2: [#933](https://github.com/owncloud/owncloud-design-system/pull/933)
* Change - Move all elements into components section: [#933](https://github.com/owncloud/owncloud-design-system/pull/933)
* Enhancement - Add tag component: [#927](https://github.com/owncloud/owncloud-design-system/pull/927)

## Details

* Change - Change license to Apache 2: [#933](https://github.com/owncloud/owncloud-design-system/pull/933)

   We've changed the license of the design system to Apache 2.

   https://github.com/owncloud/owncloud-design-system/pull/933


* Change - Move all elements into components section: [#933](https://github.com/owncloud/owncloud-design-system/pull/933)

   We've moved all elements into a new section called components.

   https://github.com/owncloud/owncloud-design-system/pull/933


* Enhancement - Add tag component: [#927](https://github.com/owncloud/owncloud-design-system/pull/927)

   We've added a tag component which is to be used for displaying various information.

   https://github.com/owncloud/owncloud-design-system/pull/927

# Changelog for [1.15.0] (2020-11-04)

The following sections list the changes in ownCloud Design System 1.15.0.

[1.15.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.13.0...v1.15.0

## Summary

* Change - Introduce mode for accordion: [#925](https://github.com/owncloud/owncloud-design-system/pull/925)
* Change - Add folder open icon: [#923](https://github.com/owncloud/owncloud-design-system/pull/923)
* Enhancement - Add color prop to buttons: [#924](https://github.com/owncloud/owncloud-design-system/pull/924)

## Details

* Change - Introduce mode for accordion: [#925](https://github.com/owncloud/owncloud-design-system/pull/925)

   The oc-accordion component was allowing both clicks on accordion items and modified property
   values for `expandedId`/`expandedIds` to modify the internal state of the accordion. We now
   introduced a property `mode`, which defaults to `click`, for selecting the mode for
   controlling the internal state. From now on it can be either click or data, but not mixed
   anymore.

   https://github.com/owncloud/owncloud-design-system/pull/925


* Change - Add folder open icon: [#923](https://github.com/owncloud/owncloud-design-system/pull/923)

   We added an icon for `folder open`

   https://github.com/owncloud/owncloud-design-system/pull/923


* Enhancement - Add color prop to buttons: [#924](https://github.com/owncloud/owncloud-design-system/pull/924)

   We've added a new property to the button component to enable setting a different color when the
   variation is set to `raw`. Available colors are `default`, `primary` and `text`.

   https://github.com/owncloud/owncloud-design-system/pull/924

# Changelog for [1.13.0] (2020-10-28)

The following sections list the changes in ownCloud Design System 1.13.0.

[1.13.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.14.0...v1.13.0

## Summary

* Change - New accordion component implementation: [#911](https://github.com/owncloud/owncloud-design-system/pull/911)

## Details

* Change - New accordion component implementation: [#911](https://github.com/owncloud/owncloud-design-system/pull/911)

   We rewrote the accordion component to remove UIKit styles and align with our own styling. Some
   accessibility aspects are already implement, for example expanding and collapsing
   accordion items by pressing space or enter already works. More will come later on.

   https://github.com/owncloud/owncloud-design-system/pull/911

# Changelog for [1.14.0] (2020-10-28)

The following sections list the changes in ownCloud Design System 1.14.0.

[1.14.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.12.2...v1.14.0

## Summary

* Change - Accordion improvements: [#916](https://github.com/owncloud/owncloud-design-system/pull/916)
* Change - Add file version icon: [#917](https://github.com/owncloud/owncloud-design-system/pull/917)

## Details

* Change - Accordion improvements: [#916](https://github.com/owncloud/owncloud-design-system/pull/916)

   We improved the accordion component to allow external control for the expanded/collapsed
   state of accordion items.

   https://github.com/owncloud/owncloud-design-system/pull/916


* Change - Add file version icon: [#917](https://github.com/owncloud/owncloud-design-system/pull/917)

   We added an icon for `file version`

   https://github.com/owncloud/owncloud-design-system/pull/917

# Changelog for [1.12.2] (2020-10-26)

The following sections list the changes in ownCloud Design System 1.12.2.

[1.12.2]: https://github.com/owncloud/owncloud-design-system/compare/v1.12.1...v1.12.2

## Summary

* Change - Adjust styles of buttons with disabled state: [#909](https://github.com/owncloud/owncloud-design-system/pull/909)
* Change - Remove fill color for document icon: [#902](https://github.com/owncloud/owncloud-design-system/pull/902)

## Details

* Change - Adjust styles of buttons with disabled state: [#909](https://github.com/owncloud/owncloud-design-system/pull/909)

   We've changed the background color of buttons with disabled state to properly differentiate
   them from all other buttons.

   https://github.com/owncloud/owncloud-design-system/pull/909


* Change - Remove fill color for document icon: [#902](https://github.com/owncloud/owncloud-design-system/pull/902)

   We've removed the hardcoded fill color of the document icon.

   https://github.com/owncloud/owncloud-design-system/pull/902

# Changelog for [1.12.1] (2020-10-05)

The following sections list the changes in ownCloud Design System 1.12.1.

[1.12.1]: https://github.com/owncloud/owncloud-design-system/compare/v1.12.0...v1.12.1

## Summary

* Bugfix - Icon variation on notification message: [#843](https://github.com/owncloud/owncloud-design-system/issues/843)
* Bugfix - Use correct input background color: [#892](https://github.com/owncloud/owncloud-design-system/pull/892)

## Details

* Bugfix - Icon variation on notification message: [#843](https://github.com/owncloud/owncloud-design-system/issues/843)

   We fixed a bug in the notification message component where the a wrong variation was applied to
   icons.

   https://github.com/owncloud/owncloud-design-system/issues/843
   https://github.com/owncloud/owncloud-design-system/pull/893


* Bugfix - Use correct input background color: [#892](https://github.com/owncloud/owncloud-design-system/pull/892)

   We've moved definition of the input background color from hook into its class.

   https://github.com/owncloud/owncloud-design-system/pull/892

# Changelog for [1.12.0] (2020-10-03)

The following sections list the changes in ownCloud Design System 1.12.0.

[1.12.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.11.1...v1.12.0

## Summary

* Change - Small design improvements: [#890](https://github.com/owncloud/owncloud-design-system/pull/890)
* Change - Add flat style for oc-loader: [#884](https://github.com/owncloud/owncloud-design-system/pull/884)
* Change - Updated checkbox and radiobutton components: [#890](https://github.com/owncloud/owncloud-design-system/pull/890)
* Enhancement - Add margin and padding utility classes: [#890](https://github.com/owncloud/owncloud-design-system/pull/890)

## Details

* Change - Small design improvements: [#890](https://github.com/owncloud/owncloud-design-system/pull/890)

   We've changed a few color design tokens. We've rotated the `link` icon 45 degrees to the left.
   We've changed the spacing values and united them in one design token for both margin and
   padding.

   https://github.com/owncloud/owncloud-design-system/pull/890


* Change - Add flat style for oc-loader: [#884](https://github.com/owncloud/owncloud-design-system/pull/884)

   The oc-loader component now has a `flat` style, which removes the border radius and shrinks the
   height. It can be enabled by setting the `flat` property on the component to `true`. With this
   the visual appearance of the loader is less prominent.

   https://github.com/owncloud/owncloud-design-system/pull/884


* Change - Updated checkbox and radiobutton components: [#890](https://github.com/owncloud/owncloud-design-system/pull/890)

   We updated the OcCheckbox and OcRadio components so that it's possible to use them properly
   with `v-model`. OcRadio had the radiobutton on the right side of the label. We moved it over to
   the left side for consistency.

   https://github.com/owncloud/owncloud-design-system/pull/890


* Enhancement - Add margin and padding utility classes: [#890](https://github.com/owncloud/owncloud-design-system/pull/890)

   We've added utility classes which can assign margin and padding to elements. In their own
   subsections of utilities section in our design system documentation can be found how to use
   them.

   https://github.com/owncloud/owncloud-design-system/pull/890

# Changelog for [1.11.1] (2020-09-21)

The following sections list the changes in ownCloud Design System 1.11.1.

[1.11.1]: https://github.com/owncloud/owncloud-design-system/compare/v1.11.0...v1.11.1

## Summary

* Change - Switch UI icon: [#882](https://github.com/owncloud/owncloud-design-system/pull/882)

## Details

* Change - Switch UI icon: [#882](https://github.com/owncloud/owncloud-design-system/pull/882)

   We've added a new icon for switching between UI versions - or resources in general.

   https://github.com/owncloud/owncloud-design-system/pull/882

# Changelog for [1.11.0] (2020-09-18)

The following sections list the changes in ownCloud Design System 1.11.0.

[1.11.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.10.1...v1.11.0

## Summary

* Bugfix - Pass event in button click handler: [#872](https://github.com/owncloud/owncloud-design-system/pull/872)
* Bugfix - Fix buttons in modal: [#876](https://github.com/owncloud/owncloud-design-system/pull/876)
* Change - Align buttons content to center: [#875](https://github.com/owncloud/owncloud-design-system/pull/875)
* Change - Use flexbox in buttons instead of css grid: [#873](https://github.com/owncloud/owncloud-design-system/pull/873)
* Change - Option for gap on oc-button: [#878](https://github.com/owncloud/owncloud-design-system/pull/878)
* Change - Option for justify-content on oc-button: [#878](https://github.com/owncloud/owncloud-design-system/pull/878)
* Change - Delete pages: [#876](https://github.com/owncloud/owncloud-design-system/pull/876)

## Details

* Bugfix - Pass event in button click handler: [#872](https://github.com/owncloud/owncloud-design-system/pull/872)

   The click event for oc-button now passes the actual click event into the handler.

   https://github.com/owncloud/owncloud-design-system/pull/872


* Bugfix - Fix buttons in modal: [#876](https://github.com/owncloud/owncloud-design-system/pull/876)

   We fixed the button layout of the oc-modal component, which was broken since we introduced our
   new oc-button.

   https://github.com/owncloud/owncloud-design-system/pull/876


* Change - Align buttons content to center: [#875](https://github.com/owncloud/owncloud-design-system/pull/875)

   We've aligned the buttons content to the center.

   https://github.com/owncloud/owncloud-design-system/pull/875


* Change - Use flexbox in buttons instead of css grid: [#873](https://github.com/owncloud/owncloud-design-system/pull/873)

   We've started using flexbox in the button component instead of css grid to position children
   correctly. This gives us more flexibility for the alignment of children.

   https://github.com/owncloud/owncloud-design-system/pull/873


* Change - Option for gap on oc-button: [#878](https://github.com/owncloud/owncloud-design-system/pull/878)

   OcButton now has a property for defining a `gap` value out of xsmall, small, normal, medium,
   large and xlarge. The sizes are the same as in our margin classes.

   https://github.com/owncloud/owncloud-design-system/pull/878


* Change - Option for justify-content on oc-button: [#878](https://github.com/owncloud/owncloud-design-system/pull/878)

   OcButton now has a property for defining a `justify-content` value out of left, center, right,
   space-between, space-around and space-evenly.

   https://github.com/owncloud/owncloud-design-system/pull/878
   https://developer.mozilla.org/en-US/docs/Web/CSS/justify-content


* Change - Delete pages: [#876](https://github.com/owncloud/owncloud-design-system/pull/876)

   We deleted the pages from our documentation. They were showcasing some of our components, but
   hard to maintain.

   https://github.com/owncloud/owncloud-design-system/pull/876

# Changelog for [1.10.1] (2020-09-17)

The following sections list the changes in ownCloud Design System 1.10.1.

[1.10.1]: https://github.com/owncloud/owncloud-design-system/compare/v1.10.0...v1.10.1

## Summary

* Bugfix - Fix wrong extend syntax: [#870](https://github.com/owncloud/owncloud-design-system/pull/870)

## Details

* Bugfix - Fix wrong extend syntax: [#870](https://github.com/owncloud/owncloud-design-system/pull/870)

   We've fixed the wrong extend syntax in oc-user-menu which was written in LESS. This syntax
   caused an issue with loading all styles and e.g. margin utility classes were not working.

   https://github.com/owncloud/owncloud-design-system/pull/870

# Changelog for [1.10.0] (2020-09-16)

The following sections list the changes in ownCloud Design System 1.10.0.

[1.10.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.9.0...v1.10.0

## Summary

* Change - Add prop to hide navigation in sidebar: [#859](https://github.com/owncloud/owncloud-design-system/pull/859)
* Change - Remove icon prop from buttons and use css grid: [#418](https://github.com/owncloud/owncloud-design-system/pull/418)
* Enhancement - Add new icons: [#858](https://github.com/owncloud/owncloud-design-system/pull/858)
* Enhancement - Add visibility off icon: [#855](https://github.com/owncloud/owncloud-design-system/pull/855)

## Details

* Change - Add prop to hide navigation in sidebar: [#859](https://github.com/owncloud/owncloud-design-system/pull/859)

   We added a prop for the sidebar component to hide the navigation. This is a breaking change, as we
   removed the internal logic hiding the navigation entirely in favor of the new prop.

   https://github.com/owncloud/owncloud-design-system/pull/859


* Change - Remove icon prop from buttons and use css grid: [#418](https://github.com/owncloud/owncloud-design-system/pull/418)

   We've removed the prop which was responsible for displaying icons in the button component. To
   add an icon into the button, it needs to be included in the slot together with the text. We've
   added css grid into the button to ensure that all child items of the button will have a predefined
   gutter between them. We've removed all UIKit button styles from the button component.

   https://github.com/owncloud/owncloud-design-system/pull/418
   https://github.com/owncloud/owncloud-design-system/pull/865


* Enhancement - Add new icons: [#858](https://github.com/owncloud/owncloud-design-system/pull/858)

   We've added new icons which can be used to symbolise adding new shares, creating links and
   shared with lists. We've also removed a wrong fill color from visibility off icon.

   https://github.com/owncloud/owncloud-design-system/pull/858
   https://github.com/owncloud/owncloud-design-system/pull/864
   https://fontawesome.com/icons/share-square?style=solid
   https://material.io/resources/icons/?search=group&icon=group_add&style=baseline


* Enhancement - Add visibility off icon: [#855](https://github.com/owncloud/owncloud-design-system/pull/855)

   We've added a new icon to represent state when visibility is off. This can be used e.g. when
   toggling the visibility of a password.

   https://github.com/owncloud/owncloud-design-system/pull/855
   https://material.io/resources/icons/?search=eye&icon=visibility_off&style=baseline

# Changelog for [1.9.0] (2020-07-28)

The following sections list the changes in ownCloud Design System 1.9.0.

[1.9.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.8.0...v1.9.0

## Summary

* Enhancement - Add new variations and sizes of progress bar: [#819](https://github.com/owncloud/owncloud-design-system/pull/819)

## Details

* Enhancement - Add new variations and sizes of progress bar: [#819](https://github.com/owncloud/owncloud-design-system/pull/819)

   We've added new variations and sizes of the progress bar component. Variations change the
   color of the progress bar and can be either `primary` or `warning`. Sizes affect the height of
   the progress bar where it can be either `default` or `small`.

   https://github.com/owncloud/owncloud-design-system/pull/819

# Changelog for [1.8.0] (2020-07-03)

The following sections list the changes in ownCloud Design System 1.8.0.

[1.8.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.7.0...v1.8.0

## Summary

* Change - Use oc-spinner css class in oc-search-bar component: [#802](https://github.com/owncloud/owncloud-design-system/pull/802)
* Enhancement - Show dropdown in autocomplete on focus: [#804](https://github.com/owncloud/owncloud-design-system/pull/804)
* Enhancement - Add lead breadcrumb variation: [#806](https://github.com/owncloud/owncloud-design-system/pull/806)
* Enhancement - Add mainContent slot to the sidebar: [#804](https://github.com/owncloud/owncloud-design-system/pull/804)
* Enhancement - Add move icon: [#807](https://github.com/owncloud/owncloud-design-system/pull/807)

## Details

* Change - Use oc-spinner css class in oc-search-bar component: [#802](https://github.com/owncloud/owncloud-design-system/pull/802)

   UiKit spinner is not supporting IE11. The css classes around oc-spinner have been introduces
   previously and now they are used in the component oc-search-bar as well.

   https://github.com/owncloud/owncloud-design-system/pull/802


* Enhancement - Show dropdown in autocomplete on focus: [#804](https://github.com/owncloud/owncloud-design-system/pull/804)

   In case the input is focused and it still has a value the dropdown will open

   https://github.com/owncloud/owncloud-design-system/pull/804


* Enhancement - Add lead breadcrumb variation: [#806](https://github.com/owncloud/owncloud-design-system/pull/806)

   We've added a lead variation to the breadcrumbs. This variation gives large font size to
   breadcrumb items.

   https://github.com/owncloud/owncloud-design-system/pull/806


* Enhancement - Add mainContent slot to the sidebar: [#804](https://github.com/owncloud/owncloud-design-system/pull/804)

   We've added slot called mainContent into the sidebar component. This slot replaces the
   navigation if defined.

   https://github.com/owncloud/owncloud-design-system/pull/804


* Enhancement - Add move icon: [#807](https://github.com/owncloud/owncloud-design-system/pull/807)

   We've added the material design folder move icon.

   https://github.com/owncloud/owncloud-design-system/pull/807

# Changelog for [1.7.0] (2020-06-17)

The following sections list the changes in ownCloud Design System 1.7.0.

[1.7.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.6.0...v1.7.0

## Summary

* Enhancement - Bring oC primary brand and interaction colors: [#546](https://github.com/owncloud/owncloud-design-system/issues/546)
* Enhancement - Improve the design of buttons: [#784](https://github.com/owncloud/owncloud-design-system/issues/784)
* Enhancement - Increase the logo clearspace: [#786](https://github.com/owncloud/owncloud-design-system/issues/786)
* Enhancement - Automatically focus modal: [#781](https://github.com/owncloud/owncloud-design-system/pull/781)
* Enhancement - Use Source Sans Pro: [#785](https://github.com/owncloud/owncloud-design-system/issues/785)

## Details

* Enhancement - Bring oC primary brand and interaction colors: [#546](https://github.com/owncloud/owncloud-design-system/issues/546)

   We've brought the ownCloud corporate identity brand and interaction colors. The primary
   brand colour is used as a background in the sidebar. Interaction colours are used in buttons and
   links.

   https://github.com/owncloud/owncloud-design-system/issues/546
   https://github.com/owncloud/owncloud-design-system/pull/791


* Enhancement - Improve the design of buttons: [#784](https://github.com/owncloud/owncloud-design-system/issues/784)

   We've added border-radius to buttons, added shadow to the primary button and adjusted
   font-weight and padding.

   https://github.com/owncloud/owncloud-design-system/issues/784
   https://github.com/owncloud/owncloud-design-system/issues/777
   https://github.com/owncloud/owncloud-design-system/pull/791


* Enhancement - Increase the logo clearspace: [#786](https://github.com/owncloud/owncloud-design-system/issues/786)

   We've increased the gutter between top corner of the sidebar and the logo. We've also decreased
   the size of the logo itself.

   https://github.com/owncloud/owncloud-design-system/issues/786
   https://github.com/owncloud/owncloud-design-system/pull/791


* Enhancement - Automatically focus modal: [#781](https://github.com/owncloud/owncloud-design-system/pull/781)

   When the modal is mounted, it receives automatically a focus. The focus is sent directly to the
   modal itself so skipping the wrapping div which works only as a background.

   https://github.com/owncloud/owncloud-design-system/pull/781


* Enhancement - Use Source Sans Pro: [#785](https://github.com/owncloud/owncloud-design-system/issues/785)

   We've started using Source Sans Pro in the default theme as the font.

   https://github.com/owncloud/owncloud-design-system/issues/785
   https://github.com/owncloud/owncloud-design-system/pull/791

# Changelog for [1.6.0] (2020-05-26)

The following sections list the changes in ownCloud Design System 1.6.0.

[1.6.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.5.0...v1.6.0

## Summary

* Bugfix - Made modal position fixed: [#769](https://github.com/owncloud/owncloud-design-system/pull/769)
* Change - Removed change and keydown events from text input: [#768](https://github.com/owncloud/owncloud-design-system/pull/768)

## Details

* Bugfix - Made modal position fixed: [#769](https://github.com/owncloud/owncloud-design-system/pull/769)

   We've made the position of modal fixed and added z-index so it is always going to be visible on top
   of the content.

   https://github.com/owncloud/owncloud-design-system/pull/769


* Change - Removed change and keydown events from text input: [#768](https://github.com/owncloud/owncloud-design-system/pull/768)

   We've removed change and keydown custom events from text input component. All listeners are
   passed to the input element so all events are still accessible. Focus and input events are still
   implemented as custom events.

   https://github.com/owncloud/owncloud-design-system/pull/768

# Changelog for [1.5.0] (2020-05-13)

The following sections list the changes in ownCloud Design System 1.5.0.

[1.5.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.4.0...v1.5.0

## Summary

* Bugfix - Do not try to emit value after confirming modal if input is missing: [#749](https://github.com/owncloud/owncloud-design-system/pull/749)
* Bugfix - Disable confirm button in the modal component if there is an error: [#741](https://github.com/owncloud/owncloud-design-system/pull/741)
* Bugfix - Use input event instead of keydown in the modal component: [#741](https://github.com/owncloud/owncloud-design-system/pull/741)
* Bugfix - Do not mutate input value prop in the modal directly: [#741](https://github.com/owncloud/owncloud-design-system/pull/741)
* Change - Do not enforce muted color as background of app bar: [#750](https://github.com/owncloud/owncloud-design-system/pull/750)
* Change - Deprecated application menu component: [#735](https://github.com/owncloud/owncloud-design-system/pull/735)
* Enhancement - Created sidebar component: [#735](https://github.com/owncloud/owncloud-design-system/pull/735)
* Enhancement - Created animation section in docs: [#735](https://github.com/owncloud/owncloud-design-system/pull/735)

## Details

* Bugfix - Do not try to emit value after confirming modal if input is missing: [#749](https://github.com/owncloud/owncloud-design-system/pull/749)

   Confirming modal resulted in an error if the modal haven't got an input. We've fixed this by not
   attempting to emit the value if the prop `hasInput` is set to false.

   https://github.com/owncloud/owncloud-design-system/pull/749


* Bugfix - Disable confirm button in the modal component if there is an error: [#741](https://github.com/owncloud/owncloud-design-system/pull/741)

   We've added check to disable confirm button in the modal component in case the input inside of
   the modal has an error.

   https://github.com/owncloud/owncloud-design-system/pull/741


* Bugfix - Use input event instead of keydown in the modal component: [#741](https://github.com/owncloud/owncloud-design-system/pull/741)

   We've started using the input event instead of the keydown in the modal component to properly
   pass the value of the input.

   https://github.com/owncloud/owncloud-design-system/pull/741


* Bugfix - Do not mutate input value prop in the modal directly: [#741](https://github.com/owncloud/owncloud-design-system/pull/741)

   We've used a local state as a v-model of the input in the modal to avoid direct mutations of the
   input value property.

   https://github.com/owncloud/owncloud-design-system/pull/741


* Change - Do not enforce muted color as background of app bar: [#750](https://github.com/owncloud/owncloud-design-system/pull/750)

   We've stopped enforcing the muted color as a background of the app bar.

   https://github.com/owncloud/owncloud-design-system/pull/750


* Change - Deprecated application menu component: [#735](https://github.com/owncloud/owncloud-design-system/pull/735)

   We've deprecated the application menu component in favor of the sidebar component.

   https://github.com/owncloud/owncloud-design-system/pull/735


* Enhancement - Created sidebar component: [#735](https://github.com/owncloud/owncloud-design-system/pull/735)

   We've created a sidebar component which is to be used as an in-app navigation and which will
   contain the branding.

   https://github.com/owncloud/owncloud-design-system/pull/735


* Enhancement - Created animation section in docs: [#735](https://github.com/owncloud/owncloud-design-system/pull/735)

   We've added an animation section in the documentation to show available animations and how to
   use them.

   https://github.com/owncloud/owncloud-design-system/pull/735

# Changelog for [1.4.0] (2020-04-29)

The following sections list the changes in ownCloud Design System 1.4.0.

[1.4.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.3.0...v1.4.0

## Summary

* Enhancement - Extended the modal component with input: [#730](https://github.com/owncloud/owncloud-design-system/pull/730)

## Details

* Enhancement - Extended the modal component with input: [#730](https://github.com/owncloud/owncloud-design-system/pull/730)

   We've added an input into the modal component which can be displayed via prop. If the input is
   displayed, the message gets overridden. The content slot can override the input. In the
   confirm event is now emitted the value of input.

   https://github.com/owncloud/owncloud-design-system/pull/730

# Changelog for [1.3.0] (2020-04-23)

The following sections list the changes in ownCloud Design System 1.3.0.

[1.3.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.2.2...v1.3.0

## Summary

* Enhancement - Added the modal component: [#723](https://github.com/owncloud/owncloud-design-system/pull/723)

## Details

* Enhancement - Added the modal component: [#723](https://github.com/owncloud/owncloud-design-system/pull/723)

   We've added the modal component together with a basic documentation.

   https://github.com/owncloud/owncloud-design-system/pull/723

# Changelog for [1.2.2] (2020-04-08)

The following sections list the changes in ownCloud Design System 1.2.2.

[1.2.2]: https://github.com/owncloud/owncloud-design-system/compare/v1.2.1...v1.2.2

## Summary

* Bugfix - Fix oc-autocomplete: [#710](https://github.com/owncloud/owncloud-design-system/pull/710)

## Details

* Bugfix - Fix oc-autocomplete: [#710](https://github.com/owncloud/owncloud-design-system/pull/710)

   We fixed a bug in OcAutocomplete which was introduced with the removal of lodash as a
   dependency.

   https://github.com/owncloud/owncloud-design-system/pull/710

# Changelog for [1.2.1] (2020-04-07)

The following sections list the changes in ownCloud Design System 1.2.1.

[1.2.1]: https://github.com/owncloud/owncloud-design-system/compare/v1.2.0...v1.2.1

## Summary

* Bugfix - Correct layout of search bar: [#706](https://github.com/owncloud/owncloud-design-system/pull/706)
* Enhancement - Bind attributes and events to input in oc-text-input: [#706](https://github.com/owncloud/owncloud-design-system/pull/706)

## Details

* Bugfix - Correct layout of search bar: [#706](https://github.com/owncloud/owncloud-design-system/pull/706)

   We've fixed layout of search bar which was broken after we introduced error message in
   oc-text-input.

   https://github.com/owncloud/owncloud-design-system/pull/706


* Enhancement - Bind attributes and events to input in oc-text-input: [#706](https://github.com/owncloud/owncloud-design-system/pull/706)

   We've binded attributes and events to input in oc-text-input so that they are passed properly
   instead of passing them to the root element.

   https://github.com/owncloud/owncloud-design-system/pull/706

# Changelog for [1.2.0] (2020-03-27)

The following sections list the changes in ownCloud Design System 1.2.0.

[1.2.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.1.1...v1.2.0

## Summary

* Bugfix - Fixed oc-text-input not appearing in generated docs: [#688](https://github.com/owncloud/owncloud-design-system/issues/688)
* Change - Render error and warning messages in oc-text-input: [#689](https://github.com/owncloud/owncloud-design-system/issues/689)

## Details

* Bugfix - Fixed oc-text-input not appearing in generated docs: [#688](https://github.com/owncloud/owncloud-design-system/issues/688)

   The oc-text-input element didn't appear in the generated docs because of a compilation
   warning.

   https://github.com/owncloud/owncloud-design-system/issues/688
   https://github.com/owncloud/owncloud-design-system/pull/690


* Change - Render error and warning messages in oc-text-input: [#689](https://github.com/owncloud/owncloud-design-system/issues/689)

   The oc-text-input element now has properties for rendering a warning or error message below
   the input field. The input field border receives the respective matching color as well. Also
   it's possible to reserve a fixed vertical space below the input element, so that an appearing
   message doesn't break the layout around the input element.

   https://github.com/owncloud/owncloud-design-system/issues/689
   https://github.com/owncloud/owncloud-design-system/pull/690

# Changelog for [1.1.1] (2020-03-18)

The following sections list the changes in ownCloud Design System 1.1.1.

[1.1.1]: https://github.com/owncloud/owncloud-design-system/compare/v1.1.0...v1.1.1

## Summary

* Bugfix - Fixed oc-icon to reload image when url changes: [#681](https://github.com/owncloud/owncloud-design-system/pull/681)

## Details

* Bugfix - Fixed oc-icon to reload image when url changes: [#681](https://github.com/owncloud/owncloud-design-system/pull/681)

   The oc-icon component will now reload its image whenever the url property has been modified

   https://github.com/owncloud/owncloud-design-system/pull/681

# Changelog for [1.1.0] (2020-03-17)

The following sections list the changes in ownCloud Design System 1.1.0.

[1.1.0]: https://github.com/owncloud/owncloud-design-system/compare/v1.0.4...v1.1.0

## Summary

* Bugfix - Made scrollbar styles consistent: [#660](https://github.com/owncloud/owncloud-design-system/pull/660)
* Bugfix - Removed width class from sidebar menu: [#668](https://github.com/owncloud/owncloud-design-system/issues/668)
* Enhancement - Added iconUrl to oc-file element: [#678](https://github.com/owncloud/owncloud-design-system/pull/678)

## Details

* Bugfix - Made scrollbar styles consistent: [#660](https://github.com/owncloud/owncloud-design-system/pull/660)

   Scrollbar styles are now more consistent between Chrome and Firefox.

   https://github.com/owncloud/owncloud-design-system/pull/660


* Bugfix - Removed width class from sidebar menu: [#668](https://github.com/owncloud/owncloud-design-system/issues/668)

   There were different values for width of the sidebar menu and it's left position when hidden.
   We've removed the width class so that the width and left position are the same and the sidebar
   menu is no more overlapping when it's state is hidden.

   https://github.com/owncloud/owncloud-design-system/issues/668
   https://github.com/owncloud/owncloud-design-system/pull/669


* Enhancement - Added iconUrl to oc-file element: [#678](https://github.com/owncloud/owncloud-design-system/pull/678)

   The oc-file element now supports passing an arbitrary URL to be displayed as a file thumbnail.
   It will fall back to the icon name in case the thumbnail could not be loaded.

   https://github.com/owncloud/owncloud-design-system/pull/678

# Changelog for [1.0.4] (2020-02-26)

The following sections list the changes in ownCloud Design System 1.0.4.

[1.0.4]: https://github.com/owncloud/owncloud-design-system/compare/v1.0.3...v1.0.4

## Summary

* Bugfix - Removed size for oc-button with raw variation: [#654](https://github.com/owncloud/owncloud-design-system/issues/654)

## Details

* Bugfix - Removed size for oc-button with raw variation: [#654](https://github.com/owncloud/owncloud-design-system/issues/654)

   Raw variation of buttons have no border, so they now also don't have a size enforced to avoid
   needless white space.

   https://github.com/owncloud/owncloud-design-system/issues/654

# Changelog for [1.0.3] (2020-02-25)

The following sections list the changes in ownCloud Design System 1.0.3.

[1.0.3]: https://github.com/owncloud/owncloud-design-system/compare/v1.0.2...v1.0.3

## Summary

* Bugfix - Removed uppercase on buttons: [#442](https://github.com/owncloud/owncloud-design-system/issues/442)

## Details

* Bugfix - Removed uppercase on buttons: [#442](https://github.com/owncloud/owncloud-design-system/issues/442)

   Buttons look nicer without the uppercase which was brought in by default by UIKit.

   https://github.com/owncloud/owncloud-design-system/issues/442
   https://github.com/owncloud/owncloud-design-system/pull/652

# Changelog for [1.0.2] (2020-02-24)

The following sections list the changes in ownCloud Design System 1.0.2.

[1.0.2]: https://github.com/owncloud/owncloud-design-system/compare/v1.0.1...v1.0.2

## Summary

* Change - Small sidebar navigation design improvements: [#583](https://github.com/owncloud/owncloud-design-system/issues/583)

## Details

* Change - Small sidebar navigation design improvements: [#583](https://github.com/owncloud/owncloud-design-system/issues/583)

   We've changed the padding of sidebar nav to small and removed text transformation to
   uppercase.

   https://github.com/owncloud/owncloud-design-system/issues/583
   https://github.com/owncloud/owncloud-design-system/pull/648

# Changelog for [1.0.1] (2020-01-30)

The following sections list the changes in ownCloud Design System 1.0.1.

[1.0.1]: https://github.com/owncloud/owncloud-design-system/compare/9c9d14acc7df90254e857e9c2de7bad8f77a238c...v1.0.1

## Summary

* Enhancement - Initial release: [#10](https://github.com/owncloud/owncloud-design-system/issues/10)

## Details

* Enhancement - Initial release: [#10](https://github.com/owncloud/owncloud-design-system/issues/10)

   We created the ownCloud Design System to provide developers with well documents and reusable
   vue.js components.

   https://github.com/owncloud/owncloud-design-system/issues/10
   https://owncloud.github.io/owncloud-design-system/

