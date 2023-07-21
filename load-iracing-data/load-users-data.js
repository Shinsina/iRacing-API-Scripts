const { MongoClient, ServerApiVersion } = require("mongodb");
const data = require("./7-21-2023-jake-user-input.json");
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
      const seasonsCol = await client.db("main").collection("users");
      const bulkOps = data.map((item) => ({
        updateOne: {
          filter: { _id: Number(item.cust_id) },
          update: {
            $set: {
              ...item,
            },
            $setOnInsert: {
              _id: Number(item.cust_id),
            },
          },
          upsert: true,
        }
      }));
      await seasonsCol.bulkWrite(bulkOps);
      console.log("User Data Bulk Write Complete!")
    } finally {
      await client.close();
    }
  }
  run().catch(console.dir);
}
