import { Resource } from '@ownclouders/web-client'
import { defaultPlugins, shallowMount } from '@ownclouders/web-test-helpers'
import AudioMetaPanel from '../../../../../src/components/SideBar/Audio/AudioMetaPanel.vue'
import { mock } from 'vitest-mock-extended'
import { Audio } from '@ownclouders/web-client/graph/generated'

describe('AudioMeta SideBar Panel', () => {
  const keys = ['title', 'duration', 'artist', 'album', 'genre', 'year', 'track', 'disc']
  it.each(keys)('shows value in panel for key "%s"', (key) => {
    const resource = mock<Resource>({
      audio: mock<Audio>({
        title: 'Jingle Bells',
        duration: 510000,
        artist: 'Horst Horstsen',
        album: 'Christmas Gedöns',
        genre: 'Christmas',
        year: 1956,
        track: 23,
        trackCount: 42,
        disc: 4,
        discCount: 12
      })
    })
    const expectedValues = {
      title: resource.audio.title,
      duration: '08:30',
      artist: resource.audio.artist,
      album: resource.audio.album,
      genre: resource.audio.genre,
      year: resource.audio.year.toString(),
      track: `${resource.audio.track} / ${resource.audio.trackCount}`,
      disc: `${resource.audio.disc} / ${resource.audio.discCount}`
    }
    const { wrapper } = createWrapper({ resource })
    expect(wrapper.find(`[data-testid="audio-panel-${key}"]`).text()).toBe(expectedValues[key])
  })
  it.each(keys)('shows "-" in panel if key "%s" has no value in provided data', (key) => {
    const emptyPhotoResourceMock = mock<Resource>({})
    const { wrapper } = createWrapper({ resource: emptyPhotoResourceMock })
    expect(wrapper.find(`[data-testid="audio-panel-${key}"]`).text()).toEqual('-')
  })
  it('shows "track" without "trackCount" if absent', () => {
    const resource = mock<Resource>({
      audio: mock<Audio>({
        duration: undefined,
        track: 23,
        trackCount: undefined
      })
    })
    const { wrapper } = createWrapper({ resource })
    expect(wrapper.find('[data-testid="audio-panel-track"]').text()).toBe(
      resource.audio.track.toString()
    )
  })
  it('shows "disc" without "discCount" if absent', () => {
    const resource = mock<Resource>({
      audio: mock<Audio>({
        duration: undefined,
        disc: 500,
        discCount: undefined
      })
    })
    const { wrapper } = createWrapper({ resource })
    expect(wrapper.find('[data-testid="audio-panel-disc"]').text()).toBe(
      resource.audio.disc.toString()
    )
  })
})

function createWrapper({ resource }: { resource: Resource }) {
  return {
    wrapper: shallowMount(AudioMetaPanel, {
      global: {
        plugins: [...defaultPlugins({})],
        provide: {
          resource
        }
      }
    })
  }
}
