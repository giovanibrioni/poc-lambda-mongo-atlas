const { MongoClient , ObjectId} = require('mongodb');
const express = require('express')
const serverless = require('serverless-http');
const bodyParser = require('body-parser')


const User = require('./user');

const dbName = "test"
const collectionName = "users"
const client = new MongoClient(process.env.MONGODB_URI);

const app = express()
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: false }))


const lambdaName = process.env.AWS_LAMBDA_FUNCTION_NAME
if (!lambdaName) {
  const port = process.env.PORT || 8080
  app.listen(port, () => {
    console.log(`app is listening on port ${port}`)
  })
}
else {
  // Use serverless to run the app on AWS Lambda
  module.exports.handler = serverless(app);
}


app.get('/v1/user/:id', async (req, res) => {
  try {
    const id = req.params.id;

    if (!ObjectId.isValid(id)) {
      return res.status(400).send({ message: 'Invalid ID' });
    }

    const user = await client.db(dbName).collection(collectionName).findOne({ _id: new ObjectId(id) });

    if (!user) {
      return res.status(404).send({ message: 'user not found' });
    }

    user.id = user._id;
    delete user._id;

    return res.status(200).json({data: user});
    
  } catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' } );
  }
});

app.post('/v1/user', async (req, res) => {
  try {
    const user = new User(req.body);
    
    const result = await client.db(dbName).collection(collectionName).insertOne(user);
    if (result.insertedId) {
      user.id = result.insertedId;
      delete user._id;
      return res.status(201).json({data: user});
    }
    console.error(error);
    res.status(500).send('Internal server error');
  } catch (error) {
    console.error(error);
    return res.status(500).send('Internal server error');
  }
});


app.get('/health', async (req, res) => {
try {
    const result = await client.db(dbName).command({ ping: 1 });
    if (!result.ok) {
      return res.status(500).send({ message: 'Internal server error' });
    }
    return res.send( {message: 'ok'} );
} catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' } );
}
});

app.get('/ping', async (req, res) => {
  return res.send( {message: 'pong'} );
  });


