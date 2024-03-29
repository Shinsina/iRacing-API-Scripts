const output = require('../3-29-2024-search-series-output.json');
const fs = require('fs');

fs.writeFileSync('./distinct-series-ids-output.json', JSON.stringify(Array.from(new Set(output.map((v) => v.series_id)))));
fs.writeFileSync('./distinct-season-ids-output.json', JSON.stringify(Array.from(new Set(output.map((v) => v.season_id)))));
