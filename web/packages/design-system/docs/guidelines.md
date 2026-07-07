## Interface Element Guidelines

### Be empathetic

Have in mind that the interface is not built only for able-bodied people but for all humans. Therefore make sure that it can be used, for example, in an non-visual way (e.g. relying on the DOM created by JavaScript) or with non-pointer devices like mouse and touch (e.g. keyboard-only or single-switch devices). Test your work in as many forms of [assistive technology](https://webaccess.berkeley.edu/resources/assistive-technology) as possible.

### Supply meaning with the right choice of elements

Regardless of the visual representation, choose a meaningful and semantic HTML element when building a component. [Find a list of elements here at MDN](https://developer.mozilla.org/en-US/docs/Web/HTML/Element). If an interface concept does not find a match in HTML, supply meaning, information about state, role, relation and name of an element via [WAI-ARIA](https://www.w3.org/TR/wai-aria/).

#### Example: Tabular Content

When outputting a large amount of structured data and related actions to said data are present, grid-like and single sets of data are ordered in columns (like in the `OcFileList` component), whereas single files and folders are ordered in rows, choose the `table` element.

#### Example: Crucial controls

A semantic choice of elements is important for the most vital elements of the ownCloud UI (such as controls related to files and folders), especially for users who do not use a pointing device. Make sure that a control is focusable and conveys its state and function. In the case of the file and folder example mentioned above differentiate between changing the state of the web-app (`<button>`, triggers non-modal dialog, editor) and changing the location (`<a href="/">`, leads to a route, showing contents of a sub folder).

### Supply precise labels

Supply adequately labelled controls especially in critical contexts like file deletion. Instead of relying on "OK" and "Cancel" supply concrete answers to the question the interface is asking the user: "Yes" and "No".

Since only very few pictograms are ambiguous, for good user experience and making the interface accessible for users with cognitive disabilities. always aim to use descriptive text in icon buttons. A rare exception is a cross icon for closing. But make sure every interactive element is at least labelled non-visually with `aria-label`.

### Avoid modal dialogs

The use of the "modal dialog" design pattern should be avoided. Its character is to force a certain interface element into the users field of view and focus, while simultaneously disabling all other present interface elements. The purpose of a modal dialog is to force a user decision. Alternatives could be: non-modal dialogs, disclosure widgets, inline-editing, undo functionality.

### Emphasize core entities

If core entities (Files, Folders, Users) are mentioned in the interface's copy (also applying to HTML e-mails sent by the system), emphasize its mentions by text format (bold, italic ) or by putting the mention in quotes.

### Choose hiding over disabling interactive elements

Under normal circumstances and apart from "Save" actions it is good practice to not disable elements that are not accessible to a user (for example due to a lack of role privileges on their account), but to hide irrelevant controls altogether.
