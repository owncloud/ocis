<template>
  <div class="account-table">
    <slot name="header" :title="title">
      <h2 class="account-table-title" v-text="title" />
    </slot>
    <oc-table-simple>
      <oc-thead :class="{ 'oc-invisible-sr': !showHead }">
        <oc-tr>
          <template v-for="field in fields" :key="typeof field === 'string' ? field : field.label">
            <oc-th v-if="typeof field === 'string'">{{ field }}</oc-th>
            <oc-th
              v-else
              :align-h="field.alignH || 'left'"
              :class="{ 'oc-invisible-sr': field.hidden }"
            >
              {{ field.label }}
            </oc-th>
          </template>
        </oc-tr>
      </oc-thead>
      <oc-tbody>
        <slot />
      </oc-tbody>
    </oc-table-simple>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

type AccountTableCell = {
  label: string
  alignH?: string
  hidden?: boolean
}

export default defineComponent({
  name: 'AccountTable',
  props: {
    title: {
      type: String,
      required: true
    },
    fields: {
      type: Array<string | AccountTableCell>,
      required: true
    },
    showHead: { type: Boolean, required: false, default: false }
  }
})
</script>

<style lang="scss">
@media (max-width: $oc-breakpoint-small-max) {
  .account-table {
    tr {
      display: block;
      padding-bottom: var(--oc-space-xsmall);
      height: 100% !important;
    }

    td {
      display: block !important;
      width: 100% !important;
      padding-top: var(--oc-space-small);
      padding-bottom: var(--oc-space-small);
    }

    h2 {
      color: var(--oc-color-text-muted);
      font-size: var(--oc-font-size-large);
      font-weight: var(--oc-font-weight-default);
    }

    .oc-select {
      width: 100%;
    }
  }
}

.account-table {
  .oc-select {
    min-width: 200px;
  }

  tr {
    border-top: 0;
    border-bottom: 1px solid var(--oc-color-border);
    height: var(--oc-size-height-table-row);
  }

  td:first-of-type {
    width: 20%;
  }

  td:not(:first-of-type) {
    color: var(--oc-color-text-muted);

    p {
      color: var(--oc-color-text-muted);
    }
  }

  @media (min-width: $oc-breakpoint-medium-default) {
    td > .checkbox-cell-wrapper {
      display: flex;
      justify-content: end;
      align-items: center;
      min-height: var(--oc-size-height-table-row);
    }
  }
}
</style>
