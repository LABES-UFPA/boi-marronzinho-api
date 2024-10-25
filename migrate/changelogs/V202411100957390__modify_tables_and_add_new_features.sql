-- Adicionando nova coluna 'validado' à tabela 'ticket_oficina'
ALTER TABLE boi_marronzinho.ticket_oficina
ADD COLUMN validado BOOLEAN DEFAULT FALSE;

COMMENT ON COLUMN boi_marronzinho.ticket_oficina.validado IS 'Indica se o elemento foi validado (TRUE) ou não (FALSE).';

-- Adicionando nova coluna 'imagem' à tabela 'oficinas'
ALTER TABLE boi_marronzinho.oficinas
ADD COLUMN imagem BYTEA;

-- Criando um novo tipo ENUM para o status de doação
CREATE TYPE status_doacao_enum AS ENUM (
    'pendente',
    'validada',
    'rejeitada'
);

-- Adicionando a coluna 'status' com o novo tipo ENUM à tabela 'troca'
ALTER TABLE boi_marronzinho.troca
ADD COLUMN status status_doacao_enum DEFAULT 'pendente';

-- Adicionando a coluna 'link_endereco' à tabela 'oficinas'
ALTER TABLE boi_marronzinho.oficinas
ADD COLUMN link_endereco VARCHAR(250);

-- Renomeando a tabela 'doacoes' para 'troca'
ALTER TABLE boi_marronzinho.doacoes
RENAME TO troca;

-- Renomeando a tabela 'itens_doacao' para 'item_troca'
ALTER TABLE boi_marronzinho.itens_doacao
RENAME TO item_troca;

-- Criando a nova tabela 'carrinho_itens'
CREATE TABLE boi_marronzinho.carrinho_itens (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES boi_marronzinho.usuarios(id) ON DELETE CASCADE,
    produto_id UUID REFERENCES boi_marronzinho.produtos(id) ON DELETE CASCADE,
    quantidade INT NOT NULL CHECK (quantidade > 0),
    preco_unitario DECIMAL(10, 2) NOT NULL CHECK (preco_unitario >= 0),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE boi_marronzinho.carrinho_itens IS 'Tabela que armazena os itens que os usuários adicionaram ao carrinho de compras.';
COMMENT ON COLUMN boi_marronzinho.carrinho_itens.preco_unitario IS 'Preço em Boicoins por unidade do produto no momento em que foi adicionado ao carrinho.';
COMMENT ON COLUMN boi_marronzinho.carrinho_itens.criado_em IS 'Data e hora em que o item foi adicionado ao carrinho.';

-- Ajustando o schema da tabela 'carrinho_itens'
ALTER TABLE boi_marronzinho.carrinho_itens
SET SCHEMA boi_marronzinho;

-- Criação da nova tabela 'evento'
CREATE TABLE boi_marronzinho.evento (
    id UUID PRIMARY KEY,
    nome VARCHAR(100) NOT NULL,
    descricao TEXT,
    data_evento TIMESTAMP NOT NULL,
    link_endereco TEXT,
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    atualizado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE boi_marronzinho.evento IS 'Tabela que armazena eventos do Boi Marronzinho, incluindo detalhes e datas.';
COMMENT ON COLUMN boi_marronzinho.evento.data_evento IS 'Data e hora em que o evento ocorrerá.';
COMMENT ON COLUMN boi_marronzinho.evento.link_endereco IS 'Link de endereço relacionado ao evento.';
COMMENT ON COLUMN boi_marronzinho.evento.criado_em IS 'Data e hora em que o registro do evento foi criado.';
COMMENT ON COLUMN boi_marronzinho.evento.atualizado_em IS 'Data e hora da última atualização do registro do evento.';

-- Função para atualizar automaticamente a coluna 'atualizado_em' na tabela 'evento'
CREATE OR REPLACE FUNCTION atualizar_data_evento() RETURNS TRIGGER AS $$
BEGIN
    NEW.atualizado_em = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger para executar a função antes de cada atualização na tabela 'evento'
CREATE TRIGGER trigger_atualizar_data_evento
BEFORE UPDATE ON boi_marronzinho.evento
FOR EACH ROW
EXECUTE FUNCTION atualizar_data_evento();

-- Removendo coluna obsoleta 'imagem' da tabela 'oficinas'
ALTER TABLE boi_marronzinho.oficinas
DROP COLUMN IF EXISTS imagem;

-- Adicionando a nova coluna 'imagem_url' à tabela 'oficinas'
ALTER TABLE boi_marronzinho.oficinas
ADD COLUMN imagem_url VARCHAR(255);

-- Removendo a coluna 'qr_code' obsoleta da tabela 'ticket_oficina'
ALTER TABLE boi_marronzinho.ticket_oficina
DROP COLUMN qr_code;

-- Adicionando nova coluna 'qr_code' na tabela 'ticket_oficina'
ALTER TABLE boi_marronzinho.ticket_oficina
ADD COLUMN qr_code VARCHAR(255);

-- Inserindo itens de doação na tabela 'item_troca'
INSERT INTO boi_marronzinho.item_troca (id, nome_item, descricao, unidade_medida, boicoins_por_unidade) VALUES
    ('b1a9f6b3-34f7-4b3d-843e-1d3e3cb5d12f', 'Óleo de cozinha usado', 'Óleo de cozinha reciclável, armazenado em garrafa PET de 2L', 'MILILITRO', 10.00),
    ('a3c5e2d9-43e2-4a5f-b9e6-6b5d33a7e2c1', 'Pote de vidro com tampa', 'Potes de vidro de qualquer tamanho, desde que acompanhados de tampa', 'UNIDADE', 5.00),
    ('d6f9a7c1-78b3-4d7e-b3e8-4d5b2a8d4e9d', 'Vasilhames plástico', 'Vasilhames plásticos diversos, como garrafas, potes e recipientes', 'UNIDADE', 3.00);
