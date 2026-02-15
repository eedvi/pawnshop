-- =====================================================
-- SEED DATA FOR PAWNSHOP
-- Datos de prueba realistas para desarrollo/testing
-- =====================================================

-- =====================================================
-- CATEGORÍAS ADICIONALES (algunas ya existen en migraciones)
-- =====================================================
INSERT INTO categories (name, slug, description, parent_id, is_active) VALUES
('Instrumentos Musicales', 'instrumentos-musicales', 'Instrumentos de música', NULL, true),
('Artículos Deportivos', 'deportivos', 'Equipos y accesorios deportivos', NULL, true)
ON CONFLICT (slug) DO NOTHING;

-- Subcategorías adicionales para Electrónicos
INSERT INTO categories (name, slug, description, parent_id, is_active)
SELECT 'Consolas', 'consolas', 'Consolas de videojuegos', id, true
FROM categories WHERE slug = 'electronicos'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description, parent_id, is_active)
SELECT 'Cámaras', 'camaras', 'Cámaras fotográficas y de video', id, true
FROM categories WHERE slug = 'electronicos'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description, parent_id, is_active)
SELECT 'Audio', 'audio', 'Equipos de audio y bocinas', id, true
FROM categories WHERE slug = 'electronicos'
ON CONFLICT (slug) DO NOTHING;

-- Subcategorías de Herramientas
INSERT INTO categories (name, slug, description, parent_id, is_active)
SELECT 'Eléctricas', 'electricas', 'Herramientas eléctricas', id, true
FROM categories WHERE slug = 'herramientas'
ON CONFLICT (slug) DO NOTHING;

INSERT INTO categories (name, slug, description, parent_id, is_active)
SELECT 'Manuales', 'manuales', 'Herramientas manuales', id, true
FROM categories WHERE slug = 'herramientas'
ON CONFLICT (slug) DO NOTHING;

-- =====================================================
-- CLIENTES
-- =====================================================
INSERT INTO customers (
    branch_id, first_name, last_name, identity_type, identity_number,
    phone, email, address, city, state, occupation, monthly_income,
    credit_limit, notes, is_active
) VALUES
(1, 'Carlos', 'López Hernández', 'dpi', '2345678901234',
 '5555-1234', 'carlos.lopez@gmail.com', '4a Calle 5-67 Zona 1', 'Guatemala', 'Guatemala',
 'Comerciante', 8000.00, 15000.00, 'Cliente frecuente, buen historial de pago', true),

(1, 'María', 'González Pérez', 'dpi', '1234567890123',
 '5555-5678', 'maria.gonzalez@hotmail.com', '12 Avenida 3-45 Zona 10', 'Guatemala', 'Guatemala',
 'Contadora', 12000.00, 20000.00, NULL, true),

(1, 'José', 'Martínez Rodríguez', 'dpi', '3456789012345',
 '5555-9012', NULL, '6a Calle 8-90 Zona 7', 'Mixco', 'Guatemala',
 'Mecánico', 6000.00, 10000.00, 'Prefiere ser contactado por teléfono', true),

(1, 'Ana', 'García Morales', 'dpi', '4567890123456',
 '5555-3456', 'ana.garcia@yahoo.com', '2a Avenida 1-23 Zona 5', 'Villa Nueva', 'Guatemala',
 'Maestra', 5500.00, 8000.00, NULL, true),

(1, 'Pedro', 'Ramírez Flores', 'dpi', '5678901234567',
 '5555-7890', NULL, '9a Calle 4-56 Zona 11', 'Guatemala', 'Guatemala',
 'Albañil', 4500.00, 7000.00, 'Trabaja en construcción', true),

(1, 'Luisa', 'Hernández Castro', 'dpi', '6789012345678',
 '5555-2345', 'luisa.h@gmail.com', '3a Avenida 7-89 Zona 2', 'Guatemala', 'Guatemala',
 'Enfermera', 7000.00, 12000.00, NULL, true),

(1, 'Roberto', 'Díaz Mendoza', 'dpi', '7890123456789',
 '5555-6789', NULL, '8a Calle 2-34 Zona 12', 'Villa Canales', 'Guatemala',
 'Empresario', 15000.00, 30000.00, 'Tiene negocio de electrodomésticos', true),

