import simpleGit from 'simple-git'
import path from 'path'
import { fileURLToPath } from 'url'
import { dirname } from 'path'
import { Command } from 'commander'
import fs from 'fs'

const __filename = fileURLToPath(import.meta.url)
const __dirname = dirname(__filename)

const git = simpleGit()
const program = new Command()

program
  .option('-d, --days <number>', 'number of days to look back')
  .option('-m, --month <number>', 'month to look back')
  .option('-y, --year <number>', 'year to look back')
  .parse(process.argv)

const options = program.opts()

function getGitDateFormat(date) {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

function getStartAndEndDate(month, year) {
  const startDate = new Date(year, month - 1, 1)
  const endDate = new Date(year, month, 0)
  return { startDate, endDate }
}

async function getRecentChanges() {
  try {
    const repoPath = path.resolve(__dirname)
    await git.cwd(repoPath)

    let formattedSinceDate, formattedUntilDate, period

    if (options.month && options.year) {
      const { startDate, endDate } = getStartAndEndDate(options.month, options.year)
      formattedSinceDate = getGitDateFormat(startDate)
      formattedUntilDate = getGitDateFormat(endDate)
      period = `${options.month}_${options.year}`
    } else if (options.days) {
      const today = new Date()
      today.setDate(today.getDate() - options.days)
      formattedSinceDate = getGitDateFormat(today)
      formattedUntilDate = getGitDateFormat(new Date())
      period = `Last_${options.days}_days`
    } else {
      console.error('Please provide either days or month and year.')
      return
    }

    const logOptions = {
      '--since': formattedSinceDate
    }
    if (formattedUntilDate) {
      logOptions['--until'] = formattedUntilDate
    }

    const logs = await git.log(logOptions)

    const csvRows = [
      ['Test-Type', 'Date', 'Tests Added', 'Tests Changed', 'Tests Deleted', 'commit-ID']
    ]

    for (const log of logs.all) {
      const e2eDiff = await git.diff([
        `${log.hash}~1`,
        log.hash,
        '--',
        'tests/e2e/cucumber/features'
      ])
      const unitDiff = await git.diff([`${log.hash}~1`, log.hash, '--', 'packages/**/*.spec.ts'])

      const e2eRow = analyzeE2eDiff(e2eDiff, log)
      const unitRow = analyzeUnitDiff(unitDiff, log)

      if (e2eRow) {
        csvRows.push(e2eRow)
      }

      if (unitRow) {
        csvRows.push(unitRow)
      }
    }

    const csvContent = csvRows.map((row) => row.join(',')).join('\n')
    const reportsDir = path.join(__dirname, 'reports')

    if (!fs.existsSync(reportsDir)) {
      fs.mkdirSync(reportsDir)
    }
    const reportFilePath = path.join(reportsDir, `QA_Activity_Report_${period}.csv`)

    fs.writeFile(reportFilePath, csvContent, (err) => {
      if (err) {
        console.error('Error writing CSV report:', err)
      } else {
        console.log(`CSV report generated successfully. You can find it in , ${reportFilePath}`)
      }
    })
  } catch (error) {
    console.error('Error:', error)
  }
}

function analyzeE2eDiff(diff, log) {
  const diffLines = diff.split('\n')

  let commitAddedTests = 0
  let commitChangedTests = 0
  let commitDeletedTests = 0
  let currentFile = null

  for (const line of diffLines) {
    if (line.startsWith('diff --git')) {
      const match = line.match(/ b\/(tests\/e2e\/cucumber\/features\/[^\s]+)/)
      if (match) {
        currentFile = match[1]
      }
    } else if (line.startsWith('+') && !line.startsWith('+++')) {
      // Consider only the addition of scenarios or features. Example: +  Scenario: activity
      if (line.includes('Scenario:')) {
        commitAddedTests++
      }
    } else if (line.startsWith('-') && !line.startsWith('---')) {
      // Consider only the deleting of scenarios or features. Example: -  Scenario: activity
      if (line.includes('Scenario:')) {
        commitDeletedTests++
      }
    } else if (line.includes('@@ Feature:') && currentFile) {
      // if line contains 'Feature', that is test change. Example @@ -17,8 +17,8 @@ Feature: Download
      commitChangedTests++
    }
  }
  if (commitAddedTests || commitChangedTests || commitDeletedTests) {
    return ['UI Test', log.date, commitAddedTests, commitChangedTests, commitDeletedTests, log.hash]
  }
}

function analyzeUnitDiff(diff, log) {
  const diffLines = diff.split('\n')

  let commitAddedTests = 0
  let commitChangedTests = 0
  let commitDeletedTests = 0
  let currentFile = null
  let inChangeBlock = false
  let inTest = false

  for (let i = 0; i < diffLines.length; i++) {
    const line = diffLines[i]

    // Detect the file being diffed
    if (line.startsWith('diff --git')) {
      const match = line.match(/ b\/(packages\/[^\s]+\.spec\.ts)/)
      if (match) {
        currentFile = match[1]
      }
    }

    // Start a new change block
    if (line.startsWith('@@')) {
      if (inChangeBlock && !inTest) {
        commitChangedTests++
      }
      inChangeBlock = true
      inTest = false
    }

    // Process changes in the current block
    if (inChangeBlock) {
      if (line.startsWith('+') && !line.startsWith('+++')) {
        if (line.includes('it(') || line.includes('it.each(')) {
          commitAddedTests++
          inTest = true
        }
      } else if (line.startsWith('-') && !line.startsWith('---')) {
        if (line.includes('it(') || line.includes('it.each(')) {
          commitDeletedTests++
          inTest = true
        }
      }
    }

    // End of a change block
    if (line === '' || i === diffLines.length - 1) {
      if (inChangeBlock && !inTest) {
        commitChangedTests++
      }
      inChangeBlock = false
    }
  }

  // Return results if there are any changes
  if (commitAddedTests || commitChangedTests || commitDeletedTests) {
    return [
      'Unit Test',
      log.date,
      commitAddedTests,
      commitChangedTests,
      commitDeletedTests,
      log.hash
    ]
  }
}

getRecentChanges()
