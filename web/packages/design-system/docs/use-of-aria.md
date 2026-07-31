# WAI-ARIA

## Introduction

ARIA stands for "Accessible Rich Internet Application" and serves as a semantic layer that adds information about a web-apps state, about functions of and relations between elements (for example as part of widgets) to assistive technology (such as screen readers). It acts as a "polyfill for HTML", so to speak and helps out in cases where HTML is not enough to semantically declare a widget's function or state. **It does not add functionality, merely announces that a certain (for example: keyboard) behaviour is implemented**. In HTML5 it has semantically "caught up" to some parts of ARIA a little bit (e.g. with elements like `<header>`, `<main>`, `<footer>` replacing `<div role="main" />` and the like) but ARIA's general aim is to supply a way to add semantics to elements, widgets or pattern that cannot be implemented with HTML alone.

> When a developer who is new to accessibility looks at the [ARIA specification](https://www.w3.org/TR/wai-aria-1.1/) it can be seem intimidating, but it doesn't need to be. Most of the ARIA roles, states and properties in that document belong on specialized widgets, and shouldn't be attempted by developers who are new to accessibility.
> ([Source](https://www.davidmacd.com/blog/wai-aria-accessbility-for-average-web-developers.html))

While ARIA can help making especially web-app interfaces more accessible, it is very easy to make a web-based project _less_ accessible by using ARIA. The general advise is to use it sparingly and responsibly, and make sure you know what you are doing.

## The rules of ARIA

Thus, W3C established common rules when using ARIA, and the first one directly relates to its possibly destructive power to accessibility:

1. "If you can use a native HTML element or attribute with the semantics and behavior you require already built in, instead of re-purposing an element and adding an ARIA role, state or property to make it accessible, then do so." - or in short: "Don't use ARIA" ([Source](https://www.w3.org/TR/using-aria/#firstrule))

2. "Do not change native semantics, unless you really have to." ([Source](https://www.w3.org/TR/using-aria/#secondrule))

3. "All interactive ARIA controls must be usable with the keyboard." ([Source](https://www.w3.org/TR/using-aria/#3rdrule))

4. "Do not use role="presentation" or aria-hidden="true" on a focusable element." ([Source](https://www.w3.org/TR/using-aria/#4thrule))

5. "All interactive elements must have an accessible name." ([Source](https://www.w3.org/TR/using-aria/#fifthrule))

## ARIA is a contract or promise to the user

It is stated above, but can't be said enough because it is one of the most common misunderstandings: **ARIA itself does not add functionality**, by using it the author just promises to implement behaviour in a certain way, or that the user agent should perceive one element differently that its original semantics convey.

## Examples

If you add a `role="button"` to a `<span>` assistive technology regard this span as a button. But that measure alone does not transform the element into a button, regarding:

- its focussability
- its ability to proxy the click event handler to also react on SPACE and RETURN key events
- its ability to react to a `disabled` attribute, disabling all event listeners and removing it from the tab order

All of the above mentioned functionality comes for free when using a `<button>` element.

One other example is proclaiming via `role="tab"` and `role="tabpanel"` that a widget is a tab component. Authors have to actually implement [the keyboard usage pattern](https://www.w3.org/TR/wai-aria-practices/examples/tabs/tabs-1/tabs.html) expected by the user, otherwise it is a promise not kept, or, stated differently, a violation of contract with the user.

## Links

- [Demystifying WAI-ARIA - 18 WAI-ARIA attributes that every web developer should know](https://www.davidmacd.com/blog/wai-aria-accessbility-for-average-web-developers.html)
- [WAI-ARIA Cheat sheet](https://www.digitala11y.com/wai-aria-1-1-cheat-sheet/)
- [Presentation slides about how easy it is to ruin a web project with ARIA; ARIA serious?](https://talks.yatil.net/47fUQW/aria-serious)
- [WAI-ARIA specs](https://www.w3.org/TR/wai-aria-1.1/)