(1, 'Carmen', 'Moreno Vásquez', 'passport', 'C12345678',
 '5555-0123', 'carmen.moreno@outlook.com', '5a Avenida 9-01 Zona 9', 'Guatemala', 'Guatemala',
 'Consultora', 20000.00, 35000.00, 'Extranjera residente', true),

(1, 'Miguel', 'Torres Ruiz', 'dpi', '8901234567890',
 '5555-4567', NULL, '1a Calle 6-78 Zona 6', 'Chinautla', 'Guatemala',
 'Electricista', 5000.00, 8000.00, NULL, true),

(1, 'Elena', 'Sánchez López', 'dpi', '9012345678901',
 '5555-8901', 'elena.sanchez@gmail.com', '7a Avenida 0-12 Zona 4', 'Guatemala', 'Guatemala',
 'Secretaria', 4500.00, 7500.00, 'Trabaja en oficina gubernamental', true);

-- =====================================================
-- ARTÍCULOS
-- =====================================================

-- Función auxiliar para insertar artículos
DO $$
DECLARE
    cat_celulares BIGINT;
    cat_laptops BIGINT;
    cat_consolas BIGINT;
    cat_camaras BIGINT;
    cat_audio BIGINT;
    cat_oro BIGINT;
    cat_plata BIGINT;
    cat_relojes BIGINT;
    cat_electricas BIGINT;
    cat_instrumentos BIGINT;
    cat_electrodomesticos BIGINT;
    cat_vehiculos BIGINT;
