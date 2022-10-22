CREATE TABLE IF NOT EXISTS public.proxy_group (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(20)

);

CREATE TABLE IF NOT EXISTS public.proxy(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ip          VARCHAR(15) NOT NULL,
    port        INT NOT NULL,
    external_ip VARCHAR(15) NOT NULL,
    country VARCHAR(20) NOT NULL,
    open_ports INT[] DEFAULT ARRAY []::INTEGER[],
    active BOOLEAN DEFAULT FALSE,
    ping INT DEFAULT -1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    checked_at TIMESTAMPTZ DEFAULT date('epoch'),
    valid_at TIMESTAMPTZ DEFAULT date('epoch'),
    bl_check INT DEFAULT 0,
    processing_status INT DEFAULT 0,
    proxy_group_id UUID NOT NULL,
    CONSTRAINT group_fk FOREIGN KEY(proxy_group_id) REFERENCES public.proxy_group(id)
);