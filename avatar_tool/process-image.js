const glob = require('glob');
const sharp = require('sharp');
const fs = require('fs');

const collection = {};

const rtxArr = []

const notJpgCollection = []

const promiseQueue = [];

glob.sync('./images/*').forEach((file, i) => {
    if(file.indexOf('.jpg') > -1) {
        console.log(`file: ${file}`);
        const rtx = /(\w+)(\(\w+\))*\.(jpg|JPG)/.exec(file)[1];
        promiseQueue.push(new Promise(resolve => {
            sharp(file).resize(300, 300).toBuffer().then(data => {
                console.log(i, file);
                const base64 = 'data:image/jpeg;base64,' + data.toString('base64');
                rtxArr.push(rtx.toLowerCase())
                collection[rtx.toLowerCase()] = base64;
                resolve();
            }).catch(err => console.error(err, file));
        }));
    }else{
        notJpgCollection.push(file)
    }
});

Promise.all(promiseQueue).then(() => {
    console.log('writing file...');
    console.log('notJpgCollection：',notJpgCollection)
    fs.writeFile('rtx.js', `var RTX = ${JSON.stringify(rtxArr)}`, (err, file) => {
        if (err) console.error(err);
    });
    fs.writeFile('avatar-300.js', `var AVATAR = ${JSON.stringify(collection)}`, (err, file) => {
        if (err) console.error(err);
    });
});
