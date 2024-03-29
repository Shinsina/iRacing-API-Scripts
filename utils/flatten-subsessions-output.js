const output = require("./3-29-2024-subsessions-output.json");
const fs = require("fs");

const test = Object.keys(output)
  .map((v) => output[v])
  .flat();
fs.writeFileSync("./flattened-subsessions-output.json", JSON.stringify(test));
