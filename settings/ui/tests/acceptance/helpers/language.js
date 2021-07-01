const filesMenu = {
  English: [
    'All files',
    'Shared with me',
    'Shared with others',
    'Shared via link',
    'Deleted files'
  ],
  Deutsch: [
    'Alle Dateien',
    'Mit mir geteilt',
    'Mit anderen geteilt',
    'Per Link geteilt',
    'Gelöschte Dateien'
  ],
  Español: [
    'Todos los archivos',
    'Compartido conmigo',
    'Compartido con otros',
    'Shared via link',
    'Archivos borrados'
  ],
  Français: [
    'Tous les fichiers',
    'Partagé avec moi',
    'Partagé avec autres',
    'Shared via link',
    'Fichiers supprimés'
  ]
}

const accountMenu = {
  English: [
    'Profile',
    'Settings',
    'Log out'
  ],
  Deutsch: [
    'Profil',
    'Einstellungen',
    'Abmelden'
  ],
  Español: [
    'Profile',
    'Configuración',
    'Salir'
  ],
  Français: [
    'Profil',
    'Settings',
    'Se déconnecter'
  ]
}

const filesListHeaderMenu = {
  English: [
    'Name',
    'Size',
    'Modified',
    'Actions'
  ],
  Deutsch: [
    'Name',
    'Größe',
    'Geändert',
    'Aktionen'
  ],
  Español: [
    'Nombre',
    'Tamaño',
    'Modificado',
    'Acciones'
  ],
  Français: [
    'Nom',
    'Taille',
    'Modifié',
    'Actions'
  ]
}

exports.getFilesMenuForLanguage = function (language) {
  const menuList = filesMenu[language]
  if (menuList === undefined) {
    throw new Error(`Menu for language ${language} is not available`)
  }
  return menuList
}

exports.getUserMenuForLanguage = function (language) {
  const menuList = accountMenu[language]
  if (menuList === undefined) {
    throw new Error(`Menu for language ${language} is not available`)
  }
  return menuList
}

exports.getFilesHeaderMenuForLanguage = function (language) {
  const menuList = filesListHeaderMenu[language]
  if (menuList === undefined) {
    throw new Error(`Menu for language ${language} is not available`)
  }
  return menuList
}
