package tenantdb

const tenantInitSQL = `
	-- ENUM Types (per-tenant schema)
	CREATE TYPE IF NOT EXISTS %s.order_status_enum AS ENUM ('pending', 'paid', 'shipped', 'delivered', 'refunded', 'cancelled');
	CREATE TYPE IF NOT EXISTS %s.payment_status_enum AS ENUM ('pending', 'success', 'failed', 'refunded');
	CREATE TYPE IF NOT EXISTS %s.currency_enum AS ENUM ('NGN', 'USD', 'EUR', 'GBP');
	CREATE TYPE IF NOT EXISTS %s.address_type_enum AS ENUM ('billing', 'shipping', 'other');
	CREATE TYPE IF NOT EXISTS %s.shipping_status_enum AS ENUM ('pending', 'shipped', 'delivered', 'returned', 'lost');

	-- Categories
	CREATE TABLE IF NOT EXISTS %s.categories (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) NOT NULL,
		description TEXT,
		parent_id UUID REFERENCES %s.categories(id) ON DELETE CASCADE,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Products
	CREATE TABLE IF NOT EXISTS %s.products (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		category_id UUID REFERENCES %s.categories(id) ON DELETE SET NULL,
		name VARCHAR(255) NOT NULL,
		description TEXT,
		price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
		sku VARCHAR(100) UNIQUE,
		is_active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Product Variants
	CREATE TABLE IF NOT EXISTS %s.product_variants (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		product_id UUID NOT NULL REFERENCES %s.products(id) ON DELETE CASCADE,
		name VARCHAR(100) NOT NULL,
		price_adjustment DECIMAL(5,2) DEFAULT 0,
		sku VARCHAR(100) UNIQUE,
		created_at TIMESTAMPTZ DEFAULT now()
	);
	-- Product Images
	CREATE TABLE IF NOT EXISTS %s.product_images (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		product_id UUID NOT NULL REFERENCES %s.products(id) ON DELETE CASCADE,
		variant_id UUID REFERENCES %s.product_variants(id) ON DELETE CASCADE,
		url TEXT NOT NULL,
		alt_text VARCHAR(255),
		is_primary BOOLEAN DEFAULT FALSE,
		order_index INT DEFAULT 0,
		created_at TIMESTAMPTZ DEFAULT now()
	);
	-- Inventory Items
	CREATE TABLE IF NOT EXISTS %s.inventory_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		product_id UUID NOT NULL REFERENCES %s.products(id) ON DELETE CASCADE,
		variant_id UUID REFERENCES %s.product_variants(id) ON DELETE SET NULL,
		location VARCHAR(100) DEFAULT 'default',
		quantity INTEGER DEFAULT 0 CHECK (quantity >= 0),
		reserved_quantity INTEGER DEFAULT 0 CHECK (quantity >= reserved_quantity),
		low_stock_threshold INTEGER DEFAULT 5,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Customers
	CREATE TABLE IF NOT EXISTS %s.customers (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		phone VARCHAR(50),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Customer Addresses
	CREATE TABLE IF NOT EXISTS %s.customer_addresses (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		customer_id UUID NOT NULL REFERENCES %s.customers(id) ON DELETE CASCADE,
		type %s.address_type_enum NOT NULL DEFAULT 'shipping',
		street VARCHAR(255) NOT NULL,
		city VARCHAR(100) NOT NULL,
		state VARCHAR(100),
		zip_code VARCHAR(20),
		country VARCHAR(100) DEFAULT 'Nigeria',
		is_default BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMPTZ DEFAULT now()
	);
	ALTER TABLE %s.customer_addresses ADD CONSTRAINT IF NOT EXISTS unique_default_address_per_type UNIQUE (customer_id, type) WHERE is_default = TRUE;
	-- Orders
	CREATE TABLE IF NOT EXISTS %s.orders (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		customer_id UUID REFERENCES %s.customers(id) ON DELETE SET NULL,
		status %s.order_status_enum NOT NULL DEFAULT 'pending',
		total_amount DECIMAL(10,2) NOT NULL,
		subtotal DECIMAL(10,2) NOT NULL,
		tax_amount DECIMAL(10,2) DEFAULT 0,
		shipping_amount DECIMAL(10,2) DEFAULT 0,
		currency %s.currency_enum DEFAULT 'NGN',
		payment_method VARCHAR(50),
		shipping_address_id UUID REFERENCES %s.customer_addresses(id),
		notes TEXT,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Order Items
	CREATE TABLE IF NOT EXISTS %s.order_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id UUID NOT NULL REFERENCES %s.orders(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES %s.products(id) ON DELETE RESTRICT,
		variant_id UUID REFERENCES %s.product_variants(id) ON DELETE SET NULL,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		unit_price DECIMAL(10,2) NOT NULL,
		total_price DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMPTZ DEFAULT now()
	);
	ALTER TABLE %s.order_items ADD COLUMN IF NOT EXISTS total_price_computed DECIMAL(10,2) GENERATED ALWAYS AS (unit_price * quantity) STORED;
	-- Order Shipments
	CREATE TABLE IF NOT EXISTS %s.order_shipments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id UUID NOT NULL REFERENCES %s.orders(id) ON DELETE CASCADE,
		carrier VARCHAR(100),
		tracking_number VARCHAR(100),
		status %s.shipping_status_enum DEFAULT 'pending',
		shipped_at TIMESTAMPTZ,
		delivered_at TIMESTAMPTZ,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Payments
	CREATE TABLE IF NOT EXISTS %s.payments (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		order_id UUID NOT NULL REFERENCES %s.orders(id) ON DELETE CASCADE,
		amount DECIMAL(10,2) NOT NULL,
		currency %s.currency_enum DEFAULT 'NGN',
		provider VARCHAR(50) NOT NULL,
		provider_payment_id VARCHAR(255),
		status %s.payment_status_enum NOT NULL DEFAULT 'pending',
		method VARCHAR(50),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Carts
	CREATE TABLE IF NOT EXISTS %s.carts (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		customer_id UUID REFERENCES %s.customers(id) ON DELETE CASCADE,
		expires_at TIMESTAMPTZ DEFAULT (now() + INTERVAL '30 days'),
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	CREATE TABLE IF NOT EXISTS %s.cart_items (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		cart_id UUID NOT NULL REFERENCES %s.carts(id) ON DELETE CASCADE,
		product_id UUID NOT NULL REFERENCES %s.products(id) ON DELETE RESTRICT,
		variant_id UUID REFERENCES %s.product_variants(id) ON DELETE SET NULL,
		quantity INTEGER NOT NULL CHECK (quantity > 0),
		created_at TIMESTAMPTZ DEFAULT now()
	);
	-- Settings
	CREATE TABLE IF NOT EXISTS %s.settings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		key VARCHAR(100) UNIQUE NOT NULL,
		value TEXT NOT NULL,
		description TEXT,
		created_at TIMESTAMPTZ DEFAULT now(),
		updated_at TIMESTAMPTZ DEFAULT now()
	);
	-- Indexes
	CREATE INDEX IF NOT EXISTS idx_products_category_id ON %s.products(category_id);
	CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON %s.inventory_items(product_id);
	CREATE INDEX IF NOT EXISTS idx_orders_customer_id ON %s.orders(customer_id);
	CREATE INDEX IF NOT EXISTS idx_orders_status ON %s.orders(status);
	CREATE INDEX IF NOT EXISTS idx_customers_email ON %s.customers(email);
	CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON %s.order_items(order_id);
	CREATE INDEX IF NOT EXISTS idx_payments_order_id_status ON %s.payments(order_id, status);
`