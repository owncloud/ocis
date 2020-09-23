const filesMenu = {
  English: [
    'All files',
    'Shared with me',
    'Shared with others',
    'Trash bin'
  ],
  Deutsch: [
    'Alle Dateien',
    'Mit mir geteilt',
    'Mit anderen geteilt',
    'Papierkorb'
  ],
  Español: [
    'Todos los archivos',
    'Compartido conmigo',
    'Compartido con otros',
    'Papelera de reciclaje'
  ],
  Français: [
    'Tous les fichiers',
    'Partagé avec moi',
    'Partagé avec autres',
    'Corbeille'
  ]
}

const accountMenu = {
  English: [
    'Manage your account',
    'Log out'
  ],
  Deutsch: [
    'Verwalten Sie Ihr Benutzerkonto',
    'Abmelden'
  ],
  Español: [
    'Administra tu cuenta',
    'Salir'
  ],
  Français: [
    'Modifier votre compte',
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
