<template>
  <div id="oc-files-context-menu">
    <oc-list
      v-for="(section, sectionIndex) in menuSections"
      :id="`oc-files-context-actions-${section.name}`"
      :key="`section-${section.name}-list`"
      class="oc-files-context-actions"
      :class="getSectionClasses(sectionIndex)"
    >
      <action-menu-item
        v-for="(action, actionIndex) in section.items"
        :key="`section-${section.name}-action-${actionIndex}`"
        :action="action"
        :appearance="appearance"
        :variation="variation"
        :action-options="actionOptions"
        :has-limited-screen-space="true"
        class="context-menu oc-files-context-action oc-px-s oc-rounded oc-menu-item-hover"
      />
    </oc-list>
  </div>
</template>

<script lang="ts" setup>
import ActionMenuItem from './ActionMenuItem.vue'
import { Action, ActionOptions } from '../../composables'

export type MenuSection = {
  name: string
  items: Action[]
}
interface Props {
  menuSections: MenuSection[]
  appearance?: string
  variation?: string
  actionOptions: ActionOptions
}
const { menuSections, appearance = 'raw', variation = 'passive' } = defineProps<Props>()
function getSectionClasses(index: number) {
  const classes: string[] = []
  if (!menuSections.length) {
    return classes
  }
  if (index < menuSections.length - 1) {
    classes.push('oc-pb-s')
  }
  if (index > 0) {
    classes.push('oc-pt-s')
  }
  if (index < menuSections.length - 1) {
    classes.push('oc-files-context-actions-border')
  }
  return classes
}
</script>

<style lang="scss">
.oc-files-context-actions {
  text-align: left;
  white-space: normal;

  > li {
    padding-left: 0 !important;
    padding-right: 0 !important;
    a,
    button,
    span {
      display: inline-flex;
      font-weight: normal !important;
      justify-content: flex-start;
      vertical-align: top;
      width: 100%;
      text-align: left;
    }
  }

  &-border {
    border-bottom: 1px solid var(--oc-color-border);
  }
}
</style>
