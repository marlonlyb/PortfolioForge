ALTER TABLE products
    ADD COLUMN IF NOT EXISTS source_markdown_url VARCHAR(2048);
