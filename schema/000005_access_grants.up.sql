CREATE TABLE access_grants (
    document_id UUID REFERENCES documents(id),
    user_id UUID REFERENCES users(id),
    UNIQUE (document_id, user_id)
);