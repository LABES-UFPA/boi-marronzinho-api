-- Enum para Tipo de Usuário
CREATE TYPE tipo_usuario_enum AS ENUM (
    'Cliente',
    'Administrador',
    'Gringo'
);

COMMENT ON TYPE tipo_usuario_enum IS 'Define os tipos de usuários permitidos no sistema.';

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
COMMENT ON COLUMN usuarios.id IS 'Identificador único do usuário (UUID fornecido manualmente).';
COMMENT ON COLUMN usuarios.nome IS 'Nome completo do usuário.';
COMMENT ON COLUMN usuarios.email IS 'Email do usuário, utilizado para login.';
COMMENT ON COLUMN usuarios.tipo_usuario IS 'Tipo de usuário: Cliente, Administrador, Gringo.';
COMMENT ON COLUMN usuarios.saldo_boicoins IS 'Saldo de Boicoins disponível na conta do usuário.';
COMMENT ON COLUMN usuarios.idioma_preferido IS 'Idioma preferido para o usuário: pt para Português, en para Inglês.';
COMMENT ON COLUMN usuarios.criado_em IS 'Data e hora de criação do registro do usuário.';

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
COMMENT ON COLUMN enderecos.id IS 'Identificador único do endereço (UUID fornecido manualmente).';
COMMENT ON COLUMN enderecos.usuario_id IS 'Referência ao usuário dono do endereço.';
COMMENT ON COLUMN enderecos.logradouro IS 'Logradouro do endereço (Rua, Avenida, etc.).';
COMMENT ON COLUMN enderecos.numero IS 'Número do endereço (casa, apartamento, etc.).';
COMMENT ON COLUMN enderecos.complemento IS 'Complemento do endereço, se aplicável (Apto, bloco, etc.).';
COMMENT ON COLUMN enderecos.bairro IS 'Bairro do endereço.';
COMMENT ON COLUMN enderecos.cidade_id IS 'Referência à cidade onde o endereço está localizado.';
COMMENT ON COLUMN enderecos.cep IS 'CEP ou código postal do endereço.';
COMMENT ON COLUMN enderecos.pais IS 'País do endereço.';
COMMENT ON COLUMN enderecos.tipo_endereco IS 'Tipo de endereço, como residencial ou comercial.';
COMMENT ON COLUMN enderecos.criado_em IS 'Data e hora de criação do registro do endereço.';

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
COMMENT ON COLUMN produtos.id IS 'Identificador único do produto (UUID fornecido manualmente).';
COMMENT ON COLUMN produtos.nome IS 'Nome do produto disponível na loja.';
COMMENT ON COLUMN produtos.descricao IS 'Descrição detalhada do produto.';
COMMENT ON COLUMN produtos.preco_boicoins IS 'Preço do produto em Boicoins.';
COMMENT ON COLUMN produtos.preco_real IS 'Preço do produto em dinheiro real (R$).';
COMMENT ON COLUMN produtos.quantidade_em_estoque IS 'Quantidade disponível em estoque do produto.';
COMMENT ON COLUMN produtos.imagem_url IS 'URL da imagem que representa o produto.';
COMMENT ON COLUMN produtos.criado_em IS 'Data e hora de criação do registro do produto.';

-- Definição de tipo ENUM para o status do pedido
CREATE TYPE status_pedido_enum AS ENUM (
    '1',  -- Em andamento
    '2',  -- Concluído
    '3',  -- Cancelado
    '4'   -- Pendente de pagamento
);

COMMENT ON TYPE status_pedido_enum IS 'Tipo ENUM que define os possíveis status dos pedidos.';

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
COMMENT ON COLUMN pontos_mapa.id IS 'Identificador único do ponto no mapa (UUID fornecido manualmente).';
COMMENT ON COLUMN pontos_mapa.nome IS 'Nome do ponto no mapa (ex: Ponto de Coleta).';
COMMENT ON COLUMN pontos_mapa.descricao IS 'Descrição do ponto, como seu propósito ou horário de funcionamento.';
COMMENT ON COLUMN pontos_mapa.latitude IS 'Latitude do ponto para exibição no mapa.';
COMMENT ON COLUMN pontos_mapa.longitude IS 'Longitude do ponto para exibição no mapa.';
COMMENT ON COLUMN pontos_mapa.imagem_url IS 'URL da imagem associada ao ponto, como uma foto ou ícone.';
COMMENT ON COLUMN pontos_mapa.criado_em IS 'Data e hora de criação do ponto no sistema.';

