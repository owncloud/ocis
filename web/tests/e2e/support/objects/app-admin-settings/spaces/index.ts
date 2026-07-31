import { Page } from '@playwright/test'
import * as po from './actions'
import { SpacesEnvironment } from '../../../environment'
import { Space } from '../../../types'
import { fileAction } from '../../../../environment/constants'

export class Spaces {
  #page: Page
  #spacesEnvironment: SpacesEnvironment

  constructor({ page }: { page: Page }) {
    this.#spacesEnvironment = new SpacesEnvironment()
    this.#page = page
  }

  getUUID({ key }: { key: string }): string {
    return this.getSpace({ key }).id
  }

  getDisplayedSpaces(): Promise<string[]> {
    return po.getDisplayedSpaces(this.#page)
  }

  getSpace({ key }: { key: string }): Space {
    return this.#spacesEnvironment.getSpace({ key })
  }

  async changeQuota({
    spaceIds,
    value,
    via
  }: {
    spaceIds: string[]
    value: string
    via: typeof fileAction.contextMenu | typeof fileAction.batchAction
  }): Promise<void> {
    await po.changeSpaceQuota({ spaceIds, value, page: this.#page, via })
  }

  async disable({
    spaceIds,
    via
  }: {
    spaceIds: string[]
    via: typeof fileAction.contextMenu | typeof fileAction.batchAction
  }): Promise<void> {
    await po.disableSpace({ spaceIds, page: this.#page, via })
  }

  async enable({
    spaceIds,
    via
  }: {
    spaceIds: string[]
    via: typeof fileAction.contextMenu | typeof fileAction.batchAction
  }): Promise<void> {
    await po.enableSpace({ spaceIds, page: this.#page, via })
  }

  async delete({
    spaceIds,
    via
  }: {
    spaceIds: string[]
    via: typeof fileAction.contextMenu | typeof fileAction.batchAction
  }): Promise<void> {
    await po.deleteSpace({ spaceIds, page: this.#page, via })
  }

  async select({ key }: { key: string }): Promise<void> {
    await po.selectSpace({ page: this.#page, id: this.getUUID({ key }) })
  }

  async renameSpaceUsingContextMenu({ key, value }: { key: string; value: string }): Promise<void> {
    await po.renameSpaceUsingContextMenu({ page: this.#page, id: this.getUUID({ key }), value })
  }

  async changeSubtitleUsingContextMenu({
    key,
    value
  }: {
    key: string
    value: string
  }): Promise<void> {
    await po.changeSpaceSubtitleUsingContextMenu({
      page: this.#page,
      id: this.getUUID({ key }),
      value
    })
  }

  async openPanel({ key }: { key: string }): Promise<void> {
    await po.openSpaceAdminSidebarPanel({ page: this.#page, id: this.getUUID({ key }) })
  }

  async openActionSideBarPanel({ action }: { action: string }): Promise<void> {
    await po.openSpaceAdminActionSidebarPanel({ page: this.#page, action })
  }

  listMembers({ filter }: { filter: string }): Promise<string[]> {
    return po.listSpaceMembers({ page: this.#page, filter })
  }
}
