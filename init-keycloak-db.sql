CREATE DATABASE keycloak;
CREATE USER keycloak_user WITH PASSWORD 'keycloak_password';

-- Grant privileges on the database
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak_user;

-- Switch to the newly created database to apply schema-level permissions
\c keycloak

-- Grant usage and creation privileges on the public schema
GRANT USAGE ON SCHEMA public TO keycloak_user;
GRANT CREATE ON SCHEMA public TO keycloak_user;

-- Grant ownership of the public schema
ALTER SCHEMA public OWNER TO keycloak_user;
