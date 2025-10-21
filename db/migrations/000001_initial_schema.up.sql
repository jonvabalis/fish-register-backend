CREATE TABLE catches (
  uuid CHAR(36),
  nickname VARCHAR(512),
  length DOUBLE PRECISION,
  weight DOUBLE PRECISION,
  comment VARCHAR(512),
  caught_at TIMESTAMP NOT NULL,
  created_at TIMESTAMP NOT NULL,
  species_uuid CHAR(36),
  locations_uuid CHAR(36),
  users_uuid CHAR(36) NOT NULL,
  rods_uuid CHAR(36),
  PRIMARY KEY(uuid)
);

CREATE TABLE users (
  uuid CHAR(36),
  username VARCHAR(512) NOT NULL,
  email VARCHAR(512) NOT NULL,
  password VARCHAR(512) NOT NULL,
  PRIMARY KEY(uuid)
);

CREATE TABLE locations (
  uuid CHAR(36),
  name VARCHAR(512) NOT NULL,
  address VARCHAR(512) NOT NULL,
  type VARCHAR(512) NOT NULL,
  PRIMARY KEY(uuid)
);

CREATE TABLE rods (
  uuid CHAR(36),
  nickname VARCHAR(512) NOT NULL,
  brand VARCHAR(512),
  purchase_place VARCHAR(512),
  users_uuid CHAR(36),
  PRIMARY KEY(uuid)
);

CREATE TABLE species (
  uuid CHAR(36),
  name VARCHAR(512) NOT NULL,
  description VARCHAR(512),
  PRIMARY KEY(uuid)
);

CREATE TABLE locations_species (
  locations_uuid CHAR(36),
  species_uuid CHAR(36),
  PRIMARY KEY(locations_uuid, species_uuid)
);

ALTER TABLE catches ADD CONSTRAINT catches_fk1 FOREIGN KEY (species_uuid) REFERENCES species(uuid);
ALTER TABLE catches ADD CONSTRAINT catches_fk2 FOREIGN KEY (locations_uuid) REFERENCES locations(uuid);
ALTER TABLE catches ADD CONSTRAINT catches_fk3 FOREIGN KEY (users_uuid) REFERENCES users(uuid) ON DELETE CASCADE;
ALTER TABLE catches ADD CONSTRAINT catches_fk4 FOREIGN KEY (rods_uuid) REFERENCES rods(uuid);
ALTER TABLE users ADD UNIQUE (username);
ALTER TABLE users ADD UNIQUE (email);
ALTER TABLE rods ADD CONSTRAINT rods_fk1 FOREIGN KEY (users_uuid) REFERENCES users(uuid) ON DELETE CASCADE;
ALTER TABLE locations_species ADD CONSTRAINT locations_species_fk1 FOREIGN KEY (locations_uuid) REFERENCES locations(uuid) ON DELETE CASCADE;
ALTER TABLE locations_species ADD CONSTRAINT locations_species_fk2 FOREIGN KEY (species_uuid) REFERENCES species(uuid) ON DELETE CASCADE;