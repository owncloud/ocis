<template>
  <div class="login">
    <div v-if="hasFailed" class="oc-login-card error-msg">
      <div class="oc-login-card-body">
        <h2 class="oc-login-card-title oc-mb-m">
          {{
            $pgettext(
              'The error message title displayed on login page when login fails due to any kind of error.',
              'Something went wrong'
            )
          }}
        </h2>
        <p class="oc-m-rm">
          {{
            $pgettext(
              'The error message displayed on login page when login fails due to any kind of error.',
              "We're having trouble connecting to the login service. If the problem continues, please contact support."
            )
          }}
        </p>
        <oc-button variation="primary" class="oc-mt-l" @click="login">
          {{
            $pgettext(
              'The action to retry the login on login page when login fails due to any kind of error.',
              'Try again'
            )
          }}
        </oc-button>
      </div>
    </div>
    <app-loading-spinner v-else />
  </div>
</template>

<script lang="ts">
import { authService } from '../services/auth'
import { queryItemAsString, useRouteQuery } from '@ownclouders/web-pkg'
import { defineComponent, ref, unref } from 'vue'
import { AppLoadingSpinner } from '@ownclouders/web-pkg'
import { captureException } from '@sentry/vue'

export default defineComponent({
  name: 'LoginPage',
  components: {
    AppLoadingSpinner
  },
  setup() {
    const redirectUrl = useRouteQuery('redirectUrl')

    const hasFailed = ref(false)

    const login = async () => {
      hasFailed.value = false

      try {
        await authService.loginUser(queryItemAsString(unref(redirectUrl)))
      } catch (e) {
        console.error(e)
        captureException(e)

        hasFailed.value = true
      }
    }

    login()

    return {
      hasFailed,
      login
    }
  }
})
</script>

<style lang="scss" scoped>
.login {
  align-content: center;
  min-height: 100dvh;
  width: 100%;
}

.error-msg {
  margin-inline: auto;
  width: min(100%, $width-xlarge-width);
}
</style>
