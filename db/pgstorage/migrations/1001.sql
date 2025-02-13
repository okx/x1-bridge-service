-- +migrate Down
DROP INDEX IF EXISTS sync.deposit_uidx;
DROP INDEX IF EXISTS sync.deposit_dest_addr_idx;
DROP INDEX IF EXISTS sync.claim_dest_addr_idx;

-- +migrate Up
CREATE UNIQUE INDEX IF NOT EXISTS deposit_uidx ON sync.deposit (network_id, deposit_cnt);
CREATE INDEX IF NOT EXISTS deposit_dest_addr_idx ON sync.deposit (dest_addr);
CREATE INDEX IF NOT EXISTS claim_dest_addr_idx ON sync.claim (dest_addr);