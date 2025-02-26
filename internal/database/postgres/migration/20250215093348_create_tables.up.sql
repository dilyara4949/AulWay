CREATE TABLE users (
                       id VARCHAR(50) PRIMARY KEY,
                       email VARCHAR(255) UNIQUE,
                       phone VARCHAR(20) UNIQUE NOT NULL,
                       password TEXT NOT NULL,
                       first_name VARCHAR(100) NULL,
                       last_name VARCHAR(100) NULL,
                       role VARCHAR(20) NOT NULL CHECK (role IN ('user', 'admin', 'manager')),
                       require_password_reset BOOLEAN DEFAULT TRUE,
                       created_at TIMESTAMP DEFAULT NOW(),
                       updated_at TIMESTAMP DEFAULT NOW(),
                       deleted_at TIMESTAMP NULL
);

-- CREATE TABLE roles (
--                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--                        name VARCHAR(50) UNIQUE NOT NULL
-- );
--
-- CREATE TABLE user_roles (
--                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--                             user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
--                             role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
--                             UNIQUE (user_id, role_id)
-- );

CREATE TABLE buses (
                       id VARCHAR(50) PRIMARY KEY,
                       number VARCHAR(50) UNIQUE NOT NULL,
                       total_seats INT NOT NULL CHECK (total_seats > 0)
);

CREATE TABLE routes (
                        id VARCHAR(50) PRIMARY KEY,
                        departure VARCHAR(100) NOT NULL,
                        destination VARCHAR(100) NOT NULL,
                        start_date TIMESTAMP NOT NULL,
                        end_date TIMESTAMP NOT NULL CHECK (end_date > start_date),
                        available_seats INT NOT NULL CHECK (available_seats >= 0),
                        bus_id VARCHAR(50) NOT NULL REFERENCES buses(id) ON DELETE SET NULL,
                        price INT NOT NULL CHECK (price >= 0),
                        created_at TIMESTAMP DEFAULT NOW(),
                        updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE tickets (
                         id VARCHAR(50) PRIMARY KEY,
                         user_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                         route_id VARCHAR(50) NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
                         price NUMERIC(10,2) NOT NULL CHECK (price >= 0),
                         status VARCHAR(20) NOT NULL CHECK (status IN ('approved', 'cancelled', 'awaiting')),
                         payment_status VARCHAR(20) NOT NULL CHECK (payment_status IN ('pending', 'paid', 'failed')),
                         qr_code TEXT, -- Stores file path or URL of QR code
                         created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE payments (
                          id VARCHAR(50) PRIMARY KEY,
                          user_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          ticket_id VARCHAR(50) UNIQUE NOT NULL REFERENCES tickets(id) ON DELETE CASCADE,
                          amount INT NOT NULL CHECK (amount >= 0),
                          status VARCHAR(50) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'successful', 'failed')),
                          transaction_id VARCHAR(100) UNIQUE NOT NULL,
                          created_at TIMESTAMP DEFAULT NOW(),
                          updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE faq (
                     id VARCHAR(50) PRIMARY KEY,
                     question TEXT NOT NULL,
                     answer TEXT NOT NULL
);

CREATE TABLE support_requests (
                                  id VARCHAR(50) PRIMARY KEY,
                                  user_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                  message TEXT NOT NULL,
                                  status VARCHAR(20) NOT NULL CHECK (status IN ('open', 'resolved')),
                                  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE settings (
                          id VARCHAR(50) PRIMARY KEY,
                          user_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                          language VARCHAR(10) NOT NULL DEFAULT 'en' CHECK (language IN ('en', 'ru', 'kk')),
                          two_factor_enabled BOOLEAN DEFAULT FALSE
);

CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_tickets_user ON tickets(user_id);
CREATE INDEX idx_routes_departure_destination ON routes(departure, destination);
CREATE UNIQUE INDEX idx_payments_transaction ON payments(transaction_id);
