import {
  RECEIVE_ERROR,
  RESET_HELLO,
  RECEIVE_HELLO,
  SERVICE_WORKER_NEW_CONTENT
} from '../actions/types';
import queryString from 'query-string';

const query = queryString.parse(document.location.search);
const flow = query.flow || '';
delete query.flow;

const defaultPathPrefix = (() => {
  const root = document.getElementById('root');
  let pathPrefix = root ? root.getAttribute('data-path-prefix') : null;
  if (!pathPrefix || pathPrefix === '__PATH_PREFIX__') {
    // Not replaced, probably we are running in debug mode or whatever. Use sane default.
    pathPrefix = '/signin/v1';
  }
  return pathPrefix;
})();

const defaultPasswordResetLink = (() => {
  const root = document.getElementById('root');
  let link = root ? root.getAttribute('passwort-reset-link') : null;
  if (!link || link === '__PASSWORD_RESET_LINK__') {
    // Not replaced, probably we are running in debug mode or whatever. Use sane default.
   link = '';
   }
  return link;
})();

const defaultState = {
  hello: null,
  branding: null,
  error: null,
  flow: flow,
  query: query,
  updateAvailable: false,
  pathPrefix: defaultPathPrefix,
  passwordResetLink: defaultPasswordResetLink
};

function commonReducer(state = defaultState, action) {
  switch (action.type) {
    case RECEIVE_ERROR:
      return Object.assign({}, state, {
        error: action.error
      });

    case RESET_HELLO:
      return Object.assign({}, state, {
        hello: null,
        branding: null
      });

    case RECEIVE_HELLO:
      return Object.assign({}, state, {
        hello: {
          state: action.state,
          username: action.username,
          displayName: action.displayName,
          details: action.hello
        },
        branding: action.hello.branding ? action.hello.branding : state.branding
      });

    case SERVICE_WORKER_NEW_CONTENT:
      return Object.assign({}, state, {
        updateAvailable: true
      });

    default:
      return state;
  }
}

export default commonReducer;
