<template>
  <div class="oc-text-center">
    <p data-testid="files-list-footer-info" class="oc-text-muted" v-text="resourceContentsText"></p>
  </div>
</template>

<script lang="ts" setup>
import { useResourceContents } from '@ownclouders/web-pkg'

const { resourceContentsText } = useResourceContents()
</script>
