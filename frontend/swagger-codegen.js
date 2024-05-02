#!/bin/env node

import fs from "fs";
import path from "path";

const modelsFile = path.join("models", "index.ts");

const models = fs.readFileSync(modelsFile, "utf8");

let newModels = `import { DeepRequired } from "@/types/utils";\n`;

newModels +=
  models
    .split("\n")
    .filter(Boolean)
    .map((line) => {
      const parts = line.split(" ");
      const importFolder = parts[3].replace(/'\.\/(.*)';/g, "$1");
      const modelType = importFolder
        .split("-")
        .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
        .join("");

      const importLine = `import { ${modelType} as Old${modelType} } from "./${importFolder}";`;
      const exportLine = `export type ${modelType} = DeepRequired<Old${modelType}>;`;
      return [importLine, exportLine];
    })
    .flat()
    .join("\n") + "\n";

newModels = newModels.replaceAll("Httperror", "HTTPError");

fs.writeFileSync(modelsFile, newModels);
