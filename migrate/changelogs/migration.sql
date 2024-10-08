-- Adiciona a coluna oficina_id em boicoins_transacoes
ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD COLUMN IF NOT EXISTS oficina_id UUID;

-- Adiciona a chave estrangeira para oficina_id em boicoins_transacoes
ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD CONSTRAINT IF NOT EXISTS boicoins_transacoes_oficina_id_fkey
FOREIGN KEY (oficina_id) REFERENCES boi_marronzinho.oficinas(id)
ON DELETE SET NULL;

-- Remove a constraint ck_transacao_unico se já existir
ALTER TABLE boi_marronzinho.boicoins_transacoes
DROP CONSTRAINT IF EXISTS ck_transacao_unico;

-- Adiciona uma nova constraint ck_transacao_unico
ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD CONSTRAINT ck_transacao_unico CHECK (
    (pedido_id IS NOT NULL AND doacao_id IS NULL AND oficina_id IS NULL) OR
    (pedido_id IS NULL AND doacao_id IS NOT NULL AND oficina_id IS NULL) OR
    (pedido_id IS NULL AND doacao_id IS NULL AND oficina_id IS NOT NULL)
);

-- Remove a coluna ponto_mapa_id e a chave estrangeira relacionada se existirem
ALTER TABLE boi_marronzinho.boicoins_transacoes
DROP COLUMN IF EXISTS ponto_mapa_id;

-- Remove qualquer chave estrangeira duplicada para usuario_id em boicoins_transacoes
ALTER TABLE boi_marronzinho.boicoins_transacoes
DROP CONSTRAINT IF EXISTS boicoins_transacoes_usuario_id_fkey;

-- Adiciona a chave estrangeira para usuario_id
ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD CONSTRAINT IF NOT EXISTS boicoins_transacoes_usuario_id_fkey
FOREIGN KEY (usuario_id) REFERENCES boi_marronzinho.usuarios(id);

-- Alterações na tabela ticket_oficina
-- Adiciona a coluna qrcode se ainda não existir
ALTER TABLE ticket_oficina
ADD COLUMN IF NOT EXISTS qrcode BYTEA;

-- Comenta a coluna qrcode
COMMENT ON COLUMN ticket_oficina.qrcode IS 'Dados da imagem do QR Code armazenados em formato binário.';

-- Renomeia a coluna data_inscricao para created_at
ALTER TABLE ticket_oficina
RENAME COLUMN data_inscricao TO created_at;

-- Adiciona uma constraint UNIQUE na coluna codigo, garantindo unicidade
ALTER TABLE ticket_oficina
ADD CONSTRAINT IF NOT EXISTS codigo_unico UNIQUE (codigo);
