const tape: HTMLElement[] = []
const isDomNode = (node: Element) => node instanceof Element
const DIRECTION_FORWARD = 'forward'
const DIRECTION_BACKWARD = 'backward'

export default {
  methods: {
    focus({ from, to, revert }: { from: HTMLElement; to: HTMLElement; revert: boolean }) {
      const direction = revert ? DIRECTION_BACKWARD : DIRECTION_FORWARD

      if (from && direction === DIRECTION_FORWARD) {
        tape.splice(0, tape.length)
      } else {
        from = document.activeElement as HTMLElement
      }

      if (direction === DIRECTION_FORWARD) {
        tape.push(from)
      }

      if (direction === DIRECTION_BACKWARD) {
        to = tape.pop()
      }

      if (isDomNode(to)) {
        to.focus()
      }
    }
  }
}
