-- Adiciona a coluna para armazenar o hash da senha, permitindo nulos temporariamente
ALTER TABLE usuarios
ADD COLUMN password_hash VARCHAR(255);

-- Preencher a coluna password_hash com um valor padrão para registros existentes
-- ATENÇÃO: Substitua 'dummy_hash' por algo apropriado se desejar obrigar a redefinição de senha.
UPDATE usuarios
SET password_hash = 'dummy_hash'
WHERE password_hash IS NULL;

-- Altera a coluna para NOT NULL após todos os valores serem preenchidos
ALTER TABLE usuarios
ALTER COLUMN password_hash SET NOT NULL;

-- Altera a coluna tipo_usuario para ter um tamanho maior e ser mais flexível
ALTER TABLE usuarios
ALTER COLUMN tipo_usuario SET DATA TYPE VARCHAR(20),
ALTER COLUMN tipo_usuario SET DEFAULT 'Cliente';

-- Adiciona a coluna last_login para registrar o último login do usuário
ALTER TABLE usuarios
ADD COLUMN last_login TIMESTAMP;

-- Adiciona a coluna password_reset_token para armazenar o token de redefinição de senha
ALTER TABLE usuarios
ADD COLUMN password_reset_token VARCHAR(255);

-- Adiciona a coluna password_reset_expires para registrar a data de expiração do token de redefinição
ALTER TABLE usuarios
ADD COLUMN password_reset_expires TIMESTAMP;

-- Adiciona a coluna deleted_at para gerenciar a deleção lógica (soft delete)
ALTER TABLE usuarios
ADD COLUMN deleted_at TIMESTAMP;

-- Atualiza o valor padrão da coluna criado_em para usar o CURRENT_TIMESTAMP
ALTER TABLE usuarios
ALTER COLUMN criado_em SET DEFAULT CURRENT_TIMESTAMP;

-- Adiciona a coluna updated_at se ela não existir
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

-- (Opcional) Atualiza o valor padrão da coluna updated_at se ela já existir
ALTER TABLE usuarios
ALTER COLUMN updated_at SET DEFAULT CURRENT_TIMESTAMP;
