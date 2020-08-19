package store

import (
	"os"
	"path/filepath"
)

const folderNameBundles = "bundles"
const folderNameValues = "values"
const folderNameAssignments = "assignments"

// buildFolderPathForBundles builds the folder path for storing settings bundles. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFolderPathForBundles(mkdir bool) string {
	folderPath := filepath.Join(s.dataPath, folderNameBundles)
	if mkdir {
		s.ensureFolderExists(folderPath)
	}
	return folderPath
}

// buildFilePathForBundle builds a unique file name from the given params. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFilePathForBundle(bundleID string, mkdir bool) string {
	extensionFolder := s.buildFolderPathForBundles(mkdir)
	return filepath.Join(extensionFolder, bundleID+".json")
}

// buildFolderPathForValues builds the folder path for storing settings values. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFolderPathForValues(mkdir bool) string {
	folderPath := filepath.Join(s.dataPath, folderNameValues)
	if mkdir {
		s.ensureFolderExists(folderPath)
	}
	return folderPath
}

// buildFilePathForValue builds a unique file name from the given params. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFilePathForValue(valueID string, mkdir bool) string {
	extensionFolder := s.buildFolderPathForValues(mkdir)
	return filepath.Join(extensionFolder, valueID+".json")
}

// buildFolderPathForRoleAssignments builds the folder path for storing role assignments. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFolderPathForRoleAssignments(mkdir bool) string {
	roleAssignmentsFolder := filepath.Join(s.dataPath, folderNameAssignments)
	if mkdir {
		s.ensureFolderExists(roleAssignmentsFolder)
	}
	return roleAssignmentsFolder
}

// buildFilePathForRoleAssignment builds a unique file name from the given params. If mkdir is true, folders in the path will be created if necessary.
func (s Store) buildFilePathForRoleAssignment(assignmentID string, mkdir bool) string {
	roleAssignmentsFolder := s.buildFolderPathForRoleAssignments(mkdir)
	return filepath.Join(roleAssignmentsFolder, assignmentID+".json")
}

// ensureFolderExists checks if the given path is an existing folder and creates one if not existing
func (s Store) ensureFolderExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0700)
		if err != nil {
			s.Logger.Err(err).Msgf("Error creating folder %v", path)
		}
	}
}
