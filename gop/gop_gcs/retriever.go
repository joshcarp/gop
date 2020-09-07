package gcs

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/joshcarp/gop/app"
	"github.com/joshcarp/gop/gen/pkg/servers/gop"

	"cloud.google.com/go/storage"
)

type downloader func(bucket, object string) (io.Reader, error)

func (a GOP) Retrieve(res *gop.Object) error {
	filename := fmt.Sprintf("%s/%s@%s", res.Repo, res.Resource, res.Version)
	if err := downloadToString(a.downloader, a.AppConfig.CacheLocation, filename, &res.Content); err != nil {
		return err
	}
	filename = fmt.Sprintf("%s/%s.pb.json@%s", res.Repo, res.Resource, res.Version)
	if err := downloadToString(a.downloader, a.AppConfig.CacheLocation, filename, res.Processed); err != nil {
		return err
	}
	return nil
}

func downloadToString(download downloader, bucketName string, filename string, target *string) error {
	file, err := download(bucketName, filename)
	if err != nil {
		return err
	}
	if err := app.ScanIntoString(target, file); err != nil {
		return err
	}
	return nil
}

func download(bucket, object string) (io.Reader, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	return client.Bucket(bucket).Object(object).NewReader(ctx)
}