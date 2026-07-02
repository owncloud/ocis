import { displayPositionedDropdown } from '../../../src/helpers'

describe('displayPositionedDropdown', () => {
  it('shows the dropdown', () => {
    const dropdown = { setProps: vi.fn(), show: vi.fn() } as any
    const ctxMenuBtn = { $el: { getBoundingClientRect: vi.fn() } } as any
    displayPositionedDropdown(dropdown, {} as MouseEvent, ctxMenuBtn)
    expect(dropdown.show).toHaveBeenCalled()
  })
  it('horizontally positions the drop at the cursor for "contextmenu"-event', () => {
    const event = { type: 'contextmenu', clientX: 100 } as MouseEvent
    let dropdownProps = { getReferenceClientRect: vi.fn(() => undefined) }
    const dropdown = {
      setProps: vi.fn((props) => {
        dropdownProps = props
      }),
      show: vi.fn()
    } as any
    const ctxMenuBtn = { $el: { getBoundingClientRect: vi.fn(() => ({})) } } as any
    displayPositionedDropdown(dropdown, event, ctxMenuBtn)
    const data = dropdownProps.getReferenceClientRect()
    expect(data.left).toEqual(event.clientX)
    expect(data.right).toEqual(event.clientX)
  })
  it('horizontally positions the drop at context menu button if no "contextmenu"-event given', () => {
    const event = { clientX: 100 } as MouseEvent
    let dropdownProps = { getReferenceClientRect: vi.fn(() => undefined) }
    const dropdown = {
      setProps: vi.fn((props) => {
        dropdownProps = props
      }),
      show: vi.fn()
    } as any
    const ctxMenuBtnX = 200
    const ctxMenuBtn = { $el: { getBoundingClientRect: vi.fn(() => ({ x: ctxMenuBtnX })) } } as any
    displayPositionedDropdown(dropdown, event, ctxMenuBtn)
    const data = dropdownProps.getReferenceClientRect()
    expect(data.left).toEqual(ctxMenuBtnX)
    expect(data.right).toEqual(ctxMenuBtnX)
  })
  it('vertically positions the drop via "clientY" if given', () => {
    const event = { clientY: 100 } as MouseEvent
    let dropdownProps = { getReferenceClientRect: vi.fn(() => undefined) }
    const dropdown = {
      setProps: vi.fn((props) => {
        dropdownProps = props
      }),
      show: vi.fn()
    } as any
    const ctxMenuBtn = { $el: { getBoundingClientRect: vi.fn(() => ({ x: 200 })) } } as any
    displayPositionedDropdown(dropdown, event, ctxMenuBtn)
    const { top, bottom } = dropdownProps.getReferenceClientRect()
    expect(top).toEqual(event.clientY)
    expect(bottom).toEqual(event.clientY)
  })
  it('vertically positions the drop via the context button position if "clientY" is 0', () => {
    const yPos = 200
    const event = {
      clientY: 0,
      target: { getBoundingClientRect: () => ({ top: yPos }) }
    } as unknown as MouseEvent
    let dropdownProps = { getReferenceClientRect: vi.fn(() => undefined) }
    const dropdown = {
      setProps: vi.fn((props) => {
        dropdownProps = props
      }),
      show: vi.fn()
    } as any
    const ctxMenuBtn = { $el: { getBoundingClientRect: vi.fn(() => ({ x: 200 })) } } as any
    displayPositionedDropdown(dropdown, event, ctxMenuBtn)
    const { top, bottom } = dropdownProps.getReferenceClientRect()
    expect(top).toEqual(yPos)
    expect(bottom).toEqual(yPos)
  })
})
