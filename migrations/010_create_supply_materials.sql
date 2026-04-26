CREATE TABLE supply_materials (
    supply_id UUID NOT NULL,
    material_id UUID NOT NULL,
    count INTEGER NOT NULL CHECK (count > 0),

    PRIMARY KEY (supply_id, material_id),

    CONSTRAINT fk_supply_materials_supply
        FOREIGN KEY (supply_id)
        REFERENCES supplies(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_supply_materials_material
        FOREIGN KEY (material_id)
        REFERENCES materials(id)
        ON DELETE CASCADE
);