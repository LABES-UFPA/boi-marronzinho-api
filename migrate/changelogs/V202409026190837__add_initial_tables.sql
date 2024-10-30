-- Definição de Enums
CREATE TYPE status_pedido_enum AS ENUM ('em_andamento', 'concluido', 'cancelado', 'pendente_pagamento');
CREATE TYPE tipo_transacao_enum AS ENUM ('recebimento_doacao', 'compra_produto', 'inscricao_oficina', 'bonus_participacao', 'estorno_oficina');
CREATE TYPE tipo_usuario_enum AS ENUM ('Usuario', 'Administrador');
CREATE TYPE status_doacao_enum AS ENUM ('pendente', 'validada', 'rejeitada');

-- Tabelas
CREATE TABLE estados (
                         id uuid NOT NULL,
                         nome varchar(100),
                         CONSTRAINT estados_pkey PRIMARY KEY (id)
);

CREATE TABLE evento (
                        id uuid NOT NULL,
                        nome varchar(100) NOT NULL,
                        descricao text,
                        data_evento timestamp NOT NULL,
                        link_endereco text,
                        criado_em timestamp DEFAULT CURRENT_TIMESTAMP,
                        atualizado_em timestamp DEFAULT CURRENT_TIMESTAMP,
                        imagem_url varchar(255),
                        CONSTRAINT evento_pkey PRIMARY KEY (id)
);

CREATE TABLE item_troca (
                            id uuid NOT NULL,
                            nome_item varchar(100),
                            descricao text,
                            unidade_medida varchar(50),
                            boicoins_por_unidade numeric(10, 2),
                            CONSTRAINT item_troca_pkey PRIMARY KEY (id)
);

