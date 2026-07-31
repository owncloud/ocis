# web-app-password-protected-folders

This web extension enhances the oCIS platform by allowing users to create password-protected folders (PPF). It provides an additional layer of security for sensitive or confidential information stored within oCIS.

Note that this extension is included and enabled by default in oCIS and not explicitly shown as an extension in the webUI.

## Features

- 🔒 Create password-protected folders within oCIS
- 🔑 Set unique passwords for each protected folder
- 🎨 Seamless integration with the oCIS user interface

## Usage

Users can create password-protected folders by following these steps:

1. Log in to your oCIS instance.
1. Select a space and browse to a location where you want to create a password protected folder.
1. Click on the **New** button and select **Password Protected Folder**.
1. Enter a name for the folder and set a unique password.
1. Click on **Create** to create the password protected folder.
1. Users will be prompted to enter the set password to access the protected folder.

## How it Works

1. When a PPF is created, it actually creates a file with a special hidden extension (.psec) that only contains a URL with a password protected share pointing to the target folder.
1. Opening the link file via the oCIS web interface automatically initiates access to the PPF. The rules for password protected shares apply.
1. The target folder for the PPF is ALWAYS created as a **subfolder** mirroring the original path of the link folder in a hidden folder in the root of the PPF creator's private space. This means that the target folder will not normally be visible to the PPF creator, unless the display of hidden files is enabled in their private space. Do not change or delete this folder unless you know what you are doing.
1. The link file and the destination folder have the same name.
1. When the **owner** of the PPF deletes the **link file**, the target folder and its contents are automatically deleted.
1. If the *owner** of the PPF deletes the parent **folder(s)** containing the **link file**, the link file and associated data are automatically deleted.
1. Deleting a PPF will move the link file and destination folder to the appropriate recycle bin of the space.

## Considerations

1. **PPF Location**

   When creating a PPF in a Space, the created link file inherits the permissions defined for the user who is granted access to the space. This has some managing and security implications to consider:

   1. Although possible, do not create a PPF in your private Space for security reasons. You would need to share the Space or link file first. Making your entire private Space public is a security issue for your private data.
   1. If you create a PPF in a Project Space where users also have delete permissions, **any** user accessing the Space with delete permissions would be able to delete the link file. Accidentally or deliberately done, no one would be able to access the PPF, even though the data of the target folder still exists.
   1. As a suggestion for a secure environment for accessing PPF's, it is recommended to create a separate project space for PPF's only, where added users will only have view permissions. However, Space Managers or defined users can or will have extended permissions. This ensures that users who are allowed to access the Space cannot "accidentally" delete or manipulate link files.

1. **Restoring a Deleted a PPF**

   It is most important to understand that resources are not accessed by their name but their internal UNIQUE ID (UUID).

   1. If a PPF is deleted by the creator, both the link file and the destination folder will be moved to the trash bin of their respective Spaces. Although you can restore both, the PPF will no longer work because the underlying UUID has changed.
   1. To recreate a deleted PPF (with the same name), you must follow the steps below exactly:
      1. If not already done, the PPF creator needs to enable the display of hidden files in their Private Space.
      1. Undelete the target folder from the trash bin.
      1. Rename the restored destination folder, the name MUST NOT be the same as before.
      1. DO NOT undelete the link file from the project space! This file can be permanently deleted.
      1. Create a new PPF. You can now use the previous PPF name as the original destination folder that was restored now has a different name.
      1. Copy or move any data from the restored and renamed destination folder to the newly created destination folder.
      1. Delete the renamed destination folder, it is no longer needed.
      1. If desired, prevent hidden files from appearing in the personal space again and prevent them from being viewed.

1. **Moving or Renaming a Link File**

   Moving or renaming a link file is prevented by the webUI, but you could copy or manipulate this file via the webUI if permissions allow, or via one of the Desktop, iOS or Android apps if `.psec' files are not excluded from the sync list. There are a few things to keep in mind:
   
   1. Link files and the destination folder share the same path. If the owner deletes a link file that has been moved, renamed or copied, the link file will be removed, but not the destination folder, because it no longer has a valid mirrored path to the destination folder and cannot find it. This operation will fail silently and nothing will be reported..
   1. Deleting a manipulated link file requires manual deletion of the target folder. Failure to do this may result in orphaned data.!!

1. **Deleting a PPF with Changed Access Rights**

   Once a PPF has been created, the owner may subsequently find that access rights to the source have been reduced due to a change in access policy. This change can occur at any time after the PPF has been created. Although the owner may (or may not) have read access, they will no longer be able to delete the PPF. To delete a PPF in such a situation, the following steps must be taken:
   
   1. Either one of the Space Managers or a user with delete privileges must delete the link file. This will make the PPF inaccessible. Note that this step does not delete the content source in the user's Personal Space.
   1. The user who created the PPF must manually delete the content source from their Personal Space.
