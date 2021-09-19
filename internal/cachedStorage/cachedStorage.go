package cachedStorage

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/vzau/thoth/internal/controllers/cdn"
	"github.com/vzau/thoth/pkg/cache"
	"github.com/vzau/thoth/pkg/storage"
	dbTypes "github.com/vzau/types/database"
)

func GetCachedFile(file dbTypes.File) (*bytes.Reader, error) {
	cacheData, exists := cache.Cache.Get(fmt.Sprintf("cdn/%d", file.ID))
	if !exists {
		data, err := storage.GetFileStream(file.Filename)
		if err != nil {
			return nil, err
		}

		defer data.Close()

		buff, err := ioutil.ReadAll(data)
		if err != nil {
			return nil, err
		}
		reader := bytes.NewReader(buff)
		cache.Cache.Set(fmt.Sprintf("cdn/%d", file.ID), cdn.FileCache{
			File: file,
			Data: buff,
		}, 0)
		return reader, nil
	} else {
		data := cdn.FileCache{}
		err := data.UnmarshalBinary([]byte(cacheData.(string)))
		if err != nil {
			cache.Cache.Delete(fmt.Sprintf("cdn/%d", file.ID))
			return GetCachedFile(file)
		}

		return bytes.NewReader(data.Data), nil
	}
}
