<template>
  <div>
    <oc-table middle divider>
      <oc-table-group>
        <oc-table-row>
          <oc-table-cell shrink type="head" />
          <oc-table-cell type="head" v-text="$gettext('Name')" />
          <oc-table-cell type="head" v-text="$gettext('Email')" />
        </oc-table-row>
      </oc-table-group>
      <oc-table-group>
        <oc-table-row v-for="account in accounts" :key="`account-list-row-${account.id}`">
          <oc-table-cell>
            <oc-avatar :userName="account.preferredName" />
          </oc-table-cell>
          <oc-table-cell v-text="account.preferredName" />
          <oc-table-cell v-text="account.mail" />
        </oc-table-row>
      </oc-table-group>
    </oc-table>
  </div>
</template>

<script>
export default {
  name: 'AccountsList',
  props: {
    accounts: {
      type: Array,
      required: true
    }
  }
}
</script>
