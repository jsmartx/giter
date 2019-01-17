#!/usr/bin/env node
'use strict';

const path = require('path');
const cp = require('child_process');

const binPath = path.join(__dirname, '../bin/giter');

cp.spawn(binPath, process.argv.slice(2), {
  cwd: process.cwd(),
  stdio: ['inherit', 'inherit', 'inherit']
});
