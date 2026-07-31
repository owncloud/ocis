export const call = function* <T>(p: Promise<T>): Generator<Promise<T>, T, T> {
  return yield p
}
