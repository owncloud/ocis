import React from 'react';
import PropTypes from 'prop-types';

import {useIntl} from 'react-intl';

const TextInput = (props) => {
    const intl = useIntl();
    const label = props.label;
    const extraClassName = props.extraClassName;

    delete props.label;
    delete props.extraClassName;

    return (
        <div>
            <label className="oc-label"
                   htmlFor={props.id}>{intl.formatMessage(label)}</label>
            <input className={`oc-input ${extraClassName ? extraClassName : ''}`} {...props}
                   placeholder={props.placeholder ? intl.formatMessage(props.placeholder) : null}/>
        </div>);
};

TextInput.propTypes = {
    placeholder: PropTypes.object,
    label: PropTypes.object,
    id: PropTypes.string,
    extraClassName: props.string,
}

export default TextInput;
