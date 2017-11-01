const express = require('express');
const app = express();
const http = require('http');
const bodyParser = require('body-parser');

if (process.env.BACK_END) {
  console.log('setting backend to ' + process.env.BACK_END);
  var back_end = process.env.BACK_END;
}
else
{
  console.log('BACK_END not set. exiting')
  process.exit()
}

if (process.env.USERNAME) {
  console.log('setting username to ' + process.env.USERNAME);
  var username = process.env.USERNAME;
}
else
{
  console.log('USERNAME variable not set. setting to admin')
  var username = 'admin';
}

if (process.env.PASSWORD) {
  console.log('setting username to ' + process.env.PASSWORD);
  var password = process.env.PASSWORD;
}
else
{
  console.log('PASSWORD variable not set. setting to admin')
  var password = 'admin';
}

app.use(bodyParser.json()); // for parsing application/json

app.get('/', (req, res) => {

  var date = new Date().toISOString()

  res.send(date + "\n" + "success" + "\n");
});

app.post('/', (req, res) => {

  if (!req.body) return res.sendStatus(400);
  console.log(req.body);

  if (!req.body.username || !req.body.password)
  {
    console.log('username or password not set');
    res.sendStatus(403);
  }

  if (req.body.username != username || req.body.password != password) {
    console.log('username or password mismatch');
    res.sendStatus(403);
  }

  var data = JSON.stringify(req.body);

  var options = {
    host: back_end,
    port: '80',
    path: '/',
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    }
  };

  callback = function(response) {

    if (response.statusCode != 200) {
      res.send({"error": "didn't receive a 200: " + response.statusCode});
    }
    else
    {
      var str = ''
      response.on('data', function (chunk) {
        str += chunk;
      });

      response.on('end', function() {
        console.log("received from backend: " + str);
        if (str == data) {
          res.send({ status: 'SUCCESS' });
        }
        else
        {
          console.log("'" + str + "' doesnt match '" + data + "'");
          res.send({ status: 'FAILED' });
        }
      });
    }
  }

  console.log('sending ' + data + ' to ' + back_end);

  var backreq = http.request(options, callback);
  backreq.write(data);
  backreq.end();
});

app.listen(process.env.APP_PORT, () => {
  console.log('listening on ' + process.env.APP_PORT)
});
