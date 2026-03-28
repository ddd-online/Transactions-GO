-- 创建账本表
CREATE TABLE IF NOT EXISTS tbl_billadm_ledger
(
    id         TEXT PRIMARY KEY,
    name       TEXT    NOT NULL,
    created_at INTEGER NOT NULL,
    updated_at INTEGER NOT NULL
);

-- 创建交易记录表
CREATE TABLE IF NOT EXISTS tbl_billadm_transaction_record
(
    transaction_id   TEXT PRIMARY KEY,
    ledger_id        TEXT    NOT NULL,
    price            INTEGER NOT NULL,
    transaction_type TEXT    NOT NULL,
    category         TEXT    NOT NULL,
    description      TEXT,
    flags            TEXT    NOT NULL DEFAULT '',
    transaction_at   INTEGER NOT NULL,
    created_at       INTEGER NOT NULL,
    updated_at       INTEGER NOT NULL
);

-- 创建交易记录标签表
CREATE TABLE IF NOT EXISTS tbl_billadm_transaction_record_tag
(
    ledger_id      TEXT NOT NULL,
    transaction_id TEXT NOT NULL,
    tag            TEXT NOT NULL
);

-- 创建消费分类表
CREATE TABLE IF NOT EXISTS tbl_billadm_category
(
    name             TEXT PRIMARY KEY,
    scope            TEXT NOT NULL,
    transaction_type TEXT NOT NULL DEFAULT ''
);

-- 创建消费标签表
CREATE TABLE IF NOT EXISTS tbl_billadm_tag
(
    name     TEXT NOT NULL,
    scope    TEXT NOT NULL,
    category TEXT NOT NULL DEFAULT '',
    UNIQUE (name, scope, category)
);


