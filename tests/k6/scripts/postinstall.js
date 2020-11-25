const shell = require('shelljs')
const path = require("path")
const fs = require('fs');
const ora = require('ora');
const axios = require('axios');

const downloadFile = async (url, name) => {
  const parsedPath = path.parse(url)
  const destDir = './dist/_files/'
  const destFile = path.join(destDir, name || parsedPath.base)

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
  });
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
  await downloadFile('https://www.sample-videos.com/img/Sample-jpg-image-50kb.jpg', 'kb_50.jpg')
  await downloadFile('http://ipv4.download.thinkbroadband.com/5MB.zip', 'mb_5.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/10MB.zip', 'mb_10.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/20MB.zip', 'mb_20.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/50MB.zip', 'mb_50.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/100MB.zip', 'mb_100.zip')
  await downloadFile('http://ipv4.download.thinkbroadband.com/200MB.zip', 'mb_200.zip')
})()
