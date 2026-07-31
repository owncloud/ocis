import { EventBus } from '../../../../src/services/eventBus'
import { SideBarEventTopics, useSideBar } from '../../../../src/composables/sideBar'
import { unref, ref } from 'vue'
import { getComposableWrapper } from '@ownclouders/web-test-helpers'
import * as localStorage from '../../../../src/composables/localStorage'

const localStorageSpy = vi
  .spyOn(localStorage, 'useLocalStorage')
  .mockImplementation(() => ref(false))

describe('useSideBar', () => {
  let eventBus: EventBus
  beforeEach(() => {
    eventBus = new EventBus()
  })
  describe('initial state', () => {
    it('should have "isSideBarOpen" as "false"', () => {
      getComposableWrapper(() => {
        const { isSideBarOpen } = useSideBar({ bus: eventBus })
        expect(unref(isSideBarOpen)).toBe(false)
      })
    })
    it('should have "sideBarActivePanel" as "null"', () => {
      getComposableWrapper(() => {
        const { sideBarActivePanel } = useSideBar({ bus: eventBus })
        expect(unref(sideBarActivePanel)).toBeNull()
      })
    })
  })
  describe('open', () => {
    it('should set "isSideBarOpen" to "true"', () => {
      getComposableWrapper(() => {
        const { isSideBarOpen } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.open)
        expect(unref(isSideBarOpen)).toBe(true)
      })
    })
    it('should set "sideBarActivePanel" to "null"', () => {
      getComposableWrapper(() => {
        const { sideBarActivePanel } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.open)
        expect(unref(sideBarActivePanel)).toBeNull()
      })
    })
  })
  describe('close', () => {
    it('should set "isSideBarOpen" to "false"', () => {
      getComposableWrapper(() => {
        const { isSideBarOpen } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.close)
        expect(unref(isSideBarOpen)).toBe(false)
      })
    })
    it('should set "sideBarActivePanel" to "null"', () => {
      getComposableWrapper(() => {
        const { sideBarActivePanel } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.close)
        expect(unref(sideBarActivePanel)).toBeNull()
      })
    })
  })
  describe('toggle', () => {
    it('should toggle "isSideBarOpen" back and forth', () => {
      getComposableWrapper(() => {
        const { isSideBarOpen } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.toggle)
        expect(unref(isSideBarOpen)).toBe(true)
        eventBus.publish(SideBarEventTopics.toggle)
        expect(unref(isSideBarOpen)).toBe(false)
      })
    })
    it('should not influence "sideBarActivePanel"', () => {
      getComposableWrapper(() => {
        const { sideBarActivePanel } = useSideBar({ bus: eventBus })
        // initial state
        eventBus.publish(SideBarEventTopics.toggle)
        expect(unref(sideBarActivePanel)).toBe(null)
        // modified active panel
        eventBus.publish(SideBarEventTopics.setActivePanel, 'SomePanel')
        eventBus.publish(SideBarEventTopics.toggle)
        expect(unref(sideBarActivePanel)).toBe('SomePanel')
      })
    })
  })
  describe('openWithPanel', () => {
    it('should set "isSideBarOpen" to "true"', () => {
      getComposableWrapper(() => {
        const { isSideBarOpen } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.openWithPanel, 'SomePanel')
        expect(unref(isSideBarOpen)).toBe(true)
      })
    })
    it('should set "sideBarActivePanel" to provided value', () => {
      getComposableWrapper(() => {
        const { sideBarActivePanel } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.openWithPanel, 'SomePanel')
        expect(unref(sideBarActivePanel)).toBe('SomePanel')
      })
    })
  })
  describe('setActivePanel', () => {
    it('should not influence "isSideBarOpen"', () => {
      getComposableWrapper(() => {
        const { isSideBarOpen } = useSideBar({ bus: eventBus })
        expect(unref(isSideBarOpen)).toBe(false)
      })
    })
    it('should set "sideBarActivePanel" to provided value', () => {
      getComposableWrapper(() => {
        const { sideBarActivePanel } = useSideBar({ bus: eventBus })
        eventBus.publish(SideBarEventTopics.setActivePanel, 'SomePanel')
        expect(unref(sideBarActivePanel)).toBe('SomePanel')
      })
    })
  })

  describe('embedMode', () => {
    beforeEach(() => {
      localStorageSpy.mockImplementationOnce(() => ref(true))
    })

    it('should use local ref when embed mode is enabled', () => {
      getComposableWrapper(
        () => {
          const { isSideBarOpen } = useSideBar({ bus: eventBus })
          expect(unref(isSideBarOpen)).toBe(false)
        },
        {
          pluginOptions: {
            piniaOptions: {
              configState: {
                options: {
                  embed: { enabled: true }
                }
              }
            }
          }
        }
      )
    })

    it('should use local storage when embed mode is disabled', () => {
      getComposableWrapper(
        () => {
          const { isSideBarOpen } = useSideBar({ bus: eventBus })
          expect(unref(isSideBarOpen)).toBe(true)
        },
        {
          pluginOptions: {
            piniaOptions: {
              configState: {
                options: {
                  embed: { enabled: false }
                }
              }
            }
          }
        }
      )
    })
  })
})
