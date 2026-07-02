import { markedExtensions } from '../../../../../src/composables/webWorkers/exportAsPdfWorker/markedExtensions'

describe('markedExtensions', () => {
  describe('superscript extension', () => {
    const superscriptExtension = markedExtensions.find((ext) => ext.name === 'superscript') as any

    describe('tokenizer function', () => {
      it('should tokenize valid superscript syntax', () => {
        const result = superscriptExtension.tokenizer('^hello^')
        expect(result).toEqual({
          type: 'sup',
          raw: '^hello^',
          text: 'hello'
        })
      })

      it('should tokenize single character superscript', () => {
        const result = superscriptExtension.tokenizer('^x^')
        expect(result).toEqual({
          type: 'sup',
          raw: '^x^',
          text: 'x'
        })
      })

      it('should tokenize numbers in superscript', () => {
        const result = superscriptExtension.tokenizer('^123^')
        expect(result).toEqual({
          type: 'sup',
          raw: '^123^',
          text: '123'
        })
      })

      it('should not tokenize when spaces are present', () => {
        expect(superscriptExtension.tokenizer('^hello world^')).toBeUndefined()
        expect(superscriptExtension.tokenizer('^ hello^')).toBeUndefined()
        expect(superscriptExtension.tokenizer('^hello ^')).toBeUndefined()
      })

      it('should not tokenize incomplete syntax', () => {
        expect(superscriptExtension.tokenizer('^hello')).toBeUndefined()
        expect(superscriptExtension.tokenizer('hello^')).toBeUndefined()
        expect(superscriptExtension.tokenizer('^')).toBeUndefined()
      })

      it('should not tokenize empty content', () => {
        expect(superscriptExtension.tokenizer('^^')).toBeUndefined()
      })
    })
  })

  describe('subscript extension', () => {
    const subscriptExtension = markedExtensions.find((ext) => ext.name === 'subscript') as any

    describe('tokenizer function', () => {
      it('should tokenize valid subscript syntax', () => {
        const result = subscriptExtension.tokenizer('~hello~')
        expect(result).toEqual({
          type: 'sub',
          raw: '~hello~',
          text: 'hello'
        })
      })

      it('should tokenize single character subscript', () => {
        const result = subscriptExtension.tokenizer('~x~')
        expect(result).toEqual({
          type: 'sub',
          raw: '~x~',
          text: 'x'
        })
      })

      it('should tokenize numbers in subscript', () => {
        const result = subscriptExtension.tokenizer('~123~')
        expect(result).toEqual({
          type: 'sub',
          raw: '~123~',
          text: '123'
        })
      })

      it('should not tokenize when spaces are present', () => {
        expect(subscriptExtension.tokenizer('~hello world~')).toBeUndefined()
        expect(subscriptExtension.tokenizer('~ hello~')).toBeUndefined()
        expect(subscriptExtension.tokenizer('~hello ~')).toBeUndefined()
      })

      it('should not tokenize incomplete syntax', () => {
        expect(subscriptExtension.tokenizer('~hello')).toBeUndefined()
        expect(subscriptExtension.tokenizer('hello~')).toBeUndefined()
        expect(subscriptExtension.tokenizer('~')).toBeUndefined()
      })

      it('should not tokenize empty content', () => {
        expect(subscriptExtension.tokenizer('~~')).toBeUndefined()
      })
    })
  })

  describe('underline extension', () => {
    const underlineExtension = markedExtensions.find((ext) => ext.name === 'underline') as any

    describe('tokenizer function', () => {
      it('should tokenize valid underline syntax', () => {
        const result = underlineExtension.tokenizer('<u>hello</u>')
        expect(result).toEqual({
          type: 'u',
          raw: '<u>hello</u>',
          text: 'hello'
        })
      })

      it('should tokenize single character underline', () => {
        const result = underlineExtension.tokenizer('<u>x</u>')
        expect(result).toEqual({
          type: 'u',
          raw: '<u>x</u>',
          text: 'x'
        })
      })

      it('should tokenize text with spaces', () => {
        const result = underlineExtension.tokenizer('<u>hello world</u>')
        expect(result).toEqual({
          type: 'u',
          raw: '<u>hello world</u>',
          text: 'hello world'
        })
      })

      it('should tokenize numbers in underline', () => {
        const result = underlineExtension.tokenizer('<u>123</u>')
        expect(result).toEqual({
          type: 'u',
          raw: '<u>123</u>',
          text: '123'
        })
      })

      it('should tokenize empty content', () => {
        const result = underlineExtension.tokenizer('<u></u>')
        expect(result).toEqual({
          type: 'u',
          raw: '<u></u>',
          text: ''
        })
      })

      it('should not tokenize incomplete syntax', () => {
        expect(underlineExtension.tokenizer('<u>hello')).toBeUndefined()
        expect(underlineExtension.tokenizer('hello</u>')).toBeUndefined()
        expect(underlineExtension.tokenizer('<u>')).toBeUndefined()
      })

      it('should handle special characters in content', () => {
        const result = underlineExtension.tokenizer('<u>hello@#$%</u>')
        expect(result).toEqual({
          type: 'u',
          raw: '<u>hello@#$%</u>',
          text: 'hello@#$%'
        })
      })
    })
  })
})
