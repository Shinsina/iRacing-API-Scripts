const output = require("../3-29-2024-search-series-output.json");
const fs = require("fs");

const test = Object.keys(output)
  .map((v) => output[v])
  .flat();
fs.writeFileSync(
  "./standings-input.json",
  JSON.stringify(
    Array.from(new Set(test.map((v) => `${v.season_id}_${v.car_class_id}`)))
  )
);
