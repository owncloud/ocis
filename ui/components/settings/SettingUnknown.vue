<template>
  <div>
    <translate :translate-params="{ type: setting.type }">Settings type not implemented: %{type}</translate>
  </div>
</template>

<script>
export default {
  name: 'SettingUnknown',
  props: {
    bundle: {
      type: Object,
      required: true
    },
    setting: {
      type: Object,
      required: true
    },
    persistedValue: {
      type: Object,
      required: false
    }
  }
}
</script>
