Enhancement: Add search MediaType filter

Add filter MediaType filter shortcuts to search for specific document types.
For example, a search query mediatype:documents will search for files with the following mimetypes:

application/msword
MimeType:application/vnd.openxmlformats-officedocument.wordprocessingml.document
MimeType:application/vnd.oasis.opendocument.text
MimeType:text/plain
MimeType:text/markdown
MimeType:application/rtf
MimeType:application/vnd.apple.pages

besides the document shorthand, it also contains following:

* file
* folder
* document
* spreadsheet
* presentation
* pdf
* image
* video
* audio
* archive

## File

## Folder

## Document:

application/msword
application/vnd.openxmlformats-officedocument.wordprocessingml.document
application/vnd.oasis.opendocument.text
text/plain
text/markdown
application/rtf
application/vnd.apple.pages

## Spreadsheet:

application/vnd.ms-excel
application/vnd.oasis.opendocument.spreadsheet
text/csv
application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
application/vnd.oasis.opendocument.spreadsheet
application/vnd.apple.numbers

## Presentations:

application/vnd.ms-powerpoint
application/vnd.openxmlformats-officedocument.presentationml.presentation
application/vnd.oasis.opendocument.presentation
application/vnd.apple.keynote

## PDF

application/pdf

## Image:

image/*

## Video:

video/*

## Audio:

audio/*

## Archive (zip ...):

application/zip
application/x-tar
application/x-gzip
application/x-7z-compressed
application/x-rar-compressed
application/x-bzip2
application/x-bzip
application/x-tgz

https://github.com/owncloud/ocis/pull/7602
https://github.com/owncloud/ocis/issues/7432
