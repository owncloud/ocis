import { TokenizerAndRendererExtension } from 'marked'

export const markedExtensions: TokenizerAndRendererExtension[] = [
  {
    name: 'superscript',
    level: 'inline',
    start(src) {
      return src.match(/\^/)?.index
    },
    tokenizer(src) {
      const match = src.match(/^\^([^\^\s]+)\^/)

      if (!match) {
        return
      }

      return {
        type: 'sup',
        raw: match[0],
        text: match[1]
      }
    },
    renderer(token) {
      return `<sup>${token.text}</sup>`
    }
  },
  {
    name: 'subscript',
    level: 'inline',
    start(src) {
      return src.match(/~/)?.index
    },
    tokenizer(src) {
      const match = src.match(/^~([^~\s]+)~/)

      if (!match) {
        return
      }

      return {
        type: 'sub',
        raw: match[0],
        text: match[1]
      }
    },
    renderer(token) {
      return `<sub>${token.text}</sub>`
    }
  },
  {
    name: 'underline',
    level: 'inline',
    start(src) {
      return src.match(/<u>/)?.index
    },
    tokenizer(src) {
      const match = src.match(/^<u>(.*?)<\/u>/)

      if (!match) {
        return
      }

      return {
        type: 'u',
        raw: match[0],
        text: match[1]
      }
    },
    renderer(token) {
      return `<u>${token.text}</u>`
    }
  }
]
