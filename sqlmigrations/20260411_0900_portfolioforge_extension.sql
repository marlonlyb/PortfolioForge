CREATE TABLE IF NOT EXISTS technologies (
    id UUID NOT NULL,
    name VARCHAR(80) NOT NULL,
    slug VARCHAR(96) NOT NULL,
    category VARCHAR(48) NOT NULL,
    icon VARCHAR(120),
    color VARCHAR(32),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT technologies_id_pk PRIMARY KEY (id),
    CONSTRAINT technologies_slug_uk UNIQUE (slug)
);

CREATE TABLE IF NOT EXISTS project_technologies (
    project_id UUID NOT NULL,
    technology_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT project_technologies_pk PRIMARY KEY (project_id, technology_id),
    CONSTRAINT project_technologies_project_fk FOREIGN KEY (project_id)
        REFERENCES products (id) ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT project_technologies_technology_fk FOREIGN KEY (technology_id)
        REFERENCES technologies (id) ON UPDATE RESTRICT ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_profiles (
    project_id UUID NOT NULL,
    business_goal TEXT,
    problem_statement TEXT,
    solution_summary TEXT,
    architecture TEXT,
    integrations JSONB NOT NULL DEFAULT '[]'::jsonb,
    ai_usage TEXT,
    technical_decisions JSONB NOT NULL DEFAULT '[]'::jsonb,
    challenges JSONB NOT NULL DEFAULT '[]'::jsonb,
    results JSONB NOT NULL DEFAULT '[]'::jsonb,
    metrics JSONB NOT NULL DEFAULT '{}'::jsonb,
    timeline JSONB NOT NULL DEFAULT '[]'::jsonb,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT project_profiles_project_pk PRIMARY KEY (project_id),
    CONSTRAINT project_profiles_project_fk FOREIGN KEY (project_id)
        REFERENCES products (id) ON UPDATE RESTRICT ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS project_media (
    id UUID NOT NULL,
    project_id UUID NOT NULL,
    media_type VARCHAR(24) NOT NULL,
    url TEXT NOT NULL,
    caption TEXT,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT project_media_id_pk PRIMARY KEY (id),
    CONSTRAINT project_media_project_fk FOREIGN KEY (project_id)
        REFERENCES products (id) ON UPDATE RESTRICT ON DELETE CASCADE,
    CONSTRAINT project_media_type_ck CHECK (media_type IN ('image', 'video', 'diagram', 'document'))
);

CREATE INDEX IF NOT EXISTS ix_project_media_project_id ON project_media (project_id, sort_order ASC);

CREATE TABLE IF NOT EXISTS contact_leads (
    id UUID NOT NULL,
    name VARCHAR(120) NOT NULL,
    email VARCHAR(160) NOT NULL,
    company VARCHAR(160),
    project_interest VARCHAR(160),
    message TEXT NOT NULL,
    status VARCHAR(24) NOT NULL DEFAULT 'new',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT contact_leads_id_pk PRIMARY KEY (id),
    CONSTRAINT contact_leads_status_ck CHECK (status IN ('new', 'reviewed', 'contacted', 'closed'))
);

CREATE INDEX IF NOT EXISTS ix_contact_leads_status_created_at ON contact_leads (status, created_at DESC);
