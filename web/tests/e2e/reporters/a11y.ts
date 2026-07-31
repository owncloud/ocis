import type {
  Reporter,
  FullConfig,
  Suite,
  TestCase,
  TestResult,
  FullResult
} from '@playwright/test/reporter'
import { AxeResults } from 'axe-core'
import * as fs from 'fs/promises'
import * as path from 'path'

interface A11yTestResult {
  test: string
  file: string
  line: number
  status: string
  duration: number
  url?: string
  violations: AxeResults['violations']
  violationCount: number
  passCount: number
  incompleteCount: number
}

interface A11yReport {
  summary: {
    totalTests: number
    totalViolations: number
    totalPasses: number
    totalIncomplete: number
    timestamp: string
    duration: number
  }
  tests: A11yTestResult[]
}

class A11yReporter implements Reporter {
  private results: A11yTestResult[] = []
  private outputFile: string
  private startTime: number = 0

  constructor(options: { outputFile?: string } = {}) {
    this.outputFile = options.outputFile || 'a11y-report.json'
  }

  onBegin(_config: FullConfig, _suite: Suite): void {
    this.startTime = Date.now()
  }

  onTestEnd(test: TestCase, result: TestResult): void {
    result.attachments.forEach((attachment) => {
      if (attachment.name !== 'accessibility-scan') {
        return
      }

      try {
        let axeResults: AxeResults

        if (attachment.body) {
          axeResults = JSON.parse(attachment.body.toString('utf-8'))
        } else if (attachment.path) {
          const content = require('fs').readFileSync(attachment.path, 'utf-8')
          axeResults = JSON.parse(content)
        } else {
          throw new Error('No accessibility scan attachment found')
        }

        this.results.push({
          test: test.title,
          file: path.relative(process.cwd(), test.location.file),
          line: test.location.line,
          status: result.status,
          duration: result.duration,
          url: axeResults.url,
          violations: axeResults.violations || [],
          violationCount: axeResults.violations?.length || 0,
          passCount: axeResults.passes?.length || 0,
          incompleteCount: axeResults.incomplete?.length || 0
        })
      } catch (error) {
        console.error(`Error parsing accessibility results for test "${test.title}":`, error)
      }
    })
  }

  async onEnd(_result: FullResult): Promise<void> {
    const duration = Date.now() - this.startTime
    const totalViolations = this.results.reduce((sum, test) => sum + test.violationCount, 0)

    const report: A11yReport = {
      summary: {
        totalTests: this.results.length,
        totalViolations,
        totalPasses: this.results.reduce((sum, test) => sum + test.passCount, 0),
        totalIncomplete: this.results.reduce((sum, test) => sum + test.incompleteCount, 0),
        timestamp: new Date().toISOString(),
        duration
      },
      tests: this.results
    }

    try {
      const outputDir = path.dirname(this.outputFile)
      await fs.mkdir(outputDir, { recursive: true })

      await fs.writeFile(this.outputFile, JSON.stringify(report, null, 2), 'utf-8')

      console.info(`\n📊 Accessibility Report Generated:`)
      console.info(`   File: ${this.outputFile}`)
      console.info(`   Tests: ${report.summary.totalTests}`)
      console.warn(`   Violations: ${report.summary.totalViolations}`)

      if (totalViolations > 0) {
        console.warn(`\n⚠️  Found ${totalViolations} accessibility violations`)
      } else {
        console.info(`\n✅ No accessibility violations found`)
      }
    } catch (error) {
      console.error('Error writing accessibility report:', error)
    }
  }
}

export default A11yReporter
