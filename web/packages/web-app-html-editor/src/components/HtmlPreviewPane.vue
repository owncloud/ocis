<template>
  <!-- Sandboxed preview. `srcdoc` runs at an opaque origin (no allow-same-origin),
       so it cannot read the shell's cookies, storage or token. The sandbox is
       intentionally minimal: only `allow-scripts` (to render self-contained pages
       and run their inline JS); `allow-forms` and `allow-popups` are deliberately
       omitted so the preview cannot submit forms or open windows for phishing or
       beaconing. The content is additionally wrapped with a strict CSP upstream.
       See DECISIONS.md D3. -->
  <iframe
    class="html-preview-pane"
    :srcdoc="content"
    sandbox="allow-scripts"
    referrerpolicy="no-referrer"
    :title="$gettext('HTML preview')"
  />
</template>

<script lang="ts" setup>
interface Props {
  content: string
}
defineProps<Props>()
</script>

<style scoped>
.html-preview-pane {
  width: 100%;
  height: 100%;
  border: none;
  background-color: var(--oc-color-background-default);
}
</style>
