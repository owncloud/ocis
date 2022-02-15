const filesMenu = {
  English: [
    'All files',
    'Shared with me',
    'Shared with others',
    'Shared via link',
    'Spaces',
    'Deleted files'
  ],
  Deutsch: [
    'Alle Dateien',
    'Mit mir geteilt',
    'Mit anderen geteilt',
    'Per Link geteilt',
    'Spaces',
    'Gelöschte Dateien'
  ],
  Español: [
    'Todos los archivos',
    'Compartido conmigo',
    'Compartido con otros',
    'Shared via link',
    'Spaces',
    'Archivos borrados'
  ],
  Français: [
    'Tous les fichiers',
    'Partagé avec moi',
    'Partagé avec autres',
    'Shared via link',
    'Spaces',
    'Fichiers supprimés'
  ]
}

const accountMenu = {
  English: [
    'N\nnull\nuser1@example.com',
    'Settings',
    'Log out',
    'Personal storage (0.2% used)\n5.06 GB of 2.85 TB used'
  ],
  Deutsch: [
    'N\nnull\nuser1@example.com',
    'Einstellungen',
    'Abmelden',
    'Persönlicher Speicher (0.2% benutzt)\n5.06 GB von 2.85 TB benutzt'
  ],
  Español: [
    'N\nnull\nuser1@example.com',
    'Configuración',
    'Salir',
    'Personal storage (0.2% used)\n5.06 GB of 2.85 TB used'
  ],
  Français: [
    'N\nnull\nuser1@example.com',
    'Settings',
    'Se déconnecter',
    'Personal storage (0.2% used)\n5.06 GB of 2.85 TB used'
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
    'Bearbeitet',
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
