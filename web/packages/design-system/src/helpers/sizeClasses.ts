const sizeClassMappings = {
  xxxsmall: 'xxxs',
  xxsmall: 'xxs',
  xsmall: 'xs',
  small: 's',
  medium: 'm',
  large: 'l',
  xlarge: 'xl',
  xxlarge: 'xxl',
  xxxlarge: 'xxxl',
  remove: 'rm'
}

export function getSizeClass(size: string) {
  return sizeClassMappings[size as keyof typeof sizeClassMappings]
}
