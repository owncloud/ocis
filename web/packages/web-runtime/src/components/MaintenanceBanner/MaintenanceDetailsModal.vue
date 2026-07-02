<template>
  <div>
    <p>
      {{
        $gettext(
          'The system is in maintenance mode. Any access will fail during this mode such as:'
        )
      }}
    </p>

    <ul>
      <li>{{ $gettext('Accessing or changing a file or settings') }}</li>
      <li>{{ $gettext('Saving a file you are currently editing (*)') }}</li>
      <li>{{ $gettext('Changing to another Space or directory') }}</li>
    </ul>

    <p>
      {{
        $gettext(
          '* ... Saving an edited file may not respond immediately to maintenance mode, depending on its duration. If you receive a warning that you cannot save, save a backup of your file locally.'
        )
      }}
    </p>

    <p>
      {{
        $gettext(
          'When maintenance mode is complete, the yellow warning bar will automatically disappear and full access will be restored. You can now save an open edited file as normal.'
        )
      }}
    </p>
  </div>
</template>
