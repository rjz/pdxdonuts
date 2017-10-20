const path = require('path');
const fs = require('fs');
const mkdirp = require('mkdirp');
const browserify = require('browserify');

const t0 = Date.now();

const root = path.resolve(__dirname, '..');
const src = path.join(root, 'src');
const dist = path.join(root, 'dist');

const entry = path.join(src, 'index.js');
const bundle = path.join(dist, 'bundle.js');

mkdirp.sync(dist);

const bundleStream = fs.createWriteStream(bundle);

const st = browserify(entry, { debug: true })
  .transform('babelify', { presets: ['env'] })
  .bundle()
  .pipe(bundleStream);

bundleStream.on('close', () => {
  const dt = Date.now() - t0;
  const prettyEntry = path.relative(root, entry);
  const prettyBundle = path.relative(root, bundle);

  console.log(`OK: Bundled ${prettyEntry} => ${prettyBundle} (${dt} ms)`);
});
