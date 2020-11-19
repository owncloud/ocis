const shell = require('shelljs')
const path = require("path")
const fs = require('fs');
const ora = require('ora');
const axios = require('axios');

const downloadFile = async (url) => {
  const parsedPath = path.parse(url)
  const destDir = './dist/_files/'
  const destFile = path.join(destDir, parsedPath.base)

  if (!fs.existsSync(destDir)) {
    shell.mkdir('-p', destDir)
  }

  if(fs.existsSync(destFile)){
    return
  }

  const spinner = ora(`downloading: ${ url }`).start();
  const { data } = await axios({
    method: "get",
    url: url,
    responseType: "stream"
  })

  const stream = fs.createWriteStream(destFile);

  data.pipe(stream);

  return new Promise((resolve, reject) => {
    data.on('error', err => {
      console.error(err);
      spinner.stop();
      reject(err);
    });

    data.on('end', () => {
      stream.end();
      spinner.stop();
      resolve();
    });
  });
}

(async () => {
  await downloadFile('http://ipv4.download.thinkbroadband.com/5MB.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/10MB.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/20MB.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/50MB.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/100MB.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/200MB.zip')
})()
