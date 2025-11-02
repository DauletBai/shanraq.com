CREATE TABLE real_estate_agencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    tagline TEXT,
    country_code CHAR(2) NOT NULL,
    website TEXT,
    logo_url TEXT,
    head_office TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE realtors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agency_id UUID NOT NULL REFERENCES real_estate_agencies(id) ON DELETE CASCADE,
    full_name TEXT NOT NULL,
    email TEXT UNIQUE,
    phone TEXT,
    languages TEXT[] NOT NULL DEFAULT '{}',
    region TEXT,
    photo_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE property_listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    agency_id UUID REFERENCES real_estate_agencies(id) ON DELETE SET NULL,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    summary TEXT,
    listing_type TEXT NOT NULL,
    country_code CHAR(2) NOT NULL,
    city TEXT,
    region TEXT,
    neighborhood TEXT,
    price NUMERIC(20,2) DEFAULT 0,
    currency CHAR(3) DEFAULT 'USD',
    bedrooms INTEGER,
    bathrooms NUMERIC(4,1),
    area_sqm NUMERIC(12,2),
    hero_image_url TEXT,
    details_url TEXT,
    tags TEXT[] NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_property_listings_country ON property_listings(country_code);
CREATE INDEX idx_property_listings_city ON property_listings(city);
CREATE INDEX idx_property_listings_listing_type ON property_listings(listing_type);
