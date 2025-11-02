DELETE FROM property_listings WHERE slug IN (
    'palm-jumeirah-sky-villa',
    'ostermalm-art-nouveau',
    'kyoto-machiya-boutique-hotel',
    'lisbon-digital-district-loft',
    'tuscany-heritage-vineyard-estate',
    'singapore-sky-garden-duplex',
    'reykjavik-geothermal-retreat',
    'cape-town-atlantic-seaboard-villa',
    'sao-paulo-innovation-hub-loft',
    'british-columbia-wilderness-lodge'
);

DELETE FROM realtors WHERE email IN (
    'layla@shanraq.com',
    'karl@nordicskyline.com',
    'maya@pacificaurban.com',
    'giulia@atlasheritage.it',
    'diego@pacificaurban.com'
);

DELETE FROM real_estate_agencies WHERE slug IN (
    'shanraq-global-realty',
    'nordic-skyline-partners',
    'pacifica-urban-advisors',
    'atlas-heritage-homes'
);
