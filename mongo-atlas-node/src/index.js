const { MongoClient , ObjectId} = require('mongodb');
const express = require('express')
const serverless = require('serverless-http');
const bodyParser = require('body-parser')


const { User, userSchema } = require('./user');

const dbName = "test"
const collectionName = "users"
const client = new MongoClient(process.env.MONGODB_URI);

const app = express()
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: false }))


const lambdaName = process.env.AWS_LAMBDA_FUNCTION_NAME
if (!lambdaName) {
  const port = process.env.PORT || 3000
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

    let user = await client.db(dbName).collection(collectionName).findOne({ _id: new ObjectId(id) });

    if (!user) {
      return res.status(404).send({ message: 'user not found' });
    }

    user = changeId(user)

    return res.status(200).json({data: user});
    
  } catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' } );
  }
});

app.post('/v1/user', async (req, res) => {
  try {
    const { error, value } = userSchema.validate(req.body);
     if (error) {
       return res.status(400).json({ message: error.details[0].message });
     }
    let user = new User(value);
    
    const result = await client.db(dbName).collection(collectionName).insertOne(user);
    if (result.insertedId) {
      user = changeId(user)
      return res.status(201).json({data: user});
    }
    console.error(error);
    res.status(500).send('Internal server error');
  } catch (error) {
    console.error(error);
    return res.status(500).send('Internal server error');
  }
});

app.get('/v1/user', async (req, res) => {
  try {
    // Parse the query parameters for page and limit
    const page = parseInt(req.query.page) || 1;
    const limit = parseInt(req.query.limit) || 10;

    // Calculate the offset and skip values
    const offset = (page - 1) * limit;
    const skip = offset;

    // Find all users and apply pagination
    const users = await client.db(dbName).collection(collectionName).find().skip(skip).limit(limit).toArray();

    // Get the total number of users
    const count = await client.db(dbName).collection(collectionName).countDocuments();

    // Calculate the total number of pages
    const pages = Math.ceil(count / limit);

    users.forEach(user => {
      user = changeId(user)
    });

    return res.status(200).json({
      data: users,
      page,
      limit,
      pages,
      count,
    });

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

const changeId = (user) => {
  user.id = user._id;
  delete user._id;
  return user
}