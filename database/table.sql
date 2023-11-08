-- Table to store information about authors
CREATE TABLE author (
    author_id INT PRIMARY KEY,
    author_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
);

INSERT INTO author              
VALUES (1, 'amren', 'amren1254@gmail.com');


-- Table to store blog posts
CREATE TABLE blogpost (
    post_id INT PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    image_path VARCHAR(255),  -- Store the file path of the image
    author_id INT,
    published_date DATE,
    FOREIGN KEY (author_id) REFERENCES author(author_id)
);

-- Table to store comments on blog posts
CREATE TABLE comments (
    comment_id INT PRIMARY KEY,
    post_id INT,
    author_id INT,
    comment_text TEXT NOT NULL,
    comment_date DATETIME,
    FOREIGN KEY (post_id) REFERENCES blogpost(post_id),
    FOREIGN KEY (author_id) REFERENCES author(author_id)
);
