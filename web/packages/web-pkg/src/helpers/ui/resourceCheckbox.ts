export const focusCheckbox = (id: string) => {
  const checkbox = document.querySelectorAll(
    `[data-item-id="${id}"] input[type=checkbox]`
  )?.[0] as HTMLInputElement
  if (checkbox) {
    checkbox.focus()
  }
}
