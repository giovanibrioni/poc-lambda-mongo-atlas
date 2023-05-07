const { MongoClient , ObjectId} = require('mongodb');
const client = new MongoClient(process.env.MONGODB_URI);
const express = require('express')
const serverless = require('serverless-http');
const app = express()

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


app.get('/users/:id', async (req, res) => {
  try {
    const id = req.params.id;

    if (!ObjectId.isValid(id)) {
      return res.status(400).send({ message: 'Invalid ID' });
    }

    const user = await client.db('test').collection('usuario').findOne({ _id: new ObjectId(id) });

    if (!user) {
      return res.status(404).send({ message: 'user not found' });
    }

    delete user._class;
    delete user.senha;

    return res.status(200).json({data: user});
  } catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' } );
  }
});

app.get('/health', async (req, res) => {
try {
    const result = await client.db('test').command({ ping: 1 });
    if (!result.ok) {
      return res.status(500).send({ message: 'Internal server error' });
    }
    return res.send( {message: 'ok'} );
} catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' } );
}
});


