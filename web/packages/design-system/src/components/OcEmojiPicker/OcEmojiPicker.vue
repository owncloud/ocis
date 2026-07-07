<template>
  <div v-if="isLoading" class="oc-flex oc-flex-center">
    <oc-spinner size="large" />
  </div>
  <div v-else ref="emojiPickerRef"></div>
</template>

<script lang="ts" setup>
import { nextTick, ref, unref, watch } from 'vue'
import { useGettext } from 'vue3-gettext'

/**
 * OcEmojiPicker component
 *
 * @component
 * @name ocemojipicker
 * @status ready
 * @release 1.0.0
 *
 * @props
 * @property {('light' | 'dark')} [theme='light'] - the theme of the emoji picker, either 'light' or 'dark'.
 *
 * @emits
 * @event emojiselect - triggered when an emoji is selected.
 * @param {string} emoji - the selected emoji's native representation.
 *
 * @event clickoutside - triggered when a click occurs outside the emoji picker.
 *
 *
 * @example
 * ```vue
 * <template>
 *   <oc-emoji-picker
 *     theme="dark"
 *     @emojiselect="onemojiselect"
 *     @clickoutside="onclickoutside"
 *   />
 * </template>
 * ```
 */

interface Props {
  theme?: 'light' | 'dark'
}
interface Emits {
  (e: 'emojiSelect', emoji: string): void
  (e: 'clickOutside'): void
}
defineOptions({
  name: 'OcEmojiPicker',
  status: 'ready',
  release: '1.0.0'
})
const { theme = 'light' } = defineProps<Props>()
const emit = defineEmits<Emits>()

const language = useGettext()
const { $gettext } = language
const emojiPickerRef = ref<HTMLElement>()

const isLoading = ref(true)

watch(
  [() => theme, () => language.current],
  async () => {
    isLoading.value = true

    await nextTick()
    const i18n = {
      search: $gettext('Search'),
      search_no_results_1: $gettext('Oh no!'),
      search_no_results_2: $gettext('That emoji couldn’t be found'),
      pick: $gettext('Pick an emoji…'),
      add_custom: $gettext('Add custom emoji'),
      categories: {
        activity: $gettext('Activity'),
        custom: $gettext('Custom'),
        flags: $gettext('Flags'),
        foods: $gettext('Food & Drink'),
        frequent: $gettext('Frequently used'),
        nature: $gettext('Animals & Nature'),
        objects: $gettext('Objects'),
        people: $gettext('Smileys & People'),
        places: $gettext('Travel & Places'),
        search: $gettext('Search Results'),
        symbols: $gettext('Symbols')
      },
      skins: {
        choose: $gettext('Choose default skin tone'),
        '1': $gettext('Default'),
        '2': $gettext('Light'),
        '3': $gettext('Medium-Light'),
        '4': $gettext('Medium'),
        '5': $gettext('Medium-Dark'),
        '6': $gettext('Dark')
      }
    }

    const data = (await import('@emoji-mart/data')).default

    const pickerOptions = {
      onEmojiSelect: (emoji: any) => emit('emojiSelect', emoji.native),
      onClickOutside: () => emit('clickOutside'),
      i18n,
      data,
      autoFocus: true,
      theme
    }

    // lazy loading to avoid loading the whole package on page load
    const { Picker } = await import('emoji-mart')
    const picker = new Picker(pickerOptions)

    isLoading.value = false
    await nextTick()

    unref(emojiPickerRef).innerHTML = ''
    unref(emojiPickerRef).appendChild(picker as any)
  },
  { immediate: true }
)
</script>
