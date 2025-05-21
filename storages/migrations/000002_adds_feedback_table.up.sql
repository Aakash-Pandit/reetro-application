CREATE TABLE feedbacks (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    message VARCHAR(100) NOT NULL,
    board_id VARCHAR(36),
    FOREIGN KEY (board_id) REFERENCES boards(id),
    created_by_id VARCHAR(36),
    FOREIGN KEY (created_by_id) REFERENCES users(id),
    created_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP(3) NOT NULL DEFAULT CURRENT_TIMESTAMP
);