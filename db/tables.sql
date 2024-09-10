CREATE TABLE "Promotion" (
    id SERIAL PRIMARY KEY,
    price VARCHAR(255),
    description VARCHAR(255),
    image_url VARCHAR(255)
)

INSERT INTO "Promotion" (price, description, image_url) VALUES ('10', '10% off', '/app/frontend/static/images/Screenshot (6).png')