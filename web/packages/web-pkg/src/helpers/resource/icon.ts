export type IconFillType = 'fill' | 'line' | 'none'
export type IconType = {
  name: string
  color?: string
  fillType?: IconFillType
}

export type ResourceIconMapping = Record<'mimeType' | 'extension', Record<string, IconType>>
export const resourceIconMappingInjectionKey = 'oc-resource-icon-mapping'

const fileIcon = {
  archive: {
    icon: { name: 'resource-type-archive', color: 'var(--oc-color-icon-archive)' },
    extensions: [
      '7z',
      'apk',
      'bz2',
      'deb',
      'gz',
      'gzip',
      'rar',
      'tar',
      'tar.bz2',
      'tar.gz',
      'tar.xz',
      'tbz2',
      'tgz',
      'zip'
    ]
  },
  audio: {
    icon: { name: 'resource-type-audio', color: 'var(--oc-color-icon-audio)' },
    extensions: [
      '3gp',
      '8svx',
      'aa',
      'aac',
      'aax',
      'act',
      'aiff',
      'alac',
      'amr',
      'ape',
      'au',
      'awb',
      'cda',
      'dss',
      'dvf',
      'flac',
      'gsm',
      'iklax',
      'ivs',
      'm4a',
      'm4b',
      'm4p',
      'mmf',
      'mogg',
      'movpkg',
      'mp3',
      'mpc',
      'msv',
      'nmf',
      'oga',
      'ogga',
      'opus',
      'ra',
      'raw',
      'rf64',
      'rm',
      'sln',
      'tta',
      'voc',
      'vox',
      'wav',
      'wma',
      'wv'
    ]
  },
  code: {
    icon: { name: 'resource-type-code', color: 'var(--oc-color-text-default)' },
    extensions: [
      'bash',
      'c++',
      'c',
      'cc',
      'cpp',
      'css',
      'feature',
      'go',
      'h',
      'hh',
      'hpp',
      'htm',
      'html',
      'java',
      'js',
      'json',
      'php',
      'pl',
      'py',
      'scss',
      'sh',
      'sh-lib',
      'sql',
      'ts',
      'xml',
      'yaml',
      'yml'
    ]
  },
  default: {
    icon: { name: 'resource-type-file', color: 'var(--oc-color-text-default)' },
    extensions: ['accdb', 'rss', 'swf']
  },
  drawio: {
    icon: { name: 'resource-type-drawio', color: 'var(--oc-color-icon-drawio)' },
    extensions: ['drawio']
  },
  document: {
    icon: { name: 'resource-type-document', color: 'var(--oc-color-icon-document)' },
    extensions: ['doc', 'docm', 'docx', 'dot', 'dotx', 'lwp', 'odt', 'one', 'wpd']
  },
  ifc: {
    icon: { name: 'resource-type-ifc', color: 'var(--oc-color-icon-ifc)' },
    extensions: ['ifc']
  },
  ipynb: {
    icon: { name: 'resource-type-jupyter', color: 'var(--oc-color-icon-jupyter)' },
    extensions: ['ipynb']
  },
  image: {
    icon: { name: 'resource-type-image', color: 'var(--oc-color-icon-image)' },
    extensions: [
      'ai',
      'cdr',
      'eot',
      'eps',
      'gif',
      'jpeg',
      'jpg',
      'otf',
      'pfb',
      'png',
      'ps',
      'psd',
      'svg',
      'ttf',
      'woff',
      'xcf'
    ]
  },
  form: {
    icon: { name: 'resource-type-form', color: 'var(--oc-color-icon-form)' },
    extensions: ['docf', 'docxf', 'oform']
  },
  markdown: {
    icon: { name: 'resource-type-markdown', color: 'var(--oc-color-icon-markdown)' },
    extensions: ['md', 'markdown']
  },
  odg: {
    icon: { name: 'resource-type-graphic', color: 'var(--oc-color-icon-graphic)' },
    extensions: ['odg']
  },
  pdf: {
    icon: { name: 'resource-type-pdf', color: 'var(--oc-color-icon-pdf)' },
    extensions: ['pdf']
  },
  presentation: {
    icon: { name: 'resource-type-presentation', color: 'var(--oc-color-icon-presentation)' },
    extensions: [
      'odp',
      'pot',
      'potm',
      'potx',
      'ppa',
      'ppam',
      'pps',
      'ppsm',
      'ppsx',
      'ppt',
      'pptm',
      'pptx'
    ]
  },
  root: {
    icon: { name: 'resource-type-root', color: 'var(--oc-color-icon-root)' },
    extensions: ['root']
  },
  spreadsheet: {
    icon: { name: 'resource-type-spreadsheet', color: 'var(--oc-color-icon-spreadsheet)' },
    extensions: ['csv', 'ods', 'xla', 'xlam', 'xls', 'xlsb', 'xlsm', 'xlsx', 'xlt', 'xltm', 'xltx']
  },
  text: {
    icon: { name: 'resource-type-text', color: 'var(--oc-color-text-default)' },
    extensions: ['cb7', 'cba', 'cbr', 'cbt', 'cbtc', 'cbz', 'cvbdl', 'eml', 'mdb', 'tex', 'txt']
  },
  url: {
    icon: { name: 'resource-type-url', color: 'var(--oc-color-text-default)' },
    extensions: ['url']
  },
  video: {
    icon: {
      name: 'resource-type-video',
      color: 'var(--oc-color-icon-video)'
    },
    extensions: ['mov', 'mp4', 'webm', 'wmv']
  },
  epub: {
    icon: { name: 'resource-type-book', color: 'var(--oc-color-icon-epub)' },
    extensions: ['epub']
  },
  board: {
    icon: { name: 'resource-type-board' },
    extensions: ['ggs']
  },
  psec: {
    icon: { name: 'file-psec', color: 'var(--oc-color-icon-folder)' },
    extensions: ['psec']
  },
  pinboard: {
    icon: { name: 'resource-type-pinboard-beta' },
    extensions: ['ggp']
  },
  excalidraw: {
    icon: { name: 'resource-type-excalidraw' },
    extensions: ['excalidraw']
  },
  visio: {
    icon: { name: 'resource-type-diagram', color: 'var(--oc-color-icon-visio)' },
    extensions: ['vsd', 'vsdm', 'vsdx', 'vss', 'vssm', 'vssx', 'vst', 'vstm', 'vstx']
  }
}

export function createDefaultFileIconMapping() {
  const fileIconMapping: Record<string, IconType> = {}

  Object.values(fileIcon).forEach((value) => {
    value.extensions.forEach((extension) => {
      fileIconMapping[extension] = value.icon
    })
  })

  return fileIconMapping
}
