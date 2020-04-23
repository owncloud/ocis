<template>
  <div>
    <translate :translate-params="{ type: setting.type }">Settings type not implemented: %{type}</translate>
  </div>
</template>

<script>
export default {
  name: 'SettingUnknown',
  props: {
    setting: {
      type: Object,
      required: true
    }
  }
}
</script>
