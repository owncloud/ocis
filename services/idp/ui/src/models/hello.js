export function newHelloRequest(flow, query) {
  const r = {};

  if (query.prompt) {
    // TODO(longsleep): Validate prompt values?
    r.prompt = query.prompt;
  }

  let selectedFlow = flow;
  switch (flow) {
    case 'oauth':
    case 'consent':
    case 'oidc':
      r.scope = query.scope || '';
      r.client_id = query.client_id || ''; // eslint-disable-line camelcase
      r.redirect_uri = query.redirect_uri || '';  // eslint-disable-line camelcase
      if (query.id_token_hint) {
        r.id_token_hint = query.id_token_hint;  // eslint-disable-line camelcase
      }
      if (query.max_age) {
        r.max_age = query.max_age;  // eslint-disable-line camelcase
      }
      if (query.claims_scope) {
        // Add additional scopes from claims request if given.
        r.scope += ' ' + query.claims_scope;
      }
      break;

    default:
      selectedFlow = null;
  }

  if (selectedFlow) {
    r.flow = selectedFlow;
  }

  return r;
}
