CREATE TYPE tipo_usuario_enum AS ENUM (
    'Cliente',
    'Administrador',
    'Gringo'
);

COMMENT ON TYPE tipo_usuario_enum IS 'Define os tipos de usuários permitidos no sistema.';

CREATE TYPE status_pedido_enum AS ENUM (
    '1',  -- Em andamento
    '2',  -- Concluído
    '3',  -- Cancelado
    '4'   -- Pendente de pagamento
);

COMMENT ON TYPE status_pedido_enum IS 'Tipo ENUM que define os possíveis status dos pedidos.';

-- Enum para Ações Administrativas
CREATE TYPE acao_adm_enum AS ENUM (
    'validacao_doacao',
    'validacao_pagamento',
    'validacao_inscricao_oficina',
    'outro'
);

COMMENT ON TYPE acao_adm_enum IS 'Tipo ENUM que define as possíveis ações administrativas.';

-- Tabela de Usuários
CREATE TABLE usuarios (
    id UUID PRIMARY KEY,
    nome VARCHAR(100),
    email VARCHAR(100) UNIQUE,
    tipo_usuario tipo_usuario_enum DEFAULT 'Cliente',
    saldo_boicoins DECIMAL(10, 2) DEFAULT 0.00,
    idioma_preferido VARCHAR(10) DEFAULT 'pt',
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE usuarios IS 'Tabela que armazena as informações dos usuários, incluindo tipo e saldo de Boicoins.';

-- Tabelas de Normalização para Estado e Cidade
CREATE TABLE estados (
    id UUID PRIMARY KEY,
    nome VARCHAR(100)
);

CREATE TABLE cidades (
    id UUID PRIMARY KEY,
    estado_id UUID REFERENCES estados(id),
    nome VARCHAR(100)
);

-- Tabela de Pontos no Mapa
CREATE TABLE pontos_mapa (
    id UUID PRIMARY KEY,
    nome VARCHAR(100),
    descricao TEXT,
    latitude DECIMAL(10, 6),
    longitude DECIMAL(10, 6),
    imagem_url VARCHAR(255),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE pontos_mapa IS 'Tabela que armazena os pontos de coleta ou distribuição exibidos no mapa do sistema.';

-- Tabela de Endereços (endereços dos usuários)
CREATE TABLE enderecos (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    logradouro VARCHAR(255),
    numero VARCHAR(10),
    complemento VARCHAR(100),
    bairro VARCHAR(100),
    cidade_id UUID REFERENCES cidades(id),
    cep VARCHAR(20),
    pais VARCHAR(100),
    tipo_endereco VARCHAR(50),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE enderecos IS 'Tabela que armazena os endereços dos usuários.';

-- Tabela de Produtos
CREATE TABLE produtos (
    id UUID PRIMARY KEY,
    nome VARCHAR(100),
    descricao TEXT,
    preco_boicoins DECIMAL(10, 2),
    preco_real DECIMAL(10, 2),
    quantidade_em_estoque INT DEFAULT 0,
    imagem_url VARCHAR(255),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE produtos IS 'Tabela que armazena os produtos disponíveis na loja.';

-- Tabela de Itens de Doação
CREATE TABLE itens_doacao (
    id UUID PRIMARY KEY,
    nome_item VARCHAR(100),
    descricao TEXT,
    unidade_medida VARCHAR(50),
    boicoins_por_unidade DECIMAL(10, 2)
);

COMMENT ON TABLE itens_doacao IS 'Tabela que lista os tipos de itens aceitos como doação.';

-- Tabela de Oficinas (Eventos)
CREATE TABLE oficinas (
    id UUID PRIMARY KEY,
    nome VARCHAR(100),
    descricao TEXT,
    preco_boicoins DECIMAL(10, 2),
    preco_real DECIMAL(10, 2),
    data_evento TIMESTAMP,
    limite_participantes INT,
    participantes_atual INT DEFAULT 0,  -- Adicionando campo para rastrear participantes atuais
    ponto_mapa_id UUID REFERENCES pontos_mapa(id),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE oficinas IS 'Tabela que armazena as oficinas e eventos oferecidos pela comunidade.';

-- Tabela de Doações
CREATE TABLE doacoes (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    item_doacao_id UUID REFERENCES itens_doacao(id),
    quantidade DECIMAL(10, 2),
    boicoins_recebidos DECIMAL(10, 2),
    data_doacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE doacoes IS 'Tabela que armazena as doações realizadas pelos usuários.';

CREATE TABLE ticket_produto (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    produto_id UUID REFERENCES produtos(id),
    codigo VARCHAR(100),  -- Código de validação ou QR Code
    quantidade INTEGER DEFAULT 1,
    data_compra TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ticket_produto IS 'Tabela que armazena os tickets de compra de produtos usando Boicoins.';

CREATE TABLE ticket_oficina (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    oficina_id UUID REFERENCES oficinas(id),
    codigo VARCHAR(100),  -- Código de validação ou QR Code
    data_inscricao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ticket_oficina IS 'Tabela que armazena os tickets de inscrição em oficinas.';

CREATE TABLE pedidos (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    produto_id UUID REFERENCES produtos(id),
    oficina_id UUID REFERENCES oficinas(id),
    status_pedido status_pedido_enum,
    endereco_id UUID REFERENCES enderecos(id),
    ponto_mapa_id UUID REFERENCES pontos_mapa(id),
    boicoins_usados DECIMAL(10, 2) DEFAULT 0.00,
    preco_real_usado DECIMAL(10, 2),
    quantidade INTEGER DEFAULT 1,
    codigo VARCHAR(100),
    data_pedido TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_conclusao TIMESTAMP 
);

COMMENT ON TABLE pedidos IS 'Tabela que armazena os pedidos de produtos e inscrições em oficinas feitos pelos usuários.';

-- Tabela de Transações de BoiCoins
CREATE TABLE boicoins_transacoes (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    quantidade DECIMAL(10, 2),
    tipo_transacao VARCHAR(50),  
    descricao TEXT,
    data_transacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    pedido_id UUID REFERENCES pedidos(id) NULL,  
    doacao_id UUID REFERENCES doacoes(id) NULL,
    ponto_mapa_id UUID REFERENCES pontos_mapa(id)
);

ALTER TABLE boicoins_transacoes
ADD CONSTRAINT ck_transacao_unico CHECK (
    (pedido_id IS NOT NULL AND doacao_id IS NULL) OR (pedido_id IS NULL AND doacao_id IS NOT NULL)
);

COMMENT ON TABLE boicoins_transacoes IS 'Tabela que armazena as transações de Boicoins, como doações e compras.';

-- Histórico de Participantes nas Oficinas
CREATE TABLE participantes_oficinas (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    oficina_id UUID REFERENCES oficinas(id),
    data_inscricao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE participantes_oficinas IS 'Tabela que registra os participantes inscritos em oficinas.';

-- Tabela de Validações de Administradores
CREATE TABLE validacoes_adm (
    id UUID PRIMARY KEY,
    administrador_id UUID REFERENCES usuarios(id),
    acao acao_adm_enum,
    descricao TEXT,
    data_validacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE validacoes_adm IS 'Tabela que armazena as ações de validação realizadas por administradores.';

-- Índices sugeridos para melhorar desempenho
CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_pedidos_usuario_id ON pedidos(usuario_id);
CREATE INDEX idx_transacoes_usuario_id ON boicoins_transacoes(usuario_id);
CREATE INDEX idx_oficinas_data_evento ON oficinas(data_evento);
CREATE INDEX idx_pontos_mapa_localizacao ON pontos_mapa(latitude, longitude);
