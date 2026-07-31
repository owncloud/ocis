import { shallowMount } from '@ownclouders/web-test-helpers'
import { AVAILABLE_SIZES } from '@ownclouders/design-system/helpers'
import ResourceIcon from '../../../../src/components/FilesList/ResourceIcon.vue'
import {
  ResourceIconMapping,
  resourceIconMappingInjectionKey
} from '../../../../src/helpers/resource'
import { Resource } from '@ownclouders/web-client'

const resourceIconMapping: ResourceIconMapping = {
  extension: {
    'not-a-real-extension': {
      name: 'resource-type-madeup-extension',
      color: 'red'
    }
  },
  mimeType: {
    'not-a-real-mimetype': {
      name: 'resource-type-file',
      color: 'var(--oc-color-text-default)'
    }
  }
}

describe('OcResourceIcon', () => {
  ;['file', 'folder', 'space'].forEach((type) => {
    match({
      type
    })
  })

  match(
    {
      type: 'file',
      extension: 'not-a-real-extension'
    },
    'with extension "not-a-real-extension"'
  )

  match(
    {
      type: 'file',
      mimeType: 'not-a-real-mimetype'
    },
    'with mimetype "not-a-real-mimetype"'
  )
})

type Size = 'xsmall' | 'small' | 'medium' | 'large' | 'xlarge' | 'xxlarge' | 'xxxlarge'

function match(resource: Partial<Resource>, additionalText?: string) {
  AVAILABLE_SIZES.forEach((size: Size) => {
    it(`renders OcIcon for resource type ${resource.type}${
      additionalText ? ` ${additionalText}` : ''
    } in size ${size}`, () => {
      const { wrapper } = getWrapper({ resource, size })
      expect(wrapper.html()).toMatchSnapshot()
    })
  })
}

function getWrapper({ resource, size }: { resource: Partial<Resource>; size: Size }) {
  return {
    wrapper: shallowMount(ResourceIcon, {
      global: {
        provide: {
          [resourceIconMappingInjectionKey]: resourceIconMapping
        }
      },
      props: {
        resource: resource as Resource,
        size
      }
    })
  }
}
