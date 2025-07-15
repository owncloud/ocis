First of all, thank you for taking the time to read this and your interest in contributing to ownCloud Infinite Scale!

The following is a set of guidelines for contributing to most of the projects hosted in the [ownCloud Organization](https://github.com/owncloud) on [GitHub](https://www.github.com). These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

For simplicity reasons, this document mostly refers to the ownCloud [Infinite Scale project](https://www.github.com/owncloud/ocis), but it should be easily transferable to other (sub)projects.

#### Table Of Contents

[I don't want to read this whole thing, I just have a question](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)

[What should I know before I get started](#what-should-i-know-before-i-get-started)
*   [ownCloud is hosted on GitHub](#owncloud-is-hosted-on-github)
*   [ownCloud Company, Engineering Partners and Community](#owncloud-company,-engineering-partners-and-community)
*   [Licensing and CLA](#licensing-and-cla)

[How Can I Contribute](#how-can-i-contribute)
*   [Help spreading the word](#help-spreading-the-word)
*   [Reporting Bugs](#reporting-bugs)
*   [Suggesting Enhancements](#suggesting-enhancements)
*   [Your First Code Contribution](#your-first-code-contribution)
*   [Pull Requests](#pull-requests)
*   [Documentation Contributions](#documentation-contributions)
*   [Internationalization](#internationalization)
*   [Deployments: Docker, Ansible and friends](#deployments-docker-ansible-and-friends)

[Styleguides](#styleguides)
*   [Git Commit Messages](#git-commit-messages)
*   [Golang Styleguide](#golang-styleguide)
*   [Web Styleguide](#web-styleguide)
*   [Documentation Styleguide](#documentation-styleguide)

[Additional Notes](#additional-notes)
*   [Issue and Pull Request Labels](#issue-and-pull-request-labels)

  ## I don't want to read this whole thing I just have a question

> **Note:** Please don't file an issue to ask a question. You'll get faster results by using the resources below.

For general questions, please refer to [ownCloud's FAQs](https://owncloud.com/faq/) or ask on the [ownCloud Central Server](https://central.owncloud.org/).

We also have a [Rocket Chat Server](https://talk.owncloud.com/channel/ocis) to answer your questions specifically about ownCloud Infinite Scale.

## What should I know before I get started

### ownCloud is hosted on GitHub

To effectively contribute to ownCloud Infinite Scale, you need a GitHub account. You can get that for free at [GitHub](https://github.com/join). You can find howtos on the internet, for example [here](https://www.wikihow.com/Create-an-Account-on-GitHub).

For other ways of contributing, for example with translations, other systems require you to have an account as well, for example [Transifex](https://www.transifex.com).

The ownCloud project follows the strict GitHub workflow of development as briefly [described here](https://guides.github.com/introduction/flow/).

### ownCloud Company, Engineering Partners and Community

ownCloud Infinite Scale is largely created by developers who are employed by the [ownCloud company](https://www.owncloud.com), which is located in Germany. It is providing support for ownCloud for customers worldwide. In addition, there are engineering partners who also work full time on ownCloud related code, for example in [CERN REVA](https://github.com/cs3org/reva/).

Because of that fact, the pace that the development is moving forward is sometimes high for people who are not willing and/or able to spend a comparable amount of time to contribute. Even though this can be a challenge, it should not scare anybody away. Here is our clear commitment that we feel honored by everybody who is interested in our work and improves it, no matter how big the contribution might be.

We as the full time devs from either organization are doing our best to listen, review and consider all changes that are brought forward following this guideline and make sense for the project.

### Licensing and CLA

We are very happy that there is *no CLA* required for most of the code of ownCloud Infinite Scale.

Currently, only for the following components you need to sign a [Contributors License Agreement](https://en.wikipedia.org/wiki/Contributor_License_Agreement):

*   [ownCloud Web](https://github.com/owncloud/web/): [Link to CLA](https://owncloud.com/contribute/join-the-development/contributor-agreement/)

Please make sure to read and understand the details of the CLA before starting to invest time on a component that requires it. If you have any questions or concerns please feel free to raise them with us.

## How Can I Contribute

There are many ways to contribute to open source projects, and all are equally valuable and appreciated.

### Help spreading the word

This way to contribute to the project can not be overestimated: People who talk about their experience with ownCloud Infinite Scale and help others with that are the key to success of the project.

There are too many ways of doing that to line them up here, but examples are answering questions in [ownCloud Central](https://central.owncloud.org/) or on [ownCloud Talk](https://talk.owncloud.com/channel/ocis), writing blog posts etc. pp.

There is no formal guideline to this, just do it :-)

### Reporting Bugs

This section guides you through submitting a bug report for ownCloud Infinite Scale. Following these guidelines helps maintainers and the community understand your report :pencil:, reproduce the behavior :computer: :computer:, and find related reports :mag_right:.

Before creating bug reports, please check [this list](#before-submitting-a-bug-report) as you might find out that you don't need to create one. When you are creating a bug report, please [include as many details as possible](#how-do-i-submit-a-good-bug-report). Fill out [the required template](https://github.com/owncloud/ocis/issues/new?Type%3ABug&template=bug_report.md), the information it asks for helps us resolve issues faster.

> **Note:** If you find a **Closed** issue that seems like it is the same thing that you're experiencing, open a new issue and include a link to the original issue in the body of your new one. If you have permission to reopen the issue, feel free to do so.

#### Before Submitting A Bug Report

*   **Make sure you are running a recent version** Usually, developers' interest in old versions of software drops very fast once a new shiny version has been released. So the general recommendation is: Use the latest released version or even the current master to reproduce problems that you might encounter. That helps a lot to attract developers attention.
*   **Determine which [repository](https://github.com/owncloud) the problem should be reported in**.
*   **Perform a [cursory search](https://github.com/search?q=+is%3Aissue+user%3Aowncloud)** with possibly a more granular filter on the repository, to see if the problem has already been reported. If it has **and the issue is still open**, add a comment to the existing issue instead of opening a new one **if you have new information**. Please abstain from adding "plus ones", except using the GitHub emojis. That might indicate how many users are affected.

#### How Do I Submit A (Good) Bug Report

Bugs are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined [which repository](https://github.com/owncloud) your bug is related to, create an issue on that repository and provide the following information by filling in [the template](https://github.com/owncloud/ocis/issues/new?Type%3ABug&template=bug_report.md).

Explain the problem and include additional details to help maintainers reproduce the problem:

*   **Use a clear and descriptive title** for the issue to identify the problem.
*   **Describe the exact steps which reproduce the problem** in as many details as possible. Start with describing, from a user perspective, what you tried to achieve, i.e. "I want to share some pictures with Grandma". When listing steps, **don't just say what you did, but explain how you did it**. For example, if you uploaded a file to ownCloud, say which client you used, which way of uploading you chose, if the name was special somehow and how big it was.
*   **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples. If you're providing snippets in the issue, use [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
*   **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
*   **Explain which behavior you expected to see instead and why.**
*   **Include screenshots and animated GIFs** which show you following the described steps and clearly demonstrate the problem. You can use [this tool](https://www.cockos.com/licecap/) to record GIFs on macOS and Windows, and [this tool](https://github.com/colinkeenan/silentcast) or [this tool](https://github.com/GNOME/byzanz) on Linux.
*   **If you report a web browser related problem**, consider to using the browser's Web developer tools (such as the debugger, console or network monitor) to check what happened. Make sure to add screenshots of the utilities if you are short of time to interpret it.
*   **If the problem wasn't triggered by a specific action**, describe what you were doing before the problem happened and share more information using the guidelines below.

Provide more context by answering these questions:

*   **Did the problem start happening recently** (e.g. after updating to a new version) or was this always a problem?
*   If the problem started happening recently, **can you reproduce the problem in an older version?** What's the most recent version in which the problem doesn't happen? You can find more information about how to set up [test environments](https://owncloud.dev/ocis/development/testing/) in the [developer documentation](https://owncloud.dev/#developer-documentation).
*   **Can you reliably reproduce the issue?** If not, provide details about how often the problem happens and under which conditions it normally happens.

Include details about your configuration and environment as asked for in the template.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for ownCloud Infinite Scale, including completely new features and minor improvements to existing functionality. Following these guidelines helps maintainers and the community understand your suggestion :pencil: and find related suggestions :mag_right:.

Before creating enhancement suggestions, please check [this list](#before-submitting-an-enhancement-suggestion) as you might find out that you don't need to create one. When you are creating an enhancement suggestion, please [include as many details as possible](#how-do-i-submit-a-good-enhancement-suggestion). Fill in [the template](https://github.com/owncloud/ocis/.github/blob/master/.github/ISSUE_TEMPLATE/feature_request.md), including the steps that you imagine you would take if the feature you're requesting existed.

#### Before Submitting An Enhancement Suggestion

*   **Check if there's already an extension or other component which provides that enhancement, even in a different way.**
*   **Perform a [cursory search](https://github.com/search?q=+is%3Aissue+user%3Aowncloud)** to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one. Feel free to use the GitHub emojis to indicate that you are in favour of an enhancement request.

#### How Do I Submit A (Good) Enhancement Suggestion

Enhancement suggestions are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined [which repository](https://github.com/owncloud) your enhancement suggestion is related to, create an issue on that repository and provide the following information:

*   **Use a clear and descriptive title** for the issue to identify the suggestion.
*   **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
*   **Provide specific examples to demonstrate the steps**. Include copy/pasteable snippets which you use in those examples, as [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
*   **Explain why this enhancement would be useful** to most ownCloud users.
*   **List some other projects or products where this enhancement exists.**

### Your First Code Contribution

Unsure where to begin contributing to ownCloud? You can start by looking through these `Needs-help` issues:

*   The [Good first issue](https://github.com/owncloud/ocis/labels/Topic%3Agood-first-issue) label marks good items to start with.
*   [Tests needed](https://github.com/owncloud/ocis/labels/Interaction%3ANeeds-tests) - issues which would benefit from a test.
*   [Help wanted issues](https://github.com/owncloud/ocis/labels/Interaction%3ANeeds-help) - issues which should be a bit more involved.

It is fine to pick one of the list following personal preference. While not perfect, number of comments is a reasonable proxy for impact a given change will have.

To find out how to set up ownCloud Infinite Scale for local development please refer to the [Developer Documentation](https://owncloud.dev/ocis/development/getting-started/). It contains a lot of information that will come in handy when starting to work on the project.

### Pull Requests

All contributions to ownClouds projects use so-called pull requests following the [GitHub PR workflow](https://guides.github.com/introduction/flow/).

Please follow these steps to have your contribution considered by the maintainers:

*   Follow all instructions in [the template](https://github.com/owncloud/ocis/blob/master/.github/pull_request_template.md)
*   Follow the [styleguides](#styleguides) where applicable
*   After you submit your pull request, verify that all [status checks](https://help.github.com/articles/about-status-checks/) are passing <details><summary>What if the status checks are failing?</summary>If a status check is failing, and you believe that the failure is unrelated to your change, please leave a comment on the pull request explaining why you believe the failure is unrelated. A maintainer will re-run the status check for you. If we conclude that the failure was a false positive, then we will open an issue to track that problem with our status check suite.</details>

While the prerequisites above must be satisfied prior to having your pull request reviewed, the reviewer(s) may ask you to complete additional design work, tests, or other changes before your pull request can be ultimately accepted.

### Documentation Contributions

ownCloud is very proud of the documentation it has, which is the work of a great team of people. Of course, also the documentation is open to contributions.

See the [Getting Started Guide](https://owncloud.dev/ocis/development/getting-started/) on how to get started. Other useful information is summarized in the [Documentation Readme](https://github.com/owncloud/docs).

### Internationalization

Our projects are getting translated into many languages to allow people from all over the world to use ownCloud in their native language. For translations, ownCloud uses [Transifex](https://www.transifex.com) as a community based collaboration platform for internationalization.

For contributions please refer to the [Transifex Resources](https://www.transifex.com/resources/) to learn how to improve ownClouds translations there.

### Deployments: Docker, Ansible and friends

Depending on the ownCloud component, there is complex deployment tooling to install in various environments. There is for example [ownCloud Ansible](https://github.com/owncloud-ansible) with Ansible resources. Contributions to that are very appreciated and follow the same guidelines as every other code contribution.

## Styleguides

To keep up with a consistent code and tooling landscape, some ownCloud modules maintain styleguides for contributions. It is mandatory to follow them in contributions.

### Git Commit Messages

*   Use the present tense ("Add feature" not "Added feature")
*   Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
*   Limit the first line to 72 characters or less
*   Reference issues and pull requests liberally after the first line
*   When only changing documentation, include `[docs-only]` in the commit title

### Golang Styleguide

Use the built-in golang code formatter before submitting the patch. Also, consulting documentation like [Effective Go](https://golang.org/doc/effective_go) or [Practical Go](http://bit.ly/gcsg-2019) helps to improve the code quality.

### Web Styleguide

Refer to related documents in the [ownCloud Web Repository](https://github.com/owncloud/web).

#### Documentation Styleguide

See the [ownCloud Documentation Styleguide](https://github.com/owncloud/docs/blob/master/docs/style-guide.md).

## Additional Notes

### Issue and Pull Request Labels

This section lists the labels we use to help us track and manage issues and pull requests. Most labels are used across all ownCloud repositories, but some are specific.

[GitHub search](https://help.github.com/articles/searching-issues/) makes it easy to use labels for finding groups of issues or pull requests you're interested in. To help you find issues and pull requests, each label can be used in search links for finding open items with that label in the ownCloud repositories.

The labels are loosely grouped by their purpose, but it's not required that every issue has a label from every group or that an issue can't have more than one label from the same group.

The list here contains all the more general categories of issues which are followed by a colon and a specific value. For example severity 1 looks like `Severity:sev1-critical`.

#### Platform

Describes the platform the issue is happening on, i.e. iOS or Windows.

#### Estimation

T-Shirt sizes for effort estimation to fix that bug or implement an enhancement. Ranges from XS to XXXL.

#### Priority

P1 to P4 (lowest) to indicate a priority. Mostly a tool for internal project management and support.

#### QA

Flags to indicate the internal QA status in terms of process and priority. Please leave alone unless you're QA ;-)

#### Severity

Severity for the product, mostly impact on user.

#### Type

The issue type, helps to structure the issues in the agile categories (Epic, Story...) but also organizational ones.

#### Topic

A general category of the topic of a ticket.

#### Category

Categorizes the issue to also indicate the type of the issue.

#### Status

The status in the ticket life cycle. Keep an eye on that one, especially for the `Waiting-for-Feedback` tag which might indicate that the reporter is asked for feedback.

#### Interaction

Another label that indicates the type of the issue.

#### Browser

Important for browser dependent web issues. It specifies the browser that shows the error.

#### Early-Adopter

Tags issues that were reported by one of the oCIS early adopters, i.e. customers and users who start using ownCloud Infinite Scale before its general availability.
