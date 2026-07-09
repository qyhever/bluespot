declare module 'fake-progress' {
  interface FakeProgressOptions {
    timeConstant?: number
    autoStart?: boolean
    parent?: FakeProgress
    parentStart?: number
    parentEnd?: number
  }

  interface FakeSubProgressOptions extends FakeProgressOptions {
    start?: number
    end?: number
  }

  export default class FakeProgress {
    progress: number

    constructor(options?: FakeProgressOptions)

    start(): void
    end(): void
    stop(): void
    createSubProgress(options: FakeSubProgressOptions): FakeProgress
    setProgress(progress: number): void
  }
}
