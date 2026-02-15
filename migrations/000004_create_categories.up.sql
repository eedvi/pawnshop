-- Create categories table
CREATE TABLE categories (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    parent_id       BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    is_active       BOOLEAN NOT NULL DEFAULT true,

    -- Timestamps
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);
CREATE INDEX idx_categories_slug ON categories(slug);

CREATE TRIGGER categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Insert default categories
INSERT INTO categories (name, slug, description) VALUES
('Electrónicos', 'electronicos', 'Dispositivos electrónicos'),
('Joyería', 'joyeria', 'Joyas y accesorios de valor'),
('Electrodomésticos', 'electrodomesticos', 'Aparatos para el hogar'),
('Herramientas', 'herramientas', 'Herramientas manuales y eléctricas'),
('Vehículos', 'vehiculos', 'Motocicletas, bicicletas y partes'),
('Otros', 'otros', 'Artículos diversos');

-- Subcategories for electronics
INSERT INTO categories (name, slug, description, parent_id)
SELECT 'Celulares', 'celulares', 'Teléfonos móviles', id FROM categories WHERE slug = 'electronicos';
INSERT INTO categories (name, slug, description, parent_id)
SELECT 'Laptops', 'laptops', 'Computadoras portátiles', id FROM categories WHERE slug = 'electronicos';
INSERT INTO categories (name, slug, description, parent_id)
SELECT 'Tablets', 'tablets', 'Tabletas electrónicas', id FROM categories WHERE slug = 'electronicos';
INSERT INTO categories (name, slug, description, parent_id)
SELECT 'Televisores', 'televisores', 'Pantallas y televisores', id FROM categories WHERE slug = 'electronicos';

-- Subcategories for jewelry
INSERT INTO categories (name, slug, description, parent_id)
SELECT 'Oro', 'oro', 'Artículos de oro', id FROM categories WHERE slug = 'joyeria';
INSERT INTO categories (name, slug, description, parent_id)
SELECT 'Plata', 'plata', 'Artículos de plata', id FROM categories WHERE slug = 'joyeria';
INSERT INTO categories (name, slug, description, parent_id)
SELECT 'Relojes', 'relojes', 'Relojes de valor', id FROM categories WHERE slug = 'joyeria';
