import yaml
import sys
import os
from datetime import date
from urllib.request import urlopen

## this python script generates based on defined variables adoc files for added, removed and deprecated
## envvars based on the env_vars.yaml that must exist in each referenced version.
## it is CRUCIAL that the version compared TO is actual - do required updates first!
## note that env_vars.yaml has been introduced with v6.0.0, comparing earlier is not possible
## note that we are always comparing from github sources and NOT local files

## when the files got created, you MUST do some post work manually like referencing services with xref:
## when running, files get recreated, existing content gets overwritten!!

## you MUST run this script from the local ocis repo root !!
## like: python3 docs/helpers/changed_envvars.py
## create a branch to prepare the changes

# CHANGE according your needs
# old is the base version to compare from
# new is the target version to compare to
# tagged versions must be of format: 'tags/v6.0.0'
# master is different, it must be:   'heads/master'
versionOld = 'tags/v6.0.0'
versionNew = 'heads/master'

# CHANGE according your needs
from_version = '5.0.0'
to_version = '7.0.0'

# CHANGE according your needs
# this will create files like 5.0.0-7.0.0-added and 5.0.0-7.0.0-removed
# this should match which versions you compare. master is ok if that is the base for a named release
nameComponent = '5.0.0-7.0.0'

# ADD new elements when a new version has been published so that it gets excluded
# array of version patterns to be excluded for added items. we dont need patch versions
excludePattern = ['pre5.0', '5.0', '6.0']

# DO NOT CHANGE
# this is the path the added/removed result is written to
adocWritePath = 'docs/services/general-info/env-var-deltas'

addedWith = {}
removedWith = {}
deprecatedWith = {}

def check_path():
	# check which path the script started. we can do this because the target path must be present
	# exit if not present
	if not os.path.exists(adocWritePath):
		print('Path not found: ' + adocWritePath)
		sys.exit()

def get_sources(versionOld, versionNew):
	# get the sources from github
	git_bleft_dir = 'https://raw.githubusercontent.com/owncloud/ocis/refs/'
	git_right_dir ='/docs/helpers/env_vars.yaml'

	urlOld = git_bleft_dir + versionOld + git_right_dir
	urlNew = git_bleft_dir + versionNew + git_right_dir

	try:
		fileOld = urlopen(urlOld).read().decode('utf-8')
		fileNew = urlopen(urlNew).read().decode('utf-8')
		return	yaml.safe_load(fileOld), yaml.safe_load(fileNew)

	except Exception as e:
		print(e)
		sys.exit()

def get_added(fileNew, excludePattern):
	# create dict with added items
	addedWith = {}
	for key, value in fileNew.items():
		if not fileNew[key]['introductionVersion'] in str(excludePattern):
			addedWith[key] = value
	return addedWith

def get_removed(fileOld, fileNew):
	# create dict with removed items
	removedWith = {}
	for key, value in fileOld.items():
		if not key in fileNew:
			removedWith[key] = value
	return removedWith

def get_deprecated(fileNew):
	# create dict with deprecated items
	deprecatedWith = {}
	for key, value in fileNew.items():
		if value['removalVersion']:
			deprecatedWith[key] = value
	return deprecatedWith

def create_adoc_start(type_text, from_version, to_version, creation_date, default):
	# create the page/table header
	a = '''// # {ftype} Variables between oCIS {ffrom} and oCIS {fto}
// commenting the headline to make it better includable

// table created per {fdate}
// the table should be recreated/updated on source () changes

[width="100%",cols="~,~,~,~",options="header"]
|===
| Service| Variable| Description| {fdefault}

'''.format(ftype = type_text, ffrom = from_version, fto = to_version, fdate = creation_date, fdefault = default)
	return a

def create_adoc_end():
	# close the table
	a = '''|===

'''
	return a

def add_adoc_line(service, variable, description, default):
	# add a table line
	a = '''| {fservice}
| {fvariable}
| {fdescription}
| {fdefault}

'''.format(fservice = service, fvariable = variable, fdescription = description, fdefault = default)
	return a

def create_table(type_text, source_dict, from_version, to_version, date_today, default = 'Default'):
	# get the table header
	a = create_adoc_start(type_text, from_version, to_version, date_today, default)
	cond_value = 'defaultValue' if default == 'Default' else 'removalVersion'
	# first add all ocis_
	for key, value in source_dict.items():
		if key.startswith('OCIS_'):
			a += add_adoc_line(
					'xref:deployment/services/env-vars-special-scope.adoc[Special Scope Envvars]',
					key,
					value['description'],
					value[cond_value]
				)
	# then add all others
	for key, value in source_dict.items():
		if not key.startswith('OCIS_'):
			a += add_adoc_line(
					'xref:{s-path}/xxx.adoc[xxx]',
					key,
					value['description'],
					value[cond_value]
				)
	# finally close the table
	a += create_adoc_end()
	return a

def write_output(a, type_text):
	# write the content to a file
	try:
		with open(adocWritePath + '/' + nameComponent + '-' + type_text + '.adoc', 'w') as file:
			file.write(a)
	except Exception as e:
		print('Failed creating ' + type_text + ' file')
		print(e)
		sys.exit()

## here are the tasks in sequence

check_path()
fileOld, fileNew = get_sources(versionOld, versionNew)
addedWith = get_added(fileNew, excludePattern)
removedWith = get_removed(fileOld, fileNew)
deprecatedWith = get_deprecated(fileNew)

a = create_table('Added', addedWith, from_version, to_version, date.today().strftime('%Y.%m.%d'))
r = create_table('Removed', removedWith, from_version, to_version, date.today().strftime('%Y.%m.%d'))
d = create_table('Deprecated', deprecatedWith, from_version, to_version, date.today().strftime('%Y.%m.%d'), 'Removalversion')

write_output(a, 'added')
write_output(r, 'removed')
write_output(d, 'deprecated')

print('Success, see files created in: ' + adocWritePath)
