import { test, request } from '@playwright/test'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'
import { login } from './support/oc'
import { seed } from './support/seed'
import { tours } from './tours'
import { config } from '../config'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const outputRoot = path.join(__dirname, 'output')
const baseURL = config.baseUrl

const shotName = (index: number, shot: string) =>
  `${String(index + 1).padStart(2, '0')}-${shot}.png`

// Seed best-effort demo data once before capturing.
test.beforeAll(async () => {
  const api = await request.newContext({ baseURL, ignoreHTTPSErrors: true })
  try {
    await seed(api)
  } finally {
    await api.dispose()
  }
})

// Capture each tour: drive the UI step by step and screenshot the resulting state.
for (const tour of tours) {
  test(`capture: ${tour.id}`, async ({ page }) => {
    await login(page)
    const dir = path.join(outputRoot, tour.id)
    fs.mkdirSync(dir, { recursive: true })
    for (let i = 0; i < tour.steps.length; i++) {
      const step = tour.steps[i]
      await step.run(page)
      await page.waitForTimeout(300)
      await page.screenshot({ path: path.join(dir, shotName(i, step.shot)) })
    }
  })
}

// Write a manifest describing the fully-captured tours and their captions, so a
// documentation site (or any consumer) can pick up the screenshots and text.
test.afterAll(() => {
  const captured = tours.filter((tour) =>
    tour.steps.every((_, i) =>
      fs.existsSync(path.join(outputRoot, tour.id, shotName(i, tour.steps[i].shot)))
    )
  )
  const manifest = {
    source: 'doc.owncloud.com/webui/latest/owncloud_web/web_for_users.html',
    tours: captured.map((tour) => ({
      id: tour.id,
      category: tour.category,
      title: tour.title,
      summary: tour.summary,
      steps: tour.steps.map((step, i) => ({
        title: step.title,
        caption: step.caption,
        screenshot: `${tour.id}/${shotName(i, step.shot)}`
      }))
    }))
  }
  fs.mkdirSync(outputRoot, { recursive: true })
  fs.writeFileSync(path.join(outputRoot, 'manifest.json'), JSON.stringify(manifest, null, 2) + '\n')
  console.log(
    `[docs-capture] captured ${captured.length}/${tours.length} tour(s) -> ${outputRoot} (+ manifest.json)`
  )
})
