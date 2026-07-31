import { Page } from '@playwright/test'
import * as po from './actions'
import { SpacesEnvironment } from '../../../environment'

export class Trashbin {
  #page: Page
  #spacesEnvironment: SpacesEnvironment

  constructor({ page }: { page: Page }) {
    this.#page = page
    this.#spacesEnvironment = new SpacesEnvironment()
  }

  async open(key: string): Promise<void> {
    const { id } = this.#spacesEnvironment.getSpace({ key })
    await po.openTrashbin({ page: this.#page, id })
  }
}
