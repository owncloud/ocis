## Text wrapping behaviour

Use one of the following classes to enforce a certain wrapping behaviour.

| Class             | Description                                                                                   |
| ----------------- | --------------------------------------------------------------------------------------------- |
| .oc-text-overflow | Sets overflow to `hidden` without resizing its container.                                     |
| .oc-text-nowrap   | Doesn't break to new lines.                                                                   |
| .oc-text-truncate | Doesn't break to new lines. Text will be truncated, showing an ellipsis instead if necessary. |
| .oc-text-break    | Text will break to new lines at word ends if it exceeds one line.                             |

## Text sizes

The ownCloud Design System uses a default font size which can be set on the html root element and then lets you use the
following text size classes to choose a size relative to the default font size.

| Class           | Description                                                                                         |
|-----------------|-----------------------------------------------------------------------------------------------------|
| .oc-text-xsmall | Sets the font size to 0.72rem. Value can be overwritten by setting the `oc-font-size-xsmall` token. |
| .oc-text-small  | Sets the font size to 0.86rem. Value can be overwritten by setting the `oc-font-size-small` token.  |
| .oc-text-medium | Sets the font size to 1rem. Value can be overwritten by setting the `oc-font-size-medium` token.    |
| .oc-text-large  | Sets the font size to 1.14rem. Value can be overwritten by setting the `oc-font-size-large` token.  |
| .oc-text-xlarge | Sets the font size to 1.29rem. Value can be overwritten by setting the `oc-font-size-xlarge` token. |