BEGIN
    SELECT id INTO cat_celulares FROM categories WHERE slug = 'celulares';
    SELECT id INTO cat_laptops FROM categories WHERE slug = 'laptops';
    SELECT id INTO cat_consolas FROM categories WHERE slug = 'consolas';
    SELECT id INTO cat_camaras FROM categories WHERE slug = 'camaras';
    SELECT id INTO cat_audio FROM categories WHERE slug = 'audio';
    SELECT id INTO cat_oro FROM categories WHERE slug = 'oro';
    SELECT id INTO cat_plata FROM categories WHERE slug = 'plata';
    SELECT id INTO cat_relojes FROM categories WHERE slug = 'relojes';
    SELECT id INTO cat_electricas FROM categories WHERE slug = 'electricas';
    SELECT id INTO cat_instrumentos FROM categories WHERE slug = 'instrumentos-musicales';
    SELECT id INTO cat_electrodomesticos FROM categories WHERE slug = 'electrodomesticos';
    SELECT id INTO cat_vehiculos FROM categories WHERE slug = 'vehiculos';

    -- Insertar artículos (status: available, collateral, sold, confiscated)
    INSERT INTO items (
        branch_id, customer_id, category_id, sku, name, description, brand, model,
        serial_number, condition, appraised_value, loan_value, sale_price, status, created_by
    ) VALUES
    -- Celulares
    (1, 1, cat_celulares, 'CEL-001', 'iPhone 13 Pro Max 256GB', 'Color grafito, incluye cargador original',
     'Apple', 'iPhone 13 Pro Max', 'DNPXYZ123456', 'good', 6500.00, 4500.00, 5500.00, 'collateral', 1),

    (1, 2, cat_celulares, 'CEL-002', 'Samsung Galaxy S22 Ultra', 'Negro fantasma, 256GB, como nuevo',
     'Samsung', 'Galaxy S22 Ultra', 'R5CNXYZ789012', 'excellent', 5500.00, 3800.00, 4800.00, 'collateral', 1),

    (1, 3, cat_celulares, 'CEL-003', 'iPhone 12 128GB', 'Azul, pequeño rayón en esquina',
     'Apple', 'iPhone 12', 'DNPABC456789', 'fair', 3500.00, 2400.00, 3000.00, 'available', 1),

    -- Laptops
    (1, 4, cat_laptops, 'LAP-001', 'MacBook Pro 14" M1 Pro', '16GB RAM, 512GB SSD, gris espacial',
     'Apple', 'MacBook Pro 14', 'C02XYZ123ABC', 'excellent', 12000.00, 8500.00, 10500.00, 'collateral', 1),

    (1, 5, cat_laptops, 'LAP-002', 'Dell XPS 15', 'Intel i7, 16GB RAM, 512GB SSD, pantalla OLED',
     'Dell', 'XPS 15 9510', 'DELLSVC456789', 'good', 8500.00, 6000.00, 7500.00, 'collateral', 1),

    (1, 1, cat_laptops, 'LAP-003', 'HP Pavilion Gaming', 'Ryzen 5, 8GB RAM, GTX 1650',
     'HP', 'Pavilion 15-ec', 'HPGAMING123456', 'good', 4500.00, 3000.00, 3800.00, 'available', 1),

    -- Consolas
    (1, 6, cat_consolas, 'CON-001', 'PlayStation 5', 'Edición disco, incluye control extra',
     'Sony', 'PlayStation 5', 'PS5ABC123456', 'excellent', 4000.00, 2800.00, 3500.00, 'collateral', 1),

    (1, 7, cat_consolas, 'CON-002', 'Nintendo Switch OLED', 'Blanca, con estuche y 2 juegos',
     'Nintendo', 'Switch OLED', 'NSOLED789012', 'good', 2500.00, 1800.00, 2200.00, 'sold', 1),

    -- Joyería - Oro
    (1, 8, cat_oro, 'ORO-001', 'Cadena de oro 18k', '45cm, eslabón cubano, 25 gramos',
     NULL, NULL, NULL, 'excellent', 8500.00, 6000.00, 7500.00, 'collateral', 1),

    (1, 9, cat_oro, 'ORO-002', 'Anillo de compromiso oro 14k', 'Con diamante 0.5 quilates',
     NULL, NULL, NULL, 'excellent', 4500.00, 3200.00, 4000.00, 'collateral', 1),

    (1, 2, cat_oro, 'ORO-003', 'Pulsera oro 18k', '18cm, eslabón cartier, 15 gramos',
     NULL, NULL, NULL, 'good', 5200.00, 3600.00, 4500.00, 'available', 1),

    -- Joyería - Relojes
    (1, 10, cat_relojes, 'REL-001', 'Rolex Submariner', 'Acero inoxidable, año 2019, con caja y papeles',
     'Rolex', 'Submariner 116610LN', 'ROL789456123', 'excellent', 75000.00, 50000.00, 65000.00, 'collateral', 1),

    (1, 3, cat_relojes, 'REL-002', 'Tag Heuer Carrera', 'Automático, acero, correa de cuero',
     'Tag Heuer', 'Carrera CV2014', 'TAG456789012', 'good', 12000.00, 8000.00, 10000.00, 'available', 1),

    -- Herramientas
    (1, 5, cat_electricas, 'HER-001', 'Rotomartillo DeWalt', '20V MAX, incluye 2 baterías y maletín',
     'DeWalt', 'DCD996P2', 'DEW123456789', 'good', 2800.00, 2000.00, 2500.00, 'collateral', 1),

    (1, 7, cat_electricas, 'HER-002', 'Sierra de mesa Makita', '10 pulgadas, 15 amp',
     'Makita', '2705', 'MAK789012345', 'fair', 3500.00, 2400.00, 3000.00, 'available', 1),

    -- Instrumentos
    (1, 4, cat_instrumentos, 'MUS-001', 'Guitarra Gibson Les Paul', 'Standard, color heritage cherry sunburst',
     'Gibson', 'Les Paul Standard', 'GIB2021LP456', 'excellent', 18000.00, 12000.00, 15000.00, 'collateral', 1),

    (1, 6, cat_instrumentos, 'MUS-002', 'Teclado Yamaha PSR-E373', '61 teclas, como nuevo',
     'Yamaha', 'PSR-E373', 'YAMKEY123456', 'excellent', 2200.00, 1500.00, 1900.00, 'sold', 1),

    -- Audio
    (1, 8, cat_audio, 'AUD-001', 'Bocinas JBL PartyBox 310', 'Bluetooth, luces LED',
     'JBL', 'PartyBox 310', 'JBLPB310XYZ', 'good', 3200.00, 2200.00, 2800.00, 'collateral', 1),

    -- Electrodomésticos
    (1, 9, cat_electrodomesticos, 'ELE-001', 'Refrigeradora Samsung', '18 pies cúbicos, French door',
     'Samsung', 'RF18A5101', 'SAMREF789012', 'good', 8500.00, 5500.00, 7000.00, 'available', 1),

    -- Vehículos
    (1, 10, cat_vehiculos, 'VEH-001', 'Motocicleta Honda CB190R', 'Año 2022, 5000 km, negra',
     'Honda', 'CB190R', 'HONDACB2022GT', 'excellent', 22000.00, 15000.00, 19000.00, 'collateral', 1);

