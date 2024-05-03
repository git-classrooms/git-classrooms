export type Reversed<T extends Record<string | number, string | number>> = {
  [x in T[keyof T]]: keyof T;
};

export const reversed = <T extends Record<string | number, string | number>>(obj: T) =>
  Object.entries(obj).reduce((acc, [key, value]) => ({ ...acc, [value]: key }), {} as Reversed<T>);

export type Simplify<T> = {
  [P in keyof T]: T[P];
  // eslint-disable-next-line @typescript-eslint/ban-types
} & {};

export type DeepRequired<T> = {
  [P in keyof T]-?: DeepRequired<T[P]>;
  // eslint-disable-next-line @typescript-eslint/ban-types
} & {};

export type SubPartial<T, K extends keyof T> = Simplify<Omit<T, K> & Partial<Pick<T, K>>>;
