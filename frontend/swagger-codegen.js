#!/bin/env node

import fs from "fs";
import path from "path";

const exceptions = ["change-owned-classroom-member-request"];
const teamOptionChange = [
  "user-classrooms",
  "get-owned-classroom-member-response",
  "get-joined-classroom-member-response",
  "get-joined-classroom-response",
];

const modelsFile = path.join("models", "index.ts");

const models = fs.readFileSync(modelsFile, "utf8");

let newModels = `import { DeepRequired, TeamPartial } from "@/types/utils";\n`;

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

      if (teamOptionChange.some((exception) => importFolder.includes(exception))) {
        editFunc = (input) => `TeamPartial<${input}>`;
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
