import { defineComponent } from 'vue'
import { useFileListHeaderPosition } from '../../../../../src/composables/fileListHeaderPosition'
import { defaultPlugins, mount } from '@ownclouders/web-test-helpers'

export const createWrapper = () =>
  mount(
    defineComponent({
      setup() {
        const { y, refresh } = useFileListHeaderPosition()

        return { y, refresh }
      },
      template: `<div>{{ y }}</div>`
    }),
    {
      attachTo: document.body,
      global: {
        plugins: [...defaultPlugins()]
      }
    }
  )

export const createAppBar = (): {
  createElement: () => void
  resize: (height: number) => void
  remove: () => void
} => {
  const createElement = () => {
    const appBar = document.createElement('div')
    appBar.setAttribute('id', 'files-app-bar')

    document.body.appendChild(appBar)
  }
  const resize = (height = 0) => {
    document.getElementById('files-app-bar').getBoundingClientRect = vi
      .fn()
      .mockReturnValue({ height })
  }
  const remove = () => {
    document.getElementById('files-app-bar').remove()
  }

  return { createElement, resize, remove }
}
