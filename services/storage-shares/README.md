# Storage-Shares Service

The storage-shares service provides a virtual storage for shares in oCIS (ownCloud Infinite Scale). It acts as a storage provider that exposes shared resources through a unified interface, allowing users to access files and folders that have been shared with them.

## Overview

The storage-shares service is a crucial component of the oCIS sharing system. It implements a virtual storage provider that aggregates and presents shared resources to users. When files or folders are shared between users, this service makes them accessible through a consistent storage interface.

### Key Features

- **Virtual Storage Provider**: Provides a unified view of shared resources
- **Share Aggregation**: Collects and presents shares from the sharing service
- **GRPC Interface**: Exposes storage operations via GRPC protocol
- **Reva Integration**: Built on top of the Reva storage framework
- **Read-Only Mode**: Supports read-only operation for specific use cases
- **Prometheus Metrics**: Built-in monitoring and metrics collection

## Architecture

The service is built using the Reva framework and consists of several key components:

- **Storage Provider**: Implements the Reva storage provider interface for shares
- **GRPC Server**: Handles storage operations and metadata requests
- **HTTP Server**: Provides additional HTTP endpoints for debugging and health checks
- **Share Resolution**: Interfaces with the sharing service to resolve share information

## Configuration

The service can be configured through environment variables or configuration files. Key configuration options include:

- **GRPC Address**: Network address for the GRPC server (default: `127.0.0.1:9154`)
- **Debug Address**: Address for debug endpoints (default: `127.0.0.1:9156`)
- **Mount ID**: Unique identifier for the storage mount point
- **Shares Provider Endpoint**: Endpoint of the sharing service
- **Read-Only Mode**: Enable/disable read-only operation
- **JWT Secret**: Secret for token validation

For detailed configuration options, see the [configuration documentation](../../docs/services/storage-shares/configuration.md).

## Usage

### Starting the Service

The service can be started using the following command:

```bash
./storage-shares server
```

### Available Commands

- `server`: Start the storage-shares service
- `health`: Check service health status
- `version`: Display version information

### Health Check

The service provides health check endpoints for monitoring:

```bash
./storage-shares health
```

## Development

### Building

To build the service:

```bash
make build
```

### Testing

To run tests:

```bash
make test
```

### Documentation

To generate documentation:

```bash
make docs-generate
```

## Integration

The storage-shares service integrates with several other oCIS services:

- **Sharing Service**: Retrieves share information and permissions
- **Gateway Service**: Registers as a storage provider
- **Storage-Users Service**: Coordinates with user storage for shared content
- **Authentication Services**: Validates user tokens and permissions

## API

The service exposes a GRPC API following the Reva storage provider specification. Key operations include:

- **ListContainer**: List shared resources
- **GetPath**: Resolve paths within shared storage
- **Stat**: Get metadata for shared resources
- **InitiateFileDownload**: Initiate download of shared files

## Monitoring

The service includes built-in Prometheus metrics for monitoring:

- Request counts and latencies
- Error rates
- Storage operations metrics
- Share resolution performance

Metrics are exposed on the debug endpoint under `/metrics`.

## Security

The service implements several security measures:

- **JWT Token Validation**: All requests require valid authentication tokens
- **Permission Checking**: Verifies user permissions for shared resources
- **TLS Support**: Optional TLS encryption for GRPC communications
- **Read-Only Mode**: Prevents modifications when enabled

## Troubleshooting

### Common Issues

1. **Service Won't Start**: Check GRPC port availability and configuration
2. **Shares Not Visible**: Verify sharing service connectivity and configuration
3. **Permission Denied**: Check JWT secret configuration and token validity
4. **Performance Issues**: Monitor metrics and check sharing service performance

### Logging

The service provides structured logging with configurable levels. Enable debug logging for troubleshooting:

```bash
STORAGE_SHARES_LOG_LEVEL=debug ./storage-shares server
```

## Contributing

When contributing to the storage-shares service:

1. Follow the existing code structure and patterns
2. Add tests for new functionality
3. Update documentation as needed
4. Ensure compatibility with the Reva framework
5. Test integration with other oCIS services

## License

This service is part of oCIS and is licensed under the Apache License 2.0.