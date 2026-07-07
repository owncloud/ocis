import { urlJoin } from '@ownclouders/web-client'

describe('proper-url-join', () => {
  it.each([
    [['http://foobar.com'], 'http://foobar.com'],
    [['http://foobar.com/'], 'http://foobar.com'],
    [['/', ''], '/'],
    [['/', 'foo'], '/foo'],
    [['/', 'foo/'], '/foo'],
    [['foo/'], '/foo'],
    [['/', undefined], '/'],
    [['', { leadingSlash: true }], '/'],
    [[undefined, { leadingSlash: true }], '/'],
    [['/', 2], '/2'],
    [['//', '/fol//der//', '//file'], '/fol/der/file'],
    [['?&@'], '/?&@']
  ])('joins %s as %s', (args: unknown, expected: string) => {
    // @ts-ignore
    expect(urlJoin(...args)).toBe(expected)
  })
})
