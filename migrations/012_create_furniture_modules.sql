CREATE TABLE furniture_modules (
    furniture_id    UUID NOT NULL,
    module_id       UUID NOT NULL,
    count           INTEGER NOT NULL CHECK (count > 0),

    PRIMARY KEY (furniture_id, module_id),

    CONSTRAINT fk_furniture_modules_furniture
        FOREIGN KEY (furniture_id)
        REFERENCES furniture(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_furniture_modules_module
        FOREIGN KEY (module_id)
        REFERENCES modules(id)
        ON DELETE CASCADE
);