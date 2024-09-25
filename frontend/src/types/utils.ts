import { z } from "zod";

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export type Reversed<T extends Record<keyof T, keyof any>> = {
  [K in T[keyof T]]: {
    [P in keyof T]: T[P] extends K ? P : never;
  }[keyof T];
  // eslint-disable-next-line @typescript-eslint/ban-types
} & {};

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const reversed = <T extends Record<keyof T, keyof any>>(obj: T) =>
  Object.fromEntries(Object.entries(obj).map(([key, value]) => [value, key])) as Reversed<T>;

export type Simplify<T> = {
  [P in keyof T]: T[P];
  // eslint-disable-next-line @typescript-eslint/ban-types
} & {};

export type DeepRequired<T> = {
  [P in keyof T]-?: DeepRequired<T[P]>;
  // eslint-disable-next-line @typescript-eslint/ban-types
} & {};

export type SubPartial<T, K extends keyof T> = Simplify<Omit<T, K> & Partial<Pick<T, K>>>;

// eslint-disable-next-line @typescript-eslint/no-explicit-any
export const zodEnumFromObjectKeys = <K extends string>(obj: Record<K, any>) => {
  const [fistKey, ...otherKeys] = Object.keys(obj) as K[];
  return z.enum([fistKey!, ...otherKeys]);
};
