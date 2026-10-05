import {
  isLeavingVaultAfterFailedStepUp,
  vaultStepUpFailedKey
} from '../../../src/helpers/vaultStepUp'

describe('isLeavingVaultAfterFailedStepUp', () => {
  beforeEach(() => {
    sessionStorage.clear()
  })

  it('is false without a failed step-up', () => {
    expect(isLeavingVaultAfterFailedStepUp()).toBe(false)
  })

  it('is true while a failed step-up is pending', () => {
    sessionStorage.setItem(vaultStepUpFailedKey, 'true')
    expect(isLeavingVaultAfterFailedStepUp()).toBe(true)
  })

  it('is false when sessionStorage is not accessible', () => {
    vi.spyOn(Storage.prototype, 'getItem').mockImplementation(() => {
      throw new Error('denied')
    })
    expect(isLeavingVaultAfterFailedStepUp()).toBe(false)
    vi.restoreAllMocks()
  })
})
