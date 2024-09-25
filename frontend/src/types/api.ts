export type MutationOptions<T> = T & {
  onError?: (error: Error) => void;
  onSuccess?: () => void;
  onSettled?: () => void;
};
