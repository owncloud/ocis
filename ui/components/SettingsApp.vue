<template>
  <div class="uk-flex uk-flex-column uk-flex-middle uk-margin-xlarge-top">
    <h1>
      Settings service...
    </h1>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
export default {
  name: 'SettingsApp',
  computed: {
    ...mapGetters({
      settingsBundles: 'Settings/settingsBundles'
    })
  },
  methods: {
    ...mapActions('Settings', [
      'fetchSettingsBundles'
    ])
  },
  mounted () {
    this.fetchSettingsBundles().then(() => {
      console.log(this.settingsBundles);
    });
  }
}
</script>
