#!/bin/env node

import fs from "fs";
import path from "path";

const exceptions = ["change-owned-classroom-member-request"];

const subPartialChange = {
  team: [
    "user-classrooms",
    "get-owned-classroom-member-response",
    "get-joined-classroom-member-response",
    "get-joined-classroom-response",
  ],
  dueDate: ["assignment", "create-assignment-request"],
};

const modelsFile = path.join("models", "index.ts");

const models = fs.readFileSync(modelsFile, "utf8");

let newModels = `import { DeepRequired, SubPartial } from "@/types/utils";\n`;

newModels +=
  models
    .split("\n")
    .filter(Boolean)
    .map((line) => {
      const parts = line.split(" ");
      const importFolder = parts[3].replace(/'\.\/(.*)';/g, "$1");

      let editFunc = (input) => input;

      if (exceptions.some((exception) => importFolder.includes(exception))) {
        return line;
      }

      const subPartials = Object.entries(subPartialChange).reduce((prev, [key, value]) => {
        if (value.includes(importFolder)) {
          return [...prev, `"${key}"`];
        }
        return prev;
      }, []);

      if (subPartials.length > 0) {
        editFunc = (input) => `SubPartial<${input}, ${subPartials.join(" | ")}>`;
      }

      const modelType = importFolder
        .split("-")
        .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
        .join("");

      const importLine = `import { ${modelType} as Old${modelType} } from "./${importFolder}";`;
      const exportLine = `export type ${modelType} = ${editFunc(`DeepRequired<Old${modelType}>`)};`;
      return [importLine, exportLine];
    })
    .flat()
    .join("\n") + "\n";

newModels = newModels.replaceAll("Httperror", "HTTPError");

fs.writeFileSync(modelsFile, newModels);
