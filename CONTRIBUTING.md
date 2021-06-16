First off, thanks for taking the time to consider to contribute ownCloud Infinite Scale!

The following is a set of guidelines for contributing to most of the projects hosted in the [ownCloud Organization](https://github.com/owncloud) on [GitHub](https://www.github.com). These are mostly guidelines, not rules. Use your best judgment, and feel free to propose changes to this document in a pull request.

For simplicity this document mostly refers to the [ocis subproject](https://www.github.com/ocis), but it should be easily transferable to other sub projects.

#### Table Of Contents


[I don't want to read this whole thing, I just have a question!!!](#i-dont-want-to-read-this-whole-thing-i-just-have-a-question)

[What should I know before I get started?](#what-should-i-know-before-i-get-started)
  * [Atom and Packages](#atom-and-packages)
  * [Atom Design Decisions](#design-decisions)

[How Can I Contribute?](#how-can-i-contribute)
  * [Reporting Bugs](#reporting-bugs)
  * [Suggesting Enhancements](#suggesting-enhancements)
  * [Your First Code Contribution](#your-first-code-contribution)
  * [Pull Requests](#pull-requests)

[Styleguides](#styleguides)
  * [Git Commit Messages](#git-commit-messages)
  * [JavaScript Styleguide](#javascript-styleguide)
  * [Documentation Styleguide](#documentation-styleguide)

[Additional Notes](#additional-notes)
  * [Issue and Pull Request Labels](#issue-and-pull-request-labels)

  ## I don't want to read this whole thing I just have a question!!!

> **Note:** [Please don't file an issue to ask a question.](https://blog.atom.io/2016/04/19/managing-the-deluge-of-atom-issues.html) You'll get faster results by using the resources below.

For general questions, please refer to [ownCloud's FAQs](https://owncloud.com/faq/) or ask on the [ownCloud Central Server](https://central.owncloud.org/).

We have a [Rocket Chat Server](https://talk.owncloud.com/channel/infinitescale) to answer your questions specifically to ownCloud Infinite Scale.

## What should I know before I get started

### ownCloud is hosted on Github.

To effectivly contribute to ownCloud Infinite Scale, you should make sure you have an Github account. You can get that for free at [Github](https://github.com/join). You can find howtos on the internet, for example [here](https://www.wikihow.com/Create-an-Account-on-GitHub).

For other ways of contributing, for example with translations, other systems require you to have an account, for example [Transifex](https://www.transifex.com)

The ownCloud development follows the strict Github based workflow of

### The ownCloud Company, the Engineering Partners and Community

ownCloud Infinite Scale is heavily developed by a number of developer that are employed by the [ownCloud company](https://www.owncloud.com), which is located in Germany, operating on the whole planet supporting customers with their ownCloud Setups. In addition there are engineering partners who might also work full time on ownCloud related code.

Because of that fact, the pace that the development is moving forward is sometimes high for people who are not willing and able to spend an comparable amount of time to contribute. Even though this can be a challenge, it should not scare anybody away. It is our clear statement that we feel honored by everybody who is interested in our work and improves it, no matter how big the contribution might be.

We are doing our best to listen to, review and consider all changes that are brought forward following this guideline and make sense for the project. That is true for the ownCloud company and also the engineering partners.

### Licensing and CLA

We are very happy that there is *no CLA* required for most of the code of ownCloud Infinite Scale.

Currently, only for the following parts you need to sign a [Contributors License Agreement](https://en.wikipedia.org/wiki/Contributor_License_Agreement) for the following components:

* [ownCloud Web](https://github.com/owncloud/web/): [Link to CLA]()

Please make sure to read and understand the details of the CLA before starting to invest time on a component that requires it. If you have any questions or concerns please feel free to raise them with us.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report for ownCloud Infinite Scale. Following these guidelines helps maintainers and the community understand your report :pencil:, reproduce the behavior :computer: :computer:, and find related reports :mag_right:.

Before creating bug reports, please check [this list](#before-submitting-a-bug-report) as you might find out that you don't need to create one. When you are creating a bug report, please [include as many details as possible](#how-do-i-submit-a-good-bug-report). Fill out [the required template](https://github.com/owncloud/.github/blob/master/.github/ISSUE_TEMPLATE/bug_report.md), the information it asks for helps us resolve issues faster.

> **Note:** If you find a **Closed** issue that seems like it is the same thing that you're experiencing, open a new issue and include a link to the original issue in the body of your new one. If you have permission to reopen the issue, feel free to do so.

#### Before Submitting A Bug Report

* **Make sure you are running a recent version** Generally developers interest in old versions of a software is dropping very fast once new shiny versions where released. So the general recommendation is: Use the latest released version or even the current master to reproduce problems that you might encounter. That helps a lot to attract developers attention.
* **Determine [which repository the problem should be reported in](#owncloud-repositories)**.
* **Perform a [cursory search](https://github.com/search?q=+is%3Aissue+user%3Aowncloud)** with possibly a more granular filter on the repository, to see if the problem has already been reported. If it has **and the issue is still open**, add a comment to the existing issue instead of opening a new one **if you have new information**. Please abstain from adding "plus ones", except using the Github emojis. That might indicate how many users are affected.

#### How Do I Submit A (Good) Bug Report?

Bugs are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined [which repository](#owncloud-repositories) your bug is related to, create an issue on that repository and provide the following information by filling in [the template](https://github.com/owncloud/ocis/.github/blob/master/.github/ISSUE_TEMPLATE/bug_report.md).

Explain the problem and include additional details to help maintainers reproduce the problem:

* **Use a clear and descriptive title** for the issue to identify the problem.
* **Describe the exact steps which reproduce the problem** in as many details as possible. When listing steps, **don't just say what you did, but explain how you did it**. For example, if you uploaded a file to ownCloud, say which client you used, which way of uploading you choose, if the name was special somehow and how big it was.
* **Provide specific examples to demonstrate the steps**. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples. If you're providing snippets in the issue, use [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
* **Describe the behavior you observed after following the steps** and point out what exactly is the problem with that behavior.
* **Explain which behavior you expected to see instead and why.**
* **Include screenshots and animated GIFs** which show you following the described steps and clearly demonstrate the problem. You can use [this tool](https://www.cockos.com/licecap/) to record GIFs on macOS and Windows, and [this tool](https://github.com/colinkeenan/silentcast) or [this tool](https://github.com/GNOME/byzanz) on Linux.
* **If you report an web browser related problem**, consider to use the browsers Web developer tools (such as the debugger, console or network monitor) to check what happened. Make sure to add screenshots of the utilities in case you are short of time to interprete it.
* **If the problem wasn't triggered by a specific action**, describe what you were doing before the problem happened and share more information using the guidelines below.

Provide more context by answering these questions:

* **Did the problem start happening recently** (e.g. after updating to a new version) or was this always a problem?
* If the problem started happening recently, **can you reproduce the problem in an older version?** What's the most recent version in which the problem doesn't happen? You can download older versions from [the releases page](https://github.com/ownCloud/ocis/releases).
* **Can you reliably reproduce the issue?** If not, provide details about how often the problem happens and under which conditions it normally happens.

Include details about your configuration and environment as ask for in the template.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion for ownCloud Infinite Scale, including completely new features and minor improvements to existing functionality. Following these guidelines helps maintainers and the community understand your suggestion :pencil: and find related suggestions :mag_right:.

Before creating enhancement suggestions, please check [this list](#before-submitting-an-enhancement-suggestion) as you might find out that you don't need to create one. When you are creating an enhancement suggestion, please [include as many details as possible](#how-do-i-submit-a-good-enhancement-suggestion). Fill in [the template](https://github.com/owncloud/ocis/.github/blob/master/.github/ISSUE_TEMPLATE/feature_request.md), including the steps that you imagine you would take if the feature you're requesting existed.

#### Before Submitting An Enhancement Suggestion

* **Check if there's already [an extension](https://marketplace) which provides that enhancement.**
* **Perform a [cursory search](https://github.com/search?q=+is%3Aissue+user%3Aowncloud)** to see if the enhancement has already been suggested. If it has, add a comment to the existing issue instead of opening a new one. Feel free ot use the Github emojis to indicate that you are in favour of an enhancement request.

#### How Do I Submit A (Good) Enhancement Suggestion?

Enhancement suggestions are tracked as [GitHub issues](https://guides.github.com/features/issues/). After you've determined [which repository](#owncloud-repositories) your enhancement suggestion is related to, create an issue on that repository and provide the following information:

* **Use a clear and descriptive title** for the issue to identify the suggestion.
* **Provide a step-by-step description of the suggested enhancement** in as many details as possible.
* **Provide specific examples to demonstrate the steps**. Include copy/pasteable snippets which you use in those examples, as [Markdown code blocks](https://help.github.com/articles/markdown-basics/#multiple-lines).
* **Explain why this enhancement would be useful** to most ownCloud users.
* **List some other projects or products where this enhancement exists.**
* **Specify which version of ownCloud you're using.**






