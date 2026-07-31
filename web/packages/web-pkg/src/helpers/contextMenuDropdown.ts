import { OcDrop } from '@ownclouders/design-system/components'
import { ComponentPublicInstance } from 'vue'

export type ContextMenuBtnClickEventData = {
  event: MouseEvent | KeyboardEvent
  dropdown: ComponentPublicInstance<typeof OcDrop>
}

export type OcDropType = typeof OcDrop

const isKeyboardEvent = (event: Event): event is KeyboardEvent => {
  return (event as any).clientY === 0
}

export const displayPositionedDropdown = (
  dropdown: ComponentPublicInstance<typeof OcDrop>,
  event: MouseEvent | KeyboardEvent,
  contextMenuButton: ComponentPublicInstance<unknown>
) => {
  const contextMenuButtonPos = contextMenuButton.$el.getBoundingClientRect()

  const yValue = isKeyboardEvent(event)
    ? (event.target as HTMLElement)?.getBoundingClientRect().top || 0
    : event.clientY

  dropdown.setProps({
    getReferenceClientRect: () => ({
      width: 0,
      height: 0,
      top: yValue,
      bottom: yValue,
      /**
       * If event type is 'contextmenu' the trigger was a right click on the table row,
       * so we render the dropdown at the position of the mouse pointer.
       * Otherwise we render the dropdown at the position of the three-dot-menu
       */
      left:
        event.type === 'contextmenu' && !isKeyboardEvent(event)
          ? event.clientX
          : contextMenuButtonPos.x,
      right:
        event.type === 'contextmenu' && !isKeyboardEvent(event)
          ? event.clientX
          : contextMenuButtonPos.x
    })
  })

  dropdown.show()
}
