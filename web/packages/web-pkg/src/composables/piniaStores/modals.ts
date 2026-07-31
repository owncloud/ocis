import { ContextualHelperData } from '@ownclouders/design-system/helpers'
import { defineStore } from 'pinia'
import { v4 as uuidV4 } from 'uuid'
import { Component, ComponentPublicInstance, computed, ref, unref } from 'vue'

export type CustomModalComponent = Component<{ modal?: Modal } & Record<string, any>>

export type CustomModalComponentEmits = {
  confirm?: () => Promise<void>
  cancel?: () => void
  'update:confirmDisabled': (value: boolean) => void
} & Record<string, any>

export type CustomModalComponentInstance = ComponentPublicInstance<
  { modal: Modal } & unknown,
  { onConfirm?: () => Promise<unknown>; onCancel?: () => unknown },
  any,
  any,
  any,
  CustomModalComponentEmits
>

export type Modal = {
  id: string
  elementId?: string
  elementClass?: string
  title: string
  variation?: string
  icon?: string
  message?: string
  cancelText?: string
  confirmDisabled?: boolean
  isLoading?: boolean
  confirmText?: string
  hideActions?: boolean
  hideConfirmButton?: boolean
  hasInput?: boolean
  inputValue?: string
  inputType?: string
  inputSelectionRange?: [number, number]
  inputLabel?: string
  inputError?: string
  inputDescription?: string
  focusTrapInitial?: string | false
  contextualHelperLabel?: string
  contextualHelperData?: ContextualHelperData
  customComponent?: CustomModalComponent
  customComponentAttrs?: () => Record<string, unknown>
  onCancel?: () => void
  onConfirm?: (value: unknown) => void | Promise<void>
  onInput?: (value: string, setError: (error: string) => void) => void
}

export const useModals = defineStore('modals', () => {
  const modals = ref<Modal[]>([])
  const activeModal = computed(() => unref(modals).at(-1))

  const getModal = (id: Modal['id']) => {
    return unref(modals).find((modal) => modal.id === id)
  }

  const dispatchModal = (data: Omit<Modal, 'id'>) => {
    const modal = { ...data, id: uuidV4() as string }
    modals.value.push(modal)
    return modal
  }

  const updateModal = <T extends Modal, K extends keyof Modal>(
    id: T['id'],
    key: K,
    value: T[K]
  ) => {
    const modal = getModal(id)
    modal[key] = value
  }

  const removeModal = (id: Modal['id']) => {
    modals.value = unref(modals).filter((modal) => modal.id !== id)
  }

  const removeAllModals = () => {
    modals.value = []
  }

  const setModalActive = (id: Modal['id']) => {
    const foundIdx = unref(modals).findIndex((modal) => modal.id === id)
    if (foundIdx < 0) {
      return
    }

    const modal = getModal(id)
    unref(modals).splice(foundIdx, 1)
    modals.value.push(modal)
  }

  return {
    modals,
    activeModal,
    dispatchModal,
    updateModal,
    removeModal,
    removeAllModals,
    setModalActive
  }
})

export type ModalStore = ReturnType<typeof useModals>
