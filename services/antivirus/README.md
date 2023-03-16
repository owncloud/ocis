# Antivirus Service

The `antivirus` service is responsible for scanning files for viruses

## Configuration

### Antivirus Scanner Type

The antivirus service currently supports `icap` and `clamav` as antivirus scanners. Use `ANTIVIRUS_SCANNER_TYPE` to configure this. 
Note that configuration depends heavily on chosen antivirus scanner. See Enviroment Variable descriptions for details.

### Maximum Scan size

Since several factors might make need necessary to limit the maximum filesize the `antivirus` service has an option to set a max scan size.
Use `ANTIVIRUS_MAX_SCAN_SIZE` to scan only that amount of bytes of a file. Obviously it is recommended to set this as high as possible, but several factors (scanner type and version, bandwith and performance issues, ...) might force to set this to a certain filesize.

### Infected File Handling

The `antivirus` service allows three different ways of handling infected files. Those can be set via the `ANTIVIRUS_INFECTED_FILE_HANDLING` envvar:
  -   `delete` (default): Infected files will be deleted immediately. Further postprocessing is cancelled.
  -   `abort`: Infected files will be kept. Further postprocessing is cancelled. Files can be manually retrieved and inspected by an admin. (Advanced option)
  -   `continue`: Infected files will be marked as infected but postprocessing continues normally. Note: Infected Files are not prevented from download. Risk of spreading viruses. (Obviously not recommended)

## Operation Modes

The `antivirus` service can scan files during postprocessing. `on demand` scanning will be added in the future.

### Postprocessing

Note: Needs to be configured via the [postprocessing service](https://github.com/owncloud/ocis/tree/master/services/postprocessing) 

The `antivirus` service will scan files during postprocessing. It listens for a postprocessing step called `"virusscan"`

### On Demand

On demand scanning is currently not supported
