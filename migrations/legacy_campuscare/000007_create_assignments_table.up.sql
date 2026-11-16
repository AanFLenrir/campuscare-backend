-- Up migration for assignments table
CREATE TABLE assignments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    complaint_id UUID NOT NULL REFERENCES complaints(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    officer_id UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    assigned_by UUID NOT NULL REFERENCES users(id) ON UPDATE CASCADE ON DELETE RESTRICT,
    assigned_at TIMESTAMP NOT NULL DEFAULT now(),
    notes TEXT
);
CREATE INDEX idx_assignments_officer ON assignments(officer_id);

