CREATE TABLE notifications (
    notification_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    recipient_internal_id BIGINT NOT NULL,
    subject VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    reference_type VARCHAR(50),
    reference_internal_id BIGINT,
    reference_public_id UUID, -- ← ADD THIS
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT notifications_public_id_unique UNIQUE (public_id),
    CONSTRAINT notifications_recipient_fk FOREIGN KEY (recipient_internal_id) REFERENCES users (internal_id) ON DELETE CASCADE
);

CREATE INDEX idx_notifications_recipient ON notifications (recipient_internal_id);

CREATE INDEX idx_notifications_reference ON notifications (
    reference_type,
    reference_internal_id
);

CREATE INDEX idx_notifications_reference_public ON notifications (reference_public_id);