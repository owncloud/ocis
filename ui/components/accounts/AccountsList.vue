<template>
  <div>
    <oc-table middle divider>
      <oc-table-group>
        <oc-table-row>
          <oc-table-cell shrink type="head" />
          <oc-table-cell type="head" v-text="$gettext('Username')" />
          <oc-table-cell type="head" v-text="$gettext('Display name')" />
          <oc-table-cell type="head" v-text="$gettext('Email')" />
          <oc-table-cell shrink type="head" v-text="$gettext('Uid number')" />
          <oc-table-cell shrink type="head" v-text="$gettext('Gid number')" />
          <oc-table-cell shrink type="head" v-text="$gettext('Enabled')" />
        </oc-table-row>
      </oc-table-group>
      <oc-table-group>
        <oc-table-row v-for="account in accounts" :key="`account-list-row-${account.id}`">
          <oc-table-cell>
            <oc-avatar :userName="account.displayName || account.preferredName" :width="35" />
          </oc-table-cell>
          <oc-table-cell v-text="account.preferredName" />
          <oc-table-cell v-text="account.displayName || '-'" />
          <oc-table-cell v-text="account.mail" />
          <oc-table-cell v-text="account.uidNumber || '-'" />
          <oc-table-cell v-text="account.gidNumber || '-'" />
          <oc-table-cell class="uk-text-center">
            <oc-icon v-if="account.accountEnabled" name="ready" variation="success" :aria-label="$gettext('Account is enabled')" />
            <oc-icon v-else name="deprecated" variation="danger" :aria-label="$gettext('Account is disabled')" />
          </oc-table-cell>
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
