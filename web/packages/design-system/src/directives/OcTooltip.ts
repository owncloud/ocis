import tippy, { Instance } from 'tippy.js'
import merge from 'deepmerge'
import __logger from '../utils/logger'

export const hideOnEsc = {
  name: 'hideOnEsc',
  defaultValue: true,
  fn({ hide }: Instance) {
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.code === 'Escape') {
        hide()
      }
    }

    return {
      onShow: () => {
        document.addEventListener('keydown', onKeyDown)
      },
      onHide: () => {
        document.removeEventListener('keydown', onKeyDown)
      }
    }
  }
}

export const ariaHidden = {
  name: 'ariaHidden',
  defaultValue: true,
  fn(instance: Instance) {
    return {
      onCreate() {
        instance.popper.setAttribute('aria-hidden', 'true')
      }
    }
  }
}

export const destroy = (_tippy: Instance) => {
  if (!_tippy) {
    return
  }

  try {
    _tippy.destroy()
  } catch (e) {
    __logger(e)
  }
}

const initOrUpdate = (el: HTMLElement & { tooltip: any }, { value = {} }: any) => {
  if (Object.prototype.toString.call(value) !== '[object Object]') {
    value = { content: value }
  }

  if ((value.content !== 0 && !value.content) || value.content === '') {
    destroy(el.tooltip)
    el.tooltip = null
    return
  }

  const props = merge.all([
    {
      ignoreAttributes: true,
      interactive: true,
      aria: {
        content: null,
        expanded: false
      }
    },
    value
  ])

  if (!el.tooltip) {
    el.tooltip = tippy(el, {
      ...props,
      plugins: [hideOnEsc, ariaHidden]
    })
    return
  }

  el.tooltip.setProps(props)
}

export default {
  name: 'OcTooltip',
  beforeMount: initOrUpdate,
  updated: initOrUpdate,
  unmounted: (el: HTMLElement & { tooltip: any }) => destroy(el.tooltip)
}
