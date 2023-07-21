const { MongoClient, ServerApiVersion } = require("mongodb");
const data = require("./7-21-2023-car-class-input.json");
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
      const seasonsCol = await client.db("main").collection("carclasses");
      const bulkOps = data.map((item) => ({
        updateOne: {
          filter: { _id: Number(item.car_class_id) },
          update: {
            $set: {
              ...item,
            },
            $setOnInsert: {
              _id: Number(item.car_class_id),
            },
          },
          upsert: true,
        }
      }));
      await seasonsCol.bulkWrite(bulkOps);
      console.log("Car Class Data Bulk Write Complete!")
    } finally {
      await client.close();
    }
  }
  run().catch(console.dir);
}
