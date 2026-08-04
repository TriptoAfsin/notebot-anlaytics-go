// One-time migration: notebot analytics MySQL  ->  notebot v2 Postgres (schema `analytics`).
// Creates the analytics.* tables and copies all rows, preserving ids and resetting sequences
// so the Go API's id-less INSERTs keep auto-incrementing. Skips the unrelated *_thesis tables.
// Idempotent: re-running DROPs + recreates the analytics tables.
//
// Run:  MYSQL_URL=... PG_URL=... node migrations/mysql-to-postgres.js
const mysql = require("mysql2/promise");
const { Client } = require("pg");

const MYSQL_URL = process.env.MYSQL_URL;
const PG_URL = process.env.PG_URL;

// MySQL table -> { pg DDL (analytics schema), column map (mysqlCol: pgCol) }
const TABLES = {
  subnamedb: {
    ddl: `id SERIAL PRIMARY KEY, sub_name varchar(100), count integer, level integer`,
    cols: { id: "id", sub_name: "sub_name", count: "count", level: "level" },
  },
  labsdb: {
    ddl: `id SERIAL PRIMARY KEY, lab_name varchar(100), level integer, count integer`,
    cols: { id: "id", lab_name: "lab_name", level: "level", count: "count" },
  },
  bot_daily_report: {
    ddl: `id BIGSERIAL PRIMARY KEY, date date, count integer DEFAULT 0, platform text`,
    cols: { id: "id", date: "date", count: "count", platform: "platform" },
  },
  bot_daily_stat: {
    ddl: `date date PRIMARY KEY, count integer DEFAULT 0, platform text`,
    cols: { date: "date", count: "count", platform: "platform" },
    noSeq: true,
  },
  game_hof: {
    ddl: `id SERIAL PRIMARY KEY, date date, score integer, email varchar(100), user_name varchar(50)`,
    cols: { id: "id", date: "date", score: "score", email: "email", user_name: "user_name" },
  },
  // Postgres folds the unquoted identifier in the Go query (game_hof_noteDino) to lowercase,
  // so the physical table must be lowercase to resolve.
  game_hof_noteDino: {
    pgName: "game_hof_notedino",
    ddl: `id SERIAL PRIMARY KEY, date date, score integer, email varchar(100), user_name varchar(50)`,
    cols: { id: "id", date: "date", score: "score", email: "email", user_name: "user_name" },
  },
  app_err_logs: {
    ddl: `id SERIAL PRIMARY KEY, date date, log text, os varchar(250), email varchar(100)`,
    cols: { id: "id", date: "date", log: "log", os: "os", email: "email" },
  },
  missed_words_table: {
    ddl: `id SERIAL PRIMARY KEY, missed_words varchar(255)`,
    cols: { id: "id", missed_words: "missed_words" },
  },
  app_users: {
    // MySQL column is imgUrl (camelCase); normalize to img_url (matches the Go INSERT).
    ddl: `id SERIAL PRIMARY KEY, email varchar(100), uni_id varchar(50), batch integer, dept varchar(50), role varchar(50), img_url varchar(255)`,
    cols: { id: "id", email: "email", uni_id: "uni_id", batch: "batch", dept: "dept", role: "role", imgUrl: "img_url" },
  },
  new_user: {
    ddl: `id SERIAL PRIMARY KEY, user_count integer`,
    cols: { id: "id", user_count: "user_count" },
  },
};

(async () => {
  const my = await mysql.createConnection(MYSQL_URL);
  const pg = new Client({ connectionString: PG_URL, ssl: { rejectUnauthorized: false } });
  await pg.connect();
  await pg.query("CREATE SCHEMA IF NOT EXISTS analytics");

  for (const [myName, spec] of Object.entries(TABLES)) {
    const pgName = spec.pgName || myName;
    const myCols = Object.keys(spec.cols);
    const pgCols = Object.values(spec.cols);

    await pg.query(`DROP TABLE IF EXISTS analytics."${pgName}" CASCADE`);
    await pg.query(`CREATE TABLE analytics."${pgName}" (${spec.ddl})`);

    const [rows] = await my.query("SELECT " + myCols.map((c) => "`" + c + "`").join(",") + " FROM `" + myName + "`");
    let inserted = 0;
    const BATCH = 500;
    for (let i = 0; i < rows.length; i += BATCH) {
      const chunk = rows.slice(i, i + BATCH);
      const values = [];
      const params = [];
      let p = 1;
      for (const r of chunk) {
        values.push(`(${myCols.map(() => `$${p++}`).join(",")})`);
        for (const c of myCols) params.push(r[c]);
      }
      await pg.query(
        `INSERT INTO analytics."${pgName}" (${pgCols.map((c) => `"${c}"`).join(",")}) VALUES ${values.join(",")}`,
        params
      );
      inserted += chunk.length;
    }

    // reset the id sequence so future id-less INSERTs continue from MAX(id)+1
    if (!spec.noSeq) {
      await pg.query(
        `SELECT setval(pg_get_serial_sequence('analytics."${pgName}"','id'), GREATEST((SELECT COALESCE(MAX(id),0) FROM analytics."${pgName}"),1))`
      );
    }
    console.log(`${myName} -> analytics.${pgName}: ${inserted} rows`);
  }

  await my.end();
  await pg.end();
  console.log("migration complete");
})().catch((e) => { console.error("MIGRATION FAILED:", e.message); process.exit(1); });
