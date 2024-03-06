const { MongoClient, ServerApiVersion } = require("mongodb");
const data = require("./3-5-2024-jake-output.json");
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
      const subsessionsCol = await client.db("main").collection("subsessions");
      const bulkOps = Object.keys(data).map((subsessionId) => ({
        updateOne: {
          filter: { _id: Number(subsessionId) },
          update: {
            $set: {
              ...data[subsessionId],
            },
            $setOnInsert: {
              _id: Number(subsessionId),
            },
          },
          upsert: true,
        }
      }));
      await subsessionsCol.bulkWrite(bulkOps);
      console.log("Subsessions Data Bulk Write Complete!")
    } finally {
      await client.close();
    }
  }
  run().catch(console.dir);
}
