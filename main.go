package main

import (
	"context"
	"flag"
	"os"
	"path"
	"sync"
	"time"
	"workerbee/storage"

	"github.com/gookit/slog"
)

func main() {
	// workerbee --root-location ./root s3://bucket2 s3://bucket3

	rootLocation := flag.String("loc", ".", "Pass root location of where you want the files to be downladed.")
	flag.Parse()
	bucketLocations := flag.Args()

	// Step 2: Create filesystem structure
	// err := createFileSystemStructure(*rootLocation, bucketLocations)
	// if err != nil {
	// 	slog.Errorf("%s", err)
	// 	os.Exit(-1)
	// }

	// Step 3: Infinitie loop to poll every n sseconds
	ctx := context.Background()

	for {
		slog.Infof("Waiting for 20 seconds")
		time.Sleep(20 * time.Second)

		for _, bl := range bucketLocations {
			fileList, err := storage.ListFilesGCS(ctx, bl)
			if err != nil {
				slog.Errorf("Error in fetching files from the bucket %s, %s", bl, err)
				continue
			}

			wg := sync.WaitGroup{}
			for _, f := range fileList {
				wg.Add(1)
				fileLocation := path.Join(*rootLocation, "incoming", bl, f)
				directory := path.Dir(fileLocation)
				err := os.MkdirAll(directory, 0755)
				if err != nil {
					slog.Errorf("Cannot create directory %s. Error - ", directory, err)
					continue
				}

				go func(ctx context.Context, bucket string, object string, fileLocation string) {
					defer wg.Done()
					err := storage.DownloadFile(ctx, bucket, object, fileLocation)
					if err != nil {
						slog.Errorf("Cannot download file %s for bucket %s. Error - ", object, bl, err)
					} else {
						slog.Infof("Downloaded file %s for bucket %s", object, bucket)
					}

				}(ctx, bl, f, fileLocation)
			}

			wg.Wait()

			slog.Infof("Completed download for bucket %s", bl)
		}

		// Write lockfile
		err := os.WriteFile(path.Join(*rootLocation, "lockfile"), []byte{}, 0755)
		if err != nil {
			slog.Errorf("Cannot write lockfile. Error - ", err)
			continue
		}

		// Delete main folder
		err = os.RemoveAll(path.Join(*rootLocation, "main"))
		if err != nil {
			slog.Errorf("Cannot remove main directory. Error - ", err)
			continue
		}

		// Rename incoming
		err = os.Rename(path.Join(*rootLocation, "incoming"), path.Join(*rootLocation, "main"))
		if err != nil {
			slog.Errorf("Cannot rename directory %s. Error - ", path.Join(*rootLocation, "incoming"), err)
			continue
		}

		// Delete lockfile
		err = os.Remove(path.Join(*rootLocation, "lockfile"))
		if err != nil {
			slog.Errorf("Cannot remove lockfile. Error - ", err)
			continue
		}

		slog.Infof("Completed downloads for all buckets")
	}

	// Step 4: Create go routines to fetch contents and get into an incoming directory

	// Step 5: Rename new directory and delete old
}
