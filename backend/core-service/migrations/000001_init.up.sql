-- =========================
-- USERS
-- =========================
CREATE TABLE IF NOT EXISTS users (
                                     id BIGSERIAL PRIMARY KEY,
                                     name TEXT NOT NULL,
                                     email TEXT NOT NULL UNIQUE,
                                     password_hash TEXT NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- =========================
-- SPECIALISTS
-- один user может быть в нескольких категориях
-- =========================
CREATE TABLE IF NOT EXISTS specialists (
                                           user_id BIGINT NOT NULL,
                                           category TEXT NOT NULL,
                                           current_load INT NOT NULL DEFAULT 0,

                                           PRIMARY KEY (user_id, category),

    CONSTRAINT fk_specialists_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

-- =========================
-- TICKETS
-- =========================
CREATE TABLE IF NOT EXISTS tickets (
                                       id BIGSERIAL PRIMARY KEY,
                                       user_id BIGINT NOT NULL,
                                       message TEXT NOT NULL,
                                       category TEXT NOT NULL,
                                       status TEXT NOT NULL,
                                       specialist_id BIGINT NULL,
                                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_tickets_user_id
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,

    CONSTRAINT fk_tickets_specialist_id
    FOREIGN KEY (specialist_id) REFERENCES users(id) ON DELETE SET NULL
    );

-- =========================
-- FIXTURES: USERS
-- (пароль: 123456 для всех)
-- =========================
INSERT INTO users (name, email, password_hash)
VALUES
    ('Иван Петров', 'ivan.petrov@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),
    ('Алексей Морозов', 'alexey.morozov@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),

    ('Мария Волкова', 'maria.volkova@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),
    ('Дмитрий Соколов', 'dmitry.sokolov@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),

    ('Екатерина Иванова', 'ekaterina.ivanova@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),
    ('Олег Кузнецов', 'oleg.kuznetsov@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),

    ('Анна Смирнова', 'anna.smirnova@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),
    ('Сергей Орлов', 'sergey.orlov@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),

    ('Наталья Федорова', 'natalya.fedorova@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),
    ('Павел Новиков', 'pavel.novikov@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),

    ('Татьяна Лебедева', 'tatyana.lebedeva@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e'),
    ('Виктор Андреев', 'viktor.andreev@support.local', '$2a$10$7QJ8eXJw1z9zqX1j0zvY6eXvS9YwK5j1QJ9k1u1xV3Kx9pF2v8m6e')
    ON CONFLICT (email) DO NOTHING;

-- =========================
-- FIXTURES: SPECIALISTS
-- =========================

-- ACCOUNT
INSERT INTO specialists (user_id, category)
SELECT id, 'ACCOUNT' FROM users WHERE email = 'ivan.petrov@support.local'
    ON CONFLICT DO NOTHING;

INSERT INTO specialists (user_id, category)
SELECT id, 'ACCOUNT' FROM users WHERE email = 'alexey.morozov@support.local'
    ON CONFLICT DO NOTHING;

-- ORDER
INSERT INTO specialists (user_id, category)
SELECT id, 'ORDER' FROM users WHERE email = 'maria.volkova@support.local'
    ON CONFLICT DO NOTHING;

INSERT INTO specialists (user_id, category)
SELECT id, 'ORDER' FROM users WHERE email = 'dmitry.sokolov@support.local'
    ON CONFLICT DO NOTHING;

-- REFUND
INSERT INTO specialists (user_id, category)
SELECT id, 'REFUND' FROM users WHERE email = 'ekaterina.ivanova@support.local'
    ON CONFLICT DO NOTHING;

INSERT INTO specialists (user_id, category)
SELECT id, 'REFUND' FROM users WHERE email = 'oleg.kuznetsov@support.local'
    ON CONFLICT DO NOTHING;

-- PAYMENT
INSERT INTO specialists (user_id, category)
SELECT id, 'PAYMENT' FROM users WHERE email = 'anna.smirnova@support.local'
    ON CONFLICT DO NOTHING;

INSERT INTO specialists (user_id, category)
SELECT id, 'PAYMENT' FROM users WHERE email = 'sergey.orlov@support.local'
    ON CONFLICT DO NOTHING;

-- DELIVERY
INSERT INTO specialists (user_id, category)
SELECT id, 'DELIVERY' FROM users WHERE email = 'natalya.fedorova@support.local'
    ON CONFLICT DO NOTHING;

INSERT INTO specialists (user_id, category)
SELECT id, 'DELIVERY' FROM users WHERE email = 'pavel.novikov@support.local'
    ON CONFLICT DO NOTHING;

-- SUPPORT
INSERT INTO specialists (user_id, category)
SELECT id, 'SUPPORT' FROM users WHERE email = 'tatyana.lebedeva@support.local'
    ON CONFLICT DO NOTHING;

INSERT INTO specialists (user_id, category)
SELECT id, 'SUPPORT' FROM users WHERE email = 'viktor.andreev@support.local'
    ON CONFLICT DO NOTHING;