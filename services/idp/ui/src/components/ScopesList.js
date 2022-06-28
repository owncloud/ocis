import React from 'react';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import { withStyles } from '@material-ui/core/styles';
import PropTypes from 'prop-types';
import Checkbox from '@material-ui/core/Checkbox';

import { injectIntl, useIntl, defineMessages, FormattedMessage } from 'react-intl';

const styles = () => ({
  row: {
    paddingTop: 0,
    paddingBottom: 0
  }
});

const scopeIDTranslations = defineMessages({
  'scope_alias_basic': {
    id: 'konnect.scopeDescription.aliasBasic',
    defaultMessage: 'Access your basic account information'
  },
  'scope_offline_access': {
    id: 'konnect.scopeDescription.offlineAccess',
    defaultMessage: 'Keep the allowed access persistently and forever'
  }
});

const ScopesList = ({scopes, meta, classes, ...rest}) => {
  const { mapping, definitions } = meta;
  const intl = useIntl()

  const rows = [];
  const known = {};

  // TODO(longsleep): Sort scopes according to priority.
  for (let scope in scopes) {
    if (!scopes[scope]) {
      continue;
    }
    let id = mapping[scope];
    if (id) {
      if (known[id]) {
        continue;
      }
      known[id] = true;
    } else {
      id = scope;
    }
    let definition = definitions[id];
    let label ;
    if (definition) {
      if (definition.id) {
        const translation = scopeIDTranslations[definition.id];
        if (translation) {
          label = intl.formatMessage(translation);
        }
      }
      if (!label) {
        label = definition.description;
      }
    }
    if (!label) {
      label = <FormattedMessage
        id="konnect.scopeDescription.scope"
        defaultMessage="Scope: {scope}"
        values={{scope}}
      />;
    }

    rows.push(
      <ListItem
        disableGutters
        dense
        key={id}
        className={classes.row}
      ><Checkbox
          checked
          disableRipple
          disabled
          className="oc-checkbox-dark"
        />
        <ListItemText primary={label} className="oc-light" />
      </ListItem>
    );
  }

  return (
    <List {...rest}>
      {rows}
    </List>
  );
};

ScopesList.propTypes = {
  classes: PropTypes.object.isRequired,

  scopes: PropTypes.object.isRequired,
  meta: PropTypes.object.isRequired
};

export default withStyles(styles)(injectIntl(ScopesList));
