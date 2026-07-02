import { Browser, BrowserContextOptions, LaunchOptions } from '@playwright/test'
import path from 'path'

import { config } from '../../../config'

export interface ActorsOptions {
  browser: Browser
  context: {
    acceptDownloads: boolean
    reportDir: string
    tracingReportDir: string
    reportVideo: boolean
    reportHar: boolean
    reportTracing: boolean
    failOnUncaughtConsoleError: boolean
  }
}

export interface ActorOptions extends ActorsOptions {
  id: string
  namespace: string
}

export const buildBrowserContextOptions = (options: ActorOptions): BrowserContextOptions => {
  const permissions = []
  // clipboard permissions are only available in chromium and chrome
  // https://github.com/microsoft/playwright/issues/13037
  if (['chromium', 'chrome'].includes(config.browser)) {
    permissions.push('clipboard-read', 'clipboard-write')
  }

  const contextOptions: BrowserContextOptions = {
    acceptDownloads: options.context.acceptDownloads,
    permissions: permissions,
    ignoreHTTPSErrors: true,
    locale: 'en-US'
  }

  if (options.context.reportVideo) {
    contextOptions.recordVideo = {
      dir: path.join(options.context.reportDir, 'playwright', 'video')
    }
  }

  if (options.context.reportHar) {
    contextOptions.recordHar = {
      path: path.join(options.context.reportDir, 'playwright', 'har', `${options.namespace}.har`)
    }
  }

  return contextOptions
}

export const getBrowserLaunchOptions = (): LaunchOptions => {
  const args = []
  if (config.browser !== 'webkit') {
    args.push('--use-fake-ui-for-media-stream', '--use-fake-device-for-media-stream')
  }

  return {
    slowMo: config.slowMo,
    args,
    firefoxUserPrefs: {
      'media.navigator.streams.fake': true,
      'media.navigator.permission.disabled': true
    },
    headless: config.headless
  }
}
