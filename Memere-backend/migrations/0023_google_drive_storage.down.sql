-- Reverse 0023_google_drive_storage: drop the Drive object index, the credential
-- store, and the storage schema (tables first so the schema drops cleanly).
DROP TABLE IF EXISTS storage.google_credentials;
DROP TABLE IF EXISTS storage.objects;
DROP SCHEMA IF EXISTS storage;
