import { createMuiTheme } from '@material-ui/core/styles'

import blueGrey from '@material-ui/core/colors/blueGrey'
import red from '@material-ui/core/colors/red'

const kopanoBlue = {
  50: '#e5f7fd',
  100: '#ccf0fc',
  200: '#99e1f9',
  300: '#66d1f6',
  400: '#33c2f3',
  500: '#00b3f0',
  600: '#14a9e7',
  700: '#1e98d7',
  800: '#188cc8',
  900: '#0d74b6',
  A100: '#96e8ff',
  A200: '#6ddaff',
  A400: '#56cbff',
  A700: '#48acf0'
}

const primaryColor = kopanoBlue
const secondaryColor = blueGrey
const errorColor = red

// All the following keys are optional.
// We try our best to provide a great default value.
const theme = createMuiTheme({
  palette: {
    primary: primaryColor,
    secondary: secondaryColor,
    error: errorColor,
    // Used by `getContrastText()` to maximize the contrast between the background and
    // the text.
    // NOTE(longsleep): KopanoBlue is too light and thus needs 2.4 contrastThreshold
    // to make sure the default 500 color is still using white text. It will
    // show warnings in development mode that the contrast is too low as W3C
    // recommends the threshold to be 3 or more. This cannot be helped.
    contrastThreshold: 2.4,
    // Used to shift a color's luminance by approximately
    // two indexes within its tonal palette.
    // E.g., shift from Red 500 to Red 300 or Red 700.
    tonalOffset: 0.2
  },
  typography: {
    useNextVariants: true,
    button: {
      textTransform: 'none',
      fontSize: '1rem'
    }
  }
})

export default theme
