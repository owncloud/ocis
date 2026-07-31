import Preview from './preview'
import List from './list'
import { Router } from 'vue-router'
import {
  CapabilityStore,
  ConfigStore,
  SearchFunction,
  SearchList,
  SearchPreview,
  SearchProvider
} from '@ownclouders/web-pkg'

function $gettext(msg: string) {
  return msg
}

export default class Provider implements SearchProvider {
  public readonly id: string
  public readonly displayName: string
  public readonly previewSearch: SearchPreview
  public readonly listSearch: SearchList
  private readonly capabilityStore: CapabilityStore

  constructor(
    capabilityStore: CapabilityStore,
    router: Router,
    searchFunction: SearchFunction,
    configStore: ConfigStore
  ) {
    this.id = 'files.sdk'
    this.displayName = $gettext('Files')
    this.previewSearch = new Preview(router, searchFunction, configStore)
    this.listSearch = new List(searchFunction)
    this.capabilityStore = capabilityStore
  }

  public get available(): boolean {
    return this.capabilityStore.davReports.includes('search-files')
  }
}
