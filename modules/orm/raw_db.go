package orm

import (
  "database/sql"
  "errors"
  "github.com/hashicorp/errwrap"
  "github.com/wgentry22/agora/modules/logg"
  "github.com/wgentry22/agora/types/config"
  "sync"
  "time"

  _ "github.com/jackc/pgx/v4/stdlib"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
)

var (
  m                              sync.Mutex
  ErrDBConnectionNotInitialized  = errors.New("database connection not initialized")
  ErrDBConnectionFailed          = errors.New("database connection failed")
  ErrORMConnectionNotInitialized = errors.New("orm not initialized")
  ErrORMConnectionFailed         = errors.New("orm connection failed")
  rawDB                          *sql.DB
  ormInstance                    *orm
  timeFunc                       = func() time.Time {
    t, err := time.Parse("Mon Jan 01, 2006 15:04:05 -0700 MST", time.Now().String())
    if err != nil {
      return time.Now()
    }

    return t
  }
  loggerConfig = config.Logging{
    Level:       "trace",
    OutputPaths: nil,
    Fields: map[string]interface{}{
      "module": "orm",
    },
  }

  logger = logg.NewLogrusLogger(loggerConfig)
)

type orm struct {
  instance *gorm.DB
}

func GetRaw() *sql.DB {
  m.Lock()
  defer m.Unlock()

  if rawDB == nil {
    panic(ErrDBConnectionNotInitialized)
  }

  if ormInstance == nil {
    panic(ErrORMConnectionNotInitialized)
  }

  return rawDB
}

func UseConfig(conf config.DB) {
  m.Lock()
  defer m.Unlock()

  makeRawDB(conf)

  makeORM(rawDB, conf)
}

func UseLoggingConfig(conf config.Logging) {
  m.Lock()
  defer m.Unlock()

  logger = logg.NewLogrusLogger(conf)

  ormInstance.instance.Logger = logg.ForGorm(logger)
}

func makeORM(db *sql.DB, conf config.DB) {
  dialector := makeDialector(db, conf)
  gormDB, err := gorm.Open(dialector, &gorm.Config{
    FullSaveAssociations: true,
    NowFunc:              timeFunc,
    DisableAutomaticPing: true,
    CreateBatchSize:      25,
    Dialector:            dialector,
    Logger:               logg.ForGorm(logger),
  })

  if err != nil {
    panic(errwrap.Wrap(ErrORMConnectionFailed, err))
  }

  ormInstance = &orm{
    instance: gormDB,
  }
}

func makeRawDB(conf config.DB) {
  raw, err := sql.Open(config.DriverName(conf), config.ConnectionString(conf))
  if err != nil {
    panic(errwrap.Wrap(ErrDBConnectionFailed, err))
  }

  rawDB = raw
}

func makeDialector(db *sql.DB, conf config.DB) gorm.Dialector {
  return postgres.New(postgres.Config{
    DriverName:           config.DriverName(conf),
    DSN:                  config.ConnectionString(conf),
    PreferSimpleProtocol: true,
    Conn:                 db,
  })
}
