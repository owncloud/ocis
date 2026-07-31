import { Resource, SpaceResource } from '@ownclouders/web-client'
import { z } from 'zod'

export const ProcessorType = z.enum(['fit', 'resize', 'fill', 'thumbnail'])

export type ProcessorType = z.infer<typeof ProcessorType>

export interface BuildQueryStringOptions {
  preview?: number
  scalingup?: number
  a?: number
  etag?: string
  dimensions?: [number, number]
  processor?: ProcessorType
}

export interface LoadPreviewOptions {
  space: SpaceResource
  resource: Resource
  dimensions?: [number, number]
  processor?: ProcessorType
}
