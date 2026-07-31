<template>
  <div :class="classes">
    <slot />
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue'

/**
 * OcGrid Component
 *
 * The Grid system allows you to arrange block elements in columns. It supports
 * configurable gutters and optional flexbox alignment for centering elements.
 *
 * @component
 * @name OcGrid
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @property {('small'|'medium'|'large'|'collapse')} [gutter='collapse'] - The Grid component comes with a default gutter that is decreased automatically from a certain breakpoint. To apply a different gutter, use one of the following: 'small', 'medium', 'large', or 'collapse'.
 * @property {boolean} [flex=false] - Center elements along the cross axis using flexbox.
 *
 * @slots
 * @slot default - Slot for grid content.
 *
 * @example
 * <template>
 *   <oc-grid gutter="large" :flex="tue" />
 * </template>
 */

interface Props {
  gutter?: 'small' | 'medium' | 'large' | 'collapse'
  flex?: boolean
}

defineOptions({
  name: 'OcGrid',
  status: 'ready',
  release: '1.0.0'
})
const { gutter = 'collapse', flex = false } = defineProps<Props>()

const classes = computed(() => {
  const c = []

  c.push('oc-grid-' + gutter)

  if (flex) {
    c.push('oc-flex')
    c.push('oc-flex-middle')
  }

  return c
})
</script>

<style lang="scss">
// Ported from UiKit until we change UI frameworks

// Name:            Grid
// Description:     Component to create responsive, fluid and nestable grids
//
// Component:       `oc-grid`
//
// Modifiers:       `oc-grid-small`
//                  `oc-grid-medium`
//                  `oc-grid-large`
//                  `oc-grid-collapse`
//                  `oc-grid-divider`
//                  `oc-grid-match`
//                  `oc-grid-stack`
//                  `oc-grid-margin`
//                  `oc-grid-margin-small`
//                  `oc-grid-margin-medium`
//                  `oc-grid-margin-large`
//                  `oc-grid-margin-collapse`
//
// Sub-modifier:    `oc-grid-item-match`
//
// States:          `oc-first-column`
//
// ========================================================================

// Variables
// ========================================================================

$grid-gutter-horizontal: 30px !default;
$grid-gutter-vertical: $grid-gutter-horizontal !default;
$grid-gutter-horizontal-l: 40px !default;
$grid-gutter-vertical-l: $grid-gutter-horizontal-l !default;

$grid-small-gutter-horizontal: 15px !default;
$grid-small-gutter-vertical: $grid-small-gutter-horizontal !default;

$grid-medium-gutter-horizontal: 30px !default;
$grid-medium-gutter-vertical: $grid-medium-gutter-horizontal !default;

$grid-large-gutter-horizontal: 40px !default;
$grid-large-gutter-vertical: $grid-large-gutter-horizontal !default;
$grid-large-gutter-horizontal-l: 70px !default;
$grid-large-gutter-vertical-l: $grid-large-gutter-horizontal-l !default;

$grid-divider-border-width: 1px !default;
$grid-divider-border: #e5e5e5 !default;

/* ========================================================================
   Component: Grid
 ========================================================================== */

/*
 * 1. Allow cells to wrap into the next line
 * 2. Reset list
 */

.oc-grid {
  display: flex;
  /* 1 */
  flex-wrap: wrap;
  /* 2 */
  margin: 0;
  padding: 0;
  list-style: none;
}

/*
 * Grid cell
 * Note: Space is allocated solely based on content dimensions, but shrinks: 0 1 auto
 * Reset margin for e.g. paragraphs
 */

.oc-grid > * {
  margin: 0;
}

/*
 * Remove margin from the last-child
 */

.oc-grid > * > :last-child {
  margin-bottom: 0;
}

/* Gutter
 ========================================================================== */

/*
 * Default
 */

/* Horizontal */
.oc-grid {
  margin-left: (-$grid-gutter-horizontal);
}
.oc-grid > * {
  padding-left: $grid-gutter-horizontal;
}

/* Vertical */
.oc-grid + .oc-grid,
.oc-grid > .oc-grid-margin,
* + .oc-grid-margin {
  margin-top: $grid-gutter-vertical;
}

