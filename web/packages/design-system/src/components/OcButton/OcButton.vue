<template>
  <component
    :is="type"
    v-bind="additionalAttributes"
    :aria-label="ariaLabel || null"
    :class="ocButton_buttonClass"
    v-on="handlers"
  >
    <oc-spinner v-if="showSpinner" size="small" class="spinner" />
    <!-- @slot Content of the button -->
    <slot />
  </component>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { RouteLocationRaw } from 'vue-router'
import { getSizeClass } from '../../helpers'

/**
 * OcButton - A versatile button component with various styling options and types.
 *
 * @prop {'button' | 'a' | 'router-link'} [type='button'] - The HTML element used for the button.
 * @prop {boolean} [disabled=false] - Whether the button is disabled.
 * @prop {'small' | 'medium' | 'large'} [size='medium'] - The size of the button.
 * @prop {string} [href=null] - When type is 'a', specifies the URL to link to.
 * @prop {'_blank' | '_self' | '_parent' | '_top'} [target=null] - When type is 'a', specifies where to open the linked document.
 * @prop {RouteLocationRaw} [to=null] - When type is 'router-link', specifies the target route.
 * @prop {string} [ariaLabel=null] - An accessible label for the button. Use only when needed to override the visible text.
 * @prop {'button' | 'submit' | 'reset'} [submit='button'] - The type attribute when rendering a button element.
 * @prop {'passive' | 'primary' | 'danger' | 'success' | 'warning' | 'brand'} [variation='passive'] - The color variation of the button.
 * @prop {'outline' | 'filled' | 'raw' | 'raw-inverse'} [appearance='outline'] - The visual style of the button.
 *   - outline: transparent with colored border
 *   - filled: colored background with contrasting text
 *   - raw: text-only button with colored text
 *   - raw-inverse: text-only button with contrasting text
 * @prop {'left' | 'center' | 'right' | 'space-around' | 'space-between' | 'space-evenly'} [justifyContent='center'] - How to align content within the button.
 * @prop {'none' | 'xsmall' | 'small' | 'medium' | 'large' | 'xlarge'} [gapSize='medium'] - Spacing between child elements.
 * @prop {boolean} [showSpinner=false] - Whether to show a loading spinner.
 *
 * @slot default - Content of the button.
 *
 * @event {MouseEvent} click - Emitted when the button is clicked.
 *
 * @example
 * ```vue
 * <!-- Basic button -->
 * <oc-button>Click me</oc-button>
 *
 * <!-- Primary filled button -->
 * <oc-button variation="primary" appearance="filled">Submit</oc-button>
 *
 * <!-- Button with icon -->
 * <oc-button>
 *   <oc-icon name="home" />
 *   Home
 * </oc-button>
 *
 * <!-- Link button -->
 * <oc-button type="a" href="https://example.com" target="_blank">External Link</oc-button>
 *
 * <!-- Button with loading state -->
 * <oc-button :show-spinner="isLoading" @click="handleClick">Save</oc-button>
 * ```
 */

interface Props {
  type?: 'button' | 'a' | 'router-link'
  disabled?: boolean
  size?: 'small' | 'medium' | 'large'
  href?: string
  target?: '_blank' | '_self' | '_parent' | '_top'
  to?: RouteLocationRaw
  ariaLabel?: string
  submit?: 'button' | 'submit' | 'reset'
  variation?: 'passive' | 'primary' | 'danger' | 'success' | 'warning' | 'brand' | string
  appearance?: 'outline' | 'filled' | 'raw' | 'raw-inverse' | string
  justifyContent?: 'left' | 'center' | 'right' | 'space-around' | 'space-between' | 'space-evenly'
  gapSize?: 'none' | 'xsmall' | 'small' | 'medium' | 'large' | 'xlarge'
  showSpinner?: boolean
}
interface Emits {
  (e: 'click', event: MouseEvent): void
}
defineOptions({
  name: 'OcButton',
  status: 'ready',
  release: '1.0.0'
})
const {
  type = 'button',
  disabled = false,
  size = 'medium',
  href = null,
  target = null,
  to = null,
  ariaLabel = null,
  submit = 'button',
  variation = 'passive',
  appearance = 'outline',
  justifyContent = 'center',
  gapSize = 'medium',
  showSpinner = false
} = defineProps<Props>()

const emit = defineEmits<Emits>()

function ocButton_onClick(event: MouseEvent) {
  emit('click', event)
}

