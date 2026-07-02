import { App } from 'vue'
import { applyCustomProp } from './helpers'

import * as components from './components'
import * as directives from './directives'

const initializeCustomProps = (tokens: string[] = [], prefix: string) => {
  for (const param in tokens) {
    applyCustomProp(prefix + param, tokens[param])
  }
}

export default {
  install(app: App, options: any = {}) {
    import('./utils/webFontLoader')

    const themeOptions = options.tokens
    initializeCustomProps(themeOptions?.breakpoints, 'breakpoint-')
    initializeCustomProps(themeOptions?.colorPalette, 'color-')
    initializeCustomProps(themeOptions?.fontSizes, 'font-size-')
    initializeCustomProps(themeOptions?.sizes, 'size-')
    initializeCustomProps(themeOptions?.spacing, 'space-')
    applyCustomProp('font-family', themeOptions?.fontFamily)

    Object.values(components).forEach((c) => app.component(c.name, c))
    Object.values(directives).forEach((d) => app.directive(d.name, d))
  }
}
