<template>
  <div class="uk-flex uk-flex-column uk-flex-middle uk-margin-xlarge-top">
    <h1>
      Settings service...
    </h1>
  </div>
</template>

<script>
export default {
  name: 'SettingsApp',
  async mounted () {
    const settingsBundles = await this.$store.dispatch('Settings/fetchSettingsBundles')
    console.log(settingsBundles)
    console.log("test")
  }
}
</script>
