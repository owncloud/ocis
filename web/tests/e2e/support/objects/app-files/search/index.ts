import { Page } from '@playwright/test'
import * as po from './actions'

export class Search {
  #page: Page

  constructor({ page }: { page: Page }) {
    this.#page = page
  }

  getSearchResultMessage(): Promise<string> {
    return po.getSearchResultMessage({ page: this.#page })
  }

  async selectTagFilter({ tag: string }: { tag: string }): Promise<void> {
    await po.selectTagFilter({ tag: string, page: this.#page })
  }

  async selectMediaTypeFilter({ mediaType: string }: { mediaType: string }): Promise<void> {
    await po.selectMediaTypeFilter({ mediaType: string, page: this.#page })
  }

  async selectlastModifiedFilter({
    lastModified: string
  }: {
    lastModified: string
  }): Promise<void> {
    await po.selectLastModifiedFilter({ lastModified: string, page: this.#page })
  }

  async clearFilter({ filter: string }: { filter: string }): Promise<void> {
    await po.clearFilter({ page: this.#page, filter: string })
  }

  async toggleSearchTitleOnly({
    enableOrDisable: string
  }: {
    enableOrDisable: string
  }): Promise<void> {
    await po.toggleSearchTitleOnly({ enableOrDisable: string, page: this.#page })
  }
}
