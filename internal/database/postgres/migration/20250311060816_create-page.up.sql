CREATE TABLE pages (
                       id SERIAL PRIMARY KEY,
                       title VARCHAR(255) NOT NULL UNIQUE, -- "about_us", "privacy_policy", "support"
                       content TEXT NOT NULL,
                       updated_at TIMESTAMP DEFAULT NOW()
);

-- Create a trigger to automatically update `updated_at` on row updates
CREATE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_pages_updated_at
    BEFORE UPDATE ON pages
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default pages
INSERT INTO pages (title, content) VALUES
                                       ('about_us', 'Welcome to our platform! We provide amazing services to help you travel with ease.'),
                                       ('privacy_policy', 'We respect your privacy and ensure your data is protected. Read more about our privacy practices here.'),
                                       ('help_support', 'For any issues or inquiries, please contact our support team at support@example.com.');
