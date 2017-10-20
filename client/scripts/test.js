const path = require('path');
const glob = require('glob');

const root = path.resolve(__dirname, '..');
const fileGlob = path.resolve(root, 'src/**/*_test.js');

glob.sync(fileGlob).forEach(require);