END $$;

-- =====================================================
-- PRÉSTAMOS
-- =====================================================

-- Préstamo 1: Activo, al día (iPhone 13 Pro Max)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000001', 1, 1, i.id,
    4500.00, 10.00, 450.00,
    4500.00, 450.00, 4950.00, 0.00,
    CURRENT_DATE - INTERVAL '15 days', CURRENT_DATE + INTERVAL '15 days', 30, 5,
    'minimum_payment', true, 450.00,
    'active', 1
FROM items i WHERE i.sku = 'CEL-001';

-- Préstamo 2: Activo (Samsung Galaxy)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000002', 1, 2, i.id,
    3800.00, 10.00, 380.00,
    3800.00, 380.00, 4180.00, 0.00,
    CURRENT_DATE - INTERVAL '10 days', CURRENT_DATE + INTERVAL '20 days', 30, 5,
    'minimum_payment', true, 380.00,
    'active', 1
FROM items i WHERE i.sku = 'CEL-002';

-- Préstamo 3: Vencido (Dell XPS)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, days_overdue, created_by
)
SELECT
    'LN-2026-000003', 1, 5, i.id,
    6000.00, 10.00, 600.00,
    6000.00, 600.00, 6600.00, 0.00,
    CURRENT_DATE - INTERVAL '45 days', CURRENT_DATE - INTERVAL '15 days', 30, 5,
    'minimum_payment', true, 600.00,
    'overdue', 15, 1
FROM items i WHERE i.sku = 'LAP-002';

-- Préstamo 4: Activo (MacBook Pro)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000004', 1, 4, i.id,
    8500.00, 10.00, 850.00,
    8500.00, 850.00, 9350.00, 0.00,
    CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE + INTERVAL '25 days', 30, 5,
    'minimum_payment', true, 850.00,
    'active', 1
FROM items i WHERE i.sku = 'LAP-001';

-- Préstamo 5: Pagado (PlayStation 5)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, paid_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000005', 1, 6, i.id,
    2800.00, 10.00, 280.00,
    0.00, 0.00, 3080.00, 3080.00,
    CURRENT_DATE - INTERVAL '35 days', CURRENT_DATE - INTERVAL '5 days', CURRENT_DATE - INTERVAL '7 days', 30, 5,
    'minimum_payment', true, 280.00,
    'paid', 1
FROM items i WHERE i.sku = 'CON-001';

-- Actualizar PS5 a disponible (fue redimido)
UPDATE items SET status = 'available' WHERE sku = 'CON-001';

-- Préstamo 6: Activo (Cadena de oro)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000006', 1, 8, i.id,
    6000.00, 10.00, 600.00,
    6000.00, 600.00, 6600.00, 0.00,
    CURRENT_DATE - INTERVAL '20 days', CURRENT_DATE + INTERVAL '10 days', 30, 5,
    'minimum_payment', true, 600.00,
    'active', 1
FROM items i WHERE i.sku = 'ORO-001';

-- Préstamo 7: Vencido (Anillo de compromiso)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, days_overdue, created_by
)
SELECT
    'LN-2026-000007', 1, 9, i.id,
    3200.00, 10.00, 320.00,
    3200.00, 320.00, 3520.00, 0.00,
    CURRENT_DATE - INTERVAL '60 days', CURRENT_DATE - INTERVAL '30 days', 30, 5,
    'minimum_payment', true, 320.00,
    'overdue', 30, 1
FROM items i WHERE i.sku = 'ORO-002';

-- Préstamo 8: Activo (Rolex)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000008', 1, 10, i.id,
    50000.00, 8.00, 4000.00,
    50000.00, 4000.00, 54000.00, 0.00,
    CURRENT_DATE - INTERVAL '30 days', CURRENT_DATE + INTERVAL '30 days', 60, 10,
    'minimum_payment', true, 2000.00,
    'active', 1
FROM items i WHERE i.sku = 'REL-001';

