const filesMenu = {
  English: [
    'All files',
    'Shared with me',
    'Shared with others',
    'Deleted files'
  ],
  Deutsch: [
    'Alle Dateien',
    'Mit mir geteilt',
    'Mit anderen geteilt',
    'Gelöschte Dateien'
  ],
  Español: [
    'Todos los archivos',
    'Compartido conmigo',
    'Compartido con otros',
    'Archivos borrados'
  ],
  Français: [
    'Tous les fichiers',
    'Partagé avec moi',
    'Partagé avec autres',
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
    'Profile',
    'Settings',
    'Abmelden'
  ],
  Español: [
    'Perfil',
    'Ajustes',
    'Salir'
  ],
  Français: [
    'Profil',
    'Paramètres',
    'Se déconnecter'
  ]
}

const filesListHeaderMenu = {
  English: [
    'Name',
    'Size',
    'Updated',
    'Actions'
  ],
  Deutsch: [
    'Name',
    'Größe',
    'Erneuert',
    'Aktionen'
  ],
  Español: [
    'Nombre',
    'Tamaño',
    'Actualizado',
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