CREATE TABLE oficinas (
                          id uuid NOT NULL,
                          nome varchar(100),
                          descricao text,
                          preco_boicoins numeric(10, 2),
                          preco_real numeric(10, 2),
                          data_evento timestamp,
                          limite_participantes integer,
                          participantes_atual integer DEFAULT 0,
                          criado_em timestamp DEFAULT CURRENT_TIMESTAMP,
                          link_endereco varchar(250),
                          imagem_url varchar(255),
                          CONSTRAINT oficinas_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_oficinas_data_evento ON oficinas USING btree (data_evento);

CREATE TABLE pontos_mapa (
                             id uuid NOT NULL,
                             nome varchar(100),
                             descricao text,
                             latitude numeric(10, 6),
                             longitude numeric(10, 6),
                             imagem_url varchar(255),
                             criado_em timestamp DEFAULT CURRENT_TIMESTAMP,
                             CONSTRAINT pontos_mapa_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_pontos_mapa_localizacao ON pontos_mapa USING btree (latitude, longitude);

CREATE TABLE produtos (
                          id uuid NOT NULL,
                          nome varchar(100),
                          descricao text,
                          preco_boicoins numeric(10, 2),
                          preco_real numeric(10, 2),
                          quantidade_em_estoque integer DEFAULT 0,
                          imagem_url varchar(255),
                          criado_em timestamp DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT produtos_pkey PRIMARY KEY (id)
);

CREATE TABLE usuarios (
                          id uuid NOT NULL,
                          email varchar(100),
                          tipo_usuario tipo_usuario_enum DEFAULT 'Usuario',
                          saldo_boicoins numeric(10, 2) DEFAULT 0.00,
                          idioma_preferido varchar DEFAULT 'pt',
                          created_at timestamp DEFAULT CURRENT_TIMESTAMP,
                          password_hash varchar(255) NOT NULL,
                          last_login timestamp,
                          password_reset_token varchar(255),
                          password_reset_expires timestamp,
                          deleted_at timestamp,
                          updated_at timestamp DEFAULT CURRENT_TIMESTAMP,
                          first_name varchar(50) NOT NULL,
                          last_name varchar(50) NOT NULL,
                          CONSTRAINT usuarios_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_usuarios_email ON usuarios USING btree (email);

CREATE TABLE carrinho_itens (
                                id uuid NOT NULL,
                                usuario_id uuid,
                                produto_id uuid,
                                quantidade integer NOT NULL,
                                preco_unitario numeric(10, 2) NOT NULL,
                                criado_em timestamp DEFAULT CURRENT_TIMESTAMP,
                                CONSTRAINT carrinho_itens_pkey PRIMARY KEY (id)
);

ALTER TABLE carrinho_itens ADD CONSTRAINT carrinho_itens_preco_unitario_check CHECK (preco_unitario >= 0);
ALTER TABLE carrinho_itens ADD CONSTRAINT carrinho_itens_quantidade_check CHECK (quantidade > 0);

CREATE TABLE cidades (
                         id uuid NOT NULL,
                         estado_id uuid,
                         nome varchar(100),
                         CONSTRAINT cidades_pkey PRIMARY KEY (id)
);

CREATE TABLE enderecos (
                           id uuid NOT NULL,
                           usuario_id uuid,
                           logradouro varchar(255),
                           numero varchar(10),
                           complemento varchar(100),
                           bairro varchar(100),
                           cidade_id uuid,
                           cep varchar(20),
                           pais varchar(100),
                           tipo_endereco varchar(50),
                           criado_em timestamp DEFAULT CURRENT_TIMESTAMP,
                           CONSTRAINT enderecos_pkey PRIMARY KEY (id)
);

CREATE TABLE participantes_oficinas (
                                        id uuid NOT NULL,
                                        usuario_id uuid,
                                        oficina_id uuid,
                                        data_inscricao timestamp DEFAULT CURRENT_TIMESTAMP,
                                        CONSTRAINT participantes_oficinas_pkey PRIMARY KEY (id)
);

CREATE TABLE pedidos (
                         id uuid NOT NULL,
                         usuario_id uuid,
                         boicoins_usados numeric(10, 2) DEFAULT 0.00,
                         preco_real_usado numeric(10, 2),
                         quantidade integer DEFAULT 1,
                         codigo varchar(100),
                         data_pedido timestamp DEFAULT CURRENT_TIMESTAMP,
                         data_conclusao timestamp,
                         status_pedido status_pedido_enum,
                         CONSTRAINT pedidos_pkey PRIMARY KEY (id)
);

CREATE INDEX idx_pedidos_usuario_id ON pedidos USING btree (usuario_id);

CREATE TABLE ticket_oficina (
                                id uuid NOT NULL,
                                usuario_id uuid,
                                oficina_id uuid,
                                codigo varchar(100),
                                created_at timestamp DEFAULT CURRENT_TIMESTAMP,
                                validado boolean DEFAULT false,
                                qr_code varchar(255),
                                CONSTRAINT ticket_oficina_pkey PRIMARY KEY (id),
                                CONSTRAINT codigo_unico UNIQUE (codigo)
);

CREATE TABLE troca (
                       id uuid NOT NULL,
                       usuario_id uuid,
                       item_troca_id uuid,
                       quantidade numeric(10, 2),
                       boicoins_recebidos numeric(10, 2),
                       data_troca timestamp DEFAULT CURRENT_TIMESTAMP,
                       status status_doacao_enum DEFAULT 'pendente',
                       CONSTRAINT trocas_pkey PRIMARY KEY (id)
);

CREATE TABLE boicoins_transacoes (
                                     id uuid NOT NULL,
                                     usuario_id uuid,
                                     quantidade numeric(10, 2),
                                     tipo_transacao tipo_transacao_enum,
                                     descricao text,
                                     data_transacao timestamp DEFAULT CURRENT_TIMESTAMP,
                                     pedido_id uuid,
                                     troca_id uuid,
                                     oficina_id uuid,
                                     CONSTRAINT boicoins_transacoes_pkey PRIMARY KEY (id)
);

ALTER TABLE boicoins_transacoes ADD CONSTRAINT ck_transacao_unico CHECK (
    ((pedido_id IS NOT NULL AND troca_id IS NULL AND oficina_id IS NULL) OR
     (pedido_id IS NULL AND troca_id IS NOT NULL AND oficina_id IS NULL) OR
     (pedido_id IS NULL AND troca_id IS NULL AND oficina_id IS NOT NULL))
    );

CREATE INDEX idx_transacoes_usuario_id ON boicoins_transacoes USING btree (usuario_id);

CREATE TABLE pedido_itens (
                              id uuid NOT NULL,
                              pedido_id uuid,
                              produto_id uuid,
                              quantidade integer NOT NULL,
                              preco_unitario double precision,
                              CONSTRAINT pedido_itens_pkey PRIMARY KEY (id)
);

-- Chaves estrangeiras
ALTER TABLE boicoins_transacoes ADD CONSTRAINT boicoins_transacoes_troca_id_fkey FOREIGN KEY (troca_id) REFERENCES troca(id);
ALTER TABLE boicoins_transacoes ADD CONSTRAINT boicoins_transacoes_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES usuarios(id);
ALTER TABLE boicoins_transacoes ADD CONSTRAINT boicoins_transacoes_oficina_id_fkey FOREIGN KEY (oficina_id) REFERENCES oficinas(id) ON DELETE SET NULL;
ALTER TABLE boicoins_transacoes ADD CONSTRAINT boicoins_transacoes_pedido_id_fkey FOREIGN KEY (pedido_id) REFERENCES pedidos(id);

ALTER TABLE carrinho_itens ADD CONSTRAINT carrinho_itens_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE;
ALTER TABLE carrinho_itens ADD CONSTRAINT carrinho_itens_produto_id_fkey FOREIGN KEY (produto_id) REFERENCES produtos(id) ON DELETE CASCADE;

ALTER TABLE cidades ADD CONSTRAINT cidades_estado_id_fkey FOREIGN KEY (estado_id) REFERENCES estados(id);

ALTER TABLE enderecos ADD CONSTRAINT enderecos_cidade_id_fkey FOREIGN KEY (cidade_id) REFERENCES cidades(id);
ALTER TABLE enderecos ADD CONSTRAINT enderecos_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

ALTER TABLE participantes_oficinas ADD CONSTRAINT participantes_oficinas_oficina_id_fkey FOREIGN KEY (oficina_id) REFERENCES oficinas(id);
ALTER TABLE participantes_oficinas ADD CONSTRAINT participantes_oficinas_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

ALTER TABLE pedido_itens ADD CONSTRAINT pedido_itens_produto_id_fkey FOREIGN KEY (produto_id) REFERENCES produtos(id);
ALTER TABLE pedido_itens ADD CONSTRAINT pedido_itens_pedido_id_fkey FOREIGN KEY (pedido_id) REFERENCES pedidos(id) ON DELETE CASCADE;

ALTER TABLE pedidos ADD CONSTRAINT pedidos_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

ALTER TABLE ticket_oficina ADD CONSTRAINT ticket_oficina_oficina_id_fkey FOREIGN KEY (oficina_id) REFERENCES oficinas(id);
ALTER TABLE ticket_oficina ADD CONSTRAINT ticket_oficina_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES usuarios(id);

ALTER TABLE troca ADD CONSTRAINT trocas_usuario_id_fkey FOREIGN KEY (usuario_id) REFERENCES usuarios(id);
ALTER TABLE troca ADD CONSTRAINT trocas_item_troca_id_fkey FOREIGN KEY (item_troca_id) REFERENCES item_troca(id);

-- Função: atualizar_data_evento
CREATE OR REPLACE FUNCTION public.atualizar_data_evento()
    RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    NEW.atualizado_em = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

-- Trigger: trigger_atualizar_data_evento
CREATE TRIGGER trigger_atualizar_data_evento
    BEFORE UPDATE ON evento
    FOR EACH ROW
EXECUTE FUNCTION public.atualizar_data_evento();

-- Função: atualizar_saldo_boicoins
CREATE OR REPLACE FUNCTION atualizar_saldo_boicoins()
    RETURNS trigger
    LANGUAGE plpgsql
AS $$
BEGIN
    IF NEW.quantidade >= 0 THEN
        UPDATE usuarios
        SET saldo_boicoins = saldo_boicoins + NEW.quantidade
        WHERE id = NEW.usuario_id;
    ELSE
        IF (SELECT saldo_boicoins FROM usuarios WHERE id = NEW.usuario_id) >= ABS(NEW.quantidade) THEN
            UPDATE usuarios
            SET saldo_boicoins = saldo_boicoins + NEW.quantidade
            WHERE id = NEW.usuario_id;
        ELSE
            RAISE EXCEPTION 'Saldo insuficiente para realizar esta transação.';
        END IF;
    END IF;
    RETURN NEW;
END;
$$;

-- Trigger: trigger_atualizar_saldo_boicoins
CREATE TRIGGER trigger_atualizar_saldo_boicoins
    AFTER INSERT OR UPDATE ON boicoins_transacoes
    FOR EACH ROW
EXECUTE FUNCTION atualizar_saldo_boicoins();

-- Comentários
COMMENT ON TABLE evento IS 'Tabela que armazena eventos do Boi Marrozinho, incluindo detalhes e datas.';
COMMENT ON COLUMN evento.data_evento IS 'Data e hora em que o evento ocorrerá.';
COMMENT ON COLUMN evento.link_endereco IS 'Link de endereço relacionado ao evento.';
COMMENT ON COLUMN evento.criado_em IS 'Data e hora em que o registro do evento foi criado.';
COMMENT ON COLUMN evento.atualizado_em IS 'Data e hora da última atualização do registro do evento.';

COMMENT ON TABLE item_troca IS 'Tabela que lista os tipos de itens aceitos para troca.';
COMMENT ON TABLE oficinas IS 'Tabela que armazena as oficinas e eventos oferecidos pela comunidade.';
COMMENT ON TABLE pontos_mapa IS 'Tabela que armazena os pontos de coleta ou distribuição exibidos no mapa do sistema.';
COMMENT ON TABLE produtos IS 'Tabela que armazena os produtos disponíveis na loja.';
COMMENT ON TABLE usuarios IS 'Tabela que armazena as informações dos usuários, incluindo tipo e saldo de Boicoins.';
COMMENT ON TABLE carrinho_itens IS 'Tabela que armazena os itens que os usuários adicionaram ao carrinho de compras.';
COMMENT ON COLUMN carrinho_itens.preco_unitario IS 'Preço em Boicoins por unidade do produto no momento em que foi adicionado ao carrinho.';
COMMENT ON COLUMN carrinho_itens.criado_em IS 'Data e hora em que o item foi adicionado ao carrinho.';
COMMENT ON TABLE enderecos IS 'Tabela que armazena os endereços dos usuários.';
COMMENT ON TABLE participantes_oficinas IS 'Tabela que registra os participantes inscritos em oficinas.';
COMMENT ON TABLE pedidos IS 'Tabela que armazena os pedidos de produtos e inscrições em oficinas feitos pelos usuários.';
COMMENT ON TABLE ticket_oficina IS 'Tabela que armazena os tickets de inscrição em oficinas.';
COMMENT ON COLUMN ticket_oficina.validado IS 'Indica se o elemento foi validado (TRUE) ou não (FALSE).';
COMMENT ON TABLE troca IS 'Tabela que armazena as trocas realizadas pelos usuários.';
COMMENT ON TABLE boicoins_transacoes IS 'Tabela que armazena as transações de Boicoins, como doações e compras.';
COMMENT ON COLUMN boicoins_transacoes.troca_id IS 'Referência para a troca que gerou a transação de Boicoins.';
COMMENT ON FUNCTION atualizar_saldo_boicoins IS 'Função que atualiza o saldo de Boicoins de um usuário após cada transação, garantindo consistência e verificando saldo insuficiente para débitos.';
