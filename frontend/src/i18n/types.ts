import { resources, defaultNS } from "./index";

declare module "i18next" {
  interface CustomTypeOptions {
    defaultNS: typeof defaultNS;
    resources: (typeof resources)["en"];
    returnNull: false;
  }
}

export type TranslationNamespace = keyof (typeof resources)["en"];
