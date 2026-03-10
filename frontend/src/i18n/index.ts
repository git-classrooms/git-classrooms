import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import enCommon from "./locales/en/common.json";
import enAuth from "./locales/en/auth.json";
import enClassroom from "./locales/en/classroom.json";
import enAssignment from "./locales/en/assignment.json";
import enTeam from "./locales/en/team.json";
import enErrors from "./locales/en/errors.json";
import enValidation from "./locales/en/validation.json";

import deCommon from "./locales/de/common.json";
import deAuth from "./locales/de/auth.json";
import deClassroom from "./locales/de/classroom.json";
import deAssignment from "./locales/de/assignment.json";
import deTeam from "./locales/de/team.json";
import deErrors from "./locales/de/errors.json";
import deValidation from "./locales/de/validation.json";

export const defaultNS = "common";
export const resources = {
  en: {
    common: enCommon,
    auth: enAuth,
    classroom: enClassroom,
    assignment: enAssignment,
    team: enTeam,
    errors: enErrors,
    validation: enValidation,
  },
  de: {
    common: deCommon,
    auth: deAuth,
    classroom: deClassroom,
    assignment: deAssignment,
    team: deTeam,
    errors: deErrors,
    validation: deValidation,
  },
} as const;

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: "en",
    defaultNS,
    detection: {
      order: ["localStorage", "navigator"],
      lookupLocalStorage: "git-classrooms-language",
      caches: ["localStorage"],
    },
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
