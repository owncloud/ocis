import React, { useCallback, useMemo, useEffect } from 'react';
import PropTypes from 'prop-types';

import { useTranslation } from 'react-i18next';

import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';

import allLocales from '../locales';

function LocaleSelect({ locales: localesProp, ...other } = {}) {
  const { i18n, ready } = useTranslation();

  const handleChange = useCallback((event) => {
    i18n.changeLanguage(event.target.value);
  }, [ i18n ])

  const locales = useMemo(() => {
    if (!localesProp) {
      return allLocales;
    }
    const supported = allLocales.filter(locale => {
      return localesProp.includes(locale.locale);
    });
    return supported;
  }, [localesProp]);

  useEffect(() => {
    if (locales) {
      const found = locales.find((locale) => {
        return locale.locale === i18n.language;
      });
      if (found) {
        // Have language -> is supported all good.
        return;
      }
      const wanted = i18n.modules.languageDetector.detectors.navigator.lookup();
      i18n.modules.languageDetector.services.languageUtils.options.supportedLngs = locales.map(locale => locale.locale);
      i18n.modules.languageDetector.services.languageUtils.options.fallbackLng = null;

      let best = i18n.modules.languageDetector.services.languageUtils.getBestMatchFromCodes(wanted);
      if (!best) {
        best = locales[0].locale;
      }

      // Auto change language to best one found if the current selected one is not enabled.
      if (i18n.language !== best) {
        i18n.changeLanguage(best);
      }
    }
  }, [i18n, locales]);

  if (!ready || !locales || locales.length < 2) {
    return null;
  }

  return <Select
    value={i18n.language}
    onChange={handleChange}
    {...other}
  >
    {locales.map(language => {
      return <MenuItem
        key={language.locale}
        value={language.locale}>
        {language.nativeName}
      </MenuItem>;
    })}
  </Select>;
}

LocaleSelect.propTypes = {
  locales: PropTypes.arrayOf(PropTypes.string),
};

export default LocaleSelect;
