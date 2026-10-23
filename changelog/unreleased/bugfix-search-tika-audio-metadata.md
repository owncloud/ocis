Bugfix: Index the audio metadata in the formats tika reports

The search service dropped the duration of audio files because tika reports
it in seconds with a fractional part, for example `225.62`, while the
search service only accepted whole numbers. For mp3 files it also dropped
the track and disc numbers, which tika reports as `7/12`, and for files
with a full release date like `2004-05-12` it dropped the year.

The search service now reads these values in the formats tika reports, so
the duration, track, disc and year of audio files are indexed.

https://github.com/owncloud/ocis/pull/13204
