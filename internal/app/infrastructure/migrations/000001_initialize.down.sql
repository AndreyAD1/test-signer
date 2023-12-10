BEGIN;

ALTER TABLE test_details DROP CONSTRAINT signature_id;

DROP TABLE test_details, signatures;

COMMIT;