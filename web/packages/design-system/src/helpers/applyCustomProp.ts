export const applyCustomProp = (key: string, value: string | undefined) => {
  if (value === undefined) {
    return
  }
  const root = document.querySelector(':root') as HTMLElement
  const prop = '--oc-' + key
  if (value === '') {
    root.style.removeProperty(prop)
  } else {
    root.style.setProperty(prop, value)
  }
}

export const removeCustomProp = (key: string) => {
  ;(document.querySelector(':root') as HTMLElement).style.removeProperty('--oc-' + key)
}
