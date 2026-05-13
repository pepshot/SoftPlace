CREATE TABLE module_materials (
    module_id   UUID NOT NULL,
    material_id UUID NOT NULL,
    count       INTEGER NOT NULL CHECK (count > 0),

    PRIMARY KEY (module_id, material_id),

    CONSTRAINT fk_module_materials_module
        FOREIGN KEY (module_id)
        REFERENCES modules(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_module_materials_material
        FOREIGN KEY (material_id)
        REFERENCES materials(id)
        ON DELETE CASCADE
);