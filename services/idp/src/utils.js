export function withClientRequestState(obj) {
  obj.state = Math.random().toString(36).substring(7);

  return obj;
}
