import { z } from 'zod'

const ExternalAppSchema = z.object({
  icon: z.string(),
  name: z.string(),
  secure_view: z.boolean().optional().default(false),
  target_ext: z.string().optional()
})
export type ExternalApp = z.infer<typeof ExternalAppSchema>

const MimeTypeSchema = z.object({
  allow_creation: z.boolean().optional(),
  app_providers: z.array(ExternalAppSchema),
  default_application: z.string().optional(),
  description: z.string().optional(),
  ext: z.string().optional(),
  mime_type: z.string(),
  name: z.string().optional()
})
export type MimeType = z.infer<typeof MimeTypeSchema>

export const MimeTypesToAppsSchema = z.object({
  'mime-types': z.array(MimeTypeSchema)
})
