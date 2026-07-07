import { omit } from 'lodash-es'
import { ContextualHelperData } from '@ownclouders/design-system/helpers'
import { ConfigStore } from '@ownclouders/web-pkg'

// just a dummy function to trick gettext tools
function $gettext(msg: string) {
  return msg
}

interface ContextualHelperOptions {
  configStore: ConfigStore
}

export const shareInviteCollaboratorHelp = (options: ContextualHelperOptions) =>
  filterContextHelper(
    {
      title: $gettext('Share with people'),
      text: $gettext(
        'Use the input field to search for users and groups. Select them to share the item.'
      ),
      list: [
        { text: $gettext('Subfolders'), headline: true },
        {
          text: $gettext(
            'If you share a folder, all of its contents and subfolders will be shared as well.'
          )
        },
        { text: $gettext('Notification'), headline: true },
        {
          text: $gettext(
            'People you share resources with will be notified via email or in-app notification.'
          )
        },
        { text: $gettext('Incognito'), headline: true },
        {
          text: $gettext('People you share resources with can not see who else has access.')
        },
        { text: $gettext('“via folder”'), headline: true },
        {
          text: $gettext(
            'The “via folder” information can be shown in the details view, if access has already been given via a parent folder. Click on the “via folder” to edit the share on its parent folder.'
          )
        }
      ],
      readMoreLink: 'https://doc.owncloud.com/go?to=webui-users-sharing'
    },
    options
  )

export const shareInviteCollaboratorHelpCern = (options: ContextualHelperOptions) =>
  filterContextHelper(
    {
      title: '',
      list: [
        { text: $gettext('Search for service or secondary Account'), headline: true },
        {
          text: $gettext(
            'To search for service or secondary accounts prefix the username with "a:" (like "a:doe") and for guest accounts prefix the username with "l:" (like "l:doe").'
          )
        }
      ]
    },
    options
  )

export const shareSpaceAddMemberHelp = (options: ContextualHelperOptions) =>
  filterContextHelper(
    {
      title: $gettext('Add members to this Space'),
      text: $gettext('Enter a name to add people or groups as members to this Space.'),
      list: [
        { text: $gettext('Member capabilities'), headline: true },
        {
          text: $gettext(
            'Members are able to see who has access to this space and access all files in this space. Read or write permissions can be set by assigning a role.'
          )
        },
        { text: $gettext('Space manager capabilities'), headline: true },
        {
          text: $gettext(
            'Members with the Manager role are able to edit all properties and content of a Space, such as adding or removing members, sharing subfolders with non-members, or creating links to share.'
          )
        }
      ],
      readMoreLink: 'https://doc.owncloud.com/go?to=webui-users-sharing'
    },
    options
  )
export const shareViaLinkHelp = (options: ContextualHelperOptions) =>
  filterContextHelper(
    {
      title: $gettext('Choose how access is granted'),
      list: [
        {
          text: $gettext(
            'No login required. Everyone with the link can access. If you share this link with people from the list "Invited people", they need to login so that their individual assigned permissions can take effect. If they are not logged-in, the permissions of the link take effect.'
          )
        }
      ],
      readMoreLink: 'https://doc.owncloud.com/go?to=webui-users-sharing'
    },
    options
  )
export const shareViaIndirectLinkHelp = (options: ContextualHelperOptions) =>
  filterContextHelper(
    {
      title: $gettext('What are indirect links?'),
      text: $gettext('Indirect links are links giving access by a parent folder.'),
      list: [
        {
          text: $gettext('How to edit indirect links'),
          headline: true
        },
        {
          text: $gettext(
            'Indirect links can only be edited in their parent folder. Click on the folder icon below the link to navigate to the parent folder.'
          )
        }
      ],
      readMoreLink: 'https://doc.owncloud.com/go?to=webui-users-sharing'
    },
    options
  )

const filterContextHelper = (
  data: ContextualHelperData,
  options?: ContextualHelperOptions
): ContextualHelperData => {
  if (options.configStore.options.contextHelpersReadMore === false) {
    return omit(data, 'readMoreLink')
  }
  return data
}

export const tagsHelper = (options: ContextualHelperOptions) =>
  filterContextHelper(
    {
      title: $gettext('Who can view tags?'),
      list: [
        {
          text: $gettext(
            'Everyone who can view the file can view its tags. Likewise, everyone who can edit the file can edit its tags.'
          )
        }
      ]
    },
    options
  )