/* Desktop and bigger */
@media (min-width: $oc-breakpoint-large-default) {
  /* Horizontal */
  .oc-grid {
    margin-left: (-$grid-gutter-horizontal-l);
  }
  .oc-grid > * {
    padding-left: $grid-gutter-horizontal-l;
  }

  /* Vertical */
  .oc-grid + .oc-grid,
  .oc-grid > .oc-grid-margin,
  * + .oc-grid-margin {
    margin-top: $grid-gutter-vertical-l;
  }
}

/*
 * Small
 */

/* Horizontal */
.oc-grid-small,
.oc-grid-column-small {
  margin-left: (-$grid-small-gutter-horizontal);
}
.oc-grid-small > *,
.oc-grid-column-small > * {
  padding-left: $grid-small-gutter-horizontal;
}

/* Vertical */
.oc-grid + .oc-grid-small,
.oc-grid + .oc-grid-row-small,
.oc-grid-small > .oc-grid-margin,
.oc-grid-row-small > .oc-grid-margin,
* + .oc-grid-margin-small {
  margin-top: $grid-small-gutter-vertical;
}

/*
 * Medium
 */

/* Horizontal */
.oc-grid-medium,
.oc-grid-column-medium {
  margin-left: (-$grid-medium-gutter-horizontal);
}
.oc-grid-medium > *,
.oc-grid-column-medium > * {
  padding-left: $grid-medium-gutter-horizontal;
}

/* Vertical */
.oc-grid + .oc-grid-medium,
.oc-grid + .oc-grid-row-medium,
.oc-grid-medium > .oc-grid-margin,
.oc-grid-row-medium > .oc-grid-margin,
* + .oc-grid-margin-medium {
  margin-top: $grid-medium-gutter-vertical;
}

/*
 * Large
 */

/* Horizontal */
.oc-grid-large,
.oc-grid-column-large {
  margin-left: (-$grid-large-gutter-horizontal);
}
.oc-grid-large > *,
.oc-grid-column-large > * {
  padding-left: $grid-large-gutter-horizontal;
}

/* Vertical */
.oc-grid + .oc-grid-large,
.oc-grid + .oc-grid-row-large,
.oc-grid-large > .oc-grid-margin,
.oc-grid-row-large > .oc-grid-margin,
* + .oc-grid-margin-large {
  margin-top: $grid-large-gutter-vertical;
}

/* Desktop and bigger */
@media (min-width: $oc-breakpoint-large-default) {
  /* Horizontal */
  .oc-grid-large,
  .oc-grid-column-large {
    margin-left: (-$grid-large-gutter-horizontal-l);
  }
  .oc-grid-large > *,
  .oc-grid-column-large > * {
    padding-left: $grid-large-gutter-horizontal-l;
  }

  /* Vertical */
  .oc-grid + .oc-grid-large,
  .oc-grid + .oc-grid-row-large,
  .oc-grid-large > .oc-grid-margin,
  .oc-grid-row-large > .oc-grid-margin,
  * + .oc-grid-margin-large {
    margin-top: $grid-large-gutter-vertical-l;
  }
}

/*
 * Collapse
 */

/* Horizontal */
.oc-grid-collapse,
.oc-grid-column-collapse {
  margin-left: 0;
}
.oc-grid-collapse > *,
.oc-grid-column-collapse > * {
  padding-left: 0;
}

/* Vertical */
.oc-grid + .oc-grid-collapse,
.oc-grid + .oc-grid-row-collapse,
.oc-grid-collapse > .oc-grid-margin,
.oc-grid-row-collapse > .oc-grid-margin {
  margin-top: 0;
}

/* Divider
 ========================================================================== */

.oc-grid-divider > * {
  position: relative;
}

.oc-grid-divider > :not(.oc-first-column)::before {
  content: '';
  position: absolute;
  top: 0;
  bottom: 0;
  border-left: $grid-divider-border-width solid $grid-divider-border;
}

/* Vertical */
.oc-grid-divider.oc-grid-stack > .oc-grid-margin::before {
  content: '';
  position: absolute;
  left: 0;
  right: 0;
  border-top: $grid-divider-border-width solid $grid-divider-border;
}

