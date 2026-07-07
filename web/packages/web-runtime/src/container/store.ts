import { NextApplication } from './application'

/**
 * main purpose of this is to keep a reference to all announced applications.
 * this is used for example in the mounted hook outside container module
 */
export const applicationStore = new Map<string, NextApplication>()
