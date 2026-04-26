CREATE TABLE shipment_garnitures (
    shipment_id UUID NOT NULL,
    garniture_id UUID NOT NULL,
    count INTEGER NOT NULL CHECK (count > 0),

    PRIMARY KEY (shipment_id, garniture_id),

    CONSTRAINT fk_shipment_garnitures_shipment
        FOREIGN KEY (shipment_id)
        REFERENCES shipments(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_shipment_garnitures_garniture
        FOREIGN KEY (garniture_id)
        REFERENCES garnitures(id)
        ON DELETE CASCADE
);