/*
 * Default
 */

/* Horizontal */
.oc-grid-divider {
  margin-left: -($grid-gutter-horizontal * 2);
}
.oc-grid-divider > * {
  padding-left: ($grid-gutter-horizontal * 2);
}

.oc-grid-divider > :not(.oc-first-column)::before {
  left: $grid-gutter-horizontal;
}

/* Vertical */
.oc-grid-divider.oc-grid-stack > .oc-grid-margin {
  margin-top: ($grid-gutter-vertical * 2);
}

.oc-grid-divider.oc-grid-stack > .oc-grid-margin::before {
  top: (-$grid-gutter-vertical);
  left: ($grid-gutter-horizontal * 2);
}

/* Desktop and bigger */
@media (min-width: $oc-breakpoint-large-default) {
  /* Horizontal */
  .oc-grid-divider {
    margin-left: -($grid-gutter-horizontal-l * 2);
  }
  .oc-grid-divider > * {
    padding-left: ($grid-gutter-horizontal-l * 2);
  }

  .oc-grid-divider > :not(.oc-first-column)::before {
    left: $grid-gutter-horizontal-l;
  }

  /* Vertical */
  .oc-grid-divider.oc-grid-stack > .oc-grid-margin {
    margin-top: ($grid-gutter-vertical-l * 2);
  }

  .oc-grid-divider.oc-grid-stack > .oc-grid-margin::before {
    top: (-$grid-gutter-vertical-l);
    left: ($grid-gutter-horizontal-l * 2);
  }
}

/*
 * Small
 */

/* Horizontal */
.oc-grid-divider.oc-grid-small,
.oc-grid-divider.oc-grid-column-small {
  margin-left: -($grid-small-gutter-horizontal * 2);
}
.oc-grid-divider.oc-grid-small > *,
.oc-grid-divider.oc-grid-column-small > * {
  padding-left: ($grid-small-gutter-horizontal * 2);
}

.oc-grid-divider.oc-grid-small > :not(.oc-first-column)::before,
.oc-grid-divider.oc-grid-column-small > :not(.oc-first-column)::before {
  left: $grid-small-gutter-horizontal;
}

/* Vertical */
.oc-grid-divider.oc-grid-small.oc-grid-stack > .oc-grid-margin,
.oc-grid-divider.oc-grid-row-small.oc-grid-stack > .oc-grid-margin {
  margin-top: ($grid-small-gutter-vertical * 2);
}

.oc-grid-divider.oc-grid-small.oc-grid-stack > .oc-grid-margin::before {
  top: (-$grid-small-gutter-vertical);
  left: ($grid-small-gutter-horizontal * 2);
}

.oc-grid-divider.oc-grid-row-small.oc-grid-stack > .oc-grid-margin::before {
  top: (-$grid-small-gutter-vertical);
}
.oc-grid-divider.oc-grid-column-small.oc-grid-stack > .oc-grid-margin::before {
  left: ($grid-small-gutter-horizontal * 2);
}

/*
 * Medium
 */

/* Horizontal */
.oc-grid-divider.oc-grid-medium,
.oc-grid-divider.oc-grid-column-medium {
  margin-left: -($grid-medium-gutter-horizontal * 2);
}
.oc-grid-divider.oc-grid-medium > *,
.oc-grid-divider.oc-grid-column-medium > * {
  padding-left: ($grid-medium-gutter-horizontal * 2);
}

.oc-grid-divider.oc-grid-medium > :not(.oc-first-column)::before,
.oc-grid-divider.oc-grid-column-medium > :not(.oc-first-column)::before {
  left: $grid-medium-gutter-horizontal;
}

/* Vertical */
.oc-grid-divider.oc-grid-medium.oc-grid-stack > .oc-grid-margin,
.oc-grid-divider.oc-grid-row-medium.oc-grid-stack > .oc-grid-margin {
  margin-top: ($grid-medium-gutter-vertical * 2);
}

.oc-grid-divider.oc-grid-medium.oc-grid-stack > .oc-grid-margin::before {
  top: (-$grid-medium-gutter-vertical);
  left: ($grid-medium-gutter-horizontal * 2);
}

