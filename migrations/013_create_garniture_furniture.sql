CREATE TABLE garniture_furniture (
    garniture_id UUID NOT NULL,
    furniture_id UUID NOT NULL,
    count INTEGER NOT NULL CHECK (count > 0),

    PRIMARY KEY (garniture_id, furniture_id),

    CONSTRAINT fk_garniture_furniture_garniture
        FOREIGN KEY (garniture_id)
        REFERENCES garnitures(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_garniture_furniture_furniture
        FOREIGN KEY (furniture_id)
        REFERENCES furniture(id)
        ON DELETE CASCADE
);