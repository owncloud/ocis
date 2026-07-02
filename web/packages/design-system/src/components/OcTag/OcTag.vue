<template>
  <component :is="type" :class="ocTagClass" :to="to" @click="ocTagClick">
    <!-- @slot Content of the tag -->
    <slot />
  </component>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import { RouteLocationRaw } from 'vue-router'
import { getSizeClass } from '../../helpers'

/**
 * @component OcTag
 * @description A versatile tag component for displaying various types of information as labels, links or buttons
 *
 * @prop {('span'|'button'|'router-link'|'a')} [type='span'] - HTML element or component to be rendered
 * @prop {string|RouteLocationRaw} [to=null] - Target location for router-link or anchor href
 * @prop {('small'|'medium'|'large')} [size='medium'] - Size of the tag
 * @prop {boolean} [rounded=false] - Whether the tag should have fully rounded corners
 *
 * @emits {MouseEvent} click - Emitted when the tag is clicked
 *
 * @slot default - Content of the tag. Can include text, icons, or other elements
 *
 * @example
 * <!-- Basic tag -->
 * <oc-tag>Basic tag</oc-tag>
 *
 * @example
 * <!-- Tag with icon -->
 * <oc-tag>
 *   <oc-icon name="links" />
 *   Shared via link
 * </oc-tag>
 *
 */

interface Props {
  type?: 'span' | 'button' | 'router-link' | 'a'
  to?: string | RouteLocationRaw
  size?: 'small' | 'medium' | 'large'
  rounded?: boolean
}

interface Emits {
  (e: 'click', event: MouseEvent): void
}

defineOptions({
  name: 'OcTag',
  status: 'ready',
  release: '2.0.0'
})

const { type = 'span', to = null, size = 'medium', rounded = false } = defineProps<Props>()

const emit = defineEmits<Emits>()

function ocTagClick(event: MouseEvent) {
  emit('click', event)
}

const ocTagClass = computed(() => {
  const classes = ['oc-tag', `oc-tag-${getSizeClass(size)}`]

  type === 'router-link' || type === 'a'
    ? classes.push('oc-tag-link')
    : classes.push(`oc-tag-${type}`)

  if (rounded) {
    classes.push('oc-tag-rounded')
  }

  return classes
})
</script>

<style lang="scss">
.oc-tag {
  align-items: center;
  background-color: var(--oc-color-background-default);
  border: 1px solid var(--oc-color-text-muted);
  border-radius: 7px;
  box-sizing: border-box;
  color: var(--oc-color-text-muted);
  display: inline-flex;
  gap: var(--oc-space-xsmall);
  text-decoration: none;

  &-s {
    font-size: 0.75rem;
    padding: var(--oc-space-xsmall);
  }

  &-m {
    font-size: 0.875rem;
    min-height: 2.125rem;
    padding: var(--oc-space-xsmall) var(--oc-space-small);
  }

  &-l {
    font-size: 1.5rem;
    min-height: 2.75rem;
    padding: var(--oc-space-small) var(--oc-space-medium);
  }

  &-rounded {
    border-radius: 99px;
    padding-left: var(--oc-space-small);
    padding-right: var(--oc-space-small);
  }

  .oc-icon > svg {
    fill: var(--oc-color-text-muted);
  }

  &-link,
  &-button {
    transition: color $transition-duration-short ease-in-out;

    .oc-icon > svg {
      transition: fill $transition-duration-short ease-in-out;
    }

    &:hover,
    &:focus {
      color: var(--oc-color-swatch-primary-hover);
      cursor: pointer;
      text-decoration: none;

      .oc-icon > svg {
        fill: var(--oc-color-swatch-primary-hover);
      }
    }
  }
}
</style>

<docs>
Component to display various information.
```js
<oc-tag>
  <oc-icon name="links" />
  Shared via link
</oc-tag>
```
## Different sizes of the tag component

```js
<div>
<oc-tag size="small">
  <oc-icon name="links" size="small" />
  Small tag
</oc-tag>
<oc-tag size="medium">
  <oc-icon name="links" size="medium" />
  Medium tag
</oc-tag>
<oc-tag size="large">
  <oc-icon name="links" size="large" />
  Large tag
</oc-tag>
</div>
```
## Different types of the tag component
The tag component can be rendered as a different element if desired. You can specify such element via property `type`.

```js
<oc-grid gutter="small" flex="true">
    <oc-tag class="oc-mr-s">
      <oc-icon name="group" />
      Shared with other people
    </oc-tag>
    <oc-tag class="oc-mr-s" type="a">
      <oc-icon name="links" />
      Shared via link
    </oc-tag>
    <oc-tag class="oc-mr-s" type="button">Expires in 2 days</oc-tag>
</oc-grid>
```
</docs>
