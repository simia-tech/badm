package badm

// SetKeyType sets the key type for the provided bucket.
func SetKeyType(configurationPath, name, value string) error {
	return updateConfiguration(configurationPath, func(c *Configuration) error {
		if c.Buckets == nil {
			c.Buckets = make(map[string]*Bucket)
		}

		b, found := c.Buckets[name]
		if found {
			b.KeyType = value
		} else {
			c.Buckets[name] = &Bucket{KeyType: value}
		}
		return nil
	})
}

// SetValueType sets the value type for the provided bucket.
func SetValueType(configurationPath, name, value string) error {
	return updateConfiguration(configurationPath, func(c *Configuration) error {
		if c.Buckets == nil {
			c.Buckets = make(map[string]*Bucket)
		}

		b, found := c.Buckets[name]
		if found {
			b.ValueType = value
		} else {
			c.Buckets[name] = &Bucket{ValueType: value}
		}
		return nil
	})
}

// Clear clears the bucket configuration.
func Clear(configurationPath, name string) error {
	return updateConfiguration(configurationPath, func(c *Configuration) error {
		if c.Buckets == nil {
			return nil
		}
		delete(c.Buckets, name)
		return nil
	})
}
