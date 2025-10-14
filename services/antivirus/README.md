# Antivirus

The `antivirus` service is responsible for scanning files for viruses.

## Memory Considerations

The antivirus service can consume considerably amounts of memory. This is relevant to provide or define sufficient memory for the deployment selected. To avoid out of memory (OOM) situations, the following equation gives a rough overview based on experiences made. The memory calculation comes without any guarantee, is intended as overview only and subject of change.

`memory limit` = `max file size` x `workers` x `factor 8 - 14`

With:
`ANTIVIRUS_WORKERS` == 1
```plaintext
 50MB file --> factor 14   --> 700MB memory
844MB file --> factor  8,3 -->   7GB memory
```

## Configuration

### Antivirus Scanner Type

The antivirus service currently supports [ICAP](https://tools.ietf.org/html/rfc3507) and [ClamAV](http://www.clamav.net/index.html) as antivirus scanners. The `ANTIVIRUS_SCANNER_TYPE` environment variable is used to select the scanner. The detailed configuration for each scanner heavily depends on the scanner type selected. See the environment variables for more details.

  -   For `icap`, only scanners using the `X-Infection-Found` header are currently supported.
  -   For `clamav` only local sockets can currently be configured.

### Maximum Scan Size

Several factors can make it necessary to limit the maximum filesize the antivirus service will use for scanning. Use the `ANTIVIRUS_MAX_SCAN_SIZE` environment variable to scan only a given amount of bytes. Obviously, it is recommended to scan the whole file, but several factors like scanner type and version, bandwidth, performance issues, etc. might make a limit necessary.

**IMPORTANT**
> Streaming of files to the virus scan service still [needs to be implemented](https://github.com/owncloud/ocis/issues/6803). To prevent OOM errors `ANTIVIRUS_MAX_SCAN_SIZE` needs to be set lower than available ram.

### Antivirus Workers

The number of concurrent scans can be increased by setting `ANTIVIRUS_WORKERS`. Be aware that this will also increase memory usage.

### Infected File Handling

The antivirus service allows three different ways of handling infected files. Those can be set via the `ANTIVIRUS_INFECTED_FILE_HANDLING` environment variable:

  -   `delete`: (default): Infected files will be deleted immediately, further postprocessing is cancelled.
  -   `abort`:  (advanced option): Infected files will be kept, further postprocessing is cancelled. Files can be manually retrieved and inspected by an admin. To identify the file for further investigation, the antivirus service logs the abort/infected state including the file ID. The file is located in the `storage/users/uploads` folder of the ocis data directory and persists until it is manually deleted by the admin via the [Manage Unfinished Uploads](https://doc.owncloud.com/ocis/next/deployment/services/s-list/storage-users.html#manage-unfinished-uploads) command.
  -   `continue`:  (obviously not recommended): Infected files will be marked via metadata as infected but postprocessing continues normally. Note: Infected Files are moved to their final destination and therefore not prevented from download which includes the risk of spreading viruses.

In all cases, a log entry is added declaring the infection and handling method and a notification via the `userlog` service sent.

### Scanner Inaccessibility

In case a scanner is not accessible by the antivirus service like a network outage, service outage or hardware outage, the antivirus service uses the `abort` case for further processing, independent of the actual setting made. In any case, an error is logged noting the inaccessibility of the scanner used.

## Operation Modes

The antivirus service can scan files during `postprocessing`. `on demand` scanning is currently not available and might be added in a future release.

### Postprocessing

The antivirus service will scan files during postprocessing. It listens for a postprocessing step called `virusscan`. This step can be added in the environment variable `POSTPROCESSING_STEPS`. Read the documentation of the [postprocessing service](https://github.com/owncloud/ocis/tree/master/services/postprocessing) for more details.

The number of concurrent scans can be increased by setting `ANTIVIRUS_WORKERS`, but be aware that this will also increase the memory usage.

### Scaling in Kubernetes

In kubernetes, `ANTIVIRUS_WORKERS` and `ANTIVIRUS_MAX_SCAN_SIZE` can be used to trigger the horizontal pod autoscaler by requesting a memory size that is below `ANTIVIRUS_MAX_SCAN_SIZE`. Keep in mind that `ANTIVIRUS_MAX_SCAN_SIZE` amount of memory might be held by `ANTIVIRUS_WORKERS` number of go routines.
