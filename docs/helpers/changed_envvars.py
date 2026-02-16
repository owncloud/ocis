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

## !! you MUST run this script from the local ocis repo root !!
## like: python3 docs/helpers/changed_envvars.py
## create a branch to prepare the changes

# CHANGE according your needs
# old is the base version to compare from which is always a tagged release
# new is the target version to compare to and can be a branch or tagged
# tagged versions must be of format:            'tags/v7.1.0'
# branches are different, master or stable-x.y: 'heads/master'
# stick with master because the follow up process is much easier and you avoid confusion
versionOld: str = 'tags/v7.3.0'
versionNew: str = 'heads/master'

# CHANGE according your needs
from_version: str = '7.3.0'
to_version: str = '8.0.0'

# CHANGE according your needs
# this will create files like 5.0.0-7.0.0-added and 5.0.0-7.0.0-removed
# this should match which versions you compare. master is ok if that is the base for a named release
nameComponent: str = '7.3.0-8.0.0'

# ADD new elements when a new version has been published so that it gets excluded
# array of version patterns to be excluded for added items. we dont need patch versions
excludePattern: list[str] = ['pre5.0', '5.0', '6.0', '6.0.0', '6.0.1', '6.1.0', '6.7', '7.0', '7.0.0', '7.1.0', '7.2.0', '7.3.0']

# DO NOT CHANGE
# this is the path the added/removed result is written to
adocWritePath: str = 'docs/services/general-info/envvars/env-var-deltas'

# define root sources from github
git_left_dir: str = 'https://raw.githubusercontent.com/owncloud/ocis/refs/'
git_right_dir: str ='/docs/helpers/env_vars.yaml'

def check_path() -> None:
	# check which path the script started. we can do this because the target path must be present
	# exit if not present
	if not os.path.exists(adocWritePath):
		print('Path not found: ' + adocWritePath)
		sys.exit()

def get_sources(versionOld: str, versionNew: str) -> (str, str):
	urlOld = git_left_dir + versionOld + git_right_dir
	urlNew = git_left_dir + versionNew + git_right_dir

	try:
		fileOld = urlopen(urlOld).read().decode('utf-8')
		fileNew = urlopen(urlNew).read().decode('utf-8')
		return	yaml.safe_load(fileOld), yaml.safe_load(fileNew)

	except Exception as e:
		print(e)
		sys.exit()

def get_added(fileNew: str, excludePattern: list[str]) -> set[str]:
	# create dict with added items
	addedWith: set[str] = {}
	for key, value in fileNew.items():
		if not fileNew[key]['introductionVersion'] in str(excludePattern):
			addedWith[key] = value
	return addedWith

def get_removed(fileOld: str, fileNew: str) -> set[str]:
	# create dict with removed items
	removedWith: set[str] = {}
	for key, value in fileOld.items():
		if not key in fileNew:
			removedWith[key] = value
	return removedWith

def get_deprecated(fileNew: str) -> set[str]:
	# create dict with deprecated items
	deprecatedWith: set[str] = {}
	for key, value in fileNew.items():
		if value['removalVersion']:
			deprecatedWith[key] = value
	return deprecatedWith

def create_adoc_start(type_text: str, from_version: str, to_version: str, creation_date: str, columns: str, closing: str) -> str:
	# create the page/table header
	# 'closing' contains variable column names dependen if added/removed ir deprecated
	a: str = '''// # {ftype} Variables between oCIS {ffrom} and oCIS {fto}
// commenting the headline to make it better includable

// table created per {fdate}
// the table should be recreated/updated on source () changes

[width="100%",cols="{fcolumns}",options="header"]
|===
| Service | Variable | Description | {fclosing}

'''.format(ftype = type_text, ffrom = from_version, fto = to_version, fdate = creation_date, fcolumns = columns, fclosing = closing)
	return a

def create_adoc_end() -> str:
	# close the table
	a: str = '''|===

'''
	return a

