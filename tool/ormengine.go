package tool

import (
	"cloudRestaurant/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var DbEngine *Orm

type Orm struct {
	*xorm.Engine
}

func OrmEngine(cfg *Config) (*Orm, error) {
	db := cfg.Db
	conn := db.User + ":" + db.Password + "@tcp(" + db.Host + ":" + db.Port + ")/" + db.DbName + "?charset=" + db.Charset
	engine, err := xorm.NewEngine(db.Driver, conn)
	if err != nil {
		return nil, err
	}
	engine.ShowSQL(db.ShowSQL)

	// 将SmsCode映射为数据库中的表结构
	if err := engine.Sync2(new(model.SmsCode), new(model.Member)); err != nil {
		return nil, err
	}

	orm := new(Orm)
	orm.Engine = engine
	DbEngine = orm

	return orm, err
}
