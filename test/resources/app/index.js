'use strict';

const path = require('path');
const fs = require('fs');
const http = require('http');

let data = ""
try {
  data = fs.readFileSync(path.join(__dirname, 'packaged_data'));
} catch (err) {
  // empty
}

const server = http.createServer(function (request, response) {
  response.writeHead(200, {"Content-Type": "text/plain"});
  response.end(data);
});

server.listen(8000);

console.log("Server running at http://127.0.0.1:8000/");
