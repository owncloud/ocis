<template>
  <oc-button
    id="new-space-menu-btn"
    v-oc-tooltip="$gettext('Request new project')"
    type="a"
    href="https://cern.service-now.com/service-portal?id=sc_cat_item&name=request-storage-space&se=CERNBox-Service"
    target="_blank"
    appearance="filled"
    variation="primary"
  >
    <oc-icon name="add" />
    <span v-translate>New project</span>
  </oc-button>
</template>
