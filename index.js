// 'use strict'

// const express = require('express')
// const bodyParser = require('body-parser')
// const request = require('request')
// const app = express()

// app.set('port', (process.env.PORT || 5000))

// // Process application/x-www-form-urlencoded
// app.use(bodyParser.urlencoded({extended: false}))

// // Process application/json
// app.use(bodyParser.json())

// // Index route
// app.get('/', function (req, res) {
// 	res.send('Hello world, I am a chat bot')
// })

// // for Facebook verification
// app.get('/webhook/', function (req, res) {
// 	if (req.query['hub.verify_token'] === 'hello world') {
// 		res.send(req.query['hub.challenge'])
// 	}
// 	res.send('Error, wrong token')
// })

// // Spin up the server
// app.listen(app.get('port'), function() {
// 	console.log('running on port', app.get('port'))
// })

// app.post('/webhook/', function (req, res) {
//     let messaging_events = req.body.entry[0].messaging
//     for (let i = 0; i < messaging_events.length; i++) {
// 		let event = req.body.entry[0].messaging[i]

// 		let sender = event.sender.id
// 		console.log("id" + sender.id);
// 	    if (event.message && event.message.text) {
// 			let text = event.message.text
// 			let input = parseMessage(text, "");
// 			let inputRequest = callAPI(input, sender);

// 		}	
// 		// Checking for attachments
// 		// if (event.message && event.message.attachments) {
// 		// 	let attachment = event.message.attachments[0];
// 		// 	// Checking if attachment is an image
// 		// 	if (attachment.type === "image") {
// 		// 		let url = attachment.payload.url;
// 		// 		console.log(url);
// 		// 		let picture = parseMessage("", url);
// 		// 		let inputRequest = callAPI(picture, sender);
// 		// 	} 
// 	//  }
//     }
//     res.sendStatus(200)
// });

// function sendTextMessage(sender, text) {
// 	console.log(sender);
// 	const token = "EAABwawgEld0BADVm3znFKiZA9jJjswSW1gNFUhqEFROOWi2qSOvQtZBtNCo1tOKdLif22pBlQ9KA5Nxio08a0smjGT44OGVqyQjUT5t0cYZABTaGNrKZCZCKdw9yWFmG8FQuktzjEZAehIvXvqJ3OdxTOGIHG3ZA9a7137eyspJZAwZDZD";
// 	console.log(token);
//     let messageData = { text:text }
//     request({
// 	    url: 'https://graph.facebook.com/v2.6/me/messages',
// 	    qs: {access_token:token},
// 	    method: 'POST',
// 		json: {
// 		    recipient: {id:sender},
// 			message: messageData,
// 		}
// 	}, function(error, response, body) {
// 		if (error) {
// 		    console.log('Error sending messages: ', error)
// 		} else if (response.body.error) {
// 		    console.log('Error: ', response.body.error)
// 	    }
//     })
// }

// // TODO: Eventually handle things that aren't text like stickers or images
// function parseMessage(text, url) {
// 	var obj = {
// 		"input": text,
// 		"url": url
// 	}
// 	return JSON.stringify(obj);
// }


// function callAPI(userInput, sender) {
// 	var https = require("https");
// 	var options = {
// 	  hostname: 'simplify.api.shaheensharifian.me',
// 	  path: '/v1/simplify/text',
// 	  method: 'POST',
// 	  headers: {
// 		  'Content-Type': 'application/json',
// 	  }
// 	};
	
// 	var req = https.request(options, function(res) {
// 		// console.log('Status: ' + res.statusCode);
// 		// console.log('Headers: ' + JSON.stringify(res.headers));
// 		res.setEncoding('utf8');
// 		var response = "";		
// 		var shit = "";
// 		res.on('data', function(body) {
// 		//   console.log('Body: ' + body);
// 			// console.log(body);
// 			shit += body;
// 			// response += parseJson(body);
// 		});
// 		res.on('end', function() {

// 			sendTextMessage(sender, parseJson(shit));
// 		  });
// 	});
// 	req.on('error', function(e) {
// 	console.log('problem with request: ' + e.message);
// 	});
// 	// write data to request body
// //   req.write('{"string": "Hello, World"}');
// 	req.write(userInput);
// 	req.end();
// }

// function parseJson(json) {
// 	let responseObj = JSON.parse(json);
// 	var array = [];
// 	responseObj.entities.forEach(function(entity) {
// 		var keyWord = entity.name;
// 		// var wordDef = entity.definition;
// 		// var wordAndDef = "This is the definition of " + keyWord + ":\n" + wordDef;
// 		array.push(keyWord);
// 	}, this);
// 	return array.join("-");
// }

// // var obj = parseMessage("", "http://cdn2-www.dogtime.com/assets/uploads/gallery/shiba-inu-puppies/shiba-inu-puppy-13.jpg");
// // var x = callAPI(obj);
// // console.log(x);

// // var shit = parseMessage("What is linux kernel?", "");
// // var callShit = callAPI(shit, );
// // console.log(callShit);

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
	if (req.query['hub.verify_token'] === 'my_voice_is_my_password_verify_me') {
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
		    sendTextMessage(sender, "Text received, echo: " + text.substring(0, 200))
	    }
    }
    res.sendStatus(200)
})

const token = "EAAQOOGWd6nEBAMiZAGLXLwEveRNhPOJQNtBDnFxDC6yS31TkrEhQgOZCHbAlqwZBBMI4nav152XaJBSrj0uElg91zHczZCh3ZC9nJYi1NB8SAZAyDhZBlVTKmi1t9FvJHh3OfmZAAHJgq1XuckmGn2lEagLYJK8npK4yyQvah9lurQZDZD";

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