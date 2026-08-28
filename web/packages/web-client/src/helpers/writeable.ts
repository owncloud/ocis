/**
 * Strips `readonly` from an object's own properties.
 *
 * The generated graph models mark response-only fields `readonly`, which is right for a
 * server response but gets in the way of the editable copies and test fixtures we build
 * from them.
 */
export type Writeable<T> = { -readonly [K in keyof T]: T[K] }

/**
 * Views a value as {@link Writeable} so a single field can be assigned.
 *
 * Use it on the object that owns the field, not on the whole tree — `readonly` is stripped
 * one level deep only, which keeps the escape hatch visible at each assignment:
 *
 * ```ts
 * writeable(user.drive.quota).total = 42
 * ```
 */
export const writeable = <T>(value: T): Writeable<T> => value
