import http from 'node:http'

// ---------------------------------------------------------------------------
// Config (all values come from environment variables)
// ---------------------------------------------------------------------------

const PORT = parseInt(process.env.PORT ?? '3030', 10)
const LLM_ENDPOINT = (process.env.LLM_ENDPOINT ?? '').replace(/\/$/, '')
const LLM_API_KEY = process.env.LLM_API_KEY ?? ''
const OCIS_URL = (process.env.OCIS_URL ?? '').replace(/\/$/, '')
/**
 * Origin (scheme://host:port) permitted to call the proxy, derived from
 * OCIS_URL. Empty when OCIS_URL is unset or unparseable, in which case the
 * origin gate cannot be enforced (the server refuses to start without
 * OCIS_URL outside tests, so this only matters in unit tests).
 */
const ALLOWED_ORIGIN = (() => {
  if (!OCIS_URL) return ''
  try {
    return new URL(OCIS_URL).origin
  } catch {
    return ''
  }
})()
/** When set, overrides the model field sent by the client. */
const LLM_MODEL = process.env.LLM_MODEL ?? ''
/** Hard ceiling on max_tokens forwarded to the LLM (default 4 096). */
const LLM_MAX_TOKENS = parseInt(process.env.LLM_MAX_TOKENS ?? '4096', 10)
/** Maximum request body the proxy will buffer in bytes (default 6 MiB — covers a 4 MiB image after base64 encoding). */
const MAX_BODY_BYTES = parseInt(process.env.MAX_BODY_BYTES ?? '6291456', 10)
/** Maximum LLM requests per user per rolling minute (default 20). */
const RATE_LIMIT_RPM = parseInt(process.env.RATE_LIMIT_RPM ?? '20', 10)
/**
 * How long to wait for the upstream LLM to finish a completion before
 * aborting (default 10 minutes). Local/self-hosted reasoning models (e.g.
 * Ollama) can legitimately take several minutes to generate a few hundred
 * tokens on modest hardware, with significant run-to-run variance (measured
 * 195s and 296s for near-identical completion lengths on the same machine),
 * so this must have real headroom rather than assume a fast, consistent
 * hosted API.
 */
const LLM_TIMEOUT_MS = parseInt(process.env.LLM_TIMEOUT_MS ?? '600000', 10)

// ---------------------------------------------------------------------------
// OIDC discovery — lazily fetched on first request, userinfo_endpoint cached
// ---------------------------------------------------------------------------

let userinfoEndpoint: string | null = null

async function resolveUserinfoEndpoint(): Promise<string> {
  if (userinfoEndpoint) return userinfoEndpoint
  const discoveryUrl = `${OCIS_URL}/.well-known/openid-configuration`
  const res = await fetch(discoveryUrl)
  if (!res.ok) {
    throw new Error(`OIDC discovery failed: ${res.status} ${res.statusText}`)
  }
  const doc = (await res.json()) as { userinfo_endpoint?: string }
  if (!doc.userinfo_endpoint) {
    throw new Error('OIDC discovery document missing userinfo_endpoint')
  }
  userinfoEndpoint = doc.userinfo_endpoint
  console.log(`[ai-llm-proxy] userinfo_endpoint: ${userinfoEndpoint}`)
  return userinfoEndpoint
}

// ---------------------------------------------------------------------------
// Token validation — returns the user's `sub` claim, or null if invalid
// ---------------------------------------------------------------------------

async function validateOcisToken(authorizationHeader: string): Promise<string | null> {
  const endpoint = await resolveUserinfoEndpoint()
  const res = await fetch(endpoint, {
    headers: { Authorization: authorizationHeader }
  })
  if (!res.ok) return null
  const info = (await res.json()) as { sub?: string }
  return typeof info.sub === 'string' ? info.sub : null
}

// ---------------------------------------------------------------------------
// Per-user sliding-window rate limiter
// ---------------------------------------------------------------------------

/** Exported so tests can reset state between runs. */
export const rateLimitWindows = new Map<string, number[]>()

/**
 * Returns true when the request is within the allowed rate, false when the
 * user has exceeded `rpm` requests in the last 60 seconds.
 * The `rpm` parameter defaults to the module-level RATE_LIMIT_RPM constant;
 * pass an explicit value in tests to avoid coupling to the env-var default.
 */
