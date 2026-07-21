#!/usr/bin/env python3
"""Minimal HTTP server that serves hosting-discovery.xml on :8080."""
import http.server
import pathlib

body = pathlib.Path('/ocis/tests/config/ci/hosting-discovery.xml').read_bytes()


class Handler(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-Type', 'text/xml')
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *a):
        pass


http.server.HTTPServer(('', 8080), Handler).serve_forever()
