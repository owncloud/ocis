import type { Resource } from '@ownclouders/web-client'

interface ReaderLocation {
  atEnd?: boolean
  atStart?: boolean
  start?: {
    cfi?: string
    displayed?: { page?: number; total?: number }
    index?: number
    percentage?: number
  }
}

export interface ReaderProgress {
  finished: boolean
  hasPosition: boolean
  percentage?: number
}

function fallbackPercentage(location: ReaderLocation, spineItemCount: number): number | undefined {
  const { displayed, index } = location.start ?? {}
  if (!spineItemCount || typeof index !== 'number') return undefined

  const page = displayed?.page ?? 1
  const totalPages = displayed?.total ?? 1
  const chapterProgress = totalPages > 0 ? (page - 1) / totalPages : 0
  const percentage = Math.round(((index + chapterProgress) / spineItemCount) * 100)

  if (location.atEnd) return 100
  if (location.atStart) return 0
  return Math.max(1, Math.min(99, percentage))
}

export function getReaderProgress(resource: Resource, spineItemCount = 0): ReaderProgress {
  if (typeof localStorage === 'undefined') return { finished: false, hasPosition: false }

  try {
    const storedValue = localStorage.getItem(`oc_epubReader_resource_${resource.id}`)
    if (!storedValue) return { finished: false, hasPosition: false }

    const location = (JSON.parse(storedValue) as { currentLocation?: ReaderLocation })
      .currentLocation
    if (!location?.start?.cfi) return { finished: false, hasPosition: false }

    const savedPercentage = location.start.percentage
    const percentage =
      typeof savedPercentage === 'number' && savedPercentage > 0
        ? Math.round(savedPercentage * 100)
        : fallbackPercentage(location, spineItemCount)
    return {
      finished: location.atEnd === true || percentage === 100,
      hasPosition: true,
      percentage:
        typeof percentage === 'number' ? Math.max(0, Math.min(100, percentage)) : undefined
    }
  } catch {
    return { finished: false, hasPosition: false }
  }
}
