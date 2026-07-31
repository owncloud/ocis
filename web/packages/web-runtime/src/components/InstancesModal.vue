<template>
  <oc-list class="instances-list">
    <li
      v-for="instance in instances"
      :key="instance.url"
      :class="['instance', instance.active && 'active']"
    >
      <div class="instance-info">
        <a :href="instance.url">{{ instance.url }}</a>
        <div>
          <oc-tag v-if="instance.primary" size="small">{{
            $pgettext(
              'The badge label for the primary instance in the instances modal available when multiple instances are enabled in oCIS',
              'Primary'
            )
          }}</oc-tag>
          <oc-tag v-if="instance.active" size="small">{{
            $pgettext(
              'The badge label for the active instance in the instances modal available when multiple instances are enabled in oCIS',
              'Active'
            )
          }}</oc-tag>
        </div>
      </div>
      <oc-button v-if="!instance.active" type="a" variation="secondary" :href="instance.url">
        {{
          $pgettext(
            'The switch action label of an instance in the instances modal available when multiple instances are enabled in oCIS',
            'Active'
          )
        }}
      </oc-button>
    </li>
  </oc-list>
</template>

<script lang="ts" setup>
import { Modal } from '@ownclouders/web-pkg'
import { useInstances } from '../composables/instances'

const props = defineProps<{
  modal: Modal
}>()

const { instances } = useInstances()
</script>

<style lang="scss" scoped>
.instances-list {
  display: grid;
  gap: var(--oc-space-small);
}

.instance {
  align-items: center;
  border: 1px solid var(--oc-color-border);
  border-radius: 5px;
  display: flex;
  gap: var(--oc-space-small);
  justify-content: space-between;
  padding: var(--oc-space-small);

  &.active {
    border-color: var(--oc-color-swatch-primary-default);
  }
}

.instance-info {
  display: grid;
  gap: var(--oc-space-small);
}
</style>
