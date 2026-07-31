export default {
  name: 'transform/ods/namespace',
  type: 'name',
  transform: (prop) => {
    return ['oc', prop.name].filter(Boolean).join('-')
  }
}
