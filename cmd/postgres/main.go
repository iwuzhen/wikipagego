/*
 * @Version: 0.0.1
 * @Author: ider
 * @Date: 2020-05-11 14:45:41
 * @LastEditors: ider
 * @LastEditTime: 2020-08-04 19:57:04
 * @Description: 
 */
 package main

 import (
		 "github.com/caarlos0/env"
		 "github.com/jmoiron/sqlx"
 
		 // "os"
		 // "database/sql"
		 "context"
		 "time"
 
		 _ "github.com/lib/pq"
		 log "github.com/sirupsen/logrus"
		 "go.mongodb.org/mongo-driver/bson"
		 "go.mongodb.org/mongo-driver/mongo"
		 "go.mongodb.org/mongo-driver/mongo/options"
 )
 
 type config struct {
		 Home         string        `env:"HOME"`
		 Port         int           `env:"PORT" envDefault:"3000"`
		 IsProduction bool          `env:"PRODUCTION"`
		 Hosts        []string      `env:"HOSTS" envSeparator:":"`
		 Duration     time.Duration `env:"DURATION"`
		 TempFolder   string        `env:"TEMP_FOLDER" envDefault:"${HOME}/tmp" envExpand:"true"`
 
		 PgUri    string `env:"PGURI" envDefault:"postgres://postgres:postgres@192.168.1.223/small_world?sslmode=disable"`
		 MongoUri string `env:"MONGOURI" envDefault:"mongodb://192.168.1.220:29001"`
 }
 
 func init() {
		 // Log as JSON instead of the default ASCII formatter.
		 log.SetFormatter(&log.TextFormatter{})
		 // Output to stdout instead of the default stderr
		 // Can be any io.Writer, see below for File example
		 //   log.SetOutput(os.Stdout)
		 // Only log the warning severity or above.
		 log.SetLevel(log.InfoLevel)
 }
 func checkErr(err error) {
	 if err != nil {
			 log.Fatal("ERROR:", err)
	 }
 }
 
 func connectPostgres(cfg config) *sqlx.DB {
	 db, err := sqlx.Open("postgres", cfg.PgUri)
	 checkErr(err)
	 err = db.Ping()
	 if err != nil {
			 log.Fatal("ping 失败", err)
	 }
	 return db
 }
 
 func connectMongo(cfg config) *mongo.Collection {
	 ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	 client, _ := mongo.NewClient(options.Client().ApplyURI(cfg.MongoUri))
	 client.Connect(ctx)
	 collection := client.Database("small_world").Collection("category_undirect_trend_level3")
	 return collection
 }
 func initDB(db *sqlx.DB, table_name string) {
 
	 // 初始化表
	 sql := `
	 DROP TABLE IF EXISTS ` + table_name + `;
 CREATE TABLE ` + table_name + `(
 id SERIAL PRIMARY KEY,
	 average_distance double precision NOT NULL,
	 clustering_coefficient double precision NOT NULL,
	 year smallint NOT NULL,
	 method VARCHAR(40) NOT NULL,
	 name VARCHAR(40) NOT NULL,
	 number_node INT NOT NULL,
	 number_edge INT NOT NULL
 );
 comment on column ` + table_name + `.id is '主键ID，自增';
 CREATE INDEX undirect_graph_year on ` + table_name + `(year);
 CREATE INDEX undirect_graph_method on ` + table_name + `(method);
 CREATE INDEX undirect_graph_name on ` + table_name + `(name);
 CREATE UNIQUE INDEX undirect_graph_name_year_method_numbernode on ` + table_name + `(name,year,method,number_node);
 `
	 db.Exec(sql)
 
 }
 func main() {
 
	 ctx, _ := context.WithTimeout(context.Background(), 360*time.Second)
 
	 table_name := "undirect_graph"
 
	 cfg := config{}
	 if err := env.Parse(&cfg); err != nil {
			 log.Warning("%+v", err)
	 }
	 log.Info(cfg.TempFolder)
	 db := connectPostgres(cfg)
	 defer db.Close()
	 initDB(db, table_name)
 
	 collection := connectMongo(cfg)
 
	 // var result bson.M
	 // collection.FindOne(ctx, bson.D{}).Decode(&result)
	 cur, err := collection.Find(ctx, bson.D{})
	 checkErr(err)
	 defer cur.Close(ctx)
 
	 sql := `insert into ` + table_name + ` (average_distance, clustering_coefficient, year, method, name, number_node, number_edge)
			 values ($1, $2, $3, $4, $5, $6, $7)`
	 tx, _ := db.Beginx()
	 stmt, _ := tx.Prepare(sql)
	 i := 1
	 for cur.Next(ctx) {
			 var result bson.M
			 err := cur.Decode(&result)
			 checkErr(err)
			 _, err = stmt.Exec(result["ad"], result["cc"], result["y"], result["m"], result["n"], result["nv"], result["ne"])
			 checkErr(err)
			 if i%1000 == 0 {
					 tx.Commit()
					 tx, _ = db.Beginx()
					 stmt, _ = tx.Prepare(sql)
					 log.Info(i)
			 }
			 i++
	 }
	 tx.Commit()
	 log.Info("commit")
 
 }
 