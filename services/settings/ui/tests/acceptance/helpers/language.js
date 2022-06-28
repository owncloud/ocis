const filesMenu = {
  English: [
    'Personal',
    'Shares',
    'Spaces\nbeta',
    'Deleted files'
  ],
  Deutsch: [
    'Persönlich',
    'Geteilt',
    'Spaces\nbeta',
    'Gelöschte Dateien'
  ],
  Español: [
    'Personal',
    'Shares',
    'Spaces\nbeta',
    'Archivos borrados'
  ],
  Français: [
    'Personal',
    'Shares',
    'Spaces\nbeta',
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
