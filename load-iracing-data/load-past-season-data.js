const { MongoClient, ServerApiVersion } = require("mongodb");
const data = require("./past-seasons-output.json");
const seasonIds = require("./past-seasons-season-ids-input.json");
require("dotenv").config()

const uri = process.env.MONGODB_URI || null;
const seasonIdsSet = new Set(seasonIds);

if (!uri) {
  console.log("Please set MONGODB_URI with a valid Mongo URI within this folder in a .env file!")
} else {
  const client = new MongoClient(uri, {
    serverApi: {
      version: ServerApiVersion.v1,
      strict: true,
      deprecationErrors: true,
    }
  });
  async function run() {
    try {
      await client.connect();
      const pastSeasonsCol = await client.db("main").collection("pastseasons");
      const items = Object.keys(data)
        .map((key) => data[key].series.seasons.filter((season) => seasonIdsSet.has(season.season_id))).flat();
      const bulkOps = items.map((item) => ({
        updateOne: {
          filter: { _id: Number(item.season_id) },
          update: {
            $set: {
              ...item,
            },
            $setOnInsert: {
              _id: Number(item.season_id),
            },
          },
          upsert: true,
        }
      }));
      await pastSeasonsCol.bulkWrite(bulkOps);
      console.log("Past Season Data Bulk Write Complete!")
    } finally {
      await client.close();
    }
  }
  run().catch(console.dir);
}
