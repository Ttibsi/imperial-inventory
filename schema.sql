CREATE TABLE status (
	id INT AUTO_INCREMENT PRIMARY KEY,
	value VARCHAR(255) NOT NULL
);

INSERT INTO status VALUES
    (NULL, "Operational"),
    (NULL, "Under Repair"),
    (NULL, "Destroyed");

CREATE TABLE spacecraft (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    class VARCHAR(255) NOT NULL,
    crew INT NOT NULL,
    image VARCHAR(255) NOT NULL,
    value DECIMAL(10, 2) NOT NULL,
    status INT,
    FOREIGN KEY (status) REFERENCES status(id) ON DELETE CASCADE
);

CREATE TABLE armament (
    id INT AUTO_INCREMENT PRIMARY KEY,
    ship_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    qty INT NOT NULL,
    FOREIGN KEY (ship_id) REFERENCES spacecraft(id) ON DELETE CASCADE
);

INSERT INTO spacecraft VALUES
    (NULL, "Devastator", "Star Destroyer", 35000, "https://static.wikia.nocookie.net/starwars/images/1/11/Executor_and_escorts.jpg/revision/latest/scale-to-width-down/1000?cb=20120105172952", 1999.99, 1),
    (NULL, "Luke's Ship", "X-wing", 2, "https://static.wikia.nocookie.net/starwars/images/0/00/Xwing-ROOCE.png/revision/latest/scale-to-width-down/1000?cb=20230516042654", 250.00, 1);

INSERT INTO armament VALUES
    (NULL, 1, "Turbo Laser", 60),
    (NULL, 1, "Ion Cannons", 60),
    (NULL, 1, "Tractor Beam", 10),
    (NULL, 2, "Front Blasters", 4);
