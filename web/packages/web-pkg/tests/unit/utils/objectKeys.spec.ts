import { objectKeys } from '../../../src/utils'

describe('objectKeys', () => {
  it('should return the correct keys', () => {
    expect(
      objectKeys({
        foo1: { bar1: { baz1: 1, baz2: 1 }, bar2: { baz1: 1, baz2: 1 }, bar3: 1 },
        foo2: 1
      })
    ).toMatchObject([
      'foo1.bar1.baz1',
      'foo1.bar1.baz2',
      'foo1.bar2.baz1',
      'foo1.bar2.baz2',
      'foo1.bar3',
      'foo2'
    ])
  })
})
