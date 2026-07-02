import { z } from 'zod'

export const AppStoreRepositorySchema = z.object({
  name: z.string(),
  url: z.string()
})
export type AppStoreRepository = z.infer<typeof AppStoreRepositorySchema>

export const AppStoreConfigSchema = z.object({
  repositories: z.array(AppStoreRepositorySchema)
})

export const AppVersionSchema = z.object({
  version: z.string(),
  minOCIS: z.string().optional(),
  url: z.string(),
  filename: z.string().optional()
})
export type AppVersion = z.infer<typeof AppVersionSchema>

export const BADGE_COLORS = ['primary', 'success', 'danger'] as const
export const AppBadgeSchema = z.object({
  label: z.string(),
  color: z.enum(BADGE_COLORS).optional().default('primary')
})
export type AppBadge = z.infer<typeof AppBadgeSchema>

export const AppAuthorSchema = z.object({
  name: z.string(),
  email: z.string().optional(),
  url: z.string().optional()
})
export type AppAuthor = z.infer<typeof AppAuthorSchema>

export const AppImageSchema = z.object({
  url: z.string(),
  caption: z.string().optional()
})
export type AppImage = z.infer<typeof AppImageSchema>

export const AppResourceSchema = z.object({
  url: z.string(),
  label: z.string(),
  icon: z.string().optional()
})
export type AppResource = z.infer<typeof AppResourceSchema>

export const RawAppSchema = z.object({
  id: z.string(),
  name: z.string(),
  subtitle: z.string(),
  badge: AppBadgeSchema.optional(),
  description: z.string().optional(),
  license: z.string(),
  versions: z.array(AppVersionSchema), // versions are expected to be sorted from newest to oldest
  authors: z.array(AppAuthorSchema),
  tags: z.array(z.string()),
  coverImage: AppImageSchema.optional(),
  screenshots: z.array(AppImageSchema).optional().default([]),
  resources: z.array(AppResourceSchema).optional().default([]) // e.g. documentation, github, etc.
})

export const AppSchema = RawAppSchema.extend({
  repository: AppStoreRepositorySchema,
  mostRecentVersion: AppVersionSchema
})
export type App = z.infer<typeof AppSchema>

export const RawAppListSchema = z.object({
  apps: z.array(RawAppSchema)
})
