import { Page } from '@playwright/test'
import { UsersEnvironment } from '../../../environment'
import { getWorld } from '../../../../environment/world'
import * as po from './actions'

export class Groups {
  #page: Page
  #usersEnvironment: UsersEnvironment

  constructor({ page }: { page: Page }) {
    this.#usersEnvironment = new UsersEnvironment()
    this.#page = page
  }

  getUUID({ key }: { key: string }): string {
    return this.#usersEnvironment.getCreatedGroup({ key }).uuid
  }

  async createGroup({ key }: { key: string }): Promise<void> {
    const world = getWorld()
    const group = this.#usersEnvironment.getGroup({ key })
    const response = await po.createGroup({ page: this.#page, key: group.displayName })
    const actualId = world.getGroupId(key)
    this.#usersEnvironment.storeCreatedGroup({
      group: {
        id: actualId,
        uuid: response['id'],
        displayName: response['displayName']
      }
    })
  }

  getDisplayedGroupsIds(): Promise<string[]> {
    return po.getDisplayedGroupsIds({ page: this.#page })
  }

  getGroupsDisplayName(): Promise<string> {
    return po.getGroupsDisplayName({ page: this.#page })
  }

  async selectGroup({ key }: { key: string }): Promise<void> {
    await po.selectGroup({ page: this.#page, uuid: this.getUUID({ key }) })
  }

  async deleteGroupUsingBatchAction({ groupIds }: { groupIds: string[] }): Promise<void> {
    await po.deleteGrouprUsingBatchAction({ page: this.#page, groupIds })
  }

  async deleteGroupUsingContextMenu({ key }: { key: string }): Promise<void> {
    await po.deleteGroupUsingContextMenu({ page: this.#page, uuid: this.getUUID({ key }) })
  }

  async changeGroup({
    key,
    attribute,
    value,
    action
  }: {
    key: string
    attribute: string
    value: string
    action: string
  }): Promise<void> {
    const displayName =
      attribute === 'displayName'
        ? this.#usersEnvironment.getGroupDisplayName({ displayName: value })
        : value

    const uuid = this.getUUID({ key })
    await po.openEditPanel({ page: this.#page, uuid, action })
    await po.changeGroup({ uuid, attribute: attribute, value: displayName, page: this.#page })
  }
}
