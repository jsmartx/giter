const download = require('download');
const ora = require('ora');
const pkg = require('../package.json');

const PLATFORM = {
  'darwin': 'darwin',
  'freebsd': 'freebsd',
  'linux': 'linux',
  'openbsd': 'openbsd',
  'win32': 'windows'
};
const ARCH = {
  'ia32': '386',
  'x64': 'amd64',
  'x32': '386'
};

function install() {
  if (!(process.arch in ARCH)) {
    console.error('Installation is not supported for this architecture: ' + process.arch);
    return;
  }

  if (!(process.platform in PLATFORM)) {
    console.error('Installation is not supported for this platform: ' + process.platform);
    return;
  }
  const platform = PLATFORM[process.platform];
  const arch = ARCH[process.arch];
  const ghURL = `https://github.com/jsmartx/giter/releases/download/v${pkg.version}/`;
  const url = `${ghURL}${pkg.name}_${pkg.version}_${platform}_${arch}.tar.gz`;

  const spinner = ora(`Downloading ${url}`).start();
  download(url, 'bin', {
    extract: true
  }).then(() => {
    spinner.succeed();
  }).catch(() => {
    spinner.fail();
  });
}

install();
