'use strict'

const express = require('express')
const bodyParser = require('body-parser')
const request = require('request')
const app = express()

app.set('port', (process.env.PORT || 5000))

// Process application/x-www-form-urlencoded
app.use(bodyParser.urlencoded({extended: false}))

// Process application/json
app.use(bodyParser.json())

// Index route
app.get('/', function (req, res) {
	res.send('Hello world, I am a chat bot')
})

// for Facebook verification
app.get('/webhook/', function (req, res) {
	if (req.query['hub.verify_token'] === process.env.VERIFY_TOKEN) {
		res.send(req.query['hub.challenge'])
	}
	res.send('Error, wrong token')
})

// Spin up the server
app.listen(app.get('port'), function() {
	console.log('running on port', app.get('port'))
})

app.post('/webhook/', function (req, res) {
    let messaging_events = req.body.entry[0].messaging
    for (let i = 0; i < messaging_events.length; i++) {
	    let event = req.body.entry[0].messaging[i]
	    let sender = event.sender.id
	    if (event.message && event.message.text) {
			let text = event.message.text
			// PARSE TEXT HERE
			let input = parseText(text);
			let inputRequest = callAPI(input);
			// POST TO SERVER
			// let response = responseHandle();

		    sendTextMessage(sender, inputRequest);
	    }
    }
    res.sendStatus(200)
});

function sendTextMessage(sender, text) {
    let messageData = { text:text }
    request({
	    url: 'https://graph.facebook.com/v2.6/me/messages',
	    qs: {access_token:token},
	    method: 'POST',
		json: {
		    recipient: {id:sender},
			message: messageData,
		}
	}, function(error, response, body) {
		if (error) {
		    console.log('Error sending messages: ', error)
		} else if (response.body.error) {
		    console.log('Error: ', response.body.error)
	    }
    })
}

// TODO: Eventually handle things that aren't text like stickers or images
function parseText(text) {
	var obj = {
		"input": text,
		"url": ""
	}
	return JSON.stringify(obj);
}

const token = process.env.FB_TOKEN;

function callAPI(userInput) {
	var https = require("https");
	var options = {
	  hostname: 'simplify.api.shaheensharifian.me',
	  path: '/v1/simplify/text',
	  method: 'POST',
	  headers: {
		  'Content-Type': 'application/json',
	  }
	};
	var req = https.request(options, function(res) {
		// console.log('Status: ' + res.statusCode);
		// console.log('Headers: ' + JSON.stringify(res.headers));
		res.setEncoding('utf8');
		res.on('data', function (body) {
		//   console.log('Body: ' + body);
			parseJson(body);
		});
	  });
	  req.on('error', function(e) {
		console.log('problem with request: ' + e.message);
	  });
	  // write data to request body
	//   req.write('{"string": "Hello, World"}');
	  req.write(userInput);
	  req.end();

}

function parseJson(json) {
	let responseObj = JSON.parse(json);
	responseObj.entities.forEach(function(entity) {
		console.log(entity.name);
	}, this);
}

var obj = parseText("What is the difference between linux kernel and shell?");
callAPI(obj);

// let test = {
// 	"input": "Hello there friend!",
// 	"url": ""
// }

// console.log(callAPI(test));

// function responseHandle(responseObj) {
// 	let json = JSON.parse(responseObj);

// }

