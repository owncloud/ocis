import { AxiosInstance } from 'axios'
import get from 'lodash-es/get'

export interface AppProviderCapability {
  apps_url?: string
  enabled?: boolean
  new_url?: string
  open_url?: string
  version?: string
}

export interface PasswordPolicyCapability {
  min_characters?: number
  max_characters?: number
  min_lowercase_characters?: number
  min_uppercase_characters?: number
  min_digits?: number
  min_special_characters?: number
}

export interface PasswordEnforcedForCapability {
  read_only?: boolean
  read_write?: boolean
  upload_only?: boolean
  read_write_delete?: boolean
}

export interface SearchPropertyCapability {
  enabled?: boolean
}
export interface LastModifiedFilterCapability extends SearchPropertyCapability {
  keywords?: string[]
}

export interface MediaTypeCapability extends SearchPropertyCapability {
  keywords?: string[]
}

/**
 * Archiver struct within the capabilities as defined in reva
 * @see https://github.com/cs3org/reva/blob/41d5a6858c2200a61736d2c165e551b9785000d1/internal/http/services/owncloud/ocs/data/capabilities.go#L105
 */
export interface ArchiverCapability {
  enabled?: boolean
  version?: string // version is just a major version, e.g. `v2`
  formats?: string[]

  archiver_url?: string
  max_num_files?: string
  max_size?: string
}

interface AuthCapability {
  mfa: {
    enabled?: boolean
    levelnames?: string[]
    session_duration?: number
  }
}

export interface Capabilities {
  capabilities: {
    auth: AuthCapability
    checksums?: {
      preferredUploadType?: string
      supportedTypes?: string[]
    }
    password_policy?: PasswordPolicyCapability
    search?: {
      property?: {
        content?: SearchPropertyCapability
        mediatype?: MediaTypeCapability
        mtime?: LastModifiedFilterCapability
        name?: MediaTypeCapability
        scope?: MediaTypeCapability
        size?: MediaTypeCapability
        tag?: MediaTypeCapability
        tags?: MediaTypeCapability
        type?: MediaTypeCapability
      }
    }
    notifications?: {
      'ocs-endpoints'?: string[]
      configurable?: boolean
    }
    core: {
      pollinterval?: number
      status?: {
        edition?: string
        installed?: boolean
        maintenance?: boolean
        needsDbUpgrade?: boolean
        product?: string
        productname?: string
        productversion?: string
        version?: string
        versionstring?: string
      }
      'support-sse'?: boolean
      'support-url-signing'?: boolean
      'webdav-root'?: string
    }
    dav: {
      chunking?: string
      chunkingParallelUploadDisabled?: boolean
      reports?: string[]
      trashbin?: string
    }
    files: {
      app_providers?: AppProviderCapability[]
      archivers?: ArchiverCapability[]
      favorites?: boolean
      full_text_search?: boolean
      permanent_deletion?: boolean
      privateLinks?: boolean
      tags?: boolean
      tus_support?: {
        extension?: string
        http_method_override?: boolean
        max_chunk_size?: number
        resumable?: string
        version?: string
      }
      undelete?: boolean
      versioning?: boolean
      thumbnail?: {
        enabled?: boolean
        version?: string
        supportedMimeTypes?: string[]
      }
    }
    files_sharing: {
      allow_custom?: boolean
      api_enabled?: boolean
      can_rename?: boolean
      default_permissions?: number
      deny_access?: boolean
      federation?: {
        incoming?: boolean
        outgoing?: boolean
      }
      group_sharing?: boolean
      public?: {
        alias?: boolean
        can_contribute?: boolean
        can_edit?: boolean
        default_permissions?: number
        enabled?: boolean
        multiple?: boolean
        password?: {
          enforced?: boolean
          enforced_for?: PasswordEnforcedForCapability
        }
        send_mail?: boolean
        supports_upload_only?: boolean
        upload?: boolean
      }
      search_min_length?: number
      user?: {
        profile_picture?: boolean
        send_mail?: boolean
        settings?: {
          enabled?: boolean
          version?: string
        }[]
      }
      quick_link?: {
        default_role?: string
      }
    }
    spaces?: {
      enabled?: boolean
      max_quota?: number
      projects?: boolean
      version?: string
      server_managed?: boolean
    }
    vault?: {
      enabled?: boolean
      vault_storage_provider?: string
    }
    graph?: {
      'personal-data-export'?: boolean
      users: {
        change_password_self_disabled?: boolean
        create_disabled?: boolean
        delete_disabled?: boolean
        read_only_attributes?: string[]
      }
      tags: {
        max_tag_length: number
      }
    }
  }
  version: {
    edition?: string
    major?: string
    minor?: string
    micro?: string
    product?: string
    productversion?: string
    string?: string
  }
}

export const GetCapabilitiesFactory = (baseURI: string, axios: AxiosInstance) => {
  const url = new URL(baseURI)
  url.pathname = [...url.pathname.split('/'), 'cloud', 'capabilities'].filter(Boolean).join('/')
  url.searchParams.append('format', 'json')
  const endpoint = url.href
  return {
    async getCapabilities(): Promise<Capabilities> {
      const response = await axios.get(endpoint)
      return get(response, 'data.ocs.data', { capabilities: null, version: null })
    }
  }
}
