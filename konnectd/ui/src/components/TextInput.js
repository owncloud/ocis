import React from 'react';
import PropTypes from 'prop-types';

import { useIntl } from 'react-intl';

const TextInput = (props) => {
  const intl = useIntl();

  return <input className="oc-input" {...props} placeholder={props.placeholder ? intl.formatMessage(props.placeholder) : null} />;
};

TextInput.propTypes = {
  placeholder: PropTypes.object,
}

export default TextInput;
