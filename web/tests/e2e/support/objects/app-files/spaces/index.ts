import { Locator, Page } from '@playwright/test'
import { SpacesEnvironment, LinksEnvironment } from '../../../environment'
import { File } from '../../../types'
import * as po from './actions'
import { spaceLocator } from './utils'
import { ICollaborator } from '../share/collaborator'

export class Spaces {
  #page: Page
  #spacesEnvironment: SpacesEnvironment
  #linksEnvironment: LinksEnvironment

  constructor({ page }: { page: Page }) {
    this.#page = page
    this.#spacesEnvironment = new SpacesEnvironment()
    this.#linksEnvironment = new LinksEnvironment()
  }

  getSpaceID({ key }: { key: string }): string {
    const { id } = this.#spacesEnvironment.getSpace({ key })
    return id
  }

  async create({
    key,
    space
  }: {
    key: string
    space: Omit<po.createSpaceArgs, 'page'>
  }): Promise<void> {
    const id = await po.createSpace({ ...space, page: this.#page })
    this.#spacesEnvironment.createSpace({
      key,
      space: { name: space.name, id }
    })
  }

  async open({ key }: { key: string }): Promise<void> {
    const { id } = this.#spacesEnvironment.getSpace({ key })
    await po.openSpace({ page: this.#page, id })
  }

  async changeName({ key, value }: { key: string; value: string }): Promise<void> {
    const { id } = this.#spacesEnvironment.getSpace({ key })
    await po.changeSpaceName({ id, value, page: this.#page })
  }

  async changeSubtitle({ key, value }: { key: string; value: string }): Promise<void> {
    const { id } = this.#spacesEnvironment.getSpace({ key })
    await po.changeSpaceSubtitle({ id, value, page: this.#page })
  }

  async changeDescription({ value }: { value: string }): Promise<void> {
    await po.changeSpaceDescription({ value, page: this.#page })
  }

  async changeQuota({ key, value }: { key: string; value: string }): Promise<void> {
    const { id } = this.#spacesEnvironment.getSpace({ key })
    await po.changeQuota({ id, value, page: this.#page })
  }

  async addMembers(args: Omit<po.SpaceMembersArgs, 'page'>): Promise<void> {
    await po.addSpaceMembers({ ...args, page: this.#page })
  }

  async removeAccessToMember(args: Omit<po.removeAccessMembersArgs, 'page'>): Promise<void> {
    await po.removeAccessSpaceMembers({ ...args, page: this.#page })
  }

  getSpaceLocator(space: string): Locator {
    const spaceID = this.getSpaceID({ key: space })
    return spaceLocator({ spaceID, page: this.#page })
  }

  async changeRoles(args: Omit<po.SpaceMembersArgs, 'page'>): Promise<void> {
    await po.changeSpaceRole({ ...args, page: this.#page })
  }

  async changeSpaceImage({ key, resource }: { key: string; resource: File }): Promise<void> {
    const { id } = this.#spacesEnvironment.getSpace({ key })
    await po.changeSpaceImage({ id, resource, page: this.#page })
  }

  async createPublicLink({ password }: { password: string }): Promise<void> {
    const url = await po.createPublicLinkForSpace({ page: this.#page, password })
    this.#linksEnvironment.createLink({
      key: 'Unnamed link',
      link: { name: 'Unnamed link', url }
    })
  }

  async addExpirationDate({
    member,
    expirationDate
  }: {
    member: Omit<ICollaborator, 'role'>
    expirationDate: string
  }): Promise<void> {
    await po.addExpirationDateToMember({ member, expirationDate, page: this.#page })
  }

  async removeExpirationDate({ member }: { member: Omit<ICollaborator, 'role'> }): Promise<void> {
    await po.removeExpirationDateFromMember({ member, page: this.#page })
  }

  downloadSpace(): Promise<string> {
    return po.downloadSpace(this.#page)
  }

  async checkSpaceActivity({ activity }: { activity: string }): Promise<void> {
    await po.checkSpaceActivity({ page: this.#page, activity })
  }
}
