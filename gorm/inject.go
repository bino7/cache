package gorm

import (
	"context"
	"fmt"
	"github.com/apex/log"
	. "github.com/go-redis/cache/v8"
	"gorm.io/gorm"
	"runtime"
	"strings"
	"sync"
	"time"
)

const DryRun = "cache:DryRun"
const SetCache = "cache:SetCache"

func InjectGormDB(db *gorm.DB, cache *Cache) {
	tryCatch(db.Callback().Create().After("gorm:after_create").Register("cache:after_create", afterCreate(cache)))
	tryCatch(db.Callback().Delete().After("gorm:after_delete").Register("cache:after_delete", afterDelete(cache)))
	tryCatch(db.Callback().Update().After("gorm:after_update").Register("cache:after_update", afterUpdate(cache)))
	tryCatch(db.Callback().Query().Before("gorm:query").Register("cache:before_query", beforeQuery(cache)))
	tryCatch(db.Callback().Query().After("gorm:query").Register("cache:after_query", afterQuery(cache)))
	tryCatch(db.Callback().Query().After("gorm:after_query").Register("cache:end_query", endQuery(cache)))
}

func tryCatch(err error) bool {
	if err != nil {
		log.Fatal(err.Error())
	}
	return true
}

func afterCreate(cache *Cache) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.RowsAffected > 0 {
			dest := db.Statement.Dest
			if item, ok := dest.(*Item); ok {
				err := cache.Set(item)
				if err != nil {
					log.Warn(err.Error())
				}
			}
		}
	}
}

func afterDelete(cache *Cache) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.RowsAffected > 0 {
		}
	}
}

func afterUpdate(cache *Cache) func(*gorm.DB) {
	return func(db *gorm.DB) {
		if db.RowsAffected > 0 {

		}
	}
}

func beforeQuery(cache *Cache) func(*gorm.DB) {
	return func(db *gorm.DB) {
		db.InstanceSet(DryRun, db.DryRun)
		db.DryRun = true
	}
}

var modelsPackage = "yxd_server/models"
var cacheKeyMap = new(sync.Map)

func getCacheKey(vars []interface{}) (string, bool) {
	key, ok := getMethodCacheKey()
	if !ok {
		return "", false
	}
	return key + ":" + valsToString(vars), true

}
func getMethodCacheKey() (string, bool) {
	pcs := make([]uintptr, 100)
	runtime.Callers(0, pcs)
	if modelsPackage != "" {
		for _, pc := range pcs {
			f := runtime.FuncForPC(pc)
			name := f.Name()
			if strings.HasPrefix(name, modelsPackage) {
				if v, ok := cacheKeyMap.Load(pc); ok {
					return v.(string), true
				}
				key := name[len(modelsPackage)+1:]
				cacheKeyMap.Store(pc, key)
				return key, true
			}
		}
	}
	return "", false
}
func valsToString(vars []interface{}) string {
	if len(vars) == 0 {
		return ""
	}
	builder := new(strings.Builder)
	for _, v := range vars {
		builder.WriteString(valToString(v))
		builder.WriteString(",")
	}
	return builder.String()
}
func valToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}
func afterQuery(cache *Cache) func(*gorm.DB) {
	return func(db *gorm.DB) {
		v, _ := db.InstanceGet(DryRun)
		db.DryRun = v.(bool)
		dest := db.Statement.Dest
		key, ok := getCacheKey(db.Statement.Vars)
		if !ok {
			key = db.Statement.SQL.String()
		}
		err := cache.Get(defaultContext(), key, dest)
		if err != nil {
			log.Warnf("CacheErr getting cache %v", err)
		} else {
			return
		}
		if !db.DryRun && db.Error == nil {
			rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, db.Statement.SQL.String(), db.Statement.Vars...)
			if err != nil {
				db.AddError(err)
				return
			}
			defer func() {
				db.AddError(rows.Close())
			}()
			db.InstanceSet(SetCache, key)
			gorm.Scan(rows, db, 0)
		}
	}
}

func endQuery(cache *Cache) func(*gorm.DB) {
	return func(db *gorm.DB) {
		key, ok := db.InstanceGet(SetCache)
		if ok {
			err := cache.Set(&Item{
				Key:   key.(string),
				Value: db.Statement.Dest,
				TTL:   100000 * time.Hour,
			})
			if err != nil {
				log.Warnf("CacheErr setting cache: %v", err)
			}
		}

	}
}

func defaultContext() context.Context {
	return context.Background()
}
