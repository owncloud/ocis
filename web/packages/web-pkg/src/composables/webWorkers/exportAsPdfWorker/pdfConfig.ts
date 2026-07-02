import { rgb } from 'pdf-lib'

export const PDF_THEME = Object.freeze({
  layout: {
    margin: 50
  },
  font: {
    baseSize: 12,
    lineHeight: 16,
    subSupSize: 9,
    codeSize: 10,
    codeLineHeight: 14,
    listBulletSize: 12,
    listItemLineHeight: 16,
    blockquoteLineHeight: 16,
    tableHeaderTextSize: 11,
    tableHeaderLineHeight: 14,
    tableCellTextSize: 10,
    tableCellLineHeight: 14,
    imageTitleSize: 10,
    h1: 24,
    h2: 18,
    h3: 15,
    h4: 12,
    h5: 10.5,
    h6: 10.2
  },
  offset: {
    subscript: -3,
    superscript: 4
  },
  color: {
    text: rgb(0, 0, 0),
    link: rgb(0, 0, 0.8),
    error: rgb(0.8, 0.2, 0.2),
    codeSpan: rgb(0.7, 0.1, 0.1),
    codeBlockBg: rgb(0.15625, 0.171875, 0.203125),
    codeBlockText: rgb(0.875, 0.875, 0.875),
    blockquoteBar: rgb(0.208, 0.702, 0.471),
    blockquoteText: rgb(0.3, 0.3, 0.3),
    tableHeaderBg: rgb(0.9, 0.9, 0.9),
    tableBorder: rgb(0.5, 0.5, 0.5),
    tableCellBg: rgb(1, 1, 1),
    hr: rgb(0.5, 0.5, 0.5),
    imagePlaceholder: rgb(0.5, 0.5, 0.5)
  },
  spacing: {
    md: 12,
    listIndent: 20,
    listGap: 8,
    sm: 8
  },
  codeBlock: {
    padding: 10
  },
  table: {
    cellPadding: 8,
    rowHeight: 30
  },
  blockquote: {
    barWidth: 3,
    barXOffset: 10,
    nestedQuoteYOffset: 10
  },
  hr: {
    thickness: 1
  },
  image: {
    contentPadding: 40
  },
  math: {
    displayModePadding: 20,
    inlineModePadding: 5
  },
  svg: {
    scaleFactor: 3
  },
  underline: {
    thickness: 1,
    offset: -3
  },
  strikeThrough: {
    thickness: 1
  }
} as const)
