const filesMenu = {
  English: [
    'Personal',
    'Shares',
    'Spaces',
    'Deleted files'
  ],
  Deutsch: [
    'Persönlich',
    'Geteilt',
    'Spaces',
    'Gelöschte Dateien'
  ],
  Español: [
    'Personal',
    'Shares',
    'Spaces',
    'Archivos borrados'
  ],
  Français: [
    'Personal',
    'Shares',
    'Spaces',
    'Fichiers supprimés'
  ]
}

const accountMenu = {
  English: [
    'U\nuser1\nuser1@example.org',
    'Settings',
    'Log out',
    'Personal storage\n0 B used'
  ],
  Deutsch: [
    'U\nuser1\nuser1@example.org',
    'Einstellungen',
    'Abmelden',
    'Persönlicher Speicherplatz\n0 B verwendet'
  ],
  Español: [
    'U\nuser1\nuser1@example.org',
    'Configuración',
    'Salir',
    'Personal storage\n0 B used'
  ],
  Français: [
    'U\nuser1\nuser1@example.org',
    'Settings',
    'Se déconnecter',
    'Personal storage\n0 B used'
  ]
}

const filesListHeaderMenu = {
  English: [
    'Name',
    'Shares',
    'Size',
    'Modified',
    'Actions'
  ],
  Deutsch: [
    'Name',
    'Geteilt',
    'Größe',
    'Bearbeitet',
    'Aktionen'
  ],
  Español: [
    'Nombre',
    'Shares',
    'Tamaño',
    'Modificado',
    'Acciones'
  ],
  Français: [
    'Nom',
    'Shares',
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