-- Préstamo 9: Activo (Rotomartillo)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000009', 1, 5, i.id,
    2000.00, 10.00, 200.00,
    2000.00, 200.00, 2200.00, 0.00,
    CURRENT_DATE - INTERVAL '8 days', CURRENT_DATE + INTERVAL '22 days', 30, 5,
    'minimum_payment', true, 200.00,
    'active', 1
FROM items i WHERE i.sku = 'HER-001';

-- Préstamo 10: Activo (Gibson Les Paul)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000010', 1, 4, i.id,
    12000.00, 10.00, 1200.00,
    12000.00, 1200.00, 13200.00, 0.00,
    CURRENT_DATE - INTERVAL '12 days', CURRENT_DATE + INTERVAL '33 days', 45, 5,
    'minimum_payment', true, 880.00,
    'active', 1
FROM items i WHERE i.sku = 'MUS-001';

-- Préstamo 11: Activo (Motocicleta Honda)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000011', 1, 10, i.id,
    15000.00, 8.00, 1200.00,
    15000.00, 1200.00, 16200.00, 0.00,
    CURRENT_DATE - INTERVAL '25 days', CURRENT_DATE + INTERVAL '35 days', 60, 10,
    'minimum_payment', true, 810.00,
    'active', 1
FROM items i WHERE i.sku = 'VEH-001';

-- Préstamo 12: Activo (Bocinas JBL)
INSERT INTO loans (
    loan_number, branch_id, customer_id, item_id,
    loan_amount, interest_rate, interest_amount,
    principal_remaining, interest_remaining, total_amount, amount_paid,
    start_date, due_date, loan_term_days, grace_period_days,
    payment_plan_type, requires_minimum_payment, minimum_payment_amount,
    status, created_by
)
SELECT
    'LN-2026-000012', 1, 8, i.id,
    2200.00, 10.00, 220.00,
    2200.00, 220.00, 2420.00, 0.00,
    CURRENT_DATE - INTERVAL '3 days', CURRENT_DATE + INTERVAL '27 days', 30, 5,
    'minimum_payment', true, 220.00,
    'active', 1
FROM items i WHERE i.sku = 'AUD-001';

-- =====================================================
-- PAGOS
-- =====================================================

-- Pagos para préstamo L-2026-000005 (pagado completamente)
INSERT INTO payments (
    payment_number, branch_id, loan_id, customer_id,
    amount, principal_amount, interest_amount, late_fee_amount,
    payment_method, status, payment_date,
    loan_balance_after, interest_balance_after,
    notes, created_by
)
SELECT
    'PY-2026-000001', 1, l.id, l.customer_id,
    1500.00, 1220.00, 280.00, 0.00,
    'cash', 'completed', CURRENT_DATE - INTERVAL '20 days',
    1580.00, 0.00,
    'Primer abono', 1
FROM loans l WHERE l.loan_number = 'LN-2026-000005';

INSERT INTO payments (
    payment_number, branch_id, loan_id, customer_id,
    amount, principal_amount, interest_amount, late_fee_amount,
    payment_method, status, payment_date,
    loan_balance_after, interest_balance_after,
    notes, created_by
)
SELECT
    'PY-2026-000002', 1, l.id, l.customer_id,
    1580.00, 1580.00, 0.00, 0.00,
    'cash', 'completed', CURRENT_DATE - INTERVAL '7 days',
    0.00, 0.00,
    'Pago final - redención de artículo', 1
FROM loans l WHERE l.loan_number = 'LN-2026-000005';

-- Pago parcial para préstamo L-2026-000001
INSERT INTO payments (
    payment_number, branch_id, loan_id, customer_id,
    amount, principal_amount, interest_amount, late_fee_amount,
    payment_method, status, payment_date,
    loan_balance_after, interest_balance_after,
    notes, created_by
)
SELECT
    'PY-2026-000003', 1, l.id, l.customer_id,
    1000.00, 550.00, 450.00, 0.00,
    'cash', 'completed', CURRENT_DATE - INTERVAL '5 days',
    3950.00, 0.00,
    'Abono', 1
FROM loans l WHERE l.loan_number = 'LN-2026-000001';

-- Actualizar balance del préstamo
UPDATE loans SET
    principal_remaining = 3950.00,
    interest_remaining = 0.00,
    amount_paid = 1000.00
