// Copyright 2016 Adam H. Leventhal. All rights reserved.
// Licensed under the Apache License, version 2.0:
// http://www.apache.org/licenses/LICENSE-2.0

'use strict';
let exec = require('child_process').exec;

exports.handler = (event, context) => {
  console.log('Received Event:', JSON.stringify(event, null, 2));

  const child = exec("./code_go_compiled " + event['nomspath'], (error, stdout, stderr) => {
    context.succeed(JSON.parse(stdout));
  });
};
