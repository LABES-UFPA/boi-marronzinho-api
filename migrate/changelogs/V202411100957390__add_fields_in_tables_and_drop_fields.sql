ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD COLUMN IF NOT EXISTS oficina_id UUID;

ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD CONSTRAINT IF NOT EXISTS boicoins_transacoes_oficina_id_fkey
FOREIGN KEY (oficina_id) REFERENCES boi_marronzinho.oficinas(id)
ON DELETE SET NULL;

ALTER TABLE boi_marronzinho.boicoins_transacoes
DROP CONSTRAINT IF EXISTS ck_transacao_unico;

ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD CONSTRAINT ck_transacao_unico CHECK (
    (pedido_id IS NOT NULL AND doacao_id IS NULL AND oficina_id IS NULL) OR
    (pedido_id IS NULL AND doacao_id IS NOT NULL AND oficina_id IS NULL) OR
    (pedido_id IS NULL AND doacao_id IS NULL AND oficina_id IS NOT NULL)
);

ALTER TABLE boi_marronzinho.boicoins_transacoes
DROP COLUMN IF EXISTS ponto_mapa_id;

ALTER TABLE boi_marronzinho.boicoins_transacoes
DROP CONSTRAINT IF EXISTS boicoins_transacoes_usuario_id_fkey;

ALTER TABLE boi_marronzinho.boicoins_transacoes
ADD CONSTRAINT IF NOT EXISTS boicoins_transacoes_usuario_id_fkey
FOREIGN KEY (usuario_id) REFERENCES boi_marronzinho.usuarios(id);

ALTER TABLE ticket_oficina
ADD COLUMN IF NOT EXISTS qrcode BYTEA;

COMMENT ON COLUMN ticket_oficina.qrcode IS 'Dados da imagem do QR Code armazenados em formato binário.';

ALTER TABLE ticket_oficina
RENAME COLUMN data_inscricao TO created_at;

ALTER TABLE ticket_oficina
ADD CONSTRAINT IF NOT EXISTS codigo_unico UNIQUE (codigo);


ALTER TABLE boi_marronzinho.ticket_oficina
ADD COLUMN validado BOOLEAN DEFAULT FALSE;

COMMENT ON COLUMN boi_marronzinho.ticket_oficina.validado IS 'Indica se o elemento foi validado (TRUE) ou não (FALSE).';