WHERE loan_number = 'LN-2026-000001';

-- Pago de solo intereses para préstamo vencido L-2026-000003
INSERT INTO payments (
    payment_number, branch_id, loan_id, customer_id,
    amount, principal_amount, interest_amount, late_fee_amount,
    payment_method, status, payment_date,
    loan_balance_after, interest_balance_after,
    notes, created_by
)
SELECT
    'PY-2026-000004', 1, l.id, l.customer_id,
    600.00, 0.00, 600.00, 0.00,
    'cash', 'completed', CURRENT_DATE - INTERVAL '10 days',
    6000.00, 0.00,
    'Pago de intereses - cliente solicita extensión', 1
FROM loans l WHERE l.loan_number = 'LN-2026-000003';

UPDATE loans SET
    interest_remaining = 0.00,
    amount_paid = 600.00
WHERE loan_number = 'LN-2026-000003';

-- =====================================================
-- VENTAS
-- =====================================================

-- Venta 1: Nintendo Switch (ya vendido)
INSERT INTO sales (
    sale_number, branch_id, item_id, customer_id,
    original_price, discount_percent, discount_amount, final_price, tax_amount, total_amount,
    payment_method, status, sale_date, notes, created_by
)
SELECT
    'SL-2026-000001', 1, i.id, 7,
    2200.00, 0.00, 0.00, 2200.00, 0.00, 2200.00,
    'cash', 'completed', CURRENT_DATE - INTERVAL '5 days', 'Venta directa', 1
FROM items i WHERE i.sku = 'CON-002';

-- Venta 2: Teclado Yamaha (ya vendido)
INSERT INTO sales (
    sale_number, branch_id, item_id, customer_id,
    original_price, discount_percent, discount_amount, final_price, tax_amount, total_amount,
    payment_method, status, sale_date, notes, created_by
)
SELECT
    'SL-2026-000002', 1, i.id, 4,
    1900.00, 5.00, 95.00, 1805.00, 0.00, 1805.00,
    'card', 'completed', CURRENT_DATE - INTERVAL '3 days', 'Pago con tarjeta de débito, descuento por pago contado', 1
FROM items i WHERE i.sku = 'MUS-002';

-- =====================================================
-- USUARIOS ADICIONALES
-- =====================================================

-- Crear usuarios adicionales para pruebas
INSERT INTO users (branch_id, role_id, email, password_hash, first_name, last_name, is_active, email_verified)
SELECT
    1,
    (SELECT id FROM roles WHERE name = 'manager'),
    'gerente@pawnshop.com',
    '$argon2id$v=19$m=65536,t=3,p=2$Rx8KhjlnxPmeOK0LY6z0ow$yg/SKFaf0DwZ0Q/xN8Ho33UE2UJGOt2wL2rRGmSeljA',
    'Juan',
    'Pérez',
    true,
    true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'gerente@pawnshop.com');

INSERT INTO users (branch_id, role_id, email, password_hash, first_name, last_name, is_active, email_verified)
SELECT
    1,
    (SELECT id FROM roles WHERE name = 'cashier'),
    'cajero@pawnshop.com',
    '$argon2id$v=19$m=65536,t=3,p=2$Rx8KhjlnxPmeOK0LY6z0ow$yg/SKFaf0DwZ0Q/xN8Ho33UE2UJGOt2wL2rRGmSeljA',
    'María',
    'Rodríguez',
    true,
    true
WHERE NOT EXISTS (SELECT 1 FROM users WHERE email = 'cajero@pawnshop.com');

-- =====================================================
-- RESUMEN DE DATOS CREADOS
-- =====================================================
-- Categorías: ~19 total (6 en migración + subcategorías + nuevas)
-- Clientes: 10
-- Artículos: 20
-- Préstamos: 12 (9 activos, 2 vencidos, 1 pagado)
-- Pagos: 4
-- Ventas: 2
-- Usuarios: 3 (admin, gerente, cajero)
--
-- Credenciales de prueba (todos con password: admin123):
-- - admin@pawnshop.com (Super Admin)
-- - gerente@pawnshop.com (Gerente)
-- - cajero@pawnshop.com (Cajero)
-- =====================================================

SELECT 'Seed completado exitosamente!' AS resultado;
