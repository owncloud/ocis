import { ShareType, ShareTypes } from '../../../../src/helpers/share'

describe('ShareTypes', () => {
  describe('getByValue', () => {
    it.each([...ShareTypes.all])(
      'finds correct ShareType by numeric value',
      (shareType: ShareType) => {
        expect(ShareTypes.getByValue(shareType.value)).toBe(shareType)
      }
    )
    it.each([-1, 5, 999])(
      'returns undefined when a numeric value has no ShareType representation',
      (value: number) => {
        expect(ShareTypes.getByValue(value)).toBeUndefined()
      }
    )
  })

  describe('getValues', () => {
    it.each([
      [
        'empty types list',
        {
          types: [],
          values: []
        }
      ],
      [
        'some types',
        {
          types: [ShareTypes.guest, ShareTypes.group],
          values: [ShareTypes.guest.value, ShareTypes.group.value]
        }
      ]
    ])('with %s', (name: string, { types, values }) => {
      expect(ShareTypes.getValues(types)).toEqual(values)
    })
  })

  describe('containsAnyValue', () => {
    it.each([
      [
        'given empty types and empty values',
        {
          types: [],
          values: [],
          result: false
        }
      ],
      [
        'given empty types and some values',
        {
          types: [],
          values: [1, 2, 3],
          result: false
        }
      ],
      [
        'given some types and empty values',
        {
          types: [ShareTypes.user, ShareTypes.group],
          values: [],
          result: false
        }
      ],
      [
        'given some types and some values without intersection',
        {
          types: [ShareTypes.user, ShareTypes.group],
          values: [ShareTypes.guest.value, ShareTypes.link.value, ShareTypes.remote.value],
          result: false
        }
      ],
      [
        'given some types and some values with partial match',
        {
          types: [ShareTypes.user, ShareTypes.group],
          values: [ShareTypes.guest.value, ShareTypes.group.value],
          result: true
        }
      ],
      [
        'given some types and some values with full match',
        {
          types: [ShareTypes.user, ShareTypes.group],
          values: [ShareTypes.user.value, ShareTypes.group.value],
          result: true
        }
      ]
    ])('%s', (name, { types, values, result }) => {
      expect(ShareTypes.containsAnyValue(types, values)).toBe(result)
    })
  })
})
