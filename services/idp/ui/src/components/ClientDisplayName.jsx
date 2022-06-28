import React from 'react';
import PropTypes from 'prop-types';

const ClientDisplayName = ({ client, ...rest }) => (
  <span {...rest}>{client.display_name ? client.display_name : client.id}</span>
);

ClientDisplayName.propTypes = {
  client: PropTypes.object.isRequired
};

export default ClientDisplayName;
