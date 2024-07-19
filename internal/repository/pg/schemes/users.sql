CREATE TABLE IF NOT EXISTS users (
    id               BIGINT NOT NULL,
    username         TEXT   DEFAULT '',
    firstname        TEXT   DEFAULT '',
    created_at       TEXT   DEFAULT '',
    is_admin         INT    DEFAULT 0,
    bot_state        TEXT   DEFAULT '',
    email            TEXT   DEFAULT '',
    ref              TEXT   DEFAULT '',
    lichka          TEXT   DEFAULT '',
    
    lives            INT    DEFAULT 0,
    step             TEXT   DEFAULT '0',

    lats_action_time TEXT   DEFAULT '',

    is_last_push     INT    DEFAULT 0,
    is_send_push     INT    DEFAULT 0,
    is_final         INT    DEFAULT 0,

    feedback_cnt     INT    DEFAULT 0,
    feedback_time    TEXT   DEFAULT '',

    not_del_mess_id  INT    DEFAULT 0,
    inst_link        TEXT   DEFAULT '',
    is_inst_push     INT    DEFAULT 0,

    PRIMARY KEY (id)
);

-------------------------------------------

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS lichka TEXT DEFAULT '';

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS not_del_mess_id INT DEFAULT 0;

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS is_inst_push INT DEFAULT 0;

ALTER TABLE users
  ADD COLUMN IF NOT EXISTS inst_link TEXT DEFAULT '';