const ocButton_buttonClass = computed(() => {
  return [
    'oc-button',
    'oc-rounded',
    `oc-button-${getSizeClass(size)}`,
    `oc-button-justify-content-${justifyContent}`,
    `oc-button-gap-${getSizeClass(gapSize)}`,
    `oc-button-${variation}`,
    `oc-button-${variation}-${appearance}`
  ]
})

const additionalAttributes = computed(() => {
  return {
    ...(href && { href: href }),
    ...(target && { target }),
    ...(to && { to }),
    ...(type === 'button' && { type: submit }),
    ...(type === 'button' && { disabled })
  }
})

const handlers = computed(() => {
  return {
    ...(type === 'button' && { click: ocButton_onClick })
  }
})
</script>

<style lang="scss">
@mixin oc-button-gap($factor) {
  gap: math.round(calc($oc-space-small * $factor / 2)) * 2;
}

@mixin oc-button-line-height($factor) {
  line-height: $oc-size-icon-default * $factor;
}

@mixin oc-button-variation($color, $hover-color, $muted-color, $contrast-color) {
  &:disabled {
    background-color: $muted-color;
  }

  &:focus:not([disabled]),
  &:hover:not([disabled]) {
    background-color: $hover-color;
  }

  &-raw,
  &-raw-inverse {
    background-color: transparent;
    border-style: none;
    font-size: var(--oc-font-size-medium);
    font-weight: normal;
    min-height: 0;
    padding: 0;

    &:focus:not([disabled]),
    &:hover:not([disabled]) {
      background-color: transparent;
    }

    &:focus:not([disabled]):not(button),
    &:hover:not([disabled]):not(button) {
      text-decoration: underline;
    }

    &:disabled {
      background-color: transparent;
      color: $muted-color;
    }
  }
  &-raw {
    color: $color;

    .oc-icon > svg {
      fill: $color;
    }

    &:focus:not([disabled]),
    &:hover:not([disabled]) {
      color: $hover-color;

      .oc-icon > svg {
        fill: $hover-color;
      }
    }
  }
  &-raw-inverse {
    color: $contrast-color;

    .oc-icon > svg {
      fill: $contrast-color;
    }

    &:focus:not([disabled]),
    &:hover:not([disabled]) {
      color: $contrast-color;

      .oc-icon > svg {
        fill: $contrast-color;
      }
    }
  }

  &-filled {
    background-color: $color;
    color: $contrast-color;

    .oc-icon > svg {
      fill: $contrast-color;
    }

    &:hover:not([disabled]),
    &:focus:not([disabled]) {
      color: $contrast-color;
      background-color: $hover-color;
      border-color: $hover-color;

      .oc-icon > svg {
        fill: $contrast-color;
      }
    }
  }

  &-outline {
    outline: 1px solid $color;
    outline-offset: -1px;
    background-color: transparent;
    color: $color;

    .oc-icon > svg {
      fill: $color;
    }

    &:disabled {
      background-color: transparent;
      color: $muted-color;
    }
  }

  &-outline:hover:not([disabled]),
  &-outline:focus:not([disabled]) {
    color: $contrast-color;
    background-color: $hover-color;
    border-color: $hover-color;

    .oc-icon > svg {
      fill: $contrast-color;
    }
  }
}

