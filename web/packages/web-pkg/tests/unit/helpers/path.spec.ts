import { getParentPaths } from '../../../src/helpers/path'

describe('build an array of parent paths from a provided path', () => {
  it('should return an empty array on an empty path', () => {
    const paths = getParentPaths('', false)
    expect(paths).toHaveLength(0)
  })

  it('should return an empty array on "/" as path', () => {
    const paths = getParentPaths('/', false)
    expect(paths).toHaveLength(0)
  })

  it('should prepend resulting paths with a "/" if none was given', () => {
    const paths = getParentPaths('a/b/c', false)
    expect(paths).toEqual(['/a/b', '/a'])
  })

  it('should make no difference between "a/b/c" and "/a/b/c" with includeCurrent=false', () => {
    const paths1 = getParentPaths('a/b/c', false)
    const paths2 = getParentPaths('/a/b/c', false)
    expect(paths1).toEqual(paths2)
  })

  it('should make no difference between "a/b/c" and "/a/b/c" with includeCurrent=true', () => {
    const paths1 = getParentPaths('a/b/c', true)
    const paths2 = getParentPaths('/a/b/c', true)
    expect(paths1).toEqual(paths2)
  })

  it('should have different results for different values of includeCurrent', () => {
    const paths1 = getParentPaths('a/b/c', true)
    const paths2 = getParentPaths('a/b/c', false)
    expect(paths1).not.toEqual(paths2)
  })

  it('should not interpret a trailing slash as yet another path segment', () => {
    const paths = getParentPaths('/a/b/c/', true)
    expect(paths).toEqual(['/a/b/c', '/a/b', '/a'])
  })

  it('should include the provided path in the result if includeCurrent=true', () => {
    const paths = getParentPaths('a/b/c', true)
    expect(paths).toEqual(['/a/b/c', '/a/b', '/a'])
  })
})
