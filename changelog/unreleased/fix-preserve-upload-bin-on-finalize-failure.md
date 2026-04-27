Bugfix: Preserve upload bin when synchronous blobstore finalization fails

When a synchronous upload finalization failed (e.g. due to a transient NFS error during
WriteBlob), decomposedfs deleted the bin file from uploads/ unconditionally. This made
the blob permanently unrecoverable — including via move-stuck-upload-blobs — and caused
affected nodes (e.g. received.json in the share manager metadata space) to fail on every
access with "no such file or directory".

Fixed by bumping reva to stable-8.0 HEAD (488fb352f5), which only cleans the bin file
when Finalize succeeds.

https://github.com/owncloud/ocis/pull/12264