.oc-button {
  @include oc-button-line-height(1);

  align-items: center;
  border: 0;
  box-sizing: border-box;
  display: inline-flex;
  font-weight: 400;
  padding: 0.5rem 0.8rem;
  text-align: left;
  text-decoration: none;

  &-justify-content {
    &-left {
      justify-content: left;
    }

    &-center {
      justify-content: center;
    }

    &-right {
      justify-content: right;
    }

    &-space-between {
      justify-content: space-between;
    }

    &-space-around {
      justify-content: space-around;
    }

    &-space-evenly {
      justify-content: space-evenly;
    }
  }

  &-gap {
    &-xs {
      @include oc-button-gap(0.5);
    }

    &-s {
      @include oc-button-gap(0.7);
    }

    &-m {
      @include oc-button-gap(1);
    }

    &-l {
      @include oc-button-gap(1.5);
    }

    &-xl {
      @include oc-button-gap(2);
    }
  }

  &:hover {
    cursor: pointer;
  }

  &-s {
    @include oc-button-line-height(0.7);

    font-size: var(--oc-font-size-small);
    min-height: 1.2rem;
  }

  &-m {
    @include oc-button-line-height(1);

    font-size: var(--oc-font-size-medium);
    min-height: $global-control-height;
  }

  &-l {
    @include oc-button-line-height(1.5);

    font-size: var(--oc-font-size-xlarge);
    min-height: 2rem;
  }

  &-passive {
    @include oc-button-variation(
      var(--oc-color-swatch-passive-default),
      var(--oc-color-swatch-passive-hover),
      var(--oc-color-swatch-passive-muted),
      var(--oc-color-swatch-passive-contrast)
    );

    &-outline {
      &:focus:not([disabled]),
      &:hover:not([disabled]) {
        color: var(--oc-color-swatch-passive-contrast);
        background-color: var(--oc-color-swatch-passive-hover-outline);
        border-color: var(--oc-color-swatch-passive-hover-outline);

        .oc-icon > svg {
          fill: var(--oc-color-swatch-passive-contrast);
        }
      }
    }
  }

  &-brand {
    @include oc-button-variation(
      var(--oc-color-swatch-brand-default),
      var(--oc-color-swatch-brand-hover),
      var(--oc-color-swatch-brand-muted),
      var(--oc-color-swatch-brand-contrast)
    );
  }

  &-primary {
    @include oc-button-variation(
      var(--oc-color-swatch-primary-default),
      var(--oc-color-swatch-primary-hover),
      var(--oc-color-swatch-primary-muted),
      var(--oc-color-swatch-primary-contrast)
    );

    &-filled {
      color: var(--oc-color-swatch-primary-contrast) !important;
      @include oc-background-primary-gradient;

      span > svg {
        fill: var(--oc-color-swatch-primary-contrast) !important;
      }
    }
  }

  &-success {
    @include oc-button-variation(
      var(--oc-color-swatch-success-default),
      var(--oc-color-swatch-success-hover),
      var(--oc-color-swatch-success-muted),
      var(--oc-color-swatch-success-contrast)
    );
  }

  &-warning {
    @include oc-button-variation(
      var(--oc-color-swatch-warning-default),
      var(--oc-color-swatch-warning-hover),
      var(--oc-color-swatch-warning-muted),
      var(--oc-color-swatch-warning-contrast)
    );
  }

  &-danger {
    @include oc-button-variation(
      var(--oc-color-swatch-danger-default),
      var(--oc-color-swatch-danger-hover),
      var(--oc-color-swatch-danger-muted),
      var(--oc-color-swatch-danger-contrast)
    );

    &-filled:hover,
    &-filled:focus {
      color: var(--oc-color-swatch-danger-default) !important;
      span > svg {
        fill: var(--oc-color-swatch-danger-default) !important;
      }
    }
  }

  &:disabled {
    cursor: default;
    opacity: 0.6;
  }

  &-group {
    display: inline-flex;
    flex-flow: row wrap;
    outline: 1px solid var(--oc-color-swatch-passive-default);
    outline-offset: -1px;
    border-radius: 5px;

    .oc-button {
      border-radius: 0;
      outline: 0;

      &:first-of-type {
        border-radius: 5px 0 0 5px;
      }

      &:last-of-type {
        border-radius: 0 5px 5px 0;
      }

      &-default {
        border-left: 0;
        border-right: 0;

        &:first-of-type {
          // TODO: Implement color variations if button group should be used again
          border-left: 1px solid var(--oc-color-swatch-primary-default);
        }

        &:last-of-type {
          // TODO: Implement color variations if button group should be used again
          border-right: 1px solid var(--oc-color-swatch-primary-default);
        }
      }
    }
  }
}
</style>

<docs>
Buttons are generally used for interface actions. Suitable for all-purpose use.

Defaults to appearance that has white background with blue border.
Filled appearance should be used only once per view for main call-to-action.
All buttons are built with a css grid which ensures that there will be a pre-defined gutter between all child items.

