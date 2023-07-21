const { MongoClient, ServerApiVersion } = require("mongodb");
const data = require("./standings-data.json");
require("dotenv").config()

const uri = process.env.MONGODB_URI || null;

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
      const standingsCol = await client.db("main").collection("standings");
      const bulkOps = data.map((item) => ({
        updateOne: {
          filter: { _id: `${item.season_id}_${item.car_class_id}_${item.season_driver_data.cust_id}` },
          update: {
            $set: {
              ...item,
            },
            $setOnInsert: {
              _id: `${item.season_id}_${item.car_class_id}_${item.season_driver_data.cust_id}`,
            },
          },
          upsert: true,
        }
      }));
      await standingsCol.bulkWrite(bulkOps);
      console.log("Standings Data Bulk Write Complete!")
    } finally {
      await client.close();
    }
  }
  run().catch(console.dir);
}
