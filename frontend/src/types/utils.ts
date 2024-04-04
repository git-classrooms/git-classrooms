export type Reversed<T extends Record<string | number, string | number>> = {
  [x in T[keyof T]]: keyof T;
};

export const reversed = <T extends Record<string | number, string | number>>(obj: T) =>
  Object.entries(obj).reduce((acc, [key, value]) => ({ ...acc, [value]: key }), {} as Reversed<T>);