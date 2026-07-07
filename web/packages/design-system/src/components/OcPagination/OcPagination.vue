<template>
  <nav class="oc-pagination" :aria-label="$gettext('Pagination')">
    <ol class="oc-pagination-list">
      <li v-if="isPrevPageAvailable" class="oc-pagination-list-item">
        <router-link
          class="oc-pagination-list-item-prev"
          :aria-label="$gettext('Go to the previous page')"
          :to="previousPageLink"
        >
          <oc-icon name="arrow-drop-left" fill-type="line" />
        </router-link>
      </li>
      <li v-for="(page, index) in displayedPages" :key="index" class="oc-pagination-list-item">
        <component :is="pageComponent(page)" :class="pageClass(page)" v-bind="bindPageProps(page)">
          {{ page }}
        </component>
      </li>
      <li v-if="isNextPageAvailable" class="oc-pagination-list-item">
        <router-link
          class="oc-pagination-list-item-next"
          :aria-label="$gettext('Go to the next page')"
          :to="nextPageLink"
        >
          <oc-icon name="arrow-drop-right" fill-type="line" />
        </router-link>
      </li>
    </ol>
  </nav>
</template>

<script lang="ts" setup>
import { computed, unref } from 'vue'
import { RouteLocation } from 'vue-router'
import { useGettext } from 'vue3-gettext'
import OcIcon from '../OcIcon/OcIcon.vue'

/**
 * @component OcPagination
 *
 * Rendering a pagination control with support for dynamic page ranges,
 * ellipsis, and navigation links. It integrates with Vue Router for seamless navigation.
 *
 * @props
 * @prop {number} pages - Total number of pages.
 * @prop {number} currentPage - The currently active page.
 * @prop {number} [maxDisplayed] - Maximum number of pages to display at once. Defaults to showing all pages.
 * @prop {RouteLocation} currentRoute - The current Vue Router route object.
 *
 * @example
 *   <OcPagination
 *     :pages="10"
 *     :currentPage="5"
 *     :maxDisplayed="5"
 *     :currentRoute="$route"
 *   />
 *
 */

defineOptions({
  name: 'OcPagination',
  status: 'ready',
  release: '7.2.0'
})
type Page = string | number

interface Props {
  pages: number
  currentPage: number
  maxDisplayed?: number
  currentRoute: RouteLocation
}

const { pages, currentPage, maxDisplayed = null, currentRoute } = defineProps<Props>()
const { $gettext } = useGettext()

function pageLabel(page: Page) {
  return $gettext('Go to page %{ page }', { page: page.toString() })
}

function isCurrentPage(page: Page) {
  return unref(computedCurrentPage) === page
}

function pageComponent(page: Page) {
  return page === '...' || isCurrentPage(page) ? 'span' : 'router-link'
}

function bindPageProps(page: Page) {
  if (page === '...') {
    return
  }

  if (isCurrentPage(page)) {
    return {
      'aria-current': 'page'
    }
  }

  const link = bindPageLink(page)

  return {
    'aria-label': pageLabel(page),
    to: link
  }
}

function pageClass(page: Page) {
  const classes = ['oc-pagination-list-item-page']

  if (isCurrentPage(page)) {
    classes.push('oc-pagination-list-item-current')
  } else if (page === '...') {
    classes.push('oc-pagination-list-item-ellipsis')
  } else {
    classes.push('oc-pagination-list-item-link')
  }

  return classes
}

function bindPageLink(page: Page) {
  return {
    name: currentRoute.name,
    query: { ...currentRoute.query, page },
    params: currentRoute.params
  }
}
const displayedPages = computed(() => {
  let pagination: Array<Page> = []

  for (let i = 0; i < pages; i++) {
    pagination.push(i + 1)
  }

  if (maxDisplayed && maxDisplayed + 1 < pages) {
    const currentPageIndex = unref(computedCurrentPage) - 1
    const indentation = Math.floor(maxDisplayed / 2)

    pagination = pagination.slice(
      Math.max(0, currentPageIndex - indentation),
      currentPageIndex + indentation + 1
    )

    if (unref(computedCurrentPage) > 2) {
      Number(pagination[0]) > 2 ? pagination.unshift(1, '...') : pagination.unshift(1)
    }

    if (unref(computedCurrentPage) < pages - 1) {
      Number(pagination[pagination.length - 1]) < pages - 1
        ? pagination.push('...', pages)
        : pagination.push(pages)
    }

    return pagination
  }

  return pagination
})

const isPrevPageAvailable = computed(() => {
  return unref(computedCurrentPage) > 1
})

const isNextPageAvailable = computed(() => {
  return unref(computedCurrentPage) < pages
})

const previousPageLink = computed(() => {
  return bindPageLink(unref(computedCurrentPage) - 1)
})

const nextPageLink = computed(() => {
  return bindPageLink(unref(computedCurrentPage) + 1)
})

const computedCurrentPage = computed(() => {
  return Math.max(1, Math.min(currentPage, pages))
})
</script>

<style lang="scss">
.oc-pagination {
  &-list {
    align-items: center;
    display: flex;
    flex-wrap: wrap;
    gap: var(--oc-space-small);
    list-style: none;
    margin: 0;
    padding: 0;

    &-item {
      &-page {
        border-radius: 4px;
        color: var(--oc-color-text-default);
        padding: var(--oc-space-xsmall) var(--oc-space-small);
        transition: background-color $transition-duration-short ease-in-out;

        &:not(span):hover {
          background-color: var(--oc-color-swatch-passive-default);
          color: var(--oc-color-text-inverse);
          text-decoration: none;
        }
      }

      &-current {
        background-color: var(--oc-color-swatch-passive-default);
        color: var(--oc-color-text-inverse);
        font-weight: bold;
      }

      &-prev,
      &-next {
        display: flex;

        > .oc-icon > svg {
          fill: var(--oc-color-text-default);
        }
      }

      &-prev {
        margin-right: var(--oc-space-small);
      }

      &-next {
        margin-left: var(--oc-space-small);
      }
    }
  }
}
</style>
