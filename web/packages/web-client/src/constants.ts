export const PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION = 'psec'

/**
 * Query param set on the public link opened inside the password-protected-folder view
 * modal. It tells the framed app instance that it is running inside that modal, so it can
 * notify the parent window when the shared folder is renamed (see the message name below).
 */
export const PASSWORD_PROTECTED_FOLDER_VIEW_QUERY = 'password-protected-folder-view'

/**
 * postMessage name emitted by the framed public-link app when the password-protected folder
 * (the public link root) is renamed inside the modal, so the parent window can keep the
 * `.psec` pointer file in sync. The `.psec` file cannot be resolved from the folder's new
 * name (the coupling is name-based), so the new name is delivered here at rename time.
 */
export const PASSWORD_PROTECTED_FOLDER_RENAMED_MESSAGE =
  'owncloud-password-protected-folder:renamed'

/**
 * List of file extensions that should be hidden from the user.
 * Hiding the extension currently leads to hiding all actions except delete.
 */
export const HIDDEN_FILE_EXTENSIONS = [PASSWORD_PROTECTED_FOLDER_FILE_EXTENSION]
