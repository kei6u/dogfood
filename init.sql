CREATE TABLE record
(
	dogfood_name varchar(50) NOT NULL,
	gram INTEGER NOT NULL,
	dog_name varchar(50) NOT NULL,
	eaten_at TIMESTAMP NOT NULL,
	PRIMARY KEY(dogfood_name, dog_name, eaten_at)
);