.oc-grid-divider.oc-grid-row-medium.oc-grid-stack > .oc-grid-margin::before {
  top: (-$grid-medium-gutter-vertical);
}
.oc-grid-divider.oc-grid-column-medium.oc-grid-stack > .oc-grid-margin::before {
  left: ($grid-medium-gutter-horizontal * 2);
}

/*
 * Large
 */

/* Horizontal */
.oc-grid-divider.oc-grid-large,
.oc-grid-divider.oc-grid-column-large {
  margin-left: -($grid-large-gutter-horizontal * 2);
}
.oc-grid-divider.oc-grid-large > *,
.oc-grid-divider.oc-grid-column-large > * {
  padding-left: ($grid-large-gutter-horizontal * 2);
}

.oc-grid-divider.oc-grid-large > :not(.oc-first-column)::before,
.oc-grid-divider.oc-grid-column-large > :not(.oc-first-column)::before {
  left: $grid-large-gutter-horizontal;
}

/* Vertical */
.oc-grid-divider.oc-grid-large.oc-grid-stack > .oc-grid-margin,
.oc-grid-divider.oc-grid-row-large.oc-grid-stack > .oc-grid-margin {
  margin-top: ($grid-large-gutter-vertical * 2);
}

.oc-grid-divider.oc-grid-large.oc-grid-stack > .oc-grid-margin::before {
  top: (-$grid-large-gutter-vertical);
  left: ($grid-large-gutter-horizontal * 2);
}

.oc-grid-divider.oc-grid-row-large.oc-grid-stack > .oc-grid-margin::before {
  top: (-$grid-large-gutter-vertical);
}
.oc-grid-divider.oc-grid-column-large.oc-grid-stack > .oc-grid-margin::before {
  left: ($grid-large-gutter-horizontal * 2);
}

/* Desktop and bigger */
@media (min-width: $oc-breakpoint-large-default) {
  /* Horizontal */
  .oc-grid-divider.oc-grid-large,
  .oc-grid-divider.oc-grid-column-large {
    margin-left: -($grid-large-gutter-horizontal-l * 2);
  }
  .oc-grid-divider.oc-grid-large > *,
  .oc-grid-divider.oc-grid-column-large > * {
    padding-left: ($grid-large-gutter-horizontal-l * 2);
  }

  .oc-grid-divider.oc-grid-large > :not(.oc-first-column)::before,
  .oc-grid-divider.oc-grid-column-large > :not(.oc-first-column)::before {
    left: $grid-large-gutter-horizontal-l;
  }

  /* Vertical */
  .oc-grid-divider.oc-grid-large.oc-grid-stack > .oc-grid-margin,
  .oc-grid-divider.oc-grid-row-large.oc-grid-stack > .oc-grid-margin {
    margin-top: ($grid-large-gutter-vertical-l * 2);
  }

  .oc-grid-divider.oc-grid-large.oc-grid-stack > .oc-grid-margin::before {
    top: (-$grid-large-gutter-vertical-l);
    left: ($grid-large-gutter-horizontal-l * 2);
  }

  .oc-grid-divider.oc-grid-row-large.oc-grid-stack > .oc-grid-margin::before {
    top: (-$grid-large-gutter-vertical-l);
  }
  .oc-grid-divider.oc-grid-column-large.oc-grid-stack > .oc-grid-margin::before {
    left: ($grid-large-gutter-horizontal-l * 2);
  }
}

/* Match child of a grid cell
 ========================================================================== */

/*
 * Behave like a block element
 * 1. Wrap into the next line
 * 2. Take the full width, at least 100%. Only if no class from the Width component is set.
 * 3. Expand width even if larger than 100%, e.g. because of negative margin (Needed for nested grids)
 */

.oc-grid-match > *,
.oc-grid-item-match {
  display: flex;
  /* 1 */
  flex-wrap: wrap;
}

.oc-grid-match > * > :not([class*='oc-width']),
.oc-grid-item-match > :not([class*='oc-width']) {
  /* 2 */
  box-sizing: border-box;
  width: 100%;
  /* 3 */
  flex: auto;
}
</style>
