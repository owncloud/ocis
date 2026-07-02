import { Resource } from '@ownclouders/web-client'
import { filterResources } from '../../../../src/helpers/resource'

describe('filterResources', () => {
  it('filters given resources by given term', () => {
    const resultset = filterResources(
      [{ name: 'foo' } as Resource, { name: 'bar' } as Resource],
      'foo'
    )

    expect(resultset).toMatchObject([{ name: 'foo' }])
    expect(resultset.length).toBe(1)
  })

  it('can limit the resultset', () => {
    const filter = filterResources(
      [{ name: 'foo' } as Resource, { name: 'foo' } as Resource],
      'foo',
      1
    )

    expect(filter).toMatchObject([{ name: 'foo' }])
    expect(filter.length).toBe(1)
  })
})
