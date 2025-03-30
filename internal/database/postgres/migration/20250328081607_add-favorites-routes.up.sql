CREATE TABLE favorite_routes (
                                 user_id VARCHAR(50) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                 route_id VARCHAR(50) NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
                                 created_at TIMESTAMP DEFAULT NOW(),
                                 PRIMARY KEY (user_id, route_id)
);
