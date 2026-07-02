import { z } from 'zod'

const CustomTranslationSchema = z.object({
  url: z.string()
})

export type CustomTranslation = z.infer<typeof CustomTranslationSchema>

const OAuth2ConfigSchema = z.object({
  apiUrl: z.string().optional(),
  authUrl: z.string().optional(),
  clientId: z.string().optional(),
  clientSecret: z.string().optional(),
  logoutUrl: z.string().optional(),
  metaDataUrl: z.string().optional(),
  url: z.string().optional()
})

export type OAuth2Config = z.infer<typeof OAuth2ConfigSchema>

const OpenIdConnectConfigSchema = z
  .object({
    authority: z.string().optional(),
    client_id: z.string().optional(),
    client_secret: z.string().optional(),
    dynamic: z.string().optional(),
    metadata_url: z.string().optional(),
    post_logout_redirect_uri: z.string().optional(),
    response_type: z.string().optional(),
    scope: z.string().optional()
  })
  .passthrough()

export type OpenIdConnectConfig = z.infer<typeof OpenIdConnectConfigSchema>

const SentryConfigSchema = z.record(z.string(), z.any())

export type SentryConfig = z.infer<typeof SentryConfigSchema>

const StyleConfigSchema = z.object({
  href: z.string().optional()
})

export type StyleConfig = z.infer<typeof StyleConfigSchema>

const ScriptConfigSchema = z.object({
  async: z.boolean().optional(),
  src: z.string().optional()
})

export type ScriptConfig = z.infer<typeof ScriptConfigSchema>

const OptionsConfigSchema = z.object({
  cernFeatures: z.boolean().optional(),
  concurrentRequests: z
    .object({
      resourceBatchActions: z.number().optional(),
      sse: z.number().optional(),
      shares: z
        .object({
          create: z.number().optional(),
          list: z.number().optional()
        })
        .optional()
    })
    .optional(),
  contextHelpers: z.boolean().optional(),
  contextHelpersReadMore: z.boolean().optional(),
  defaultExtension: z.string().optional(),
  disabledExtensions: z.array(z.string()).optional(),
  disableFeedbackLink: z.boolean().optional(),
  accountEditLink: z
    .object({
      href: z.string().optional()
    })
    .optional(),
  editor: z
    .object({
      autosaveEnabled: z.boolean().optional(),
      autosaveInterval: z.number().optional(),
      openAsPreview: z.union([z.boolean(), z.array(z.string())]).optional()
    })
    .optional(),
  embed: z
    .object({
      enabled: z.boolean().optional(),
      target: z.string().optional(),
      messagesOrigin: z.string().optional(),
      delegateAuthentication: z.boolean().optional(),
      delegateAuthenticationOrigin: z.string().optional(),
      fileTypes: z.array(z.string()).optional(),
      chooseFileName: z.boolean().optional(),
      chooseFileNameSuggestion: z.string().optional()
    })
    .optional(),
  feedbackLink: z
    .object({
      ariaLabel: z.string().optional(),
      description: z.string().optional(),
      href: z.string().optional()
    })
    .optional(),
  isRunningOnEos: z.boolean().optional(),
  loginUrl: z.string().optional(),
  logoutUrl: z.string().optional(),
  ocm: z
    .object({
      openRemotely: z.boolean().optional()
    })
    .optional(),
  routing: z
    .object({
      fullShareOwnerPaths: z.boolean().optional(),
      idBased: z.boolean().optional()
    })
    .optional(),
  runningOnEos: z.boolean().optional(),
  tokenStorageLocal: z.boolean().optional(),
  upload: z
    .object({
      companionUrl: z.string().optional()
    })
    .optional(),
  userListRequiresFilter: z.boolean().optional(),
  hideLogo: z.boolean().optional(),
  hideAppSwitcher: z.boolean().optional(),
  hideAccountMenu: z.boolean().optional(),
  hideNavigation: z.boolean().optional(),
  defaultLanguage: z.string().optional()
})

export type OptionsConfig = z.infer<typeof OptionsConfigSchema>

const ExternalApp = z.object({
  id: z.string(),
  path: z.string(),
  config: z.record(z.string(), z.unknown()).optional()
})

export const RawConfigSchema = z.object({
  server: z.string(),
  theme: z.string(),
  options: OptionsConfigSchema,
  apps: z.array(z.string()).optional(),
  external_apps: z.array(ExternalApp).optional(),
  customTranslations: z.array(CustomTranslationSchema).optional(),
  auth: OAuth2ConfigSchema.optional(),
  openIdConnect: OpenIdConnectConfigSchema.optional(),
  sentry: SentryConfigSchema.optional(),
  scripts: z.array(ScriptConfigSchema).optional(),
  styles: z.array(StyleConfigSchema).optional()
})

export type RawConfig = z.infer<typeof RawConfigSchema>
