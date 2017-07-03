package badm

import "fmt"

// ListBuckets lists the buckets in the database.
func ListBuckets(configurationPath string) error {
	c, err := readConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("read configuration: %v", err)
	}

	db, err := openDatabase(c)
	if err != nil {
		return fmt.Errorf("open database: %v", err)
	}
	defer db.close()

	if err := db.listBuckets(); err != nil {
		return fmt.Errorf("list buckets: %v", err)
	}

	return nil
}