export function checkRateLimit(sub: string, rpm = RATE_LIMIT_RPM): boolean {
  const now = Date.now()
  const timestamps = (rateLimitWindows.get(sub) ?? []).filter((t) => now - t < 60_000)
  if (timestamps.length >= rpm) return false
  timestamps.push(now)
  rateLimitWindows.set(sub, timestamps)
  return true
}

// ---------------------------------------------------------------------------
// Body helpers
// ---------------------------------------------------------------------------

export class BodyTooLargeError extends Error {}

/**
 * Reads the request body, enforcing a hard byte limit.
 * Throws BodyTooLargeError when the limit is exceeded.
 */
export function readBody(req: http.IncomingMessage, maxBytes: number): Promise<string> {
  return new Promise((resolve, reject) => {
    const chunks: Buffer[] = []
    let totalBytes = 0
    req.on('data', (chunk: Buffer) => {
      totalBytes += chunk.length
      if (totalBytes > maxBytes) {
        reject(new BodyTooLargeError('Request body too large'))
        return
      }
      chunks.push(chunk)
    })
    req.on('end', () => resolve(Buffer.concat(chunks).toString('utf8')))
    req.on('error', reject)
  })
}

// ---------------------------------------------------------------------------
// Body sanitization
// ---------------------------------------------------------------------------

interface ClientBody {
  model?: unknown
  messages?: unknown
  max_tokens?: unknown
  stream?: unknown
  temperature?: unknown
  [key: string]: unknown
}

export interface SanitizedBody {
  model: string
  messages: unknown[]
  max_tokens: number
  temperature?: number
}

/**
 * Validates and sanitizes the client JSON body. Returns the clean object or
 * an `{ error }` descriptor when the body is unacceptable.
 *
 * - `modelOverride`: when non-empty, replaces whatever model the client sent
 * - `maxTokensLimit`: client's max_tokens is silently clamped to this ceiling
 *
 * Only the fields the LLM actually needs are forwarded; everything else is
 * dropped to prevent the client from injecting unexpected parameters.
 */
