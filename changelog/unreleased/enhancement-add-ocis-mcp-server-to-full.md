Enhancement: Add oCIS MCP Server to the ocis_full deployment example

Added the oCIS MCP Server (https://github.com/owncloud/ocis-mcp-server) as an optional
service to the ocis_full deployment example. It exposes oCIS as a set of AI-accessible
tools over the Model Context Protocol, so AI assistants such as Claude can manage users,
spaces, files and shares through natural language without the user having to build and
run the MCP server locally. It is disabled by default and can be enabled by uncommenting
`OCIS_MCP=:mcp.yml` in `.env`.

https://github.com/owncloud/ocis/pull/12603
