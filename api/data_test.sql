drop table if exists transaction_categories_to_budget_definitions;
drop table if exists transactions;
drop table if exists transaction_cursors;
drop table if exists transaction_categories;
drop table if exists budget_definitions;
drop table if exists budgets;
drop table if exists plaid_items;
drop table if exists "sessions";
drop table if exists users cascade;

create table if not exists users (
    id uuid primary key default gen_random_uuid(),
    user_name text unique not null,
    password_hash varchar(255) not null 
);

create table if not exists sessions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    token_hash varchar(255) unique not null
);

create table if not exists plaid_items (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    item_id text not null,
    access_token text not null
);

create table if not exists budgets (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    "name" text not null,
    color text
);

create table if not exists budget_definitions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    budget_id uuid not null references budgets(id) on delete cascade,
    "name" text not null,
    allocation int not null
);

create table if not exists transaction_cursors (
    id uuid primary key default gen_random_uuid(),
    user_id uuid unique not null references users(id) on delete cascade,
    "cursor" varchar(256) not null
);

create table if not exists transaction_categories (
    id uuid primary key default gen_random_uuid(),
    "primary" text not null,
    detailed text unique not null,
    "description" text unique not null
);

create table if not exists transactions (
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id) on delete cascade,
    item_id uuid not null references plaid_items(id) on delete cascade,
    transaction_category_detailed text references transaction_categories(detailed) on delete cascade,
    -- fields pulled out from "data"
    account_id text unique not null,
    transaction_id text unique not null,
    amount numeric not null,
    "date" date not null,

    data jsonb not null
    -- account_id text unique not null,
    -- amount numeric not null,
    -- iso_currency_code text,
    -- unofficial_currency_code text,
    -- check_number text,
    -- "date" date not null,
    -- "datetime" timestamp,
    -- "location" jsonb not null,
    -- "name" text,
    -- merchant_name text,
    -- original_description text,
    -- payment_meta jsonb not null,
    -- pending boolean not null,
    -- pending_transaction_id text,
    -- account_owner text,
    -- transaction_id text unique not null,
    -- logo_url text,
    -- website text,
    -- authorized_date date,
    -- authorized_datetime timestamp,
    -- payment_channel text,
    -- transaction_code text,
    -- personal_finance_category_icon_url text,
    -- counterparties jsonb,
    -- merchant_entity_id text
);

create table if not exists transaction_categories_to_budget_definitions (
    id uuid primary key default gen_random_uuid(),
    definition_id uuid not null references budget_definitions(id) on delete cascade,
    category_id uuid not null references transaction_categories(id) on delete cascade
);

