Bugfix: Restart Postprocessing

In case the postprocessing service cannot find the specified upload when restarting postprocessing, it will now send a
`RestartPostprocessing` event to retrigger complete postprocessing

https://github.com/owncloud/ocis/pull/6726
