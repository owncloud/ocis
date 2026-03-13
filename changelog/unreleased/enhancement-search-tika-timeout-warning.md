Enhancement: Warn about Tika server timeout mismatch at startup

The search service now checks the Tika server's taskTimeoutMillis at
startup and logs a warning if it appears to be at the default (120s),
which is lower than the typical tesseract OCR parser timeout (300s).
This mismatch causes the Tika watchdog to kill child processes during
legitimate long-running OCR, leading to cascading restarts and
extraction failures.

https://github.com/owncloud/ocis/pull/XXXXX
