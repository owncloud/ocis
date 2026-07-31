import { computed } from 'vue'
import { omit } from 'lodash-es'
import { ContextualHelperData } from '@ownclouders/design-system/helpers'
import { useConfigStore } from '@ownclouders/web-pkg'
import { useGettext } from 'vue3-gettext'

export const useContextualHelpers = () => {
  const { $gettext } = useGettext()
  const configStore = useConfigStore()

  const filterContextHelper = (data: ContextualHelperData): ContextualHelperData => {
    if (configStore.options.contextHelpersReadMore === false) {
      return omit(data, 'readMoreLink')
    }
    return data
  }

  const shareInviteCollaboratorHelp = computed<ContextualHelperData>(() =>
    filterContextHelper({
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
        { text: $gettext('"via folder"'), headline: true },
        {
          text: $gettext(
            'The "via folder" information can be shown in the details view, if access has already been given via a parent folder. Click on the "via folder" to edit the share on its parent folder.'
          )
        }
      ],
      readMoreLink: 'https://doc.owncloud.com/go?to=webui-users-sharing'
    })
  )

  const shareInviteCollaboratorHelpCern = computed<ContextualHelperData>(() =>
    filterContextHelper({
      title: '',
      list: [
        { text: $gettext('Search for service or secondary Account'), headline: true },
        {
          text: $gettext(
            'To search for service or secondary accounts prefix the username with "a:" (like "a:doe") and for guest accounts prefix the username with "l:" (like "l:doe").'
          )
        }
      ]
    })
  )

  const shareSpaceAddMemberHelp = computed<ContextualHelperData>(() =>
    filterContextHelper({
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
    })
  )

  const shareViaLinkHelp = computed<ContextualHelperData>(() =>
    filterContextHelper({
      title: $gettext('Choose how access is granted'),
      list: [
        {
          text: $gettext(
            'No login required. Everyone with the link can access. If you share this link with people from the list "Invited people", they need to login so that their individual assigned permissions can take effect. If they are not logged-in, the permissions of the link take effect.'
          )
        }
      ],
      readMoreLink: 'https://doc.owncloud.com/go?to=webui-users-sharing'
    })
  )

  const shareViaIndirectLinkHelp = computed<ContextualHelperData>(() =>
    filterContextHelper({
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
    })
  )

  const tagsHelper = computed<ContextualHelperData>(() =>
    filterContextHelper({
      title: $gettext('Who can view tags?'),
      list: [
        {
          text: $gettext(
            'Everyone who can view the file can view its tags. Likewise, everyone who can edit the file can edit its tags.'
          )
        }
      ]
    })
  )

  return {
    shareInviteCollaboratorHelp,
    shareInviteCollaboratorHelpCern,
    shareSpaceAddMemberHelp,
    shareViaLinkHelp,
    shareViaIndirectLinkHelp,
    tagsHelper
  }
}
