export const isNonNullable = <T>(val: T): val is NonNullable<T> => val != null