-- Tabela de Itens de Doação (lista os itens possíveis para doação)
CREATE TABLE itens_doacao (
    id UUID PRIMARY KEY,
    nome_item VARCHAR(100),
    descricao TEXT,
    unidade_medida VARCHAR(50),
    boicoins_por_unidade DECIMAL(10, 2)
);

COMMENT ON TABLE itens_doacao IS 'Tabela que lista os tipos de itens aceitos como doação.';
COMMENT ON COLUMN itens_doacao.id IS 'Identificador único do item de doação (UUID fornecido manualmente).';
COMMENT ON COLUMN itens_doacao.nome_item IS 'Nome do item aceito para doação (Óleo, Garrafa PET, etc.).';
COMMENT ON COLUMN itens_doacao.descricao IS 'Descrição do item de doação.';
COMMENT ON COLUMN itens_doacao.unidade_medida IS 'Unidade de medida do item (Litros, Unidades, Quilos, etc.).';
COMMENT ON COLUMN itens_doacao.boicoins_por_unidade IS 'Quantidade de Boicoins que o usuário recebe por unidade doada.';

-- Tabela de Doações (generalizada)
CREATE TABLE doacoes (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    item_doacao_id UUID REFERENCES itens_doacao(id),
    quantidade DECIMAL(10, 2),
    boicoins_recebidos DECIMAL(10, 2),
    data_doacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE doacoes IS 'Tabela que armazena as doações realizadas pelos usuários.';
COMMENT ON COLUMN doacoes.id IS 'Identificador único da doação (UUID fornecido manualmente).';
COMMENT ON COLUMN doacoes.usuario_id IS 'Referência ao usuário que realizou a doação.';
COMMENT ON COLUMN doacoes.item_doacao_id IS 'Referência ao item de doação.';
COMMENT ON COLUMN doacoes.quantidade IS 'Quantidade do item doado.';
COMMENT ON COLUMN doacoes.boicoins_recebidos IS 'Quantidade de Boicoins recebida em troca pela doação.';
COMMENT ON COLUMN doacoes.data_doacao IS 'Data e hora em que a doação foi realizada.';

-- Tabela de Pedidos (compra de produtos)
CREATE TABLE pedidos (
    id UUID PRIMARY KEY,
    usuario_id UUID REFERENCES usuarios(id),
    produto_id UUID REFERENCES produtos(id),
    status_pedido status_pedido_enum,
    endereco_id UUID REFERENCES enderecos(id),
    ponto_mapa_id UUID REFERENCES pontos_mapa(id), -- Conexão com a tabela pontos_mapa para entrega/retirada
    boicoins_usados DECIMAL(10, 2) DEFAULT 0.00,
    preco_real_usado DECIMAL(10, 2),
    data_pedido TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    data_conclusao TIMESTAMP
);

COMMENT ON TABLE pedidos IS 'Tabela que armazena os pedidos de produtos feitos pelos usuários.';
COMMENT ON COLUMN pedidos.id IS 'Identificador único do pedido (UUID fornecido manualmente).';
COMMENT ON COLUMN pedidos.usuario_id IS 'Referência ao usuário que fez o pedido.';
COMMENT ON COLUMN pedidos.produto_id IS 'Referência ao produto solicitado no pedido.';
COMMENT ON COLUMN pedidos.status_pedido IS 'Status atual do pedido (em andamento, concluído, cancelado, etc.).';
COMMENT ON COLUMN pedidos.endereco_id IS 'Referência ao endereço de entrega associado ao pedido.';
COMMENT ON COLUMN pedidos.ponto_mapa_id IS 'Referência ao ponto de entrega ou retirada no mapa.';
COMMENT ON COLUMN pedidos.boicoins_usados IS 'Quantidade de Boicoins usados no pedido.';
COMMENT ON COLUMN pedidos.preco_real_usado IS 'Valor em dinheiro real usado no pedido.';
COMMENT ON COLUMN pedidos.data_pedido IS 'Data e hora de criação do pedido.';
COMMENT ON COLUMN pedidos.data_conclusao IS 'Data e hora de conclusão do pedido.';

-- Tabela de Transações de BoiCoins (histórico de transações)
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

-- Restrição para garantir que apenas um dos campos (pedido_id ou doacao_id) seja preenchido
ALTER TABLE boicoins_transacoes
ADD CONSTRAINT ck_transacao_unico CHECK (
    (pedido_id IS NOT NULL AND doacao_id IS NULL) OR (pedido_id IS NULL AND doacao_id IS NOT NULL)
);

COMMENT ON TABLE boicoins_transacoes IS 'Tabela que armazena as transações de Boicoins, como doações e compras.';
COMMENT ON COLUMN boicoins_transacoes.id IS 'Identificador único da transação de Boicoins (UUID fornecido manualmente).';
COMMENT ON COLUMN boicoins_transacoes.usuario_id IS 'Referência ao usuário que realizou a transação.';
COMMENT ON COLUMN boicoins_transacoes.quantidade IS 'Quantidade de Boicoins transacionada.';
COMMENT ON COLUMN boicoins_transacoes.tipo_transacao IS 'Tipo de transação: doação, compra, troca, etc.';
COMMENT ON COLUMN boicoins_transacoes.descricao IS 'Descrição da transação, detalhando sua natureza.';
COMMENT ON COLUMN boicoins_transacoes.data_transacao IS 'Data e hora em que a transação foi realizada.';
COMMENT ON COLUMN boicoins_transacoes.pedido_id IS 'Referência ao pedido relacionado, caso a transação esteja ligada a uma compra.';
COMMENT ON COLUMN boicoins_transacoes.doacao_id IS 'Referência à doação relacionada, caso a transação esteja ligada a uma doação.';
COMMENT ON COLUMN boicoins_transacoes.ponto_mapa_id IS 'Referência ao ponto no mapa associado à transação.';

-- Tabela de Oficinas (Eventos)
CREATE TABLE oficinas (
    id UUID PRIMARY KEY,
    nome VARCHAR(100),
    descricao TEXT,
    preco_boicoins DECIMAL(10, 2),
    preco_real DECIMAL(10, 2),
    data_evento TIMESTAMP,
    limite_participantes INT,
    ponto_mapa_id UUID REFERENCES pontos_mapa(id),
    criado_em TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE oficinas IS 'Tabela que armazena as oficinas e eventos oferecidos pela comunidade.';
COMMENT ON COLUMN oficinas.id IS 'Identificador único da oficina (UUID fornecido manualmente).';
COMMENT ON COLUMN oficinas.nome IS 'Nome da oficina ou evento.';
COMMENT ON COLUMN oficinas.descricao IS 'Descrição detalhada da oficina.';
COMMENT ON COLUMN oficinas.preco_boicoins IS 'Custo da oficina em Boicoins.';
COMMENT ON COLUMN oficinas.preco_real IS 'Custo da oficina em dinheiro real (R$).';
COMMENT ON COLUMN oficinas.data_evento IS 'Data e hora em que a oficina será realizada.';
COMMENT ON COLUMN oficinas.limite_participantes IS 'Número máximo de participantes permitidos na oficina.';
COMMENT ON COLUMN oficinas.ponto_mapa_id IS 'Referência ao local da oficina no mapa.';
COMMENT ON COLUMN oficinas.criado_em IS 'Data e hora de criação do registro da oficina.';

-- Enum para Ações Administrativas
CREATE TYPE acao_adm_enum AS ENUM (
    'validacao_doacao',
    'validacao_pagamento',
    'validacao_inscricao_oficina',
    'outro'
);

-- Tabela de Validações de Administradores (para ações diversas)
CREATE TABLE validacoes_adm (
    id UUID PRIMARY KEY,
    administrador_id UUID REFERENCES usuarios(id),
    acao acao_adm_enum,
    descricao TEXT,
    data_validacao TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE validacoes_adm IS 'Tabela que armazena as ações de validação realizadas por administradores.';
COMMENT ON COLUMN validacoes_adm.id IS 'Identificador único da validação (UUID fornecido manualmente).';
COMMENT ON COLUMN validacoes_adm.administrador_id IS 'Referência ao administrador que realizou a validação.';
COMMENT ON COLUMN validacoes_adm.acao IS 'Ação que foi validada (ex: confirmação de doação, pagamento, etc.).';
COMMENT ON COLUMN validacoes_adm.descricao IS 'Descrição detalhada da ação realizada pelo administrador.';
COMMENT ON COLUMN validacoes_adm.data_validacao IS 'Data e hora em que a validação foi realizada.';

-- Índices sugeridos para melhorar desempenho
CREATE INDEX idx_usuarios_email ON usuarios(email);
CREATE INDEX idx_pedidos_usuario_id ON pedidos(usuario_id);
CREATE INDEX idx_transacoes_usuario_id ON boicoins_transacoes(usuario_id);
CREATE INDEX idx_oficinas_data_evento ON oficinas(data_evento);
CREATE INDEX idx_pontos_mapa_localizacao ON pontos_mapa(latitude, longitude);
