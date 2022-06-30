import React from 'react';

import {
  Fade,
  CircularProgress,
 } from '@material-ui/core';
 import { makeStyles } from '@material-ui/core/styles';

 const useStyles = makeStyles(() => ({
  root: {
    position: 'fixed',
    top: '50%',
    left: '50%',
    transform: 'translate(-50%, -50%)',
  },
 }));

const Spinner = () => {
  const classes = useStyles();

  return <div className={classes.root}>
    <Fade
      in
      style={{
        transitionDelay: '800ms',
      }}
      unmountOnExit
    >
      <CircularProgress size={70} thickness={1}/>
    </Fade>
  </div>;
}

export default Spinner;
