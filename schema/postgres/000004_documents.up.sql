CREATE TABLE documents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    mime VARCHAR(64) NOT NULL,
    file_path VARCHAR(255),
    is_file BOOLEAN NOT NULL,
    is_public BOOLEAN NOT NULL,
    document_data TEXT,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);