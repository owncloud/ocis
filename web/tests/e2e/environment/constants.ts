export const fileAction = {
  contextMenu: 'context menu',
  batchAction: 'batch actions',
  sideBarPanel: 'sidebar panel',
  quickAction: 'QUICK_ACTION',
  urlNavigation: 'URL_NAVIGATION',
  singleShareView: 'SINGLE_SHARE_VIEW',
  previewTopBar: 'PREVIEW_TOPBAR',
  keyboard: 'KEYBOARD',
  dropDownMenu: 'DROP_DOWN_MENU',
  dragDrop: 'DRAG_DROP',
  dragDropBreadcrumb: 'DRAG_DROP_BREADCRUMB'
} as const

export const client = {
  mobile: 'mobile',
  desktop: 'desktop'
} as const

export const application = {
  textEditor: 'texteditor',
  pdfViewer: 'pdfviewer',
  mediaViewer: 'mediaviewer',
  collabora: 'Collabora',
  onlyOffice: 'OnlyOffice'
} as const

export const searchScope = {
  allFiles: 'all files',
  currentFolder: 'current folder'
} as const

export const searchFilter = {
  mediaType: 'mediaType',
  tags: 'tags',
  lastModified: 'lastModified',
  fullText: 'fullText'
} as const

export const resourcePage = {
  searchList: 'search list',
  filesList: 'files list',
  shares: 'Shares',
  trashbin: 'trashbin'
} as const

export const shareIndicator = {
  linkDirect: 'link-direct',
  linkIndirect: 'link-indirect',
  userDirect: 'user-direct',
  userIndirect: 'user-indirect'
} as const
