import { useModals, useUserStore } from '@ownclouders/web-pkg'
import { storeToRefs } from 'pinia'
import { computed, unref } from 'vue'
import { useGettext } from 'vue3-gettext'
import InstancesModal from '../../components/InstancesModal.vue'

/** Limits the number of instances displayed in the user menu list */
const INLINE_INSTANCES_LIMIT = 3

export function useInstances() {
  const { dispatchModal } = useModals()
  const { $pgettext } = useGettext()

  const userStore = useUserStore()
  const { user } = storeToRefs(userStore)

  const currentInstance = computed(() => window.location.origin)

  const instances = computed(() => {
    if (unref(user).instances.length < 1) {
      return []
    }

    return unref(user)
      .instances.map((instance) => ({
        ...instance,
        active: instance.url === unref(currentInstance)
      }))
      .toSorted((a, b) => {
        if (a.active !== b.active) {
          return a.active ? -1 : 1
        }

        if (a.primary !== b.primary) {
          return a.primary ? -1 : 1
        }

        return a.url.localeCompare(b.url)
      })
  })

  const inlineInstances = computed(() => unref(instances).slice(0, INLINE_INSTANCES_LIMIT))

  const canOpenInstancesModal = computed(() => unref(instances).length > INLINE_INSTANCES_LIMIT)

  function showInstancesModal() {
    dispatchModal({
      title: $pgettext(
        'The instances modal title available when multiple instances are enabled in oCIS',
        'Instances'
      ),
      customComponent: InstancesModal,
      customComponentAttrs: () => ({
        modal: {
          id: 'instances-modal'
        }
      }),
      hideConfirmButton: true,
      cancelText: $pgettext(
        'The close instances modal action label in the instances modal available when multiple instances are enabled in oCIS',
        'Close'
      )
    })
  }

  return {
    instances,
    inlineInstances,
    canOpenInstancesModal,
    showInstancesModal
  }
}
