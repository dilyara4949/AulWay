DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_tickets_user;
DROP INDEX IF EXISTS idx_routes_departure;
DROP INDEX IF EXISTS idx_routes_destination;
DROP INDEX IF EXISTS idx_payments_transaction;

DROP TABLE IF EXISTS settings;
DROP TABLE IF EXISTS support_requests;
DROP TABLE IF EXISTS faq;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS tickets;
DROP TABLE IF EXISTS routes;
DROP TABLE IF EXISTS buses;
DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS roles;
DROP TABLE IF EXISTS users;
