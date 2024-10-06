-- Criação do ENUM para o tipo de transação Boicoins
CREATE TYPE tipo_transacao_enum AS ENUM (
    'recebimento_doacao',
    'compra_produto',
    'inscricao_oficina',
    'bonus_participacao'
);

ALTER TABLE boicoins_transacoes
ALTER COLUMN tipo_transacao TYPE tipo_transacao_enum
USING tipo_transacao::tipo_transacao_enum;

CREATE OR REPLACE FUNCTION atualizar_saldo_boicoins() RETURNS TRIGGER AS $$
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
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_atualizar_saldo_boicoins
AFTER INSERT OR UPDATE ON boicoins_transacoes
FOR EACH ROW
EXECUTE FUNCTION atualizar_saldo_boicoins();

COMMENT ON TYPE tipo_transacao_enum IS 'Enum para definir os diferentes tipos de transações de Boicoins.';
COMMENT ON FUNCTION atualizar_saldo_boicoins IS 'Função que atualiza o saldo de Boicoins de um usuário após cada transação, garantindo consistência e verificando saldo insuficiente para débitos.';
