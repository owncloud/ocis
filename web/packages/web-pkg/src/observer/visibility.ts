enum states {
  enter,
  exit
}

interface CbackOptions {
  element: Element
  callCount: number
  unobserve: () => void
}

type cback = (options: CbackOptions) => void

type unobserve = (element: Element) => void

interface Cbacks {
  onEnter?: cback
  onExit?: cback
}

class Target {
  private state: states
  private observeEnter: boolean
  private observeExit: boolean
  private onEnterCallCount: number
  private onExitCallCount: number
  private readonly onEnter?: cback
  private readonly onExit?: cback
  public readonly threshold: number
  public readonly unobserver: unobserve

  constructor(unobserver: unobserve, threshold: number, cbacks: Cbacks) {
    this.unobserver = unobserver
    this.threshold = threshold
    this.onEnter = cbacks.onEnter
    this.onExit = cbacks.onExit
    this.observeEnter = !!cbacks.onEnter
    this.observeExit = !!cbacks.onExit
    this.onEnterCallCount = 0
    this.onExitCallCount = 0
  }

  public request(state: states, element: Element) {
    const sharedProps = {
      element: element,
      unobserve: () => this.unobserve(state, element)
    }

    if (state === states.enter && this.observeEnter && this.onEnter) {
      this.onEnterCallCount++
      this.onEnter({
        callCount: this.onEnterCallCount,
        ...sharedProps
      })
    } else if (
      this.state === states.enter &&
      state === states.exit &&
      this.observeExit &&
      this.onExit
    ) {
      this.onExitCallCount++
      this.onExit({
        callCount: this.onExitCallCount,
        ...sharedProps
      })
    }

    this.state = state
  }

  private unobserve(state: states, element: Element) {
    if (state === states.enter) {
      this.observeEnter = false
    } else if (state === states.exit) {
      this.observeExit = false
    }
    if (!this.observeEnter && !this.observeExit) {
      this.unobserver(element)
    }
  }
}

interface VisibilityObserverOptions {
  root?: Element | Document | null
  rootMargin?: string
  threshold?: number
}

export class VisibilityObserver {
  private targets: WeakMap<Element, Target>
  private readonly intersectionObserver: IntersectionObserver
  private readonly options: VisibilityObserverOptions

  constructor(options: VisibilityObserverOptions = {}) {
    this.options = {
      root: options.root,
      rootMargin: options.rootMargin,
      threshold: options.threshold || 0
    }
    this.targets = new WeakMap<Element, Target>()
    this.intersectionObserver = new IntersectionObserver(this.trigger.bind(this), this.options)
  }

  public observe(element: Element, cbacks: Cbacks = {}, threshold?: number): void {
    if (!cbacks.onEnter && !cbacks.onExit) {
      return
    }

    this.targets.set(
      element,
      new Target(this.unobserve.bind(this), threshold || this.options.threshold || 0, {
        onEnter: cbacks.onEnter,
        onExit: cbacks.onExit
      })
    )
    this.intersectionObserver.observe(element)
  }

  public unobserve(element: Element): void {
    this.targets.delete(element)
    this.intersectionObserver.unobserve(element)
  }

  public disconnect(): void {
    this.targets = new WeakMap<Element, Target>()
    this.intersectionObserver.disconnect()
  }

  private trigger(entries: IntersectionObserverEntry[]) {
    entries.forEach((entry: IntersectionObserverEntry) => {
      const observedElement = this.targets.get(entry.target)

      if (!observedElement) {
        return
      }

      observedElement.request(
        entry.isIntersecting && entry.intersectionRatio > observedElement.threshold
          ? states.enter
          : states.exit,
        entry.target
      )
    })
  }
}
