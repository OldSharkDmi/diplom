-- migrations/014_occupancy.sql
CREATE OR REPLACE VIEW v_train_occupancy AS
SELECT t.uid,
       s.delay_min,
       o.occupancy,
       s.updated_at
FROM   train_statuses s
           JOIN   trains t          ON t.train_run_id = s.train_run_id
           LEFT   JOIN occupancy_predictions o USING (train_run_id);
