#!/usr/bin/env node

var gettextParser = require("gettext-parser");

const args = process.argv.slice(2);

var input = require('fs').readFileSync(args[0]);
var po = gettextParser.po.parse(input);

Object.entries(po.translations[""]).map(([context, v]) => {
  if (v.msgid) {
    if (!v.comments) {
      v.comments = {};
    }
    v.comments.extracted = "From: " + (v.comments.reference || '');
  }
});

delete po.headers["PO-Revision-Date"];
delete po.headers["Language"];
delete po.headers["Plural-Forms"];
delete po.headers["mime-version"];
po.headers["MIME-Version"] = "1.0";

var output = gettextParser.po.compile(po);
require('fs').writeFileSync(args[1], output);
