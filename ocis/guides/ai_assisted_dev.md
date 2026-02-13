---
title: "Building oCIS Web Extensions with AI-Assisted Development"
date: 2026-01-21T16:00:00+02:00
weight: 9
geekdocRepo: https://github.com/owncloud/ocis
geekdocEditPath: edit/master/docs/ocis/guides
geekdocFilePath: ai_assisted_dev.md
geekdocCollapseSection: true
---

{{< toc >}}

## Introduction

This guide demonstrates how to build Infinite Scale (oCIS) web extensions using AI-assisted development (sometimes called "vibe coding") with Claude AI. By following this approach, review-ready extensions can be created without writing code manually—the AI handles syntax, boilerplate, and implementation while the developer focuses on architecture and user experience.

**What you'll learn:**
- Setting up MCP connectors so Claude can access your server
- Using Claude to install development tools and manage your environment
- The five-phase workflow for AI-assisted extension development
- Practical tips for supervising AI-generated code
- How to contribute extensions back to the oCIS community

**Time investment:** ~10-14 hours per extension (concept to review-ready)

**Cost:** ~$100 USD for two extensions (using the Claude Max subscription)

## Prerequisites

Before starting, ensure you have the following:

- A running oCIS instance, see the [Admin Docs](https://doc.owncloud.com/ocis/latest/depl-examples/ubuntu-compose/ubuntu-compose-prod.html) for a deployment example
- [Claude Max subscription](https://claude.com/pricing/max) (starts with $100 US/month at the time of writing, other plans availabe)
- SSH access to your server (for MCP connectivity)
- Full root access to the server unless Claude Code, Claude MCP connector, oCIS and other tools all run as another user

That's it. Claude will help install everything else.

## Step 1: Configure MCP Connector

{{< hint type=important title="This Is the Hardest Part" >}}
Setting up the MCP connector is the most challenging step. Once this is done, Claude handles everything else—including installing development tools, managing git, and building extensions.
{{< /hint >}}

MCP (Model Context Protocol) allows Claude to connect directly to a server, enabling both the Claude web interface and Claude Code to access files, run commands, and debug oCIS deployments.

### Installation

Follow the official MCP documentation to set up a connector:

- **[MCP Quickstart Guide](https://modelcontextprotocol.io/quickstart)**
- **[Connect to Local Servers](https://modelcontextprotocol.io/docs/develop/connect-local-servers)**
- **[Claude Desktop MCP Setup](https://support.claude.com/en/articles/10949351-getting-started-with-local-mcp-servers-on-claude-desktop)**

### Verify Connection

Once configured, verify Claude can access the server by asking:

```text
Can you list the contents of /home on my server?
```

If Claude successfully lists directories, proceed to the next step.

## Step 2: Let Claude Set Up Your Environment

With MCP connected, Claude can install the development tools needed. Simply ask:

```text
Please install Claude Code and Git on my server. I'll need these for 
developing an oCIS web extension.
```

Claude will:
- Detect the operating system
- Install Claude Code using the appropriate method
- Install Git if not already present
- Configure any necessary paths or permissions

## The Five-Phase Workflow

This workflow has been refined over several projects and consistently produces review-ready code.

### Phase 1: Architecture (~15-30 minutes)

Start in the Claude web interface (not Claude Code). Describe the project and have a conversation:

```text
I want to build an oCIS web extension that displays photos in a timeline 
grouped by their EXIF capture date. It should have infinite scroll starting 
from today and going backwards, show camera metadata when you click on a photo, 
and include an interactive map view showing where photos were taken.

What questions do you have about the requirements?
```

Claude will ask clarifying questions about scope, user experience, and technical constraints. This back-and-forth shapes the architecture before any code is written.

### Phase 2: Deep Research (~30 minutes)

Ask Claude to research the existing oCIS architecture:

```text
Please research how oCIS web extensions work. I need to understand:
- The extension manifest format
- How extensions integrate with the oCIS web UI
- What APIs are available (Graph API, WebDAV, etc.)
- How existing extensions like json-viewer are structured
- The AMD module format requirements

Examine the oCIS documentation and existing extensions to synthesise 
a technical approach for my photo gallery extension.
```

Claude reads documentation, examines code patterns, and produces an architectural plan.

### Phase 3: Scaffolding (~15 minutes)

Ask Claude to create the project structure on the server:

```text
Please create the project structure for my photo-addon extension:
- Git repository initialised
- Build configuration (Vite with AMD output)
- TypeScript configuration
- manifest.json for oCIS
- Basic Vue 3 component structure
- A starter CLAUDE.md file documenting the project for future sessions
```

By the end of this phase, a working (empty) extension exists.

### Phase 4: Implementation (~5-6 hours)

Login to the server and change to the folder where the new project was created, then run:

```text
claude "Please review the CLAUDE.md project file, build the web extension,
and deploy it to our server. Install and configure any dependencies required.
Use 'sudo' where needed to modify or configure the oCIS environment."
```
Once complete, work through features incrementally. Example prompts:

**Feature requests:**

```text
Please adjust the logic to start by searching for today's date (unless 
overwritten by the dropdowns at the top), then display them on the screen. 
If the images don't fill the screen, load up another day, until the screen 
is filled plus some buffer. Then when I scroll, bring in more images.
```

```text
In the map view, when I hover over a dot, the thumbnail shows up above 
off the screen. Can you fix the positioning?
```

**Bug fixes:**

```text
I'm seeing this error in the console:
[paste error here]

Can you investigate and fix it?
```

### Phase 5: Polish (~3 hours)

Once core features work, request comprehensive cleanup:

**Performance optimisation:**

```text
Are there any performance improvements that we can make?
```

**Bug hunting:**

```text
Are there any memory leaks or bugs that can be found?
```

**Code quality:**

```text
Are there any inefficiencies we can improve on? Unused variables? 
Duplicate functions? Excessive if statements? Overwriting CSS styles? 
Overly complex logic?
```

**Error handling:**

```text
Are there any error handling gaps or functions that may fall through 
without handling all scenarios?
```

**Documentation:**

```text
Are there any missing comments? Or complex code that could use 
clear documentation?
```

**Internationalisation:**

```text
Are there any visible text strings that are not yet set up for translation?
```

**Comprehensive cleanup:**

```text
Pulling from the backlog list, can you do a good code refactor to simplify, 
add all required unit tests, and add i18n support, then commit and push the 
code with the right comments.
```

## Active Supervision

{{< hint type=warning title="This is NOT 'Set It and Forget It'" >}}
While no coding is required, this approach demands active supervision. Developers cannot simply give Claude a prompt and walk away.
{{< /hint >}}

### When to Intervene

There are times when Claude Code may:
- Deviate from requirements (with good intentions)
- Undo something it had just done
- Go down a path that doesn't match the architecture

The decisions are always logical, but sometimes the harder path is required. Expect to intervene and redirect approximately ~5% of the time.

**Watch for:**
- Changes that don't match architectural decisions
- Unnecessary complexity being added
- Repeated undo/redo cycles
- Deviation from oCIS extension patterns

### Code Cleanup Is Essential

AI-generated code works, but it can accumulate:
- Redundant logic
- Unused variables
- Overly complex conditionals
- Duplicate CSS styles
- Missing good practice for the orders of directives
- Missing the definition of required interactive handlers

When submitting a pull request (PR), issues highlighted during CI steps, such as linting, can be resolved by explicitly asking for code refactoring and optimisation.

### UI Testing Is Manual

The AI can write unit tests, but exploratory testing requires human judgment:
- Click around the interface
- Try unexpected inputs
- Test with different user accounts (search in oCIS is space-scoped!)
- Check mobile layouts
- Look for edge cases

## Debugging Techniques

When issues arise, these techniques help Claude understand and fix problems:

### Console Output

Open Chrome DevTools, copy console errors, and paste into Claude:

```text
I'm seeing this error in the console:
[paste error here]

Can you investigate and fix it?
```

### HAR Files

For performance issues, export HAR files (network traffic captures):

1. Open DevTools > Network tab
2. Reproduce the issue
3. Right-click > Save all as HAR
4. Share the HAR file with Claude for analysis

### HTML/CSS Inspection

Copy specific HTML/CSS snippets from DevTools:

```text
The thumbnail hover popup is positioned incorrectly. Here's the current 
CSS and HTML structure:
[paste code here]

Can you fix the positioning so the popup appears above the cursor 
but stays within the viewport?
```

{{< hint type=tip title="Chrome Extension Limitations" >}}
Claude has a Chrome extension, but it may be limited for this use case. Manually copying console output and HTML/CSS snippets into Claude is often more effective.
{{< /hint >}}

## Example Web Extensions Built with This Approach

### Photo Gallery Extension (web-app-photo-addon)

Features:
- **Timeline View**: Photos grouped by EXIF capture date with infinite scroll
- **Pinch-to-zoom calendar**: Navigate between day, month, and year views
- **EXIF Metadata Panel**: Camera make/model, aperture, ISO, focal length
- **Map View**: Interactive Leaflet map with marker clustering
- **Lightbox**: Full-screen viewer with keyboard navigation

Tech stack: Vue 3 Composition API, TypeScript, Leaflet.js, AMD module format

### Advanced Search Extension (web-app-advanced-search)

Features:
- **Photo Metadata Filters**: Search by camera, date range, aperture, ISO
- **Filter Chips**: Visual indicators with one-click removal
- **Search Saving**: Save and retrieve search queries
- **KQL Search**: Direct advanced query editing
- **Index Statistics**: View indexing status and configuration

### Backend Changes in oCIS

To enable photo metadata search, changes were contributed to oCIS core:
- Photo field mappings in the KQL parser
- `Store=true` configuration in Bleve for field retrieval
- WebDAV properties for photo metadata

This was the trickiest part—understanding why searches returned empty results. The answer: Bleve was indexing fields but not storing them for retrieval.

## Tips for oCIS Development

### 1. Start with Conversation, Then Architecture

Don't jump straight into implementation. Spend ~30-60 minutes in regular chat, discussing the problem space and having Claude research the codebase. Those early conversations shape the architecture that makes implementation smooth.

### 2. Work Incrementally

Ask for one feature at a time. Get it working, then move to the next. This creates natural checkpoints and makes debugging much easier.

### 3. Use Browser DevTools Liberally

Console errors, network requests, and performance profiles are incredibly valuable context. Copy and paste them into AI conversations.

### 4. Understand oCIS Extension Requirements

Extensions must:
- Use CJS module format (not ES modules)
- Have a proper `manifest.json`
- Use the oCIS web SDK for API access

Claude learns these constraints during the research phase and applies them consistently.

### 5. Test with Multiple Accounts

Search in oCIS is space-scoped. What works for an admin account might not work for regular users.

### 6. Read Maintainer Feedback Carefully

PR reviews teach things no documentation mentions. The oCIS team's suggestions can significantly improve code quality.

## Contributing to oCIS or web-extensions

Once a web-extension is ready, Claude can handle the entire contribution workflow:

- Install and configure Git
- Create forks of the oCIS repositories
- Create feature branches
- Generate commits with appropriate messages
- Push changes to GitHub
- Create pull requests with descriptions and screenshots

Example prompts:

```text
Please fork the owncloud/web-extensions repository to my GitHub account 
and clone it to my server.
```

```text
Create a new branch called feat/web-app-photo-addon and commit our 
extension with an appropriate commit message.
```

```text
Please create a pull request for our photo-addon extension. Include 
screenshots from the /screenshots directory and a clear description 
of what the extension does.
```

**Contributions made by this approach:**
- **oCIS core**: Photo metadata search backend
- **web-extensions**: Photo gallery extension
- **web-extensions**: Advanced search extension

{{< hint type=note title="CLA Requirement" >}}
Contributors need to sign the ownCloud Contributor License Agreement (CLA) for PRs to be accepted. Claude can help find the CLA link and explain the process.
{{< /hint >}}

## Notes

Additional steps can be required in the oCIS repository for changes in the web-extensions repository, in particular:

- If the web extension requires it, add proposals to the back-end code in the oCIS codebase.
- In order to make the new web extension more easily accessible to the public, the `ocis_full` deployment example in the oCIS repository should be updated. The new oCIS version will then reference it automatically and the documentation team will update the release notes and documentation accordingly.

## Summary

Using this approach, complete web-extensions can be built without writing code manually. No git commands need to be typed manually. No pnpm commands either.

Every commit, every push, every PR—all generated through AI assistance. The backend changes to oCIS core were approved and merged. The web-extensions demonstrate the viability of this development approach.

The barrier to entry for meaningful open source contribution is significantly lower with AI-assisted development. For those interested in contributing to oCIS but feeling intimidated by the codebase, this approach is worth considering.

Please note that, as with all contributions, this approach requires a full review process from the code maintainers, with iterations required to finalise all steps. To ensure stability and maintainability, all feedback must be incorporated into the proposed changes.

---

## Resources

| Resource | URL |
|----------|-----|
| MCP Quickstart | https://modelcontextprotocol.io/quickstart |
| MCP Local Servers | https://modelcontextprotocol.io/docs/develop/connect-local-servers |
| Claude Desktop MCP | https://support.claude.com/en/articles/10949351-getting-started-with-local-mcp-servers-on-claude-desktop |
| oCIS Web Extensions | https://github.com/owncloud/web-extensions |
| oCIS Admin Documentation | https://doc.owncloud.com/ocis/ |
| oCIS Developer Documentation | https://owncloud.dev |
---

*This guide was written with AI assistance and reviewed by the ownCloud documentation team for accuracy.*