def add_adoc_line_1(service: str, variable: str, description: str, value: str) -> str:
	# add a table line for added/removed
	# the dummy values are only here to have the same number of parameters as add_adoc_line_2
	a: str = '''| {fservice}
| {fvariable}
| {fdescription}
| {fvalue}

'''.format(fservice = service, fvariable = variable, fdescription = description, fvalue = value)
	return a

def add_adoc_line_2(service, variable, description, removalVersion, deprecationInfo):
	# add a table line for deprecated, this has different columns
	a: str = '''| {fservice}
| {fvariable}
| {fdescription}
| {fremovalVersion}
| {fdeprecationInfo}

'''.format(fservice = service, fvariable = variable, fdescription = description, fremovalVersion = removalVersion, fdeprecationInfo = deprecationInfo)
	return a

def create_table(type_text: str, source_dict: set[str], from_version: str, to_version: str, date_today: str, type: bool = False) -> str:
	# get the table header
	columns: str = '~,~,~,~' if type == False else '~,~,~,~,~'
	closing: str = 'Default' if type == False else 'Removal Version | Deprecation Info'
	a: str = create_adoc_start(type_text, from_version, to_version, date_today, columns, closing)

	if not type:
	# added and removed envvars
		# first add all ocis_
		for key, value in source_dict.items():
			if key.startswith('OCIS_'):
				a += add_adoc_line_1(
						# note that any envvar starting with ocis cant be assigned to a service automatically
						# the xref must be corrected in the output file manually
						'xref:deployment/services/env-vars-special-scope.adoc[Special Scope Envvars]',
						key,
						value['description'],
						value['defaultValue']
					)
		# then add all others
		for key, value in source_dict.items():
			if not key.startswith('OCIS_'):
				a += add_adoc_line_1(
						'xref:{s-path}/xxx.adoc[xxx]',
						key,
						value['description'],
						value['defaultValue']
					)
	else:
	# deprecated envvars
		# first add all ocis_
		for key, value in source_dict.items():
			if key.startswith('OCIS_'):
				a += add_adoc_line_2(
						'xref:deployment/services/env-vars-special-scope.adoc[Special Scope Envvars]',
						key,
						value['description'],
						value['removalVersion'],
						value['deprecationInfo']
					)
		# then add all others
		for key, value in source_dict.items():
			if not key.startswith('OCIS_'):
				a += add_adoc_line_2(
						'xref:{s-path}/xxx.adoc[xxx]',
						key,
						value['description'],
						value['removalVersion'],
						value['deprecationInfo']
					)

	# finally close the table
	a += create_adoc_end()
	return a

def write_output(a: str, type_text: str) -> None:
	# write the content to a file
	try:
		with open(adocWritePath + '/' + nameComponent + '-' + type_text + '.adoc', 'w') as file:
			file.write(a)
	except Exception as e:
		print('Failed creating ' + type_text + ' file')
		print(e)
		sys.exit()

def main() -> None:
	## here are the tasks in sequence

	fileOld: str = ''
	fileNew: str = ''
	addedWith: set[str] = {}
	removedWith: set[str] = {}
	deprecatedWith: set[str] = {}

	check_path()
	fileOld, fileNew = get_sources(versionOld, versionNew)
	addedWith = get_added(fileNew, excludePattern)
	removedWith = get_removed(fileOld, fileNew)
	deprecatedWith = get_deprecated(fileNew)

	a: str = create_table('Added', addedWith, from_version, to_version, date.today().strftime('%Y.%m.%d'))
	r: str = create_table('Removed', removedWith, from_version, to_version, date.today().strftime('%Y.%m.%d'))
	d: str = create_table('Deprecated', deprecatedWith, from_version, to_version, date.today().strftime('%Y.%m.%d'), True)

	write_output(a, 'added')
	write_output(r, 'removed')
	write_output(d, 'deprecated')

	print('Success, see files created in: ' + adocWritePath)

if __name__ == '__main__':
	main()
