package main

import (
	"fmt"
	"os"
	"path"
)

func createFileSystemStructure(rootLocation string, bucketLocations []string) error {
	_, err := os.Stat(rootLocation)
	if err != nil {
		p := path.Join(rootLocation, "main")
		err := os.MkdirAll(p, 0755)
		if err != nil {
			return fmt.Errorf("unable to create the root directory %s", err)
		}
	}

	// pbl := parseBucketLocations(bucketLocations)

	/* A GCS bucket name has the structure gs://bucket-name, however while reading the SDK natively handles it
	and expects just 'bucket-name' to be passed in.
	*/
	for _, bl := range bucketLocations {
		bucketDirLoc := path.Join(rootLocation, bl)
		_, err := os.Stat(bucketDirLoc)
		if err != nil {
			err := os.Mkdir(bucketDirLoc, 0755)
			if err != nil {
				return fmt.Errorf("unable to create a directory for the bucket %s. The following error occured %s", bl, err)
			}
		}
	}

	return nil
}

// func parseBucketLocations(bucketLocations []string) []string {
// 	parsedBucketLocations := []string{}
// 	for _, bl := range bucketLocations {
// 		parsedBucketLocations = append(parsedBucketLocations, strings.Split(bl, "://")[1])
// 	}

// 	return parsedBucketLocations
// }
