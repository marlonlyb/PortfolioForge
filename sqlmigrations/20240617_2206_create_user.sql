CREATE TABLE users (
    id UUID NOT NULL,
    email VARCHAR(254) NOT NULL,
    password VARCHAR(72) NOT NULL,
    is_admin BOOL NOT NULL DEFAULT FALSE,
    details JSONB NOT NULL,
    created_at INTEGER NOT NULL DEFAULT EXTRACT(EPOCH FROM now())::int,
    updated_at INTEGER,
    CONSTRAINT user_id_pk PRIMARY KEY (id),
    CONSTRAINT user_email_uk UNIQUE (email)
);

/* Éste comando despues de crear la tabla */
COMMENT ON TABLE users IS 'Storage the admins and customers for the eMLB';

/* por ahora no hay creación de usuario con is_admin = true */
UPDATE public.users SET updated_at = 1722009915 WHERE id = '34847678-4b6a-11ef-bebe-1826499730cf';