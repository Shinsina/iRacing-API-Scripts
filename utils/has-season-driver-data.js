const output = require('../jack-standings-output.json');
const fs = require('fs');

fs.writeFileSync('./jack-standings-output.json', JSON.stringify(output.filter((v) => v.season_driver_data)));
