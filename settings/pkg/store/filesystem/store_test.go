package store

import (
	"os"
	"path/filepath"
)

const (
	// account UUIDs
	accountUUID1 = "c4572da7-6142-4383-8fc6-efde3d463036"
	//accountUUID2 = "e11f9769-416a-427d-9441-41a0e51391d7"
	//accountUUID3 = "633ecd77-1980-412a-8721-bf598a330bb4"

	// extension names
	extension1 = "test-extension-1"
	extension2 = "test-extension-2"

	// bundle ids
	bundle1 = "2f06addf-4fd2-49d5-8f71-00fbd3a3ec47"
	bundle2 = "2d745744-749c-4286-8e92-74a24d8331c5"
	bundle3 = "d8fd27d1-c00b-4794-a658-416b756a72ff"

	// setting ids
	setting1 = "c7ebbc8b-d15a-4f2e-9d7d-d6a4cf858d1a"
	setting2 = "3fd9a3d9-20b7-40d4-9294-b22bb5868c10"
	setting3 = "24bb9535-3df4-42f1-a622-7c0562bec99f"

	// value ids
	value1 = "fd3b6221-dc13-4a22-824d-2480495f1cdb"
	value2 = "2a0bd9b0-ca1d-491a-8c56-d2ddfd68ded8"
	//value3 = "b42702d2-5e4d-4d73-b133-e1f9e285355e"

	dataRoot = "/tmp/herecomesthesun"
)

func burnRoot() {
	os.RemoveAll(filepath.Join(dataRoot, "values"))
	os.RemoveAll(filepath.Join(dataRoot, "bundles"))
}
