package badm

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// ListKeys lists the keys in the provided bucket.
func ListKeys(configurationPath, name string) error {
	c, err := readConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("read configuration: %v", err)
	}

	db, err := openDatabase(c)
	if err != nil {
		return fmt.Errorf("open database: %v", err)
	}
	defer db.close()

	if err := db.listKeys(name); err != nil {
		return fmt.Errorf("list keys: %v", err)
	}

	return nil
}

// ListValues lists the values in the provided bucket.
func ListValues(configurationPath, name string) error {
	c, err := readConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("read configuration: %v", err)
	}

	db, err := openDatabase(c)
	if err != nil {
		return fmt.Errorf("open database: %v", err)
	}
	defer db.close()

	if err := db.listValues(name); err != nil {
		return fmt.Errorf("list values: %v", err)
	}

	return nil
}

// ListKeyValues lists the keys and values in the provided bucket.
func ListKeyValues(configurationPath, name string) error {
	c, err := readConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("read configuration: %v", err)
	}

	db, err := openDatabase(c)
	if err != nil {
		return fmt.Errorf("open database: %v", err)
	}
	defer db.close()

	if err := db.listKeyValues(name); err != nil {
		return fmt.Errorf("list key values: %v", err)
	}

	return nil
}

type database struct {
	configuration *Configuration
	db            *bolt.DB
}

func openDatabase(c *Configuration) (*database, error) {
	if c.SelectedPath == "" {
		return nil, fmt.Errorf("no database file selected")
	}

	db, err := bolt.Open(c.SelectedPath, 0644, nil)
	if err != nil {
		return nil, fmt.Errorf("open bolt: %v", err)
	}

	return &database{
		configuration: c,
		db:            db,
	}, nil
}

func (db *database) listBuckets() error {
	return db.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			fmt.Printf("%s (key-type: %s / value-type: %s)\n", name, db.configuration.KeyType(string(name)), db.configuration.ValueType(string(name)))
			return nil
		})
	})
}

func (db *database) listKeys(name string) error {
	return db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("could not find bucket %s", name)
		}

		t, err := TypeFor(db.configuration.KeyType(name))
		if err != nil {
			return err
		}

		return bucket.ForEach(func(key, _ []byte) error {
			v, err := t(key)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", v)
			return nil
		})
	})
}

func (db *database) listValues(name string) error {
	return db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("could not find bucket %s", name)
		}

		t, err := TypeFor(db.configuration.ValueType(name))
		if err != nil {
			return err
		}

		return bucket.ForEach(func(_, value []byte) error {
			v, err := t(value)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n", v)
			return nil
		})
	})
}

func (db *database) listKeyValues(name string) error {
	return db.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(name))
		if bucket == nil {
			return fmt.Errorf("could not find bucket %s", name)
		}

		kt, err := TypeFor(db.configuration.KeyType(name))
		if err != nil {
			return err
		}
		vt, err := TypeFor(db.configuration.ValueType(name))
		if err != nil {
			return err
		}

		return bucket.ForEach(func(key, value []byte) error {
			kv, err := kt(key)
			if err != nil {
				return err
			}
			vv, err := vt(value)
			if err != nil {
				return err
			}
			fmt.Printf("%s\n%s\n\n", kv, vv)
			return nil
		})
	})
}

func (db *database) close() error {
	return db.db.Close()
}
