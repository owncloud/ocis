import { textUtils } from '../../../src/helpers/textUtils'

describe('textUtils', () => {
  // textUtils compares two strings to provide a natural sort
  describe('naturalSortCompare', () => {
    it('sorts a list naturally', () => {
      const actual = ['Brian Murphy', 'Alice Hansen', 'grp1', 'grp11']

      actual.sort(textUtils.naturalSortCompare)

      expect(actual.join(',')).toBe('Alice Hansen,Brian Murphy,grp1,grp11')
    })

    /**
     * variations
     *
     [
       { firstString: 'b', secondString: 'a' },
       { firstString: 'bb', secondString: 'ba' },
       { firstString: '1', secondString: '0' }
     ]
     */
    it.todo('should return negative integer if "b" comes before "a"')

    /**
     * variations
     *
     [
       { firstString: 'a', secondString: 'b' },
       { firstString: 'aa', secondString: 'ab' },
       { firstString: '0', secondString: '1' }
     ]
     */
    it.todo('should return positive integer if "a" comes before "b"')

    /**
     * variations
     *
     [
       { firstString: '0', secondString: '0' },
       { firstString: 'a', secondString: 'a' },
       { firstString: 'aa', secondString: 'aa' },
       { firstString: 'ab', secondString: 'ab' }
     ]
     */
    it.todo('should return 0 if the provided strings are identical')
  })
})
