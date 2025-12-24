-- Test data for unwise application
-- This script populates the database with sample books and highlights for testing

-- Insert sample books
INSERT INTO books (id, title, author, category, source, cover_image_url, unique_url, updated) VALUES
(1, 'The Pragmatic Programmer', 'Andrew Hunt, David Thomas', 'books', 'kindle', 'https://images-na.ssl-images-amazon.com/images/I/51W1sBPO7tL.jpg', 'https://www.goodreads.com/book/show/4099.The_Pragmatic_Programmer', datetime('now')),
(2, 'Clean Code', 'Robert C. Martin', 'books', 'kindle', 'https://images-na.ssl-images-amazon.com/images/I/41xShlnTZTL.jpg', 'https://www.goodreads.com/book/show/3735293-clean-code', datetime('now')),
(3, 'Design Patterns', 'Erich Gamma, Richard Helm, Ralph Johnson, John Vlissides', 'books', 'kindle', 'https://images-na.ssl-images-amazon.com/images/I/51szD9HC9pL.jpg', 'https://www.goodreads.com/book/show/85009.Design_Patterns', datetime('now')),
(4, 'Refactoring', 'Martin Fowler', 'books', 'kindle', 'https://images-na.ssl-images-amazon.com/images/I/41gOt6FfxPL.jpg', 'https://www.goodreads.com/book/show/44936.Refactoring', datetime('now'));

-- Insert sample highlights for The Pragmatic Programmer
INSERT INTO highlights (id, book_id, text, note, chapter, location, url, updated) VALUES
(1, 1, 'Care About Your Craft: Why spend your life developing software unless you care about doing it well?', 'This is the foundation of being a good developer', 'Chapter 1: A Pragmatic Philosophy', 42, 'https://www.goodreads.com/book/show/4099.The_Pragmatic_Programmer', datetime('now')),
(2, 1, 'Think! About Your Work: Turn off the autopilot and take control. Constantly critique and appraise your work.', 'Always be mindful and intentional', 'Chapter 1: A Pragmatic Philosophy', 45, 'https://www.goodreads.com/book/show/4099.The_Pragmatic_Programmer', datetime('now')),
(3, 1, 'Don''t Live with Broken Windows: Fix bad designs, wrong decisions, and poor code when you see them.', 'Technical debt compounds quickly', 'Chapter 2: A Pragmatic Approach', 89, 'https://www.goodreads.com/book/show/4099.The_Pragmatic_Programmer', datetime('now')),
(4, 1, 'You Can''t Write Perfect Software: Accept it as an axiom of life. Embrace it. Celebrate it.', NULL, 'Chapter 4: Pragmatic Paranoia', 156, 'https://www.goodreads.com/book/show/4099.The_Pragmatic_Programmer', datetime('now'));

-- Insert sample highlights for Clean Code
INSERT INTO highlights (id, book_id, text, note, chapter, location, url, updated) VALUES
(5, 2, 'Clean code is code that has been taken care of. Someone has taken the time to keep it simple and orderly.', 'Definition of clean code', 'Chapter 1: Clean Code', 23, 'https://www.goodreads.com/book/show/3735293-clean-code', datetime('now')),
(6, 2, 'The ratio of time spent reading versus writing is well over 10 to 1. We are constantly reading old code as part of the effort to write new code.', 'Reading code is more common than writing', 'Chapter 1: Clean Code', 34, 'https://www.goodreads.com/book/show/3735293-clean-code', datetime('now')),
(7, 2, 'Functions should do one thing. They should do it well. They should do it only.', 'Single Responsibility Principle for functions', 'Chapter 3: Functions', 78, 'https://www.goodreads.com/book/show/3735293-clean-code', datetime('now')),
(8, 2, 'The first rule of functions is that they should be small. The second rule of functions is that they should be smaller than that.', NULL, 'Chapter 3: Functions', 72, 'https://www.goodreads.com/book/show/3735293-clean-code', datetime('now'));

-- Insert sample highlights for Design Patterns
INSERT INTO highlights (id, book_id, text, note, chapter, location, url, updated) VALUES
(9, 3, 'Program to an interface, not an implementation.', 'Key principle of object-oriented design', 'Chapter 1: Introduction', 56, 'https://www.goodreads.com/book/show/85009.Design_Patterns', datetime('now')),
(10, 3, 'Favor object composition over class inheritance.', 'Composition provides more flexibility', 'Chapter 1: Introduction', 58, 'https://www.goodreads.com/book/show/85009.Design_Patterns', datetime('now')),
(11, 3, 'Design patterns are descriptions of communicating objects and classes that are customized to solve a general design problem in a particular context.', NULL, 'Chapter 1: Introduction', 34, 'https://www.goodreads.com/book/show/85009.Design_Patterns', datetime('now')),
(12, 3, 'Encapsulate what varies.', 'Isolate changing parts from stable parts', 'Chapter 1: Introduction', 62, 'https://www.goodreads.com/book/show/85009.Design_Patterns', datetime('now'));

-- Insert sample highlights for Refactoring
INSERT INTO highlights (id, book_id, text, note, chapter, location, url, updated) VALUES
(13, 4, 'Any fool can write code that a computer can understand. Good programmers write code that humans can understand.', 'Code is for humans first, computers second', 'Chapter 1: Refactoring, a First Example', 28, 'https://www.goodreads.com/book/show/44936.Refactoring', datetime('now')),
(14, 4, 'When you find you have to add a feature to a program, and the program''s code is not structured in a convenient way to add the feature, first refactor the program to make it easy to add the feature, then add the feature.', 'Refactor before adding features', 'Chapter 2: Principles in Refactoring', 45, 'https://www.goodreads.com/book/show/44936.Refactoring', datetime('now')),
(15, 4, 'Refactoring is a controlled technique for improving the design of an existing code base.', NULL, 'Chapter 2: Principles in Refactoring', 42, 'https://www.goodreads.com/book/show/44936.Refactoring', datetime('now')),
(16, 4, 'Before you start refactoring, check that you have a solid suite of tests. These tests must be self-checking.', 'Tests are essential for safe refactoring', 'Chapter 4: Building Tests', 89, 'https://www.goodreads.com/book/show/44936.Refactoring', datetime('now'));
