#!/usr/bin/env node

const { program } = require('commander');

program.command('m3u8Decrypt')
    .argument('<string>', 'string to split')
    .action((str, options) => {
        // const urls = decrypt(str)
        // const { url, ext } = urls[0]
        // const result = {
        //     plain_url: url,
        //     ts_url_part: {
        //         host: ext.host.substring(ext.host.indexOf("://") + 3),
        //         path: ext.path,
        //         query: ext.param,
        //     }
        // }
        //
        // console.log(JSON.stringify(result));
    });

program.command("keyURLWrapper").argument('<plainKey>', 'key uri')
    .action((uri, options) => {
        console.log(uri)
    })

program.command("keyDecrypt").argument('<plainKey>', 'plainKey')
    .action((keyBase64, options) => {
        const KeyBuf = Buffer.from(keyBase64, 'base64')
        console.log(Buffer.from(KeyBuf).toString('base64'));
    })

program.parse();