-- Seed categories based off of csv found here https://plaid.com/docs/api/products/transactions/#categoriesget
delete from transaction_categories;
insert into transaction_categories (id, "primary", detailed, description) values
    (default,'INCOME','INCOME_DIVIDENDS','Dividends from investment accounts'),
    (default,'INCOME','INCOME_INTEREST_EARNED','Income from interest on savings accounts'),
    (default,'INCOME','INCOME_RETIREMENT_PENSION','Income from pension payments'),
    (default,'INCOME','INCOME_TAX_REFUND','Income from tax refunds'),
    (default,'INCOME','INCOME_UNEMPLOYMENT','Income from unemployment benefits, including unemployment insurance and healthcare'),
    (default,'INCOME','INCOME_WAGES','Income from salaries, gig-economy work, and tips earned'),
    (default,'INCOME','INCOME_OTHER_INCOME','Other miscellaneous income, including alimony, social security, child support, and rental'),
    (default,'TRANSFER_IN','TRANSFER_IN_CASH_ADVANCES_AND_LOANS','Loans and cash advances deposited into a bank account'),
    (default,'TRANSFER_IN','TRANSFER_IN_DEPOSIT','Cash, checks, and ATM deposits into a bank account'),
    (default,'TRANSFER_IN','TRANSFER_IN_INVESTMENT_AND_RETIREMENT_FUNDS','Inbound transfers to an investment or retirement account'),
    (default,'TRANSFER_IN','TRANSFER_IN_SAVINGS','Inbound transfers to a savings account'),
    (default,'TRANSFER_IN','TRANSFER_IN_ACCOUNT_TRANSFER','General inbound transfers from another account'),
    (default,'TRANSFER_IN','TRANSFER_IN_OTHER_TRANSFER_IN','Other miscellaneous inbound transactions'),
    (default,'TRANSFER_OUT','TRANSFER_OUT_INVESTMENT_AND_RETIREMENT_FUNDS','Transfers to an investment or retirement account, including investment apps such as Acorns, Betterment'),
    (default,'TRANSFER_OUT','TRANSFER_OUT_SAVINGS','Outbound transfers to savings accounts'),
    (default,'TRANSFER_OUT','TRANSFER_OUT_WITHDRAWAL','Withdrawals from a bank account'),
    (default,'TRANSFER_OUT','TRANSFER_OUT_ACCOUNT_TRANSFER','General outbound transfers to another account'),
    (default,'TRANSFER_OUT','TRANSFER_OUT_OTHER_TRANSFER_OUT','Other miscellaneous outbound transactions'),
    (default,'LOAN_PAYMENTS','LOAN_PAYMENTS_CAR_PAYMENT','Car loans and leases'),
    (default,'LOAN_PAYMENTS','LOAN_PAYMENTS_CREDIT_CARD_PAYMENT','Payments to a credit card. These are positive amounts for credit card subtypes and negative for depository subtypes'),
    (default,'LOAN_PAYMENTS','LOAN_PAYMENTS_PERSONAL_LOAN_PAYMENT','Personal loans, including cash advances and buy now pay later repayments'),
    (default,'LOAN_PAYMENTS','LOAN_PAYMENTS_MORTGAGE_PAYMENT','Payments on mortgages'),
    (default,'LOAN_PAYMENTS','LOAN_PAYMENTS_STUDENT_LOAN_PAYMENT','Payments on student loans. For college tuition, refer to "General Services - Education"'),
    (default,'LOAN_PAYMENTS','LOAN_PAYMENTS_OTHER_PAYMENT','Other miscellaneous debt payments'),
    (default,'BANK_FEES','BANK_FEES_ATM_FEES','Fees incurred for out-of-network ATMs'),
    (default,'BANK_FEES','BANK_FEES_FOREIGN_TRANSACTION_FEES','Fees incurred on non-domestic transactions'),
    (default,'BANK_FEES','BANK_FEES_INSUFFICIENT_FUNDS','Fees relating to insufficient funds'),
    (default,'BANK_FEES','BANK_FEES_INTEREST_CHARGE','Fees incurred for interest on purchases, including not-paid-in-full or interest on cash advances'),
    (default,'BANK_FEES','BANK_FEES_OVERDRAFT_FEES','Fees incurred when an account is in overdraft'),
    (default,'BANK_FEES','BANK_FEES_OTHER_BANK_FEES','Other miscellaneous bank fees'),
    (default,'ENTERTAINMENT','ENTERTAINMENT_CASINOS_AND_GAMBLING','Gambling, casinos, and sports betting'),
    (default,'ENTERTAINMENT','ENTERTAINMENT_MUSIC_AND_AUDIO','Digital and in-person music purchases, including music streaming services'),
    (default,'ENTERTAINMENT','ENTERTAINMENT_SPORTING_EVENTS_AMUSEMENT_PARKS_AND_MUSEUMS','Purchases made at sporting events, music venues, concerts, museums, and amusement parks'),
    (default,'ENTERTAINMENT','ENTERTAINMENT_TV_AND_MOVIES','In home movie streaming services and movie theaters'),
    (default,'ENTERTAINMENT','ENTERTAINMENT_VIDEO_GAMES','Digital and in-person video game purchases'),
    (default,'ENTERTAINMENT','ENTERTAINMENT_OTHER_ENTERTAINMENT','Other miscellaneous entertainment purchases, including night life and adult entertainment'),
    (default,'FOOD_AND_DRINK','FOOD_AND_DRINK_BEER_WINE_AND_LIQUOR','Beer, Wine & Liquor Stores'),
    (default,'FOOD_AND_DRINK','FOOD_AND_DRINK_COFFEE','Purchases at coffee shops or cafes'),
    (default,'FOOD_AND_DRINK','FOOD_AND_DRINK_FAST_FOOD','Dining expenses for fast food chains'),
    (default,'FOOD_AND_DRINK','FOOD_AND_DRINK_GROCERIES','Purchases for fresh produce and groceries, including farmers'' markets'),
    (default,'FOOD_AND_DRINK','FOOD_AND_DRINK_RESTAURANT','Dining expenses for restaurants, bars, gastropubs, and diners'),
    (default,'FOOD_AND_DRINK','FOOD_AND_DRINK_VENDING_MACHINES','Purchases made at vending machine operators'),
    (default,'FOOD_AND_DRINK','FOOD_AND_DRINK_OTHER_FOOD_AND_DRINK','Other miscellaneous food and drink, including desserts, juice bars, and delis'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_BOOKSTORES_AND_NEWSSTANDS','Books, magazines, and news '),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_CLOTHING_AND_ACCESSORIES','Apparel, shoes, and jewelry'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_CONVENIENCE_STORES','Purchases at convenience stores'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_DEPARTMENT_STORES','Retail stores with wide ranges of consumer goods, typically specializing in clothing and home goods'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_DISCOUNT_STORES','Stores selling goods at a discounted price'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_ELECTRONICS','Electronics stores and websites'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_GIFTS_AND_NOVELTIES','Photo, gifts, cards, and floral stores'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_OFFICE_SUPPLIES','Stores that specialize in office goods'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_ONLINE_MARKETPLACES','Multi-purpose e-commerce platforms such as Etsy, Ebay and Amazon'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_PET_SUPPLIES','Pet supplies and pet food'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_SPORTING_GOODS','Sporting goods, camping gear, and outdoor equipment'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_SUPERSTORES','Superstores such as Target and Walmart, selling both groceries and general merchandise'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_TOBACCO_AND_VAPE','Purchases for tobacco and vaping products'),
    (default,'GENERAL_MERCHANDISE','GENERAL_MERCHANDISE_OTHER_GENERAL_MERCHANDISE','Other miscellaneous merchandise, including toys, hobbies, and arts and crafts'),
    (default,'HOME_IMPROVEMENT','HOME_IMPROVEMENT_FURNITURE','Furniture, bedding, and home accessories'),
    (default,'HOME_IMPROVEMENT','HOME_IMPROVEMENT_HARDWARE','Building materials, hardware stores, paint, and wallpaper'),
    (default,'HOME_IMPROVEMENT','HOME_IMPROVEMENT_REPAIR_AND_MAINTENANCE','Plumbing, lighting, gardening, and roofing'),
    (default,'HOME_IMPROVEMENT','HOME_IMPROVEMENT_SECURITY','Home security system purchases'),
    (default,'HOME_IMPROVEMENT','HOME_IMPROVEMENT_OTHER_HOME_IMPROVEMENT','Other miscellaneous home purchases, including pool installation and pest control'),
    (default,'MEDICAL','MEDICAL_DENTAL_CARE','Dentists and general dental care'),
    (default,'MEDICAL','MEDICAL_EYE_CARE','Optometrists, contacts, and glasses stores'),
    (default,'MEDICAL','MEDICAL_NURSING_CARE','Nursing care and facilities'),
    (default,'MEDICAL','MEDICAL_PHARMACIES_AND_SUPPLEMENTS','Pharmacies and nutrition shops'),
    (default,'MEDICAL','MEDICAL_PRIMARY_CARE','Doctors and physicians'),
    (default,'MEDICAL','MEDICAL_VETERINARY_SERVICES','Prevention and care procedures for animals'),
    (default,'MEDICAL','MEDICAL_OTHER_MEDICAL','Other miscellaneous medical, including blood work, hospitals, and ambulances'),
    (default,'PERSONAL_CARE','PERSONAL_CARE_GYMS_AND_FITNESS_CENTERS','Gyms, fitness centers, and workout classes'),
    (default,'PERSONAL_CARE','PERSONAL_CARE_HAIR_AND_BEAUTY','Manicures, haircuts, waxing, spa/massages, and bath and beauty products '),
    (default,'PERSONAL_CARE','PERSONAL_CARE_LAUNDRY_AND_DRY_CLEANING','Wash and fold, and dry cleaning expenses'),
    (default,'PERSONAL_CARE','PERSONAL_CARE_OTHER_PERSONAL_CARE','Other miscellaneous personal care, including mental health apps and services'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_ACCOUNTING_AND_FINANCIAL_PLANNING','Financial planning, and tax and accounting services'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_AUTOMOTIVE','Oil changes, car washes, repairs, and towing'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_CHILDCARE','Babysitters and daycare'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_CONSULTING_AND_LEGAL','Consulting and legal services'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_EDUCATION','Elementary, high school, professional schools, and college tuition'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_INSURANCE','Insurance for auto, home, and healthcare'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_POSTAGE_AND_SHIPPING','Mail, packaging, and shipping services'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_STORAGE','Storage services and facilities'),
    (default,'GENERAL_SERVICES','GENERAL_SERVICES_OTHER_GENERAL_SERVICES','Other miscellaneous services, including advertising and cloud storage '),
    (default,'GOVERNMENT_AND_NON_PROFIT','GOVERNMENT_AND_NON_PROFIT_DONATIONS','Charitable, political, and religious donations'),
    (default,'GOVERNMENT_AND_NON_PROFIT','GOVERNMENT_AND_NON_PROFIT_GOVERNMENT_DEPARTMENTS_AND_AGENCIES','Government departments and agencies, such as driving licences, and passport renewal'),
    (default,'GOVERNMENT_AND_NON_PROFIT','GOVERNMENT_AND_NON_PROFIT_TAX_PAYMENT','Tax payments, including income and property taxes'),
    (default,'GOVERNMENT_AND_NON_PROFIT','GOVERNMENT_AND_NON_PROFIT_OTHER_GOVERNMENT_AND_NON_PROFIT','Other miscellaneous government and non-profit agencies'),
    (default,'TRANSPORTATION','TRANSPORTATION_BIKES_AND_SCOOTERS','Bike and scooter rentals'),
    (default,'TRANSPORTATION','TRANSPORTATION_GAS','Purchases at a gas station'),
    (default,'TRANSPORTATION','TRANSPORTATION_PARKING','Parking fees and expenses'),
    (default,'TRANSPORTATION','TRANSPORTATION_PUBLIC_TRANSIT','Public transportation, including rail and train, buses, and metro'),
    (default,'TRANSPORTATION','TRANSPORTATION_TAXIS_AND_RIDE_SHARES','Taxi and ride share services'),
    (default,'TRANSPORTATION','TRANSPORTATION_TOLLS','Toll expenses'),
    (default,'TRANSPORTATION','TRANSPORTATION_OTHER_TRANSPORTATION','Other miscellaneous transportation expenses'),
    (default,'TRAVEL','TRAVEL_FLIGHTS','Airline expenses'),
    (default,'TRAVEL','TRAVEL_LODGING','Hotels, motels, and hosted accommodation such as Airbnb'),
    (default,'TRAVEL','TRAVEL_RENTAL_CARS','Rental cars, charter buses, and trucks'),
    (default,'TRAVEL','TRAVEL_OTHER_TRAVEL','Other miscellaneous travel expenses'),
    (default,'RENT_AND_UTILITIES','RENT_AND_UTILITIES_GAS_AND_ELECTRICITY','Gas and electricity bills'),
    (default,'RENT_AND_UTILITIES','RENT_AND_UTILITIES_INTERNET_AND_CABLE','Internet and cable bills'),
    (default,'RENT_AND_UTILITIES','RENT_AND_UTILITIES_RENT','Rent payment'),
    (default,'RENT_AND_UTILITIES','RENT_AND_UTILITIES_SEWAGE_AND_WASTE_MANAGEMENT','Sewage and garbage disposal bills'),
    (default,'RENT_AND_UTILITIES','RENT_AND_UTILITIES_TELEPHONE','Cell phone bills'),
    (default,'RENT_AND_UTILITIES','RENT_AND_UTILITIES_WATER','Water bills'),
    (default,'RENT_AND_UTILITIES','RENT_AND_UTILITIES_OTHER_UTILITIES','Other miscellaneous utility bills');









































































































