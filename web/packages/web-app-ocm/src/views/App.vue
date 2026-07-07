<template>
  <div class="sciencemesh">
    <div class="oc-flex oc-flex-column sciencemesh-wrapper">
      <div class="oc-flex sciencemesh-top">
        <div id="sciencemesh-invite" class="oc-width-1-2">
          <outgoing-invitations />
        </div>
        <div id="sciencemesh-accept-invites" class="oc-width-1-2">
          <incoming-invitations @highlight-new-connections="highlightNewConnections" />
        </div>
      </div>
      <div id="sciencemesh-connections">
        <connections-panel
          v-model:connections="connections"
          :highlighted-connections="highlightedConnections.map((c) => c.id)"
          :loading="loadingConnections"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, onUnmounted, ref, unref, Ref } from 'vue'
import ConnectionsPanel from './ConnectionsPanel.vue'
import IncomingInvitations from './IncomingInvitations.vue'
import OutgoingInvitations from './OutgoingInvitations.vue'
import {
  useClientService,
  useScrollTo,
  FederatedConnection,
  useMessages
} from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'
import { buildConnection } from '../functions'

const { showMessage } = useMessages()
const { scrollToResource } = useScrollTo()
const clientSerivce = useClientService()
const { $gettext } = useGettext()

const connections: Ref<FederatedConnection[]> = ref([])
const highlightedConnections: Ref<FederatedConnection[]> = ref([])
const highlightNewConnectionsInterval = ref(null)
const loadingConnections = ref(true)

const findAcceptedUsers = async () => {
  try {
    const { data: acceptedUsers } = await clientSerivce.httpAuthenticated.get<
      FederatedConnection[]
    >('/sciencemesh/find-accepted-users')
    loadingConnections.value = false
    connections.value = acceptedUsers.map(buildConnection)
  } catch {
    connections.value = []
    loadingConnections.value = false
  }
}

const highlightNewConnections = async () => {
  const oldConnections = [...unref(connections)]
  await findAcceptedUsers()
  if (oldConnections.length < unref(connections).length) {
    highlightedConnections.value = unref(connections).filter(
      (connection) => !oldConnections.map((c) => c.id).includes(connection.id)
    )
    if (unref(highlightedConnections).length === 1) {
      scrollToResource(unref(highlightedConnections)[0].id)
      showMessage({
        title: $gettext('New federated connection'),
        status: 'success',
        desc: $gettext('You can share with and recieve shares from %{user} now', {
          user: unref(highlightedConnections)[0].display_name
        })
      })
    } else if (unref(highlightedConnections).length > 1) {
      const newConnections = unref(highlightedConnections)
        .map((c) => c.display_name)
        .join(', ')

      showMessage({
        title: $gettext('New federated connections'),
        status: 'success',
        desc: $gettext('You can share with and receive shares from %{ connections } now', {
          connections: newConnections
        })
      })
    }
  }
}

onMounted(async () => {
  await findAcceptedUsers()
  loadingConnections.value = false
  highlightNewConnectionsInterval.value = setInterval(() => {
    highlightNewConnections()
  }, 10 * 1000)
})

onUnmounted(() => {
  clearInterval(unref(highlightNewConnectionsInterval))
})
</script>

<style lang="scss" scoped>
.sciencemesh {
  background-color: var(--oc-color-background-hover);
  overflow: auto;
}
.sciencemesh-wrapper {
  height: 100%;
}
.sciencemesh-top {
  max-height: 360px;
  @media (max-width: $oc-breakpoint-large-default) {
    flex-direction: column;
    height: auto;
    max-height: unset;
  }
}
#sciencemesh-invite,
#sciencemesh-accept-invites {
  margin: var(--oc-space-small);
  overflow: auto;
}
#sciencemesh-invite,
#sciencemesh-accept-invites,
#sciencemesh-connections {
  margin: var(--oc-space-small);
  background-color: var(--oc-color-background-default);
  border-radius: 15px;
  padding: var(--oc-space-small);
  @media (max-width: $oc-breakpoint-large-default) {
    width: auto;
  }
}
#sciencemesh-connections {
  flex: 1;
}
#sciencemesh-invite {
  overflow: auto;
  @media (max-width: $oc-breakpoint-large-default) {
    margin-bottom: var(--oc-space-small);
  }
}
</style>