## Accessibility
### Distinction when to use a `
<button>`, when to use an `<a>`

  Regardless of the visual representation the following differentiation should be made if in doubt what to choose.

  - an anchor/link changes the location, be it internally on the website, or to another resource/document/route.
  - a button does change the state of the website, e.g.: adds, deletes, opens, ...

  ### Accessible name
  The accessible name ([explainer for the
  term](https://developer.paciellogroup.com/blog/2017/04/what-is-an-accessible-name/)) of a button is derived from
  several sources. Right now, only two of them are relevant:

  1. The value of the `aria-label` attribute
  2. The text between the opening and closing tag, `
  <button>This text here</button>
  `

  When an aria-label attribute exists, its value will override the button text. So in this case, the accessible name
  would be "foo": `
  <button aria-label="foo">Bar</button>
  `, although visual users will see "Bar". This difference between accessible name and visual name is a problem for a
  certain type of assistive technology ([explainer for the term](https://en.wikipedia.org/wiki/Assistive_technology)),
  this is why the [WCAG success criterion 2.5.3, "Label in name"](https://www.w3.org/TR/WCAG21/#label-in-name) exists.
  This difference should be avoided, if it can't, W3C recommends that the accessible name should start with visible
  label.

  ### Icon-only button

  Every button has to have an accessible name. It cannot be provided by a text between the button tags (for example,
  because it is an icon button; see "Upload" example below), the accessible name has to be provided by `aria-label`.

  ## Examples

  ```js
  <template>
    <div>

      <h3 class="oc-heading-divider">
        Button variations
      </h3>
      <div>
        <oc-button v-for="variation in variations" :variation="variation.title"
                   :class="{ 'oc-background-brand': variation.title === 'inverse' }" class="oc-mr-s oc-mb-s">
          {{ variation.title }}
        </oc-button>
      </div>

      <h3 class="oc-heading-divider">
        Button sizes
      </h3>
      <div class="oc-flex">
        <oc-button size="large" class="oc-mr-s">Large</oc-button>
        <oc-button class="oc-mr-s oc-align-self-center">Medium (default)</oc-button>
        <oc-button size="small" class="oc-align-self-center">Small</oc-button>
      </div>

      <h3 class="oc-heading-divider">
        Button with icons
      </h3>
      <div class="oc-flex">
        <oc-button class="oc-mr-s">
          <oc-icon name="home-2"/>
          Home
        </oc-button>
        <oc-button variation="primary" class="oc-mr-s">
          Select
          <oc-icon name="arrow-drop-down" fill-type="line"/>
        </oc-button>
        <oc-button variation="primary" aria-label="Upload a file">
          <oc-icon name="upload-cloud-2" fill-type="line"/>
        </oc-button>
      </div>

      <h3 class="oc-heading-divider">
        Using buttons in a group
      </h3>
      <div class="oc-button-group">
        <oc-button variation="primary" appearance="filled">Hello</oc-button>
        <oc-button>
          <oc-icon name="folder"/>
          Demo Button
        </oc-button>
        <oc-button variation="danger" appearance="filled">Delete</oc-button>
      </div>

      <h3 class="oc-heading-divider">
        All available button combinations
      </h3>
      <oc-table-simple>
        <oc-thead>
          <oc-tr>
            <oc-th>Variation &amp; usage</oc-th>
            <oc-th>Outline (default)</oc-th>
            <oc-th>Filled</oc-th>
            <oc-th>Raw</oc-th>
          </oc-tr>
        </oc-thead>
        <oc-tbody>
          <oc-tr v-for="variation in variations">
            <oc-td>
              {{ variation.description }}
            </oc-td>
            <oc-td v-for="appearance in appearances" :class="{ 'oc-background-brand': variation.title === 'inverse' }"
                   class="oc-p-m">
              <oc-button :variation="variation.title" :appearance="appearance" class="oc-mb-s">
                {{ variation.title }}
                <oc-icon name="home-2"/>
              </oc-button>
              <oc-button :variation="variation.title" :appearance="appearance" disabled>
                {{ variation.title }}
                <oc-icon name="home-2"/>
              </oc-button>
            </oc-td>
          </oc-tr>
        </oc-tbody>
      </oc-table-simple>
    </div>
  </template>
  <script>
    export default {
      computed: {
        variations() {
          return [
            {
              title: "passive",
              description: "Default appearance of interactive elements (links and buttons)"
            },
            {
              title: "primary",
              description: "Used as a visual clue to the main action on each page"
            },
            {
              title: "success",
              description: "Useful as confirmations or to finish a multi-step process"
            },
            {
              title: "warning",
              description: "Used as a visual hint for certain actions, like publicly sharing a file or removing a share"

            },
            {
              title: "danger",
              description: "Use for triggering possible destructive and non-reversible actions, like permanently deleting a file"
            }
          ]
        },
        appearances() {
          return [
            "outline",
            "filled",
            "raw",
            "raw-inverse"
          ]
        },
      },
    }
  </script>
  ```
</docs>
