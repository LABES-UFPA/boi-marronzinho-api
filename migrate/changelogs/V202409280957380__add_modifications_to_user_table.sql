ALTER TABLE usuarios
ADD COLUMN password_hash VARCHAR(255);

UPDATE usuarios
SET password_hash = 'dummy_hash'
WHERE password_hash IS NULL;

ALTER TABLE usuarios
ALTER COLUMN password_hash SET NOT NULL;

ALTER TABLE usuarios
ALTER COLUMN tipo_usuario SET DATA TYPE VARCHAR(20),
ALTER COLUMN tipo_usuario SET DEFAULT 'Cliente';

ALTER TABLE usuarios
ADD COLUMN last_login TIMESTAMP;

ALTER TABLE usuarios
ADD COLUMN password_reset_token VARCHAR(255);

ALTER TABLE usuarios
ADD COLUMN password_reset_expires TIMESTAMP;

ALTER TABLE usuarios
ADD COLUMN deleted_at TIMESTAMP;

ALTER TABLE usuarios
RENAME COLUMN criado_em TO created_at;

ALTER TABLE usuarios
ALTER COLUMN created_at SET DEFAULT CURRENT_TIMESTAMP;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name='usuarios' AND column_name='updated_at'
    ) THEN
        ALTER TABLE usuarios ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
    END IF;
END $$;

ALTER TABLE usuarios
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;

ALTER TABLE usuarios
ADD COLUMN first_name VARCHAR(50) NOT NULL,
ADD COLUMN last_name VARCHAR(50) NOT NULL;

UPDATE usuarios
SET first_name = split_part(nome, ' ', 1),
    last_name = CASE
        WHEN length(split_part(nome, ' ', 2)) > 0 THEN substring(nome from position(' ' in nome) + 1)
        ELSE 'Sobrenome não informado'
    END;

ALTER TABLE usuarios
DROP COLUMN nome;