export function sanitizeBody(
  raw: ClientBody,
  modelOverride: string,
  maxTokensLimit: number
): SanitizedBody | { error: string } {
  const model = modelOverride || (typeof raw.model === 'string' ? raw.model : '')
  if (!model) return { error: 'model is required' }

  if (!Array.isArray(raw.messages)) return { error: 'messages must be an array' }

  const requestedTokens =
    typeof raw.max_tokens === 'number' && raw.max_tokens > 0 ? raw.max_tokens : maxTokensLimit

  const result: SanitizedBody = {
    model,
    messages: raw.messages as unknown[],
    max_tokens: Math.min(requestedTokens, maxTokensLimit)
  }

  if (typeof raw.temperature === 'number') result.temperature = raw.temperature

  return result
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

function setCorsHeaders(res: http.ServerResponse): void {
  res.setHeader('Access-Control-Allow-Origin', OCIS_URL || '*')
  res.setHeader('Access-Control-Allow-Methods', 'POST, OPTIONS')
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization')
  res.setHeader('Access-Control-Max-Age', '86400')
}

/**
 * Authoritative server-side origin check.
 *
 * CORS response headers (set by `setCorsHeaders`) are advisory and enforced
 * only by the browser, so they cannot stop a non-browser or hostile client
 * from reaching the LLM. This is the server-side gate that actually rejects
 * requests from an unexpected origin.
 *
 * A cross-origin browser request always carries an `Origin` header (and so
 * does a same-origin non-GET request, which is what the oCIS web extensions
 * issue), so a genuine call from the oCIS UI always presents
 * `Origin: <OCIS_URL origin>`. Any other present origin is a cross-site
 * request and is rejected. A request with no `Origin` header cannot be a
 * cross-site browser request and is still gated by the mandatory oCIS bearer
 * token, so it is allowed through to token validation.
 *
 * @param origin   the incoming request's `Origin` header value, if any
 * @param expected the allowed origin (scheme://host:port); empty disables the check
 * @returns true when the request may proceed
 */
export function isOriginAllowed(origin: string | undefined, expected: string): boolean {
  if (!expected) return true
  if (origin !== undefined && origin !== expected) return false
  return true
}

function sendJson(res: http.ServerResponse, status: number, body: unknown): void {
  const payload = JSON.stringify(body)
  res.writeHead(status, { 'Content-Type': 'application/json' })
  res.end(payload)
}

// ---------------------------------------------------------------------------
// Request handler
// ---------------------------------------------------------------------------

async function handleRequest(req: http.IncomingMessage, res: http.ServerResponse): Promise<void> {
  setCorsHeaders(res)

  // Preflight
  if (req.method === 'OPTIONS') {
    res.writeHead(204)
    res.end()
    return
  }

  // Only POST /v1/chat/completions
  if (req.method !== 'POST' || req.url !== '/v1/chat/completions') {
    sendJson(res, 404, { error: 'Not found' })
    return
  }

  // Authoritative server-side origin check (CORS headers above are advisory).
  // Reject cross-site requests before doing any work.
  const originHeader = req.headers['origin']
  const origin = Array.isArray(originHeader) ? originHeader[0] : originHeader
  if (!isOriginAllowed(origin, ALLOWED_ORIGIN)) {
    sendJson(res, 403, { error: 'Origin not allowed' })
    return
  }

  const authHeader = req.headers['authorization']
  if (!authHeader) {
    sendJson(res, 401, { error: 'Missing Authorization header' })
    return
  }

  let userSub: string | null
  try {
    userSub = await validateOcisToken(authHeader)
  } catch (err) {
    console.error('[ai-llm-proxy] Token validation error:', err)
    sendJson(res, 502, { error: 'Could not reach oCIS to validate token' })
    return
  }

  if (!userSub) {
    sendJson(res, 401, { error: 'Invalid or expired oCIS token' })
    return
  }

  if (!checkRateLimit(userSub)) {
    sendJson(res, 429, { error: 'Rate limit exceeded. Please slow down.' })
    return
  }

  let rawBodyStr: string
  try {
    rawBodyStr = await readBody(req, MAX_BODY_BYTES)
  } catch (err) {
    if (err instanceof BodyTooLargeError) {
      sendJson(res, 413, { error: 'Request body too large' })
    } else {
      sendJson(res, 400, { error: 'Could not read request body' })
    }
    return
  }

  let parsed: ClientBody
  try {
    parsed = JSON.parse(rawBodyStr) as ClientBody
  } catch {
    sendJson(res, 400, { error: 'Invalid JSON body' })
    return
  }

  const sanitized = sanitizeBody(parsed, LLM_MODEL, LLM_MAX_TOKENS)
  if ('error' in sanitized) {
    sendJson(res, 400, { error: sanitized.error })
    return
  }

  const llmHeaders: Record<string, string> = {
    'Content-Type': 'application/json'
  }
  if (LLM_API_KEY) {
    llmHeaders['Authorization'] = `Bearer ${LLM_API_KEY}`
  }

  let llmRes: Response
  try {
    llmRes = await fetch(`${LLM_ENDPOINT}/chat/completions`, {
      method: 'POST',
      headers: llmHeaders,
      body: JSON.stringify(sanitized),
      signal: AbortSignal.timeout(LLM_TIMEOUT_MS)
    })
  } catch (err) {
    console.error('[ai-llm-proxy] LLM request error:', err)
    sendJson(res, 502, { error: 'Could not reach LLM endpoint' })
    return
  }

  const llmBody = await llmRes.text()
  res.writeHead(llmRes.status, {
    'Content-Type': llmRes.headers.get('content-type') ?? 'application/json'
  })
  res.end(llmBody)
}

// ---------------------------------------------------------------------------
// Server — only started when run directly, not when imported in tests
// ---------------------------------------------------------------------------

if (process.env.NODE_ENV !== 'test') {
  if (!LLM_ENDPOINT) {
    console.error('LLM_ENDPOINT is required')
    process.exit(1)
  }
  if (!OCIS_URL) {
    console.error('OCIS_URL is required')
    process.exit(1)
  }

  const server = http.createServer((req, res) => {
    handleRequest(req, res).catch((err) => {
      console.error('[ai-llm-proxy] Unhandled error:', err)
      if (!res.headersSent) {
        sendJson(res, 500, { error: 'Internal server error' })
      }
    })
  })

  server.listen(PORT, () => {
    console.log(`[ai-llm-proxy] Listening on: ${PORT}`)
    console.log(`[ai-llm-proxy] LLM endpoint: ${LLM_ENDPOINT}`)
    console.log(`[ai-llm-proxy] oCIS URL:     ${OCIS_URL}`)
    console.log(`[ai-llm-proxy] LLM timeout:  ${LLM_TIMEOUT_MS}ms`)
  })
}
